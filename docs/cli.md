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