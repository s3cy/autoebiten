---
name: using-autoebiten
description: Use when automating Ebitengine games via CLI input injection, taking screenshots, or writing tests with the testkit package. REQUIRED when user mentions controlling games, sending inputs to games, or testing Ebitengine games. Use to decide between black-box (process-based) and white-box (in-process) testing approaches.
---

# Using AutoEbiten

## Overview

AutoEbiten automates Ebitengine games through input injection, screenshots, and scripted sequences. It provides two integration methods (Patch for existing games, Library for new games) and two testing modes (Black-box via RPC, White-box via Mock).

**ALWAYS follow this decision flow before taking action:**

1. Determine if writing tests or doing CLI automation
2. Determine integration method (find evidence first!)
3. Choose appropriate tool based on findings

---

## Step 1: Determine Your Purpose

**CRITICAL: Identify your goal before choosing any approach:**

| Purpose | Tool | Description |
|---------|------|-------------|
| CLI Automation | `autoebiten` CLI | Ad-hoc control of running games |
| Writing Tests | `testkit` package | Programmatic Go tests |

**Decision rule:**
- If the user wants to "test", or "automate testing" → Use **testkit**
- If the user wants to "run", "control", "send inputs", or "take screenshots" without mentioning tests → Use **CLI**

---

## Step 2: Determine Integration Method

**REQUIRED: Find evidence in the project before proceeding.**

Run these checks IN ORDER using Grep and Read tools:

### Check 1: Look for `replace` directive in go.mod

```bash
# Command to run:
grep "replace github.com/hajimehoshi/ebiten/v2" go.mod
```

**If FOUND → Patch Method confirmed**
- Evidence: `replace github.com/hajimehoshi/ebiten/v2 => /path/to/local/ebiten`
- The game uses a patched version of Ebiten
- No code changes needed in game

**If NOT FOUND → Continue to Check 2**

### Check 2: Look for autoebiten library calls in game code

```bash
# Command to run:
grep -r "autoebiten\.Update\|autoebiten\.Capture\|autoebiten\.IsKeyPressed" --include="*.go" .
```

**If FOUND → Library Method confirmed**
- Evidence: `autoebiten.Update()`, `autoebiten.Capture()`, or `autoebiten.IsKeyPressed()` calls
- Game directly uses the autoebiten library

**If NOT FOUND → No integration yet**
- Need to choose integration method based on whether game already exists

### Integration Method Summary

| Method | Detection | Code Changes |
|--------|-----------|--------------|
| **Patch** | `replace` directive in go.mod | None |
| **Library** | `autoebiten.Update()` calls in code | Required |

---

## Branch A: CLI Automation

If purpose is CLI automation (NOT writing tests):

### A1: Verify Game is Running

```bash
# Check connection
autoebiten ping

# If no game found, start it:
./mygame &
autoebiten ping
```

### A2: Send Inputs

**Keyboard:**
```bash
# Press a key once
autoebiten input --key KeySpace --action press

# Hold for 10 ticks (~167ms at 60 TPS)
autoebiten input --key KeyW --action hold --duration_ticks 10
```

**Mouse:**
```bash
# Move cursor
autoebiten mouse --action position -x 100 -y 200

# Click
autoebiten mouse --action press --button MouseButtonLeft
```

**Wheel:**
```bash
autoebiten wheel -y -3
```

### A3: Take Screenshot

```bash
autoebiten screenshot --output shot.png
```

### A4: Run Scripts

```bash
# From file
autoebiten run --script script.json

# Inline
autoebiten run --inline '{"version":"1.0","commands":[{"input":{"action":"press","key":"KeySpace"}}]}'
```

---

## Branch B: Writing Tests (testkit)

If purpose is writing tests:

### B1: Choose Testing Mode

| Mode | Process | Communication | Best For |
|------|---------|---------------|----------|
| Black-Box | Separate | RPC via socket | Full integration testing |
| White-Box | Same | Direct calls | Fast unit testing of game logic |

### B2: Black-Box Testing

Use `testkit.Launch()` to start game in separate process.

**Prerequisites Check:**

1. **Is binary built?**
   ```bash
   ls -la ./mygame
   ```
   If not found: `go build -o ./mygame ./cmd/mygame`

2. **Need StateQuery? Check for StateExporter:**
   ```bash
   grep -r "RegisterStateExporter" --include="*.go" .
   ```
   If not found, add to game code:
   ```go
   autoebiten.RegisterStateExporter("gamestate", &gameInstance)
   ```

**Example:**

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

func TestPlayerMovement(t *testing.T) {
    game := testkit.Launch(t, "./mygame",
        testkit.WithTimeout(30*time.Second))
    defer game.Shutdown()

    // Wait for game to be ready
    ready := game.WaitFor(func() bool {
        return game.Ping() == nil
    }, 5*time.Second)
    require.True(t, ready)

    // Get initial position
    initialX, err := game.StateQuery("gamestate", "Player.X")
    require.NoError(t, err)

    // Send input
    err = game.HoldKey(ebiten.KeyArrowRight, 10)
    require.NoError(t, err)

    // Verify position changed
    newX, err := game.StateQuery("gamestate", "Player.X")
    require.NoError(t, err)
    assert.Greater(t, newX.(float64), initialX.(float64))
}
```

### B3: White-Box Testing

Use `testkit.NewMock()` to test game logic in same process.

**Example:**

```go
package mygame

import (
    "testing"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/s3cy/autoebiten"
    "github.com/s3cy/autoebiten/testkit"
    "github.com/stretchr/testify/assert"
)

func TestPlayerMovesRight(t *testing.T) {
    autoebiten.SetMode(autoebiten.InjectionOnly)

    game := NewGame()
    mock := testkit.NewMock(t, game)

    initialX := game.Player.X

    mock.InjectKeyPress(ebiten.KeyArrowRight)
    mock.Ticks(10)

    assert.Greater(t, game.Player.X, initialX)
}
```

---

## Common Mistakes

| Mistake | Evidence | Fix |
|---------|----------|-----|
| StateQuery fails | `grep RegisterStateExporter` returns nothing | Add `autoebiten.RegisterStateExporter("name", &game)` |
| Binary not found | `ls ./mygame` fails | Build: `go build -o ./mygame ./cmd/mygame` |
| Patch method panics | `replace` in go.mod AND `autoebiten.Update()` calls | Remove `autoebiten.Update()` calls (patch handles automatically) |

---

## Quick Reference

### Black-Box API

| Common Methods | Description |
|--------|-------------|
| `Launch(t, binary, opts...)` | Start game process |
| `game.Shutdown()` | Stop game |
| `game.Ping()` | Check responsive |
| `game.WaitFor(fn, timeout)` | Poll until true |
| `game.HoldKey(key, ticks)` | Hold key N ticks |
| `game.StateQuery(name, path)` | Query exported state |

### White-Box API

| Common Methods | Description |
|--------|-------------|
| `NewMock(t, game)` | Create mock |
| `mock.InjectKeyPress(key)` | Buffer press |
| `mock.Tick()` / `mock.Ticks(n)` | Advance |

### Key Constants

Common: `KeyA`-`KeyZ`, `Key0`-`Key9`, `KeySpace`, `KeyEnter`, `KeyArrowUp`, `KeyArrowDown`, `KeyArrowLeft`, `KeyArrowRight`

List all: `autoebiten keys`

### Ticks vs FPS

1 tick = 1 `Update()` call. Ebiten runs at 60 TPS by default, independent of FPS.
