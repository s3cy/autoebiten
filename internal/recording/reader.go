package recording

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// Reader reads recording files.
type Reader struct {
	pid int
}

// NewReader creates a reader for the given PID.
func NewReader(pid int) *Reader {
	return &Reader{pid: pid}
}

// ReadAll reads all entries from the recording file.
// Returns empty slice if file doesn't exist.
func (r *Reader) ReadAll() ([]Entry, error) {
	path := Path(r.pid)

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

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var entry Entry
		if err := json.Unmarshal(line, &entry); err != nil {
			// Skip corrupted entries
			continue
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read recording file: %w", err)
	}

	return entries, nil
}