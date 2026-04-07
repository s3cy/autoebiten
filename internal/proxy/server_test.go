package proxy

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/s3cy/autoebiten/internal/output"
	"github.com/s3cy/autoebiten/internal/rpc"
)

// mockGameClient is a mock RPC client for testing.
type mockGameClient struct {
	responses map[string]*rpc.RPCResponse
	callCount int
}

func newMockGameClient() *mockGameClient {
	return &mockGameClient{
		responses: make(map[string]*rpc.RPCResponse),
	}
}

func (m *mockGameClient) SendRequest(req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	m.callCount++
	if resp, ok := m.responses[req.Method]; ok {
		return resp, nil
	}
	return &rpc.RPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  json.RawMessage(`{"success": true}`),
	}, nil
}

func (m *mockGameClient) Close() error {
	return nil
}

func (m *mockGameClient) SetResponse(method string, resp *rpc.RPCResponse) {
	m.responses[method] = resp
}

// createTestOutputManager creates a test OutputManager with temp files.
// Opens existing files in append mode, creates new files if they don't exist.
func createTestOutputManager(t *testing.T, logPath, snapPath string) *output.OutputManager {
	// Open or create log file in append mode
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		t.Fatalf("Failed to open log file: %v", err)
	}
	return output.NewOutputManager(logFile, logPath, snapPath)
}

func TestNewServer(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	mock := newMockGameClient()
	paths := &output.FilePath{
		Log:        logPath,
		Snapshot:   snapPath,
		LaunchSock: filepath.Join(tmpDir, "test-launch.sock"),
	}
	outputMgr := createTestOutputManager(t, logPath, snapPath)

	server := NewServer(mock, outputMgr, paths)
	if server == nil {
		t.Fatal("NewServer returned nil")
	}
	if server.gameClient != mock {
		t.Error("gameClient not set correctly")
	}
	if server.outputFiles != paths {
		t.Error("outputFiles not set correctly")
	}
}

func TestForwardRequest(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	// Create initial log with some content
	os.WriteFile(logPath, []byte("initial line\n"), 0600)
	os.WriteFile(snapPath, []byte("initial line\n"), 0600)

	mock := newMockGameClient()
	paths := &output.FilePath{
		Log:        logPath,
		Snapshot:   snapPath,
		LaunchSock: filepath.Join(tmpDir, "test-launch.sock"),
	}
	outputMgr := createTestOutputManager(t, logPath, snapPath)

	server := NewServer(mock, outputMgr, paths)

	// Add new content to log (simulating game output during command)
	os.WriteFile(logPath, []byte("initial line\nnew line\n"), 0600)

	// Test forwarding a request
	req := &rpc.RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "input",
		Params:  json.RawMessage(`{"action":"press","key":"KeySpace"}`),
	}

	resp, err := server.ForwardRequest(req)
	if err != nil {
		t.Fatalf("ForwardRequest failed: %v", err)
	}

	// Verify response
	if resp == nil {
		t.Fatal("Expected non-nil response")
	}
	if resp.JSONRPC != "2.0" {
		t.Errorf("JSONRPC = %q, want %q", resp.JSONRPC, "2.0")
	}
	if resp.Output == "" {
		t.Error("Expected non-empty diff output")
	}

	// Verify snapshot was updated
	snapContent, _ := os.ReadFile(snapPath)
	if string(snapContent) != "initial line\nnew line\n" {
		t.Errorf("Snapshot not updated correctly: %q", snapContent)
	}

	// Verify game client was called
	if mock.callCount != 1 {
		t.Errorf("Game client called %d times, want 1", mock.callCount)
	}
}

