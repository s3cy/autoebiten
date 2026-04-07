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

// Response wraps an RPC response with captured output.
// This is the format returned by the proxy to CLI clients.
type Response struct {
	rpc.RPCResponse
	Output string `json:"output,omitempty"` // Diff output from game
}

// Server wraps game RPC calls and captures output.
type Server struct {
	gameClient  GameClient
	outputFiles *output.FilePath
	mu          sync.Mutex
}

// GameClient is the interface for game RPC clients.
type GameClient interface {
	SendRequest(req *rpc.RPCRequest) (*rpc.RPCResponse, error)
	Close() error
}

// NewServer creates a new proxy server.
func NewServer(gameClient GameClient, outputFiles *output.FilePath) *Server {
	return &Server{
		gameClient:  gameClient,
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
func (s *Server) ForwardRequest(req *rpc.RPCRequest) (*Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Forward request to game
	gameResp, err := s.gameClient.SendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("game request failed: %w", err)
	}

	// Generate diff directly from file paths
	diff := output.GenerateDiff(s.outputFiles.Snapshot, s.outputFiles.Log)

	// Copy log to snapshot for next command
	if err := output.CopyFile(s.outputFiles.Snapshot, s.outputFiles.Log); err != nil {
		return nil, fmt.Errorf("failed to write snapshot: %w", err)
	}

	// Build proxy response
	resp := &Response{
		RPCResponse: *gameResp,
		Output:      diff,
	}

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
	resp := &Response{
		RPCResponse: rpc.ErrorResponse(id, code, message),
	}
	sendResponse(conn, resp)
}

func sendResponse(conn net.Conn, resp *Response) {
	encoder := json.NewEncoder(conn)
	encoder.Encode(resp)
}
