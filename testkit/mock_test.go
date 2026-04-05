package testkit

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// mockGame is a simple game for testing Mock.
type mockGame struct {
	updateCount int
	lastError   error
}

func (g *mockGame) Update() error {
	g.updateCount++
	return g.lastError
}

// TestMockNewMock tests creating a Mock.
func TestMockNewMock(t *testing.T) {
	game := &mockGame{}
	mock := NewMock(t, game)

	if mock == nil {
		t.Fatal("NewMock returned nil")
	}

	if mock.game != game {
		t.Error("mock game mismatch")
	}
}

// TestMockInjectKeyPress tests key press injection.
func TestMockInjectKeyPress(t *testing.T) {
	game := &mockGame{}
	mock := NewMock(t, game)

	mock.InjectKeyPress(ebiten.KeyA)
	mock.InjectKeyPress(ebiten.KeyB)

	// Verify actions are buffered
	if len(mock.actions) != 2 {
		t.Errorf("expected 2 actions, got %d", len(mock.actions))
	}
}

// TestMockInjectKeyRelease tests key release injection.
func TestMockInjectKeyRelease(t *testing.T) {
	game := &mockGame{}
	mock := NewMock(t, game)

	mock.InjectKeyRelease(ebiten.KeyA)

	if len(mock.actions) != 1 {
		t.Errorf("expected 1 action, got %d", len(mock.actions))
	}
}

// TestMockInjectMousePosition tests mouse position injection.
func TestMockInjectMousePosition(t *testing.T) {
	game := &mockGame{}
	mock := NewMock(t, game)

	mock.InjectMousePosition(100, 200)

	if len(mock.actions) != 1 {
		t.Errorf("expected 1 action, got %d", len(mock.actions))
	}
}

// TestMockInjectMouseButton tests mouse button injection.
func TestMockInjectMouseButton(t *testing.T) {
	game := &mockGame{}
	mock := NewMock(t, game)

	mock.InjectMouseButtonPress(ebiten.MouseButtonLeft)
	mock.InjectMouseButtonRelease(ebiten.MouseButtonLeft)
	mock.InjectMouseButtonPress(ebiten.MouseButtonRight)

	if len(mock.actions) != 3 {
		t.Errorf("expected 3 actions, got %d", len(mock.actions))
	}
}

// TestMockInjectWheel tests wheel injection.
func TestMockInjectWheel(t *testing.T) {
	game := &mockGame{}
	mock := NewMock(t, game)

	mock.InjectWheel(1.5, -2.5)

	if len(mock.actions) != 1 {
		t.Errorf("expected 1 action, got %d", len(mock.actions))
	}
}

// TestMockTick tests the Tick method.
func TestMockTick(t *testing.T) {
	game := &mockGame{}
	mock := NewMock(t, game)

	initialCount := game.updateCount
	mock.Tick()

	if game.updateCount != initialCount+1 {
		t.Errorf("expected update count %d, got %d", initialCount+1, game.updateCount)
	}
}

// TestMockTicks tests the Ticks method.
func TestMockTicks(t *testing.T) {
	game := &mockGame{}
	mock := NewMock(t, game)

	mock.Ticks(5)

	if game.updateCount != 5 {
		t.Errorf("expected update count 5, got %d", game.updateCount)
	}
}

// TestMockTickClearsInputs tests that Tick clears input buffers.
func TestMockTickClearsInputs(t *testing.T) {
	game := &mockGame{}
	mock := NewMock(t, game)

	mock.InjectKeyPress(ebiten.KeyA)
	mock.InjectKeyRelease(ebiten.KeyB)
	mock.InjectMouseButtonPress(ebiten.MouseButtonLeft)
	mock.InjectWheel(1.0, 2.0)

	if len(mock.actions) != 4 {
		t.Errorf("expected 4 actions before Tick, got %d", len(mock.actions))
	}

	mock.Tick()

	// Verify actions buffer is cleared
	if len(mock.actions) != 0 {
		t.Errorf("expected 0 actions after Tick, got %d", len(mock.actions))
	}
}

// TestMockGameAccess tests accessing the underlying game.
func TestMockGameAccess(t *testing.T) {
	game := &mockGame{}
	mock := NewMock(t, game)

	if mock.Game() != game {
		t.Error("Game() returned wrong game")
	}
}
