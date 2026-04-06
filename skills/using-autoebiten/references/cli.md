# CLI Reference

> Purpose: Complete guide to controlling Ebitengine games via the autoebiten CLI
> Audience: Developers using CLI for automation, testing, or AI-driven gameplay

---

## Quick Decision

**Integration method:**
```
├─ Already using ebiten.IsKeyPressed() in your game?
│  └─→ YES: Use Patch method (no code changes)
│
└─ Writing a new game or willing to modify input handling?
   └─→ Use Library method (direct API)
```

**Input mode (Library integration):**
```
├─ Want only CLI-injected inputs?
│  └─→ InjectionOnly (SetMode(autoebiten.InjectionOnly))
│
├─ Want CLI + real keyboard/mouse combined?
│  └─→ InjectionFallback (default)
│
└─ Disable CLI control temporarily?
   └─→ Passthrough (SetMode(autoebiten.Passthrough))
```

---

## Overview

autoebiten is a CLI tool that enables automation of Ebitengine games through:

- **Input injection** - Send keyboard, mouse, and wheel inputs programmatically
- **Screenshots** - Capture game state visually
- **Script execution** - Run complex automation sequences
- **Custom commands** - Execute game-specific actions
- **State queries** - Read internal game state for testing

Communication uses JSON-RPC over Unix sockets. Each game instance creates a socket at `/tmp/autoebiten/autoebiten-{PID}.sock`.

---

## Integration Methods

### Patch Method (Zero Code Changes)

For existing games using Ebiten's native input functions (`ebiten.IsKeyPressed`, etc.).

**Requirements:**
- Ebiten v2.9.9 (patch tested with this version)
- Local Ebiten clone with patch applied

**Steps:**
1. Clone Ebiten: `git clone https://github.com/hajimehoshi/ebiten.git`
2. Checkout v2.9.9: `git checkout v2.9.9`
3. Apply patch: `git apply /path/to/autoebiten/ebiten.patch`
4. Add replace directive to game's go.mod

### Library Method (Direct API)

For new games or when you control the source code.

**Requirements:**
- Import `github.com/s3cy/autoebiten`
- Replace ebiten input calls with autoebiten wrappers

**Key functions:**
- `autoebiten.Update()` - Process RPC commands (call in Update)
- `autoebiten.Capture(screen)` - Enable screenshots (call in Draw)
- `autoebiten.IsKeyPressed(key)` - Input query wrapper

**Build modes:**
- Default: RPC server enabled (automation available)
- `-tags release`: RPC disabled (for shipping to players)

---

## API Reference (Library Method)

### Core Functions

#### Update

**Signature:** `func Update() bool`

**Purpose:** Process pending RPC commands and return whether game should continue.

**Parameters:** None

**Returns:**
- `true` - Continue running
- `false` - Exit requested (return error from Update)

**Example:**
```go
func (g *Game) Update() error {
    if !autoebiten.Update() {
        return fmt.Errorf("exit requested")
    }
    // Your update logic
    return nil
}
```

**Notes:** Must be called every tick. In Patch method, this happens automatically.

---

#### Capture

**Signature:** `func Capture(screen image.Image)`

**Purpose:** Enable screenshot capture for CLI requests.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| screen | image.Image | The screen image from Draw |

**Returns:** None

**Example:**
```go
func (g *Game) Draw(screen *ebiten.Image) {
    // Your draw logic
    autoebiten.Capture(screen) // Call at end
}
```

**Notes:** Call at the end of Draw when screen is complete.

---

#### SetMode

**Signature:** `func SetMode(mode Mode)`

**Purpose:** Configure how injected and real inputs are combined.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| mode | Mode | One of: InjectionOnly, InjectionFallback, Passthrough |

**Returns:** None

**Example:**
```go
autoebiten.SetMode(autoebiten.InjectionOnly)
```

**Notes:** Default is InjectionFallback.

---

#### GetMode

**Signature:** `func GetMode() Mode`

**Purpose:** Get current input mode.

**Returns:**
- Current Mode value

---

### Input Query Functions

#### IsKeyPressed

**Signature:** `func IsKeyPressed(key ebiten.Key) bool`

**Purpose:** Check if key is pressed (injected or real based on mode).

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| key | ebiten.Key | Key constant (KeyA, KeySpace, KeyArrowUp, etc.) |

**Returns:**
- `true` if key is pressed

**Example:**
```go
if autoebiten.IsKeyPressed(ebiten.KeySpace) {
    player.Jump()
}
```

**Notes:** In Patch method, use `ebiten.IsKeyPressed()` directly (patched).

