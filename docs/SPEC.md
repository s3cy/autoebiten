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
│  (input, mouse, wheel, screenshot, run, ping, custom)       │
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
│  │  autoebiten.Register() → custom command handlers    │    │
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

### Deep Integration via Ebiten Patch

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
autoebiten mouse --action press -x 100 -y 200 --button MouseButtonLeft
autoebiten wheel -x 0 -y -3
autoebiten screenshot --output shot.png
autoebiten run --script script.json
autoebiten ping
autoebiten keys
autoebiten mouse_buttons
autoebiten get_mouse_position
autoebiten get_wheel_position
autoebiten list_custom
autoebiten custom getPlayerInfo
autoebiten state --name gamestate --path Player.Health
autoebiten wait-for --condition "state:gamestate:Player.Health == 100" --timeout 10s
autoebiten version
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
| `version` | Get CLI and game library versions |
| `exit` | Request game to exit |
| `get_mouse_position` | Get injected mouse cursor position |
| `get_wheel_position` | Get injected wheel position |
| `list_custom_commands` | List registered custom commands |
| `custom` | Execute a custom command |

**Note:** The CLI `state` and `wait-for` commands use the `custom` RPC method internally. The `state` command queries registered state exporters (via `.state.<name>` prefix), and `wait-for` polls using the `custom` method until a condition is met.

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
| `mouse` | `action` (position/press/release/hold), `x`, `y`, `button`, `duration_ticks` | Mouse input. Defaults to `hold` when `button` is provided without `action`. Inject "-x 0 -y 0" to restore real inputs |
| `wheel` | `x`, `y` | Mouse wheel movement. Inject "-x 0 -y 0" to restore real inputs |
| `screenshot` | `output` | Capture screenshot |
| `delay` | `ms` | Wait in milliseconds |
| `repeat` | `times`, `commands` | Repeat block N times |

### Input Modes

- `InjectionOnly` - Returns only injected input results
- `InjectionFallback` (default) - Returns injected results if available, otherwise falls back to ebiten's native input
- `Passthrough` - Only uses ebiten's native input handling

### Input Function Support

The library provides mode-aware wrappers for both direct input queries and inpututil functions:

**Direct Input:**
- `IsKeyPressed(key Key) bool`
- `CursorPosition() (x, y int)` - Returns injected position or falls back to real position based on mode
- `Wheel() (xoff, yoff float64)` - Returns injected offset or falls back to real offset based on mode
- `IsMouseButtonPressed(mouseButton MouseButton) bool`

**Note on Position Queries:** The `get_mouse_position` and `get_wheel_position` CLI commands (and their RPC equivalents) only return the **injected** values that were previously set via the `mouse` or `wheel` commands. They do not return the real OS cursor/wheel position. To retrieve real positions, use `ebiten.CursorPosition()` and `ebiten.Wheel()` directly in your game code when in `Passthrough` or `InjectionFallback` mode.

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

## Testkit Package

The `testkit` package provides a Go testing framework for Ebiten games that use autoebiten.

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      testkit Package                         │
├───────────────────────────┬─────────────────────────────────┤
│      Black-Box (Game)     │       White-Box (Mock)          │
│  ┌─────────────────────┐  │  ┌─────────────────────────┐    │
│  │ Separate Process    │  │  │ Same Process            │    │
│  │ (game binary)       │  │  │                         │    │
│  │         │           │  │  │ ┌─────────────────┐     │    │
│  │         ▼           │  │  │ │ Game struct     │     │    │
│  │ ┌─────────────┐     │  │  │ │ (user code)     │     │    │
│  │ │ RPC Server  │     │  │  │ └────────┬────────┘     │    │
│  │ │ (in game)   │◄────┼──┼──┤          │              │    │
│  │ └─────────────┘     │  │  │          ▼              │    │
│  │         ▲           │  │  │ ┌─────────────────┐     │    │
│  │         │           │  │  │ │ Mock RPC Client │     │    │
│  │ ┌─────────────┐     │  │  │ │ (simulates CLI) │     │    │
│  │ │ RPC Client  │     │  │  │ └─────────────────┘     │    │
│  │ │ (in test)   │     │  │  │                         │    │
│  │ └─────────────┘     │  │  └─────────────────────────┘    │
│  └─────────────────────┘  │                                 │
└───────────────────────────┴─────────────────────────────────┘
```

### Black-Box Mode (Game)

Launches a game binary in a separate process and controls it via RPC.

```go
func Launch(t *testing.T, binaryPath string, opts ...Option) *Game
```

Key capabilities:
- Input injection (`PressKey`, `HoldKey`, `MoveMouse`, etc.)
- Screenshots (`Screenshot`, `ScreenshotToFile`, `ScreenshotBase64`)
- State queries via autoebiten `RegisterStateExporter` (`StateQuery`)
- Custom command execution (`RunCustom`)
- Lifecycle management (`Shutdown`, `Ping`, `WaitFor`)

### White-Box Mode (Mock)

Simulates autoebiten's RPC server for testing game logic in-process.

```go
func NewMock(t *testing.T, game GameUpdate) *Mock
```

Key capabilities:
- Input injection (`InjectKeyPress`, `InjectMousePosition`, etc.)
- Tick execution (`Tick`, `Ticks`)

### State Exporter

Games can expose internal state for black-box testing using `autoebiten.RegisterStateExporter`:

```go
// In your game
autoebiten.RegisterStateExporter("gamestate", &gameState)
```

Tests can then query state via dot-notation paths:

```go
value, err := game.StateQuery("gamestate", "Player.X")
```

Supports paths like `Player.X`, `Inventory.0.Name`, `Skills.Sword`.

### Documentation

See the [testkit package documentation](testkit/doc.go) for complete API reference and usage examples.