package autoui

import (
	"image"
	"image/color"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Default highlight configuration
const (
	defaultHighlightDuration = 3 * time.Second
)

var (
	defaultHighlightColor = color.RGBA{255, 0, 0, 255} // Red
)

// Highlight represents a visual highlight rectangle for debugging.
type Highlight struct {
	// Rect is the screen rectangle to highlight.
	Rect image.Rectangle

	// Color is the highlight color.
	Color color.Color

	// ExpiresAt is the time when this highlight should be removed.
	ExpiresAt time.Time
}

// highlightManager manages visual highlight rectangles with expiry.
type highlightManager struct {
	mu         sync.RWMutex
	highlights []Highlight
	duration   time.Duration
	color      color.Color
}

// newHighlightManager creates a new highlight manager with default settings.
func newHighlightManager() *highlightManager {
	return &highlightManager{
		highlights: make([]Highlight, 0),
		duration:   defaultHighlightDuration,
		color:      defaultHighlightColor,
	}
}

// add creates a new highlight for the given rectangle with the configured duration and color.
func (m *highlightManager) add(rect image.Rectangle) {
	m.mu.Lock()
	defer m.mu.Unlock()

	h := Highlight{
		Rect:      rect,
		Color:     m.color,
		ExpiresAt: time.Now().Add(m.duration),
	}
	m.highlights = append(m.highlights, h)
}

// clear removes all highlights immediately.
func (m *highlightManager) clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.highlights = make([]Highlight, 0)
}

// getActive returns all non-expired highlights and removes expired ones from the internal slice.
func (m *highlightManager) getActive() []Highlight {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()

	// Filter out expired highlights
	active := make([]Highlight, 0, len(m.highlights))
	for _, h := range m.highlights {
		if h.ExpiresAt.After(now) {
			active = append(active, h)
		}
	}

	// Update internal slice to remove expired
	m.highlights = active

	return active
}

// draw renders all active highlights onto the screen using vector.StrokeRect.
func (m *highlightManager) draw(screen *ebiten.Image) {
	active := m.getActive()

	strokeWidth := float32(2.0)
	antialias := true

	for _, h := range active {
		vector.StrokeRect(
			screen,
			float32(h.Rect.Min.X),
			float32(h.Rect.Min.Y),
			float32(h.Rect.Dx()),
			float32(h.Rect.Dy()),
			strokeWidth,
			h.Color,
			antialias,
		)
	}
}

// Global highlight manager instance
var globalHighlightManager = newHighlightManager()

// drawHighlightsCallback is the callback registered with integrate for patch method.
// It type asserts image.Image to *ebiten.Image for vector drawing.
func drawHighlightsCallback(screen image.Image) {
	if ebiScreen, ok := screen.(*ebiten.Image); ok {
		globalHighlightManager.draw(ebiScreen)
	}
}

// DrawHighlights draws all active highlights onto the screen.
// This is the public API to be called in the game's Draw method after ui.Draw(screen).
func DrawHighlights(screen *ebiten.Image) {
	globalHighlightManager.draw(screen)
}

// SetHighlightDuration configures the duration for new highlights.
func SetHighlightDuration(d time.Duration) {
	globalHighlightManager.mu.Lock()
	defer globalHighlightManager.mu.Unlock()

	globalHighlightManager.duration = d
}

// AddHighlight adds a highlight for the given rectangle using the global manager.
func AddHighlight(rect image.Rectangle) {
	globalHighlightManager.add(rect)
}

// ClearHighlights removes all highlights immediately using the global manager.
func ClearHighlights() {
	globalHighlightManager.clear()
}