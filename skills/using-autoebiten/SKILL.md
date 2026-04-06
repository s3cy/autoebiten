---
name: using-autoebiten
description: Automate Ebitengine games via CLI input injection, screenshots, and scripted sequences. Use when controlling games, sending inputs to games, testing Ebitengine games, or writing tests with the testkit package. Provides two integration methods (Patch for existing games, Library for new games) and two testing modes (Black-box via RPC, White-box via Mock).
---

# Using AutoEbiten

AutoEbiten automates Ebitengine games through input injection, screenshots, and scripted sequences.

## Decision Flow

**Step 1: Determine purpose**

| Purpose | Tool |
|---------|------|
| CLI automation (send inputs, screenshots without tests) | CLI → See [references/cli.md](references/cli.md) |
| Writing Go tests | testkit → See [references/testkit.md](references/testkit.md) |

**Step 2: Determine integration method** (detect before proceeding)

Run checks in order:

1. Check for `replace github.com/hajimehoshi/ebiten/v2` in go.mod
   - **Found** → Patch method (no code changes needed)
2. Check for `autoebiten.Update()` or `autoebiten.IsKeyPressed()` in code
   - **Found** → Library method (already integrated)
3. **Neither found** → Choose based on whether game exists
   - Existing game with `ebiten.IsKeyPressed()` → Patch method
   - New game or willing to modify → Library method

| Method | Detection | Code Changes |
|--------|-----------|--------------|
| Patch | `replace` directive in go.mod | None (patched Ebiten) |
| Library | `autoebiten.Update()` calls | Required |

---

## Reference Documents

- **[CLI Reference](references/cli.md)** - Complete CLI command guide, integration setup, tutorials, examples
- **[Testkit Reference](references/testkit.md)** - Go testing framework, Black-box vs White-box modes, API reference

---

## Finding Examples

The full project with examples is at **github.com/s3cy/autoebiten**.

Examples locations in the repository:
- `examples/` - Integration examples (Patch and Library methods)
- `examples/testkit` - Testing examples (Black-box and White-box)

To browse examples:
```bash
git clone https://github.com/s3cy/autoebiten
cd autoebiten
ls examples/ examples/testkit/
```

---

## Quick Reference

### CLI Commands

```bash
autoebiten ping                     # Check connection
autoebiten input --key KeySpace     # Press key
autoebiten mouse -x 100 -y 200      # Move cursor
autoebiten screenshot               # Capture screen
autoebiten run --script script.json # Run automation
autoebiten keys                     # List key names
```

### Testkit Quick Start

**Black-box (separate process):**
```go
game := testkit.Launch(t, "./mygame")
defer game.Shutdown()
game.HoldKey(ebiten.KeyArrowRight, 10)
```

**White-box (in-process):**
```go
mock := testkit.NewMock(t, game)
mock.InjectKeyPress(ebiten.KeyArrowRight)
mock.Ticks(10)
```

### Ticks vs Time

1 tick = 1 `Update()` call. Ebiten runs at 60 TPS by default.
- 6 ticks ≈ 100ms
- 60 ticks ≈ 1 second

---

## Common Mistakes

| Issue | Check | Fix |
|-------|-------|-----|
| StateQuery fails | `grep RegisterStateExporter` returns nothing | Add `autoebiten.RegisterStateExporter("name", &game)` |
| Binary not found | `ls ./mygame` fails | Build: `go build -o ./mygame ./cmd/mygame` |
| Patch panics | `replace` in go.mod AND `autoebiten.Update()` calls | Remove `autoebiten.Update()` (patch handles it) |