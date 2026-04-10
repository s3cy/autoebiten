package docgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
	ctx := NewContext()
	assert.Nil(t, ctx.GameSession)
	assert.Nil(t, ctx.Config)
}

func TestContextSetConfig(t *testing.T) {
	ctx := NewContext()
	cfg := &Config{GameDir: "examples/test"}
	ctx.SetConfig(cfg)
	assert.Equal(t, "examples/test", ctx.Config.GameDir)
}