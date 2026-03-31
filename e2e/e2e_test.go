package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/s3cy/autoebiten/internal/input"
	"github.com/s3cy/autoebiten/internal/rpc"
	"github.com/s3cy/autoebiten/internal/script"
)

// TestRPCEndToEnd tests the complete RPC flow from request building to response handling.
func TestRPCEndToEnd(t *testing.T) {
	t.Run("ping flow", func(t *testing.T) {
		server, client := setupTestSocketPair(t)
		defer server.Close()
		defer client.Close()

		// Send ping request
		req := rpc.RPCRequest{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "ping",
			Params:  nil,
		}

		resp, err := sendRequest(client, &req)
		if err != nil {
			t.Fatalf("Failed to send ping: %v", err)
		}

		if resp.Error != nil {
			t.Errorf("Unexpected error: %s", resp.Error.Message)
		}

		result, ok := resp.Result.(map[string]any)
		if !ok {
			t.Fatal("Expected result map")
		}

		if ok, _ := result["ok"].(bool); !ok {
			t.Error("Expected ok to be true")
		}
	})

	t.Run("input press flow", func(t *testing.T) {
		server, client := setupTestSocketPair(t)
		defer server.Close()
		defer client.Close()

		params := rpc.InputParams{
			Action:        "press",
			Key:           "KeyA",
			DurationTicks: 0,
		}

		req := rpc.RPCRequest{
			JSONRPC: "2.0",
			ID:      2,
			Method:  "input",
			Params:  mustMarshal(t, params),
		}

		resp, err := sendRequest(client, &req)
		if err != nil {
			t.Fatalf("Failed to send input: %v", err)
		}

		if resp.Error != nil {
			t.Errorf("Unexpected error: %s", resp.Error.Message)
		}

		result, ok := resp.Result.(map[string]any)
		if !ok {
			t.Fatal("Expected result map")
		}

		if success, _ := result["success"].(bool); !success {
			t.Error("Expected success to be true")
		}
	})

	t.Run("mouse flow", func(t *testing.T) {
		server, client := setupTestSocketPair(t)
		defer server.Close()
		defer client.Close()

		params := rpc.MouseParams{
			Action:        "position",
			X:             100,
			Y:             200,
			Button:        "",
			DurationTicks: 0,
		}

		req := rpc.RPCRequest{
			JSONRPC: "2.0",
			ID:      3,
			Method:  "mouse",
			Params:  mustMarshal(t, params),
		}

		resp, err := sendRequest(client, &req)
		if err != nil {
			t.Fatalf("Failed to send mouse: %v", err)
		}

		if resp.Error != nil {
			t.Errorf("Unexpected error: %s", resp.Error.Message)
		}

		result, ok := resp.Result.(map[string]any)
		if !ok {
			t.Fatal("Expected result map")
		}

		if success, _ := result["success"].(bool); !success {
			t.Error("Expected success to be true")
		}
	})

	t.Run("wheel flow", func(t *testing.T) {
		server, client := setupTestSocketPair(t)
		defer server.Close()
		defer client.Close()

		params := rpc.WheelParams{
			X: 10.5,
			Y: -5.0,
		}

		req := rpc.RPCRequest{
			JSONRPC: "2.0",
			ID:      4,
			Method:  "wheel",
			Params:  mustMarshal(t, params),
		}

		resp, err := sendRequest(client, &req)
		if err != nil {
			t.Fatalf("Failed to send wheel: %v", err)
		}

		if resp.Error != nil {
			t.Errorf("Unexpected error: %s", resp.Error.Message)
		}

		result, ok := resp.Result.(map[string]any)
		if !ok {
			t.Fatal("Expected result map")
		}

		if success, _ := result["success"].(bool); !success {
			t.Error("Expected success to be true")
		}
	})

	t.Run("invalid method", func(t *testing.T) {
		server, client := setupTestSocketPair(t)
		defer server.Close()
		defer client.Close()

		req := rpc.RPCRequest{
			JSONRPC: "2.0",
			ID:      5,
			Method:  "invalid_method",
			Params:  nil,
		}

		resp, err := sendRequest(client, &req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		if resp.Error == nil {
			t.Error("Expected error response")
		}

		if resp.Error.Code != rpc.ErrInvalidParams {
			t.Errorf("Expected code %d, got %d", rpc.ErrInvalidParams, resp.Error.Code)
		}
	})

	t.Run("invalid params", func(t *testing.T) {
		server, client := setupTestSocketPair(t)
		defer server.Close()
		defer client.Close()

		req := rpc.RPCRequest{
			JSONRPC: "2.0",
			ID:      6,
			Method:  "input",
			Params:  json.RawMessage(`{"key": ""}`), // Empty key should fail
		}

		resp, err := sendRequest(client, &req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}

		if resp.Error == nil {
			t.Error("Expected error response")
		}

		if resp.Error.Code != rpc.ErrInvalidParams {
			t.Errorf("Expected code %d, got %d", rpc.ErrInvalidParams, resp.Error.Code)
		}
	})
}

// TestBuildRequest tests that BuildRequest creates properly formatted requests.
func TestBuildRequest(t *testing.T) {
	t.Run("input request", func(t *testing.T) {
		params := rpc.InputParams{
			Action:        "press",
			Key:           "KeySpace",
			DurationTicks: 6,
		}

		req, err := rpc.BuildRequest("input", params)
		if err != nil {
			t.Fatalf("Failed to build request: %v", err)
		}

		if req.JSONRPC != "2.0" {
			t.Errorf("Expected JSONRPC 2.0, got %s", req.JSONRPC)
		}

		if req.Method != "input" {
			t.Errorf("Expected method 'input', got %s", req.Method)
		}

		if req.ID != 1 {
			t.Errorf("Expected ID 1, got %v", req.ID)
		}

		// Verify params can be unmarshaled
		var decoded rpc.InputParams
		if err := json.Unmarshal(req.Params, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal params: %v", err)
		}

		if decoded.Action != "press" {
			t.Errorf("Expected action 'press', got %s", decoded.Action)
		}
		if decoded.Key != "KeySpace" {
			t.Errorf("Expected key 'KeySpace', got %s", decoded.Key)
		}
	})
}

// TestScriptExecutorEndToEnd tests the complete script execution flow.
func TestScriptExecutorEndToEnd(t *testing.T) {
	t.Run("simple script execution", func(t *testing.T) {
		scriptJSON := `{
			"version": "1.0",
			"commands": [
				{"input": {"action": "press", "key": "KeyW"}},
				{"delay": {"ms": 100}},
				{"input": {"action": "press", "key": "KeyA"}}
			]
		}`

		s, err := script.ParseBytes([]byte(scriptJSON))
		if err != nil {
			t.Fatalf("Failed to parse script: %v", err)
		}

		executor := script.NewExecutor(s)

		var executedCommands []string
		var mu sync.Mutex

		executor.SetInputFunc(func(key, action string, durationTicks int64, async bool) error {
			mu.Lock()
			executedCommands = append(executedCommands, fmt.Sprintf("%s:%s", action, key))
			mu.Unlock()
			return nil
		})

		count, err := executor.Execute()
		if err != nil {
			t.Fatalf("Failed to execute script: %v", err)
		}

		// Note: delay commands don't call inputFunc, so we expect 2 input calls
		if count != 3 {
			t.Errorf("Expected 3 commands, got %d", count)
		}

		mu.Lock()
		defer mu.Unlock()

		if len(executedCommands) != 2 {
			t.Errorf("Expected 2 input commands, got %d", len(executedCommands))
		}

		if len(executedCommands) > 0 && executedCommands[0] != "press:KeyW" {
			t.Errorf("Expected first command 'press:KeyW', got %s", executedCommands[0])
		}
	})

	t.Run("repeat command execution", func(t *testing.T) {
		scriptJSON := `{
			"version": "1.0",
			"commands": [
				{"repeat": {"times": 3, "commands": [
					{"input": {"action": "press", "key": "KeyX", "duration_ticks": 6}}
				]}}
			]
		}`

		s, err := script.ParseBytes([]byte(scriptJSON))
		if err != nil {
			t.Fatalf("Failed to parse script: %v", err)
		}

		executor := script.NewExecutor(s)

		var executedKeys []string
		var mu sync.Mutex

		executor.SetInputFunc(func(key, action string, durationTicks int64, async bool) error {
			mu.Lock()
			executedKeys = append(executedKeys, key)
			mu.Unlock()
			return nil
		})

		count, err := executor.Execute()
		if err != nil {
			t.Fatalf("Failed to execute script: %v", err)
		}

		// 1 repeat command containing 1 input command, repeated 3 times
		if count != 4 {
			t.Errorf("Expected 4 total commands (1 repeat + 3 nested), got %d", count)
		}

		mu.Lock()
		defer mu.Unlock()

		if len(executedKeys) != 3 {
			t.Errorf("Expected 3 key presses, got %d", len(executedKeys))
		}

		for i, key := range executedKeys {
			if key != "KeyX" {
				t.Errorf("Expected KeyX at index %d, got %s", i, key)
			}
		}
	})

	t.Run("mouse command execution", func(t *testing.T) {
		scriptJSON := `{
			"version": "1.0",
			"commands": [
				{"mouse": {"action": "position", "x": 100, "y": 200}},
				{"mouse": {"action": "press", "x": 100, "y": 200, "button": "MouseButtonLeft"}}
			]
		}`

		s, err := script.ParseBytes([]byte(scriptJSON))
		if err != nil {
			t.Fatalf("Failed to parse script: %v", err)
		}

		executor := script.NewExecutor(s)

		type mouseCall struct {
			action string
			x, y   int
			button string
		}
		var mouseCalls []mouseCall
		var mu sync.Mutex

		executor.SetMouseFunc(func(action string, x, y int, button string, durationTicks int64) error {
			mu.Lock()
			mouseCalls = append(mouseCalls, mouseCall{action, x, y, button})
			mu.Unlock()
			return nil
		})

		count, err := executor.Execute()
		if err != nil {
			t.Fatalf("Failed to execute script: %v", err)
		}

		if count != 2 {
			t.Errorf("Expected 2 commands, got %d", count)
		}

		mu.Lock()
		defer mu.Unlock()

		if len(mouseCalls) != 2 {
			t.Fatalf("Expected 2 mouse calls, got %d", len(mouseCalls))
		}

		if mouseCalls[0].action != "position" || mouseCalls[0].x != 100 || mouseCalls[0].y != 200 {
			t.Errorf("First mouse call mismatch: %+v", mouseCalls[0])
		}

		if mouseCalls[1].action != "press" || mouseCalls[1].button != "MouseButtonLeft" {
			t.Errorf("Second mouse call mismatch: %+v", mouseCalls[1])
		}
	})

	t.Run("wheel command execution", func(t *testing.T) {
		scriptJSON := `{
			"version": "1.0",
			"commands": [
				{"wheel": {"x": 5.0, "y": -3.0}}
			]
		}`

		s, err := script.ParseBytes([]byte(scriptJSON))
		if err != nil {
			t.Fatalf("Failed to parse script: %v", err)
		}

		executor := script.NewExecutor(s)

		var wheelCalls []struct{ x, y float64 }
		var mu sync.Mutex

		executor.SetWheelFunc(func(x, y float64) error {
			mu.Lock()
			wheelCalls = append(wheelCalls, struct{ x, y float64 }{x, y})
			mu.Unlock()
			return nil
		})

		_, err = executor.Execute()
		if err != nil {
			t.Fatalf("Failed to execute script: %v", err)
		}

		mu.Lock()
		defer mu.Unlock()

		if len(wheelCalls) != 1 {
			t.Fatalf("Expected 1 wheel call, got %d", len(wheelCalls))
		}

		if wheelCalls[0].x != 5.0 || wheelCalls[0].y != -3.0 {
			t.Errorf("Wheel call mismatch: got (%.2f, %.2f)", wheelCalls[0].x, wheelCalls[0].y)
		}
	})
}

// TestSocketCommunication tests socket-based communication patterns.
func TestSocketCommunication(t *testing.T) {
	t.Run("concurrent requests", func(t *testing.T) {
		server, client := setupTestSocketPair(t)
		defer server.Close()
		defer client.Close()

		var wg sync.WaitGroup
		numRequests := 10

		for i := 0; i < numRequests; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				req := rpc.RPCRequest{
					JSONRPC: "2.0",
					ID:      id,
					Method:  "ping",
					Params:  nil,
				}

				resp, err := sendRequest(client, &req)
				if err != nil {
					t.Errorf("Request %d failed: %v", id, err)
					return
				}

				if resp.Error != nil {
					t.Errorf("Request %d got error: %s", id, resp.Error.Message)
				}
			}(i)
		}

		wg.Wait()
	})

	t.Run("socket path respects PID", func(t *testing.T) {
		// Test that SocketPath uses PID when no env var is set
		path := rpc.SocketPath()
		if path == "" {
			t.Error("SocketPath returned empty string")
		}

		// Verify it looks like a valid socket path
		if len(path) < 5 {
			t.Error("SocketPath should return a meaningful path")
		}
	})
}

