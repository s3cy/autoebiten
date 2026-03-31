package server

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/s3cy/autoebiten/internal/input"
	"github.com/s3cy/autoebiten/internal/rpc"
)

var (
	inputMu    sync.Mutex
	inputQueue []*rpc.InputParams
)

func queueInputResult(params *rpc.InputParams) {
	inputMu.Lock()
	defer inputMu.Unlock()
	inputQueue = append(inputQueue, params)
}

func processInputResults() {
	// Collect and clear queue
	inputMu.Lock()
	if len(inputQueue) == 0 {
		inputMu.Unlock()
		return
	}
	queue := inputQueue
	inputQueue = nil
	inputMu.Unlock()

	// Process inputs and send responses
	for _, req := range queue {
		if req.Conn != nil {
			// Sync mode - send response via connection
			rpcResp := rpc.RPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
			}

			rpcResp.Result = &rpc.InputResult{Success: true}

			// Send response
			if err := json.NewEncoder(req.Conn).Encode(rpcResp); err != nil {
				fmt.Fprintf(os.Stderr, "failed to send input response: %v\n", err)
			}
		}
	}
}

// processInputRequest processes an input request and returns the result.
// If async is false, it queues the result. Update() will call processInputs to send the response in the next tick.
func processInputRequest(params *rpc.InputParams) (any, error) {
	if params.Key == "" {
		return nil, fmt.Errorf("key is required")
	}

	key, ok := input.LookupKey(params.Key)
	if !ok {
		return nil, fmt.Errorf("unknown key: %s", params.Key)
	}

	it := input.NewInputTimeFromTick(Tick(), globalServerHandler.subtick.Add(1))

	switch params.Action {
	case "press":
		input.Get().InjectKeyPress(key, it)
	case "release":
		input.Get().InjectKeyRelease(key, it)
	case "hold":
		duration := max(params.DurationTicks, 1)
		input.Get().InjectKeyHold(key, it, duration)
	default:
		return nil, fmt.Errorf("unknown action: %s", params.Action)
	}

	if !params.Async {
		queueInputResult(params)
		return nil, nil
	}

	return &rpc.InputResult{Success: true}, nil
}
