package server

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/s3cy/autoebiten/internal/input"
	"github.com/s3cy/autoebiten/internal/rpc"
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
	if params.Key == "" {
		return nil, fmt.Errorf("key is required")
	}

	key, ok := input.LookupKey(params.Key)
	if !ok {
		return nil, fmt.Errorf("unknown key: %s", params.Key)
	}

	it := input.NewInputTimeFromTick(Tick(), h.subtick.Add(1))

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

	return &rpc.InputResult{Success: true}, nil
}

// HandleMouse handles mouse command.
func (h *serverHandler) HandleMouse(params *rpc.MouseParams) (any, error) {
	if params.Action == "position" || params.X > 0 || params.Y > 0 {
		input.Get().InjectCursorMove(params.X, params.Y)
	}
	if params.Action == "position" {
		return &rpc.MouseResult{Success: true}, nil
	}

	if params.Button == "" {
		return nil, fmt.Errorf("button is required for action %s", params.Action)
	}

	btn, ok := input.LookupMouseButton(params.Button)
	if !ok {
		return nil, fmt.Errorf("unknown button: %s", params.Button)
	}

	it := input.NewInputTimeFromTick(Tick(), h.subtick.Add(1))

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

	return &rpc.MouseResult{Success: true}, nil
}

// HandleWheel handles wheel command.
func (h *serverHandler) HandleWheel(params *rpc.WheelParams) (any, error) {
	input.Get().InjectWheelMove(params.X, params.Y)
	return &rpc.WheelResult{Success: true}, nil
}

// HandleScreenshot handles screenshot command.
func (h *serverHandler) HandleScreenshot(params *rpc.ScreenshotParams) (any, error) {
	path := params.Output

	if path == "" {
		path = fmt.Sprintf("screenshot_%s.png", time.Now().Format("20060102150405"))
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	if !params.Async {
		// Sync mode: queue request. ProcessScreenshots in Draw() will send the response
		var conn net.Conn
		if params.Conn != nil {
			conn = params.Conn
		}
		queueScreenshot(params.ID, path, conn)
		return nil, nil
	}

	// Async mode: fire-and-forget, no connection needed
	queueScreenshot(params.ID, path, nil)
	return &rpc.ScreenshotResult{Success: true, Path: path}, nil
}

// HandlePing handles ping command.
func (h *serverHandler) HandlePing() (any, error) {
	return &rpc.PingResult{OK: true}, nil
}

// HandleExit handles exit command.
func (h *serverHandler) HandleExit() {
	exitingMu.Lock()
	exiting = true
	exitingMu.Unlock()
}

var _ rpc.Handler = (*serverHandler)(nil)
