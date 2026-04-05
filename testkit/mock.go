package testkit

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/s3cy/autoebiten/internal/input"
	"github.com/s3cy/autoebiten/internal/server"
)

// GameUpdate is the interface required for Mock white-box testing.
// Games must implement at least an Update method.
type GameUpdate interface {
	Update() error
}

// Mock provides white-box testing control over a game running in the same process.
// It injects inputs directly into the game's update loop without requiring RPC.
type Mock struct {
	t    *testing.T
	game GameUpdate

	// Input state buffers
	actions []func(vi *input.VirtualInput, inputTime input.InputTime)
}

// NewMock creates a new Mock controller for white-box testing.
// The provided game must implement at least Update() error.
//
// The Mock is automatically cleaned up when the test ends via t.Cleanup().
func NewMock(t *testing.T, game GameUpdate) *Mock {
	t.Helper()

	m := &Mock{
		t:    t,
		game: game,
	}

	return m
}

// InjectKeyPress buffers a key press event to be applied on the next Tick.
func (m *Mock) InjectKeyPress(key ebiten.Key) {
	m.actions = append(m.actions, func(vi *input.VirtualInput, inputTime input.InputTime) {
		vi.InjectKeyPress(input.Key(key), inputTime)
	})
}

// InjectKeyRelease buffers a key release event to be applied on the next Tick.
func (m *Mock) InjectKeyRelease(key ebiten.Key) {
	m.actions = append(m.actions, func(vi *input.VirtualInput, inputTime input.InputTime) {
		vi.InjectKeyRelease(input.Key(key), inputTime)
	})
}

// InjectMousePosition sets the mouse cursor position.
func (m *Mock) InjectMousePosition(x, y int) {
	m.actions = append(m.actions, func(vi *input.VirtualInput, inputTime input.InputTime) {
		vi.InjectCursorMove(x, y)
	})
}

// InjectMouseButtonPress buffers a mouse button press event.
func (m *Mock) InjectMouseButtonPress(button ebiten.MouseButton) {
	m.actions = append(m.actions, func(vi *input.VirtualInput, inputTime input.InputTime) {
		vi.InjectMouseButtonPress(input.MouseButton(button), inputTime)
	})
}

// InjectMouseButtonRelease buffers a mouse button release event.
func (m *Mock) InjectMouseButtonRelease(button ebiten.MouseButton) {
	m.actions = append(m.actions, func(vi *input.VirtualInput, inputTime input.InputTime) {
		vi.InjectMouseButtonRelease(input.MouseButton(button), inputTime)
	})
}

// InjectWheel sets the wheel scroll delta.
func (m *Mock) InjectWheel(x, y float64) {
	m.actions = append(m.actions, func(vi *input.VirtualInput, inputTime input.InputTime) {
		vi.InjectWheelMove(x, y)
	})
}

// Tick advances the game by one tick, applying buffered inputs.
// This calls the game's Update() method once.
func (m *Mock) Tick() {
	server.IncrementTick()
	tick := server.Tick()

	vi := input.Get()

	for _, action := range m.actions {
		action(vi, input.NewInputTimeFromTick(tick, server.IncrementSubtick()))
	}

	// Call game update
	if err := m.game.Update(); err != nil {
		m.t.Errorf("testkit: game update failed: %v", err)
	}

	// Clear input buffers
	m.actions = m.actions[:0]
}

// Ticks advances the game by N ticks.
// This is equivalent to calling Tick() N times.
func (m *Mock) Ticks(n int) {
	for i := 0; i < n; i++ {
		m.Tick()
	}
}

// Game returns the underlying game instance.
// This allows direct access to game state for assertions.
func (m *Mock) Game() GameUpdate {
	return m.game
}
