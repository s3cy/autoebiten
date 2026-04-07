# Documentation Rewrite Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create LLM-optimized documentation for autoebiten CLI and testkit, with comprehensive API references, tutorials, and examples.

**Architecture:** Two main documentation files (cli.md, testkit.md) following a self-contained, flat structure. Each file combines API reference, tutorial, and examples. Source code extraction ensures accuracy - no reliance on existing markdown.

**Tech Stack:** Go (Ebitengine), Markdown documentation, JSON scripts

---

## File Structure

| File | Action | Lines |
|------|--------|-------|
| `docs/SPEC.md` | Move from root, update | Existing |
| `docs/cli.md` | Create | ~400-500 |
| `docs/testkit.md` | Create | ~300-400 |
| `README.md` | Add navigation section | ~20 new |
| `examples/scripts/basic.json` | Create | ~15 |
| `examples/scripts/mouse.json` | Create | ~20 |
| `examples/scripts/complex.json` | Create | ~25 |
| `examples/state_exporter/main.go` | Create | ~80 |
| `examples/testkit/black_box_test.go` | Create | ~60 |
| `examples/testkit/white_box_test.go` | Create | ~40 |

---

## Task 1: Move SPEC.md to docs/SPEC.md

**Files:**
- Move: `SPEC.md` → `docs/SPEC.md`
- Update: README.md (add navigation link)

- [ ] **Step 1: Create docs directory if needed**

```bash
mkdir -p docs
```

- [ ] **Step 2: Move SPEC.md**

```bash
git mv SPEC.md docs/SPEC.md
```

- [ ] **Step 3: Verify the move**

```bash
ls -la docs/SPEC.md
```

Expected: File exists at docs/SPEC.md

- [ ] **Step 4: Commit**

```bash
git add docs/SPEC.md
git commit -m "docs: move SPEC.md to docs directory"
```

---

## Task 2: Create docs/cli.md - Quick Decision and Overview Sections

**Files:**
- Create: `docs/cli.md`

**Source files to reference:**
- `cmd/autoebiten/main.go` - CLI commands and flags
- `autoebiten.go` - Mode types and SetMode/GetMode
- `integrate/integrate.go` - Patch method integration functions
- `autoebiten_default.go` - Library method Update/Capture
- `internal/cli/commands.go` - Command executor

- [ ] **Step 1: Create docs/cli.md with header and overview**

```markdown
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

**Input mode:**
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

Communication uses JSON-RPC over Unix sockets. Each game instance creates a socket at `/tmp/autoebiten/autoebiten-{socket-hash}.sock`.

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

```

Write this content to `docs/cli.md` (first section).

- [ ] **Step 2: Verify file created**

```bash
head -50 docs/cli.md
```

Expected: Header and Overview sections visible

---

## Task 3: Create docs/cli.md - API Reference Section

**Files:**
- Edit: `docs/cli.md`

- [ ] **Step 1: Add Core Functions API Reference**

Extract signatures from `integrate/integrate.go` and `autoebiten.go`.

```markdown
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

**Signature:** `func Capture(screen *ebiten.Image)`

**Purpose:** Enable screenshot capture for CLI requests.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| screen | *ebiten.Image | The screen image from Draw |

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

```

Append this to `docs/cli.md`.

- [ ] **Step 2: Add Custom Commands API Reference**

Extract from `custom_command.go` and `state_exporter.go`.

```markdown
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

```

Append this to `docs/cli.md`.

---

## Task 4: Create docs/cli.md - CLI Commands Section

**Files:**
- Edit: `docs/cli.md`

- [ ] **Step 1: Add CLI Commands Reference**

Extract from `cmd/autoebiten/main.go` command definitions.

```markdown
---

## CLI Commands

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

# Release a held key
autoebiten input --key KeyW --action release
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
# Move cursor to position
autoebiten mouse --action position -x 100 -y 200

# Click at current position
autoebiten mouse --action press --button MouseButtonLeft

# Click at specific position
autoebiten mouse -x 100 -y 200 --button MouseButtonLeft

# Hold right button for 20 ticks
autoebiten mouse --button MouseButtonRight --duration_ticks 20

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

# Scroll down 2 units
autoebiten wheel -y 2

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
| --async | Return immediately |
| --no-record | Skip recording |

**Examples:**
```bash
# Save to auto-generated filename
autoebiten screenshot

