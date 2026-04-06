# Recording and Replay Feature Design

**Date:** 2025-04-06  
**Status:** Approved

## Overview

Add automatic recording of all CLI commands sent to a game, with replay and script export capabilities. Multiple CLI instances targeting the same game write to the same recording file using file locking.

## Goals

- Automatically record all CLI commands to a file per game process
- Support replay with timing and optional speed scaling
- Allow export to script format for editing
- Handle concurrent writes from multiple CLI instances

## Non-Goals

- Recording game state changes (only CLI commands)
- Visual recording (video)
- Recording across game restarts (one file per process)

## CLI Commands

| Command | Description |
|---------|-------------|
| `autoebiten input --key KeyA` | Auto-records to `/tmp/autoebiten/recording-{PID}.jsonl` after server response |
| `autoebiten input --key KeyA --no-record` | Executes command but skips recording |
| `autoebiten clear_recording` | Clears recording file for current game |
| `autoebiten replay` | Generates script from recording with original timing and executes via `run` |
| `autoebiten replay --speed 2` | Generates script with 2x speed scaling (delays halved) |
| `autoebiten replay --dump script.json` | Generates script file without execution |

## Recording Flow

```
CLI sends command
    ↓
Server processes
    ↓
CLI receives response
    ↓
CLI appends command + timestamp to recording file (with flock)
```

## File Format

**Location:** `/tmp/autoebiten/recording-{PID}.jsonl`

**Format:** NDJSON (newline-delimited JSON), one entry per command:

```json
{"timestamp":"2025-04-06T14:30:22.123456789Z","command":{"input":{"action":"press","key":"KeyA"}}}
{"timestamp":"2025-04-06T14:30:22.456789012Z","command":{"mouse":{"action":"position","x":100,"y":200}}}
{"timestamp":"2025-04-06T14:30:23.000000000Z","command":{"screenshot":{"output":"shot.png"}}}
```

- `timestamp`: RFC3339Nano format for nanosecond precision
- `command`: Same JSON structure as script commands

## Concurrency Strategy

Use Unix advisory file locking (`flock`) to coordinate multiple CLI instances:

1. Open file with `O_APPEND` flag
2. Acquire exclusive lock with `flock(fd, LOCK_EX)`
3. Write entry + newline
4. Release lock with `flock(fd, LOCK_UN)`
5. Close file

This ensures atomic appends even with multiple concurrent writers.

## Commands Not Recorded

These query-only commands are not recorded:
- `ping`, `version`
- `keys`, `mouse_buttons`
- `get_mouse_position`, `get_wheel_position`
- `list_custom`
- `clear_recording`, `replay` (meta-commands)
- `schema`

## Replay Behavior

The `replay` command:

1. **Read** the recording file line by line
2. **Parse** each entry into timestamp + command
3. **Calculate** delays between consecutive commands (current - previous timestamp)
4. **Scale** delays by speed factor (delay / speed)
5. **Generate** script structure:
   - Start with first command
   - Insert `{"delay": {"ms": N}}` between commands
   - Add subsequent command
6. **Execute** via existing `run` logic (or dump if `--dump` specified)

### Speed Scaling Example

Original recording:
```
T0: input KeyA
T0+500ms: mouse position
T0+800ms: screenshot
```

Generated script with `--speed 2`:
```json
{
  "version": "1.0",
  "commands": [
    {"input": {"action": "press", "key": "KeyA"}},
    {"delay": {"ms": 250}},
    {"mouse": {"action": "position", "x": 100, "y": 200}},
    {"delay": {"ms": 150}},
    {"screenshot": {"output": "shot.png"}}
  ]
}
```

## Architecture

### New Package: `internal/recording/`

```
internal/recording/
├── recorder.go      # Append commands to file with flock
├── reader.go        # Read and parse recording entries
├── scriptgen.go     # Generate script from recording entries
└── recording_test.go
```

### `recorder.go`

```go
// Recorder handles appending commands to recording file
type Recorder struct {
    pid int
}

// NewRecorder creates a recorder for the given PID
func NewRecorder(pid int) *Recorder

// Record appends a command with current timestamp
// Uses flock for concurrency safety
func (r *Recorder) Record(cmd script.CommandWrapper) error

// Clear removes the recording file for the given PID
func Clear(pid int) error

// Path returns the recording file path for a PID
func Path(pid int) string
```

### `reader.go`

```go
// Entry represents a single recorded command
type Entry struct {
    Timestamp time.Time
    Command   script.CommandWrapper
}

// Reader reads recording files
type Reader struct {
    pid int
}

// NewReader creates a reader for the given PID
func NewReader(pid int) *Reader

// ReadAll reads all entries from the recording file
func (r *Reader) ReadAll() ([]Entry, error)
```

### `scriptgen.go`

```go
// Generator creates scripts from recording entries
type Generator struct {
    speed float64  // 1.0 = original, 2.0 = 2x speed, 0.5 = half speed
}

// NewGenerator creates a generator with speed multiplier
func NewGenerator(speed float64) *Generator

// Generate creates a Script from recording entries
func (g *Generator) Generate(entries []recording.Entry) (*script.Script, error)
```

### CLI Integration

**`cmd/autoebiten/main.go` changes:**

1. Add `--no-record` flag (applies to commands that record)
2. Add `clear_recording` subcommand
3. Add `replay` subcommand with `--speed` and `--dump` flags

**`internal/cli/commands.go` changes:**

Modify command methods to record after successful execution:

```go
func (e *CommandExecutor) RunInputCommand(...) error {
    // ... existing execution code ...

    // Record after successful execution
    if shouldRecord {
        recorder := recording.NewRecorder(rpc.GetTargetPID())
        recorder.Record(&script.InputCmd{...})
    }

    return nil
}
```

## Edge Cases

### Empty Recording

If `replay` is called with no recording file or empty file:
- Return error: "no recording found for game"

### Missing Recording File

If recording file doesn't exist when `replay` is called:
- Same error as empty recording

### Corrupted Entry

If a line in the recording file is invalid JSON:
- Skip that entry and log warning
- Continue processing remaining entries

### Speed = 0

If `--speed 0` is passed:
- Return error: "speed must be greater than 0"

### Very Large Recordings

For recordings with thousands of commands:
- Stream read file line by line (don't load all into memory)
- Generator should handle large entry sets efficiently

## Testing Strategy

1. **Unit tests:**
   - `recorder.Record()` writes valid JSON
   - `recorder.Record()` handles concurrent writes correctly
   - `reader.ReadAll()` parses all entries
   - `generator.Generate()` creates correct delays
   - Speed scaling produces correct delays

2. **Integration tests:**
   - Multiple CLI instances write to same file
   - Replay generates executable script
   - Clear removes file correctly

## Security Considerations

- Recording file written to `/tmp/autoebiten/` (world-writable directory)
- Use PID in filename to avoid collisions
- No sensitive data expected in commands (keys, mouse positions)
- File locking only advisory (other processes can bypass)

## Future Extensions

- `--start-at` and `--end-at` flags for replay (time ranges)
- Pause/resume recording
- Multiple recording slots per game
- Compress old recordings
