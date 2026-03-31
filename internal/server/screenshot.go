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
	for _, req := range queue {
		// No more pitfalls. Yes, Go 1.22
		go respondScreenshot(req, data, err)
	}
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

func respondScreenshot(params *rpc.ScreenshotParams, data []byte, err error) {
	if params.Async {
		// Async fire-and-forget, just save to file if path provided
		if params.Output != "" && err == nil {
			os.WriteFile(params.Output, data, 0644)
		}
		return
	}

	// Sync mode - send response via connection
	rpcResp := rpc.RPCResponse{
		JSONRPC: "2.0",
		ID:      params.ID,
	}

	if err != nil {
		rpcResp.Error = &rpc.RPCError{
			Code:    rpc.ErrScreenshotFailed,
			Message: err.Error(),
		}
	} else {
		var result rpc.ScreenshotResult
		result.Success = true
		if data != nil && params.Base64 {
			result.Data = base64.StdEncoding.EncodeToString(data)
		}
		if params.Output != "" {
			result.Path = params.Output
			os.WriteFile(params.Output, data, 0644)
		}
		resultJSON, _ := json.Marshal(result)
		rpcResp.Result = resultJSON
	}

	// Send response and signal completion
	if err := json.NewEncoder(params.Conn).Encode(rpcResp); err != nil {
		fmt.Fprintf(os.Stderr, "failed to send screenshot response: %v\n", err)
	}
}