---

#### IsMouseButtonPressed

**Signature:** `func IsMouseButtonPressed(button ebiten.MouseButton) bool`

**Purpose:** Check if mouse button is pressed.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| button | ebiten.MouseButton | MouseButtonLeft, MouseButtonRight, MouseButtonMiddle |

**Returns:**
- `true` if button is pressed

---

#### CursorPosition

**Signature:** `func CursorPosition() (x, y int)`

**Purpose:** Get mouse cursor position (injected or real based on mode).

**Returns:**
- x, y - Cursor coordinates in pixels

**Example:**
```go
x, y := autoebiten.CursorPosition()
```

---

#### Wheel

**Signature:** `func Wheel() (x, y float64)`

**Purpose:** Get mouse wheel scroll delta.

**Returns:**
- x - Horizontal scroll (negative=left, positive=right)
- y - Vertical scroll (negative=down, positive=up)

---

### inpututil Wrappers

#### IsKeyJustPressed

**Signature:** `func IsKeyJustPressed(key ebiten.Key) bool`

**Purpose:** Check if key was just pressed this tick.

---

#### IsKeyJustReleased

**Signature:** `func IsKeyJustReleased(key ebiten.Key) bool`

**Purpose:** Check if key was just released this tick.

---

#### KeyPressDuration

**Signature:** `func KeyPressDuration(key ebiten.Key) int`

**Purpose:** Get how many ticks a key has been held.

**Returns:**
- Duration in ticks (0 if not pressed)

---

#### IsMouseButtonJustPressed

**Signature:** `func IsMouseButtonJustPressed(button ebiten.MouseButton) bool`

**Purpose:** Check if mouse button was just pressed this tick.

---

#### IsMouseButtonJustReleased

**Signature:** `func IsMouseButtonJustReleased(button ebiten.MouseButton) bool`

**Purpose:** Check if mouse button was just released this tick.

---

#### MouseButtonPressDuration

**Signature:** `func MouseButtonPressDuration(button ebiten.MouseButton) int`

**Purpose:** Get how many ticks a mouse button has been held.

---

### Custom Commands

#### Register

**Signature:** `func Register(name string, handler func(CommandContext))`

**Purpose:** Register a custom command callable from CLI.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| name | string | Command name (e.g., "getPlayerInfo") |
| handler | func(CommandContext) | Handler function |

**Example:**
```go
autoebiten.Register("heal", func(ctx autoebiten.CommandContext) {
    player.Health = min(player.Health+20, 100)
    ctx.Respond(fmt.Sprintf("Healed to %d", player.Health))
})
```

---

#### Unregister

**Signature:** `func Unregister(name string) bool`

**Purpose:** Remove a registered custom command.

**Returns:**
- `true` if command was found and removed

---

#### ListCustomCommands

**Signature:** `func ListCustomCommands() []string`

**Purpose:** Get names of all registered custom commands.

---

#### CommandContext

**Methods:**
- `Request() string` - Get request data from CLI
- `Respond(response string)` - Send response back to CLI

---

#### RegisterStateExporter

**Signature:** `func RegisterStateExporter(name string, root any)`

**Purpose:** Expose game state for reflection-based queries.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| name | string | Exporter name (e.g., "gamestate") |
| root | any | Pointer to game state struct |

**Example:**
```go
type GameState struct {
    Player  Player
    Enemies []Enemy
}

autoebiten.RegisterStateExporter("gamestate", &gameInstance)
```

**Notes:** Query via CLI: `autoebiten state --name gamestate --path Player.Health`

**Supported paths:**
- `Player.X` - Struct field
- `Inventory.0.Name` - Array/slice index
- `Skills.Sword` - Map key

---

## CLI Commands

### Global Flags

These flags apply to all commands:

| Flag | Description |
|------|-------------|
| --pid, -p | Target game process PID (auto-detected if not specified) |

If multiple games are running and --pid is not specified, autoebiten will list available games and exit with an error.

### Input Control

#### input

Send keyboard input to the game.

```bash
autoebiten input --key <KeyName> --action <Action> [--duration_ticks N]
```

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| --key, -k | required | Key name (KeyA, KeySpace, KeyArrowUp) |
| --action, -a | hold | Action: press, release, hold |
| --duration_ticks, -d | 6 | Ticks to hold (for hold action) |
| --async | false | Return immediately without waiting |
| --no-record | false | Skip recording this command |

**Actions:**
- `press` - Press and immediately release
- `release` - Release a held key
- `hold` - Press and hold for duration_ticks

