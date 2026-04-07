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
	HandleCustom(params *CustomParams) (any, error)
	HandleListCustomCommands() (any, error)
}

// marshalResult marshals a result value to json.RawMessage.
func marshalResult(result any) json.RawMessage {
	if result == nil {
		return nil
	}
	data, err := json.Marshal(result)
	if err != nil {
		return nil
	}
	return data
}

// ProcessRequest processes an RPC request and returns the response.
func ProcessRequest(req *Request, handler Handler) RPCResponse {
	id := req.Req.ID

	switch req.Req.Method {
	case "ping":
		result, err := handler.HandlePing()
		if err != nil {
			return ErrorResponse(id, ErrInvalidParams, err.Error())
		}
		return RPCResponse{JSONRPC: "2.0", ID: id, Result: marshalResult(result)}

	case "get_mouse_position":
		result, err := handler.HandleGetMousePosition()
		if err != nil {
			return ErrorResponse(id, ErrInvalidParams, err.Error())
		}
		return RPCResponse{JSONRPC: "2.0", ID: id, Result: marshalResult(result)}

	case "get_wheel_position":
		result, err := handler.HandleGetWheelPosition()
		if err != nil {
			return ErrorResponse(id, ErrInvalidParams, err.Error())
		}
		return RPCResponse{JSONRPC: "2.0", ID: id, Result: marshalResult(result)}

	case "exit":
		handler.HandleExit()
		return RPCResponse{JSONRPC: "2.0", ID: id, Result: marshalResult(map[string]bool{"success": true})}

	case "custom":
		var params CustomParams
		if err := decodeParams(req.Req.Params, &params); err != nil {
			return ErrorResponse(id, ErrInvalidParams, fmt.Sprintf("invalid params: %v", err))
		}
		params.ID = req.Req.ID // Track request ID for response routing
		params.Conn = req.Conn // Store connection for response
		result, err := handler.HandleCustom(&params)
		if err != nil {
			return ErrorResponse(id, ErrInvalidParams, err.Error())
		}
		return RPCResponse{JSONRPC: "2.0", ID: id, Result: marshalResult(result)}

	case "list_custom_commands":
		result, err := handler.HandleListCustomCommands()
		if err != nil {
			return ErrorResponse(id, ErrInvalidParams, err.Error())
		}
		return RPCResponse{JSONRPC: "2.0", ID: id, Result: marshalResult(result)}

	case "input":
		var params InputParams
		if err := decodeParams(req.Req.Params, &params); err != nil {
			return ErrorResponse(id, ErrInvalidParams, fmt.Sprintf("invalid params: %v", err))
		}
		params.ID = req.Req.ID // Track request ID for response routing
		params.Conn = req.Conn // Store connection for response
		result, err := handler.HandleInput(&params)
		if err != nil {
			return ErrorResponse(id, ErrInvalidParams, err.Error())
		}
		// For sync inputs, the response is sent by ProcessInputs in Update()
		// so we return nil result to indicate no immediate response
		if result == nil {
			return RPCResponse{}
		}
		return RPCResponse{JSONRPC: "2.0", ID: id, Result: marshalResult(result)}

	case "mouse":
		var params MouseParams
		if err := decodeParams(req.Req.Params, &params); err != nil {
			return ErrorResponse(id, ErrInvalidParams, fmt.Sprintf("invalid params: %v", err))
		}
		params.ID = req.Req.ID // Track request ID for response routing
		params.Conn = req.Conn // Store connection for response
		result, err := handler.HandleMouse(&params)
		if err != nil {
			return ErrorResponse(id, ErrInvalidParams, err.Error())
		}
		// For sync mouse commands, the response is sent by processMouseResults in Update()
		// so we return nil result to indicate no immediate response
		if result == nil {
			return RPCResponse{}
		}
		return RPCResponse{JSONRPC: "2.0", ID: id, Result: marshalResult(result)}

	case "wheel":
		var params WheelParams
		if err := decodeParams(req.Req.Params, &params); err != nil {
			return ErrorResponse(id, ErrInvalidParams, fmt.Sprintf("invalid params: %v", err))
		}
		params.ID = req.Req.ID // Track request ID for response routing
		params.Conn = req.Conn // Store connection for response
		result, err := handler.HandleWheel(&params)
		if err != nil {
			return ErrorResponse(id, ErrInvalidParams, err.Error())
		}
		// For sync wheel commands, the response is sent by processWheelResults in Update()
		// so we return nil result to indicate no immediate response
		if result == nil {
			return RPCResponse{}
		}
		return RPCResponse{JSONRPC: "2.0", ID: id, Result: marshalResult(result)}

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
			return ErrorResponse(id, ErrScreenshotFailed, err.Error())
		}
		// For sync screenshots, the response is sent by ProcessScreenshots in Draw()
		// so we return nil result to indicate no immediate response
		if result == nil {
			return RPCResponse{}
		}
		return RPCResponse{JSONRPC: "2.0", ID: id, Result: marshalResult(result)}

	default:
		return ErrorResponse(id, ErrInvalidParams, fmt.Sprintf("unknown method: %s", req.Req.Method))
	}
}

func decodeParams(data json.RawMessage, v any) error {
	return json.Unmarshal(data, v)
}

func ErrorResponse(id any, code int, message string) RPCResponse {
	return RPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
		},
	}
}
