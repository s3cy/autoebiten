package proxy

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/s3cy/autoebiten/internal/output"
	"github.com/s3cy/autoebiten/internal/rpc"
)

// mockGameClient is a mock RPC client for testing.
type mockGameClient struct {
	responses   map[string]*rpc.RPCResponse
	callCount   int
	shouldError error
	mu          sync.Mutex
}

func newMockGameClient() *mockGameClient {
	return &mockGameClient{
		responses: make(map[string]*rpc.RPCResponse),
	}
}

func (m *mockGameClient) SendRequest(req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callCount++
	if m.shouldError != nil {
		return nil, m.shouldError
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

func (m *mockGameClient) Close() error {
	return nil
}

func (m *mockGameClient) SetResponse(method string, resp *rpc.RPCResponse) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses[method] = resp
}

func (m *mockGameClient) SetError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldError = err
}

func (m *mockGameClient) GetCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.callCount
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

// ==================== UnifiedHandler Tests ====================

func setupUnifiedHandlerTest(t *testing.T) (*UnifiedHandler, *mockGameClient, *output.OutputManager, string, string) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	mock := newMockGameClient()
	outputMgr := createTestOutputManager(t, logPath, snapPath)

	// Create initial log file content
	os.WriteFile(logPath, []byte("initial line\n"), 0600)
	os.WriteFile(snapPath, []byte("initial line\n"), 0600)

	handler := NewUnifiedHandler(outputMgr, nil)

	return handler, mock, outputMgr, logPath, snapPath
}

func TestNewUnifiedHandler(t *testing.T) {
	_, _, outputMgr, _, _ := setupUnifiedHandlerTest(t)

	handler := NewUnifiedHandler(outputMgr, nil)

	if handler == nil {
		t.Fatal("NewUnifiedHandler returned nil")
	}
	if handler.GetState() != StateWaiting {
		t.Errorf("Initial state = %v, want %v", handler.GetState(), StateWaiting)
	}
	if handler.outputMgr != outputMgr {
		t.Error("outputMgr not set correctly")
	}
	if handler.gameClient != nil {
		t.Error("gameClient should be nil initially")
	}
	if handler.launchError != nil {
		t.Error("launchError should be nil initially")
	}
}

func TestUnifiedHandlerStateWaiting(t *testing.T) {
	handler, _, _, logPath, _ := setupUnifiedHandlerTest(t)

	// Add new content to log to generate diff
	os.WriteFile(logPath, []byte("initial line\nnew line\n"), 0600)

	// Create pipe to simulate connection
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	// Send request
	req := &rpc.RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "input",
		Params:  json.RawMessage(`{"action":"press","key":"KeySpace"}`),
	}

	// Process in goroutine
	go handler.ProcessRequest(serverConn, req)

	// Read response
	decoder := json.NewDecoder(clientConn)
	var resp rpc.RPCResponse
	if err := decoder.Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify error response
	if resp.Error == nil {
		t.Fatal("Expected error response in waiting state")
	}
	if resp.Error.Code != rpc.ErrGameNotConnected {
		t.Errorf("Error code = %d, want %d", resp.Error.Code, rpc.ErrGameNotConnected)
	}
	if resp.Error.Message != "game not connected" {
		t.Errorf("Error message = %q, want %q", resp.Error.Message, "game not connected")
	}

	// Verify extra fields
	if resp.Extra == nil {
		t.Fatal("Expected Extra field in response")
	}

	// Diff should be present (new line added)
	diff, ok := resp.Extra["diff"].(string)
	if !ok {
		t.Fatal("Expected diff in Extra field")
	}
	if diff == "" {
		t.Error("Expected non-empty diff")
	}

	// proxy_error should be empty in waiting state
	proxyErr, ok := resp.Extra["proxy_error"].(string)
	if !ok {
		t.Fatal("Expected proxy_error in Extra field")
	}
	if proxyErr != "" {
		t.Errorf("proxy_error should be empty in waiting state, got %q", proxyErr)
	}
}

func TestUnifiedHandlerTransitionToConnected(t *testing.T) {
	handler, mock, _, _, _ := setupUnifiedHandlerTest(t)

	handler.TransitionToConnected(mock)

	if handler.GetState() != StateConnected {
		t.Errorf("State = %v, want %v", handler.GetState(), StateConnected)
	}
	if handler.gameClient != mock {
		t.Error("gameClient not set correctly")
	}
}

