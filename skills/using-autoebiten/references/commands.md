# CLI Commands Reference

> Purpose: Complete reference for all autoebiten CLI commands
> Audience: Developers using the CLI for game automation

---

## Global Flags

These flags apply to all commands:

| Flag | Description |
|------|-------------|
| --pid, -p | Target game process PID (auto-detected if not specified) |

If multiple games are running and --pid is not specified, autoebiten will list available games and exit with an error.

---

## Game Control

### launch

Launch a game with output capture and crash diagnostics.

```bash
autoebiten launch -- ./game [args...]
```

**Flags:**
| Flag | Default | Description |
|------|---------|-------------|
| --timeout | 10s | Timeout waiting for game RPC server (e.g., 10s, 30s, 1m) |

**Examples:**
```bash
# Launch a game
autoebiten launch -- ./mygame

# Launch with game arguments
autoebiten launch -- ./mygame --level 1 --debug

# Launch with longer timeout
autoebiten launch --timeout 30s -- ./mygame
```

**Crash Diagnostics:**
When a game crashes or fails to start, the launch proxy captures error information and game output. CLI commands will show:
- Log diff since last command in `<log_diff>` tags
- Error details from the game exit in `<proxy_error>` tags

```bash
autoebiten launch -- ./mygame
autoebiten ping
# If game crashed, shows:
# <log_diff>
# --- snapshot ...
# +++ current ...
# @@ ... @@
# +panic: runtime error: index out of range
# </log_diff>
# <proxy_error>
# game exited: exit status 2
# </proxy_error>
# Error: game not connected

---

#### exit

Send exit signal to gracefully terminate the game.

```bash
autoebiten exit
```

**Examples:**
```bash
# Exit a game running directly
./mygame &
autoebiten exit

# Exit a game launched via launch
autoebiten launch -- ./mygame &
autoebiten exit
```

---

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
- `press` - Press key down (does NOT release - key stays pressed until `release` action)
- `release` - Release a held key
- `hold` - Press and hold for duration_ticks, then auto-release

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
- `press` - Press button down (does NOT release - button stays pressed until `release` action)
- `release` - Release button
- `hold` - Press and hold for duration_ticks, then auto-release (default when --button is set)

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
| --output, -o | Output file path relative to game's working directory (auto-generated if not set) |
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
