package output

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCarriageReturnWriterBasicOverwrite(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single overwrite",
			input:    "Loading 10%\rLoading 20%\n",
			expected: "Loading 20%\n",
		},
		{
			name:     "multiple overwrites",
			input:    "0%\r25%\r50%\r75%\r100%\n",
			expected: "100%\n",
		},
		{
			name:     "no overwrite just newline",
			input:    "Hello World\n",
			expected: "Hello World\n",
		},
		{
			name:     "incomplete line flushed at end",
			input:    "Loading 50%",
			expected: "Loading 50%\n",
		},
		{
			name:     "overwrite then incomplete line",
			input:    "Loading 0%\rLoading 100%",
			expected: "Loading 100%\n",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "multiple lines",
			input:    "Line1\nLine2\rLine2Modified\nLine3\n",
			expected: "Line1\nLine2Modified\nLine3\n",
		},
		{
			name:     "multiple carriage returns before newline",
			input:    "abc\rdef\rghi\n",
			expected: "ghi\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &bytes.Buffer{}
			writer := NewCarriageReturnWriter(dst)
			writer.Write([]byte(tt.input))
			writer.Flush()
			assert.Equal(t, tt.expected, dst.String())
		})
	}
}

func TestOutputManagerWrite(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	logFile, err := os.Create(logPath)
	require.NoError(t, err)

	manager := NewOutputManager(logFile, logPath, snapPath)

	// Write some data
	n, err := manager.Write([]byte("Hello\n"))
	assert.NoError(t, err)
	assert.Equal(t, 6, n)

	// Read file to verify
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)
	assert.Equal(t, "Hello\n", string(content))
}

func TestOutputManagerDiffAndUpdateSnapshot(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	// Create initial snapshot
	os.WriteFile(snapPath, []byte("old line\n"), 0600)

	// Create log with new content
	logFile, err := os.Create(logPath)
	require.NoError(t, err)
	logFile.WriteString("new line\n")
	logFile.Close()

	// Reopen log file for manager
	logFile, err = os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND, 0600)
	require.NoError(t, err)

	manager := NewOutputManager(logFile, logPath, snapPath)

	diff, err := manager.DiffAndUpdateSnapshot()
	require.NoError(t, err)
	assert.Contains(t, diff, "new line")

	// Verify snapshot was updated
	snapContent, err := os.ReadFile(snapPath)
	require.NoError(t, err)
	assert.Equal(t, "new line\n", string(snapContent))
}

func TestOutputManagerDiffAndUpdateSnapshotNoChanges(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	// Same content in both
	os.WriteFile(snapPath, []byte("same content\n"), 0600)
	logFile, err := os.Create(logPath)
	require.NoError(t, err)
	logFile.WriteString("same content\n")
	logFile.Close()

	logFile, err = os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND, 0600)
	require.NoError(t, err)

	manager := NewOutputManager(logFile, logPath, snapPath)

	diff, err := manager.DiffAndUpdateSnapshot()
	require.NoError(t, err)
	assert.Empty(t, diff) // No diff when identical
}

func TestOutputManagerDiffAndUpdateSnapshotEmptyFiles(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")
	snapPath := filepath.Join(tmpDir, "test-snapshot.log")

	// Empty log file
	logFile, err := os.Create(logPath)
	require.NoError(t, err)

	manager := NewOutputManager(logFile, logPath, snapPath)

	diff, err := manager.DiffAndUpdateSnapshot()
	require.NoError(t, err)
	assert.Empty(t, diff) // Empty diff when both empty
}