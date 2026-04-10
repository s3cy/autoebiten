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
```

---

### exit

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

## Input Control

### input

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

Example output:
```text
Error: no running game found
Usage:
  autoebiten input [flags]

Flags:
  -a, --action string        Action: press, release, or hold (default "hold")
      --async                Async mode: return immediately without waiting for the input to be processed
  -d, --duration_ticks int   Duration in ticks for hold action (default 6)
  -h, --help                 help for input
  -k, --key string           Key name (e.g., KeyA, KeySpace, KeyArrowUp)
      --no-record            Skip recording this command

Global Flags:
  -p, --pid int   Target game process PID (auto-detected if not specified)

Error: no running game found
```

---

### mouse

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

Example output:
```text
Error: no running game found
Usage:
  autoebiten mouse [flags]

Flags:
  -a, --action string        Action: position, press, release, or hold (defaults to position, or hold when --button is used)
      --async                Async mode: return immediately without waiting for the input to be processed
  -b, --button string        Mouse button (e.g., MouseButtonLeft, MouseButtonRight)
  -d, --duration_ticks int   Duration in ticks for hold action (default 6)
  -h, --help                 help for mouse
      --no-record            Skip recording this command
  -x, --x int                X coordinate
  -y, --y int                Y coordinate

Global Flags:
  -p, --pid int   Target game process PID (auto-detected if not specified)

Error: no running game found
```

---

### wheel

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

Example output:
```text
Error: no running game found
Usage:
  autoebiten wheel [flags]

Flags:
      --async       Async mode: return immediately without waiting for the input to be processed
  -h, --help        help for wheel
      --no-record   Skip recording this command
  -x, --x float     Horizontal scroll (negative=left, positive=right)
  -y, --y float     Vertical scroll (negative=down, positive=up)

Global Flags:
  -p, --pid int   Target game process PID (auto-detected if not specified)

Error: no running game found
```

---

## Screenshot

### screenshot

Capture the game window.

```bash
autoebiten screenshot [--output <Path>] [--base64]
```

**Flags:**
| Flag | Description |
|------|-------------|
| --output, -o | Output file path (auto-generated if not set) |
| --base64 | Output as base64 string instead of file |
| --async, -a | Return immediately |
| --no-record | Skip recording |

**Examples:**
```bash
autoebiten screenshot
autoebiten screenshot --output capture.png
autoebiten screenshot --base64
```

Example output:
```text
Error: no running game found
Usage:
  autoebiten screenshot [flags]

Flags:
  -a, --async           Async mode: return immediately without waiting for capture
      --base64          Output screenshot as base64 string instead of saving to a file
  -h, --help            help for screenshot
      --no-record       Skip recording this command
  -o, --output string   Output file path (optional)

Global Flags:
  -p, --pid int   Target game process PID (auto-detected if not specified)

Error: no running game found
```

---

## Script Execution

### run

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

Example output:
```text
Error: no running game found
Usage:
  autoebiten run [flags]

Flags:
  -h, --help            help for run
      --inline string   Inline JSON script string
  -s, --script string   Path to script file

Global Flags:
  -p, --pid int   Target game process PID (auto-detected if not specified)

Error: no running game found
```

---

## Status and Info

### ping

Check if game is running and responsive.

```bash
autoebiten ping
```

Example output:
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

---

### version

Print CLI and game library versions.

```bash
autoebiten version
```

Example output:
```text
CLI version:    v0.7.1-0.20260409094100-ac392531f833
failed to connect to socket /tmp/autoebiten/autoebiten-48340.sock: dial unix /tmp/autoebiten/autoebiten-48340.sock: connect: no such file or directory
Game version:   not connected
```

---

### keys

List all available key names.

```bash
autoebiten keys
```

Example output:
```text
Error: no running game found
Usage:
  autoebiten keys [flags]

Flags:
  -h, --help   help for keys

Global Flags:
  -p, --pid int   Target game process PID (auto-detected if not specified)

Error: no running game found
```

---

### mouse_buttons

List all available mouse button names.

```bash
autoebiten mouse_buttons
```

Example output:
```text
Error: no running game found
Usage:
  autoebiten mouse_buttons [flags]

Flags:
  -h, --help   help for mouse_buttons

Global Flags:
  -p, --pid int   Target game process PID (auto-detected if not specified)

