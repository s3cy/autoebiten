---
name: Launch Socket for Crash Diagnostics
description: Enable CLI commands to query launch process for error diagnostics when game crashes before or after RPC connection
type: project
---

# Launch Socket for Crash Diagnostics

## Problem

When running `autoebiten launch -- ./mygame &` in background, the user has no visibility into crash scenarios:
- Game crashes before RPC connection (e.g., executable not found, instant crash)
- Game crashes after RPC connection

Current behavior: launch exits immediately on pre-RPC crash, or waits 30s silently on post-RPC crash. CLI commands like `ping` cannot retrieve error information.

## Solution

Create a single launch socket that persists throughout the launch lifecycle. CLI commands connect to this socket and receive error diagnostics (Go error chain + game output log diff) when game is not connected.

## Design

### Socket Naming

**Single launch socket**: `autoebiten-{LAUNCH_PID}-launch.sock`
- Created by launch process before game starts
- Persists throughout launch lifecycle
- LAUNCH_PID is the autoebiten launch process PID

**Game socket**: `autoebiten-{LAUNCH_PID}.sock`
- Game's RPC socket, using LAUNCH_PID from `AUTOEBITEN_SOCKET` env
- Launch process sets `AUTOEBITEN_SOCKET=/tmp/autoebiten/autoebiten-{LAUNCH_PID}.sock` in game's environment

**Why use LAUNCH_PID for both**:
- Game socket uses LAUNCH_PID (not GAME_PID) because we control it via env
- This keeps naming consistent and simplifies discovery
- One launch process = one game process (1:1 relationship)

### CLI Discovery Logic

**With --pid**:
1. Try `autoebiten-{PID}-launch.sock` (user may specify LAUNCH_PID)
2. Fallback to `autoebiten-{PID}.sock` (direct game connection)

**With AUTOEBITEN_SOCKET env** (path = `{path}`):
1. Try `{path}-launch.sock`
2. Fallback to `{path}.sock`

**Auto-select** (no --pid, no env):
1. Scan all sockets in `/tmp/autoebiten/`
2. Dedup: for same PID, if both `launch.sock` and `.sock` exist → show only `launch`
3. Show socket type in list: "PID 12345: launch" / "PID 12345: game"
4. If exactly one unique PID → use its preferred socket (launch > game)
5. If multiple unique PIDs → show list with types, require `--pid`

### Proxy State Machine

```go
type ProxyState int

const (
    StateWaiting ProxyState = iota  // Waiting for game RPC connection
    StateConnected ProxyState = 1    // Game connected, proxy active
    StateCrashed  ProxyState = 2    // Game crashed/exited
)
```

**State transitions**:
- `[Launch starts]` → `StateWaiting`
- `StateWaiting` → `StateConnected` (game RPC connects successfully)
- `StateWaiting` → `StateCrashed` (game exits before RPC)
- `StateConnected` → `StateCrashed` (game exits after RPC)
- `StateCrashed` → `[Launch terminates]` (CLI query received OR 30s timeout)

**Error accumulation**:
- `launchError error` field in handler, wrapped throughout lifecycle using `fmt.Errorf("context: %w", err)`
- Error points: game start failure, game exit, RPC timeout, game request failure

**State behavior**:
| State | RPC Request Handling |
|-------|---------------------|
| Waiting | Return empty proxy_error + log_diff + "Error: game not connected" |
| Connected | Forward to game, return response with log_diff |
| Crashed | Return proxy_error (accumulated) + log_diff + "Error: game not connected" |

### Unified Handler

**Handler structure** (rewrite `internal/proxy/server.go`):
```go
type UnifiedHandler struct {
    state       ProxyState
    gameClient  GameClient      // nil until StateConnected
    outputMgr   *OutputManager
    launchError error           // accumulated error chain
    onCrashed   func()          // callback to signal launch exit
    mu          sync.Mutex
}
```

**Request processing**:
1. Lock mutex
2. Check state:
   - `Waiting/Crashed`: generate diff, build error response, signal exit if crashed
   - `Connected`: forward to game, if request fails → transition to crashed, return error
3. Unlock mutex
4. Return response with Extra fields ("diff", optionally "proxy_error")

**Error response**:
- RPC error code: `-32001` (ErrGameNotConnected)
- Extra["proxy_error"]: accumulated error chain (empty if waiting, populated if crashed)
- Extra["diff"]: log diff (always generated)

### Launch Lifecycle

