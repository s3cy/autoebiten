# testkit Design Document

**Date:** 2026-04-05  
**Status:** Approved  
**Author:** Claude Code

## Overview

The `testkit` package provides a Go testing framework for Ebiten games that use autoebiten. It enables two testing modes:

1. **Black-Box Mode (`Game`)** — Launches game in separate process, controls via RPC
2. **White-Box Mode (`Mock`)** — Tests game logic in same process with mocked inputs

## Goals

- Make e2e testing as easy as writing unit tests
- Support both integration testing (full game) and unit testing (game logic)
- Enable screenshot-based assertions
- Allow inspection of game internal state
- Follow Go testing idioms (use `testing.T`, `defer cleanup()`)

## Non-Goals

- Replace the CLI (complementary to `autoebiten` commands)
- Support non-Ebiten games
- Provide visual debugging tools

## Architecture

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

## API Reference

### Black-Box Mode — `Game`

Controls a running game process via RPC.

#### Construction

```go
// Launch starts a game binary and waits for it to be ready.
func Launch(t *testing.T, binary string, opts ...Option) (*Game, error)

// Options
type Option func(*config)

func WithArgs(args ...string) Option           // Pass args to game binary
func WithEnv(key, value string) Option          // Set environment variables
func WithTimeout(d time.Duration) Option        // Startup timeout (default: 10s)
func WithWorkingDir(dir string) Option          // Working directory
```

#### Lifecycle

```go
// Shutdown terminates the game process.
func (g *Game) Shutdown() error

// Ping checks if game is responsive.
func (g *Game) Ping() error

// WaitFor blocks until condition is met or timeout.
func (g *Game) WaitFor(ctx context.Context, fn func() (bool, error)) error
```

#### Input

```go
func (g *Game) PressKey(key input.Key) error
func (g *Game) ReleaseKey(key input.Key) error
func (g *Game) HoldKey(key input.Key, durationTicks int64) error
func (g *Game) MoveMouse(x, y int) error
func (g *Game) PressMouseButton(btn input.MouseButton) error
func (g *Game) ReleaseMouseButton(btn input.MouseButton) error
func (g *Game) ScrollWheel(x, y float64) error
```

#### Screenshots

```go
func (g *Game) Screenshot() (*image.RGBA, error)
func (g *Game) ScreenshotToFile(path string) error
func (g *Game) ScreenshotBase64() (string, error)
```

#### Custom Commands

```go
func (g *Game) RunCustom(name, request string) (string, error)
func (g *Game) ListCustom() ([]string, error)

// Query exported state (requires StateExporter registration)
func (g *Game) StateQuery(path string) (any, error)
```

### White-Box Mode — `Mock`

Simulates autoebiten's RPC server for testing game logic in-process.

#### Construction

```go
func NewMock(t *testing.T) *Mock
```

#### Input Injection

```go
func (m *Mock) InjectKeyPress(key input.Key)
func (m *Mock) InjectKeyRelease(key input.Key)
func (m *Mock) InjectMousePosition(x, y int)
func (m *Mock) InjectMouseButton(btn input.MouseButton, pressed bool)
func (m *Mock) InjectWheel(x, y float64)
```

#### Custom Commands

```go
func (m *Mock) RegisterCustom(name string, handler CustomHandler)
```

#### Tick Execution

```go
// Tick runs one game update tick.
func (m *Mock) Tick(g interface{ Update() error }) error

// Ticks runs multiple ticks.
func (m *Mock) Ticks(g interface{ Update() error }, n int) error
```

### State Exporter

Reflection-based state exposure for black-box testing.

```go
// StateExporter creates a custom command handler that exports game state.
func StateExporter(v interface{}) func(autoebiten.CommandContext)
```

**Reflection rules:**
- Only exported fields (uppercase first letter)
- Primitives returned as-is
- Structs navigable via dot notation
- Slices/arrays: indexable (`Player.Inventory.0.Name`)
- Maps: keyable if string keys (`Player.Skills.Sword`)

## Usage Examples

### Black-Box Test

```go
func TestPlayerMovement(t *testing.T) {
    game := testkit.Launch(t, "./mygame")
    defer game.Shutdown()
    
    // Initial position
    x, _ := game.StateQuery("Player.X")
    assert.Equal(t, 0, x)
    
    // Press D for 10 ticks
    game.HoldKey(ebiten.KeyD, 10)
    
    // Verify movement
    x, _ = game.StateQuery("Player.X")
    assert.Equal(t, 10, x)
}
```

### White-Box Test

```go
func TestPlayerTakesDamage(t *testing.T) {
    g := NewGame()
    mock := testkit.NewMock(t)
    
    // Simulate enemy attack
    g.EnemyAttacks()
    
    // Run one tick
    mock.Tick(g)
    
    // Direct state access
    assert.Equal(t, 90, g.Player.Health)
}
```

### State Exporter Registration

```go
// In game code
func main() {
    g := &Game{}
    
    // Expose state for testing
    autoebiten.Register("testkit.state", testkit.StateExporter(g))
    
    ebiten.RunGame(g)
}
```

## Error Handling

| Scenario | Behavior |
|----------|----------|
| Game fails to start | `Launch` returns error with process stderr |
| Game crashes during test | Methods return error; `Shutdown` is no-op |
| RPC timeout | Returns `context.DeadlineExceeded` (default 5s) |
| State query path not found | Returns `ErrPathNotFound` |
| Screenshot failure | Returns wrapped error from RPC |

## Testing the Testkit

Example games in `testkit/internal/testgames/`:

- `simple/` — Minimal game for basic tests
- `stateful/` — Game with StateExporter for state tests
- `custom/` — Game with custom commands

## Implementation Notes

### Package Structure

```
testkit/
├── doc.go              // Package documentation
├── game.go             // Black-box Game implementation
├── game_test.go        // Tests using testgames
├── mock.go             // White-box Mock implementation
├── mock_test.go
├── state.go            // StateExporter implementation
├── state_test.go
└── internal/
    └── testgames/
        ├── simple/
        │   └── main.go
        ├── stateful/
        │   └── main.go
        └── custom/
            └── main.go
```

### Dependencies

- `github.com/s3cy/autoebiten/internal/rpc` — RPC client
- `github.com/s3cy/autoebiten/internal/input` — Input types
- Standard library: `os/exec`, `context`, `testing`, `image`

## Future Considerations

- Parallel test execution (games use unique socket paths)
- Screenshot comparison helpers (diff, pixel assertions)
- Record/replay functionality for regression tests
- CI/CD integration helpers
