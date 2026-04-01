package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/s3cy/autoebiten/internal/input"
	"github.com/s3cy/autoebiten/internal/rpc"
	"github.com/s3cy/autoebiten/internal/version"
)

var (
	exiting   bool
	exitingMu sync.RWMutex

	// socketServerOnce ensures the socket server is started only once
	socketServerOnce sync.Once

	// reqChan receives requests from the socket server
	reqChan <-chan *rpc.Request

	globalServerHandler = &serverHandler{}
)

// Update processes RPC commands from the socket.
// Call this in your game loop's Update function.
//
// Returns false if the game should exit (exit command received).
func Update() bool {
	processInputResults()
	processMouseResults()
	processWheelResults()

	globalServerHandler.incrementTick()

	// Ensure socket server is started (one-time initialization)
	socketServerOnce.Do(startSocketServer)

	// Check if we should exit
	exitingMu.RLock()
	if exiting {
		exitingMu.RUnlock()
		return false
	}
	exitingMu.RUnlock()

	// Process any pending request synchronously (non-blocking)
Loop:
	for {
		select {
		case req, ok := <-reqChan:
			if !ok {
				// Channel closed, reset and let next call restart
				reqChan = nil
				socketServerOnce = sync.Once{}
			} else {
				processRequest(req)
			}
		default:
			// No pending requests
			break Loop
		}
	}

	return true
}

func Tick() int64 {
	return globalServerHandler.tick.Load()
}

func processRequest(req *rpc.Request) {
	resp := rpc.ProcessRequest(req, globalServerHandler)
	// If result is nil (async fire-and-forget), don't send a response
	if resp.Result == nil && resp.Error == nil {
		return
	}
	go func() {
		encoder := json.NewEncoder(req.Conn)
		if err := encoder.Encode(resp); err != nil {
			fmt.Fprintf(os.Stderr, "autoebiten: failed to encode response: %v\n", err)
		}
	}()
}

// startSocketServer starts the socket server and stores the request channel.
func startSocketServer() {
	var err error
	reqChan, err = rpc.Serve()
	if err != nil {
		fmt.Fprintf(os.Stderr, "autoebiten: socket server error: %v\n", err)
	}
}

// serverHandler implements rpc.Handler for socket connections.
type serverHandler struct {
	tick    atomic.Int64
	subtick atomic.Int64
}

func (h *serverHandler) incrementTick() {
	h.tick.Add(1)
	h.subtick.Store(0)
}

// HandleInput handles input command.
func (h *serverHandler) HandleInput(params *rpc.InputParams) (any, error) {
	return processInputRequest(params)
}

// HandleMouse handles mouse command.
func (h *serverHandler) HandleMouse(params *rpc.MouseParams) (any, error) {
	return processMouseRequest(params)
}

// HandleWheel handles wheel command.
func (h *serverHandler) HandleWheel(params *rpc.WheelParams) (any, error) {
	return processWheelRequest(params)
}

// HandleScreenshot handles screenshot command.
func (h *serverHandler) HandleScreenshot(params *rpc.ScreenshotParams) (any, error) {
	if !params.Base64 {
		if params.Output == "" {
			params.Output = fmt.Sprintf("screenshot_%s.png", time.Now().Format("20060102150405"))
		}

		var err error
		params.Output, err = filepath.Abs(params.Output)
		if err != nil {
			return nil, fmt.Errorf("invalid path: %w", err)
		}
	}

	if !params.Async {
		// Sync mode: queue request. ProcessScreenshots in Draw() will send the response
		queueScreenshot(params)
		return nil, nil
	}

	// Async mode: fire-and-forget, no connection needed
	queueScreenshot(params)
	return &rpc.ScreenshotResult{Success: true, Path: params.Output}, nil
}

// HandlePing handles ping command.
func (h *serverHandler) HandlePing() (any, error) {
	return &rpc.PingResult{OK: true, Version: version.VersionForLibrary()}, nil
}

// HandleGetMousePosition handles get_mouse_position command.
func (h *serverHandler) HandleGetMousePosition() (any, error) {
	x, y := input.Get().CursorPosition()
	return &rpc.GetMousePositionResult{X: x, Y: y}, nil
}

// HandleGetWheelPosition handles get_wheel_position command.
func (h *serverHandler) HandleGetWheelPosition() (any, error) {
	x, y := input.Get().Wheel()
	return &rpc.GetWheelPositionResult{X: x, Y: y}, nil
}

// HandleExit handles exit command.
func (h *serverHandler) HandleExit() {
	exitingMu.Lock()
	exiting = true
	exitingMu.Unlock()
}

var _ rpc.Handler = (*serverHandler)(nil)
