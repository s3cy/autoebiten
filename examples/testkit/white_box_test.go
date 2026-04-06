package testkit_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/s3cy/autoebiten"
	"github.com/s3cy/autoebiten/examples/state_exporter"
	"github.com/s3cy/autoebiten/testkit"
	"github.com/stretchr/testify/assert"
)

// TestPlayerMovesRight demonstrates white-box testing with Mock.
func TestPlayerMovesRight(t *testing.T) {
	// Set mode to only use injected inputs
	autoebiten.SetMode(autoebiten.InjectionOnly)

	// Create game instance in same process
	game := state_exporter.NewGame()
	mock := testkit.NewMock(t, game)

	initialX := game.State.Player.X

	// Inject input
	mock.InjectKeyPress(ebiten.KeyArrowRight)
	mock.Ticks(10)

	// Check state directly (no RPC needed)
	assert.Greater(t, game.State.Player.X, initialX)
}

// TestPlayerTakesDamage demonstrates damage logic testing.
func TestPlayerTakesDamage(t *testing.T) {
	autoebiten.SetMode(autoebiten.InjectionOnly)

	game := state_exporter.NewGame()
	mock := testkit.NewMock(t, game)

	initialHealth := game.State.Player.Health

	// Simulate damage key
	mock.InjectKeyPress(ebiten.KeyD)
	mock.Tick()

	assert.Less(t, game.State.Player.Health, initialHealth)
}

// TestComboInput demonstrates multiple inputs in one tick.
func TestComboInput(t *testing.T) {
	autoebiten.SetMode(autoebiten.InjectionOnly)

	game := state_exporter.NewGame()
	mock := testkit.NewMock(t, game)

	// Inject multiple inputs before tick
	mock.InjectKeyPress(ebiten.KeyArrowUp)
	mock.InjectKeyPress(ebiten.KeyArrowRight)
	mock.Ticks(5)

	// Player should have moved diagonally
	assert.Greater(t, game.State.Player.X, 100.0)
	assert.Less(t, game.State.Player.Y, 100.0)
}

// TestMouseInput demonstrates mouse position injection.
func TestMouseInput(t *testing.T) {
	autoebiten.SetMode(autoebiten.InjectionOnly)

	game := state_exporter.NewGame()
	mock := testkit.NewMock(t, game)

	// Inject mouse position
	mock.InjectMousePosition(200, 150)
	mock.Tick()

	// In a real game, you'd check mouse-dependent behavior
	// For this example, we just verify injection works
	x, y := autoebiten.CursorPosition()
	assert.Equal(t, 200, x)
	assert.Equal(t, 150, y)
}

// TestWheelInput demonstrates wheel scroll injection.
func TestWheelInput(t *testing.T) {
	autoebiten.SetMode(autoebiten.InjectionOnly)

	game := state_exporter.NewGame()
	mock := testkit.NewMock(t, game)

	mock.InjectWheel(0, -5)
	mock.Tick()

	x, y := autoebiten.Wheel()
	assert.Equal(t, 0.0, x)
	assert.Equal(t, -5.0, y)
}