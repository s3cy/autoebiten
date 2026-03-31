package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"net"
	"os"
	"sync"

	"github.com/s3cy/autoebiten/internal/rpc"
)

// ScreenshotRequest represents a pending screenshot request.
type ScreenshotRequest struct {
	ID   any      // RPC request ID for response
	Path string   // output path, empty for base64 return
	Conn net.Conn // connection to send response on (nil for async)
}

var (
	screenshotMu    sync.Mutex
	screenshotQueue []ScreenshotRequest
)

// queueScreenshot queues a screenshot request to be processed in Draw.
// Call this from Update(). For sync mode, pass conn and done channel.
func queueScreenshot(id any, path string, conn net.Conn) {
	screenshotMu.Lock()
	defer screenshotMu.Unlock()
	screenshotQueue = append(screenshotQueue, ScreenshotRequest{
		ID:   id,
		Path: path,
		Conn: conn,
	})
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

func respondScreenshots(queue []ScreenshotRequest, data []byte, err error) {
	for _, req := range queue {
		if req.Conn == nil {
			// Async fire-and-forget, just save to file if path provided
			if req.Path != "" && err == nil {
				os.WriteFile(req.Path, data, 0644)
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
			if req.Path != "" {
				result["path"] = req.Path
				os.WriteFile(req.Path, data, 0644)
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
