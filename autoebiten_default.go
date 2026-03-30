//go:build !release

package autoebiten

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/s3cy/autoebiten/internal/input"
	"github.com/s3cy/autoebiten/internal/server"
)

// Capture processes screenshots for injection.
func Capture(screen *ebiten.Image) {
	server.ProcessScreenshots(screen)
}

// Update runs the internal update loop.
func Update() bool {
	return server.Update()
}

// IsKeyPressed returns whether the key is pressed, respecting the current mode.
func IsKeyPressed(key ebiten.Key) bool {
	switch currentMode {
	case InjectionOnly:
		return input.Get().IsKeyPressed(key, server.Tick())
	case Passthrough:
		return ebiten.IsKeyPressed(key)
	case InjectionFallback:
		if input.Get().IsKeyPressed(key, server.Tick()) {
			return true
		}
		return ebiten.IsKeyPressed(key)
	}
	return false
}

// CursorPosition returns the cursor position, respecting the current mode.
func CursorPosition() (x, y int) {
	switch currentMode {
	case InjectionOnly:
		return input.Get().CursorPosition()
	case Passthrough:
		return ebiten.CursorPosition()
	case InjectionFallback:
		cx, cy := input.Get().CursorPosition()
		if cx != 0 || cy != 0 {
			return cx, cy
		}
		return ebiten.CursorPosition()
	}
	return 0, 0
}

// Wheel returns the mouse wheel scroll amount, respecting the current mode.
func Wheel() (x, y float64) {
	switch currentMode {
	case InjectionOnly:
		return input.Get().Wheel()
	case Passthrough:
		return ebiten.Wheel()
	case InjectionFallback:
		wx, wy := input.Get().Wheel()
		if wx != 0 || wy != 0 {
			return wx, wy
		}
		return ebiten.Wheel()
	}
	return 0, 0
}

// IsMouseButtonPressed returns whether the mouse button is pressed,
// respecting the current mode.
func IsMouseButtonPressed(button ebiten.MouseButton) bool {
	switch currentMode {
	case InjectionOnly:
		return input.Get().IsMouseButtonPressed(button, server.Tick())
	case Passthrough:
		return ebiten.IsMouseButtonPressed(button)
	case InjectionFallback:
		if input.Get().IsMouseButtonPressed(button, server.Tick()) {
			return true
		}
		return ebiten.IsMouseButtonPressed(button)
	}
	return false
}
