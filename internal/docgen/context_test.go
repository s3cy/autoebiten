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

func TestContextAddOutput(t *testing.T) {
	ctx := NewContext()

	ctx.AddOutput("first output")
	ctx.AddOutput("second output")
	ctx.AddOutput("third output")

	// Verify outputs were added (we can't access outputs directly, but GetOutputs returns them)
	outputs := ctx.GetOutputs()
	assert.Len(t, outputs, 3)
	assert.Equal(t, "first output", outputs[0])
	assert.Equal(t, "second output", outputs[1])
	assert.Equal(t, "third output", outputs[2])
}

func TestContextGetOutputsReturnsCopy(t *testing.T) {
	ctx := NewContext()
	ctx.AddOutput("original")

	// Get outputs and modify the returned slice
	outputs := ctx.GetOutputs()
	outputs[0] = "modified"

	// Verify that the internal state was not changed
	originalOutputs := ctx.GetOutputs()
	assert.Equal(t, "original", originalOutputs[0],
		"GetOutputs should return a copy, not the internal slice")
}

func TestContextGetOutputsEmpty(t *testing.T) {
	ctx := NewContext()
	outputs := ctx.GetOutputs()
	assert.Empty(t, outputs)
}
