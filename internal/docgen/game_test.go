package docgen

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDelay(t *testing.T) {
	start := time.Now()
	Delay("100ms")
	elapsed := time.Since(start)

	assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(100))
}

func TestBuildGame(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Use testkit's internal test games for testing
	testgameDir := filepath.Join("..", "..", "testkit", "internal", "testgames", "simple")

	binaryPath, err := buildGame(testgameDir)
	require.NoError(t, err)
	defer os.Remove(binaryPath)

	// Verify binary exists
	assert.FileExists(t, binaryPath)

	// Verify binary is executable
	info, err := os.Stat(binaryPath)
	require.NoError(t, err)
	assert.False(t, info.IsDir())
}

func TestBuildGameInvalidDir(t *testing.T) {
	_, err := buildGame("/nonexistent/path")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to build game")
}

func TestLaunchGameNoConfig(t *testing.T) {
	// Test with nil context
	_, err := LaunchGame(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config.GameDir not set")

	// Test with nil config
	ctx := &Context{}
	_, err = LaunchGame(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config.GameDir not set")

	// Test with empty GameDir
	ctx.SetConfig(&Config{GameDir: ""})
	_, err = LaunchGame(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config.GameDir not set")
}

func TestLaunchGameIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	testgameDir := filepath.Join("..", "..", "testkit", "internal", "testgames", "simple")
	ctx := NewContext()
	ctx.SetConfig(&Config{GameDir: testgameDir})

	session, err := LaunchGame(ctx)
	require.NoError(t, err)
	defer EndGame(session)

	// Verify session is initialized
	assert.NotNil(t, session.game)
	assert.NotEmpty(t, session.socketPath)
	assert.Equal(t, ctx, session.ctx)

	// Verify game is responsive
	err = session.game.Ping()
	assert.NoError(t, err)
}

func TestEndGameNilSession(t *testing.T) {
	err := EndGame(nil)
	assert.NoError(t, err)

	err = EndGame(&GameSession{})
	assert.NoError(t, err)
}

func TestEndGameCleanup(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	testgameDir := filepath.Join("..", "..", "testkit", "internal", "testgames", "simple")
	ctx := NewContext()
	ctx.SetConfig(&Config{GameDir: testgameDir})

	session, err := LaunchGame(ctx)
	require.NoError(t, err)

	// Get binary path for cleanup verification
	binaryName := filepath.Base(testgameDir) + "_docgen"
	binaryPath := filepath.Join(testgameDir, binaryName)
	assert.FileExists(t, binaryPath)

	// End session
	err = EndGame(session)
	assert.NoError(t, err)

	// Verify binary was removed
	assert.NoFileExists(t, binaryPath)
}