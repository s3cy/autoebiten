package output

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
)

// CarriageReturnWriter interprets \r as "move cursor to line start, overwrite".
// It writes to dst only the final visual state of each line.
type CarriageReturnWriter struct {
	dst     io.Writer
	lineBuf []byte
}

// NewCarriageReturnWriter creates a new carriage return interpreting writer.
func NewCarriageReturnWriter(dst io.Writer) *CarriageReturnWriter {
	return &CarriageReturnWriter{
		dst:     dst,
		lineBuf: make([]byte, 0, 1024),
	}
}

// Write processes input bytes, interpreting \r as line overwrite.
func (w *CarriageReturnWriter) Write(data []byte) (int, error) {
	for _, b := range data {
		switch b {
		case '\r':
			// Move cursor to line start: clear line buffer
			w.lineBuf = w.lineBuf[:0]
		case '\n':
			// Flush completed line
			if len(w.lineBuf) > 0 {
				w.dst.Write(w.lineBuf)
			}
			w.dst.Write([]byte{'\n'})
			w.lineBuf = w.lineBuf[:0]
		default:
			w.lineBuf = append(w.lineBuf, b)
		}
	}
	return len(data), nil
}

// Flush writes any remaining incomplete line (with newline appended).
func (w *CarriageReturnWriter) Flush() error {
	if len(w.lineBuf) > 0 {
		if _, err := w.dst.Write(w.lineBuf); err != nil {
			return err
		}
		if _, err := w.dst.Write([]byte{'\n'}); err != nil {
			return err
		}
		w.lineBuf = w.lineBuf[:0]
	}
	return nil
}

// OutputManager coordinates mutex-protected write, diff, and snapshot operations.
type OutputManager struct {
	mu           sync.Mutex
	logFile      *os.File
	outputPath   string
	snapshotPath string
}

// NewOutputManager creates a new output manager.
func NewOutputManager(logFile *os.File, outputPath, snapshotPath string) *OutputManager {
	return &OutputManager{
		logFile:      logFile,
		outputPath:   outputPath,
		snapshotPath: snapshotPath,
	}
}

// Write writes data to the log file with mutex protection and flush.
func (m *OutputManager) Write(data []byte) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	n, err := m.logFile.Write(data)
	if err != nil {
		return n, err
	}
	m.logFile.Sync()
	return n, nil
}

// DiffAndUpdateSnapshot generates diff between snapshot and log, then copies log to snapshot.
// Returns the diff string.
func (m *OutputManager) DiffAndUpdateSnapshot() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Flush log file to ensure diff sees latest content
	m.logFile.Sync()

	// Read both files content
	snapshotContent, snapshotErr := os.ReadFile(m.snapshotPath)
	outputContent, outputErr := os.ReadFile(m.outputPath)

	// Treat non-existent files as empty content
	if os.IsNotExist(snapshotErr) {
		snapshotContent = nil
	} else if snapshotErr != nil {
		return "", snapshotErr
	}
	if os.IsNotExist(outputErr) {
		outputContent = nil
	} else if outputErr != nil {
		return "", outputErr
	}

	// Generate diff only if content differs
	var diff string
	if !bytes.Equal(snapshotContent, outputContent) {
		diff = generateUnifiedDiff(snapshotContent, outputContent)
	}

	// Copy output to snapshot
	if err := writeSnapshot(m.snapshotPath, outputContent); err != nil {
		return diff, err
	}

	return diff, nil
}

// generateUnifiedDiff creates a unified diff between old and new content.
func generateUnifiedDiff(oldContent, newContent []byte) string {
	oldLines := splitLines(oldContent)
	newLines := splitLines(newContent)

	// Build simple diff output
	result := "--- snapshot\n+++ current\n"

	// Simple approach: show all removed lines, then all added lines
	oldSet := make(map[string]bool)
	for _, line := range oldLines {
		oldSet[line] = true
	}

	newSet := make(map[string]bool)
	for _, line := range newLines {
		newSet[line] = true
	}

	// Lines in old but not in new (removed)
	for _, line := range oldLines {
		if !newSet[line] {
			result += "-" + line + "\n"
		}
	}

	// Lines in new but not in old (added)
	for _, line := range newLines {
		if !oldSet[line] {
			result += "+" + line + "\n"
		}
	}

	return result
}

// splitLines splits content into lines without trailing newlines.
func splitLines(content []byte) []string {
	if len(content) == 0 {
		return nil
	}
	s := string(content)
	if s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}