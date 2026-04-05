package testkit

import (
	"testing"

	"github.com/s3cy/autoebiten/internal/input"
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

	mock.InjectKeyPress(input.KeyA)
	mock.InjectKeyPress(input.KeyB)

	// Verify inputs are buffered
	if len(mock.keyPresses) != 2 {
		t.Errorf("expected 2 key presses, got %d", len(mock.keyPresses))
	}
}

// TestMockInjectKeyRelease tests key release injection.
func TestMockInjectKeyRelease(t *testing.T) {
	game := &mockGame{}
	mock := NewMock(t, game)

	mock.InjectKeyRelease(input.KeyA)

	if len(mock.keyReleases) != 1 {
		t.Errorf("expected 1 key release, got %d", len(mock.keyReleases))
	}
}

// TestMockInjectMousePosition tests mouse position injection.
func TestMockInjectMousePosition(t *testing.T) {
	game := &mockGame{}
	mock := NewMock(t, game)

	mock.InjectMousePosition(100, 200)

	if mock.mousePos.x != 100 || mock.mousePos.y != 200 {
		t.Errorf("expected mouse pos (100, 200), got (%d, %d)", mock.mousePos.x, mock.mousePos.y)
	}
}

// TestMockInjectMouseButton tests mouse button injection.
func TestMockInjectMouseButton(t *testing.T) {
	game := &mockGame{}
	mock := NewMock(t, game)

	mock.InjectMouseButtonPress(input.MouseButtonLeft)
	mock.InjectMouseButtonRelease(input.MouseButtonLeft)
	mock.InjectMouseButtonPress(input.MouseButtonRight)

	if len(mock.mouseButtons) != 3 {
		t.Errorf("expected 3 mouse button events, got %d", len(mock.mouseButtons))
	}
}

// TestMockInjectWheel tests wheel injection.
func TestMockInjectWheel(t *testing.T) {
	game := &mockGame{}
	mock := NewMock(t, game)

	mock.InjectWheel(1.5, -2.5)

	if mock.wheelDelta.x != 1.5 || mock.wheelDelta.y != -2.5 {
		t.Errorf("expected wheel delta (1.5, -2.5), got (%f, %f)", mock.wheelDelta.x, mock.wheelDelta.y)
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

	mock.InjectKeyPress(input.KeyA)
	mock.InjectKeyRelease(input.KeyB)
	mock.InjectMouseButtonPress(input.MouseButtonLeft)
	mock.InjectWheel(1.0, 2.0)

	mock.Tick()

	// Verify all buffers are cleared
	if len(mock.keyPresses) != 0 {
		t.Error("keyPresses buffer not cleared")
	}
	if len(mock.keyReleases) != 0 {
		t.Error("keyReleases buffer not cleared")
	}
	if len(mock.mouseButtons) != 0 {
		t.Error("mouseButtons buffer not cleared")
	}
	if mock.wheelDelta.x != 0 || mock.wheelDelta.y != 0 {
		t.Error("wheelDelta buffer not cleared")
	}
}

// TestMockRegisterCustom tests custom command registration.
func TestMockRegisterCustom(t *testing.T) {
	game := &mockGame{}
	mock := NewMock(t, game)

	called := false
	mock.RegisterCustom("test", func(req string) string {
		called = true
		return "response: " + req
	})

	// Verify handler is registered
	if _, ok := mock.customHandlers["test"]; !ok {
		t.Error("custom handler not registered")
	}

	// Call the handler
	response := mock.RunCustom("test", "hello")

	if !called {
		t.Error("custom handler not called")
	}
	if response != "response: hello" {
		t.Errorf("unexpected response: %s", response)
	}
}

// TestMockRunCustomNotFound tests calling an unregistered custom command.
func TestMockRunCustomNotFound(t *testing.T) {
	game := &mockGame{}
	mock := NewMock(t, game)

	response := mock.RunCustom("nonexistent", "request")

	if response != "" {
		t.Errorf("expected empty response for unknown command, got: %s", response)
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
