//go:build !release

package autoebiten

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/s3cy/autoebiten/integrate"
)

func panicIfPatched() {
	if integrate.IsPatched {
		panic("This function should NOT be called when Ebiten is patched." +
			"The patch handles this automatically. Remove it from your game.")
	}
}

// Capture processes screenshots for injection.
// Panic when using the patch integration method.
func Capture(screen image.Image) {
	panicIfPatched()
	integrate.Capture(screen)
}

// Update runs the internal update loop.
// Panic when using the patch integration method.
func Update() bool {
	panicIfPatched()
	return integrate.Update()
}

// IsKeyPressed wraps ebiten.IsKeyPressed, respecting the current mode.
// Panic when using the patch integration method.
func IsKeyPressed(key ebiten.Key) bool {
	panicIfPatched()
	switch currentMode {
	case InjectionOnly:
		return integrate.IsKeyPressed(integrate.Key(key))
	case Passthrough:
		return ebiten.IsKeyPressed(key)
	case InjectionFallback:
		if integrate.IsKeyPressed(integrate.Key(key)) {
			return true
		}
		return ebiten.IsKeyPressed(key)
	}
	return false
}

// CursorPosition wraps ebiten.CursorPosition, respecting the current mode.
// Panic when using the patch integration method.
func CursorPosition() (x, y int) {
	panicIfPatched()
	switch currentMode {
	case InjectionOnly:
		return integrate.CursorPosition()
	case Passthrough:
		return ebiten.CursorPosition()
	case InjectionFallback:
		cx, cy := integrate.CursorPosition()
		if cx != 0 || cy != 0 {
			return cx, cy
		}
		return ebiten.CursorPosition()
	}
	return 0, 0
}

// Wheel wraps ebiten.Wheel, respecting the current mode.
// Panic when using the patch integration method.
func Wheel() (x, y float64) {
	panicIfPatched()
	switch currentMode {
	case InjectionOnly:
		return integrate.Wheel()
	case Passthrough:
		return ebiten.Wheel()
	case InjectionFallback:
		wx, wy := integrate.Wheel()
		if wx != 0 || wy != 0 {
			return wx, wy
		}
		return ebiten.Wheel()
	}
	return 0, 0
}

// IsMouseButtonPressed wraps ebiten.IsMouseButtonPressed, respecting the current mode.
// Panic when using the patch integration method.
func IsMouseButtonPressed(button ebiten.MouseButton) bool {
	panicIfPatched()
	switch currentMode {
	case InjectionOnly:
		return integrate.IsMouseButtonPressed(integrate.MouseButton(button))
	case Passthrough:
		return ebiten.IsMouseButtonPressed(button)
	case InjectionFallback:
		if integrate.IsMouseButtonPressed(integrate.MouseButton(button)) {
			return true
		}
		return ebiten.IsMouseButtonPressed(button)
	}
	return false
}

// IsKeyJustPressed wraps inpututil.IsKeyJustPressed, respecting the current mode.
// Panic when using the patch integration method.
func IsKeyJustPressed(key ebiten.Key) bool {
	panicIfPatched()
	switch currentMode {
	case InjectionOnly:
		return integrate.IsKeyJustPressed(integrate.Key(key))
	case Passthrough:
		return inpututil.IsKeyJustPressed(key)
	case InjectionFallback:
		if integrate.IsKeyJustPressed(integrate.Key(key)) {
			return true
		}
		return inpututil.IsKeyJustPressed(key)
	}
	return false
}

// IsKeyJustReleased wraps inpututil.IsKeyJustReleased, respecting the current mode.
// Panic when using the patch integration method.
func IsKeyJustReleased(key ebiten.Key) bool {
	panicIfPatched()
	switch currentMode {
	case InjectionOnly:
		return integrate.IsKeyJustReleased(integrate.Key(key))
	case Passthrough:
		return inpututil.IsKeyJustReleased(key)
	case InjectionFallback:
		if integrate.IsKeyJustReleased(integrate.Key(key)) {
			return true
		}
		return inpututil.IsKeyJustReleased(key)
	}
	return false
}

// KeyPressDuration wraps inpututil.KeyPressDuration, respecting the current mode.
// Panic when using the patch integration method.
func KeyPressDuration(key ebiten.Key) int {
	panicIfPatched()
	switch currentMode {
	case InjectionOnly:
		return integrate.KeyPressDuration(integrate.Key(key))
	case Passthrough:
		return inpututil.KeyPressDuration(key)
	case InjectionFallback:
		d := integrate.KeyPressDuration(integrate.Key(key))
		if d > 0 {
			return d
		}
		return inpututil.KeyPressDuration(key)
	}
	return 0
}

// IsMouseButtonJustPressed wraps inpututil.IsMouseButtonJustPressed, respecting the current mode.
// Panic when using the patch integration method.
func IsMouseButtonJustPressed(button ebiten.MouseButton) bool {
	panicIfPatched()
	switch currentMode {
	case InjectionOnly:
		return integrate.IsMouseButtonJustPressed(integrate.MouseButton(button))
	case Passthrough:
		return inpututil.IsMouseButtonJustPressed(button)
	case InjectionFallback:
		if integrate.IsMouseButtonJustPressed(integrate.MouseButton(button)) {
			return true
		}
		return inpututil.IsMouseButtonJustPressed(button)
	}
	return false
}

// IsMouseButtonJustReleased wraps inpututil.IsMouseButtonJustReleased, respecting the current mode.
// Panic when using the patch integration method.
func IsMouseButtonJustReleased(button ebiten.MouseButton) bool {
	panicIfPatched()
	switch currentMode {
	case InjectionOnly:
		return integrate.IsMouseButtonJustReleased(integrate.MouseButton(button))
	case Passthrough:
		return inpututil.IsMouseButtonJustReleased(button)
	case InjectionFallback:
		if integrate.IsMouseButtonJustReleased(integrate.MouseButton(button)) {
			return true
		}
		return inpututil.IsMouseButtonJustReleased(button)
	}
	return false
}

// MouseButtonPressDuration wraps inpututil.MouseButtonPressDuration, respecting the current mode.
// Panic when using the patch integration method.
func MouseButtonPressDuration(button ebiten.MouseButton) int {
	panicIfPatched()
	switch currentMode {
	case InjectionOnly:
		return integrate.MouseButtonPressDuration(integrate.MouseButton(button))
	case Passthrough:
		return inpututil.MouseButtonPressDuration(button)
	case InjectionFallback:
		d := integrate.MouseButtonPressDuration(integrate.MouseButton(button))
		if d > 0 {
			return d
		}
		return inpututil.MouseButtonPressDuration(button)
	}
	return 0
}
