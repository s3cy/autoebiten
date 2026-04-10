# Integration Guide

> Purpose: Setup and integration guide for autoebiten — choose your method and integrate with your game
> Audience: Developers setting up autoebiten for the first time

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

**CLI Usage:**

Check if custom commands are available:
```json
Error: no running game found
Usage:
  autoebiten list_custom [flags]

Flags:
  -h, --help   help for list_custom

Global Flags:
  -p, --pid int   Target game process PID (auto-detected if not specified)

Error: no running game found
```

Execute a custom command:
```text
Error: no running game found
Usage:
  autoebiten custom [name] [flags]

Flags:
  -h, --help             help for custom
  -n, --name string      Custom command name
      --no-record        Skip recording this command
  -r, --request string   Request data to pass to the command

Global Flags:
  -p, --pid int   Target game process PID (auto-detected if not specified)

Error: no running game found
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

**Important:** State Exporter uses reflection internally. Only **exported fields** (capitalized names) are accessible. Unexported fields (lowercase) will not be queryable.

**Path Navigation Rules:**
- Use **Go field names**
- JSON tags like `json:"player_name"` are ignored for path navigation
- Interface fields can be queried directly, but cannot traverse into nested fields

**Example with JSON tags:**
```go
type Player struct {
    Name string `json:"player_name"`  // Query as "Name", NOT "player_name"
    Health int    `json:"hp"`          // Query as "Health", NOT "hp"
}
```

**Supported paths:**
- `Player.X` - Struct field (must be exported/capitalized)
- `Player.Name` - Works even with `json:"player_name"` tag
- `Inventory.0.Name` - Array/slice index
- `Skills.Sword` - Map key
- `Entity` - Interface field itself (returns the stored value)
- `Entity.Field` - NOT supported (cannot traverse into interfaces)

**CLI Usage:**

Use the `state` command to query exported state:

```bash
autoebiten state --name gamestate --path Player.Health
autoebiten state --name gamestate --path Enemies.0.Name
```

---

## CLI Connection Verification

After integrating autoebiten into your game, verify the connection:

```text
Error: no running game found
Usage:
  autoebiten ping [flags]

Flags:
  -h, --help   help for ping

Global Flags:
  -p, --pid int   Target game process PID (auto-detected if not specified)

Error: no running game found
```
