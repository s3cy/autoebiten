# autoebiten

CLI tool for automating Ebitengine games via input injection, screenshots, scripted sequences, and custom commands.

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
autoebiten mouse --action position -x 100 -y 200
autoebiten mouse --action press --button MouseButtonLeft
autoebiten mouse --button MouseButtonLeft  # defaults to hold action
autoebiten mouse -x 100 -y 200 --button MouseButtonLeft  # move to positen, then trigger the button

# Scroll wheel
autoebiten wheel -y -3

# Get injected positions (returns last injected values, not real OS positions)
autoebiten get_mouse_position
autoebiten get_wheel_position

# Take a screenshot
autoebiten screenshot --output shot.png

# Check connection
autoebiten ping

# List available keys
autoebiten keys

# List mouse buttons
autoebiten mouse_buttons
```

## Custom Commands

Games can register custom commands that can be invoked from the CLI:

### Game Side

```go
// In your game initialization
autoebiten.Register("getPlayerInfo", func(ctx autoebiten.CommandContext) {
    info := fmt.Sprintf("Health: %d, Mana: %d", playerHealth, playerMana)
    ctx.Respond(info)
})

autoebiten.Register("heal", func(ctx autoebiten.CommandContext) {
    playerHealth = min(playerHealth+20, 100)
    ctx.Respond(fmt.Sprintf("Healed to %d", playerHealth))
})
```

The `CommandContext` provides:
- `Request() string` - The request data sent from CLI
- `Respond(response string)` - Send response back to CLI (can be called immediately or deferred)

### CLI Side

```bash
# List available custom commands
autoebiten list_custom

# Execute a custom command
autoebiten custom getPlayerInfo

# Execute with request data
autoebiten custom echo --request "hello world"
```

See [examples/custom_commands](examples/custom_commands/main.go) for a complete example.

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

// inpututil wrappers (respect input mode)
autoebiten.IsKeyJustPressed(key ebiten.Key)       // bool
autoebiten.IsKeyJustReleased(key ebiten.Key)      // bool
autoebiten.KeyPressDuration(key ebiten.Key)       // int (ticks)
autoebiten.IsMouseButtonJustPressed(button ebiten.MouseButton)  // bool
autoebiten.IsMouseButtonJustReleased(button ebiten.MouseButton) // bool
autoebiten.MouseButtonPressDuration(button ebiten.MouseButton)  // int (ticks)

// Configure input mode
autoebiten.SetMode(autoebiten.InjectionFallback) // default: injected + real input

// Custom commands
autoebiten.Register(name string, handler func(CommandContext)) // Register a custom command
autoebiten.Unregister(name string) bool                        // Remove a custom command
autoebiten.ListCustomCommands() []string                       // List registered commands
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

**Note:** When using the library integration method, you need to use `autoebiten.IsKeyPressed()` and other wrapper functions instead of Ebiten's native functions. This may require changes if your game already uses Ebiten's input functions directly.

## Deeper Integration (Patching Ebiten)

For games that use Ebiten's native input functions directly without modification, autoebiten can be integrated at the Ebiten engine level via a patch. This approach requires no code changes to your game.

### Setup Instructions

1. Clone the Ebiten repository to a local directory:
```bash
git clone https://github.com/hajimehoshi/ebiten.git /path/to/ebiten
cd /path/to/ebiten
git checkout v2.9.9  # Ensure correct version
```

2. Apply the autoebiten patch:
```bash
git apply /path/to/autoebiten/ebiten.patch
```

3. In your game's `go.mod`, add a replace directive to use the local Ebiten:
```go
replace github.com/hajimehoshi/ebiten/v2 => /path/to/ebiten
```

4. Build and run your game normally - no code changes required!

### Limitations

- **Version dependent:** The patch relies on Ebiten's internal implementation details and has only been tested with Ebiten v2.9.9. It may not apply cleanly to other versions.
- **Maintenance:** You will need to re-apply the patch when updating Ebiten versions.
- **Compatibility:** Patch may require adjustment for significant Ebiten version changes.

## License

MIT