package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/s3cy/autoebiten/internal/output"
	"github.com/s3cy/autoebiten/internal/proxy"
	"github.com/s3cy/autoebiten/internal/rpc"
)

// TestOutputCapture tests the output capture functionality.
func TestOutputCapture(t *testing.T) {
	t.Run("derive paths from socket", func(t *testing.T) {
		socketPath := "/tmp/autoebiten/autoebiten-12345.sock"
		paths := output.DerivePaths(socketPath)

		if paths.Log != "/tmp/autoebiten/autoebiten-12345-output.log" {
			t.Errorf("Log path = %q, want %q", paths.Log, "/tmp/autoebiten/autoebiten-12345-output.log")
		}
		if paths.Snapshot != "/tmp/autoebiten/autoebiten-12345-snapshot.log" {
			t.Errorf("Snapshot path = %q, want %q", paths.Snapshot, "/tmp/autoebiten/autoebiten-12345-snapshot.log")
		}
		if paths.LaunchSock != "/tmp/autoebiten/autoebiten-12345-launch.sock" {
			t.Errorf("LaunchSock path = %q, want %q", paths.LaunchSock, "/tmp/autoebiten/autoebiten-12345-launch.sock")
		}
	})

	t.Run("generate diff", func(t *testing.T) {
		tmpDir := t.TempDir()
		snapshotPath := filepath.Join(tmpDir, "snapshot")
		currentPath := filepath.Join(tmpDir, "current")

		os.WriteFile(snapshotPath, []byte("line1\nline2"), 0600)
		os.WriteFile(currentPath, []byte("line1\nline2\nline3"), 0600)

		diff := output.GenerateDiff(snapshotPath, currentPath)
		if diff == "" {
			t.Error("Expected non-empty diff")
		}
		if !strings.Contains(diff, "@@") {
			t.Error("Expected hunk header in diff")
		}
		if !strings.Contains(diff, "+line3") {
			t.Error("Expected added line in diff")
		}
	})
}

// TestLaunchSocketPath tests the launch socket path derivation.
func TestLaunchSocketPath(t *testing.T) {
	t.Run("derive from game socket", func(t *testing.T) {
		gameSocket := "/tmp/autoebiten/autoebiten-12345.sock"
		paths := output.DerivePaths(gameSocket)
		launchSocket := paths.LaunchSock

		expected := "/tmp/autoebiten/autoebiten-12345-launch.sock"
		if launchSocket != expected {
			t.Errorf("LaunchSock = %q, want %q", launchSocket, expected)
		}
	})

	t.Run("derive from custom socket", func(t *testing.T) {
		gameSocket := "/var/run/mygame.sock"
		paths := output.DerivePaths(gameSocket)
		launchSocket := paths.LaunchSock

		expected := "/var/run/mygame-launch.sock"
		if launchSocket != expected {
			t.Errorf("LaunchSock = %q, want %q", launchSocket, expected)
		}
	})
}

// TestProxyResponse tests proxy response structure.
func TestProxyResponse(t *testing.T) {
	resp := proxy.Response{
		RPCResponse: rpc.RPCResponse{
			JSONRPC: "2.0",
			ID:      1,
			Result:  []byte(`{"success": true}`),
		},
		Output: "--- snapshot\n+++ current\n@@ -1 +1,2 @@\n line1\n+line2",
	}

	if resp.JSONRPC != "2.0" {
		t.Errorf("JSONRPC = %q, want %q", resp.JSONRPC, "2.0")
	}
	if resp.Output == "" {
		t.Error("Expected non-empty output")
	}
}

// TestLaunchCommandExists tests that the launch command exists.
func TestLaunchCommandExists(t *testing.T) {
	// Build the CLI
	cliPath := filepath.Join(t.TempDir(), "autoebiten")
	buildCmd := exec.Command("go", "build", "-o", cliPath, "./cmd/autoebiten")
	buildCmd.Dir = "../"
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build CLI: %v\nOutput: %s", err, output)
	}

	// Test that launch command shows usage when called without args
	cmd := exec.Command(cliPath, "launch")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected error when running launch without args")
	}
	if !strings.Contains(string(output), "no game command provided") {
		t.Errorf("Expected 'no game command provided' in output, got: %s", output)
	}
}

// TestLogFileCreation tests that log files are created with proper permissions.
func TestLogFileCreation(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	file, err := output.CreateLogFile(logPath)
	if err != nil {
		t.Fatalf("CreateLogFile failed: %v", err)
	}
	file.Close()

	// Check file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}

	// Check permissions (0600)
	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("File permissions = %o, want 0600", info.Mode().Perm())
	}
}

