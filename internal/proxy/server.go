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

// Server wraps game RPC calls and captures output.
type Server struct {
	gameClient  GameClient
	outputMgr   *output.OutputManager
	outputFiles *output.FilePath
	mu          sync.Mutex
}

// GameClient is the interface for game RPC clients.
type GameClient interface {
	SendRequest(req *rpc.RPCRequest) (*rpc.RPCResponse, error)
	Close() error
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