Error: no running game found
```

---

### get_mouse_position

Get injected mouse cursor position.

```bash
autoebiten get_mouse_position
```

Example output:
```text
Error: no running game found
Usage:
  autoebiten get_mouse_position [flags]

Flags:
  -h, --help   help for get_mouse_position

Global Flags:
  -p, --pid int   Target game process PID (auto-detected if not specified)

Error: no running game found
```

---

### get_wheel_position

Get injected wheel position.

```bash
autoebiten get_wheel_position
```

Example output:
```text
Error: no running game found
Usage:
  autoebiten get_wheel_position [flags]

Flags:
  -h, --help   help for get_wheel_position

Global Flags:
  -p, --pid int   Target game process PID (auto-detected if not specified)

Error: no running game found
```

---

### schema

Output JSON Schema for script files.

```bash
autoebiten schema > autoebiten-schema.json
```

Example output:
```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://github.com/s3cy/autoebiten/internal/script/script-schema",
  "$ref": "#/$defs/ScriptSchema",
  "$defs": {
    "CommandSchema": {
      "oneOf": [
        {
          "required": [
            "input"
          ],
          "title": "input"
        },
        {
          "required": [
            "mouse"
          ],
          "title": "mouse"
        },
        {
          "required": [
            "wheel"
          ],
          "title": "wheel"
        },
        {
          "required": [
            "screenshot"
          ],
          "title": "screenshot"
        },
        {
          "required": [
            "delay"
          ],
          "title": "delay"
        },
        {
          "required": [
            "custom"
          ],
          "title": "custom"
        },
        {
          "required": [
            "state"
          ],
          "title": "state"
        },
        {
          "required": [
            "wait"
          ],
          "title": "wait"
        },
        {
          "required": [
            "repeat"
          ],
          "title": "repeat"
        }
      ],
      "properties": {
        "input": {
          "$ref": "#/$defs/InputCmd",
          "description": "Inject keyboard input"
        },
        "mouse": {
          "$ref": "#/$defs/MouseCmd",
          "description": "Inject mouse input"
        },
        "wheel": {
          "$ref": "#/$defs/WheelCmd",
          "description": "Inject wheel/scroll input"
        },
        "screenshot": {
          "$ref": "#/$defs/ScreenshotCmd",
          "description": "Capture game screenshot"
        },
        "delay": {
          "$ref": "#/$defs/DelayCmd",
          "description": "Pause execution for a duration"
        },
        "custom": {
          "$ref": "#/$defs/CustomCmd",
          "description": "Execute a registered custom command"
        },
        "state": {
          "$ref": "#/$defs/StateCmd",
          "description": "Query game state via registered exporter"
        },
        "wait": {
          "$ref": "#/$defs/WaitCmd",
          "description": "Wait for condition to be met"
        },
        "repeat": {
          "$ref": "#/$defs/RepeatSchema",
          "description": "Repeat a block of commands"
        }
      },
      "additionalProperties": false,
      "type": "object",
      "title": "Command",
      "description": "A single command to execute"
    },
    "CustomCmd": {
      "properties": {
        "name": {
          "type": "string",
          "description": "Name of the registered custom command"
        },
        "request": {
          "type": "string",
          "description": "Optional request data to pass to the custom command"
        }
      },
      "additionalProperties": false,
      "type": "object"
    },
    "DelayCmd": {
      "properties": {
        "ms": {
          "type": "integer",
          "minimum": 0,
          "description": "Milliseconds to wait"
        }
      },
      "additionalProperties": false,
      "type": "object"
    },
    "InputCmd": {
      "properties": {
        "action": {
          "type": "string",
          "enum": [
            "press",
            "release",
            "hold"
          ],
          "description": "Action to perform"
        },
        "key": {
          "type": "string",
          "description": "Key name (use 'autoebiten keys' to list all)"
        },
        "duration_ticks": {
          "type": "integer",
          "description": "Duration in game ticks for hold action",
          "default": 6
        },
        "async": {
          "type": "boolean",
          "description": "Return immediately without waiting",
          "default": false
        }
      },
      "additionalProperties": false,
      "type": "object"
    },
    "MouseCmd": {
      "properties": {
        "action": {
          "type": "string",
          "enum": [
            "position",
            "press",
            "release",
            "hold"
          ],
          "description": "Action to perform (default is position or hold when button is used)"
        },
        "x": {
          "type": "integer",
          "description": "X coordinate",
          "default": 0
        },
        "y": {
          "type": "integer",
          "description": "Y coordinate",
          "default": 0
        },
        "button": {
          "type": "string",
          "description": "Mouse button (use 'autoebiten mouse_buttons' to list all)"
        },
        "duration_ticks": {
          "type": "integer",
          "description": "Duration in game ticks for hold action",
          "default": 6
        },
        "async": {
          "type": "boolean",
          "description": "Return immediately without waiting",
          "default": false
        }
      },
      "additionalProperties": false,
      "type": "object"
    },
    "RepeatSchema": {
      "properties": {
        "times": {
          "type": "integer",
          "minimum": 1,
          "description": "Number of times to repeat"
        },
        "commands": {
          "items": {
            "$ref": "#/$defs/CommandSchema"
          },
          "type": "array",
          "description": "Commands to repeat"
        }
      },
      "additionalProperties": false,
      "type": "object"
    },
    "ScreenshotCmd": {
      "properties": {
        "output": {
          "type": "string",
          "description": "Output file path (optional"
        },
        "base64": {
          "type": "boolean",
          "description": "Return screenshot as base64 string in the response instead of saving to a file",
          "default": false
        },
        "async": {
          "type": "boolean",
          "description": "Return immediately without waiting",
          "default": false
        }
      },
      "additionalProperties": false,
      "type": "object"
    },
    "ScriptSchema": {
      "properties": {
        "$schema": {
          "type": "string",
          "format": "uri-reference",
          "description": "JSON Schema URI for this document"
        },
        "version": {
          "type": "string",
          "enum": [
            "1.0"
          ],
          "description": "Script format version"
        },
        "commands": {
          "items": {
            "$ref": "#/$defs/CommandSchema"
          },
          "type": "array",
          "description": "List of commands to execute in order"
        }
      },
      "additionalProperties": false,
      "type": "object"
    },
    "StateCmd": {
      "properties": {
        "name": {
          "type": "string",
          "description": "State exporter name"
        },
        "path": {
          "type": "string",
          "description": "Dot-notation path to query"
        }
      },
      "additionalProperties": false,
      "type": "object",
      "required": [
        "name",
        "path"
      ]
    },
    "WaitCmd": {
      "properties": {
        "condition": {
          "type": "string",
          "description": "Condition to poll for (e.g."
        },
        "timeout": {
          "type": "string",
          "description": "Maximum wait duration (e.g."
        },
        "interval": {
          "type": "string",
          "description": "Poll interval (default 100ms)"
        },
        "verbose": {
          "type": "boolean",
          "description": "Print errors during polling"
        }
      },
      "additionalProperties": false,
      "type": "object",
      "required": [
        "condition",
        "timeout"
      ]
    },
    "WheelCmd": {
      "properties": {
        "x": {
          "type": "number",
          "description": "Horizontal scroll",
          "default": 0
        },
        "y": {
          "type": "number",
          "description": "Vertical scroll",
          "default": 0
        },
        "async": {
          "type": "boolean",
          "description": "Return immediately without waiting",
          "default": false
        }
      },
      "additionalProperties": false,
      "type": "object"
    }
  },
  "title": "AutoEbiten Script",
  "description": "JSON script format for automating Ebitengine games"
}
```

---

## Custom Commands

### list_custom

List registered custom commands.

```bash
autoebiten list_custom
```

Example output:
```text
Error: no running game found
Usage:
  autoebiten list_custom [flags]

Flags:
  -h, --help   help for list_custom

Global Flags:
  -p, --pid int   Target game process PID (auto-detected if not specified)

Error: no running game found
```

---

### custom

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

Example output for getPlayerInfo:
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

Example output for echo:
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

## State Queries

### state

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

### wait-for

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

## Recording

### clear_recording

Clear the recording file for current game.

```bash
autoebiten clear_recording
```

Example output:
```text
Error: no running game found
Usage:
  autoebiten clear_recording [flags]

Flags:
  -h, --help   help for clear_recording

Global Flags:
  -p, --pid int   Target game process PID (auto-detected if not specified)

Error: no running game found
```

---

### replay

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

Example output for dump:
```json
Error: no running game found
Usage:
  autoebiten replay [flags]

Flags:
  -d, --dump string   Dump script to file instead of executing
  -h, --help          help for replay
  -s, --speed float   Speed multiplier (default 1.0) (default 1)

Global Flags:
  -p, --pid int   Target game process PID (auto-detected if not specified)

Error: no running game found
```
