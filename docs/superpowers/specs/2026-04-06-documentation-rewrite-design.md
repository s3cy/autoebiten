# Documentation Rewrite Design

> Purpose: Create LLM-optimized documentation for autoebiten CLI and testkit
> Date: 2026-04-06

---

## Overview

Rewrite autoebiten documentation with two comprehensive, self-contained files optimized for both LLMs and human readers. Each file combines API reference, tutorial, and examples in one place.

---

## Critical Constraint

**DO NOT rely on existing README.md, SPEC.md, or SKILL.md content. They are outdated.**

When writing documentation:
- Reference actual source code via `Read` and `Grep` tools
- Extract API signatures from Go source files (`*.go`)
- Verify function behavior by reading implementation code
- Cross-reference with test files for usage examples
- Only use existing docs for structure ideas, NOT for content accuracy

**Source of truth:** Go source code in the repository, NOT markdown files.

---

## Directory Layout

```
/Users/s3cy/Desktop/go/autoebiten/
├── README.md                    (entry point - stays at root, update navigation)
├── docs/
│   ├── SPEC.md                  (moved from root, update)
│   ├── cli.md                   (NEW: CLI API + Tutorial + Examples)
│   └── testkit.md               (NEW: Testkit API + Tutorial + Examples)
│
├── skills/using-autoebiten/
│   └── SKILL.md                 (NOT moved - stays in skills directory)
│
└── examples/                    (code examples)
    ├── simple/main.go           (existing)
    ├── custom_commands/main.go  (existing)
    ├── scripts/                 (NEW: JSON script examples)
    ├── state_exporter/          (NEW: state exporter demo)
    └── testkit/                 (NEW: test examples)
        ├── black_box_test.go
        └── white_box_test.go
```

---

## File Structure Pattern

Each documentation file (`cli.md`, `testkit.md`) follows this LLM-optimized template:

```markdown
# [CLI/Testkit] Reference

> Purpose: [Single sentence describing scope]
> Audience: [Who should read this]

---

## Quick Decision

[Text-based decision tree]

**If you want [X]:** → Use [Approach A]
**If you want [Y]:** → Use [Approach B]

---

## API Reference

[All public functions, alphabetically organized]

### FunctionName

**Signature:** `func Name(params) returns`

**Purpose:** [One line]

**Parameters:**
| Name | Type | Description |
|------|------|-------------|

**Returns:** [What each return means]

**Example:**
```go
[Code block]
```

**Notes:** [Edge cases, common mistakes, related functions]

---

## Tutorial

### Step 1: [Title]

**Goal:** [What you'll accomplish]

**Prerequisites:** [What you need]

**Action:**
```bash
[Commands]
```

**Expected result:** [What you should see]

**Troubleshooting:** [Common issues]

---

## Examples

### [Example Name]

**Scenario:** [What problem this solves]

**Code:**
```go
[Code with inline comments]
```

**How to run:**
```bash
[Commands]
```

**Key concepts:** [What this demonstrates]
```

---

## cli.md Content (~400-500 lines)

### Quick Decision Section

**Integration method:**
- Already using `ebiten.IsKeyPressed()` → Patch method (no code changes)
- Writing a new game → Library method (direct API)

**Input mode:**
- Want only CLI inputs → InjectionOnly
- Want CLI + real input → InjectionFallback (default)
- Disable CLI control → Passthrough

### API Reference Section

**Core Functions:**
- `Update()` - Process RPC commands
- `Capture()` - Screenshot capture
- `SetMode()` - Configure input mode

**Input Functions:**
- `IsKeyPressed(key Key) bool`
- `IsMouseButtonPressed(button MouseButton) bool`
- `CursorPosition() (x, y int)`
- `Wheel() (x, y float64)`
- `IsKeyJustPressed(key Key) bool`
- `IsKeyJustReleased(key Key) bool`
- `KeyPressDuration(key Key) int`
- `IsMouseButtonJustPressed(button MouseButton) bool`
- `IsMouseButtonJustReleased(button MouseButton) bool`
- `MouseButtonPressDuration(button MouseButton) int`

**Custom Commands:**
- `Register(name string, handler func(CommandContext))`
- `Unregister(name string) bool`
- `ListCustomCommands() []string`

**State Exporter:**
- `RegisterStateExporter(name string, exporter any)`

**Build Tags:**
- Default (dev): RPC server enabled
- `-tags release`: RPC server disabled

### Tutorial Section

1. Choose Integration Method (Patch vs Library decision flow)
2. Patch Method Setup (clone Ebiten, apply patch, go.mod replace)
3. Library Method Setup (import, modify Update/Draw, build modes)
4. Verify Connection (ping, understanding socket paths)
5. Basic CLI Commands (input, mouse, wheel, screenshot)
6. Understanding Ticks (TPS vs FPS, duration_ticks parameter)
7. Custom Commands Intro (Register, CommandContext)

### Examples Section

