# autoebiten

CLI tool for automating Ebitengine games via input injection, screenshots, and scripted sequences.

## Installation

```bash
go install github.com/s3cy/autoebiten/cmd/autoebiten@latest
```

## Quick Start

### 1. Add the library to your game

```go
import "github.com/s3cy/autoebiten"

func (g *Game) Update() error {
    if !autoebiten.Update() {
        return errors.New("exit requested")
    }
    // Your update logic
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    // Your draw logic
    autoebiten.Capture(screen) // Call at the end
}
```

### 2. Run your game

```bash
./my-game
```

### 3. Control it via CLI

```bash
# Press a key
autoebiten input --key KeyW --action press

# Hold a key for 6 ticks (default)
autoebiten input --key KeySpace --action hold

# Move mouse and click
autoebiten mouse --action position --x 100 --y 200
autoebiten mouse --action press --button MouseButtonLeft

# Scroll wheel
autoebiten wheel --y -3

# Take a screenshot
autoebiten screenshot --output shot.png

# Check connection
autoebiten ping

# List available keys
autoebiten keys

# List mouse buttons
autoebiten mouse_buttons
```

## Scripted Automation

Create a JSON script for complex sequences:

```json
{
  "version": "1.0",
  "commands": [
    {"input": {"action": "press", "key": "KeyW"}},
    {"delay": {"ms": 100}},
    {"repeat": {"times": 3, "commands": [
      {"input": {"action": "press", "key": "KeyA"}},
      {"delay": {"ms": 200}}
    ]}}
  ]
}
```

Run it with:

```bash
autoebiten run --script script.json
```

## Multiple Game Instances

Each game instance uses a PID-based socket at `/tmp/autoebiten/autoebiten-{PID}.sock`.

Target a specific game:

```bash
autoebiten --pid 12345 input --key KeySpace --action press
```

Or set the socket path manually:

```bash
AUTOEBITEN_SOCKET=/tmp/autoebiten/autoebiten-12345.sock autoebiten ping
```

## Public API

```go
// In your game loop
autoebiten.Update()        // Process RPC commands, returns false on exit
autoebiten.Capture(screen) // Capture screenshot (call in Draw)

// Query input state
autoebiten.IsKeyPressed(key ebiten.Key)           // bool
autoebiten.IsMouseButtonPressed(button ebiten.MouseButton) // bool
autoebiten.CursorPosition()                       // (x, y int)
autoebiten.Wheel()                               // (x, y float64)

// Configure input mode
autoebiten.SetMode(autoebiten.InjectionFallback) // default: injected + real input
```

### Input Modes

The library operates in different input modes to control how real and injected inputs are combined:

| Mode | Description |
|------|-------------|
| `InjectionOnly` | Only CLI-injected input is recognized |
| `InjectionFallback` | Injected input takes priority; falls back to real input (default) |
| `Passthrough` | All input passes through to ebiten directly; no injection |

### Build Tags

autoebiten uses build tags to control which implementation is included:

- **Default** (`go run` or `go build`): Full RPC server enabled. Your game listens for CLI commands via Unix socket.
- **Release** (`go build -tags release`): RPC server disabled. All functions are no-ops that delegate directly to ebiten. Use this when shipping your game to players.

```bash
# Development (default)
go build ./cmd/autoebiten

# Release build (no CLI automation)
go build -tags release ./cmd/autoebiten
```

## License

MIT