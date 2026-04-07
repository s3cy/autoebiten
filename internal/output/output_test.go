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

func TestGenerateDiff(t *testing.T) {
	tests := []struct {
		name     string
		snapshot string
		current  string
		wantDiff bool // whether we expect some diff content
	}{
		{
			name:     "no changes",
			snapshot: "hello\nworld",
			current:  "hello\nworld",
			wantDiff: false,
		},
		{
			name:     "added lines",
			snapshot: "hello",
			current:  "hello\nworld\ntest",
			wantDiff: true,
		},
		{
			name:     "removed lines",
			snapshot: "hello\nworld\ntest",
			current:  "hello",
			wantDiff: true,
		},
		{
			name:     "modified lines",
			snapshot: "hello\nworld",
			current:  "hello\nchanged",
			wantDiff: true,
		},
		{
			name:     "empty to content",
			snapshot: "",
			current:  "new line",
			wantDiff: true,
		},
		{
			name:     "content to empty",
			snapshot: "old line",
			current:  "",
			wantDiff: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp files
			tmpDir := t.TempDir()
			snapshotPath := filepath.Join(tmpDir, "snapshot")
			currentPath := filepath.Join(tmpDir, "current")

			os.WriteFile(snapshotPath, []byte(tt.snapshot), 0600)
			os.WriteFile(currentPath, []byte(tt.current), 0600)

			diff := GenerateDiff(snapshotPath, currentPath)

			if tt.wantDiff {
				if diff == "" {
					t.Errorf("Expected non-empty diff for %q vs %q", tt.snapshot, tt.current)
				}
				// Check for diff format markers
				if !strings.Contains(diff, "---") {
					t.Errorf("Diff missing '---' header")
				}
				if !strings.Contains(diff, "+++") {
					t.Errorf("Diff missing '+++' header")
				}
				if !strings.Contains(diff, "@@") {
					t.Errorf("Diff missing hunk header")
				}
			} else {
				if diff != "" {
					t.Errorf("Expected empty diff, got:\n%s", diff)
				}
			}
		})
	}
}

func TestGenerateDiffWithCarriageReturn(t *testing.T) {
	// Test that diff correctly handles carriage returns
	snapshot := "Loading: 50%"
	current := "Loading: 0%\rLoading: 50%\rLoading: 100%"

	// Create temp files
	tmpDir := t.TempDir()
	snapshotPath := filepath.Join(tmpDir, "snapshot")
	currentPath := filepath.Join(tmpDir, "current")

	os.WriteFile(snapshotPath, []byte(snapshot), 0600)
	os.WriteFile(currentPath, []byte(current), 0600)

	diff := GenerateDiff(snapshotPath, currentPath)

	// Should show the final state after processing \r
	if diff == "" {
		t.Errorf("Expected non-empty diff")
	}
	// The diff should show Loading: 100% (final state)
	if !strings.Contains(diff, "Loading: 100%") {
		t.Errorf("Diff should show final state 'Loading: 100%%', got:\n%s", diff)
	}
}
