package docgen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
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

func TestLoadConfigNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/config.yaml")
	assert.Error(t, err)
}

func TestNormalize(t *testing.T) {
	rules := []NormalizeRule{
		{Pattern: `_addr="[^"]*"`, Replace: `_addr="<ADDR>"`},
		{Pattern: `\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+`, Replace: `<TIMESTAMP>`},
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normalize address",
			input:    `<Button _addr="0x14000123456" id="btn"/>`,
			expected: `<Button _addr="<ADDR>" id="btn"/>`,
		},
		{
			name:     "normalize timestamp",
			input:    `2026-04-09 14:30:15.123 game started`,
			expected: `<TIMESTAMP> game started`,
		},
		{
			name:     "normalize multiple",
			input:    `<Button _addr="0x14000123456"/> at 2026-04-09 14:30:15.123`,
			expected: `<Button _addr="<ADDR>"/> at <TIMESTAMP>`,
		},
		{
			name:     "no matches",
			input:    `plain text`,
			expected: `plain text`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Normalize(tt.input, rules)
			assert.Equal(t, tt.expected, result)
		})
	}
}
