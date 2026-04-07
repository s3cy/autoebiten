package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/s3cy/autoebiten/internal/output"
	"github.com/s3cy/autoebiten/internal/proxy"
	"github.com/s3cy/autoebiten/internal/rpc"
)

// LaunchOptions contains options for the launch command.
type LaunchOptions struct {
	GameCmd  string
	GameArgs []string
	Timeout  time.Duration // Timeout waiting for game RPC server
}

// LaunchCommand handles the `autoebiten launch` functionality.
type LaunchCommand struct {
	options      *LaunchOptions
	outputFiles  *output.FilePath
	gameProc     *os.Process
	proxyServer  *proxy.Server
	proxyHandler *proxy.Handler
	listener     net.Listener
	gameExited   chan struct{}
	done         chan struct{}
}

// NewLaunchCommand creates a new launch command handler.
func NewLaunchCommand(options *LaunchOptions) *LaunchCommand {
	return &LaunchCommand{
		options:    options,
		gameExited: make(chan struct{}),
		done:       make(chan struct{}),
	}
}

// Run executes the launch command.
// This method blocks until the game exits or is terminated.
func (lc *LaunchCommand) Run() error {
	// Create game command with pipes (must be done before Start())
	gameCmd, stdoutPipe, stderrPipe, err := lc.createGameCommand()
	if err != nil {
		return fmt.Errorf("failed to create game command: %w", err)
	}

	// Start the game
	if err := gameCmd.Start(); err != nil {
		return fmt.Errorf("failed to start game: %w", err)
	}

	// Set target PID to the game process and store reference
	rpc.SetTargetPID(gameCmd.Process.Pid)
	lc.gameProc = gameCmd.Process

	// Now derive output file paths using the game PID
	gameSocketPath := rpc.SocketPath()
	lc.outputFiles = output.DerivePaths(gameSocketPath)

	// Create log file
	logFile, err := output.CreateLogFile(lc.outputFiles.Log)
	if err != nil {
		lc.terminateGame()
		return fmt.Errorf("failed to create log file: %w", err)
	}
	// Note: we don't defer close here - it needs to stay open for the tee goroutines

	// Tee stdout/stderr to both terminal and log file
	go lc.teeOutput(stdoutPipe, os.Stdout, logFile)
	go lc.teeOutput(stderrPipe, os.Stderr, logFile)

	// Monitor game exit
	go func() {
		gameCmd.Wait()
		close(lc.gameExited)
	}()

	// Wait for game RPC server to be ready (with timeout)
	gameClient, err := lc.waitForGameRPC()
	if err != nil {
		lc.cleanup()
		lc.terminateGame()
		return fmt.Errorf("failed to connect to game RPC server: %w", err)
	}

	// Create proxy server
	lc.proxyServer = proxy.NewServer(gameClient, lc.outputFiles)
	lc.proxyHandler = proxy.NewHandler(lc.proxyServer)

	// Start proxy RPC server
	if err := lc.startProxyServer(); err != nil {
		lc.cleanup()
		return fmt.Errorf("failed to start proxy server: %w", err)
	}

	// Setup signal handling
	lc.setupSignalHandling()

	// Wait for game to exit or termination signal
	lc.waitForExit()

	return nil
}

// waitForGameRPC polls the game's RPC server until it's ready or timeout.
func (lc *LaunchCommand) waitForGameRPC() (*rpc.Client, error) {
	// Default timeout if not specified
	timeout := lc.options.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	interval := 100 * time.Millisecond
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		// Try to connect
		client, err := rpc.NewClient()
		if err == nil {
			// Try to ping to verify it's really ready
			req, _ := rpc.BuildRequest("ping", nil)
			resp, err := client.SendRequest(req)
			if err == nil && resp.Error == nil {
				return client, nil
			}
			// Ping failed, close and retry
			client.Close()
		}

		// Wait before next attempt
		time.Sleep(interval)
	}

	return nil, fmt.Errorf("timeout after %v waiting for game RPC server", timeout)
}

// createGameCommand creates the game command with pipes set up.
// Returns the command and stdout/stderr pipes (must call Start() after this).
func (lc *LaunchCommand) createGameCommand() (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	cmd := exec.Command(lc.options.GameCmd, lc.options.GameArgs...)

	// Pass through all environment variables
	cmd.Env = os.Environ()

	// Pass through stdin for interactive games
	cmd.Stdin = os.Stdin

	// Create pipes for stdout and stderr (must be done before Start())
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	return cmd, stdoutPipe, stderrPipe, nil
}

