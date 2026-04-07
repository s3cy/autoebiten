package output

import (
	"fmt"
	"os"
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
