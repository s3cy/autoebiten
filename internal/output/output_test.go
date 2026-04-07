package output

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDerivePaths(t *testing.T) {
	tests := []struct {
		name       string
		socketPath string
		wantLog    string
		wantSnap   string
		wantLaunch string
	}{
		{
			name:       "default socket path",
			socketPath: "/tmp/autoebiten/autoebiten-12345.sock",
			wantLog:    "/tmp/autoebiten/autoebiten-12345-output.log",
			wantSnap:   "/tmp/autoebiten/autoebiten-12345-snapshot.log",
			wantLaunch: "/tmp/autoebiten/autoebiten-12345-launch.sock",
		},
		{
			name:       "custom socket path",
			socketPath: "/var/run/mygame.sock",
			wantLog:    "/var/run/mygame-output.log",
			wantSnap:   "/var/run/mygame-snapshot.log",
			wantLaunch: "/var/run/mygame-launch.sock",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paths := DerivePaths(tt.socketPath)
			if paths.Log != tt.wantLog {
				t.Errorf("Log = %q, want %q", paths.Log, tt.wantLog)
			}
			if paths.Snapshot != tt.wantSnap {
				t.Errorf("Snapshot = %q, want %q", paths.Snapshot, tt.wantSnap)
			}
			if paths.LaunchSock != tt.wantLaunch {
				t.Errorf("LaunchSock = %q, want %q", paths.LaunchSock, tt.wantLaunch)
			}
		})
	}
}

func TestCreateLogFile(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	// Test creating log file
	file, err := CreateLogFile(logPath)
	if err != nil {
		t.Fatalf("CreateLogFile failed: %v", err)
	}
	file.Close()

	// Check file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Errorf("Log file was not created")
	}

	// Check permissions
	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("File permissions = %o, want 0600", info.Mode().Perm())
	}

	// Test writing and appending
	file, err = CreateLogFile(logPath)
	if err != nil {
		t.Fatalf("CreateLogFile (append) failed: %v", err)
	}
	file.WriteString("line1\n")
	file.Close()

	file, err = CreateLogFile(logPath)
	if err != nil {
		t.Fatalf("CreateLogFile (reopen) failed: %v", err)
	}
	file.WriteString("line2\n")
	file.Close()

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	// File was truncated on reopen, so only line2 should exist
	if !strings.Contains(string(content), "line2") {
		t.Errorf("Expected file to contain line2 after truncate")
	}
}

func TestWriteSnapshot(t *testing.T) {
	tmpDir := t.TempDir()
	snapPath := filepath.Join(tmpDir, "snapshot.log")

	content := []byte("test content\nline2")
	err := writeSnapshot(snapPath, content)
	if err != nil {
		t.Fatalf("WriteSnapshot failed: %v", err)
	}

	// Read back and verify
	readContent, err := os.ReadFile(snapPath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if string(readContent) != string(content) {
		t.Errorf("Written content = %q, want %q", readContent, content)
	}

	// Check permissions
	info, err := os.Stat(snapPath)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("File permissions = %o, want 0600", info.Mode().Perm())
	}
}