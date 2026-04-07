# Design: Carriage Return Interpretation + Mutex-Protected Output

**Date:** 2026-04-07
**Status:** Approved

## Problem

The current implementation writes raw game stdout/stderr to a log file, including `\r` (carriage return) characters. This causes:

1. **Log file doesn't match terminal visual state** — `\r` in terminal means "move cursor to line start, overwrite". The file keeps all bytes, making it unreadable.
2. **Diff workaround needed** — `output.go:GenerateDiff` uses `sed` to strip `\r` content as a post-processing hack.
3. **No synchronization** — `teeOutput` writes to log file concurrently while `ForwardRequest` reads for diff/snapshot. No mutex coordination.

## Goal

Make the log file reflect "what the terminal showed" (human-readable, snapshot-comparable) while ensuring safe concurrent access for write, copy, and diff operations.

## Solution Overview

| Component | Purpose |
|-----------|---------|
| `CarriageReturnWriter` | Interprets `\r` as line overwrite, produces "visual state" output |
| `OutputManager` | Mutex-protected coordinator for write/diff/snapshot operations |

## Architecture

### CarriageReturnWriter

Custom `io.Writer` that interprets `\r` before writing to file.

**Algorithm:**
```
For each byte in input:
  if byte == '\r':
    clear lineBuf (simulate cursor moving to line start)
  elif byte == '\n':
    dst.Write(lineBuf + '\n')  // flush completed line
    clear lineBuf
  else:
    append byte to lineBuf

At end of stream:
  if lineBuf has content:
    dst.Write(lineBuf + '\n')  // flush incomplete line
```

**Example:**
- Input: `"Loading 10%\rLoading 20%\nDone\n"`
- Output: `"Loading 20%\nDone\n"` (matches terminal visual state)

**Edge cases:**
- Multiple `\r` before `\n` → only last segment survives (correct overwrite)
- Stream ends without `\n` → flush incomplete line with `\n` appended

### OutputManager

Mutex-protected coordinator for all file operations.

**Struct:**
```go
type OutputManager struct {
    mu           sync.Mutex
    logFile      *os.File
    outputPath   string
    snapshotPath string
}
```

**Methods:**
- `Write(data []byte)` — Lock, write to log file, flush (`Sync()`)
- `DiffAndUpdateSnapshot()` — Lock, flush, generate diff, copy log to snapshot, return diff string

### Integration Flow

**Write path (game → file):**
```
Game stdout/stderr
    ↓ (bufio.Reader)
teeOutput (launch.go)
    ↓ (raw bytes)
CarriageReturnWriter
    ↓ (interpreted bytes)
OutputManager.Write()
    ↓ (mutex.Lock → file.Write → file.Sync → mutex.Unlock)
log file
```

**Diff/snapshot path (RPC request):**
```
CLI sends RPC request
    ↓
ForwardRequest()
    ↓
OutputManager.DiffAndUpdateSnapshot()
    ↓ (mutex.Lock → flush → diff → copy → mutex.Unlock)
Return diff to CLI
```

## File Changes

| File | Change |
|------|--------|
| `internal/output/manager.go` | **New**: `OutputManager` + `CarriageReturnWriter` |
| `internal/output/output.go` | Remove `GenerateDiff`, `CopyFile`; keep `FilePath`, `DerivePaths`, `CreateLogFile` |
| `internal/cli/launch.go` | Replace `logFile.Write()` with `OutputManager.Write()` via `CarriageReturnWriter` |
| `internal/proxy/server.go` | Replace `output.GenerateDiff` + `output.CopyFile` with `OutputManager.DiffAndUpdateSnapshot()` |

## Testing

**Unit tests for `CarriageReturnWriter`:**
- `\r` clears line buffer (overwrite)
- `\n` flushes line buffer
- Multiple `\r` before `\n` → only last segment survives
- Stream ends without `\n` → flush incomplete line
- Empty input → no write

**Unit tests for `OutputManager`:**
- Concurrent `Write()` calls are serialized (mutex works)
- `DiffAndUpdateSnapshot()` produces correct diff after writes
- `DiffAndUpdateSnapshot()` updates snapshot correctly

**Integration test:**
- Launch game, capture output with `\r`, verify log file matches terminal visual state
- Send RPC request, verify diff is correct, verify snapshot updated