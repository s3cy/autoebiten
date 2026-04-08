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
	"sync"
	"syscall"
	"time"

	"github.com/s3cy/autoebiten/internal/output"
	"github.com/s3cy/autoebiten/internal/proxy"
	"github.com/s3cy/autoebiten/internal/rpc"
)

var WaitForExitTimeout = 30 * time.Second

// LaunchOptions contains options for the launch command.
type LaunchOptions struct {
	GameCmd  string
	GameArgs []string
	Timeout  time.Duration // Timeout waiting for game RPC server
}

// LaunchCommand handles the `autoebiten launch` functionality.
type LaunchCommand struct {
	options     *LaunchOptions
	outputFiles *output.FilePath
	outputMgr   *output.OutputManager
	gameProc    *os.Process
	handler     *proxy.UnifiedHandler
	listener    net.Listener
	gameExited  chan struct{}
	crashed     chan struct{}
	crashedOnce sync.Once
	done        chan struct{}
	doneOnce    sync.Once
}

// NewLaunchCommand creates a new launch command handler.
func NewLaunchCommand(options *LaunchOptions) *LaunchCommand {
	return &LaunchCommand{
		options:    options,
		gameExited: make(chan struct{}),
		crashed:    make(chan struct{}),
		done:       make(chan struct{}),
	}
}

// gameSocketPath returns the path for the game socket.
// Format: autoebiten-{LAUNCH_PID}.sock
// This is set as AUTOEBITEN_SOCKET env var for the game process.
func (lc *LaunchCommand) gameSocketPath() string {
	return rpc.SocketPath()
}

// createLaunchSocket creates and listens on the launch socket.
func (lc *LaunchCommand) createLaunchSocket(path string) (net.Listener, error) {
	// Ensure socket directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create socket directory: %w", err)
	}

	// Remove existing socket if present
	os.Remove(path)

	// Create listener
	listener, err := net.Listen("unix", path)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on socket %s: %w", path, err)
	}

	// Set socket permissions
	if err := os.Chmod(path, 0777); err != nil {
		listener.Close()
		return nil, fmt.Errorf("failed to set socket permissions: %w", err)
	}

	return listener, nil
}

// onCrashedCallback is called when the handler is in Crashed state and receives a CLI query.
// It signals the launch command to exit immediately.
func (lc *LaunchCommand) onCrashedCallback() {
	lc.crashedOnce.Do(func() {
		close(lc.crashed)
	})
}

// createHandler creates the UnifiedHandler with the onCrashed callback.
func (lc *LaunchCommand) createHandler() *proxy.UnifiedHandler {
	return proxy.NewUnifiedHandler(lc.outputMgr, lc.onCrashedCallback)
}

func (lc *LaunchCommand) Run() error {
	// Step 1: Create launch socket BEFORE starting the game. In step 2 we will ask the game
	// to use our PID for the socket path so we can predict it here.
	rpc.SetTargetPID(os.Getpid())
	lc.outputFiles = output.DerivePaths(lc.gameSocketPath())
	listener, err := lc.createLaunchSocket(lc.outputFiles.LaunchSock)
	if err != nil {
		return fmt.Errorf("failed to create launch socket: %w", err)
	}
	lc.listener = listener

	// Create log file
	logFile, err := output.CreateLogFile(lc.outputFiles.Log)
	if err != nil {
		lc.cleanup()
		return fmt.Errorf("failed to create log file: %w", err)
	}

	// Create OutputManager
	lc.outputMgr = output.NewOutputManager(logFile, lc.outputFiles.Log, lc.outputFiles.Snapshot)

	// Create the UnifiedHandler with onCrashed callback
	lc.handler = lc.createHandler()

	// Step 2: Create game command with pipes (must be done before Start())
	gameCmd, stdoutPipe, stderrPipe, err := lc.createGameCommand()
	if err != nil {
		lc.cleanup()
		return fmt.Errorf("failed to create game command: %w", err)
	}

	// Step 3: Start the game
	if err := gameCmd.Start(); err != nil {
		lc.handler.TransitionToCrashed(fmt.Errorf("failed to start game: %w", err))
		// Start accept loop to allow CLI to query for error
		go lc.acceptLoop()
		// Wait for CLI query before cleaning up
		lc.waitForExit("Failed to start game: " + err.Error())
		return fmt.Errorf("failed to start game: %w", err)
	}

	lc.gameProc = gameCmd.Process

	// Tee stdout/stderr through CarriageReturnWriter to OutputManager
	stdoutWriter := output.NewCarriageReturnWriter(lc.outputMgr)
	stderrWriter := output.NewCarriageReturnWriter(lc.outputMgr)
	go lc.teeOutput(stdoutPipe, os.Stdout, stdoutWriter)
	go lc.teeOutput(stderrPipe, os.Stderr, stderrWriter)

	// Monitor game exit in a goroutine
	go func() {
		gameCmd.Wait()
		close(lc.gameExited)
	}()

	// Step 4: Start accept loop in background (before waiting for RPC)
	go lc.acceptLoop()

	// Step 5: Wait for game RPC server to be ready (with timeout)
	gameClient, err := lc.waitForGameRPC()
	if err != nil {
		lc.handler.TransitionToCrashed(err)
		// Wait for CLI query before cleaning up
		lc.waitForExit("Failed to connect to game RPC server: " + err.Error())
		return fmt.Errorf("failed to connect to game RPC server: %w", err)
	}

	// Step 6: Transition to Connected state
	lc.handler.TransitionToConnected(gameClient)

	// Setup signal handling
	lc.setupSignalHandling()

	// Step 7: Wait for game exit
	<-lc.gameExited

	// Step 8: Wait for CLI query or timeout
	lc.waitForExit("Game exited")

	return nil
}