func TestUnifiedHandlerStateConnected(t *testing.T) {
	handler, mock, _, logPath, _ := setupUnifiedHandlerTest(t)

	// Set up mock response
	mock.SetResponse("input", &rpc.RPCResponse{
		JSONRPC: "2.0",
		ID:      1,
		Result:  json.RawMessage(`{"success": true}`),
	})

	// Transition to connected
	handler.TransitionToConnected(mock)

	// Add new content to log
	os.WriteFile(logPath, []byte("initial line\nnew line\n"), 0600)

	// Create pipe
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	// Send request
	req := &rpc.RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "input",
		Params:  json.RawMessage(`{"action":"press","key":"KeySpace"}`),
	}

	// Process
	go handler.ProcessRequest(serverConn, req)

	// Read response
	decoder := json.NewDecoder(clientConn)
	var resp rpc.RPCResponse
	if err := decoder.Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify success response
	if resp.Error != nil {
		t.Errorf("Unexpected error: %v", resp.Error)
	}
	if resp.Result == nil {
		t.Error("Expected result in response")
	}

	// Verify diff is included
	if resp.Extra == nil {
		t.Fatal("Expected Extra field")
	}
	diff, ok := resp.Extra["diff"].(string)
	if !ok {
		t.Fatal("Expected diff in Extra field")
	}
	if diff == "" {
		t.Error("Expected non-empty diff")
	}

	// Verify game client was called
	if mock.GetCallCount() != 1 {
		t.Errorf("Game client called %d times, want 1", mock.GetCallCount())
	}
}

func TestUnifiedHandlerConnectedToCrashed(t *testing.T) {
	handler, mock, _, logPath, _ := setupUnifiedHandlerTest(t)

	// Set up mock to return error
	mock.SetError(fmt.Errorf("connection reset"))

	// Transition to connected
	handler.TransitionToConnected(mock)

	// Add new content to log
	os.WriteFile(logPath, []byte("initial line\nerror line\n"), 0600)

	// Create pipe
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	// Send request
	req := &rpc.RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "input",
		Params:  json.RawMessage(`{"action":"press","key":"KeySpace"}`),
	}

	// Process
	go handler.ProcessRequest(serverConn, req)

	// Read response
	decoder := json.NewDecoder(clientConn)
	var resp rpc.RPCResponse
	if err := decoder.Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Should return error
	if resp.Error == nil {
		t.Fatal("Expected error response after game failure")
	}
	if resp.Error.Code != rpc.ErrGameNotConnected {
		t.Errorf("Error code = %d, want %d", resp.Error.Code, rpc.ErrGameNotConnected)
	}

	// Handler should now be in crashed state
	if handler.GetState() != StateCrashed {
		t.Errorf("State = %v, want %v", handler.GetState(), StateCrashed)
	}

	// proxy_error should contain the error
	proxyErr, ok := resp.Extra["proxy_error"].(string)
	if !ok {
		t.Fatal("Expected proxy_error in Extra field")
	}
	if proxyErr == "" {
		t.Error("Expected non-empty proxy_error after crash")
	}
	if !contains(proxyErr, "connection reset") {
		t.Errorf("proxy_error should contain 'connection reset', got %q", proxyErr)
	}
}

func TestUnifiedHandlerStateCrashed(t *testing.T) {
	var onCrashedCalledMu sync.Mutex
	onCrashedCalled := false
	onCrashed := func() {
		onCrashedCalledMu.Lock()
		defer onCrashedCalledMu.Unlock()
		onCrashedCalled = true
	}

	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	os.WriteFile(logPath, []byte("initial line\n"), 0600)
	os.WriteFile(snapPath, []byte("initial line\n"), 0600)

	outputMgr := createTestOutputManager(t, logPath, snapPath)
	handler := NewUnifiedHandler(outputMgr, onCrashed)

	// Transition to crashed with error
	handler.TransitionToCrashed(fmt.Errorf("game process exited with code 1"))

	// Add new content
	os.WriteFile(logPath, []byte("initial line\ncrash output\n"), 0600)

	// Create pipe
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	// Send request
	req := &rpc.RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "input",
	}

	// Process
	go handler.ProcessRequest(serverConn, req)

	// Read response
	decoder := json.NewDecoder(clientConn)
	var resp rpc.RPCResponse
	if err := decoder.Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Should return error
	if resp.Error == nil {
		t.Fatal("Expected error response in crashed state")
	}

	// proxy_error should contain accumulated error
	proxyErr, ok := resp.Extra["proxy_error"].(string)
	if !ok {
		t.Fatal("Expected proxy_error in Extra field")
	}
	if !contains(proxyErr, "game process exited") {
		t.Errorf("proxy_error should contain crash error, got %q", proxyErr)
	}

	// Wait for onCrashed to be called (called async)
	time.Sleep(100 * time.Millisecond)
	onCrashedCalledMu.Lock()
	if !onCrashedCalled {
		t.Error("Expected onCrashed callback to be called")
	}
	onCrashedCalledMu.Unlock()
}