# Save to specific file
autoebiten screenshot --output capture.png

# Get base64 (useful for scripts)
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

Get injected mouse cursor position (not real OS position).

```bash
autoebiten get_mouse_position
```

---

#### get_wheel_position

Get injected wheel position (not real OS position).

```bash
autoebiten get_wheel_position
```

---

#### schema

Output JSON Schema for script files (for IDE validation).

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

**Examples:**
```bash
# Execute without data
autoebiten custom getPlayerInfo

# Execute with data
autoebiten custom echo --request "hello world"
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
- value: JSON value (number, string, boolean)

**Examples:**
```bash
# Wait for health to reach 100
autoebiten wait-for --condition "state:gamestate:Player.Health == 100" --timeout 10s

# Wait for X position > 50
autoebiten wait-for --condition "state:gamestate:Player.X > 50" --timeout 30s
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

```

Append this to `docs/cli.md`.

---

## Task 5: Create docs/cli.md - Tutorial Section

**Files:**
- Edit: `docs/cli.md`

- [ ] **Step 1: Add Tutorial Section**

```markdown
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

```

Append this to `docs/cli.md`.

---

## Task 6: Create docs/cli.md - Examples Section

**Files:**
- Edit: `docs/cli.md`

- [ ] **Step 1: Add Examples Section**

```markdown
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

```

Append this to `docs/cli.md`.

---

## Task 7: Create docs/testkit.md - Complete File

**Files:**
- Create: `docs/testkit.md`

**Source files:**
- `testkit/game.go` - Black-box Game type
- `testkit/mock.go` - White-box Mock type
- `testkit/options.go` - Options
- `testkit/doc.go` - Package doc
- `testkit/errors.go` - Errors

- [ ] **Step 1: Create docs/testkit.md with full content**

```markdown
# Testkit Reference

> Purpose: Go testing framework for autoebiten games
> Audience: Developers writing automated tests for Ebitengine games

---

## Quick Decision

**Testing mode:**
```
├─ Need full integration (real game process, real socket)?
│  └─→ Black-box: testkit.Launch()
│
└─ Need fast unit tests (game logic only)?
   └─→ White-box: testkit.NewMock()
```

**State verification:**
```
├─ Need to query internal state?
│  └─→ Add autoebiten.RegisterStateExporter() to game
│      Then use game.StateQuery()
│
└─ Only need visual verification?
   └─→ Use game.Screenshot() or ScreenshotToFile()
```

---

## Overview

The testkit package provides two testing modes:

- **Black-box (Game):** Launches game in separate process, controls via RPC. Tests the full system.
- **White-box (Mock):** Tests game logic in same process with mocked inputs. Fast unit testing.

Both modes use similar APIs for input injection, making tests easy to write and maintain.

---

## API Reference

### Black-Box Testing (Game)

#### Launch

**Signature:** `func Launch(t *testing.T, binaryPath string, opts ...Option) *Game`

**Purpose:** Start game in separate process and return controller.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| t | *testing.T | Test context |
| binaryPath | string | Path to compiled game binary |
| opts | ...Option | Optional configuration |

**Returns:**
- `*Game` - Controller for the launched game

**Options:**
- `WithTimeout(d time.Duration)` - Set timeout (default 30s)
- `WithArgs(args ...string)` - Add command-line arguments
- `WithEnv(key, value string)` - Set environment variable

**Example:**
```go
func TestGame(t *testing.T) {
    game := testkit.Launch(t, "./mygame",
        testkit.WithTimeout(30*time.Second))
    defer game.Shutdown()
}
```

**Notes:** Game binary must be built before testing.

---

#### Shutdown

**Signature:** `func (g *Game) Shutdown() error`

**Purpose:** Gracefully stop the game process.

**Returns:**
- error if shutdown fails

**Notes:** Automatically called via t.Cleanup(). Call manually for early cleanup.

---

#### Ping

**Signature:** `func (g *Game) Ping() error`

**Purpose:** Check if game is responsive.

**Returns:**
- nil if responsive
- error if not running or not responding

---

#### WaitFor

**Signature:** `func (g *Game) WaitFor(fn func() bool, timeout time.Duration) bool`

