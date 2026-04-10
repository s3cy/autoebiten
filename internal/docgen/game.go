package docgen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/s3cy/autoebiten/testkit"
)

// GameSession wraps testkit.Game with docgen-specific state.
type GameSession struct {
	game       *testkit.Game
	t          *testing.T
	socketPath string
	ctx        *Context
}

// buildGame compiles a game binary and returns its path.
func buildGame(gameDir string) (string, error) {
	binaryName := filepath.Base(gameDir) + "_docgen"
	binaryPath := filepath.Join(gameDir, binaryName)

	cmd := exec.Command("go", "build", "-o", binaryName, ".")
	cmd.Dir = gameDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build game: %w\n%s", err, output)
	}

	return binaryPath, nil
}

// LaunchGame builds and starts a game, returning a session.
func LaunchGame(ctx *Context, args ...string) (*GameSession, error) {
	if ctx == nil || ctx.Config == nil || ctx.Config.GameDir == "" {
		return nil, fmt.Errorf("config.GameDir not set")
	}

	// Build binary
	binaryPath, err := buildGame(ctx.Config.GameDir)
	if err != nil {
		return nil, err
	}

	// Create test context (testkit requires *testing.T)
	t := &testing.T{}

	// Build options
	opts := []testkit.Option{testkit.WithTimeout(30 * time.Second)}
	if len(args) > 0 {
		opts = append(opts, testkit.WithArgs(args...))
	}

	// Launch via testkit
	game := testkit.Launch(t, binaryPath, opts...)

	// Wait for ready
	ready := game.WaitFor(func() bool {
		return game.Ping() == nil
	}, 5*time.Second)

	if !ready {
		game.Shutdown()
		os.Remove(binaryPath)
		return nil, fmt.Errorf("game failed to start")
	}

	return &GameSession{
		game:       game,
		t:          t,
		socketPath: game.SocketPath(),
		ctx:        ctx,
	}, nil
}

// EndGame shuts down the game and cleans up.
func EndGame(session *GameSession) error {
	if session == nil || session.game == nil {
		return nil
	}

	session.game.Shutdown()

	// Remove built binary
	if session.ctx != nil && session.ctx.Config != nil && session.ctx.Config.GameDir != "" {
		binaryName := filepath.Base(session.ctx.Config.GameDir) + "_docgen"
		binaryPath := filepath.Join(session.ctx.Config.GameDir, binaryName)
		os.Remove(binaryPath)
	}

	return nil
}

// Delay pauses execution for the specified duration.
// Used for crash scenarios where game crashes after N seconds.
func Delay(duration string) {
	d, err := time.ParseDuration(duration)
	if err != nil {
		d = 1 * time.Second // default fallback
	}
	time.Sleep(d)
}