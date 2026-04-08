package proxy

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/s3cy/autoebiten/internal/output"
	"github.com/s3cy/autoebiten/internal/rpc"
)

// ProxyState represents the state of the proxy handler.
type ProxyState int

const (
	StateWaiting ProxyState = iota   // Waiting for game RPC connection
	StateConnected                   // Game connected, proxy active
	StateCrashed                     // Game crashed/exited
)

// String returns the string representation of the state.
func (s ProxyState) String() string {
	switch s {
	case StateWaiting:
		return "waiting"
	case StateConnected:
		return "connected"
	case StateCrashed:
		return "crashed"
	default:
		return "unknown"
	}
}

// GameClient is the interface for game RPC clients.
type GameClient interface {
	SendRequest(req *rpc.RPCRequest) (*rpc.RPCResponse, error)
	Close() error
}

// UnifiedHandler handles all proxy RPC requests with state machine.
type UnifiedHandler struct {
	state       ProxyState
	gameClient  GameClient
	outputMgr   *output.OutputManager
	launchError error
	onCrashed   func()
	mu          sync.Mutex
}

// NewUnifiedHandler creates a new unified handler starting in Waiting state.
func NewUnifiedHandler(outputMgr *output.OutputManager, onCrashed func()) *UnifiedHandler {
	return &UnifiedHandler{
		state:     StateWaiting,
		outputMgr: outputMgr,
		onCrashed: onCrashed,
	}
}

// TransitionToConnected transitions the handler to Connected state.
func (h *UnifiedHandler) TransitionToConnected(gameClient GameClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.state = StateConnected
	h.gameClient = gameClient
}

// TransitionToCrashed transitions the handler to Crashed state with an error.
func (h *UnifiedHandler) TransitionToCrashed(err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.state = StateCrashed
	if h.launchError == nil {
		h.launchError = err
	} else {
		h.launchError = fmt.Errorf("%v; %v", h.launchError, err)
	}
}

// GetState returns the current state (for testing).
func (h *UnifiedHandler) GetState() ProxyState {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.state
}

// ProcessRequest processes an RPC request based on current state.
func (h *UnifiedHandler) ProcessRequest(conn net.Conn, req *rpc.RPCRequest) {
	h.mu.Lock()
	defer h.mu.Unlock()

	switch h.state {
	case StateWaiting:
		h.handleWaitingRequest(conn, req)
	case StateConnected:
		h.handleConnectedRequest(conn, req)
	case StateCrashed:
		h.handleCrashedRequest(conn, req)
	}
}

func (h *UnifiedHandler) handleWaitingRequest(conn net.Conn, req *rpc.RPCRequest) {
	// Generate diff (always do this)
	diff, err := h.outputMgr.DiffAndUpdateSnapshot()
	if err != nil {
		diff = ""
	}

	// Build error response
	resp := rpc.ErrorResponse(req.ID, rpc.ErrGameNotConnected, "game not connected")
	resp.Extra = map[string]any{
		"diff":        diff,
		"proxy_error": "", // Empty in waiting state
	}

	h.sendResponse(conn, &resp)
}

func (h *UnifiedHandler) handleConnectedRequest(conn net.Conn, req *rpc.RPCRequest) {
	// Generate diff before forwarding
	diff, err := h.outputMgr.DiffAndUpdateSnapshot()
	if err != nil {
		diff = ""
	}

	// Forward to game
	gameResp, err := h.gameClient.SendRequest(req)
	if err != nil {
		// Game request failed - transition to crashed
		h.state = StateCrashed
		h.launchError = fmt.Errorf("game request failed: %w", err)

		resp := rpc.ErrorResponse(req.ID, rpc.ErrGameNotConnected, "game not connected")
		resp.Extra = map[string]any{
			"diff":        diff,
			"proxy_error": h.launchError.Error(),
		}
		h.sendResponse(conn, &resp)
		return
	}

	// Add diff to response
	if gameResp.Extra == nil {
		gameResp.Extra = make(map[string]any)
	}
	gameResp.Extra["diff"] = diff

	h.sendResponse(conn, gameResp)
}

