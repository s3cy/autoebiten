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