**Purpose:** Poll until condition is true or timeout.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| fn | func() bool | Condition function (called every 100ms) |
| timeout | time.Duration | Maximum wait time |

**Returns:**
- `true` if condition met
- `false` if timeout

**Example:**
```go
ready := game.WaitFor(func() bool {
    return game.Ping() == nil
}, 5*time.Second)
```

---

#### PressKey

**Signature:** `func (g *Game) PressKey(key ebiten.Key) error`

**Purpose:** Send key press (press and release).

---

#### HoldKey

**Signature:** `func (g *Game) HoldKey(key ebiten.Key, ticks int64) error`

**Purpose:** Hold key for specified ticks.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| key | ebiten.Key | Key to hold |
| ticks | int64 | Duration in ticks |

---

#### MoveMouse

**Signature:** `func (g *Game) MoveMouse(x, y int) error`

**Purpose:** Move cursor to position.

---

#### PressMouseButton

**Signature:** `func (g *Game) PressMouseButton(button ebiten.MouseButton) error`

**Purpose:** Click mouse button.

---

#### ScrollWheel

**Signature:** `func (g *Game) ScrollWheel(x, y float64) error`

**Purpose:** Inject wheel scroll.

---

#### Screenshot

**Signature:** `func (g *Game) Screenshot() (*image.RGBA, error)`

**Purpose:** Capture screen as image.

**Returns:**
- Screen image
- error if capture fails

---

#### ScreenshotToFile

**Signature:** `func (g *Game) ScreenshotToFile(path string) error`

**Purpose:** Save screenshot to file.

---

#### ScreenshotBase64

**Signature:** `func (g *Game) ScreenshotBase64() (string, error)`

**Purpose:** Get screenshot as base64 string.

---

#### StateQuery

**Signature:** `func (g *Game) StateQuery(name string, path string) (any, error)`

**Purpose:** Query game state via registered exporter.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| name | string | State exporter name |
| path | string | Dot-notation path |

**Returns:**
- Value at path (type depends on state)
- error if path not found

**Example:**
```go
health, err := game.StateQuery("gamestate", "Player.Health")
require.NoError(t, err)
assert.Equal(t, 100, health)
```

**Notes:** Requires `RegisterStateExporter` in game code.

---

#### RunCustom

**Signature:** `func (g *Game) RunCustom(name, request string) (string, error)`

**Purpose:** Execute custom command.

**Returns:**
- Response string from command
- error if command fails

---

### White-Box Testing (Mock)

#### NewMock

**Signature:** `func NewMock(t *testing.T, game GameUpdate) *Mock`

**Purpose:** Create mock controller for in-process testing.

**Parameters:**
| Name | Type | Description |
|------|------|-------------|
| t | *testing.T | Test context |
| game | GameUpdate | Game instance (must have Update method) |

**Returns:**
- `*Mock` - Controller for mocked inputs

**Example:**
```go
func TestLogic(t *testing.T) {
    autoebiten.SetMode(autoebiten.InjectionOnly)

    game := NewMyGame()
    mock := testkit.NewMock(t, game)

    mock.InjectKeyPress(ebiten.KeyD)
    mock.Ticks(10)

    assert.Equal(t, 10, game.Player.X)
}
```

---

#### InjectKeyPress

**Signature:** `func (m *Mock) InjectKeyPress(key ebiten.Key)`

**Purpose:** Buffer key press for next Tick.

---

#### InjectKeyRelease

**Signature:** `func (m *Mock) InjectKeyRelease(key ebiten.Key)`

**Purpose:** Buffer key release for next Tick.

---

#### InjectMousePosition

**Signature:** `func (m *Mock) InjectMousePosition(x, y int)`

**Purpose:** Set mouse position for next Tick.

---

#### InjectMouseButtonPress

**Signature:** `func (m *Mock) InjectMouseButtonPress(button ebiten.MouseButton)`

**Purpose:** Buffer mouse press for next Tick.

---

#### InjectMouseButtonRelease

**Signature:** `func (m *Mock) InjectMouseButtonRelease(button ebiten.MouseButton)`

**Purpose:** Buffer mouse release for next Tick.

---

#### InjectWheel

**Signature:** `func (m *Mock) InjectWheel(x, y float64)`

**Purpose:** Set wheel delta for next Tick.

