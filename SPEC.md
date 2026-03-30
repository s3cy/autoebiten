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
├── inpututil/                   # Input utilities (build-tagged)
├── examples/
│   └── simple/
│       └── main.go              # Example game
├── e2e/
│   └── e2e_test.go              # End-to-end tests
├── autoebiten.go                # Mode configuration
├── autoebiten_default.go        # Default build (with RPC server)
├── autoebiten_release.go        # Release build (no-op stubs)
├── go.mod
└── SPEC.md
```

## Public API (`autoebiten` package)

```go
// Update processes RPC commands from the socket.
// Call this in your game loop's Update function.
// Returns false if the game should exit (exit command received).
func Update() bool

// Capture captures screenshots.
// Call this at the end of your game loop's Draw function.
func Capture(screen *ebiten.Image)

// IsKeyPressed returns whether the key is pressed (real + injected).
func IsKeyPressed(key ebiten.Key) bool

// IsMouseButtonPressed returns whether the mouse button is pressed.
func IsMouseButtonPressed(button ebiten.MouseButton) bool

// CursorPosition returns the current cursor position (real + injected).
func CursorPosition() (x, y int)

// Wheel returns the wheel state.
func Wheel() (x, y float64)

// SetMode sets the input handling mode (InjectionOnly, InjectionFallback, Passthrough).
func SetMode(mode Mode)

// GetMode returns the current input handling mode.
func GetMode() Mode
```

### Input Modes

- `InjectionOnly` - Returns only injected input results
- `InjectionFallback` (default) - Returns injected results if available, otherwise falls back to ebiten's native input
- `Passthrough` - Only uses ebiten's native input handling

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