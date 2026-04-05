package testkit

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/s3cy/autoebiten"
	"github.com/s3cy/autoebiten/internal/input"
	"github.com/s3cy/autoebiten/internal/rpc"
)

// keyToString converts an ebiten.Key to its string representation.
func keyToString(key ebiten.Key) string {
	// Reverse lookup in StringKeyMap using ebiten.Key
	for name, k := range input.StringKeyMap {
		if input.Key(key) == k {
			return name
		}
	}
	panic(fmt.Sprintf("unhandled key: %v", key))
}

// mouseButtonToString converts an ebiten.MouseButton to its string representation.
func mouseButtonToString(btn ebiten.MouseButton) string {
	// Reverse lookup in StringMouseButtonMap using ebiten.MouseButton
	for name, b := range input.StringMouseButtonMap {
		if input.MouseButton(btn) == b {
			return name
		}
	}
	panic(fmt.Sprintf("unhandled mouse button: %v", btn))
}

// Game provides black-box testing control over a game running in a separate process.
// It communicates with the game via JSON-RPC over Unix sockets.
type Game struct {
	t      *testing.T
	config *config

	// Process control
	process    *os.Process
	cmd        *exec.Cmd
	socketPath string

	// RPC client
	client *rpc.Client
}

// Launch starts a game binary in a separate process and returns a Game controller.
// The game must be an autoebiten-enabled game that listens on the Unix socket.
//
// The Game is automatically cleaned up when the test ends via t.Cleanup().
// Callers can also manually call Shutdown() for early cleanup.
//
// Options:
//   - WithTimeout: sets timeout for operations (default 30s)
//   - WithArgs: adds command-line arguments
//   - WithEnv: sets environment variables
func Launch(t *testing.T, binaryPath string, opts ...Option) *Game {
	t.Helper()

	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	// Validate binary exists
	if _, err := os.Stat(binaryPath); err != nil {
		t.Fatalf("testkit: game binary not found: %s", binaryPath)
	}

	// Create unique socket path
	socketPath := filepath.Join(rpc.DefaultSocketDir, fmt.Sprintf("autoebiten-testkit-%d-%d.sock", time.Now().UnixNano(), os.Getpid()))

	// Build command
	cmd := exec.Command(binaryPath, cfg.args...)
	cmd.Env = buildEnv(cfg, socketPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start process
	if err := cmd.Start(); err != nil {
		t.Fatalf("testkit: failed to start game: %v", err)
	}

	g := &Game{
		t:          t,
		config:     cfg,
		cmd:        cmd,
		process:    cmd.Process,
		socketPath: socketPath,
	}

	// Wait for socket to appear
	if err := waitForSocket(socketPath, cfg.timeout); err != nil {
		g.kill()
		t.Fatalf("testkit: game did not create socket in time: %v", err)
	}

	// Create RPC client
	client, err := rpc.NewClient()
	if err != nil {
		g.kill()
		t.Fatalf("testkit: failed to create RPC client: %v", err)
	}
	g.client = client

	// Register cleanup
	t.Cleanup(func() {
		g.Shutdown()
	})

	return g
}

// waitForSocket polls for the socket file to appear.
func waitForSocket(path string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for socket %s", path)
		case <-ticker.C:
			if _, err := os.Stat(path); err == nil {
				return nil
			}
		}
	}
}

// kill forcefully terminates the process.
func (g *Game) kill() {
	if g.process != nil {
		g.process.Kill()
		g.process.Wait()
		g.process = nil
	}
	if g.client != nil {
		g.client.Close()
		g.client = nil
	}
	os.Remove(g.socketPath)
}

// Shutdown gracefully terminates the game process.
// It first attempts SIGTERM, then falls back to SIGKILL if needed.
func (g *Game) Shutdown() error {
	if g.process == nil {
		return nil
	}

	// Try graceful shutdown first
	if err := g.process.Signal(syscall.SIGTERM); err == nil {
		// Wait for process to exit with timeout
		done := make(chan struct{})
		go func() {
			g.process.Wait()
			close(done)
		}()

		select {
		case <-done:
			// Process exited cleanly
		case <-time.After(5 * time.Second):
			// Timeout, force kill
			g.process.Kill()
			g.process.Wait()
		}
	}

	// Close RPC client
	if g.client != nil {
		g.client.Close()
		g.client = nil
	}

	// Clean up socket file
	os.Remove(g.socketPath)

	g.process = nil
	return nil
}