// waitForGameRPC polls the game's RPC server until it's ready or timeout.
func (lc *LaunchCommand) waitForGameRPC() (proxy.GameClient, error) {
	// Default timeout if not specified
	timeout := lc.options.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-lc.gameExited:
			return nil, fmt.Errorf("game exited")
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout after %v waiting for game RPC server", timeout)
		default:
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
			// Wait a bit before retrying
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// createGameCommand creates the game command with pipes set up.
// Returns the command and stdout/stderr pipes (must call Start() after this).
func (lc *LaunchCommand) createGameCommand() (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	cmd := exec.Command(lc.options.GameCmd, lc.options.GameArgs...)

	// Pass through all environment variables
	cmd.Env = os.Environ()

	// Set AUTOEBITEN_SOCKET to the game socket path
	gameSock := lc.gameSocketPath()
	cmd.Env = append(cmd.Env, "AUTOEBITEN_SOCKET="+gameSock)

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

// teeOutput copies data from src to both dst1 (terminal) and dst2 (managed writer).
func (lc *LaunchCommand) teeOutput(src io.Reader, dst1 *os.File, dst2 io.Writer) {
	reader := bufio.NewReader(src)
	for {
		data, err := reader.ReadBytes('\n')
		if len(data) > 0 {
			dst1.Write(data) // Terminal gets raw bytes (it interprets \r)
			dst2.Write(data) // CarriageReturnWriter + OutputManager
		}
		if err != nil {
			if err == io.EOF {
				// Flush any remaining data at stream end
				remaining, _ := reader.ReadBytes('\n')
				if len(remaining) > 0 {
					dst1.Write(remaining)
					dst2.Write(remaining)
				}
				// Flush the CarriageReturnWriter
				if flusher, ok := dst2.(interface{ Flush() error }); ok {
					flusher.Flush()
				}
			}
			break
		}
	}
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
			lc.handler.ProcessRequest(conn, &req)
			lc.doneOnce.Do(func() {
				close(lc.done)
			})
			return
		}

		lc.handler.ProcessRequest(conn, &req)
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
		lc.doneOnce.Do(func() {
			close(lc.done)
		})
	}()
}

// waitForExit waits for the done signal, crashed signal, or timeout.
// It exits immediately when CLI queries after crash (via crashed channel).
func (lc *LaunchCommand) waitForExit(message string) {
	fmt.Println(message)
	go printWaitingMessage()

	ctx, cancel := context.WithTimeout(context.Background(), WaitForExitTimeout)
	defer cancel()

	select {
	case <-lc.done:
		fmt.Println("Exit command received, exiting immediately.")
	case <-lc.crashed:
		fmt.Println("CLI queried after crash, exiting immediately.")
	case <-ctx.Done():
		fmt.Println("Timeout reached, exiting.")
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
	// Close listener
	if lc.listener != nil {
		lc.listener.Close()
	}

	// Remove launch socket
	if lc.outputFiles != nil {
		os.Remove(lc.outputFiles.LaunchSock)
	}

	// Remove log and snapshot files
	if lc.outputFiles != nil {
		os.Remove(lc.outputFiles.Log)
		os.Remove(lc.outputFiles.Snapshot)
	}
}

func printWaitingMessage() {
	timeout := WaitForExitTimeout
	interval := 1 * time.Second
	for {
		fmt.Printf("\033[2K\rWait %s for CLI to read final output... (Ctrl-C to interrupt)", timeout)
		time.Sleep(interval)
		timeout = timeout - interval
		if timeout <= 0 {
			break
		}
	}
}
