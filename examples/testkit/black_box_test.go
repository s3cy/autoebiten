package testkit_test

import (
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/s3cy/autoebiten/testkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPlayerMovement demonstrates black-box testing with StateQuery.
// Requires a game built from examples/state_exporter/cmd.
func TestPlayerMovement(t *testing.T) {
	// Launch game in separate process (binary built from cmd directory)
	game := testkit.Launch(t, "./examples/state_exporter/cmd/state_exporter",
		testkit.WithTimeout(30*time.Second))
	defer game.Shutdown()

	// Wait for game to be ready
	ready := game.WaitFor(func() bool {
		return game.Ping() == nil
	}, 5*time.Second)
	require.True(t, ready, "game should be ready within timeout")

	// Get initial position
	x, err := game.StateQuery("gamestate", "Player.X")
	require.NoError(t, err)
	initialX := x.(float64)

	// Move player right for 10 ticks
	err = game.HoldKey(ebiten.KeyArrowRight, 10)
	require.NoError(t, err)

	// Verify position changed
	x, err = game.StateQuery("gamestate", "Player.X")
	require.NoError(t, err)
	assert.Greater(t, x.(float64), initialX, "player should have moved right")
}

// TestHealthModification demonstrates custom commands and state verification.
func TestHealthModification(t *testing.T) {
	game := testkit.Launch(t, "./examples/state_exporter/cmd/state_exporter")
	defer game.Shutdown()

	game.WaitFor(func() bool { return game.Ping() == nil }, 5*time.Second)

	// Get initial health
	health, err := game.StateQuery("gamestate", "Player.Health")
	require.NoError(t, err)
	require.Equal(t, 100.0, health, "initial health should be 100")

	// Heal via custom command
	resp, err := game.RunCustom("heal", "")
	require.NoError(t, err)
	assert.Contains(t, resp, "Healed")

	// Health unchanged (already at max)
	health, err = game.StateQuery("gamestate", "Player.Health")
	require.NoError(t, err)
	assert.Equal(t, 100.0, health)
}

// TestScreenshotCapture demonstrates visual verification.
func TestScreenshotCapture(t *testing.T) {
	game := testkit.Launch(t, "./examples/state_exporter/cmd/state_exporter")
	defer game.Shutdown()

	game.WaitFor(func() bool { return game.Ping() == nil }, 5*time.Second)

	// Capture screenshot
	img, err := game.Screenshot()
	require.NoError(t, err)
	require.NotNil(t, img)

	// Verify dimensions
	bounds := img.Bounds()
	assert.Equal(t, 640, bounds.Dx())
	assert.Equal(t, 480, bounds.Dy())
}

// TestEnemyStateQuery demonstrates querying array/slice state.
func TestEnemyStateQuery(t *testing.T) {
	game := testkit.Launch(t, "./examples/state_exporter/cmd/state_exporter")
	defer game.Shutdown()

	game.WaitFor(func() bool { return game.Ping() == nil }, 5*time.Second)

	// Query first enemy name
	name, err := game.StateQuery("gamestate", "Enemies.0.Name")
	require.NoError(t, err)
	assert.Equal(t, "Goblin", name)

	// Query second enemy health
	health, err := game.StateQuery("gamestate", "Enemies.1.Health")
	require.NoError(t, err)
	assert.Equal(t, 50.0, health)
}