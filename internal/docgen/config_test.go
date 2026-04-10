package docgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

