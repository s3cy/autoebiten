package testkit

import (
	"encoding/json"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/s3cy/autoebiten/internal/rpc"
)

// TestMouseActionConstants verifies that mouse action constants have the expected values.
// This prevents accidental changes to the action strings that would break RPC communication.
func TestMouseActionConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"position", mouseActionPosition, "position"},
		{"press", mouseActionPress, "press"},
		{"release", mouseActionRelease, "release"},
		{"hold", mouseActionHold, "hold"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("mouse action constant %q = %q, want %q", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

// TestMouseActionValues validates that all expected mouse action values are defined.
// This test guards against:
// 1. Typos in action strings (e.g., "move" instead of "position")
// 2. Missing action constants
// 3. Incorrect action values that would cause RPC failures
func TestMouseActionValues(t *testing.T) {
	// Valid mouse actions as defined by the RPC protocol
	validActions := map[string]bool{
		"position": true,
		"press":    true,
		"release":  true,
		"hold":     true,
	}

	// Verify all constants match valid actions
	actions := []string{
		mouseActionPosition,
		mouseActionPress,
		mouseActionRelease,
		mouseActionHold,
	}

	for _, action := range actions {
		if !validActions[action] {
			t.Errorf("mouse action %q is not a valid RPC action", action)
		}
	}
}

// TestMouseParamsSerialization verifies that MouseParams serializes correctly
// with the action constants.
func TestMouseParamsSerialization(t *testing.T) {
	tests := []struct {
		name   string
		params rpc.MouseParams
	}{
		{
			name:   "position action",
			params: rpc.MouseParams{Action: mouseActionPosition, X: 100, Y: 200},
		},
		{
			name:   "press action",
			params: rpc.MouseParams{Action: mouseActionPress, Button: "MouseButtonLeft"},
		},
		{
			name:   "release action",
			params: rpc.MouseParams{Action: mouseActionRelease, Button: "MouseButtonRight"},
		},
		{
			name:   "hold action",
			params: rpc.MouseParams{Action: mouseActionHold, Button: "MouseButtonLeft", DurationTicks: 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.params)
			if err != nil {
				t.Fatalf("failed to marshal MouseParams: %v", err)
			}

			// Verify it round-trips correctly
			var result rpc.MouseParams
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("failed to unmarshal MouseParams: %v", err)
			}

			if result.Action != tt.params.Action {
				t.Errorf("action mismatch: got %q, want %q", result.Action, tt.params.Action)
			}
		})
	}
}

// TestMouseButtonToString validates mouse button conversion returns valid button names.
func TestMouseButtonToString(t *testing.T) {
	buttons := []ebiten.MouseButton{
		ebiten.MouseButtonLeft,
		ebiten.MouseButtonRight,
		ebiten.MouseButtonMiddle,
	}

	validButtons := map[string]bool{
		"MouseButtonLeft":   true,
		"MouseButtonRight":  true,
		"MouseButtonMiddle": true,
		"MouseButton0":      true,
		"MouseButton1":      true,
		"MouseButton2":      true,
	}

	for _, btn := range buttons {
		result := mouseButtonToString(btn)
		if !validButtons[result] {
			t.Errorf("mouseButtonToString(%v) returned invalid button name: %q", btn, result)
		}
	}
}

// TestMouseActionPositionSpecifically guards against the specific bug where
// MoveMouse was using "move" instead of "position".
func TestMouseActionPositionSpecifically(t *testing.T) {
	// This test specifically guards against the bug fixed in:
	// "move" should really be "position"
	if mouseActionPosition != "position" {
		t.Errorf("mouseActionPosition = %q, want %q - this would break MoveMouse()", mouseActionPosition, "position")
	}
}