func TestForwardRequestWithEmptyLog(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	mock := newMockGameClient()
	paths := &output.FilePath{
		Log:        logPath,
		Snapshot:   snapPath,
		LaunchSock: filepath.Join(tmpDir, "test-launch.sock"),
	}
	outputMgr := createTestOutputManager(t, logPath, snapPath)

	server := NewServer(mock, outputMgr, paths)

	// Create log during command
	os.WriteFile(logPath, []byte("new output\n"), 0600)

	req := &rpc.RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "ping",
	}

	resp, err := server.ForwardRequest(req)
	if err != nil {
		t.Fatalf("ForwardRequest failed: %v", err)
	}

	if resp.Output == "" {
		t.Error("Expected non-empty diff for new output")
	}
}

func TestForwardRequestNoChanges(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	// Same content in both
	os.WriteFile(logPath, []byte("same content\n"), 0600)
	os.WriteFile(snapPath, []byte("same content\n"), 0600)

	mock := newMockGameClient()
	paths := &output.FilePath{
		Log:        logPath,
		Snapshot:   snapPath,
		LaunchSock: filepath.Join(tmpDir, "test-launch.sock"),
	}
	outputMgr := createTestOutputManager(t, logPath, snapPath)

	server := NewServer(mock, outputMgr, paths)

	req := &rpc.RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "ping",
	}

	resp, err := server.ForwardRequest(req)
	if err != nil {
		t.Fatalf("ForwardRequest failed: %v", err)
	}

	// Diff should be empty when no changes
	if resp.Output != "" {
		t.Errorf("Expected empty diff for unchanged content, got: %q", resp.Output)
	}
}

func TestForwardRequestVariousMethods(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	mock := newMockGameClient()
	paths := &output.FilePath{
		Log:        logPath,
		Snapshot:   snapPath,
		LaunchSock: filepath.Join(tmpDir, "test-launch.sock"),
	}
	outputMgr := createTestOutputManager(t, logPath, snapPath)

	server := NewServer(mock, outputMgr, paths)

	// Test forwarding various methods
	tests := []struct {
		name   string
		method string
		params json.RawMessage
	}{
		{"input", "input", json.RawMessage(`{"action":"press","key":"KeySpace"}`)},
		{"mouse", "mouse", json.RawMessage(`{"action":"press","x":100,"y":200}`)},
		{"wheel", "wheel", json.RawMessage(`{"x":0,"y":-3}`)},
		{"screenshot", "screenshot", json.RawMessage(`{"output":"test.png"}`)},
		{"ping", "ping", nil},
		{"get_mouse_position", "get_mouse_position", nil},
		{"get_wheel_position", "get_wheel_position", nil},
		{"list_custom_commands", "list_custom_commands", nil},
		{"custom", "custom", json.RawMessage(`{"name":"test","request":"data"}`)},
		{"exit", "exit", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset log for each test
			os.WriteFile(logPath, []byte(tt.name+" output\n"), 0600)
			os.WriteFile(snapPath, []byte{}, 0600)

			req := &rpc.RPCRequest{
				JSONRPC: "2.0",
				ID:      1,
				Method:  tt.method,
				Params:  tt.params,
			}

			resp, err := server.ForwardRequest(req)
			if err != nil {
				t.Fatalf("%s failed: %v", tt.name, err)
			}
			if resp == nil {
				t.Fatal("Expected non-nil response")
			}
		})
	}
}

func TestCleanup(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	// Create files
	os.WriteFile(logPath, []byte("log content"), 0600)
	os.WriteFile(snapPath, []byte("snapshot content"), 0600)

	paths := &output.FilePath{
		Log:        logPath,
		Snapshot:   snapPath,
		LaunchSock: filepath.Join(tmpDir, "test-launch.sock"),
	}
	outputMgr := createTestOutputManager(t, logPath, snapPath)

	server := NewServer(nil, outputMgr, paths)
	err := server.Cleanup()
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Verify files are removed
	if _, err := os.Stat(logPath); !os.IsNotExist(err) {
		t.Error("Log file should be removed")
	}
	if _, err := os.Stat(snapPath); !os.IsNotExist(err) {
		t.Error("Snapshot file should be removed")
	}
}

