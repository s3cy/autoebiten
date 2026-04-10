package docgen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigDefaults(t *testing.T) {
	cfg := &Config{}
	assert.Equal(t, "", cfg.GameDir)
	assert.Empty(t, cfg.Normalize)
}

func TestNormalizeRule(t *testing.T) {
	rules := []NormalizeRule{
		{Pattern: `PID=\d+`, Replace: "PID=<PID>"},
		{Pattern: `\d{4}-\d{2}-\d{2}`, Replace: "<TIMESTAMP>"},
	}
	cfg := &Config{Normalize: rules}
	assert.Len(t, cfg.Normalize, 2)
	assert.Equal(t, `PID=\d+`, cfg.Normalize[0].Pattern)
}

// Legacy tests for LoadConfig (will be removed in Task 14)
func TestLoadConfigLegacy(t *testing.T) {
	// Create temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `game_dir: examples/autoui
normalize:
  - pattern: '_addr="[^"]*"'
    replace: '_addr="<ADDR>"'
  - pattern: '\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+'
    replace: '<TIMESTAMP>'
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	config, err := LoadConfig(configPath)
	require.NoError(t, err)

	assert.Equal(t, "examples/autoui", config.GameDir)
	assert.Len(t, config.Normalize, 2)
	assert.Equal(t, `_addr="[^"]*"`, config.Normalize[0].Pattern)
	assert.Equal(t, "_addr=\"<ADDR>\"", config.Normalize[0].Replace)
}

func TestLoadConfigNotFoundLegacy(t *testing.T) {
	_, err := LoadConfig("/nonexistent/config.yaml")
	assert.Error(t, err)
}

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