func TestUnifiedHandlerTransitionToCrashedAccumulatesErrors(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	os.WriteFile(logPath, []byte(""), 0600)
	os.WriteFile(snapPath, []byte(""), 0600)

	outputMgr := createTestOutputManager(t, logPath, snapPath)
	handler := NewUnifiedHandler(outputMgr, nil)

	// Transition to crashed multiple times
	handler.TransitionToCrashed(fmt.Errorf("first error"))
	handler.TransitionToCrashed(fmt.Errorf("second error"))
	handler.TransitionToCrashed(fmt.Errorf("third error"))

	// Check state
	if handler.GetState() != StateCrashed {
		t.Errorf("State = %v, want %v", handler.GetState(), StateCrashed)
	}

	// Create pipe for request
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	req := &rpc.RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "ping",
	}

	go handler.ProcessRequest(serverConn, req)

	decoder := json.NewDecoder(clientConn)
	var resp rpc.RPCResponse
	if err := decoder.Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Check all errors are accumulated
	proxyErr := resp.Extra["proxy_error"].(string)
	if !contains(proxyErr, "first error") {
		t.Errorf("Expected first error in proxy_error, got %q", proxyErr)
	}
	if !contains(proxyErr, "second error") {
		t.Errorf("Expected second error in proxy_error, got %q", proxyErr)
	}
	if !contains(proxyErr, "third error") {
		t.Errorf("Expected third error in proxy_error, got %q", proxyErr)
	}
}

func TestUnifiedHandlerConcurrentRequests(t *testing.T) {
	handler, mock, _, logPath, _ := setupUnifiedHandlerTest(t)

	// Transition to connected
	mock.SetResponse("ping", &rpc.RPCResponse{
		JSONRPC: "2.0",
		ID:      nil,
		Result:  json.RawMessage(`{"ok": true}`),
	})
	handler.TransitionToConnected(mock)

	// Run concurrent requests
	const numRequests = 10
	var wg sync.WaitGroup
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
		go func(n int) {
			defer wg.Done()

			// Add unique content for each request
			os.WriteFile(logPath, []byte(fmt.Sprintf("output %d\n", n)), 0600)

			clientConn, serverConn := net.Pipe()
			defer clientConn.Close()
			defer serverConn.Close()

			req := &rpc.RPCRequest{
				JSONRPC: "2.0",
				ID:      n,
				Method:  "ping",
			}

			go handler.ProcessRequest(serverConn, req)

			// Read response
			decoder := json.NewDecoder(clientConn)
			var resp rpc.RPCResponse
			if err := decoder.Decode(&resp); err != nil {
				t.Errorf("Request %d: Failed to decode: %v", n, err)
				return
			}

			if resp.Error != nil {
				t.Errorf("Request %d: Unexpected error: %v", n, resp.Error)
			}
		}(i)
	}

	wg.Wait()

	// All requests should have been processed
	if mock.GetCallCount() != numRequests {
		t.Errorf("Expected %d calls, got %d", numRequests, mock.GetCallCount())
	}
}

func TestProxyStateString(t *testing.T) {
	tests := []struct {
		state    ProxyState
		expected string
	}{
		{StateWaiting, "waiting"},
		{StateConnected, "connected"},
		{StateCrashed, "crashed"},
		{ProxyState(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.state.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsInternal(s, substr))
}

func containsInternal(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ==================== Legacy Server Tests ====================

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
	// Extract diff from Extra field
	diff, _ := resp.Extra["diff"].(string)
	if diff == "" {
		t.Error("Expected non-empty diff output")
	}

	// Verify snapshot was updated
	snapContent, _ := os.ReadFile(snapPath)
	if string(snapContent) != "initial line\nnew line\n" {
		t.Errorf("Snapshot not updated correctly: %q", snapContent)
	}

	// Verify game client was called
	if mock.GetCallCount() != 1 {
		t.Errorf("Game client called %d times, want 1", mock.GetCallCount())
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

	// Extract diff from Extra field
	diff, _ := resp.Extra["diff"].(string)
	if diff == "" {
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

	// Extract diff from Extra field
	diff, _ := resp.Extra["diff"].(string)
	// Diff should be empty when no changes
	if diff != "" {
		t.Errorf("Expected empty diff for unchanged content, got: %q", diff)
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
	if mock.GetCallCount() != 3 {
		t.Errorf("Expected 3 calls, got %d", mock.GetCallCount())
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
	var resp rpc.RPCResponse
	if err := decoder.Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resp.JSONRPC != "2.0" {
		t.Errorf("JSONRPC = %q, want %q", resp.JSONRPC, "2.0")
	}
	// Extract diff from Extra field
	diff, _ := resp.Extra["diff"].(string)
	if diff == "" {
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
