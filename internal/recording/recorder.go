package recording

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/unix"

	"github.com/s3cy/autoebiten/internal/script"
)

var (
	// RecordingDir is the directory for recording files.
	RecordingDir = "/tmp/autoebiten"
)

// Path returns the recording file path for a PID.
func Path(pid int) string {
	return filepath.Join(RecordingDir, fmt.Sprintf("recording-%d.jsonl", pid))
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
	pid int
}

// NewRecorder creates a recorder for the given PID.
func NewRecorder(pid int) *Recorder {
	return &Recorder{pid: pid}
}

// Record appends a command with current timestamp.
// Uses flock for concurrency safety.
func (r *Recorder) Record(cmd script.CommandWrapper) error {
	path := Path(r.pid)

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

// Clear removes the recording file for the given PID.
func Clear(pid int) error {
	path := Path(pid)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to clear recording: %w", err)
	}
	return nil
}
