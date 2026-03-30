//go:build !release

package inpututil

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/s3cy/autoebiten"
	"github.com/s3cy/autoebiten/internal/input"
	"github.com/s3cy/autoebiten/internal/server"
)

func IsKeyJustPressed(key ebiten.Key) bool {
	switch autoebiten.GetMode() {
	case autoebiten.InjectionOnly:
		return input.Get().IsKeyJustPressed(key, server.Tick())
	case autoebiten.Passthrough:
		return inpututil.IsKeyJustPressed(key)
	case autoebiten.InjectionFallback:
		if input.Get().IsKeyJustPressed(key, server.Tick()) {
			return true
		}
		return inpututil.IsKeyJustPressed(key)
	}
	return false
}

func IsKeyJustReleased(key ebiten.Key) bool {
	switch autoebiten.GetMode() {
	case autoebiten.InjectionOnly:
		return input.Get().IsKeyJustReleased(key, server.Tick())
	case autoebiten.Passthrough:
		return inpututil.IsKeyJustReleased(key)
	case autoebiten.InjectionFallback:
		if input.Get().IsKeyJustReleased(key, server.Tick()) {
			return true
		}
		return inpututil.IsKeyJustReleased(key)
	}
	return false
}

func KeyPressDuration(key ebiten.Key) int {
	switch autoebiten.GetMode() {
	case autoebiten.InjectionOnly:
		return int(input.Get().KeyPressDuration(key, server.Tick()))
	case autoebiten.Passthrough:
		return inpututil.KeyPressDuration(key)
	case autoebiten.InjectionFallback:
		d := input.Get().KeyPressDuration(key, server.Tick())
		if d > 0 {
			return int(d)
		}
		return inpututil.KeyPressDuration(key)
	}
	return 0
}

func IsMouseButtonJustPressed(button ebiten.MouseButton) bool {
	switch autoebiten.GetMode() {
	case autoebiten.InjectionOnly:
		return input.Get().IsMouseButtonJustPressed(button, server.Tick())
	case autoebiten.Passthrough:
		return inpututil.IsMouseButtonJustPressed(button)
	case autoebiten.InjectionFallback:
		if input.Get().IsMouseButtonJustPressed(button, server.Tick()) {
			return true
		}
		return inpututil.IsMouseButtonJustPressed(button)
	}
	return false
}

func IsMouseButtonJustReleased(button ebiten.MouseButton) bool {
	switch autoebiten.GetMode() {
	case autoebiten.InjectionOnly:
		return input.Get().IsMouseButtonJustReleased(button, server.Tick())
	case autoebiten.Passthrough:
		return inpututil.IsMouseButtonJustReleased(button)
	case autoebiten.InjectionFallback:
		if input.Get().IsMouseButtonJustReleased(button, server.Tick()) {
			return true
		}
		return inpututil.IsMouseButtonJustReleased(button)
	}
	return false
}

func MouseButtonPressDuration(button ebiten.MouseButton) int {
	switch autoebiten.GetMode() {
	case autoebiten.InjectionOnly:
		return int(input.Get().MouseButtonPressDuration(button, server.Tick()))
	case autoebiten.Passthrough:
		return inpututil.MouseButtonPressDuration(button)
	case autoebiten.InjectionFallback:
		d := input.Get().MouseButtonPressDuration(button, server.Tick())
		if d > 0 {
			return int(d)
		}
		return inpututil.MouseButtonPressDuration(button)
	}
	return 0
}
