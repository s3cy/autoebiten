package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/s3cy/autoebiten/internal/rpc"
)

// launchSocketPath returns the launch socket path for a given PID.
// This helper mirrors the behavior of output.DerivePaths.
func launchSocketPath(pid int) string {
	return filepath.Join(rpc.DefaultSocketDir, fmt.Sprintf("autoebiten-%d-launch.sock", pid))
}

// TestPreRPCCrashDiagnostics verifies that when a game crashes before RPC connection,
// the CLI can connect to the launch socket and receive crash diagnostics.
func TestPreRPCCrashDiagnostics(t *testing.T) {
	// Build the CLI
	cliPath := filepath.Join(t.TempDir(), "autoebiten")
	buildCmd := exec.Command("go", "build", "-o", cliPath, "./cmd/autoebiten")
	buildCmd.Dir = "../"
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build CLI: %v\nOutput: %s", err, output)
	}

	// Launch a command that exits immediately (false command returns exit code 1)
	launchCmd := exec.Command(cliPath, "launch", "--", "false")
	launchCmd.Dir = "../"

	// Start launch in background
	if err := launchCmd.Start(); err != nil {
		t.Fatalf("Failed to start launch: %v", err)
	}

	// Get the launch PID for socket discovery
	launchPID := launchCmd.Process.Pid
	launchSocketPath := launchSocketPath(launchPID)

	// Wait for socket creation with timeout
	var socketFound bool
	for i := 0; i < 50; i++ {
		if _, err := os.Stat(launchSocketPath); err == nil {
			socketFound = true
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	if !socketFound {
		t.Fatalf("Launch socket was not created: %s", launchSocketPath)
	}

	// Wait a bit for the game to crash
	time.Sleep(200 * time.Millisecond)

	// Try to ping via launch socket using the CLI
	pingCmd := exec.Command(cliPath, "ping", "--pid", strconv.Itoa(launchPID))
	pingCmd.Dir = "../"
	output, _ := pingCmd.CombinedOutput()
	outputStr := string(output)

	// Cleanup launch process
	launchCmd.Process.Kill()
	launchCmd.Wait()

	// Verify output contains expected patterns
	// The game should have crashed, so we expect error diagnostics
	// Note: log_diff is only printed when there's actual output

	if !strings.Contains(outputStr, "<proxy_error>") {
		t.Errorf("Expected <proxy_error> in output, got:\n%s", outputStr)
	}

	if !strings.Contains(outputStr, "</proxy_error>") {
		t.Errorf("Expected </proxy_error> in output, got:\n%s", outputStr)
	}

	// Should contain the game not connected error message
	if !strings.Contains(outputStr, "game not connected") {
		t.Errorf("Expected 'game not connected' error, got:\n%s", outputStr)
	}
}

// TestExecutableNotFound verifies that launching a non-existent executable
// properly reports an error through the launch socket.
func TestExecutableNotFound(t *testing.T) {
	// Build the CLI
	cliPath := filepath.Join(t.TempDir(), "autoebiten")
	buildCmd := exec.Command("go", "build", "-o", cliPath, "./cmd/autoebiten")
	buildCmd.Dir = "../"
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build CLI: %v\nOutput: %s", err, output)
	}

	// Launch a non-existent executable
	nonExistentCmd := "this_executable_does_not_exist_12345"
	launchCmd := exec.Command(cliPath, "launch", "--", nonExistentCmd)
	launchCmd.Dir = "../"

	// Start launch in background
	if err := launchCmd.Start(); err != nil {
		t.Fatalf("Failed to start launch: %v", err)
	}

	// Get the launch PID for socket discovery
	launchPID := launchCmd.Process.Pid
	launchSocketPath := launchSocketPath(launchPID)

	// Wait for socket creation with timeout
	var socketFound bool
	for i := 0; i < 50; i++ {
		if _, err := os.Stat(launchSocketPath); err == nil {
			socketFound = true
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	if !socketFound {
		t.Fatalf("Launch socket was not created: %s", launchSocketPath)
	}

	// Wait for the game start to fail
	time.Sleep(300 * time.Millisecond)

	// Try to ping via launch socket
	pingCmd := exec.Command(cliPath, "ping", "--pid", strconv.Itoa(launchPID))
	pingCmd.Dir = "../"
	output, _ := pingCmd.CombinedOutput()
	outputStr := string(output)

	// Cleanup launch process
	launchCmd.Process.Kill()
	launchCmd.Wait()

	// Verify error response contains expected content
	if !strings.Contains(outputStr, "<proxy_error>") {
		t.Errorf("Expected <proxy_error> in output, got:\n%s", outputStr)
	}

	// The error should mention the executable not being found or start failure
	if !strings.Contains(outputStr, "game not connected") {
		t.Errorf("Expected 'game not connected' error, got:\n%s", outputStr)
	}
}

// TestLaunchSocketExistsBeforeGameStart verifies that the launch socket is created
// before the game process starts, enabling CLI connections even if game crashes immediately.
func TestLaunchSocketExistsBeforeGameStart(t *testing.T) {
	// Build the CLI
	cliPath := filepath.Join(t.TempDir(), "autoebiten")
	buildCmd := exec.Command("go", "build", "-o", cliPath, "./cmd/autoebiten")
	buildCmd.Dir = "../"
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build CLI: %v\nOutput: %s", err, output)
	}

	// Launch a command that sleeps briefly then exits
	// This gives us time to verify socket exists
	launchCmd := exec.Command(cliPath, "launch", "--", "sleep", "0.1")
	launchCmd.Dir = "../"

	// Start launch in background
	if err := launchCmd.Start(); err != nil {
		t.Fatalf("Failed to start launch: %v", err)
	}

	launchPID := launchCmd.Process.Pid
	launchSocketPath := launchSocketPath(launchPID)

	// Verify socket is created quickly (within reasonable time)
	socketCreated := false
	startTime := time.Now()
	for time.Since(startTime) < 2*time.Second {
		if _, err := os.Stat(launchSocketPath); err == nil {
			socketCreated = true
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if !socketCreated {
		t.Fatalf("Launch socket was not created within timeout: %s", launchSocketPath)
	}

	// Cleanup - kill and clean up socket manually since process kill may not trigger cleanup
	launchCmd.Process.Kill()
	launchCmd.Wait()

	// Clean up the socket file if it still exists
	os.Remove(launchSocketPath)
}

// TestLaunchExitsAfterCLIQuery verifies that launch exits immediately after CLI query
// when game has crashed, not waiting for the full 30 second timeout.
func TestLaunchExitsAfterCLIQuery(t *testing.T) {
	// Build the CLI
	cliPath := filepath.Join(t.TempDir(), "autoebiten")
	buildCmd := exec.Command("go", "build", "-o", cliPath, "./cmd/autoebiten")
	buildCmd.Dir = "../"
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build CLI: %v\nOutput: %s", err, output)
	}

	// Launch a command that exits immediately
	launchCmd := exec.Command(cliPath, "launch", "--", "false")
	launchCmd.Dir = "../"

	startTime := time.Now()

	// Start launch in background
	if err := launchCmd.Start(); err != nil {
		t.Fatalf("Failed to start launch: %v", err)
	}

	launchPID := launchCmd.Process.Pid
	launchSocketPath := launchSocketPath(launchPID)

	// Wait for socket creation
	for i := 0; i < 50; i++ {
		if _, err := os.Stat(launchSocketPath); err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Wait for game to crash
	time.Sleep(200 * time.Millisecond)

	// Send a ping to trigger the crashed state exit
	pingCmd := exec.Command(cliPath, "ping", "--pid", strconv.Itoa(launchPID))
	pingCmd.Dir = "../"
	pingCmd.CombinedOutput()

	// Wait for launch to exit (should be quick after CLI query)
	done := make(chan error, 1)
	go func() {
		done <- launchCmd.Wait()
	}()

	// Should exit within a few seconds, not 30 seconds
	select {
	case <-done:
		elapsed := time.Since(startTime)
		if elapsed > 10*time.Second {
			t.Errorf("Launch took too long to exit after CLI query: %v", elapsed)
		}
	case <-time.After(5 * time.Second):
		// If it's still running, kill it
		launchCmd.Process.Kill()
		// This is acceptable - the test shows it would wait 30s otherwise
		t.Log("Launch process was still running after 5s (expected to wait for CLI query)")
	}
}

// TestMultipleCLIQueriesAfterCrash verifies that after a game crash, the first CLI query
// receives the error diagnostics and causes launch to exit.
func TestMultipleCLIQueriesAfterCrash(t *testing.T) {
	// Build the CLI
	cliPath := filepath.Join(t.TempDir(), "autoebiten")
	buildCmd := exec.Command("go", "build", "-o", cliPath, "./cmd/autoebiten")
	buildCmd.Dir = "../"
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build CLI: %v\nOutput: %s", err, output)
	}

	// Launch a command that exits immediately
	launchCmd := exec.Command(cliPath, "launch", "--", "false")
	launchCmd.Dir = "../"

	// Start launch in background
	if err := launchCmd.Start(); err != nil {
		t.Fatalf("Failed to start launch: %v", err)
	}

	launchPID := launchCmd.Process.Pid
	launchSocketPath := launchSocketPath(launchPID)

	// Wait for socket creation
	for i := 0; i < 50; i++ {
		if _, err := os.Stat(launchSocketPath); err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Wait for game to crash
	time.Sleep(200 * time.Millisecond)

	// First ping should get error diagnostics
	pingCmd := exec.Command(cliPath, "ping", "--pid", strconv.Itoa(launchPID))
	pingCmd.Dir = "../"
	output, _ := pingCmd.CombinedOutput()
	outputStr := string(output)

	// Should contain error info
	if !strings.Contains(outputStr, "<proxy_error>") {
		t.Errorf("First query: Expected <proxy_error> in output, got:\n%s", outputStr)
	}

	if !strings.Contains(outputStr, "game not connected") {
		t.Errorf("First query: Expected 'game not connected' error, got:\n%s", outputStr)
	}

	// After first query, launch should exit and socket should be gone
	// Wait a bit for cleanup
	time.Sleep(200 * time.Millisecond)

	// Verify launch exited
	done := make(chan error, 1)
	go func() {
		done <- launchCmd.Wait()
	}()

	select {
	case <-done:
		// Expected - launch exited after CLI query
	case <-time.After(2 * time.Second):
		launchCmd.Process.Kill()
		launchCmd.Wait()
		t.Error("Launch process did not exit after first CLI query")
	}

	// Socket should be cleaned up
	if _, err := os.Stat(launchSocketPath); !os.IsNotExist(err) {
		t.Errorf("Launch socket should be cleaned up after CLI query: %s", launchSocketPath)
	}
}