func (h *UnifiedHandler) handleCrashedRequest(conn net.Conn, req *rpc.RPCRequest) {
	// Generate diff
	diff, err := h.outputMgr.DiffAndUpdateSnapshot()
	if err != nil {
		diff = ""
	}

	// Build error response with accumulated error
	proxyErr := ""
	if h.launchError != nil {
		proxyErr = h.launchError.Error()
	}

	resp := rpc.ErrorResponse(req.ID, rpc.ErrGameNotConnected, "game not connected")
	resp.Extra = map[string]any{
		"diff":        diff,
		"proxy_error": proxyErr,
	}

	// Signal that CLI has queried us (trigger exit)
	if h.onCrashed != nil {
		go h.onCrashed()
	}

	h.sendResponse(conn, &resp)
}

func (h *UnifiedHandler) sendResponse(conn net.Conn, resp *rpc.RPCResponse) {
	encoder := json.NewEncoder(conn)
	encoder.Encode(resp)
}

// ==================== Legacy Server and Handler (for backward compatibility) ====================

// Server wraps game RPC calls and captures output.
type Server struct {
	gameClient  GameClient
	outputMgr   *output.OutputManager
	outputFiles *output.FilePath
	mu          sync.Mutex
}

// NewServer creates a new proxy server.
func NewServer(gameClient GameClient, outputMgr *output.OutputManager, outputFiles *output.FilePath) *Server {
	return &Server{
		gameClient:  gameClient,
		outputMgr:   outputMgr,
		outputFiles: outputFiles,
	}
}

// Close closes the game client connection.
func (s *Server) Close() error {
	if s.gameClient != nil {
		return s.gameClient.Close()
	}
	return nil
}

// Cleanup removes log and snapshot files.
func (s *Server) Cleanup() error {
	var errs []error
	if err := os.Remove(s.outputFiles.Log); err != nil && !os.IsNotExist(err) {
		errs = append(errs, err)
	}
	if err := os.Remove(s.outputFiles.Snapshot); err != nil && !os.IsNotExist(err) {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return fmt.Errorf("cleanup errors: %v", errs)
	}
	return nil
}

// ForwardRequest forwards an RPC request to the game and returns a wrapped response with output.
// This is the main proxy method - it doesn't know or care about specific RPC methods.
func (s *Server) ForwardRequest(req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Forward request to game
	resp, err := s.gameClient.SendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("game request failed: %w", err)
	}

	// Generate diff and update snapshot
	diff, err := s.outputMgr.DiffAndUpdateSnapshot()
	if err != nil {
		return nil, fmt.Errorf("failed to generate diff: %w", err)
	}

	// Build proxy response
	if resp.Extra == nil {
		resp.Extra = make(map[string]any)
	}
	resp.Extra["diff"] = diff

	return resp, nil
}

// Handler handles incoming proxy RPC requests.
type Handler struct {
	server *Server
}

// NewHandler creates a new proxy handler.
func NewHandler(server *Server) *Handler {
	return &Handler{server: server}
}

// ProcessRequest processes a proxy RPC request and sends the response.
// The proxy is transparent - it forwards any method to the game and wraps the response.
func (h *Handler) ProcessRequest(conn net.Conn, req *rpc.RPCRequest) {
	// Forward to game and capture output
	resp, err := h.server.ForwardRequest(req)
	if err != nil {
		sendErrorResponse(conn, req.ID, rpc.ErrInvalidParams, err.Error())
		return
	}

	// Send response back to CLI
	sendResponse(conn, resp)
}

func sendErrorResponse(conn net.Conn, id any, code int, message string) {
	resp := rpc.ErrorResponse(id, code, message)
	sendResponse(conn, &resp)
}

func sendResponse(conn net.Conn, resp *rpc.RPCResponse) {
	encoder := json.NewEncoder(conn)
	encoder.Encode(resp)
}
