package script

import (
	"encoding/json"
	"fmt"
	"os"
)

// Parse reads and parses a script file.
func Parse(path string) (*Script, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return ParseBytes(data)
}

// ParseString parses a script from a JSON string.
func ParseString(s string) (*Script, error) {
	return ParseBytes([]byte(s))
}

// ParseBytes parses a script from JSON bytes.
func ParseBytes(data []byte) (*Script, error) {
	var script Script
	if err := json.Unmarshal(data, &script); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if script.Version == "" {
		script.Version = "1.0" // Default to 1.0 if not specified
	}

	return &script, nil
}