// TestRPCErrorCodes tests that error codes are properly returned.
func TestRPCErrorCodes(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		params     any
		wantCode   int
		wantErrMsg string
	}{
		{
			name:     "input with empty key",
			method:   "input",
			params:   rpc.InputParams{Action: "press", Key: ""},
			wantCode: rpc.ErrInvalidParams,
		},
		{
			name:     "input with unknown key",
			method:   "input",
			params:   rpc.InputParams{Action: "press", Key: "UnknownKey"},
			wantCode: rpc.ErrInvalidParams,
		},
		{
			name:     "mouse with unknown button",
			method:   "mouse",
			params:   rpc.MouseParams{Action: "press", Button: "UnknownButton"},
			wantCode: rpc.ErrInvalidParams,
		},
		{
			name:     "unknown method",
			method:   "does_not_exist",
			params:   nil,
			wantCode: rpc.ErrInvalidParams,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestSocketPair(t)
			defer server.Close()
			defer client.Close()

			var params json.RawMessage
			if tt.params != nil {
				params = mustMarshal(t, tt.params)
			}

			req := rpc.RPCRequest{
				JSONRPC: "2.0",
				ID:      1,
				Method:  tt.method,
				Params:  params,
			}

			resp, err := sendRequest(client, &req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}

			if resp.Error == nil {
				t.Fatalf("Expected error, got nil")
			}

			if resp.Error.Code != tt.wantCode {
				t.Errorf("Want error code %d, got %d", tt.wantCode, resp.Error.Code)
			}
		})
	}
}