---

#### Tick

**Signature:** `func (m *Mock) Tick()`

**Purpose:** Advance game by one tick, applying buffered inputs.

**Notes:** Calls game's Update() once.

---

#### Ticks

**Signature:** `func (m *Mock) Ticks(n int)`

**Purpose:** Advance game by N ticks.

---

#### Game

**Signature:** `func (m *Mock) Game() GameUpdate`

**Purpose:** Get underlying game instance for assertions.

---

### GameUpdate Interface

**Definition:**
```go
type GameUpdate interface {
    Update() error
}
```

Games must implement at least `Update() error` for white-box testing.

---

### Errors

| Error | Description |
|-------|-------------|
| ErrGameNotRunning | Operation on non-running game |
| ErrTimeout | Operation timed out |
| ErrInvalidState | Invalid game state |

---

## Tutorial

### Step 1: Choose Testing Mode

**Black-box checklist:**
- [ ] Need real game process
- [ ] Need RPC communication
- [ ] Need screenshot capture

**White-box checklist:**
- [ ] Fast unit tests
- [ ] Direct state access
- [ ] No external dependencies

---

### Step 2: Black-Box Setup

**Goal:** Launch game in separate process.

**Prerequisites:**
- Built game binary

**Actions:**

```bash
# Build game
go build -o ./mygame ./cmd/mygame
```

```go
func TestGame(t *testing.T) {
    game := testkit.Launch(t, "./mygame")
    defer game.Shutdown()

    // Wait for ready
    ready := game.WaitFor(func() bool {
        return game.Ping() == nil
    }, 5*time.Second)
    require.True(t, ready)
}
```

---

### Step 3: White-Box Setup

**Goal:** Test game logic in-process.

```go
func TestPlayerMovement(t *testing.T) {
    autoebiten.SetMode(autoebiten.InjectionOnly)

    game := NewMyGame()
    mock := testkit.NewMock(t, game)

    // Inject input
    mock.InjectKeyPress(ebiten.KeyArrowRight)
    mock.Ticks(10)

    // Check state directly
    assert.Greater(t, game.Player.X, 0.0)
}
```

---

### Step 4: Adding State Exporter

**Goal:** Enable StateQuery for black-box tests.

**In game:**
```go
type GameState struct {
    Player struct {
        X, Y float64
    }
}

func main() {
    state := &GameState{}
    autoebiten.RegisterStateExporter("gamestate", state)
    // ...
}
```

**In test:**
```go
x, err := game.StateQuery("gamestate", "Player.X")
require.NoError(t, err)
assert.Equal(t, 10.0, x)
```

---

### Step 5: Writing Assertions

**Goal:** Verify game behavior.

```go
// Using testify
func TestPlayerHealth(t *testing.T) {
    game := testkit.Launch(t, "./mygame")
    defer game.Shutdown()

    // Get initial state
    health, err := game.StateQuery("gamestate", "Player.Health")
    require.NoError(t, err)
    initial := health.(float64)

    // Apply damage
    _, err = game.RunCustom("damage", "")
    require.NoError(t, err)

    // Verify change
    health, err = game.StateQuery("gamestate", "Player.Health")
    require.NoError(t, err)
    assert.Less(t, health.(float64), initial)
}
```

---

## Examples

### Black-Box: Player Movement

**Scenario:** Test movement via real game process.

```go
package mygame_test

import (
    "testing"
    "time"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/s3cy/autoebiten/testkit"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestPlayerMovesRight(t *testing.T) {
    game := testkit.Launch(t, "./mygame",
        testkit.WithTimeout(30*time.Second))
    defer game.Shutdown()

    // Wait for ready
    require.True(t, game.WaitFor(func() bool {
        return game.Ping() == nil
    }, 5*time.Second))

    // Get initial position
    x, err := game.StateQuery("gamestate", "Player.X")
    require.NoError(t, err)
    initialX := x.(float64)

    // Move right
    err = game.HoldKey(ebiten.KeyArrowRight, 10)
    require.NoError(t, err)

    // Verify moved
    x, err = game.StateQuery("gamestate", "Player.X")
    require.NoError(t, err)
    assert.Greater(t, x.(float64), initialX)
}
```

---

### White-Box: Game Logic

**Scenario:** Fast unit test without RPC.

