# Game Output Capture Implementation Plan

## Overview

Enable autoebiten CLI to capture game stdout/stderr and display changes using `diff -u` format between commands.

## Reference

- Design: [docs/superpowers/specs/2026-04-07-game-output-capture-design.md](../specs/2026-04-07-game-output-capture-design.md)

## Implementation Steps

### Phase 1: Core Infrastructure

#### Step 1: Create Output Capture Package

Create `internal/output/` package for log and snapshot file management.

**Files to create:**
- `internal/output/output.go` - Core types and file path derivation
- `internal/output/output_test.go` - Unit tests

**Key functions:**
```go
// FilePath returns paths for log, snapshot, and launch socket
type FilePath struct {
    Log       string  // autoebiten-{PID}-output.log
    Snapshot  string  // autoebiten-{PID}-snapshot.log
    LaunchSock string // autoebiten-{PID}-launch.sock
}

func DerivePaths(socketPath string) *FilePath
func CreateLogFile(path string) (*os.File, error)
func ReadSnapshot(path string) ([]byte, error)  // Returns empty if missing
func WriteSnapshot(path string, content []byte) error
func ProcessCarriageReturn(content []byte) []byte  // Handle \r to show final state
func GenerateDiff(snapshot, current []byte) string  // Uses diff -u format
```

**Testing:**
- Test path derivation from socket path
- Test snapshot read/write with missing files
- Test carriage return processing
- Test diff generation (verify `diff -u` format)

---

#### Step 2: Create Launch Proxy RPC Server

Create proxy RPC server that wraps game commands and captures output.

**Files to create:**
- `internal/proxy/server.go` - Proxy RPC server implementation
- `internal/proxy/server_test.go` - Unit tests

**Key types:**
```go
type ProxyServer struct {
    gameClient   *rpc.Client
    outputFiles  *output.FilePath
    mu           sync.Mutex
}

type ProxyHandler struct {
    server *ProxyServer
}
```

**RPC Methods (proxy prefix):**
- `proxy_input` → forward to game `input`
- `proxy_mouse` → forward to game `mouse`
- `proxy_wheel` → forward to game `wheel`
- `proxy_screenshot` → forward to game `screenshot`
- `proxy_ping` → forward to game `ping`
- `proxy_custom` → forward to game `custom`
- `proxy_exit` → forward to game `exit`, then cleanup and terminate

**Each proxy method flow:**
1. Read snapshot file
2. Forward command to game RPC
3. Wait for response
4. Read current log file
5. Generate diff
6. Write snapshot for next command
7. Return response with diff output

**Testing:**
- Mock game client, test proxy forwarding
- Test diff output in response
- Test concurrent requests are serialized

---

### Phase 2: Launch Command

#### Step 3: Create Launch Command

Add `autoebiten launch -- ./game [args...]` command to CLI.

**Files to modify:**
- `cmd/autoebiten/main.go` - Add launch command registration
- `internal/cli/launch.go` - Launch command implementation (NEW)
- `internal/cli/launch_test.go` - Unit tests

**Launch command behavior:**
1. Parse game command and arguments from `--` separator
2. Create log file and snapshot file paths
3. Launch game as subprocess with stdout/stderr piped
4. Pipe stdout/stderr to both terminal (delegate) and log file
5. Start proxy RPC server on launch socket
6. Wait for game to exit or interrupt signal
7. On game crash: wait 30s for CLI to read final output
8. Cleanup files on exit

**Key challenges:**
- Stdout/stderr delegation while capturing (use `io.MultiWriter` or tee pattern)
- Handle terminal TTY passthrough for interactive games
- Signal handling (Ctrl+C terminates immediately)

**Testing:**
- Test game subprocess launching
- Test output capture with mock subprocess
- Test cleanup on exit/crash

---

### Phase 3: CLI Integration

#### Step 4: Add Launch Proxy Detection to CLI Commands

Modify CLI commands to auto-detect and use launch proxy when available.

**Files to modify:**
- `internal/cli/commands.go` - Add proxy detection logic
- `internal/cli/state.go` - Add proxy detection
- `internal/cli/wait.go` - Add proxy detection
- `internal/rpc/socket.go` - Add `LaunchSocketPath()` helper