// Ping checks if the game is responsive.
// Returns an error if the game is not running or not responding.
func (g *Game) Ping() error {
	if g.client == nil {
		return ErrGameNotRunning
	}

	req, err := rpc.BuildRequest("ping", nil)
	if err != nil {
		return fmt.Errorf("failed to build ping request: %w", err)
	}

	resp, err := g.client.SendRequest(req)
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("ping error: %s", resp.Error.Message)
	}

	var result rpc.PingResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return fmt.Errorf("failed to unmarshal ping result: %w", err)
	}

	if !result.OK {
		return fmt.Errorf("ping returned not OK")
	}

	return nil
}

// WaitFor polls the provided function until it returns true or the timeout is reached.
// The function is called every 100ms.
func (g *Game) WaitFor(fn func() bool, timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			if fn() {
				return true
			}
		}
	}
}

// PressKey sends a key press event to the game.
func (g *Game) PressKey(key ebiten.Key) error {
	return g.sendInput("press", key, 0)
}

// ReleaseKey sends a key release event to the game.
func (g *Game) ReleaseKey(key ebiten.Key) error {
	return g.sendInput("release", key, 0)
}

// HoldKey holds a key down for the specified number of ticks.
func (g *Game) HoldKey(key ebiten.Key, ticks int64) error {
	return g.sendInput("hold", key, ticks)
}