```go
package mygame

import (
    "testing"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/s3cy/autoebiten"
    "github.com/s3cy/autoebiten/testkit"
    "github.com/stretchr/testify/assert"
)

func TestPlayerTakesDamage(t *testing.T) {
    autoebiten.SetMode(autoebiten.InjectionOnly)

    game := NewGame()
    mock := testkit.NewMock(t, game)

    initialHealth := game.Player.Health

    // Simulate damage key
    mock.InjectKeyPress(ebiten.KeyD)
    mock.Tick()

    assert.Less(t, game.Player.Health, initialHealth)
}

func TestComboInput(t *testing.T) {
    autoebiten.SetMode(autoebiten.InjectionOnly)

    game := NewGame()
    mock := testkit.NewMock(t, game)

    // Combo: A + B
    mock.InjectKeyPress(ebiten.KeyA)
    mock.InjectKeyPress(ebiten.KeyB)
    mock.Tick()

    assert.True(t, game.ComboActivated)
}
```

---

### State Query: Health Verification

**Scenario:** Query nested state via reflection.

```go
func TestInventoryAccess(t *testing.T) {
    game := testkit.Launch(t, "./mygame")
    defer game.Shutdown()

    // Query array index
    name, err := game.StateQuery("gamestate", "Inventory.0.Name")
    require.NoError(t, err)
    assert.Equal(t, "Sword", name)

    // Query nested field
    count, err := game.StateQuery("gamestate", "Inventory.0.Count")
    require.NoError(t, err)
    assert.Equal(t, 5, count)
}
```
```

Write this to `docs/testkit.md`.

- [ ] **Step 2: Verify file created**

```bash
wc -l docs/testkit.md
```

Expected: ~300-400 lines

---

## Task 8: Create JSON Script Examples

**Files:**
- Create: `examples/scripts/basic.json`
- Create: `examples/scripts/mouse.json`
- Create: `examples/scripts/complex.json`

- [ ] **Step 1: Create examples/scripts directory**

```bash
mkdir -p examples/scripts
```

- [ ] **Step 2: Create basic.json**

```json
{
  "version": "1.0",
  // Basic keyboard automation
  "commands": [
    {"input": {"action": "press", "key": "KeySpace"}},
    {"delay": {"ms": 500}},
    {"input": {"action": "hold", "key": "KeyW", "duration_ticks": 30}},
    {"screenshot": {"output": "basic-result.png"}}
  ]
}
```

Write to `examples/scripts/basic.json`.

- [ ] **Step 3: Create mouse.json**

```json
{
  "version": "1.0",
  // Mouse automation example
  "commands": [
    // Move cursor to center
    {"mouse": {"action": "position", "x": 320, "y": 240}},
    {"delay": {"ms": 100}},
    // Click left button
    {"mouse": {"action": "press", "button": "MouseButtonLeft"}},
    {"delay": {"ms": 200}},
    // Drag: move while holding
    {"mouse": {"action": "position", "x": 400, "y": 300}},
    {"mouse": {"action": "hold", "button": "MouseButtonLeft", "duration_ticks": 10}},
    // Scroll down
    {"wheel": {"x": 0, "y": 3}},
    {"screenshot": {"output": "mouse-result.png"}}
  ]
}
```

Write to `examples/scripts/mouse.json`.

- [ ] **Step 4: Create complex.json**

```json
{
  "version": "1.0",
  // Complex automation: combat sequence
  "commands": [
    /* Approach enemy */
    {"input": {"action": "hold", "key": "KeyW", "duration_ticks": 20}},
    {"delay": {"ms": 300}},
    /* Attack combo - repeat 3 times */
    {"repeat": {"times": 3, "commands": [
      {"input": {"action": "press", "key": "KeySpace"}},
      {"delay": {"ms": 100}}
    ]}},
    /* Wait for enemy defeat */
    {"wait": {"condition": "state:gamestate:Enemy.Health == 0", "timeout": "10s"}},
    /* Loot */
    {"mouse": {"action": "position", "x": 200, "y": 150}},
    {"mouse": {"button": "MouseButtonLeft"}},
    {"delay": {"ms": 500}},
    /* Check player health */
    {"state": {"name": "gamestate", "path": "Player.Health"}},
    /* Final screenshot */
    {"screenshot": {"output": "combat-done.png"}}
  ]
}
```

Write to `examples/scripts/complex.json`.

- [ ] **Step 5: Verify examples**

```bash
ls -la examples/scripts/
```

Expected: 3 JSON files

---

## Task 9: Create State Exporter Example

**Files:**
- Create: `examples/state_exporter/main.go`

- [ ] **Step 1: Create examples/state_exporter directory**

```bash
mkdir -p examples/state_exporter
```

- [ ] **Step 2: Create main.go**

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

// GameState holds all game data for state queries.
type GameState struct {
    Player struct {
        X      float64
        Y      float64
        Health int
        Mana   int
    }
    Enemies []Enemy
    Score   int
}

type Enemy struct {
    Name   string
    Health int
    X, Y   float64
}

// Game implements ebiten.Game interface.
type Game struct {
    state GameState
}

func NewGame() *Game {
    g := &Game{
        state: GameState{
            Player: struct {
                X      float64
                Y      float64
                Health int
                Mana   int
            }{X: 100, Y: 100, Health: 100, Mana: 50},
            Enemies: []Enemy{
                {Name: "Goblin", Health: 30, X: 300, Y: 200},
                {Name: "Orc", Health: 50, X: 400, Y: 300},
            },
            Score: 0,
        },
    }

    // Register state exporter for StateQuery
    autoebiten.RegisterStateExporter("gamestate", &g.state)

    // Register custom commands
    autoebiten.Register("heal", func(ctx autoebiten.CommandContext) {
        old := g.state.Player.Health
        g.state.Player.Health = min(g.state.Player.Health+20, 100)
        ctx.Respond(fmt.Sprintf("Healed from %d to %d", old, g.state.Player.Health))
    })

    return g
}

func (g *Game) Update() error {
    if !autoebiten.Update() {
        return fmt.Errorf("exit requested")
    }

    // Movement
    speed := 2.0
    if autoebiten.IsKeyPressed(ebiten.KeyArrowRight) {
        g.state.Player.X += speed
    }
    if autoebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
        g.state.Player.X -= speed
    }
    if autoebiten.IsKeyPressed(ebiten.KeyArrowUp) {
        g.state.Player.Y -= speed
    }
    if autoebiten.IsKeyPressed(ebiten.KeyArrowDown) {
        g.state.Player.Y += speed
    }

    // Damage
    if autoebiten.IsKeyPressed(ebiten.KeyD) {
        g.state.Player.Health = max(g.state.Player.Health-5, 0)
    }

    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    screen.Fill(color.RGBA{0x10, 0x20, 0x40, 0xff})

    msg := "=== State Exporter Demo ===\n\n"
    msg += fmt.Sprintf("Player: (%.1f, %.1f)\n", g.state.Player.X, g.state.Player.Y)
    msg += fmt.Sprintf("Health: %d  Mana: %d\n", g.state.Player.Health, g.state.Player.Mana)
    msg += fmt.Sprintf("Score: %d\n", g.state.Score)
    msg += "\nEnemies:\n"
    for i, e := range g.state.Enemies {
        msg += fmt.Sprintf("  %d: %s (HP:%d)\n", i, e.Name, e.Health)
    }
    msg += "\nCLI Commands:\n"
    msg += "  state --name gamestate --path Player.Health\n"
    msg += "  state --name gamestate --path Enemies.0.Name\n"

    ebitenutil.DebugPrint(screen, msg)
    autoebiten.Capture(screen)
}

func (g *Game) Layout(_, _ int) (int, int) {
    return screenWidth, screenHeight
}

func main() {
    ebiten.SetWindowSize(screenWidth, screenHeight)
    ebiten.SetWindowTitle("State Exporter Demo")

    if err := ebiten.RunGame(NewGame()); err != nil {
        log.Fatal(err)
    }
}
```

