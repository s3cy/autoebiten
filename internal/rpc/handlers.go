package rpc

import (
	"encoding/json"
	"fmt"
)

// Handler is the interface for RPC command handlers.
type Handler interface {
	HandleInput(params *InputParams) (any, error)
	HandleMouse(params *MouseParams) (any, error)
	HandleWheel(params *WheelParams) (any, error)
	HandleScreenshot(params *ScreenshotParams) (any, error)
	HandlePing() (any, error)
	HandleGetMousePosition() (any, error)
	HandleGetWheelPosition() (any, error)
	HandleExit()
}

// DefaultHandler is a default no-op handler for testing.
type DefaultHandler struct{}

// NewDefaultHandler creates a new DefaultHandler.
func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{}
}

// HandleInput handles input command.
func (h *DefaultHandler) HandleInput(params *InputParams) (any, error) {
	return &InputResult{Success: true}, nil
}

// HandleMouse handles mouse command.
func (h *DefaultHandler) HandleMouse(params *MouseParams) (any, error) {
	return &MouseResult{Success: true}, nil
}

// HandleWheel handles wheel command.
func (h *DefaultHandler) HandleWheel(params *WheelParams) (any, error) {
	return &WheelResult{Success: true}, nil
}

// HandleScreenshot handles screenshot command.
func (h *DefaultHandler) HandleScreenshot(params *ScreenshotParams) (any, error) {
	return &ScreenshotResult{Success: true, Path: params.Output}, nil
}

// HandlePing handles ping command.
func (h *DefaultHandler) HandlePing() (any, error) {
	return &PingResult{OK: true}, nil
}

// HandleGetMousePosition handles get_mouse_position command.
func (h *DefaultHandler) HandleGetMousePosition() (any, error) {
	return &GetMousePositionResult{X: 0, Y: 0}, nil
}

// HandleGetWheelPosition handles get_wheel_position command.
func (h *DefaultHandler) HandleGetWheelPosition() (any, error) {
	return &GetWheelPositionResult{X: 0, Y: 0}, nil
}

// HandleExit handles exit command.
func (h *DefaultHandler) HandleExit() {}

// ProcessRequest processes an RPC request and returns the response.
func ProcessRequest(req *Request, handler Handler) RPCResponse {
	id := req.Req.ID

	switch req.Req.Method {
	case "ping":
		result, err := handler.HandlePing()
		if err != nil {
			return errorResponse(id, ErrInvalidParams, err.Error())
		}
		return RPCResponse{JSONRPC: "2.0", ID: id, Result: result}

	case "get_mouse_position":
		result, err := handler.HandleGetMousePosition()
		if err != nil {
			return errorResponse(id, ErrInvalidParams, err.Error())
		}
		return RPCResponse{JSONRPC: "2.0", ID: id, Result: result}

	case "get_wheel_position":
		result, err := handler.HandleGetWheelPosition()
		if err != nil {
			return errorResponse(id, ErrInvalidParams, err.Error())
		}
		return RPCResponse{JSONRPC: "2.0", ID: id, Result: result}

	case "exit":
		handler.HandleExit()
		return RPCResponse{JSONRPC: "2.0", ID: id, Result: map[string]bool{"success": true}}

	case "input":
		var params InputParams
		if err := decodeParams(req.Req.Params, &params); err != nil {
			return errorResponse(id, ErrInvalidParams, fmt.Sprintf("invalid params: %v", err))
		}
		result, err := handler.HandleInput(&params)
		if err != nil {
			return errorResponse(id, ErrInvalidParams, err.Error())
		}
		return RPCResponse{JSONRPC: "2.0", ID: id, Result: result}

	case "mouse":
		var params MouseParams
		if err := decodeParams(req.Req.Params, &params); err != nil {
			return errorResponse(id, ErrInvalidParams, fmt.Sprintf("invalid params: %v", err))
		}
		result, err := handler.HandleMouse(&params)
		if err != nil {
			return errorResponse(id, ErrInvalidParams, err.Error())
		}
		return RPCResponse{JSONRPC: "2.0", ID: id, Result: result}

	case "wheel":
		var params WheelParams
		if err := decodeParams(req.Req.Params, &params); err != nil {
			return errorResponse(id, ErrInvalidParams, fmt.Sprintf("invalid params: %v", err))
		}
		result, err := handler.HandleWheel(&params)
		if err != nil {
			return errorResponse(id, ErrInvalidParams, err.Error())
		}
		return RPCResponse{JSONRPC: "2.0", ID: id, Result: result}

	case "screenshot":
		var params ScreenshotParams
		if err := decodeParams(req.Req.Params, &params); err != nil {
			// Screenshot params are optional, continue with empty
			params = ScreenshotParams{}
		}
		params.ID = req.Req.ID // Track request ID for response routing
		params.Conn = req.Conn // Store connection for response
		result, err := handler.HandleScreenshot(&params)
		if err != nil {
			return errorResponse(id, ErrScreenshotFailed, err.Error())
		}
		// For sync screenshots, the response is sent by ProcessScreenshots in Draw()
		// so we return nil result to indicate no immediate response
		if result == nil {
			return RPCResponse{}
		}
		return RPCResponse{JSONRPC: "2.0", ID: id, Result: result}

	default:
		return errorResponse(id, ErrInvalidParams, fmt.Sprintf("unknown method: %s", req.Req.Method))
	}
}

func decodeParams(data json.RawMessage, v any) error {
	return json.Unmarshal(data, v)
}

func errorResponse(id any, code int, message string) RPCResponse {
	return RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
		},
	}
}