**Examples:**
```bash
# Press space once
autoebiten input --key KeySpace --action press

# Hold W for 10 ticks (~167ms at 60 TPS)
autoebiten input --key KeyW --action hold --duration_ticks 10
```

---

#### mouse

Send mouse input to the game.

```bash
autoebiten mouse --action <Action> [-x N] [-y N] [--button <ButtonName>]
```

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| --action, -a | position | Action: position, press, release, hold |
| --x, -x | 0 | X coordinate |
| --y, -y | 0 | Y coordinate |
| --button, -b | none | Button: MouseButtonLeft, MouseButtonRight, MouseButtonMiddle |
| --duration_ticks, -d | 6 | Ticks to hold |
| --async | false | Return immediately |
| --no-record | false | Skip recording |

**Actions:**
- `position` - Move cursor to (x, y)
- `press` - Press button at current position
- `release` - Release button
- `hold` - Press and hold for duration_ticks (default when --button is set)

**Examples:**
```bash
# Move cursor
autoebiten mouse -x 100 -y 200

# Click
autoebiten mouse --button MouseButtonLeft

# Click at position
autoebiten mouse -x 100 -y 200 --button MouseButtonLeft

# Restore real mouse input
autoebiten mouse -x 0 -y 0
```

---

#### wheel

Send scroll wheel input.

```bash
autoebiten wheel -x <X> -y <Y>
```

**Flags:**
| Flag | Description |
|------|-------------|
| --x, -x | Horizontal scroll (negative=left, positive=right) |
| --y, -y | Vertical scroll (negative=down, positive=up) |
| --async | Return immediately |
| --no-record | Skip recording |

**Examples:**
```bash
# Scroll up 3 units
autoebiten wheel -y -3

# Restore real wheel input
autoebiten wheel -x 0 -y 0
```

---

### Screenshot

#### screenshot

Capture the game window.

```bash
autoebiten screenshot [--output <Path>] [--base64]
```

**Flags:**
| Flag | Description |
|------|-------------|
| --output, -o | Output file path (auto-generated if not set) |
| --base64 | Output as base64 string instead of file |
| --async, -a | Return immediately |
| --no-record | Skip recording |

**Examples:**
```bash
autoebiten screenshot
autoebiten screenshot --output capture.png
autoebiten screenshot --base64
```

---

### Script Execution

#### run

Execute a JSON script.

```bash
autoebiten run --script <Path>
autoebiten run --inline '<JSON>'
```

**Flags:**
| Flag | Description |
|------|-------------|
| --script, -s | Path to script file |
| --inline | Inline JSON string |

**Examples:**
```bash
# From file
autoebiten run --script automation.json

# Inline
autoebiten run --inline '{"version":"1.0","commands":[{"input":{"key":"KeySpace"}}]}'
```

---

### Status and Info

#### ping

Check if game is running and responsive.

```bash
autoebiten ping
```

---

#### version

Print CLI and game library versions.

```bash
autoebiten version
```

---

#### keys

List all available key names.

```bash
autoebiten keys
```

---

#### mouse_buttons

List all available mouse button names.

```bash
autoebiten mouse_buttons
```

---

#### get_mouse_position

Get injected mouse cursor position.

```bash
autoebiten get_mouse_position
```

---

#### get_wheel_position

Get injected wheel position.

```bash
autoebiten get_wheel_position
```

---

#### schema

Output JSON Schema for script files.

```bash
autoebiten schema > autoebiten-schema.json
```

---

### Custom Commands

#### list_custom

List registered custom commands.

```bash
autoebiten list_custom
```

---

#### custom

Execute a custom command.

```bash
autoebiten custom <Name> [--request <Data>]
```

**Flags:**
| Flag | Description |
|------|-------------|
| --name, -n | Command name |
| --request, -r | Request data to pass |
| --no-record | Skip recording |

**Examples:**
```bash
autoebiten custom getPlayerInfo
autoebiten custom echo --request "hello"
```

---

### State Queries

#### state

Query game state via registered exporter.

```bash
autoebiten state --name <ExporterName> --path <Path>
```

**Flags:**
| Flag | Description |
|------|-------------|
| --name | State exporter name (required) |
| --path | Dot-notation path (required) |
| --no-record | Skip recording |

**Examples:**
```bash
autoebiten state --name gamestate --path Player.Health
autoebiten state --name gamestate --path Inventory.0.Name
```

---

#### wait-for

Wait for a condition to be met.

```bash
autoebiten wait-for --condition "<Condition>" --timeout <Duration>
```