Write to `examples/state_exporter/main.go`.

- [ ] **Step 3: Verify example**

```bash
go build ./examples/state_exporter/
```

Expected: Build succeeds

---

## Task 10: Create Testkit Examples

**Files:**
- Create: `examples/testkit/black_box_test.go`
- Create: `examples/testkit/white_box_test.go`

- [ ] **Step 1: Create examples/testkit directory**

```bash
mkdir -p examples/testkit
```

- [ ] **Step 2: Create black_box_test.go**

```go
package testkit_test

import (
    "testing"
    "time"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/s3cy/autoebiten/testkit"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// TestPlayerMovement demonstrates black-box testing with StateQuery.
// Requires a game built from examples/state_exporter.
func TestPlayerMovement(t *testing.T) {
    // Launch game in separate process
    game := testkit.Launch(t, "./examples/state_exporter/state_exporter",
        testkit.WithTimeout(30*time.Second))
    defer game.Shutdown()

    // Wait for game to be ready
    ready := game.WaitFor(func() bool {
        return game.Ping() == nil
    }, 5*time.Second)
    require.True(t, ready, "game should be ready within timeout")

    // Get initial position
    x, err := game.StateQuery("gamestate", "Player.X")
    require.NoError(t, err)
    initialX := x.(float64)

    // Move player right for 10 ticks
    err = game.HoldKey(ebiten.KeyArrowRight, 10)
    require.NoError(t, err)

    // Verify position changed
    x, err = game.StateQuery("gamestate", "Player.X")
    require.NoError(t, err)
    assert.Greater(t, x.(float64), initialX, "player should have moved right")
}

// TestHealthModification demonstrates custom commands and state verification.
func TestHealthModification(t *testing.T) {
    game := testkit.Launch(t, "./examples/state_exporter/state_exporter")
    defer game.Shutdown()

    game.WaitFor(func() bool { return game.Ping() == nil }, 5*time.Second)

    // Get initial health
    health, err := game.StateQuery("gamestate", "Player.Health")
    require.NoError(t, err)
    require.Equal(t, 100.0, health, "initial health should be 100")

    // Heal via custom command
    resp, err := game.RunCustom("heal", "")
    require.NoError(t, err)
    assert.Contains(t, resp, "Healed")

    // Health unchanged (already at max)
    health, err = game.StateQuery("gamestate", "Player.Health")
    require.NoError(t, err)
    assert.Equal(t, 100.0, health)
}

// TestScreenshotCapture demonstrates visual verification.
func TestScreenshotCapture(t *testing.T) {
    game := testkit.Launch(t, "./examples/state_exporter/state_exporter")
    defer game.Shutdown()

    game.WaitFor(func() bool { return game.Ping() == nil }, 5*time.Second)

    // Capture screenshot
    img, err := game.Screenshot()
    require.NoError(t, err)
    require.NotNil(t, img)

    // Verify dimensions
    bounds := img.Bounds()
    assert.Equal(t, 640, bounds.Dx())
    assert.Equal(t, 480, bounds.Dy())

    // Save to file
    err = game.ScreenshotToFile("testdata/screenshot.png")
    require.NoError(t, err)
}

// TestEnemyStateQuery demonstrates querying array/slice state.
func TestEnemyStateQuery(t *testing.T) {
    game := testkit.Launch(t, "./examples/state_exporter/state_exporter")
    defer game.Shutdown()

    game.WaitFor(func() bool { return game.Ping() == nil }, 5*time.Second)

    // Query first enemy name
    name, err := game.StateQuery("gamestate", "Enemies.0.Name")
    require.NoError(t, err)
    assert.Equal(t, "Goblin", name)

    // Query second enemy health
    health, err := game.StateQuery("gamestate", "Enemies.1.Health")
    require.NoError(t, err)
    assert.Equal(t, 50.0, health)
}
```