```go
func (lc *LaunchCommand) Run() error {
    // 1. Create launch socket (before game)
    launchSock := fmt.Sprintf("autoebiten-%d-launch.sock", os.Getpid())
    lc.createSocket(launchSock)
    
    // 2. Create game command with AUTOEBITEN_SOCKET env
    gameSock := fmt.Sprintf("autoebiten-%d.sock", os.Getpid())
    cmd.Env = append(os.Environ(), fmt.Sprintf("AUTOEBITEN_SOCKET=%s", gameSock))
    
    // 3. Start game, capture stdout/stderr to OutputManager
    if err := cmd.Start(); err != nil {
        handler.TransitionToCrashed(fmt.Errorf("failed to start game: %w", err))
        lc.waitForExit()  // Wait for CLI query or 30s
        return err
    }
    
    // 4. Monitor game exit in goroutine
    go func() {
        cmd.Wait()
        handler.TransitionToCrashed(fmt.Errorf("game exited: %s", cmd.ProcessState))
        close(lc.gameExited)
    }()
    
    // 5. Wait for game RPC (with timeout)
    client, err := lc.waitForGameRPC()
    if err != nil {
        // Error already set by exit monitor or timeout
        lc.waitForExit()
        return err
    }
    
    // 6. Transition to connected
    handler.TransitionToConnected(client)
    
    // 7. Wait for game exit
    <-lc.gameExited
    
    // 8. Wait for CLI query or timeout
    lc.waitForExit()
}
```

**Wait behavior after crash**:
- `waitForExit()` waits for: CLI query signal OR exit command OR 30s timeout
- Exit immediately after CLI receives error response (no unnecessary waiting)

### CLI Output Format

**Unified response handler**:
```go
func (e *CommandExecutor) handleResponse(resp *rpc.RPCResponse, onSuccess func()) {
    // Always print log_diff if present
    if diff, ok := resp.Extra["diff"].(string); ok && diff != "" {
        fmt.Fprintf(os.Stdout, "<log_diff>\n%s\n</log_diff>\n", diff)
    }
    
    // Always print proxy_error if present
    if proxyError, ok := resp.Extra["proxy_error"].(string); ok && proxyError != "" {
        fmt.Fprintf(os.Stdout, "<proxy_error>\n%s\n</proxy_error>\n", proxyError)
    }
    
    // Handle result/error
    if resp.Error != nil {
        fmt.Fprintf(os.Stderr, "Error: %s\n", resp.Error.Message)
    } else {
        onSuccess()
    }
}
```

**Example output (crash)**:
```
<proxy_error>
game exited before RPC connection: exit status 127
</proxy_error>
<log_diff>
--- snapshot (empty)
+++ current 2024-01-08 10:30:15.123
@@ -0,0 +1,2 @@
+./mygame: command not found
+
</log_diff>
Error: game not connected
```

**Example output (waiting, no crash yet)**:
```
<log_diff>
</log_diff>
Error: game not connected
```

**Example output (success)**:
```
<log_diff>
--- snapshot ...
+++ current ...
@@ ...
+game output line
</log_diff>
game is running
```

### RPC Error Codes

Add to `internal/rpc/messages.go`:
```go
const ErrGameNotConnected = -32001
```

### Edge Cases

| Scenario | Behavior |
|----------|----------|
| Game exits instantly | Socket exists, StateCrashed, CLI gets error + log_diff |
| Executable not found | `Start()` fails, StateCrashed, CLI gets wrapped error |
| Game hangs, RPC timeout | StateCrashed with timeout error, CLI gets error + log_diff |
| Game crashes post-RPC | Connected → Crashed, next CLI command gets error + log_diff |
| CLI connects during RPC wait | StateWaiting, empty proxy_error, "game not connected" |
| Multiple CLI queries | Each gets same proxy_error + updated log_diff |
| Launch cleanup | Remove launch socket on exit |

### Files to Modify

| File | Changes |
|------|---------|
| `internal/cli/launch.go` | Rewrite with state machine, socket creation, lifecycle |
| `internal/cli/commands.go` | Replace per-command output with `handleResponse` callback |
| `internal/cli/writer.go` | Simplify or remove (handleResponse replaces methods) |
| `internal/proxy/server.go` | Rewrite as `UnifiedHandler` with state machine |
| `internal/rpc/socket.go` | Update discovery logic for launch socket detection |
| `internal/rpc/messages.go` | Add `ErrGameNotConnected` error code |

### Implementation Principles

- **Simplicity first**: Rewrite rather than workaround
- **Single socket**: No prelaunch/launch transition complexity
- **Unified handler**: One handler for all states
- **Unified output**: One `handleResponse` for all commands