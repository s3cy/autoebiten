//go:build release

package inpututil

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func IsKeyJustPressed(key ebiten.Key) bool {
	return inpututil.IsKeyJustPressed(key)
}

func IsKeyJustReleased(key ebiten.Key) bool {
	return inpututil.IsKeyJustReleased(key)
}

func KeyPressDuration(key ebiten.Key) int {
	return inpututil.KeyPressDuration(key)
}

func IsMouseButtonJustPressed(button ebiten.MouseButton) bool {
	return inpututil.IsMouseButtonJustPressed(button)
}

func IsMouseButtonJustReleased(button ebiten.MouseButton) bool {
	return inpututil.IsMouseButtonJustReleased(button)
}

func MouseButtonPressDuration(button ebiten.MouseButton) int {
	return inpututil.MouseButtonPressDuration(button)
}
