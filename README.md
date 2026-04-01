# autoebiten

CLI tool for automating Ebitengine games via input injection, screenshots, scripted sequences, and custom commands.

Offers two integration methods: **Patch** (zero-code changes for existing games using Ebiten's native input) and **Library** (direct API integration for new projects).

## Installation

```bash
go install github.com/s3cy/autoebiten/cmd/autoebiten@latest
```

## Which Integration Method?

```
├─ Already using ebiten's native input functions (ebiten.IsKeyPressed, etc.)?
│  └─→ YES: Use the Patch method. No code changes required!
│
└── Writing a new game or willing to modify input handling code?
   └─→ Use the Library method.
```

**Trade-offs:**

| Aspect | Patch | Library |
|--------|-------|---------|
| Code changes | None required | Must use `autoebiten.IsKeyPressed()` etc. |
| Ebiten version | Locked to v2.9.9 (may need custom patch for other versions) | Works with any version |
| Third-party libs | Compatible with libs that use Ebiten input | Incompatible with libs using raw Ebiten input |
| Implementation | Hack - relies on Ebiten internals | Clean - direct API integration |

> ⚠️ Don't mix Patch and Library methods. If Ebiten is patched, calling autoebiten.Update(), autoebiten.Capture(), and any input related methods will panic. With the Patch method, these calls happen automatically inside Ebiten.

## Quick Start (Patch Method)

No code changes required! Patch Ebiten to add automation capabilities.

### 1. Clone and patch Ebiten

```bash
git clone https://github.com/hajimehoshi/ebiten.git /path/to/ebiten
cd /path/to/ebiten
git checkout v2.9.9  # Ensure correct version
git apply /path/to/autoebiten/ebiten.patch
go mod tidy
```

### 2. Update your game's go.mod

```go
replace github.com/hajimehoshi/ebiten/v2 => /path/to/local/ebiten
```

### 3. Build and run your game normally

```bash
go build ./cmd/my-game
./my-game
```

### 4. Control it via CLI

**Note on ticks:** 1 tick = 1 `Update()` call. Ebiten runs at 60 ticks per second (TPS) by default. This is **not** the same as FPS—your game can render at 120 FPS on a high refresh monitor while still running at 60 TPS. The default `--duration_ticks 6` means a key is held for 6 Update() calls (~100ms at 60 TPS). See [TPS vs FPS](https://github.com/tinne26/tps-vs-fps) for details.

```bash
# Press a key
autoebiten input --key KeyW --action press

# Hold a key for 6 ticks (default)
autoebiten input --key KeySpace --action hold

# Move mouse and click
autoebiten mouse --action position -x 100 -y 200
autoebiten mouse --action press --button MouseButtonLeft
autoebiten mouse --button MouseButtonLeft  # defaults to hold action
autoebiten mouse -x 100 -y 200 --button MouseButtonLeft  # move to position, then trigger the button

# Scroll wheel
autoebiten wheel -y -3

# Get injected positions (returns last injected values, not real OS positions)
autoebiten get_mouse_position
autoebiten get_wheel_position

# Take a screenshot
autoebiten screenshot --output shot.png

# Run a script file
autoebiten run --script script.json

# Run an inline script
autoebiten run --inline '{"version":"1.0","commands":[{"input":{"action":"press","key":"KeySpace"}}]}'

# Get JSON Schema for IDE support
autoebiten schema > autoebiten-schema.json

# Check connection
autoebiten ping

# Print version information
autoebiten version

# List available keys
autoebiten keys

# List mouse buttons
autoebiten mouse_buttons
```

### Patch Limitations

- **What it modifies:** The patch wraps Ebiten's public input API functions (`IsKeyPressed`, `CursorPosition`, `Wheel`, etc.) to check for injected input before falling back to real OS input. It does not modify Ebiten's internal `inputstate` package.
- **Version dependent:** The patch relies on Ebiten's internal implementation details and has only been tested with Ebiten v2.9.9. It may not apply cleanly to other versions.
- **Custom patches:** To use a different Ebiten version, you may need to write your own patch by adapting the changes in `ebiten.patch` to your target version's source code.
- **Maintenance:** You will need to re-apply the patch when updating Ebiten versions.
- **Release builds:** To ship your game to players without automation capabilities, simply remove the `replace` directive from your `go.mod` and use the official Ebiten module.

## Library Method

For new games or when you prefer explicit integration, use the library directly.

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

**Note:** When using the library integration method, you need to use `autoebiten.IsKeyPressed()` and other wrapper functions instead of Ebiten's native functions.

### 2. Build and run

```bash
# Development (default)
go build ./cmd/my-game
./my-game

# Release build (no CLI automation)
go build -tags release ./cmd/my-game
```

### 3. Control via CLI

Same CLI commands as the Patch method above.

### Build Tags

- **Default** (`go run` or `go build`): Full RPC server enabled. Your game listens for CLI commands via Unix socket.
- **Release** (`go build -tags release`): RPC server disabled. All functions are no-ops that delegate directly to ebiten. Use this when shipping your game to players.

### Input Modes

The library operates in different input modes to control how real and injected inputs are combined:

| Mode | Description |
|------|-------------|
| `InjectionOnly` | Only CLI-injected input is recognized |
| `InjectionFallback` | Per-key fallback: if a key/button has injected state, return it; otherwise check real input (default) |
| `Passthrough` | All input passes through to ebiten directly; no injection |

**Note on `InjectionFallback`:** Each input is checked independently. For example, if you inject 'KeyS' via CLI while physically holding 'KeyW' on your keyboard, `IsKeyPressed(KeyS)` returns the injected state and `IsKeyPressed(KeyW)` returns the real keyboard state. They do not interfere with each other.

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

Run from a file:

```bash
autoebiten run --script script.json
```

Or pass an inline JSON string:

```bash
autoebiten run --inline '{"version":"1.0","commands":[{"input":{"action":"press","key":"KeySpace"}}]}'
```

### JSON Schema

Generate the JSON Schema for IDE autocompletion and validation:

```bash
autoebiten schema > autoebiten-schema.json
```

The schema defines all available commands, their fields, and valid values. Use it with editors that support JSON Schema for intelligent autocomplete and inline validation.

### Script Error Handling

Scripts execute commands sequentially. **If any command fails, script execution stops immediately** and the error is returned. There is no continue-on-error or try/catch mechanism—fix the issue and re-run the script.

## Platform Support

- **macOS**: Fully supported
- **Linux**: Fully supported
- **Windows**: Not supported (uses Unix domain sockets)

## Multiple Game Instances

Each game instance uses a PID-based socket at `/tmp/autoebiten/autoebiten-{PID}.sock`.

**Auto-detection:** If `--pid` is not provided, autoebiten automatically detects a running game. If only one game is running, it uses that. If multiple games are running, it lists them and exits with an error—you must specify which one with `--pid`.

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

## License

MIT