**Flags:**
| Flag | Description |
|------|-------------|
| --condition | Condition expression (required) |
| --timeout | Maximum wait duration (required, e.g., 10s, 5m) |
| --interval | Poll interval (default 100ms) |
| --verbose, -v | Print errors during polling |
| --no-record | Skip recording |

**Condition format:**
```
<type>:<name>:<path> <operator> <value>
```

- type: `state` or `custom`
- name: exporter name or custom command name
- path: dot-notation path or request string
- operator: ==, !=, <, >, <=, >=
- value: JSON value

**Examples:**
```bash
autoebiten wait-for --condition "state:gamestate:Player.Health == 100" --timeout 10s
```

---

### Recording

#### clear_recording

Clear the recording file for current game.

```bash
autoebiten clear_recording
```

---

#### replay

Replay recorded commands.

```bash
autoebiten replay [--speed N] [--dump <Path>]
```

**Flags:**
| Flag | Description |
|------|-------------|
| --speed, -s | Speed multiplier (default 1.0) |
| --dump, -d | Dump script to file instead of executing |

**Examples:**
```bash
# Replay at normal speed
autoebiten replay

# Replay at 2x speed
autoebiten replay --speed 2

# Dump script without executing
autoebiten replay --dump script.json
```

---

## Tutorial

### Step 1: Choose Integration Method

**Decision tree:**
1. Check your game's go.mod for `replace github.com/hajimehoshi/ebiten/v2`
2. If found → Patch method already applied
3. If not found → Check for `autoebiten.Update()` in code
4. If found → Library method already integrated
5. If neither → Choose based on needs

**Patch method checklist:**
- [ ] Clone Ebiten locally
- [ ] Checkout v2.9.9
- [ ] Apply `ebiten.patch`
- [ ] Add replace directive to go.mod

**Library method checklist:**
- [ ] Import autoebiten package
- [ ] Add `autoebiten.Update()` to Update()
- [ ] Add `autoebiten.Capture()` to Draw()
- [ ] Replace input calls with wrappers

---

### Step 2: Patch Method Setup

**Goal:** Enable automation without modifying game code.

**Prerequisites:**
- Ebiten v2.9.9 compatible game
- Git installed

**Actions:**

```bash
# Clone Ebiten
git clone https://github.com/hajimehoshi/ebiten.git /path/to/ebiten
cd /path/to/ebiten
git checkout v2.9.9

# Apply patch (from autoebiten repo root)
git apply /path/to/autoebiten/ebiten.patch
go mod tidy
```

**Add to game's go.mod:**

```go
replace github.com/hajimehoshi/ebiten/v2 => /path/to/local/ebiten
```

**Build and run:**

```bash
go build ./cmd/mygame
./mygame
```

**Expected:** Game runs normally. Automation socket created at startup.

**Troubleshooting:**
- Patch doesn't apply: Check Ebiten version matches v2.9.9
- Import errors: Run `go mod tidy` in both repos

---

### Step 3: Library Method Setup

**Goal:** Integrate autoebiten directly into game code.

**Prerequisites:**
- Control over game source code

**Actions:**

Add import:

```go
import "github.com/s3cy/autoebiten"
```

Modify Update:

```go
func (g *Game) Update() error {
    if !autoebiten.Update() {
        return fmt.Errorf("exit requested")
    }
    // Your logic here
    return nil
}
```

Modify Draw:

```go
func (g *Game) Draw(screen *ebiten.Image) {
    // Your drawing here
    autoebiten.Capture(screen) // Call at end
}
```

Replace input calls:

```go
// Before
if ebiten.IsKeyPressed(ebiten.KeySpace) { ... }

// After (Library method)
if autoebiten.IsKeyPressed(ebiten.KeySpace) { ... }
```

**Build modes:**

```bash
# Development (automation enabled)
go build ./cmd/mygame

# Release (no automation)
go build -tags release ./cmd/mygame
```

---

### Step 4: Verify Connection

**Goal:** Confirm CLI can communicate with game.

**Actions:**

```bash
# Start game
./mygame &

# Check connection
autoebiten ping
```

**Expected:** Output: `game is running`

**Troubleshooting:**
- "connection failed": Game not started or socket not created
- Multiple games running: Use `--pid` to specify

---

### Step 5: Basic CLI Commands

**Goal:** Send simple inputs and capture screenshot.

**Actions:**

```bash
# Press a key
autoebiten input --key KeySpace --action press

# Move mouse
autoebiten mouse -x 100 -y 200

# Click
autoebiten mouse --button MouseButtonLeft

# Take screenshot
autoebiten screenshot --output test.png
```

