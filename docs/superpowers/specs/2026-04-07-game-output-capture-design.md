# Game Output Capture Design

## Overview

Enable autoebiten CLI to capture game stdout/stderr and display changes using `diff -u` format between commands.

## Goals

1. Capture game output (stdout/stderr) when launched via `autoebiten launch`
2. Show users what happened in the game between CLI commands using `diff -u`
3. Handle `\r` (carriage return) correctly by showing final state
4. Support concurrent CLI commands without race conditions
5. Delegate stdin/stdout/stderr so interactive games work normally
6. Respect `AUTOEBITEN_SOCKET` environment variable for file locations

## Non-Goals

1. Capture output from externally-started games (only works with `autoebiten launch`)
2. Log rotation or truncation (keep it simple)

## Architecture

### Launch as Proxy Pattern

The `autoebiten launch` command acts as an intermediary between CLI commands and the game:

```
┌─────────────────┐      ┌──────────────────┐      ┌─────────────┐
│   CLI Commands  │ ───→ │ autoebiten launch │ ───→ │  Game RPC   │
│  (input, ping)  │      │   (Proxy/Owner)   │      │   Server    │
└─────────────────┘      └──────────────────┘      └─────────────┘
                                │
                                ↓ captures stdout/stderr
                         ┌─────────────┐
                         │  log file   │
                         │ + snapshot  │
                         └─────────────┘
```

**Benefits:**
- Serializes snapshot access (no concurrency issues)
- Centralizes log management
- Clean separation of concerns

## Components

### 1. Launch Command

**Command:** `autoebiten launch -- ./game [args...]`

**Behavior:**
1. Launches game as child process with stdin/stdout/stderr delegated to terminal
2. Captures stdout/stderr to log file using pipe + tee pattern
3. Starts RPC server on launch socket
4. Forwards commands to game's RPC server
5. If game crashes: keeps launch process running for 30s to allow CLI to read final output
6. On game exit/crash + 30s timeout OR explicit CLI termination request: cleans up files and exits
7. On interrupt (Ctrl+C): terminates game immediately and cleans up

**File Locations:**

File paths follow the same pattern as the recording feature, derived from `rpc.SocketPath()`:

```
Socket:      {dir}/autoebiten-{PID}.sock
Log file:    {dir}/autoebiten-{PID}-output.log
Snapshot:    {dir}/autoebiten-{PID}-snapshot.log
Launch sock: {dir}/autoebiten-{PID}-launch.sock
```

This respects the `AUTOEBITEN_SOCKET` environment variable.

### 2. Proxy RPC Server

The launch process exposes an RPC server that wraps game commands:

**Methods:**
- `proxy_input` → forwards to game `input`, returns output + result
- `proxy_mouse` → forwards to game `mouse`, returns output + result
- `proxy_wheel` → forwards to game `wheel`, returns output + result
- `proxy_screenshot` → forwards to game `screenshot`, returns output + result
- `proxy_ping` → forwards to game `ping`, returns output + result
- `proxy_custom` → forwards to game `custom`, returns output + result (covers state and wait-for CLI commands)
- `proxy_exit` → forwards to game `exit`, cleans up files and exits

**Output Capture Flow:**
1. Read snapshot file (copy of log at last command)
2. Send command to game RPC server
3. Wait for response
4. Read current log file
5. Run `diff -u snapshot.log current.log` to generate diff
6. Display diff output
7. Copy current log to snapshot file for next command

**Note:** Each CLI command is serialized through the launch proxy, eliminating concurrency issues.

### 3. Enhanced CLI Commands

All existing commands (input, mouse, ping, etc.) check for launch proxy first:

**Auto-detection Logic:**
1. Derive launch socket path from `rpc.SocketPath()` (append `-launch.sock`)
2. If launch socket exists: connect to launch proxy
3. If not exists: connect directly to game (with warning)

**Output Format:**

Uses `diff -u` format to show changes between snapshot and current log state:

**Single command example:**
```
$ autoebiten input --key KeySpace
--- snapshot	2026-04-07 10:00:00.000000000 +0000
+++ current	2026-04-07 10:00:01.000000000 +0000
@@ -1,2 +1,4 @@
 [INFO] Game started
-Loading: 50%
+Loading: 100%
+KeySpace pressed
+Player jumped!

OK: input press KeySpace
```

**Script command example** (continuous stream, diff shown after each command):
```
$ autoebiten run --script test.json
--- snapshot	2026-04-07 10:00:00.000000000 +0000
+++ current	2026-04-07 10:00:01.000000000 +0000
@@ -0,0 +1,2 @@
+[INFO] Game started
+Loading: 100%

OK: game is running
--- snapshot	2026-04-07 10:00:01.000000000 +0000
+++ current	2026-04-07 10:00:02.000000000 +0000
@@ -2,2 +2,4 @@
 [INFO] Game started
 Loading: 100%
+KeySpace pressed
+Player jumped!

OK: input press KeySpace
```

When game started externally (no output capture):
```
$ autoebiten ping
Tip: Use 'autoebiten launch -- ./game' to capture game output between commands.

OK: game is running
```

## Handling Carriage Return (\r)

Games use `\r` for progress bars and animations:

```
Loading: 0%\rLoading: 50%\rLoading: 100%\nDone!\n
```

**Approach:**
- The log file contains the raw output (with `\r` characters)
- When reading the log, we process `\r` to get final line states
- `diff -u` naturally shows modified lines with `-` and `+` markers
- Users see the diff showing `Loading: 50%` changed to `Loading: 100%`

**Example:**
```
$ autoebiten ping
--- snapshot	2026-04-07 10:00:00.000000000 +0000
+++ current	2026-04-07 10:00:01.000000000 +0000
@@ -0,0 +1,3 @@
+[INFO] Game started
+Loading: 100%
+Done!

OK: game is running
```

## File Format

### Log File

Plain text file containing captured stdout/stderr:
```
[INFO] Game started
[DEBUG] Initializing renderer
Loading: 100%
[INFO] Server listening on :8080
```

### Snapshot File

Copy of log file at last command completion. Used as the "before" state for `diff -u`.

**Location:** Same directory as log file, with `-snapshot.log` suffix.

**Example:**
- Log file: `/tmp/autoebiten/autoebiten-12345-output.log`
- Snapshot: `/tmp/autoebiten/autoebiten-12345-snapshot.log`

## Error Handling

| Scenario | Behavior |
|----------|----------|
| Game crashes | Launch process keeps running, waiting for CLI to read final output. Auto-terminates after 30s if not read. |
| Launch proxy dies | CLI commands fall back to direct game connection with friendly tip about 'autoebiten launch' |
| Log file missing | Treat as empty, diff shows no changes |
| Snapshot file missing | Treat as empty file, diff shows all lines as new |

## Implementation Plan

1. **Create launch command** - New cobra command, process management
2. **Create proxy RPC server** - Wrap game RPC, handle diff logic
3. **Modify CLI commands** - Add launch proxy detection
4. **Add file management** - Create/cleanup log and snapshot files
5. **Testing** - Unit tests for diff logic, integration tests for proxy

## Security Considerations

1. Log files written to `/tmp/autoebiten/` with restrictive permissions (0600)
2. No sensitive data filtering (games should not log secrets)
3. Snapshot files contain same data as logs (same permissions)

## Future Enhancements (Out of Scope)

1. Log rotation for long-running games
2. Filtering/sanitizing log output
3. Structured logging (JSON format)
4. Web UI for viewing logs
