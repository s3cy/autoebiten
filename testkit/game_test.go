package testkit

import (
	"os"
	"testing"
	"time"
)

// TestGameLaunchShutdown tests that Launch starts a game and Shutdown terminates it.
func TestGameLaunchShutdown(t *testing.T) {
	// Skip if no test binary available
	binaryPath := getSimpleTestGameBinary()
	if binaryPath == "" {
		t.Skip("test game binary not available")
	}

	game := Launch(t, binaryPath, WithTimeout(10*time.Second))

	// Verify game is responsive
	if err := game.Ping(); err != nil {
		t.Fatalf("game ping failed: %v", err)
	}

	// Shutdown
	if err := game.Shutdown(); err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}

	// Verify game is no longer responsive
	if err := game.Ping(); err == nil {
		t.Fatal("expected ping to fail after shutdown")
	}
}

// TestGameCleanup tests that t.Cleanup() properly shuts down the game.
func TestGameCleanup(t *testing.T) {
	// This test verifies cleanup happens - we can't really test
	// the cleanup itself without a subtest, but we verify Launch works
	binaryPath := getSimpleTestGameBinary()
	if binaryPath == "" {
		t.Skip("test game binary not available")
	}

	// Launch creates cleanup automatically
	_ = Launch(t, binaryPath, WithTimeout(5*time.Second))
	// Game will be cleaned up when test ends
}

// TestGameWaitFor tests the WaitFor helper.
func TestGameWaitFor(t *testing.T) {
	binaryPath := getSimpleTestGameBinary()
	if binaryPath == "" {
		t.Skip("test game binary not available")
	}

	game := Launch(t, binaryPath, WithTimeout(10*time.Second))
	defer game.Shutdown()

	// Test successful wait
	called := 0
	result := game.WaitFor(func() bool {
		called++
		return called >= 3
	}, 1*time.Second)

	if !result {
		t.Error("expected WaitFor to return true")
	}
	if called < 3 {
		t.Errorf("expected at least 3 calls, got %d", called)
	}

	// Test timeout
	result = game.WaitFor(func() bool {
		return false
	}, 100*time.Millisecond)

	if result {
		t.Error("expected WaitFor to return false on timeout")
	}
}

// TestGamePing tests the Ping method.
func TestGamePing(t *testing.T) {
	binaryPath := getSimpleTestGameBinary()
	if binaryPath == "" {
		t.Skip("test game binary not available")
	}

	game := Launch(t, binaryPath, WithTimeout(10*time.Second))
	defer game.Shutdown()

	// Multiple pings should succeed
	for i := 0; i < 3; i++ {
		if err := game.Ping(); err != nil {
			t.Fatalf("ping %d failed: %v", i+1, err)
		}
	}
}

// getSimpleTestGameBinary returns the path to the simple test game binary.
// Returns empty string if not available.
func getSimpleTestGameBinary() string {
	// Try to find the binary in common locations
	candidates := []string{
		"./internal/testgames/simple/simple",
		"./internal/testgames/simple/main",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	// Try to build it
	return ""
}
