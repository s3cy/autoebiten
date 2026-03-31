package rpc

import (
	"encoding/json"
	"net"
)

// RPCRequest represents a JSON-RPC request.
type RPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// RPCResponse represents a JSON-RPC response.
type RPCResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      any       `json:"id"`
	Result  any       `json:"result,omitempty"`
	Error   *RPCError `json:"error,omitempty"`
}

// RPCError represents a JSON-RPC error.
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error codes.
const (
	ErrConnectionFailed = -32000
	ErrInvalidParams    = -32001
	ErrScriptFailed     = -32002
	ErrScreenshotFailed = -32003
	ErrGameNotRunning   = -32004
)

// InputParams represents input command parameters.
type InputParams struct {
	Action        string `json:"action"`
	Key           string `json:"key"`
	DurationTicks int64  `json:"duration_ticks"`
	Async         bool   `json:"async"`
	// ID is an internal field used by the server to track the request.
	// This is not part of the JSON-RPC protocol.
	ID any `json:"-"`
	// Conn is an internal field used by the server to send responses.
	// This is not part of the JSON-RPC protocol.
	Conn net.Conn `json:"-"`
}

// MouseParams represents mouse command parameters.
type MouseParams struct {
	Action        string `json:"action"`
	X             int    `json:"x"`
	Y             int    `json:"y"`
	Button        string `json:"button"`
	DurationTicks int64  `json:"duration_ticks"`
	Async         bool   `json:"async"`
	// ID is an internal field used by the server to track the request.
	// This is not part of the JSON-RPC protocol.
	ID any `json:"-"`
	// Conn is an internal field used by the server to send responses.
	// This is not part of the JSON-RPC protocol.
	Conn net.Conn `json:"-"`
}

// WheelParams represents wheel command parameters.
type WheelParams struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Async bool `json:"async"`
	// ID is an internal field used by the server to track the request.
	// This is not part of the JSON-RPC protocol.
	ID any `json:"-"`
	// Conn is an internal field used by the server to send responses.
	// This is not part of the JSON-RPC protocol.
	Conn net.Conn `json:"-"`
}

// ScreenshotParams represents screenshot command parameters.
type ScreenshotParams struct {
	Output string `json:"output"`
	Async  bool   `json:"sync"`
	// ID is an internal field used by the server to track the request.
	// This is not part of the JSON-RPC protocol.
	ID any `json:"-"`
	// Conn is an internal field used by the server to send responses.
	// This is not part of the JSON-RPC protocol.
	Conn net.Conn `json:"-"`
}

// InputResult represents input command result.
type InputResult struct {
	Success bool `json:"success"`
}

// MouseResult represents mouse command result.
type MouseResult struct {
	Success bool `json:"success"`
}

// WheelResult represents wheel command result.
type WheelResult struct {
	Success bool `json:"success"`
}

// ScreenshotResult represents screenshot command result.
type ScreenshotResult struct {
	Success bool   `json:"success"`
	Path    string `json:"path,omitempty"`
	Data    string `json:"data,omitempty"`
}

// PingResult represents ping command result.
type PingResult struct {
	OK bool `json:"ok"`
}

// GetMousePositionResult represents get_mouse_position command result.
type GetMousePositionResult struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// GetWheelPositionResult represents get_wheel_position command result.
type GetWheelPositionResult struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
