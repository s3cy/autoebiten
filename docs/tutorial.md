# Tutorial and Examples

> Purpose: Step-by-step tutorial and practical examples for using autoebiten
> Audience: Developers learning autoebiten through hands-on examples

---

## Tutorial

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

**Option A: Using launch (recommended)**
```bash
# Start game with output capture and crash diagnostics
autoebiten launch -- ./mygame &
```

**Option B: Running game directly**
```bash
# Start game in background
./mygame &

# Check connection
autoebiten ping
```

**Expected:** Output `game is running`

**Troubleshooting:**
- "connection failed": Game not started or socket not created
- "timeout waiting for RPC": Game takes longer to initialize, use `--timeout 30s`
- "game not connected": Game crashed; check `<log_diff>` for error output
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
# Build the demo
go build -o demo

# Option 1: Use launch for automation with output capture
autoebiten launch -- ./demo &
# In another terminal (or the same after launch starts):
autoebiten input --key KeySpace --action press

# Option 2: Run game directly
./demo &
autoebiten input --key KeySpace --action press
```

---

### Launch and Output Capture

**Scenario:** Automate game with visibility into game responses and crash detection.

**Using launch:**
```bash
# Start game with launch proxy
autoebiten launch -- ./mygame

# In another terminal, CLI commands connect through the proxy:
autoebiten ping
# Output shows unified diff of game output since last command:
# <log_diff>
# --- snapshot (empty)
# +++ current 2026-04-07 23:42:11.844
# @@ -0,0 +1,7 @@
# +2026-04-07 23:42:11.743 simple[37214:10490862] [CAMetalLayer nextDrawable] ...
# </log_diff>
# OK: game is running

autoebiten input --key KeySpace --action hold
# Shows diff of new output since ping command
```

**Crash Diagnostics:**
When the game crashes, commands show error context:
```bash
autoebiten ping
# <log_diff>
# --- snapshot ...
# +++ current ...
# @@ ... @@
# +panic: runtime error: index out of range [3] with length 3
# +main.(*Game).Update(0x14000112060)
# +    /path/to/game/main.go:42
# </log_diff>
# <proxy_error>
# game exited: exit status 2
# </proxy_error>
# Error: game not connected
```

**Key files:**
- Log: `/tmp/autoebiten/autoebiten-{PID}-output.log` - raw game output
- Launch socket: `/tmp/autoebiten/autoebiten-{PID}-launch.sock` - proxy connection

---

### Crash Diagnostic Testing

**Scenario:** Test crash detection before and after RPC connection.

**Example game:** `examples/crash_diagnostic/main.go`

This demo game can be configured to crash at different points:

```bash
# Build the example
cd examples/crash_diagnostic
go build -o crash_demo

# Test 1: Normal operation (no crash)
autoebiten launch -- ./crash_demo
autoebiten ping
# OK: game is running
autoebiten exit

# Test 2: Crash BEFORE RPC connection
autoebiten launch -- ./crash_demo --crash-before-rpc
autoebiten ping
# <log_diff>
# --- snapshot ...
# +++ current ...
# @@ ... @@
# +Starting crash diagnostic demo...
# +Flags: crash-before-rpc=true, crash-after-rpc=false
# +Initialization complete
# +About to crash before RPC connection!
# +panic: intentional crash before RPC connection
# </log_diff>
# <proxy_error>
# game exited: exit status 2
# </proxy_error>
# Error: game not connected

# Test 3: Crash AFTER RPC connection
autoebiten launch -- ./crash_demo --crash-after-rpc
autoebiten ping
# OK: game is running
sleep 4  # Wait for game to crash (~3 seconds)
autoebiten ping
# <log_diff>
# --- snapshot ...
# +++ current ...
# @@ ... @@
# +Game running... tick 60
# +Game running... tick 120
# +Game running... tick 180
# +panic: intentional crash after RPC connection
# </log_diff>
# <proxy_error>
# game exited: exit status 2
# </proxy_error>
# Error: game not connected
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