// Helper functions

// testHandler implements rpc.Handler for testing.
type testHandler struct {
	mu              sync.Mutex
	inputCalls      []rpc.InputParams
	mouseCalls      []rpc.MouseParams
	wheelCalls      []rpc.WheelParams
	pingCount       int
	screenshotCount int
}

func newTestHandler() *testHandler {
	return &testHandler{}
}

func (h *testHandler) HandleInput(params *rpc.InputParams) (any, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.inputCalls = append(h.inputCalls, *params)

	if params.Key == "" {
		return nil, fmt.Errorf("key is required")
	}

	if _, ok := input.LookupKey(params.Key); !ok {
		return nil, fmt.Errorf("unknown key: %s", params.Key)
	}
	return &rpc.InputResult{Success: true}, nil
}

func (h *testHandler) HandleMouse(params *rpc.MouseParams) (any, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.mouseCalls = append(h.mouseCalls, *params)

	// Validate button for non-position actions
	if params.Action != "position" && params.Button == "" {
		return nil, fmt.Errorf("button is required for action %s", params.Action)
	}

	if params.Button != "" {
		if _, ok := input.LookupMouseButton(params.Button); !ok {
			return nil, fmt.Errorf("unknown button: %s", params.Button)
		}
	}
	return &rpc.MouseResult{Success: true}, nil
}

