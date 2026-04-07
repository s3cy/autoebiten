package e2e

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/s3cy/autoebiten/internal/output"
)

func TestCarriageReturnWriterIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "game.log")
	snapPath := filepath.Join(tmpDir, "snapshot.log")

	logFile, err := os.Create(logPath)
	require.NoError(t, err)

	outputMgr := output.NewOutputManager(logFile, logPath, snapPath)
	writer := output.NewCarriageReturnWriter(outputMgr)

	// Simulate game output with progress bar
	writer.Write([]byte("Loading: 0%"))
	writer.Write([]byte("\rLoading: 25%"))
	writer.Write([]byte("\rLoading: 50%"))
	writer.Write([]byte("\rLoading: 75%"))
	writer.Write([]byte("\rLoading: 100%\n"))
	writer.Write([]byte("Done!\n"))
	writer.Flush()

	// Verify log file shows final visual state
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)

	// Should show final progress state and Done
	assert.Contains(t, string(content), "Loading: 100%\n")
	assert.Contains(t, string(content), "Done!\n")
	assert.NotContains(t, string(content), "Loading: 0%")
	assert.NotContains(t, string(content), "Loading: 25%")

	// Generate diff and update snapshot
	diff, err := outputMgr.DiffAndUpdateSnapshot()
	require.NoError(t, err)
	assert.Contains(t, diff, "Loading: 100%")

	// Verify snapshot was updated
	snapContent, err := os.ReadFile(snapPath)
	require.NoError(t, err)
	assert.Equal(t, string(content), string(snapContent))
}

func TestOutputManagerConcurrentWriteAndDiff(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "game.log")
	snapPath := filepath.Join(tmpDir, "snapshot.log")

	logFile, err := os.Create(logPath)
	require.NoError(t, err)

	outputMgr := output.NewOutputManager(logFile, logPath, snapPath)
	writer := output.NewCarriageReturnWriter(outputMgr)

	// Write in background (simulating game output)
	writeDone := make(chan bool)
	go func() {
		for i := range 5 {
			writer.Write([]byte("Progress: " + strconv.Itoa(i*20) + "%\r"))
			time.Sleep(10 * time.Millisecond)
		}
		writer.Write([]byte("Progress: 100%\nComplete!\n"))
		writer.Flush()
		writeDone <- true
	}()

	// Wait for writing to complete
	<-writeDone

	// Now diff
	diff, err := outputMgr.DiffAndUpdateSnapshot()
	require.NoError(t, err)
	assert.Contains(t, diff, "Complete!")
}

func TestCarriageReturnWriterMultipleLines(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "game.log")
	snapPath := filepath.Join(tmpDir, "snapshot.log")

	logFile, err := os.Create(logPath)
	require.NoError(t, err)

	outputMgr := output.NewOutputManager(logFile, logPath, snapPath)
	writer := output.NewCarriageReturnWriter(outputMgr)

	// Simulate multiple lines with overwrites
	writer.Write([]byte("Line 1: Start\n"))
	writer.Write([]byte("Progress: 0%\rProgress: 50%\rProgress: 100%\n"))
	writer.Write([]byte("Line 3: End\n"))
	writer.Flush()

	// Verify log file shows final visual state
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)

	expected := "Line 1: Start\nProgress: 100%\nLine 3: End\n"
	assert.Equal(t, expected, string(content))
}

func TestCarriageReturnWriterEmptyOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "game.log")
	snapPath := filepath.Join(tmpDir, "snapshot.log")

	logFile, err := os.Create(logPath)
	require.NoError(t, err)

	outputMgr := output.NewOutputManager(logFile, logPath, snapPath)
	writer := output.NewCarriageReturnWriter(outputMgr)

	// Test clearing line with \r before newline
	writer.Write([]byte("This will be cleared\r\n"))
	writer.Write([]byte("This stays\n"))
	writer.Flush()

	content, err := os.ReadFile(logPath)
	require.NoError(t, err)

	// \r followed by \n should result in just \n (empty line cleared)
	expected := "\nThis stays\n"
	assert.Equal(t, expected, string(content))
}