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
	mouseMu    sync.Mutex
	mouseQueue []*rpc.MouseParams
)

func queueMouseResult(params *rpc.MouseParams) {
	mouseMu.Lock()
	defer mouseMu.Unlock()
	mouseQueue = append(mouseQueue, params)
}

func processMouseResults() {
	// Collect and clear queue
	mouseMu.Lock()
	if len(mouseQueue) == 0 {
		mouseMu.Unlock()
		return
	}
	queue := mouseQueue
	mouseQueue = nil
	mouseMu.Unlock()

	// Process mouse results and send responses
	for _, req := range queue {
		// Sync mode - send response via connection
		rpcResp := rpc.RPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  &rpc.MouseResult{Success: true},
		}

		go func() {
			// Send response
			if err := json.NewEncoder(req.Conn).Encode(rpcResp); err != nil {
				fmt.Fprintf(os.Stderr, "failed to send mouse response: %v\n", err)
			}
		}()
	}
}

// processMouseRequest processes a mouse request and returns the result.
// If async is false, it queues the result. Update() will call processMouseResults to send the response in the next tick.
func processMouseRequest(params *rpc.MouseParams) (any, error) {
	if params.Action == "position" {
		input.Get().InjectCursorMove(params.X, params.Y)

		if !params.Async {
			queueMouseResult(params)
			return nil, nil
		}
		return &rpc.MouseResult{Success: true}, nil
	}

	if params.Button == "" {
		return nil, fmt.Errorf("button is required for action %s", params.Action)
	}

	btn, ok := input.LookupMouseButton(params.Button)
	if !ok {
		return nil, fmt.Errorf("unknown button: %s", params.Button)
	}

	if params.X != 0 || params.Y != 0 {
		input.Get().InjectCursorMove(params.X, params.Y)
	}

	it := input.NewInputTimeFromTick(Tick(), globalServerHandler.subtick.Add(1))

	switch params.Action {
	case "press":
		input.Get().InjectMouseButtonPress(btn, it)
	case "release":
		input.Get().InjectMouseButtonRelease(btn, it)
	case "hold":
		duration := max(params.DurationTicks, 1)
		input.Get().InjectMouseButtonHold(btn, it, duration)
	default:
		return nil, fmt.Errorf("unknown action: %s", params.Action)
	}

	if !params.Async {
		queueMouseResult(params)
		return nil, nil
	}

	return &rpc.MouseResult{Success: true}, nil
}
