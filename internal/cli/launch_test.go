package cli

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/s3cy/autoebiten/internal/output"
	"github.com/s3cy/autoebiten/internal/rpc"
)

// TestLaunchSocketPath verifies launch socket path generation.
func TestLaunchSocketPath(t *testing.T) {
	lc := &LaunchCommand{
		options: &LaunchOptions{
			GameCmd: "test",
			Timeout: 5 * time.Second,
		},
	}

	path := lc.launchSocketPath()
	expected := rpc.LaunchSocketPath(os.Getpid())

	if path != expected {
		t.Errorf("launchSocketPath() = %q, want %q", path, expected)
	}
}

// TestGameSocketPath verifies game socket path generation.
func TestGameSocketPath(t *testing.T) {
	lc := &LaunchCommand{
		options: &LaunchOptions{
			GameCmd: "test",
			Timeout: 5 * time.Second,
		},
	}

	path := lc.gameSocketPath()
	expected := rpc.SocketPath()

	if path != expected {
		t.Errorf("gameSocketPath() = %q, want %q", path, expected)
	}
}

// TestCreateLaunchSocket verifies socket creation.
func TestCreateLaunchSocket(t *testing.T) {
	tmpDir := t.TempDir()
	socketPath := filepath.Join(tmpDir, "test-launch.sock")

	lc := &LaunchCommand{
		options: &LaunchOptions{
			GameCmd: "test",
			Timeout: 5 * time.Second,
		},
	}

	listener, err := lc.createLaunchSocket(socketPath)
	if err != nil {
		t.Fatalf("createLaunchSocket failed: %v", err)
	}
	defer listener.Close()

	// Verify socket file exists
	if _, err := os.Stat(socketPath); os.IsNotExist(err) {
		t.Error("Socket file was not created")
	}
}

// TestLaunchCommandStructure verifies the LaunchCommand struct has required fields.
func TestLaunchCommandStructure(t *testing.T) {
	lc := NewLaunchCommand(&LaunchOptions{
		GameCmd: "test",
		Timeout: 5 * time.Second,
	})

	// Verify all required channels are initialized
	if lc.gameExited == nil {
		t.Error("gameExited channel not initialized")
	}
	if lc.crashed == nil {
		t.Error("crashed channel not initialized")
	}
	if lc.done == nil {
		t.Error("done channel not initialized")
	}
}

// TestLaunchSocketBeforeGame verifies launch socket is created before game starts.
func TestLaunchSocketBeforeGame(t *testing.T) {
	// Use a shorter path to avoid macOS socket path length limits
	socketPath := "/tmp/autoebiten-test-launch.sock"

	// Create socket
	lc := &LaunchCommand{
		options: &LaunchOptions{
			GameCmd: "test",
			Timeout: 5 * time.Second,
		},
	}

	listener, err := lc.createLaunchSocket(socketPath)
	if err != nil {
		t.Fatalf("createLaunchSocket failed: %v", err)
	}
	defer listener.Close()

	// Verify socket exists and is listening
	if _, err := os.Stat(socketPath); os.IsNotExist(err) {
		t.Error("Launch socket should exist before game starts")
	}
}

// TestOnCrashedCallback verifies the onCrashed callback closes the crashed channel.
func TestOnCrashedCallback(t *testing.T) {
	lc := NewLaunchCommand(&LaunchOptions{
		GameCmd: "test",
		Timeout: 5 * time.Second,
	})

	// Set up test to verify channel is closed
	done := make(chan bool, 1)

	go func() {
		select {
		case <-lc.crashed:
			done <- true
		case <-time.After(100 * time.Millisecond):
			done <- false
		}
	}()

	// Call onCrashed callback
	lc.onCrashedCallback()

	result := <-done
	if !result {
		t.Error("onCrashedCallback should close the crashed channel")
	}
}

// TestCreateGameCommandEnv verifies AUTOEBITEN_SOCKET is set.
func TestCreateGameCommandEnv(t *testing.T) {
	lc := NewLaunchCommand(&LaunchOptions{
		GameCmd:  "echo",
		GameArgs: []string{"hello"},
		Timeout:  5 * time.Second,
	})

	cmd, _, _, err := lc.createGameCommand()
	if err != nil {
		t.Fatalf("createGameCommand failed: %v", err)
	}

	// Check that AUTOEBITEN_SOCKET is in environment
	found := false
	expectedPath := lc.gameSocketPath()
	for _, env := range cmd.Env {
		if env == "AUTOEBITEN_SOCKET="+expectedPath {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("AUTOEBITEN_SOCKET not set to %q in environment", expectedPath)
	}
}

// TestUnifiedHandlerIntegration verifies handler is created with onCrashed callback.
func TestUnifiedHandlerIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatalf("Failed to create log file: %v", err)
	}
	defer logFile.Close()

	outputMgr := output.NewOutputManager(logFile, logPath, snapPath)

	lc := NewLaunchCommand(&LaunchOptions{
		GameCmd: "test",
		Timeout: 5 * time.Second,
	})
	lc.outputMgr = outputMgr

	// Create handler with callback
	handler := lc.createHandler()

	if handler == nil {
		t.Fatal("createHandler returned nil")
	}

	// Verify handler is in Waiting state
	if handler.GetState() != 0 { // StateWaiting = 0
		t.Errorf("Handler state = %v, want StateWaiting", handler.GetState())
	}
}

// TestWaitForExitChannels verifies waitForExit responds to channels.
func TestWaitForExitChannels(t *testing.T) {
	t.Run("exits on done signal", func(t *testing.T) {
		lc := NewLaunchCommand(&LaunchOptions{
			GameCmd: "test",
			Timeout: 100 * time.Millisecond,
		})

		// Signal done
		go func() {
			time.Sleep(10 * time.Millisecond)
			close(lc.done)
		}()

		// Should exit quickly
		start := time.Now()
		lc.waitForExit("")
		elapsed := time.Since(start)

		if elapsed > 50*time.Millisecond {
			t.Errorf("waitForExit took too long: %v", elapsed)
		}
	})

	t.Run("exits on crashed signal", func(t *testing.T) {
		lc := NewLaunchCommand(&LaunchOptions{
			GameCmd: "test",
			Timeout: 100 * time.Millisecond,
		})

		// Signal crashed
		go func() {
			time.Sleep(10 * time.Millisecond)
			close(lc.crashed)
		}()

		// Should exit quickly
		start := time.Now()
		lc.waitForExit("")
		elapsed := time.Since(start)

		if elapsed > 50*time.Millisecond {
			t.Errorf("waitForExit took too long on crashed: %v", elapsed)
		}
	})
}
