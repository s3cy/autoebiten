---
name: using-autoebiten
description: Automate Ebitengine games via CLI input injection, screenshots, and scripted sequences. Use when controlling games, sending inputs to games, testing Ebitengine games, or writing tests with the testkit package. Provides two integration methods (Patch for existing games, Library for new games) and two testing modes (Black-box via RPC, White-box via Mock). TRIGGER for any Ebitengine game automation, testing, or input injection task - even if user just mentions "game automation", "ebiten testing", or "control game programmatically".
---

# Using AutoEbiten

AutoEbiten automates Ebitengine games through input injection, screenshots, and scripted sequences.

## Decision Flow

**Step 1: Determine purpose**

| Purpose | Action |
|---------|------|
| CLI automation (send inputs, screenshots without tests) | **READ [references/cli.md](references/cli.md)** for full CLI guide |
| Writing Go tests with testkit | **READ [references/testkit.md](references/testkit.md)** for testing framework |

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

## When to Read References

**READ references/integration.md when:**
- Setting up Patch or Library integration
- Implementing custom commands
- Setting up State Exporter for state queries
- Understanding API functions (Update, Capture, input wrappers)

**READ references/commands.md when:**
- Using CLI commands beyond basic input/screenshot
- Looking up command flags and options
- Debugging connection issues
- Writing automation scripts

**READ references/tutorial.md when:**
- Following step-by-step integration tutorial
- Looking for practical examples
- Understanding crash diagnostics

**READ references/autoui.md when:**
- Automating EbitenUI widget trees
- Querying widgets by coordinates, attributes, or XPath
- Invoking widget methods (Click, SetText) without hard-coded coordinates
- Writing tests that interact with UI widgets

**READ references/testkit.md when:**
- Choosing between Black-box (Launch) vs White-box (Mock)
- Writing test assertions with StateQuery
- Setting up State Exporter for white-box tests
- Understanding ticks vs time
- Writing game logic unit tests

---

## Quick Start (Minimal)

### CLI Commands

```bash
autoebiten ping                     # Check connection
autoebiten input --key KeySpace     # Press key
autoebiten screenshot               # Capture screen
```

### Testkit (Black-box)

```go
game := testkit.Launch(t, "./mygame")
defer game.Shutdown()
game.HoldKey(ebiten.KeyArrowRight, 10)
```

### Ticks vs Time

1 tick = 1 `Update()` call. Ebiten runs at 60 TPS by default.
- 6 ticks ≈ 100ms
- 60 ticks ≈ 1 second

---

## Common Mistakes

| Issue | Check | Fix |
|-------|-------|-----|
| Key stuck pressed after `press` action | Using `--action press` alone | `press` only presses DOWN. It does NOT release. Use `hold` for press+hold+release, or manually call `press` then `release` |
| CLI command fails with connection error | `autoebiten ping` fails | Game likely crashed. Restart game: `./mygame &`. Check game logs for crash cause |
| StateQuery fails | `grep RegisterStateExporter` returns nothing | Add `autoebiten.RegisterStateExporter("name", &game)` |
| StateQuery returns empty/missing | Field is lowercase (unexported) | Capitalize field name. See **State Exporter** section in reference docs for full rules |
| Binary not found | `ls ./mygame` fails | Build: `go build -o ./mygame ./cmd/mygame` |
| Patch panics | `replace` in go.mod AND `autoebiten.Update()` calls | Remove `autoebiten.Update()` (patch handles it) |

---

## Reference Documents

For full documentation, **READ the appropriate reference file** based on your task:

- **[Integration Guide](references/integration.md)** - Patch vs Library integration, API reference, custom commands, state exporter setup
- **[CLI Commands](references/commands.md)** - Complete CLI command reference with all flags and options
- **[Tutorial & Examples](references/tutorial.md)** - Step-by-step tutorial and practical examples
- **[autoui Reference](references/autoui.md)** - EbitenUI widget automation, XPath queries, widget method invocation
- **[Testkit Reference](references/testkit.md)** - Black-box vs White-box testing, test API, state queries, testing patterns

---

## Examples Repository

The full project with runnable examples is at **github.com/s3cy/autoebiten**.

To browse examples:
```bash
git clone https://github.com/s3cy/autoebiten
cd autoebiten
ls examples/ examples/testkit/
```