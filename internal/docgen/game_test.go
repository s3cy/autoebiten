package docgen

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDelay(t *testing.T) {
	start := time.Now()
	Delay("100ms")
	elapsed := time.Since(start)

	assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(100))
}