**Auto-detection logic:**
```go
func detectProxyOrDirect() (*rpc.Client, error) {
    launchSock := rpc.LaunchSocketPath(rpc.SocketPath())
    if fileExists(launchSock) {
        return rpc.NewClientWithPath(launchSock), nil
    }
    // No proxy, connect directly with tip
    fmt.Println("Tip: Use 'autoebiten launch -- ./game' to capture game output between commands.")
    return rpc.NewClient(), nil
}
```

**Testing:**
- Test proxy detection with/without launch socket
- Test direct connection fallback

---

#### Step 5: Display Diff Output in CLI Responses

Modify CLI Writer to display diff output before command result.

**Files to modify:**
- `internal/cli/writer.go` - Add diff output formatting

**Output format:**
```
--- snapshot    2026-04-07 10:00:00.000000000 +0000
+++ current     2026-04-07 10:00:01.000000000 +0000
@@ -1,2 +1,4 @@
 [INFO] Game started
-Loading: 50%
+Loading: 100%
+KeySpace pressed
+Player jumped!

OK: input press KeySpace
```

**Testing:**
- Test diff output formatting
- Test empty diff (no changes)

---

### Phase 4: Process Management

#### Step 6: Handle Game Exit and Cleanup

Implement graceful shutdown and cleanup logic.

**Files to modify:**
- `internal/cli/launch.go` - Add exit handling
- `internal/proxy/server.go` - Add cleanup methods

**Scenarios:**
| Scenario | Behavior |
|----------|----------|
| Game exits normally | Wait 30s for CLI reads, then cleanup |
| Game crashes | Same as above |
| Ctrl+C on launch | Terminate game immediately, cleanup |
| `proxy_exit` called | Forward to game, cleanup after 30s |

**Cleanup actions:**
- Close proxy RPC server
- Delete log file
- Delete snapshot file
- Delete launch socket

**Testing:**
- Test cleanup on normal exit
- Test cleanup on crash
- Test cleanup on interrupt

---

### Phase 5: End-to-End Testing

#### Step 7: Integration Tests

Add end-to-end tests for the complete flow.

**Files to modify:**
- `e2e/e2e_test.go` - Add launch/output capture tests

**Test cases:**
1. Launch game, send commands, verify diff output
2. Game crash scenario - verify final output captured
3. Multiple commands in sequence - verify snapshot updates
4. Interactive game with carriage returns

---

## File Summary

### New Files
- `internal/output/output.go`
- `internal/output/output_test.go`
- `internal/proxy/server.go`
- `internal/proxy/server_test.go`
- `internal/cli/launch.go`
- `internal/cli/launch_test.go`

### Modified Files
- `cmd/autoebiten/main.go` - Add launch command
- `internal/cli/commands.go` - Add proxy detection
- `internal/cli/state.go` - Add proxy detection
- `internal/cli/wait.go` - Add proxy detection
- `internal/cli/writer.go` - Add diff output formatting
- `internal/rpc/socket.go` - Add `LaunchSocketPath()`
- `e2e/e2e_test.go` - Add integration tests

## Dependencies

No new external dependencies. Uses standard Go packages:
- `os/exec` for subprocess management
- `io` for tee pattern
- `syscall` for signal handling

## Testing Strategy

1. **Unit tests** for each component (output package, proxy server)
2. **Integration tests** with mock game subprocess
3. **E2E tests** with actual test game binary

## Risks and Mitigations

| Risk | Mitigation |
|------|------------|
| Terminal TTY handling for interactive games | Use `os.Stdin` passthrough, only capture stdout/stderr |
| Race conditions with concurrent CLI commands | Proxy serializes all requests through single mutex |
| Log file permissions | Use 0600 permissions for security |
| Large log files for long-running games | Document non-goal: no rotation in this version |

## Success Criteria

1. `autoebiten launch -- ./game` captures output
2. CLI commands show diff output when proxy is active
3. Game crashes still allow reading final output
4. No race conditions or deadlocks
5. All tests pass with 80%+ coverage