---

### Step 6: Understanding Ticks

**Goal:** Use duration_ticks correctly.

**Key concept:** 1 tick = 1 `Update()` call. Ebiten runs at 60 TPS by default.

**Duration calculation:**
- 6 ticks ≈ 100ms (default)
- 10 ticks ≈ 167ms
- 60 ticks ≈ 1 second

**Example:**

```bash
# Hold key for 1 second
autoebiten input --key KeyW --action hold --duration_ticks 60
```

**Note:** TPS ≠ FPS. Game can render at 120 FPS while running at 60 TPS.

---

### Step 7: Custom Commands Intro

**Goal:** Add game-specific commands.

**In game code:**

```go
autoebiten.Register("getPlayerInfo", func(ctx autoebiten.CommandContext) {
    ctx.Respond(fmt.Sprintf("Health: %d", player.Health))
})
```

**From CLI:**

```bash
autoebiten custom getPlayerInfo
```

**Expected:** Output: `Health: 100`

---

## Examples

### Library Integration

**Scenario:** New game with direct autoebiten integration.

**Code:**
```go
package main

import (
    "fmt"
    "image/color"
    "log"

    "github.com/s3cy/autoebiten"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
    screenWidth  = 640
    screenHeight = 480
)

type Game struct{}

func (g *Game) Update() error {
    // Process CLI commands
    if !autoebiten.Update() {
        return fmt.Errorf("exit requested")
    }

    // Check injected or real input
    if autoebiten.IsKeyPressed(ebiten.KeySpace) {
        fmt.Println("Space pressed!")
    }

    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    screen.Fill(color.RGBA{0x00, 0x00, 0x66, 0xff})
    ebitenutil.DebugPrint(screen, "autoebiten demo")
    autoebiten.Capture(screen) // Enable screenshots
}

func (g *Game) Layout(_, _ int) (int, int) {
    return screenWidth, screenHeight
}

func main() {
    ebiten.SetWindowSize(screenWidth, screenHeight)
    ebiten.SetWindowTitle("Demo")
    if err := ebiten.RunGame(&Game{}); err != nil {
        log.Fatal(err)
    }
}
```

**How to run:**
```bash
go build -o demo
./demo &
autoebiten input --key KeySpace --action press
```

---

### Custom Commands

**Scenario:** Add game-specific actions callable from CLI.

**Code:**
```go
// In game initialization
func NewGame() *Game {
    g := &Game{
        PlayerHealth: 100,
    }

    // Register commands
    autoebiten.Register("heal", func(ctx autoebiten.CommandContext) {
        g.PlayerHealth = min(g.PlayerHealth+20, 100)
        ctx.Respond(fmt.Sprintf("Healed to %d", g.PlayerHealth))
    })

    autoebiten.Register("damage", func(ctx autoebiten.CommandContext) {
        g.PlayerHealth = max(g.PlayerHealth-10, 0)
        ctx.Respond(fmt.Sprintf("Health: %d", g.PlayerHealth))
    })

    return g
}
```

**CLI usage:**
```bash
autoebiten custom heal
autoebiten custom damage
```

---

### State Exporter

**Scenario:** Expose game state for testing and verification.

**Code:**
```go
type GameState struct {
    Player struct {
        X      float64
        Y      float64
        Health int
    }
    Inventory []string
}

func main() {
    game := &GameState{}
    autoebiten.RegisterStateExporter("gamestate", game)

    // ... game loop
}
```

**CLI usage:**
```bash
autoebiten state --name gamestate --path Player.Health
autoebiten state --name gamestate --path Inventory.0
```

---

### Scripted Automation

**Scenario:** Complex input sequences.

**Script:**
```json
{
  "version": "1.0",
  "commands": [
    {"input": {"action": "press", "key": "KeyW"}},
    {"delay": {"ms": 100}},
    {"repeat": {"times": 3, "commands": [
      {"input": {"action": "press", "key": "KeyA"}},
      {"delay": {"ms": 200}}
    ]}},
    {"screenshot": {"output": "final.png"}}
  ]
}
```

**How to run:**
```bash
autoebiten run --script automation.json
```

---

### Multiple Game Instances

**Scenario:** Control specific game when multiple are running.

**Detection:**
```bash
# Auto-detect (fails if multiple)
autoebiten ping

# Specify PID
autoebiten --pid 12345 input --key KeySpace

# Or set socket path
AUTOEBITEN_SOCKET=/tmp/autoebiten/autoebiten-custom.sock autoebiten ping
```