// TestProxyServerCleanup tests proxy server cleanup functionality.
func TestProxyServerCleanup(t *testing.T) {
	tmpDir := t.TempDir()
	paths := &output.FilePath{
		Log:        filepath.Join(tmpDir, "test.log"),
		Snapshot:   filepath.Join(tmpDir, "test-snapshot.log"),
		LaunchSock: filepath.Join(tmpDir, "test-launch.sock"),
	}

	// Create files
	os.WriteFile(paths.Log, []byte("log content"), 0600)
	os.WriteFile(paths.Snapshot, []byte("snapshot content"), 0600)

	server := proxy.NewServer(nil, paths)
	err := server.Cleanup()
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Verify files are removed
	if _, err := os.Stat(paths.Log); !os.IsNotExist(err) {
		t.Error("Log file should be removed")
	}
	if _, err := os.Stat(paths.Snapshot); !os.IsNotExist(err) {
		t.Error("Snapshot file should be removed")
	}
}

// TestDiffOutputFormat tests that diff output has proper format.
func TestDiffOutputFormat(t *testing.T) {
	tmpDir := t.TempDir()
	snapshotPath := filepath.Join(tmpDir, "snapshot")
	currentPath := filepath.Join(tmpDir, "current")

	os.WriteFile(snapshotPath, []byte("[INFO] Game started"), 0600)
	os.WriteFile(currentPath, []byte("[INFO] Game started\nLoading: 100%"), 0600)

	diff := output.GenerateDiff(snapshotPath, currentPath)

	// Check for unified diff format markers
	if !strings.Contains(diff, "---") {
		t.Error("Expected '---' header")
	}
	if !strings.Contains(diff, "+++") {
		t.Error("Expected '+++' header")
	}
	if !strings.Contains(diff, "@@") {
		t.Error("Expected hunk header '@@'")
	}
	if !strings.Contains(diff, "+Loading: 100%") {
		t.Error("Expected added line with + prefix")
	}
}

// TestEmptyDiff tests that no changes produces empty diff.
func TestEmptyDiff(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "same")

	os.WriteFile(filePath, []byte("same content\nmultiple lines"), 0600)

	diff := output.GenerateDiff(filePath, filePath)

	if diff != "" {
		t.Errorf("Expected empty diff for identical files, got: %s", diff)
	}
}

// TestCarriageReturnInDiff tests diff with carriage return sequences.
func TestCarriageReturnInDiff(t *testing.T) {
	tmpDir := t.TempDir()
	snapshotPath := filepath.Join(tmpDir, "snapshot")
	currentPath := filepath.Join(tmpDir, "current")

	// Progress bar style output
	os.WriteFile(snapshotPath, []byte(""), 0600)
	os.WriteFile(currentPath, []byte("Loading: 0%\rLoading: 50%\rLoading: 100%\nDone!"), 0600)

	diff := output.GenerateDiff(snapshotPath, currentPath)

	// Should show final state
	if !strings.Contains(diff, "Loading: 100%") {
		t.Errorf("Expected 'Loading: 100%%' in diff, got:\n%s", diff)
	}
	if strings.Contains(diff, "Loading: 0%") {
		t.Error("Should not show intermediate state 'Loading: 0%'")
	}
	if strings.Contains(diff, "Loading: 50%") {
		t.Error("Should not show intermediate state 'Loading: 50%'")
	}
}

// TestConcurrentProxyRequests tests that proxy handles concurrent requests safely.
func TestConcurrentProxyRequests(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	// Create initial log
	os.WriteFile(logPath, []byte("initial\n"), 0600)

	paths := &output.FilePath{
		Log:        logPath,
		Snapshot:   snapPath,
		LaunchSock: filepath.Join(tmpDir, "test-launch.sock"),
	}

	// Create a mock game client that simulates delays
	mock := &delayedMockClient{delay: 10 * time.Millisecond}
	server := proxy.NewServer(mock, paths)

	// Run concurrent requests
	done := make(chan bool, 3)
	for i := range make([]struct{}, 3) {
		go func() {
			// Append to log to simulate game output
			f, _ := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0600)
			f.WriteString("output\n")
			f.Close()

			req := &rpc.RPCRequest{
				JSONRPC: "2.0",
				ID:      i,
				Method:  "ping",
			}
			_, err := server.ForwardRequest(req)
			if err != nil {
				t.Errorf("Concurrent request failed: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all
	for range make([]struct{}, 3) {
		<-done
	}

	// All requests should complete without race conditions
	if mock.callCount != 3 {
		t.Errorf("Expected 3 calls, got %d", mock.callCount)
	}
}

// delayedMockClient simulates a game client with delay.
type delayedMockClient struct {
	callCount int
	delay     time.Duration
}

func (m *delayedMockClient) SendRequest(req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	m.callCount++
	time.Sleep(m.delay)
	return &rpc.RPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  []byte(`{"ok": true, "version": "test"}`),
	}, nil
}

func (m *delayedMockClient) Close() error {
	return nil
}
