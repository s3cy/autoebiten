package server

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/s3cy/autoebiten/internal/custom"
	"github.com/s3cy/autoebiten/internal/rpc"
)

// HandleCustom handles custom command.
func (h *serverHandler) HandleCustom(params *rpc.CustomParams) (any, error) {
	if params.Name == "" {
		return nil, fmt.Errorf("command name is required")
	}

	handler := custom.Get(params.Name)
	if handler == nil {
		return nil, fmt.Errorf("unknown custom command: %s", params.Name)
	}

	// Create command context with deferred response capability
	ctx := custom.NewContext(params.Request, func(resp string) {
		result, _ := json.Marshal(&rpc.CustomResult{Response: resp})
		rpcResp := rpc.RPCResponse{
			JSONRPC: "2.0",
			ID:      params.ID,
			Result:  result,
		}

		// Send response
		if err := json.NewEncoder(params.Conn).Encode(rpcResp); err != nil {
			fmt.Fprintf(os.Stderr, "failed to send custom response: %v\n", err)
		}
	})

	// Execute the handler
	handler(ctx)

	return nil, nil
}

// HandleListCustomCommands returns a list of registered custom commands.
func (h *serverHandler) HandleListCustomCommands() (any, error) {
	return custom.List(), nil
}
