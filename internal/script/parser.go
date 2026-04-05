package script

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

// regex to match // single-line and /* */ multi-line comments
var (
	singleLineComment = regexp.MustCompile(`(?m)//.*$`)
	multiLineComment  = regexp.MustCompile(`(?s)/\*.*?\*/`)
)

// stripComments removes // and /* */ style comments from JSONC.
// Note: This simple implementation may incorrectly strip // inside string values
// (e.g., URLs like "http://example.com"). For such cases, use the URL-encoded
// form or place the value in an external config file.
func stripComments(data []byte) []byte {
	// Remove multi-line comments first (/* */)
	data = multiLineComment.ReplaceAll(data, nil)
	// Then remove single-line comments (//)
	data = singleLineComment.ReplaceAll(data, nil)
	return data
}

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

// ParseBytes parses a script from JSONC bytes (JSON with Comments).
func ParseBytes(data []byte) (*Script, error) {
	// Strip comments before parsing
	data = stripComments(data)

	var script Script
	if err := json.Unmarshal(data, &script); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if script.Version == "" {
		script.Version = "1.0" // Default to 1.0 if not specified
	}

	return &script, nil
}
