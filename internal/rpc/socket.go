package rpc

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

const (
	// DefaultSocketDir is the default directory for the Unix socket.
	DefaultSocketDir = "/tmp/autoebiten"
	// DefaultSocketName is the default socket filename.
	DefaultSocketName = "autoebiten.sock"
)

var targetPID int

// SetTargetPID sets the target process PID for socket connection.
// This allows the CLI to connect to a specific game instance.
func SetTargetPID(pid int) {
	targetPID = pid
}

// GetTargetPID returns the current target PID.
func GetTargetPID() int {
	if targetPID > 0 {
		return targetPID
	}
	return os.Getpid()
}

// GameInfo represents a running game instance.
type GameInfo struct {
	PID  int
	Name string
}

// findRunningGames scans for running game processes with autoebiten sockets.
func findRunningGames() ([]GameInfo, error) {
	// Find all socket files in the autoebiten directory
	socketDir := filepath.Dir(SocketPath())

	entries, err := os.ReadDir(socketDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read socket directory: %w", err)
	}

	var games []GameInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasPrefix(name, "autoebiten-") || !strings.HasSuffix(name, ".sock") {
			continue
		}

		// Extract PID from filename: autoebiten-{PID}.sock
		pidStr := strings.TrimPrefix(name, "autoebiten-")
		pidStr = strings.TrimSuffix(pidStr, ".sock")
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}

		// Check if process is actually running via syscall
		if err := syscall.Kill(pid, 0); err != nil {
			// Process is dead, remove stale socket file
			socketPath := filepath.Join(socketDir, name)
			os.Remove(socketPath)
			continue
		}

		// Try to get process name
		procName := fmt.Sprintf("process-%d", pid)
		cmdLine, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "comm=").Output()
		if err == nil {
			procName = strings.TrimSpace(string(cmdLine))
		}

		games = append(games, GameInfo{PID: pid, Name: procName})
	}

	// Sort by PID for consistent output
	sort.Slice(games, func(i, j int) bool {
		return games[i].PID < games[j].PID
	})

	return games, nil
}

// AutoSelectGame finds a single game or returns an error if multiple are found.
func AutoSelectGame() (*GameInfo, error) {
	games, err := findRunningGames()
	if err != nil {
		return nil, err
	}

	switch len(games) {
	case 0:
		return nil, fmt.Errorf("no running game found")
	case 1:
		return &games[0], nil
	default:
		var msg strings.Builder
		msg.WriteString("multiple game instances running, use --pid to specify one:\n")
		for _, g := range games {
			fmt.Fprintf(&msg, "  %d: %s\n", g.PID, g.Name)
		}
		return nil, fmt.Errorf("%s", msg.String())
	}
}

// SocketPath returns the path to the Unix socket.
// It respects the AUTOEBITEN_SOCKET environment variable.
// By default, uses PID-based naming to allow multiple game instances.
// If SetTargetPID was called, uses that PID for the socket path.
func SocketPath() string {
	if path := os.Getenv("AUTOEBITEN_SOCKET"); path != "" {
		return path
	}
	if targetPID > 0 {
		return filepath.Join(DefaultSocketDir, fmt.Sprintf("autoebiten-%d.sock", targetPID))
	}
	return filepath.Join(DefaultSocketDir, fmt.Sprintf("autoebiten-%d.sock", os.Getpid()))
}

// EnsureSocketDir ensures the socket directory exists.
func EnsureSocketDir() error {
	dir := filepath.Dir(SocketPath())
	return os.MkdirAll(dir, 0755)
}

// Request represents an incoming RPC request with its connection for writing responses.
type Request struct {
	Req  *RPCRequest
	Conn net.Conn
}

// Serve starts a JSON-RPC server on the Unix socket.
// It runs the accept loop in a background goroutine and returns a channel
// of incoming requests. Each request must be processed by calling ProcessRequest,
// and the response written to the request's Conn.
//
// The caller should range over the returned channel to process requests one by one.
// The channel is closed when the listener is closed.
func Serve() (<-chan *Request, error) {
	if err := EnsureSocketDir(); err != nil {
		return nil, fmt.Errorf("failed to create socket directory: %w", err)
	}

	path := SocketPath()
	// Remove existing socket file if present
	os.Remove(path)

	listener, err := net.Listen("unix", path)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on socket %s: %w", path, err)
	}

	// Set socket permissions to allow communication
	if err := os.Chmod(path, 0777); err != nil {
		listener.Close()
		return nil, fmt.Errorf("failed to set socket permissions: %w", err)
	}

	reqChan := make(chan *Request)

	go func() {
		defer close(reqChan)
		defer listener.Close()

		for {
			conn, err := listener.Accept()
			if err != nil {
				if netErr, ok := err.(*net.OpError); ok && netErr.Err.Error() == "use of closed network connection" {
					break
				}
				fmt.Fprintf(os.Stderr, "autoebiten: accept error: %v\n", err)
				continue
			}

			go handleConnection(conn, reqChan)
		}
	}()

	return reqChan, nil
}

func handleConnection(conn net.Conn, reqChan chan<- *Request) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	for {
		var req RPCRequest
		if err := decoder.Decode(&req); err != nil {
			if err == io.EOF {
				return
			}
			fmt.Fprintf(os.Stderr, "autoebiten: decode error: %v\n", err)
			return
		}

		reqChan <- &Request{Req: &req, Conn: conn}
	}
}

// Client is a JSON-RPC client that connects to a Unix socket.
type Client struct {
	conn    net.Conn
	encoder *json.Encoder
	decoder *json.Decoder
	mu      sync.Mutex
}

func NewClientWithPath(path string) (*Client, error) {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to socket %s: %w", path, err)
	}

	return &Client{
		conn:    conn,
		encoder: json.NewEncoder(conn),
		decoder: json.NewDecoder(conn),
	}, nil
}

func NewClient() (*Client, error) {
	return NewClientWithPath(SocketPath())
}

// Close closes the client connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// SendRequest sends a request and waits for a response.
func (c *Client) SendRequest(req *RPCRequest) (*RPCResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.encoder.Encode(req); err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	var resp RPCResponse
	if err := c.decoder.Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resp, nil
}

// SendRequestSocket is a convenience function that creates a client,
// sends a request, and closes the connection.
// Use this for single requests like CLI commands.
func SendRequestSocket(req *RPCRequest) (*RPCResponse, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return client.SendRequest(req)
}

// BuildRequest creates a new RPC request with the given method and params.
func BuildRequest(method string, params any) (*RPCRequest, error) {
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %w", err)
	}

	return &RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  method,
		Params:  paramsJSON,
	}, nil
}
