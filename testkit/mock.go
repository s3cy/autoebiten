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
	keyPresses   []input.Key
	keyReleases  []input.Key
	mousePos     struct{ x, y int }
	mouseButtons []struct {
		button  input.MouseButton
		pressed bool
	}
	wheelDelta struct{ x, y float64 }
}

// NewMock creates a new Mock controller for white-box testing.
// The provided game must implement at least Update() error.
//
// The Mock is automatically cleaned up when the test ends via t.Cleanup().
func NewMock(t *testing.T, game GameUpdate) *Mock {
	t.Helper()

	m := &Mock{
		t:           t,
		game:        game,
		keyPresses:  make([]input.Key, 0),
		keyReleases: make([]input.Key, 0),
		mouseButtons: make([]struct {
			button  input.MouseButton
			pressed bool
		}, 0),
	}

	// Register cleanup (no-op for mock, but maintains symmetry with Game)
	t.Cleanup(func() {
		// Clear any pending inputs
		m.keyPresses = nil
		m.keyReleases = nil
		m.mouseButtons = nil
	})

	return m
}

// InjectKeyPress buffers a key press event to be applied on the next Tick.
func (m *Mock) InjectKeyPress(key ebiten.Key) {
	m.keyPresses = append(m.keyPresses, input.Key(key))
}

// InjectKeyRelease buffers a key release event to be applied on the next Tick.
func (m *Mock) InjectKeyRelease(key ebiten.Key) {
	m.keyReleases = append(m.keyReleases, input.Key(key))
}

// InjectMousePosition sets the mouse cursor position.
func (m *Mock) InjectMousePosition(x, y int) {
	m.mousePos.x = x
	m.mousePos.y = y
}

// InjectMouseButtonPress buffers a mouse button press event.
func (m *Mock) InjectMouseButtonPress(button ebiten.MouseButton) {
	m.mouseButtons = append(m.mouseButtons, struct {
		button  input.MouseButton
		pressed bool
	}{button: input.MouseButton(button), pressed: true})
}

// InjectMouseButtonRelease buffers a mouse button release event.
func (m *Mock) InjectMouseButtonRelease(button ebiten.MouseButton) {
	m.mouseButtons = append(m.mouseButtons, struct {
		button  input.MouseButton
		pressed bool
	}{button: input.MouseButton(button), pressed: false})
}

// InjectWheel sets the wheel scroll delta.
func (m *Mock) InjectWheel(x, y float64) {
	m.wheelDelta.x = x
	m.wheelDelta.y = y
}

// Tick advances the game by one tick, applying buffered inputs.
// This calls the game's Update() method once.
func (m *Mock) Tick() {
	// Get current tick from server
	tick := server.Tick()
	inputTime := input.NewInputTimeFromTick(tick, 0)

	vi := input.Get()

	// Apply buffered key presses
	for _, key := range m.keyPresses {
		vi.InjectKeyPress(key, inputTime)
	}

	// Apply buffered key releases
	for _, key := range m.keyReleases {
		vi.InjectKeyRelease(key, inputTime)
	}

	// Apply mouse position
	vi.InjectCursorMove(m.mousePos.x, m.mousePos.y)

	// Apply buffered mouse buttons
	for _, btn := range m.mouseButtons {
		if btn.pressed {
			vi.InjectMouseButtonPress(btn.button, inputTime)
		} else {
			vi.InjectMouseButtonRelease(btn.button, inputTime)
		}
	}

	// Apply wheel delta
	vi.InjectWheelMove(m.wheelDelta.x, m.wheelDelta.y)

	// Call game update
	if err := m.game.Update(); err != nil {
		m.t.Errorf("testkit: game update failed: %v", err)
	}

	// Clear input buffers
	m.keyPresses = m.keyPresses[:0]
	m.keyReleases = m.keyReleases[:0]
	m.mouseButtons = m.mouseButtons[:0]
	m.wheelDelta.x = 0
	m.wheelDelta.y = 0
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
