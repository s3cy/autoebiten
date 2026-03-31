# autoebiten Specification

## Overview

A utility package for Ebitengine that enables AI agents to automate games via CLI. Provides input injection, screenshot capture, and script execution capabilities.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     AI Agent / CLI                          │
└─────────────────────┬───────────────────────────────────────┘
                      │ JSON-RPC over Unix socket
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                    autoebiten CLI                           │
│         (input, mouse, wheel, screenshot, run, ping)       │
└─────────────────────┬───────────────────────────────────────┘
                      │ Unix socket (/tmp/autoebiten/autoebiten-{PID}.sock)
                      │ or $AUTOEBITEN_SOCKET
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                      Game Process                           │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                autoebiten package                   │    │
│  │  autoebiten.Update() → processes RPC commands       │    │
│  │  autoebiten.Capture() → screenshot capture          │    │
│  │  VirtualInput blends real + injected input          │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

## Components

### 1. Game Library (`autoebiten` package)

User imports and integrates into their game:

```go
import "github.com/s3cy/autoebiten"

func (g *Game) Update() error {
    if !autoebiten.Update() {
        return errors.New("exit requested")
    }
    // Use autoebiten.IsKeyPressed() instead of ebiten.IsKeyPressed()
    if autoebiten.IsKeyPressed(ebiten.KeySpace) {
        fmt.Println("Space is pressed!")
    }
}

func (g *Game) Draw(screen *ebiten.Image) {
    ...

    autoebiten.Capture(screen) // Screenshot captured at the end (screen is valid)
}
```

### Alternative: Deep Integration via Ebiten Patch

For games that cannot modify input calls, autoebiten provides a patch for Ebiten itself. The patch modifies Ebiten's internal `gameforui.go`, `input.go`, and `inpututil/inpututil.go` to call into the `integrate` package:

```go
// Ebiten's internal input functions call integrate.*
func IsKeyPressed(key Key) bool {
    if integrate.IsKeyPressed(integrate.Key(key)) {
        return true
    }
    return inputstate.Get().IsKeyPressed(ui.Key(key))
}
```

This requires:
1. Cloning Ebiten locally
2. Applying `ebiten.patch`
3. Using `replace` directive in go.mod
4. No game code changes required

See README.md for detailed patching instructions.

### 2. CLI Tool (`autoebiten`)

```bash
autoebiten input --key KeyA --action press
autoebiten input --key KeySpace --action hold --duration_ticks 6
autoebiten mouse --action press --x 100 --y 200 --button MouseButtonLeft
autoebiten wheel --x 0 --y -3
autoebiten screenshot --output shot.png
autoebiten run --script script.json
autoebiten ping
autoebiten keys
autoebiten mouse_buttons
```

## JSON-RPC Protocol

**Transport**: Unix domain socket

**Socket Path**: `/tmp/autoebiten/autoebiten-{PID}.sock` (configurable via `AUTOEBITEN_SOCKET` environment variable)

### Request Format
```json
{"jsonrpc": "2.0", "id": 1, "method": "input", "params": {...}}
```

### Response Format
```json
{"jsonrpc": "2.0", "id": 1, "result": {"success": true}}
```

### Error Format
```json
{"jsonrpc": "2.0", "id": 1, "error": {"code": -32001, "message": "invalid params"}}
```

### Methods

| Method | Description |
|--------|-------------|
| `input` | Inject keyboard input (press/release/hold) |
| `mouse` | Inject mouse input (position/press/release/hold) |
| `wheel` | Inject mouse wheel movement |
| `screenshot` | Capture game window to file |
| `ping` | Health check |
| `exit` | Request game to exit |

### Error Codes

| Code | Name | Description |
|------|------|-------------|
| -32000 | `ERR_CONNECTION_FAILED` | Could not connect to game |
| -32001 | `ERR_INVALID_PARAMS` | Invalid parameters |
| -32002 | `ERR_SCRIPT_FAILED` | Script execution failed |
| -32003 | `ERR_SCREENSHOT_FAILED` | Screenshot capture failed |
| -32004 | `ERR_GAME_NOT_RUNNING` | Game process not running |

