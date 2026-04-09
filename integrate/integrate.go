package integrate

import (
	"image"

	"github.com/s3cy/autoebiten/internal/input"
	"github.com/s3cy/autoebiten/internal/server"
)

type Key = input.Key
type MouseButton = input.MouseButton

// IsPatched indicates whether the game uses a patched version of Ebiten.
var IsPatched = false

// drawHighlightsFunc is the registered callback for drawing highlight overlays.
var drawHighlightsFunc func(screen image.Image)

// RegisterDrawHighlights registers a callback for drawing highlight overlays.
// The callback receives screen as image.Image; caller must type assert to *ebiten.Image.
// Called by autoui during Register() to enable highlights in patch method.
func RegisterDrawHighlights(fn func(screen image.Image)) {
	drawHighlightsFunc = fn
}

// AfterDraw handles all post-draw operations after game.Draw() completes.
// Currently handles: screenshots (ProcessScreenshots) and highlight overlays.
// Called by patched ebiten in DrawOffscreen().
func AfterDraw(screen image.Image) {
	Capture(screen)
	if drawHighlightsFunc != nil {
		drawHighlightsFunc(screen)
	}
}

// Capture processes screenshots for injection.
func Capture(screen image.Image) {
	server.ProcessScreenshots(screen)
}

// Update runs the internal update loop.
func Update() bool {
	return server.Update()
}

// IsKeyPressed returns whether the key is pressed, respecting the current mode.
func IsKeyPressed(key Key) bool {
	return input.Get().IsKeyPressed(key, server.Tick())
}

// CursorPosition returns the cursor position, respecting the current mode.
func CursorPosition() (x, y int) {
	return input.Get().CursorPosition()
}

// Wheel returns the mouse wheel scroll amount, respecting the current mode.
func Wheel() (x, y float64) {
	return input.Get().Wheel()
}

// IsMouseButtonPressed returns whether the mouse button is pressed,
// respecting the current mode.
func IsMouseButtonPressed(button MouseButton) bool {
	return input.Get().IsMouseButtonPressed(button, server.Tick())
}

func IsKeyJustPressed(key Key) bool {
	return input.Get().IsKeyJustPressed(key, server.Tick())
}

func IsKeyJustReleased(key Key) bool {
	return input.Get().IsKeyJustReleased(key, server.Tick())
}

func KeyPressDuration(key Key) int {
	return int(input.Get().KeyPressDuration(key, server.Tick()))
}

func IsMouseButtonJustPressed(button MouseButton) bool {
	return input.Get().IsMouseButtonJustPressed(button, server.Tick())
}

func IsMouseButtonJustReleased(button MouseButton) bool {
	return input.Get().IsMouseButtonJustReleased(button, server.Tick())
}

func MouseButtonPressDuration(button MouseButton) int {
	return int(input.Get().MouseButtonPressDuration(button, server.Tick()))
}
