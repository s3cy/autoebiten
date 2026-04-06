package recording

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// Reader reads recording files.
type Reader struct {
	socketPath string
}

// NewReaderFromSocket creates a reader using a socket path.
// The recording file path is derived from the socket path.
func NewReaderFromSocket(socketPath string) *Reader {
	return &Reader{socketPath: socketPath}
}

// ReadAll reads all entries from the recording file.
// Returns empty slice if file doesn't exist.
func (r *Reader) ReadAll() ([]Entry, error) {
	path := PathFromSocket(r.socketPath)

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Entry{}, nil
		}
		return nil, fmt.Errorf("failed to open recording file: %w", err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)

	lineno := 0
	for scanner.Scan() {
		lineno += 1
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var entry Entry
		if err := json.Unmarshal(line, &entry); err != nil {
			return nil, fmt.Errorf("entry corrupted at line %d: %w", lineno, err)
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read recording file: %w", err)
	}

	return entries, nil
}
