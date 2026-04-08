package rpc

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/s3cy/autoebiten/internal/output"
)

func TestLaunchSocketPath(t *testing.T) {
	tests := []struct {
		name     string
		pid      int
		expected string
	}{
		{
			name:     "returns correct launch socket path for PID",
			pid:      12345,
			expected: filepath.Join(DefaultSocketDir, "autoebiten-12345-launch.sock"),
		},
		{
			name:     "returns correct launch socket path for PID 1",
			pid:      1,
			expected: filepath.Join(DefaultSocketDir, "autoebiten-1-launch.sock"),
		},
		{
			name:     "returns correct launch socket path for large PID",
			pid:      999999,
			expected: filepath.Join(DefaultSocketDir, "autoebiten-999999-launch.sock"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LaunchSocketPath(tt.pid)
			if result != tt.expected {
				t.Errorf("LaunchSocketPath(%d) = %s, want %s", tt.pid, result, tt.expected)
			}
		})
	}
}

func TestFindRunningGames_WithLaunchSockets(t *testing.T) {
	// Create a temporary directory for test sockets
	tempDir, err := os.MkdirTemp("", "autoebiten-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test socket files
	// Note: We use the current process PID so syscall.Kill returns success
	currentPID := os.Getpid()

	// Create a regular socket file
	regularSocket := filepath.Join(tempDir, fmt.Sprintf("autoebiten-%d.sock", currentPID))
	if err := os.WriteFile(regularSocket, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to create regular socket: %v", err)
	}

	// Create a launch socket file for the same PID
	launchSocket := filepath.Join(tempDir, fmt.Sprintf("autoebiten-%d-launch.sock", currentPID))
	if err := os.WriteFile(launchSocket, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to create launch socket: %v", err)
	}

	// Temporarily override the socket directory using env var
	originalSocket := os.Getenv("AUTOEBITEN_SOCKET")
	defer os.Setenv("AUTOEBITEN_SOCKET", originalSocket)

	// Set up to use our temp directory
	os.Setenv("AUTOEBITEN_SOCKET", filepath.Join(tempDir, "test.sock"))

	games, err := findRunningGames()
	if err != nil {
		t.Fatalf("findRunningGames() error = %v", err)
	}

	// Should find exactly one game (deduplicated), and it should prefer launch socket
	if len(games) != 1 {
		t.Errorf("findRunningGames() returned %d games, want 1", len(games))
	}

	if len(games) > 0 && games[0].PID != currentPID {
		t.Errorf("findRunningGames() returned PID %d, want %d", games[0].PID, currentPID)
	}
}

func TestFindRunningGames_DeduplicationPrefersLaunch(t *testing.T) {
	// This test verifies that when both regular and launch sockets exist,
	// only the launch socket is reported

	tempDir, err := os.MkdirTemp("", "autoebiten-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	currentPID := os.Getpid()

	// Create both socket types
	regularSocket := filepath.Join(tempDir, fmt.Sprintf("autoebiten-%d.sock", currentPID))
	launchSocket := filepath.Join(tempDir, fmt.Sprintf("autoebiten-%d-launch.sock", currentPID))

	if err := os.WriteFile(regularSocket, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to create regular socket: %v", err)
	}
	if err := os.WriteFile(launchSocket, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to create launch socket: %v", err)
	}

	// Temporarily override the socket directory using env var
	originalSocket := os.Getenv("AUTOEBITEN_SOCKET")
	defer os.Setenv("AUTOEBITEN_SOCKET", originalSocket)

	os.Setenv("AUTOEBITEN_SOCKET", filepath.Join(tempDir, "test.sock"))

	games, err := findRunningGames()
	if err != nil {
		t.Fatalf("findRunningGames() error = %v", err)
	}

	// Should be deduplicated to 1
	if len(games) != 1 {
		t.Errorf("Expected 1 game after deduplication, got %d", len(games))
	}
}

func TestFindRunningGames_RemovesStaleSockets(t *testing.T) {
	// Create a temporary directory for test sockets
	tempDir, err := os.MkdirTemp("", "autoebiten-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Use a PID that definitely doesn't exist (very high number)
	deadPID := 99999999

	// Create stale socket files for the dead PID
	regularSocket := filepath.Join(tempDir, fmt.Sprintf("autoebiten-%d.sock", deadPID))
	launchSocket := filepath.Join(tempDir, fmt.Sprintf("autoebiten-%d-launch.sock", deadPID))

	if err := os.WriteFile(regularSocket, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to create regular socket: %v", err)
	}
	if err := os.WriteFile(launchSocket, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to create launch socket: %v", err)
	}

	// Verify files exist
	if _, err := os.Stat(regularSocket); err != nil {
		t.Fatalf("Regular socket should exist before test")
	}
	if _, err := os.Stat(launchSocket); err != nil {
		t.Fatalf("Launch socket should exist before test")
	}

	// Temporarily override the socket directory using env var
	originalSocket := os.Getenv("AUTOEBITEN_SOCKET")
	defer os.Setenv("AUTOEBITEN_SOCKET", originalSocket)

	os.Setenv("AUTOEBITEN_SOCKET", filepath.Join(tempDir, "test.sock"))

	games, err := findRunningGames()
	if err != nil {
		t.Fatalf("findRunningGames() error = %v", err)
	}

	// Should find no games (dead process)
	if len(games) != 0 {
		t.Errorf("Expected 0 games for dead process, got %d", len(games))
	}

	// Both socket files should be cleaned up
	if _, err := os.Stat(regularSocket); !os.IsNotExist(err) {
		t.Error("Regular socket should have been removed")
	}
	if _, err := os.Stat(launchSocket); !os.IsNotExist(err) {
		t.Error("Launch socket should have been removed")
	}
}

func TestFindRunningGames_MultipleProcesses(t *testing.T) {
	// Test with multiple PIDs - only current process should be found
	tempDir, err := os.MkdirTemp("", "autoebiten-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	currentPID := os.Getpid()
	deadPID1 := 99999997
	deadPID2 := 99999998

	// Create sockets for current process (both types)
	regularSocket := filepath.Join(tempDir, fmt.Sprintf("autoebiten-%d.sock", currentPID))
	launchSocket := filepath.Join(tempDir, fmt.Sprintf("autoebiten-%d-launch.sock", currentPID))
	os.WriteFile(regularSocket, []byte{}, 0644)
	os.WriteFile(launchSocket, []byte{}, 0644)

	// Create sockets for dead processes
	deadRegular1 := filepath.Join(tempDir, fmt.Sprintf("autoebiten-%d.sock", deadPID1))
	deadLaunch1 := filepath.Join(tempDir, fmt.Sprintf("autoebiten-%d-launch.sock", deadPID1))
	os.WriteFile(deadRegular1, []byte{}, 0644)
	os.WriteFile(deadLaunch1, []byte{}, 0644)

	deadRegular2 := filepath.Join(tempDir, fmt.Sprintf("autoebiten-%d.sock", deadPID2))
	os.WriteFile(deadRegular2, []byte{}, 0644)

	// Temporarily override the socket directory using env var
	originalSocket := os.Getenv("AUTOEBITEN_SOCKET")
	defer os.Setenv("AUTOEBITEN_SOCKET", originalSocket)

	os.Setenv("AUTOEBITEN_SOCKET", filepath.Join(tempDir, "test.sock"))

	games, err := findRunningGames()
	if err != nil {
		t.Fatalf("findRunningGames() error = %v", err)
	}

	// Should find only the current process
	if len(games) != 1 {
		t.Errorf("Expected 1 game (current process), got %d", len(games))
	}

	if len(games) > 0 && games[0].PID != currentPID {
		t.Errorf("Expected PID %d, got %d", currentPID, games[0].PID)
	}

	// Dead process sockets should be cleaned up
	if _, err := os.Stat(deadRegular1); !os.IsNotExist(err) {
		t.Error("Dead process regular socket should have been removed")
	}
	if _, err := os.Stat(deadLaunch1); !os.IsNotExist(err) {
		t.Error("Dead process launch socket should have been removed")
	}
	if _, err := os.Stat(deadRegular2); !os.IsNotExist(err) {
		t.Error("Dead process regular socket 2 should have been removed")
	}
}

func TestOutputDerivePaths_LaunchSocket(t *testing.T) {
	// Verify that output.DerivePaths correctly generates launch socket paths
	socketPath := filepath.Join(DefaultSocketDir, "autoebiten-12345.sock")
	paths := output.DerivePaths(socketPath)

	expectedLaunchSock := filepath.Join(DefaultSocketDir, "autoebiten-12345-launch.sock")
	if paths.LaunchSock != expectedLaunchSock {
		t.Errorf("DerivePaths().LaunchSock = %s, want %s", paths.LaunchSock, expectedLaunchSock)
	}
}
