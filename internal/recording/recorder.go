package recording

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/sys/unix"

	"github.com/s3cy/autoebiten/internal/script"
)

// PathFromSocket derives the recording file path from a socket path.
// The recording file is placed in the same directory with a name based on the socket.
// Example: socket "autoebiten-12345.sock" → recording "autoebiten-12345-recording.jsonl"
func PathFromSocket(socketPath string) string {
	dir := filepath.Dir(socketPath)
	base := filepath.Base(socketPath)

	// Remove .sock extension if present
	name := strings.TrimSuffix(base, ".sock")

	return filepath.Join(dir, fmt.Sprintf("%s-recording.jsonl", name))
}

// entryCommandWrapper is used for JSON marshaling/unmarshaling with discriminator.
type entryCommandWrapper struct {
	Input      *script.InputCmd      `json:"input,omitempty"`
	Mouse      *script.MouseCmd      `json:"mouse,omitempty"`
	Wheel      *script.WheelCmd      `json:"wheel,omitempty"`
	Screenshot *script.ScreenshotCmd `json:"screenshot,omitempty"`
	Delay      *script.DelayCmd      `json:"delay,omitempty"`
	Custom     *script.CustomCmd     `json:"custom,omitempty"`
	Repeat     *script.RepeatCmd     `json:"repeat,omitempty"`
	State      *script.StateCmd      `json:"state,omitempty"`
	Wait       *script.WaitCmd       `json:"wait,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for Entry.
// It wraps the command in a discriminator object (e.g., {"input": {...}}).
func (e Entry) MarshalJSON() ([]byte, error) {
	type Alias Entry
	aux := &struct {
		Command json.RawMessage `json:"command"`
		Alias
	}{
		Alias: Alias(e),
	}

	// Create wrapper based on command type
	var wrapper entryCommandWrapper
	switch cmd := e.Command.(type) {
	case *script.InputCmd:
		wrapper.Input = cmd
	case *script.MouseCmd:
		wrapper.Mouse = cmd
	case *script.WheelCmd:
		wrapper.Wheel = cmd
	case *script.ScreenshotCmd:
		wrapper.Screenshot = cmd
	case *script.DelayCmd:
		wrapper.Delay = cmd
	case *script.CustomCmd:
		wrapper.Custom = cmd
	case *script.RepeatCmd:
		wrapper.Repeat = cmd
	case *script.StateCmd:
		wrapper.State = cmd
	case *script.WaitCmd:
		wrapper.Wait = cmd
	default:
		return nil, fmt.Errorf("unknown command type: %T", e.Command)
	}

	cmdBytes, err := json.Marshal(wrapper)
	if err != nil {
		return nil, err
	}
	aux.Command = cmdBytes

	return json.Marshal(aux)
}

// UnmarshalJSON implements custom JSON unmarshaling for Entry.
// It extracts the command from the discriminator object.
func (e *Entry) UnmarshalJSON(data []byte) error {
	type Alias Entry
	aux := &struct {
		Command json.RawMessage `json:"command"`
		Alias
	}{
		Alias: Alias{
			Timestamp: e.Timestamp,
		},
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Unmarshal command using the script package's helper
	cmd, err := script.UnmarshalCommand(aux.Command)
	if err != nil {
		return err
	}

	e.Timestamp = aux.Alias.Timestamp
	e.Command = cmd

	return nil
}

// Entry represents a single recorded command with timestamp.
type Entry struct {
	Timestamp time.Time             `json:"timestamp"`
	Command   script.CommandWrapper `json:"command"`
}

// Recorder handles appending commands to recording file.
type Recorder struct {
	socketPath string
}

// NewRecorderFromSocket creates a recorder using a socket path.
// The recording file path is derived from the socket path.
func NewRecorderFromSocket(socketPath string) *Recorder {
	return &Recorder{socketPath: socketPath}
}

// Record appends a command with current timestamp.
// Uses flock for concurrency safety.
func (r *Recorder) Record(cmd script.CommandWrapper) error {
	path := PathFromSocket(r.socketPath)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create recording directory: %w", err)
	}

	// Open file with append mode
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open recording file: %w", err)
	}
	defer f.Close()

	// Acquire exclusive lock
	if err := unix.Flock(int(f.Fd()), unix.LOCK_EX); err != nil {
		return fmt.Errorf("failed to lock recording file: %w", err)
	}
	defer unix.Flock(int(f.Fd()), unix.LOCK_UN)

	// Create entry
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Command:   cmd,
	}

	// Marshal and write
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal entry: %w", err)
	}

	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write entry: %w", err)
	}
	if err := f.Sync(); err != nil {
		return fmt.Errorf("failed to sync recording file: %w", err)
	}
	return nil
}

// Clear removes the recording file derived from the socket path.
func Clear(socketPath string) error {
	path := PathFromSocket(socketPath)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to clear recording: %w", err)
	}
	return nil
}
