package cli

import (
	"testing"
)

func TestRunStateCommand(t *testing.T) {
	// Note: This is a basic structure test. Full integration tests
	// require a running game with registered state exporter.
	//
	// The state command internally calls the custom command with
	// ".state.<name>" prefix, which is tested via integration tests.

	executor := NewCommandExecutor()
	if executor == nil {
		t.Fatal("NewCommandExecutor returned nil")
	}

	// Test that the method exists and has correct signature
	// This should fail to compile if RunStateCommand is not defined
	err := executor.RunStateCommand("test", "path.to.value")
	_ = err // We expect this to fail (no game connected) but method should exist
}