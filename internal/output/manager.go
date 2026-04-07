package output

import (
	"io"
	"os"
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