## Script Format (JSON)

```json
{
  "version": "1.0",
  "commands": [
    {"input": {"action": "press", "key": "KeyW"}},
    {"delay": {"ms": 100}},
    {"repeat": {"times": 3, "commands": [
      {"input": {"action": "press", "key": "KeyA", "duration_ticks": 6}},
      {"delay": {"ms": 200}}
    ]}}
  ]
}
```

### Command Types

| Command | Fields | Description |
|---------|--------|-------------|
| `input` | `action` (press/release/hold), `key`, `duration_ticks` | Keyboard input |
| `mouse` | `action` (position/press/release/hold), `x`, `y`, `button`, `duration_ticks` | Mouse input |
| `wheel` | `x`, `y` | Mouse wheel movement |
| `screenshot` | `output` | Capture screenshot |
| `delay` | `ms` | Wait in milliseconds |
| `repeat` | `times`, `commands` | Repeat block N times |

## File Structure

```
autoebiten/
├── cmd/
│   └── autoebiten/
│       └── main.go              # CLI entry point
├── internal/
│   ├── cli/
│   │   ├── commands.go          # CLI command executor
│   │   └── writer.go            # Output formatting
│   ├── input/
│   │   ├── input.go             # VirtualInput, key/mouse state
│   │   ├── keys.go              # Key constant mappings
│   │   ├── mouse_buttons.go     # Mouse button constants
│   │   └── input_time.go        # Tick-based input timing
│   ├── rpc/
│   │   ├── handlers.go          # Command handlers registry
│   │   ├── messages.go          # RPC request/response types
│   │   └── socket.go            # Unix socket server and client
│   ├── script/
│   │   ├── ast.go               # Script AST nodes
│   │   ├── parser.go            # JSON script parser
│   │   ├── executor.go          # Script execution engine
│   │   └── parser_test.go       # Parser tests
│   └── server/
│       ├── server.go            # RPC request processing
│       ├── screenshot.go        # Screenshot capture
│       └── tick.go              # Tick management
├── integrate/
│   └── integrate.go             # Low-level integration API for Ebiten patch
├── examples/
│   └── simple/
│       └── main.go              # Example game
├── e2e/
│   └── e2e_test.go              # End-to-end tests
├── ebiten.patch                 # Patch for Ebiten v2.9.9 deep integration
├── autoebiten.go                # Mode configuration
├── autoebiten_default.go        # Default build (with RPC server)
├── autoebiten_release.go        # Release build (no-op stubs)
├── go.mod
├── README.md
└── SPEC.md
```

### Input Modes

- `InjectionOnly` - Returns only injected input results
- `InjectionFallback` (default) - Returns injected results if available, otherwise falls back to ebiten's native input
- `Passthrough` - Only uses ebiten's native input handling

### Input Function Support

The library provides mode-aware wrappers for both direct input queries and inpututil functions:

**Direct Input:**
- `IsKeyPressed(key Key) bool`
- `CursorPosition() (x, y int)`
- `Wheel() (xoff, yoff float64)`
- `IsMouseButtonPressed(mouseButton MouseButton) bool`

**inpututil Wrappers:**
- `IsKeyJustPressed(key Key) bool`
- `IsKeyJustReleased(key Key) bool`
- `KeyPressDuration(key Key) int`
- `IsMouseButtonJustPressed(button MouseButton) bool`
- `IsMouseButtonJustReleased(button MouseButton) bool`
- `MouseButtonPressDuration(button MouseButton) int`

## Dependencies

```go
require (
    github.com/hajimehoshi/ebiten/v2 v2.9.x
    github.com/spf13/cobra v1.10.x
)
```

## Testing

- Unit tests for script parser, executor, RPC handlers (`internal/script/parser_test.go`)
- End-to-end tests for RPC communication, CLI commands, and script execution (`e2e/e2e_test.go`)
- Integration tests require a running game process