Write to `examples/testkit/black_box_test.go`.

- [ ] **Step 3: Create white_box_test.go**

```go
package testkit_test

import (
    "testing"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/s3cy/autoebiten"
    "github.com/s3cy/autoebiten/examples/state_exporter"
    "github.com/s3cy/autoebiten/testkit"
    "github.com/stretchr/testify/assert"
)

// TestPlayerMovesRight demonstrates white-box testing with Mock.
func TestPlayerMovesRight(t *testing.T) {
    // Set mode to only use injected inputs
    autoebiten.SetMode(autoebiten.InjectionOnly)

    // Create game instance in same process
    game := state_exporter.NewGame()
    mock := testkit.NewMock(t, game)

    initialX := game.state.Player.X

    // Inject input
    mock.InjectKeyPress(ebiten.KeyArrowRight)
    mock.Ticks(10)

    // Check state directly (no RPC needed)
    assert.Greater(t, game.state.Player.X, initialX)
}

// TestPlayerTakesDamage demonstrates damage logic testing.
func TestPlayerTakesDamage(t *testing.T) {
    autoebiten.SetMode(autoebiten.InjectionOnly)

    game := state_exporter.NewGame()
    mock := testkit.NewMock(t, game)

    initialHealth := game.state.Player.Health

    // Simulate damage key
    mock.InjectKeyPress(ebiten.KeyD)
    mock.Tick()

    assert.Less(t, game.state.Player.Health, initialHealth)
}

// TestComboInput demonstrates multiple inputs in one tick.
func TestComboInput(t *testing.T) {
    autoebiten.SetMode(autoebiten.InjectionOnly)

    game := state_exporter.NewGame()
    mock := testkit.NewMock(t, game)

    // Inject multiple inputs before tick
    mock.InjectKeyPress(ebiten.KeyArrowUp)
    mock.InjectKeyPress(ebiten.KeyArrowRight)
    mock.Ticks(5)

    // Player should have moved diagonally
    assert.Greater(t, game.state.Player.X, 100.0)
    assert.Less(t, game.state.Player.Y, 100.0)
}

// TestMouseInput demonstrates mouse position injection.
func TestMouseInput(t *testing.T) {
    autoebiten.SetMode(autoebiten.InjectionOnly)

    game := state_exporter.NewGame()
    mock := testkit.NewMock(t, game)

    // Inject mouse position
    mock.InjectMousePosition(200, 150)
    mock.Tick()

    // In a real game, you'd check mouse-dependent behavior
    // For this example, we just verify injection works
    x, y := autoebiten.CursorPosition()
    assert.Equal(t, 200, x)
    assert.Equal(t, 150, y)
}

// TestWheelInput demonstrates wheel scroll injection.
func TestWheelInput(t *testing.T) {
    autoebiten.SetMode(autoebiten.InjectionOnly)

    game := state_exporter.NewGame()
    mock := testkit.NewMock(t, game)

    mock.InjectWheel(0, -5)
    mock.Tick()

    x, y := autoebiten.Wheel()
    assert.Equal(t, 0.0, x)
    assert.Equal(t, -5.0, y)
}
```

