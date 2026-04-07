package output

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/aymanbagabas/go-udiff"
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

	// Read files content
	snapshotContent, snapshotModTime, err := readFileAndModTime(m.snapshotPath)
	if err != nil {
		return "", err
	}
	outputContent, outputModTime, err := readFileAndModTime(m.outputPath)
	if err != nil {
		return "", err
	}

	diff := generateUnifiedDiff("snapshot "+snapshotModTime, "current "+outputModTime, snapshotContent, outputContent)

	// Copy output to snapshot
	if err := writeSnapshot(m.snapshotPath, outputContent); err != nil {
		return "", err
	}

	return diff, nil
}

func readFileAndModTime(path string) ([]byte, string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		// Treat non-existent file as empty content
		if os.IsNotExist(err) {
			return nil, "(empty)", nil
		}
		return nil, "", err
	}
	info, err := os.Stat(path)
	if err != nil {
		return nil, "", err
	}
	return content, info.ModTime().Format("2006-01-02 15:04:05.000"), nil
}

// generateUnifiedDiff creates a unified diff between old and new content.
func generateUnifiedDiff(oldLabel, newLabel string, oldContent, newContent []byte) string {
	edits := udiff.Bytes(oldContent, newContent)
	unified, err := udiff.ToUnified(oldLabel, newLabel, string(oldContent), edits, udiff.DefaultContextLines)
	if err != nil {
		// Can't happen: edits are consistent.
		log.Fatalf("internal error in udiff: %v", err)
	}
	return unified
}
