package integrate

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAfterDrawNoCallback tests AfterDraw when no callback is registered.
func TestAfterDrawNoCallback(t *testing.T) {
	// Reset callback to nil
	drawHighlightsFunc = nil

	// Create a dummy image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// AfterDraw should not panic when no callback is registered
	AfterDraw(img)
	// No assertion needed - just verify it doesn't panic
}

// TestAfterDrawWithCallback tests AfterDraw invokes registered callback.
func TestAfterDrawWithCallback(t *testing.T) {
	// Track if callback was invoked
	callbackInvoked := false
	var receivedImage image.Image

	// Register a test callback
	RegisterDrawHighlights(func(screen image.Image) {
		callbackInvoked = true
		receivedImage = screen
	})

	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// Call AfterDraw
	AfterDraw(img)

	// Verify callback was invoked with correct image
	assert.True(t, callbackInvoked, "Callback should be invoked")
	assert.Equal(t, img, receivedImage, "Callback should receive the same image")

	// Reset callback
	drawHighlightsFunc = nil
}

// TestRegisterDrawHighlights tests callback registration.
func TestRegisterDrawHighlights(t *testing.T) {
	// Reset
	drawHighlightsFunc = nil

	// Register first callback
	callback1 := func(screen image.Image) {}
	RegisterDrawHighlights(callback1)

	// Verify it's registered
	assert.NotNil(t, drawHighlightsFunc, "Callback should be registered")

	// Register second callback (should replace first)
	callback2 := func(screen image.Image) {}
	RegisterDrawHighlights(callback2)

	// Verify second callback is registered (simple pointer comparison won't work,
	// but we can verify it's not nil)
	assert.NotNil(t, drawHighlightsFunc, "Second callback should be registered")

	// Reset
	drawHighlightsFunc = nil
}