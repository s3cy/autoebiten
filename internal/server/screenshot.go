package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"os"
	"sync"

	"github.com/s3cy/autoebiten/internal/rpc"
)

var (
	screenshotMu    sync.Mutex
	screenshotQueue []*rpc.ScreenshotParams
)

// queueScreenshot queues a screenshot request to be processed in Draw.
// Call this from Update(). For sync mode, pass conn and done channel.
func queueScreenshot(params *rpc.ScreenshotParams) {
	screenshotMu.Lock()
	defer screenshotMu.Unlock()
	screenshotQueue = append(screenshotQueue, params)
}

// ProcessScreenshots processes pending screenshot requests.
// Call this from the game's Draw function.
func ProcessScreenshots(screen image.Image) {
	// Collect and clear queue
	screenshotMu.Lock()
	if len(screenshotQueue) == 0 {
		screenshotMu.Unlock()
		return
	}
	queue := screenshotQueue
	screenshotQueue = nil
	screenshotMu.Unlock()

	data, err := captureScreen(screen)
	go respondScreenshots(queue, data, err)
}

// captureScreen captures the current screen to PNG format.
func captureScreen(screen image.Image) ([]byte, error) {
	if screen == nil {
		return nil, nil
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, screen); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func respondScreenshots(queue []*rpc.ScreenshotParams, data []byte, err error) {
	for _, req := range queue {
		if req.Conn == nil {
			// Async fire-and-forget, just save to file if path provided
			if req.Output != "" && err == nil {
				os.WriteFile(req.Output, data, 0644)
			}
			continue
		}

		// Sync mode - send response via connection
		rpcResp := rpc.RPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
		}

		if err != nil {
			rpcResp.Error = &rpc.RPCError{
				Code:    rpc.ErrScreenshotFailed,
				Message: err.Error(),
			}
		} else {
			result := map[string]any{"success": true}
			if req.Output != "" {
				result["path"] = req.Output
				os.WriteFile(req.Output, data, 0644)
			} else if data != nil {
				result["data"] = base64.StdEncoding.EncodeToString(data)
			}
			rpcResp.Result = result
		}

		// Send response and signal completion
		if err := json.NewEncoder(req.Conn).Encode(rpcResp); err != nil {
			fmt.Fprintf(os.Stderr, "failed to send screenshot response: %v\n", err)
		}
	}
}
