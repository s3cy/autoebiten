package autoui

import (
	"image/color"
	"time"
)

// Config holds autoui configuration.
type Config struct {
	HighlightDuration time.Duration
	HighlightColor    color.Color
}

// DefaultConfig returns default configuration.
func DefaultConfig() Config {
	return Config{
		HighlightDuration: 3 * time.Second,
		HighlightColor:    color.RGBA{255, 0, 0, 255},
	}
}

var config = DefaultConfig()

// SetConfig updates autoui configuration.
func SetConfig(c Config) {
	config = c
	SetHighlightDuration(c.HighlightDuration)
}