Write to `examples/testkit/white_box_test.go`.

- [ ] **Step 4: Create testdata directory for screenshots**

```bash
mkdir -p examples/testkit/testdata
```

---

## Task 11: Update README.md Navigation

**Files:**
- Edit: `README.md`

- [ ] **Step 1: Add Documentation section at top of README.md**

After the first description paragraph, add:

```markdown
## Documentation

- **CLI Automation** - Control games via command line → [docs/cli.md](docs/cli.md)
- **Game Testing** - Write automated tests → [docs/testkit.md](docs/testkit.md)
- **Technical Specs** - Architecture & protocol → [docs/SPEC.md](docs/SPEC.md)
```

---

## Task 12: Final Verification and Commit

- [ ] **Step 1: Verify all files exist**

```bash
ls -la docs/cli.md docs/testkit.md docs/SPEC.md examples/scripts/*.json examples/state_exporter/main.go examples/testkit/*.go
```

Expected: All files present

- [ ] **Step 2: Run go build on examples**

```bash
go build ./examples/state_exporter/
```

Expected: Build succeeds

- [ ] **Step 3: Commit all changes**

```bash
git add docs/ examples/scripts/ examples/state_exporter/ examples/testkit/ README.md
git commit -m "docs: rewrite documentation with LLM-optimized structure

- Move SPEC.md to docs/SPEC.md
- Create docs/cli.md: comprehensive CLI API reference, tutorial, examples
- Create docs/testkit.md: testkit API reference, tutorial, examples
- Add JSON script examples in examples/scripts/
- Add state_exporter demo game
- Add testkit example tests
- Update README.md with documentation navigation"
```

---

## Self-Review Checklist

After completing all tasks:

- [ ] Verify docs/cli.md has ~400-500 lines
- [ ] Verify docs/testkit.md has ~300-400 lines
- [ ] Check no placeholders (TBD, TODO) remain
- [ ] Verify all API signatures match source code
- [ ] Verify example code compiles
- [ ] Verify README.md links work