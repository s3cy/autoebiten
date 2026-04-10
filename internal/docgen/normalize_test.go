package docgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		rules    []NormalizeRule
		expected string
	}{
		{
			name:     "no rules",
			input:    "hello world",
			rules:    nil,
			expected: "hello world",
		},
		{
			name:     "single rule - replace PID",
			input:    "process PID=12345 started",
			rules:    []NormalizeRule{{Pattern: `PID=\d+`, Replace: "PID=<PID>"}},
			expected: "process PID=<PID> started",
		},
		{
			name:  "multiple rules",
			input: "PID=12345 at 2024-01-15",
			rules: []NormalizeRule{
				{Pattern: `PID=\d+`, Replace: "PID=<PID>"},
				{Pattern: `\d{4}-\d{2}-\d{2}`, Replace: "<TIMESTAMP>"},
			},
			expected: "PID=<PID> at <TIMESTAMP>",
		},
		{
			name:     "trims whitespace",
			input:    "  hello world  ",
			rules:    nil,
			expected: "hello world",
		},
		{
			name:     "complex pattern - quoted address",
			input:    `server listening on _addr="127.0.0.1:8080"`,
			rules:    []NormalizeRule{{Pattern: `_addr="[^"]*"`, Replace: `_addr="<ADDR>"`}},
			expected: `server listening on _addr="<ADDR>"`,
		},
		{
			name:     "regex with timestamp",
			input:    "event at 2024-01-15 10:30:45.123 occurred",
			rules:    []NormalizeRule{{Pattern: `\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+`, Replace: "<TIMESTAMP>"}},
			expected: "event at <TIMESTAMP> occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Normalize(tt.input, tt.rules)
			assert.Equal(t, tt.expected, result)
		})
	}
}