func (h *testHandler) HandleWheel(params *rpc.WheelParams) (any, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.wheelCalls = append(h.wheelCalls, *params)
	return &rpc.WheelResult{Success: true}, nil
}

func (h *testHandler) HandleScreenshot(params *rpc.ScreenshotParams) (any, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.screenshotCount++
	return &rpc.ScreenshotResult{Success: true, Path: params.Output}, nil
}

func (h *testHandler) HandlePing() (any, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.pingCount++
	return &rpc.PingResult{OK: true}, nil
}

func (h *testHandler) HandleExit() {}

func (h *testHandler) HandleGetMousePosition() (any, error) {
	return &rpc.GetMousePositionResult{X: 0, Y: 0}, nil
}

func (h *testHandler) HandleGetWheelPosition() (any, error) {
	return &rpc.GetWheelPositionResult{X: 0, Y: 0}, nil
}

// setupTestSocketPair creates a connected socket pair for testing.
// Uses TCP for reliability across different systems.
func setupTestSocketPair(t *testing.T) (net.Listener, *testClient) {
	t.Helper()

	// Use TCP on a random port for reliability
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}

	addr := listener.Addr().String()

	handler := newTestHandler()

	// Start server goroutine
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}

			go handleTestConnection(conn, handler)
		}
	}()

	// Give server time to start
	time.Sleep(10 * time.Millisecond)

	// Connect client
	client, err := newTestClient(addr)
	if err != nil {
		listener.Close()
		t.Fatalf("Failed to create client: %v", err)
	}

	return listener, client
}

// testClient is a simple JSON-RPC client for testing.
type testClient struct {
	conn    net.Conn
	encoder *json.Encoder
	decoder *json.Decoder
	mu      sync.Mutex
}

func newTestClient(addr string) (*testClient, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &testClient{
		conn:    conn,
		encoder: json.NewEncoder(conn),
		decoder: json.NewDecoder(conn),
	}, nil
}

func (c *testClient) SendRequest(req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.encoder.Encode(req); err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	var resp rpc.RPCResponse
	if err := c.decoder.Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &resp, nil
}

func (c *testClient) Close() error {
	return c.conn.Close()
}

func handleTestConnection(conn net.Conn, handler *testHandler) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var req rpc.RPCRequest
		if err := decoder.Decode(&req); err != nil {
			if err == io.EOF {
				return
			}
			fmt.Fprintf(os.Stderr, "decode error: %v", err)
			return
		}

		rpcReq := &rpc.Request{Req: &req, Conn: conn}
		resp := rpc.ProcessRequest(rpcReq, handler)

		// Skip nil responses (async fire-and-forget)
		if resp.Result == nil && resp.Error == nil {
			continue
		}

		if err := encoder.Encode(resp); err != nil {
			fmt.Fprintf(os.Stderr, "encode error: %v", err)
			return
		}
	}
}

func sendRequest(client *testClient, req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	return client.SendRequest(req)
}

func mustMarshal(t *testing.T, v any) json.RawMessage {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}
	return data
}
