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
	wheelMu    sync.Mutex
	wheelQueue []*rpc.WheelParams
)

func queueWheelResult(params *rpc.WheelParams) {
	wheelMu.Lock()
	defer wheelMu.Unlock()
	wheelQueue = append(wheelQueue, params)
}

func processWheelResults() {
	// Collect and clear queue
	wheelMu.Lock()
	if len(wheelQueue) == 0 {
		wheelMu.Unlock()
		return
	}
	queue := wheelQueue
	wheelQueue = nil
	wheelMu.Unlock()

	// Process wheel results and send responses
	for _, req := range queue {
		// Sync mode - send response via connection
		rpcResp := rpc.RPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  &rpc.WheelResult{Success: true},
		}

		go func() {
			// Send response
			if err := json.NewEncoder(req.Conn).Encode(rpcResp); err != nil {
				fmt.Fprintf(os.Stderr, "failed to send wheel response: %v\n", err)
			}
		}()
	}
}

// processWheelRequest processes a wheel request and returns the result.
// If async is false, it queues the result. Update() will call processWheelResults to send the response in the next tick.
func processWheelRequest(params *rpc.WheelParams) (any, error) {
	input.Get().InjectWheelMove(params.X, params.Y)

	if !params.Async {
		queueWheelResult(params)
		return nil, nil
	}

	return &rpc.WheelResult{Success: true}, nil
}