1. **Library Integration** - Document existing `examples/simple/main.go`
2. **Patch Method Setup** - New walkthrough with sample game
3. **Custom Commands** - Document existing `examples/custom_commands/main.go`
4. **Scripted Automation** - New JSON scripts in `examples/scripts/`
5. **State Exporter** - New example in `examples/state_exporter/`
6. **Complex Input Sequences** - Script combining keyboard + mouse + timing
7. **Multiple Instances** - PID targeting, AUTOEBITEN_SOCKET

---

## testkit.md Content (~300-400 lines)

### Quick Decision Section

**Testing mode:**
- Need full integration (real game process) → Black-box (`Launch`)
- Need fast unit tests (game logic only) → White-box (`NewMock`)
- Need to verify internal state → Add `RegisterStateExporter` to game

### API Reference Section

**Black-Box (Game type):**
- `Launch(t *testing.T, binaryPath string, opts ...Option) *Game`
- `Shutdown()`
- `Ping() error`
- `WaitFor(fn func() bool, timeout time.Duration) bool`
- `PressKey(key ebiten.Key) error`
- `HoldKey(key ebiten.Key, ticks int) error`
- `MoveMouse(x, y int) error`
- `ClickMouse(button ebiten.MouseButton) error`
- `HoldMouse(button ebiten.MouseButton, ticks int) error`
- `Screenshot() (*ebiten.Image, error)`
- `ScreenshotToFile(path string) error`
- `ScreenshotBase64() (string, error)`
- `StateQuery(name string, path string) (any, error)`
- `RunCustom(name string, request string) (string, error)`
- `InjectWheel(x, y float64) error`

**White-Box (Mock type):**
- `NewMock(t *testing.T, game GameUpdate) *Mock`
- `InjectKeyPress(key ebiten.Key)`
- `InjectKeyHold(key ebiten.Key, ticks int)`
- `InjectMousePosition(x, y int)`
- `InjectMouseClick(button ebiten.MouseButton)`
- `InjectMouseHold(button ebiten.MouseButton, ticks int)`
- `InjectWheel(x, y float64)`
- `Tick()`
- `Ticks(n int)`

**Options:**
- `WithTimeout(d time.Duration) Option`
- `WithSocketPath(path string) Option`

**GameUpdate interface:** (define what game must implement for Mock)

### Tutorial Section

1. Choose Testing Mode (Black-box vs White-box decision)
2. Black-Box Setup (build binary, Launch, lifecycle)
3. White-Box Setup (NewMock, in-process testing)
4. Adding State Exporter (RegisterStateExporter for StateQuery)
5. Writing Assertions (testify patterns, WaitFor polling)

### Examples Section

1. **Black-Box: Player Movement Test** - Launch, HoldKey, StateQuery
2. **White-Box: Game Logic Test** - NewMock, InjectKeyPress, Ticks
3. **State Query: Health Verification** - RegisterStateExporter, dot-notation paths

---

## README.md Updates

Add navigation section at the top:

```markdown
## Documentation

### CLI Automation
Control games via command line → **[docs/cli.md](docs/cli.md)**

### Game Testing
Write automated tests → **[docs/testkit.md](docs/testkit.md)**

### Technical Specs
Architecture & protocol → **[docs/SPEC.md](docs/SPEC.md)**
```

Keep existing Quick Start and Installation sections. Link to detailed docs instead of duplicating content.

---

## Files to Create

| File | Lines | Source |
|------|-------|--------|
| `docs/cli.md` | ~400-500 | Extract from Go source files |
| `docs/testkit.md` | ~300-400 | Extract from Go source files |

## Files to Move

| From | To |
|------|-----|
| `SPEC.md` | `docs/SPEC.md` |

## Files to Update

| File | Change |
|------|--------|
| `README.md` | Add documentation navigation section |

## New Example Code

| Directory | Files | Purpose |
|-----------|-------|---------|
| `examples/scripts/` | `basic.json`, `mouse.json`, `complex.json` | JSON script examples |
| `examples/state_exporter/` | `main.go` | StateQuery demo game |
| `examples/testkit/` | `black_box_test.go`, `white_box_test.go` | Testkit usage examples |

---

## Implementation Order

1. Move `SPEC.md` to `docs/SPEC.md` then update
2. Create `docs/cli.md` (extract API from source code)
3. Create `docs/testkit.md` (extract API from source code)
4. Create example code (`scripts/`, `state_exporter/`, `testkit/`)
5. Update `README.md` navigation

---

## LLM Optimizations Applied

| Optimization | Rationale |
|--------------|-----------|
| Flat structure (2 main docs) | Complete context in one read |
| Text-based decision trees | Reliable conditional parsing |
| Alphabetical function ordering | Predictable location pattern |
| Tables for parameters | Structured data extraction |
| Self-contained sections | No cross-file dependencies |
| Code before explanation | Concrete examples first |
| Explicit "Notes" field | Edge cases visibility |