// sendInput sends an input command to the game.
func (g *Game) sendInput(action string, key ebiten.Key, ticks int64) error {
	if g.client == nil {
		return ErrGameNotRunning
	}

	params := rpc.InputParams{
		Action:        action,
		Key:           keyToString(key),
		DurationTicks: ticks,
	}

	req, err := rpc.BuildRequest("input", params)
	if err != nil {
		return fmt.Errorf("failed to build input request: %w", err)
	}

	resp, err := g.client.SendRequest(req)
	if err != nil {
		return fmt.Errorf("input request failed: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("input error: %s", resp.Error.Message)
	}

	return nil
}

// MoveMouse moves the mouse cursor to the specified position.
func (g *Game) MoveMouse(x, y int) error {
	return g.sendMouse("move", x, y, "", 0)
}

// PressMouseButton sends a mouse button press event.
func (g *Game) PressMouseButton(button ebiten.MouseButton) error {
	return g.sendMouse("press", 0, 0, mouseButtonToString(button), 0)
}

// ReleaseMouseButton sends a mouse button release event.
func (g *Game) ReleaseMouseButton(button ebiten.MouseButton) error {
	return g.sendMouse("release", 0, 0, mouseButtonToString(button), 0)
}

// sendMouse sends a mouse command to the game.
func (g *Game) sendMouse(action string, x, y int, button string, ticks int64) error {
	if g.client == nil {
		return ErrGameNotRunning
	}

	params := rpc.MouseParams{
		Action:        action,
		X:             x,
		Y:             y,
		Button:        button,
		DurationTicks: ticks,
	}

	req, err := rpc.BuildRequest("mouse", params)
	if err != nil {
		return fmt.Errorf("failed to build mouse request: %w", err)
	}

	resp, err := g.client.SendRequest(req)
	if err != nil {
		return fmt.Errorf("mouse request failed: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("mouse error: %s", resp.Error.Message)
	}

	return nil
}

// ScrollWheel sends a wheel scroll event.
func (g *Game) ScrollWheel(x, y float64) error {
	if g.client == nil {
		return ErrGameNotRunning
	}

	params := rpc.WheelParams{
		X: x,
		Y: y,
	}

	req, err := rpc.BuildRequest("wheel", params)
	if err != nil {
		return fmt.Errorf("failed to build wheel request: %w", err)
	}

	resp, err := g.client.SendRequest(req)
	if err != nil {
		return fmt.Errorf("wheel request failed: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("wheel error: %s", resp.Error.Message)
	}

	return nil
}

// Screenshot captures the current game screen and returns it as an image.
func (g *Game) Screenshot() (*image.RGBA, error) {
	if g.client == nil {
		return nil, ErrGameNotRunning
	}

	params := rpc.ScreenshotParams{
		Base64: true,
	}

	req, err := rpc.BuildRequest("screenshot", params)
	if err != nil {
		return nil, fmt.Errorf("failed to build screenshot request: %w", err)
	}

	resp, err := g.client.SendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("screenshot request failed: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("screenshot error: %s", resp.Error.Message)
	}

	var result rpc.ScreenshotResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal screenshot result: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("screenshot failed")
	}

	// Decode base64 data
	data, err := base64.StdEncoding.DecodeString(result.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode screenshot data: %w", err)
	}

	// Parse PNG data (assuming PNG format)
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode screenshot image: %w", err)
	}

	rgba, ok := img.(*image.RGBA)
	if !ok {
		// Convert to RGBA
		bounds := img.Bounds()
		rgba = image.NewRGBA(bounds)
		draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	}

	return rgba, nil
}

// ScreenshotToFile captures the current game screen and saves it to a file.
func (g *Game) ScreenshotToFile(path string) error {
	if g.client == nil {
		return ErrGameNotRunning
	}

	params := rpc.ScreenshotParams{
		Output: path,
	}

	req, err := rpc.BuildRequest("screenshot", params)
	if err != nil {
		return fmt.Errorf("failed to build screenshot request: %w", err)
	}

	resp, err := g.client.SendRequest(req)
	if err != nil {
		return fmt.Errorf("screenshot request failed: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("screenshot error: %s", resp.Error.Message)
	}

	var result rpc.ScreenshotResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return fmt.Errorf("failed to unmarshal screenshot result: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("screenshot failed")
	}

	return nil
}

// ScreenshotBase64 captures the current game screen and returns it as a base64 string.
func (g *Game) ScreenshotBase64() (string, error) {
	if g.client == nil {
		return "", ErrGameNotRunning
	}

	params := rpc.ScreenshotParams{
		Base64: true,
	}

	req, err := rpc.BuildRequest("screenshot", params)
	if err != nil {
		return "", fmt.Errorf("failed to build screenshot request: %w", err)
	}

	resp, err := g.client.SendRequest(req)
	if err != nil {
		return "", fmt.Errorf("screenshot request failed: %w", err)
	}

	if resp.Error != nil {
		return "", fmt.Errorf("screenshot error: %s", resp.Error.Message)
	}

	var result rpc.ScreenshotResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal screenshot result: %w", err)
	}

	if !result.Success {
		return "", fmt.Errorf("screenshot failed")
	}

	return result.Data, nil
}

// RunCustom executes a custom command registered by the game.
func (g *Game) RunCustom(name, request string) (string, error) {
	if g.client == nil {
		return "", ErrGameNotRunning
	}

	params := rpc.CustomParams{
		Name:    name,
		Request: request,
	}

	req, err := rpc.BuildRequest("custom", params)
	if err != nil {
		return "", fmt.Errorf("failed to build custom request: %w", err)
	}

	resp, err := g.client.SendRequest(req)
	if err != nil {
		return "", fmt.Errorf("custom request failed: %w", err)
	}

	if resp.Error != nil {
		return "", fmt.Errorf("custom error: %s", resp.Error.Message)
	}

	var result rpc.CustomResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal custom result: %w", err)
	}

	return result.Response, nil
}

// StateQuery queries game state via reflection-based path.
// The path uses dot notation, e.g., "Player.X", "Inventory.0.Name".
func (g *Game) StateQuery(name string, path string) (any, error) {
	response, err := g.RunCustom(autoebiten.StateExporterPathPrefix+name, path)
	if err != nil {
		return nil, err
	}

	// Response is JSON-encoded value
	var value any
	if err := json.Unmarshal([]byte(response), &value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal state query result: %w", err)
	}

	return value, nil
}
