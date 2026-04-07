package output

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FilePath holds the paths for output capture files.
type FilePath struct {
	Log        string // autoebiten-{PID}-output.log
	Snapshot   string // autoebiten-{PID}-snapshot.log
	LaunchSock string // autoebiten-{PID}-launch.sock
}

// DerivePaths derives output file paths from a socket path.
// Socket format: {dir}/autoebiten-{PID}.sock
// Output files:  {dir}/autoebiten-{PID}-output.log
//
//	{dir}/autoebiten-{PID}-snapshot.log
//	{dir}/autoebiten-{PID}-launch.sock
func DerivePaths(socketPath string) *FilePath {
	dir := filepath.Dir(socketPath)
	base := filepath.Base(socketPath)

	// Remove .sock extension to get "autoebiten-{PID}"
	baseWithoutExt := strings.TrimSuffix(base, ".sock")

	return &FilePath{
		Log:        filepath.Join(dir, baseWithoutExt+"-output.log"),
		Snapshot:   filepath.Join(dir, baseWithoutExt+"-snapshot.log"),
		LaunchSock: filepath.Join(dir, baseWithoutExt+"-launch.sock"),
	}
}

// CreateLogFile creates the log file with secure permissions.
// Returns the file handle for writing.
func CreateLogFile(path string) (*os.File, error) {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create file with restricted permissions (0600)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_TRUNC, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	return file, nil
}

// writeSnapshot writes content to the snapshot file.
func writeSnapshot(path string, content []byte) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create snapshot directory: %w", err)
	}

	// Write with restricted permissions (0600)
	if err := os.WriteFile(path, content, 0600); err != nil {
		return fmt.Errorf("failed to write snapshot: %w", err)
	}

	return nil
}

// CopyFile copies src to dst, creating dst if needed.
// Returns error if src doesn't exist (not os.IsNotExist).
func CopyFile(dst, src string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		if os.IsNotExist(err) {
			// Source doesn't exist - write empty file
			return writeSnapshot(dst, nil)
		}
		return fmt.Errorf("failed to read source: %w", err)
	}
	return writeSnapshot(dst, content)
}

// GenerateDiff generates a unified diff between snapshot and current file paths.
// Uses diff -u directly on files with sed to handle carriage returns.
func GenerateDiff(snapshotPath, currentPath string) string {
	// Use bash process substitution to:
	// 1. Process carriage returns (keep only content after last \r per line)
	// 2. Run diff -u on the processed content
	cmd := exec.Command("bash", "-c",
		"diff -u <(sed 's/^.*\\r//' '"+snapshotPath+"' 2>/dev/null || echo '') "+
			"<(sed 's/^.*\\r//' '"+currentPath+"' 2>/dev/null || echo '')")

	output, err := cmd.Output()
	if err != nil {
		// diff returns exit code 1 when files differ (expected)
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return string(output)
		}
		// If files don't exist or are empty, handle gracefully
		return ""
	}

	return string(output)
}
