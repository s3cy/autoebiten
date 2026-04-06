package recording

import (
	"fmt"

	"github.com/s3cy/autoebiten/internal/script"
)

// Generator creates scripts from recording entries.
type Generator struct {
	speed float64 // 1.0 = original, 2.0 = 2x speed, 0.5 = half speed
}

// NewGenerator creates a generator with speed multiplier.
func NewGenerator(speed float64) *Generator {
	return &Generator{speed: speed}
}

// Generate creates a Script from recording entries.
func (g *Generator) Generate(entries []Entry) (*script.Script, error) {
	if len(entries) == 0 {
		return nil, fmt.Errorf("no recording entries")
	}

	if g.speed <= 0 {
		return nil, fmt.Errorf("speed must be greater than 0")
	}

	s := &script.Script{
		Version:  "1.0",
		Commands: []script.CommandWrapper{},
	}

	// Add first command
	s.Commands = append(s.Commands, entries[0].Command)

	// Add delays between subsequent commands
	for i := 1; i < len(entries); i++ {
		delay := entries[i].Timestamp.Sub(entries[i-1].Timestamp)
		delayMs := int64(delay.Milliseconds()) / int64(g.speed)

		if delayMs > 0 {
			s.Commands = append(s.Commands, &script.DelayCmd{Ms: delayMs})
		}

		s.Commands = append(s.Commands, entries[i].Command)
	}

	return s, nil
}