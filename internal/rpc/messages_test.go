package rpc

import "testing"

func TestErrorCodes(t *testing.T) {
	tests := []struct {
		code     int
		expected int
	}{
		{ErrConnectionFailed, -32000},
		{ErrInvalidParams, -32001},
		{ErrScriptFailed, -32002},
		{ErrScreenshotFailed, -32003},
		{ErrGameNotRunning, -32004},
		{ErrGameNotConnected, -32005},
	}

	for _, tt := range tests {
		if tt.code != tt.expected {
			t.Errorf("Expected error code %d, got %d", tt.expected, tt.code)
		}
	}
}
