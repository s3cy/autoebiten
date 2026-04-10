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


**Expected output:**
```
OK: game is running
```


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
autoebiten input --key KeyH --action press
```

**Output:**
```
OK: input press KeyH
```

```bash
autoebiten mouse --action position -x 100 -y 200
```

**Output:**
```
OK: mouse position at (100, 200)
```

```bash
autoebiten screenshot
```

**Output:**
```
OK: screenshot saved to /Users/s3cy/Desktop/go/autoebiten/.worktrees/doc-template-system/screenshot_20260410174313.png
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
autoebiten input --key KeySpace --action hold --duration_ticks 30
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
autoebiten custom --name getPlayerInfo
```

**Expected output:**
```
OK: Health: 100, Mana: 50
```


---

## Examples

### Library Integration

**Goal:** Full integration for testing and automation.

**Complete example:**

```go
package main

import (
    "fmt"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/s3cy/autoebiten"
)

type Game struct {
    PlayerX, PlayerY float64
    Health           int
}

func (g *Game) Update() error {
    // Handle automation exit
    if !autoebiten.Update() {
        return fmt.Errorf("exit requested")
    }

    // Use wrapped input functions
    if autoebiten.IsKeyPressed(ebiten.KeyArrowRight) {
        g.PlayerX += 2
    }
    if autoebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
        g.PlayerX -= 2
    }

    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    // Your rendering code
    // ...

    // Capture for screenshots (call at end)
    autoebiten.Capture(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return 640, 480
}

func main() {
    game := &Game{Health: 100}

    // Register custom commands
    autoebiten.Register("heal", func(ctx autoebiten.CommandContext) {
        game.Health = min(game.Health+20, 100)
        ctx.Respond(fmt.Sprintf("Healed to %d", game.Health))
    })

    // Register state exporter
    autoebiten.RegisterStateExporter("gamestate", game)

    ebiten.RunGame(game)
}
```

---

### Launch and Output Capture

**Scenario:** Capture game output between commands.

**Using launch:**

```bash
# Start with output capture
autoebiten launch -- ./mygame

# Commands show output changes
autoebiten ping
# OK: game is running

# Output is captured between commands
autoebiten input --key KeySpace
# Shows any new output from game

# Exit cleans up
autoebiten exit
```

**Benefits:**
- Captures stdout/stderr between commands
- Shows crash diagnostics automatically
- Maintains connection state

---

### Crash Diagnostic Testing

**Scenario:** Test crash detection before and after RPC connection.

**Example game:** `examples/crash_diagnostic/main.go`

This demo game can be configured to crash at different points:

```bash
# Build the example
cd examples/crash_diagnostic
go build -o crash_demo
```

**Test 1: Normal operation (no crash)**


```
OK: game is running
```


**Test 2: Crash BEFORE RPC connection**

When the game crashes before RPC connection, the launch command captures the error:

```text
Starting crash diagnostic demo...
Flags: crash-before-rpc=true, crash-after-rpc=false
Initialization complete
About to crash before RPC connection!
panic: intentional crash before RPC connection
Error: game failed to start
```

**Test 3: Crash AFTER RPC connection**

```bash
autoebiten ping
```
```
OK: game is running
```
Wait for crash (~3 seconds), then:
```bash
sleep 4
autoebiten ping
```
```text
<log_diff>
--- snapshot ...
+++ current ...
@@ ... @@
+panic: runtime error: index out of range
</log_diff>
<proxy_error>
game exited: exit status 2
</proxy_error>
Error: game not connected
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
autoebiten custom --name heal
```

**Output:**
```
OK: Healed from 100 to 100
```

```bash
autoebiten custom --name damage
```

**Output:**
```
OK: Damaged from 100 to 90
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