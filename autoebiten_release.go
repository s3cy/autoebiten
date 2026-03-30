//go:build release

package autoebiten

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Capture is a no-op in release mode.
func Capture(screen *ebiten.Image) {
	// No-op in release mode
}

// Update is a no-op in release mode.
func Update() bool {
	return true
}

// IsKeyPressed wraps ebiten.IsKeyPressed directly in release mode.
func IsKeyPressed(key ebiten.Key) bool {
	return ebiten.IsKeyPressed(key)
}

// CursorPosition wraps ebiten.CursorPosition directly in release mode.
func CursorPosition() (x, y int) {
	return ebiten.CursorPosition()
}

// Wheel wraps ebiten.Wheel directly in release mode.
func Wheel() (x, y float64) {
	return ebiten.Wheel()
}

// IsMouseButtonPressed wraps ebiten.IsMouseButtonPressed directly in release mode.
func IsMouseButtonPressed(button ebiten.MouseButton) bool {
	return ebiten.IsMouseButtonPressed(button)
}