func TestCleanupNonExistentFiles(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "nonexistent.log")
	snapPath := filepath.Join(tmpDir, "nonexistent-snapshot.log")

	paths := &output.FilePath{
		Log:        logPath,
		Snapshot:   snapPath,
		LaunchSock: filepath.Join(tmpDir, "test-launch.sock"),
	}
	outputMgr := createTestOutputManager(t, logPath, snapPath)

	server := NewServer(nil, outputMgr, paths)
	err := server.Cleanup()
	if err != nil {
		t.Fatalf("Cleanup should succeed for non-existent files: %v", err)
	}
}

func TestConcurrentRequests(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	mock := newMockGameClient()
	paths := &output.FilePath{
		Log:        logPath,
		Snapshot:   snapPath,
		LaunchSock: filepath.Join(tmpDir, "test-launch.sock"),
	}
	outputMgr := createTestOutputManager(t, logPath, snapPath)

	server := NewServer(mock, outputMgr, paths)

	// Simulate concurrent requests
	done := make(chan bool, 3)
	for i := range make([]struct{}, 3) {
		go func(n int) {
			os.WriteFile(logPath, []byte(fmt.Sprintf("output %d\n", n)), 0600)
			req := &rpc.RPCRequest{
				JSONRPC: "2.0",
				ID:      n,
				Method:  "ping",
			}
			_, err := server.ForwardRequest(req)
			if err != nil {
				t.Errorf("Concurrent request %d failed: %v", n, err)
			}
			done <- true
		}(i)
	}

	// Wait for all to complete
	for range make([]struct{}, 3) {
		<-done
	}

	// All 3 should have succeeded (serialized via mutex)
	if mock.callCount != 3 {
		t.Errorf("Expected 3 calls, got %d", mock.callCount)
	}
}

func TestHandlerProcessRequest(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	mock := newMockGameClient()
	paths := &output.FilePath{
		Log:        logPath,
		Snapshot:   snapPath,
		LaunchSock: filepath.Join(tmpDir, "test-launch.sock"),
	}
	outputMgr := createTestOutputManager(t, logPath, snapPath)

	server := NewServer(mock, outputMgr, paths)
	handler := NewHandler(server)

	// Create a pipe to simulate a connection
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	// Write log content
	os.WriteFile(logPath, []byte("game output\n"), 0600)

	// Send a request through the handler
	req := &rpc.RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "proxy_ping", // Will be stripped to "ping"
	}

	// Process in background
	go handler.ProcessRequest(serverConn, req)

	// Read response
	decoder := json.NewDecoder(clientConn)
	var resp Response
	if err := decoder.Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.JSONRPC != "2.0" {
		t.Errorf("JSONRPC = %q, want %q", resp.JSONRPC, "2.0")
	}
	if resp.Output == "" {
		t.Error("Expected non-empty output in response")
	}
}

// trackingMockClient tracks method calls
type trackingMockClient struct {
	onSendRequest func(req *rpc.RPCRequest)
	responses     map[string]*rpc.RPCResponse
}

func (m *trackingMockClient) SendRequest(req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	if m.onSendRequest != nil {
		m.onSendRequest(req)
	}
	if resp, ok := m.responses[req.Method]; ok {
		return resp, nil
	}
	return &rpc.RPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  json.RawMessage(`{"success": true}`),
	}, nil
}

func (m *trackingMockClient) Close() error {
	return nil
}

func (m *trackingMockClient) SetResponse(method string, resp *rpc.RPCResponse) {
	if m.responses == nil {
		m.responses = make(map[string]*rpc.RPCResponse)
	}
	m.responses[method] = resp
}

// delayedMockClient simulates a game client with delay.
type delayedMockClient struct {
	callCount int
	delay     time.Duration
}

func (m *delayedMockClient) SendRequest(req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	m.callCount++
	time.Sleep(m.delay)
	return &rpc.RPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  json.RawMessage(`{"ok": true, "version": "test"}`),
	}, nil
}

func (m *delayedMockClient) Close() error {
	return nil
}