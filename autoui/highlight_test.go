package autoui

import (
	"image"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestAddHighlight tests that adding a highlight creates an entry with proper expiry.
func TestAddHighlight(t *testing.T) {
	// Create a fresh manager for testing
	m := newHighlightManager()

	// Add a highlight
	rect := image.Rect(10, 20, 110, 50)
	m.add(rect)

	// Get active highlights - should have 1
	active := m.getActive()
	assert.Len(t, active, 1, "Expected 1 active highlight")

	// Verify the highlight properties
	h := active[0]
	assert.Equal(t, rect, h.Rect, "Rect should match")
	assert.Equal(t, defaultHighlightColor, h.Color, "Color should be default red")

	// Verify expiry time is approximately 3 seconds from now
	expectedExpiry := time.Now().Add(defaultHighlightDuration)
	// Allow 100ms tolerance for test execution time
	assert.WithinDuration(t, expectedExpiry, h.ExpiresAt, 100*time.Millisecond,
		"ExpiresAt should be approximately 3 seconds from now")
}

// TestClearHighlights tests that clear removes all highlights.
func TestClearHighlights(t *testing.T) {
	m := newHighlightManager()

	// Add multiple highlights
	m.add(image.Rect(10, 10, 110, 40))
	m.add(image.Rect(120, 10, 220, 40))
	m.add(image.Rect(230, 10, 330, 40))

	// Verify we have 3 highlights
	active := m.getActive()
	assert.Len(t, active, 3, "Expected 3 highlights before clear")

	// Clear all highlights
	m.clear()

	// Verify all are removed
	active = m.getActive()
	assert.Len(t, active, 0, "Expected 0 highlights after clear")
}

// TestHighlightExpiry tests that expired highlights are filtered out.
func TestHighlightExpiry(t *testing.T) {
	m := newHighlightManager()

	// Set a very short duration for testing
	m.duration = 50 * time.Millisecond

	// Add a highlight
	m.add(image.Rect(10, 10, 110, 40))

	// Should be active immediately
	active := m.getActive()
	assert.Len(t, active, 1, "Expected 1 active highlight immediately after add")

	// Wait for expiry
	time.Sleep(100 * time.Millisecond)

	// Should be filtered out now
	active = m.getActive()
	assert.Len(t, active, 0, "Expected 0 active highlights after expiry")
}

// TestHighlightExpiryMixed tests that only expired highlights are removed.
func TestHighlightExpiryMixed(t *testing.T) {
	m := newHighlightManager()
	m.duration = 100 * time.Millisecond

	// Add first highlight (will expire quickly)
	m.add(image.Rect(10, 10, 110, 40))

	// Wait half the duration
	time.Sleep(50 * time.Millisecond)

	// Add second highlight (still active)
	m.add(image.Rect(120, 10, 220, 40))

	// Both should be active now
	active := m.getActive()
	assert.Len(t, active, 2, "Expected 2 active highlights before expiry")

	// Wait for first to expire (but not second)
	time.Sleep(60 * time.Millisecond)

	// Only second should remain
	active = m.getActive()
	assert.Len(t, active, 1, "Expected 1 active highlight after first expires")
	if len(active) > 0 {
		assert.Equal(t, image.Rect(120, 10, 220, 40), active[0].Rect,
			"Remaining highlight should be the second one")
	}
}

// TestSetHighlightDuration tests configuring the highlight duration.
func TestSetHighlightDuration(t *testing.T) {
	// Reset to default for clean test
	m := newHighlightManager()

	// Verify default
	assert.Equal(t, defaultHighlightDuration, m.duration,
		"Default duration should be 3 seconds")

	// Set custom duration
	customDuration := 5 * time.Second
	SetHighlightDuration(customDuration)

	// Verify it's set on global manager
	assert.Equal(t, customDuration, globalHighlightManager.duration,
		"Duration should be updated to 5 seconds")
}

// TestHighlightThreadSafety tests concurrent access to the highlight manager.
func TestHighlightThreadSafety(t *testing.T) {
	m := newHighlightManager()

	// Run concurrent operations
	done := make(chan bool)

	// Writer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			m.add(image.Rect(i, i, i+100, i+50))
		}
		done <- true
	}()

	// Reader goroutine
	go func() {
		for i := 0; i < 100; i++ {
			_ = m.getActive()
		}
		done <- true
	}()

	// Clear goroutine
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(10 * time.Millisecond)
			m.clear()
		}
		done <- true
	}()

	// Wait for all goroutines
	for i := 0; i < 3; i++ {
		<-done
	}

	// Test completed without race (race detector will catch issues)
	t.Log("Concurrent test completed without race conditions")
}

// TestHighlightColor tests that highlight color is properly set.
func TestHighlightColor(t *testing.T) {
	m := newHighlightManager()

	m.add(image.Rect(0, 0, 100, 50))

	active := m.getActive()
	assert.Len(t, active, 1)

	// Color should be the default red
	assert.Equal(t, defaultHighlightColor, active[0].Color)
}