// teeOutput copies data from src to both dst1 (terminal) and dst2 (log file).
func (lc *LaunchCommand) teeOutput(src io.Reader, dst1, dst2 *os.File) {
	reader := bufio.NewReader(src)
	for {
		data, err := reader.ReadBytes('\n')
		if len(data) > 0 {
			dst1.Write(data)
			dst2.Write(data)
		}
		if err != nil {
			if err != io.EOF {
				// Read remaining data
				remaining, _ := reader.ReadBytes('\n')
				if len(remaining) > 0 {
					dst1.Write(remaining)
					dst2.Write(remaining)
				}
			}
			break
		}
	}
}

// startProxyServer starts the proxy RPC server on the launch socket.
func (lc *LaunchCommand) startProxyServer() error {
	// Ensure socket directory exists
	dir := filepath.Dir(lc.outputFiles.LaunchSock)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create socket directory: %w", err)
	}

	// Remove existing socket if present
	os.Remove(lc.outputFiles.LaunchSock)

	// Create listener
	listener, err := net.Listen("unix", lc.outputFiles.LaunchSock)
	if err != nil {
		return fmt.Errorf("failed to listen on socket %s: %w", lc.outputFiles.LaunchSock, err)
	}
	lc.listener = listener

	// Set socket permissions
	if err := os.Chmod(lc.outputFiles.LaunchSock, 0777); err != nil {
		listener.Close()
		return fmt.Errorf("failed to set socket permissions: %w", err)
	}

	// Start accept loop in background
	go lc.acceptLoop()

	return nil
}

// acceptLoop accepts incoming RPC connections.
func (lc *LaunchCommand) acceptLoop() {
	for {
		conn, err := lc.listener.Accept()
		if err != nil {
			// Check if listener was closed
			if netErr, ok := err.(*net.OpError); ok && netErr.Err.Error() == "use of closed network connection" {
				return
			}
			fmt.Fprintf(os.Stderr, "autoebiten: accept error: %v\n", err)
			continue
		}

		go lc.handleConnection(conn)
	}
}

// handleConnection handles a single RPC connection.
func (lc *LaunchCommand) handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	for {
		var req rpc.RPCRequest
		if err := decoder.Decode(&req); err != nil {
			if err == io.EOF {
				return
			}
			fmt.Fprintf(os.Stderr, "autoebiten: decode error: %v\n", err)
			return
		}

		// Handle exit specially - it should trigger cleanup
		if req.Method == "exit" {
			lc.proxyHandler.ProcessRequest(conn, &req)
			close(lc.done)
			return
		}

		lc.proxyHandler.ProcessRequest(conn, &req)
	}
}

// setupSignalHandling handles Ctrl+C and other signals.
func (lc *LaunchCommand) setupSignalHandling() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nReceived interrupt signal, terminating game...")
		lc.terminateGame()
		close(lc.done)
	}()
}

// waitForExit waits for the game to exit or termination signal.
func (lc *LaunchCommand) waitForExit() {
	// Wait for game to exit
	<-lc.gameExited

	fmt.Println("Game exited, waiting 30s for CLI to read final output...")

	// Wait for either: done signal (from exit command or interrupt) or 30s timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	select {
	case <-lc.done:
		// Explicit termination requested
	case <-ctx.Done():
		// Timeout reached
	}

	// Cleanup
	lc.cleanup()
}

// terminateGame terminates the game process.
func (lc *LaunchCommand) terminateGame() {
	if lc.gameProc != nil {
		lc.gameProc.Signal(syscall.SIGTERM)
		// Give it a moment to terminate gracefully
		time.Sleep(100 * time.Millisecond)
		lc.gameProc.Kill()
	}
}

// cleanup removes all temporary files.
func (lc *LaunchCommand) cleanup() {
	// Close proxy server
	if lc.proxyServer != nil {
		lc.proxyServer.Close()
	}

	// Close listener
	if lc.listener != nil {
		lc.listener.Close()
	}

	// Remove launch socket
	os.Remove(lc.outputFiles.LaunchSock)

	// Remove log and snapshot files
	lc.proxyServer.Cleanup()
}
