package autoebiten

// Mode represents the input handling mode.
type Mode int

const (
	// InjectionOnly mode returns only injected input results.
	InjectionOnly Mode = iota
	// InjectionFallback mode returns injected results if available,
	// otherwise falls back to ebiten's native input handling.
	InjectionFallback
	// Passthrough mode only uses ebiten's native input handling.
	Passthrough
)

var currentMode = InjectionFallback

// SetMode sets the input handling mode.
func SetMode(mode Mode) {
	currentMode = mode
}

// GetMode returns the current input handling mode.
func GetMode() Mode {
	return currentMode
}
