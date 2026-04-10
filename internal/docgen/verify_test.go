package docgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyOutputs(t *testing.T) {
	// All identical - should pass
	err := VerifyOutputs("OK: running", "OK: running", "OK: running")
	assert.NoError(t, err)
}

func TestVerifyOutputsMismatch(t *testing.T) {
	// Different outputs - should fail
	err := VerifyOutputs("OK: running", "Error: failed", "OK: running")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "output mismatch")
}

func TestVerifyOutputsEmpty(t *testing.T) {
	// No outputs - should pass
	err := VerifyOutputs()
	assert.NoError(t, err)
}