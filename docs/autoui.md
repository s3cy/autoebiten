# autoui Reference

> Purpose: EbitenUI widget automation for CLI, tests, and LLM agents
> Audience: CLI users, test writers, LLM agents automating EbitenUI games

---

## Quick Decision

**Querying widgets:**
```
├─ Need widget at coordinates? → autoui.at command
├─ Need widgets by attribute? → autoui.find command
├─ Need complex query? → autoui.xpath command
├─ Need to check existence? → autoui.exists command (returns JSON)
└─ Need full tree? → autoui.tree command
```

**Acting on widgets:**
```
├─ Need to click/interact? → autoui.call command
├─ Need to set text? → autoui.call SetText method
└─ Need visual debugging? → autoui.highlight command
```

---

## Overview

autoui provides EbitenUI automation via CLI commands:

- Widget tree inspection → XML export
- Widget search → coordinates, attributes, XPath
- Method invocation → reflection-based calls
- Visual debugging → highlight rectangles

**Key concepts:**

1. **WidgetInfo**: Internal representation (type, rect, state, customData)
2. **XML output**: Widget tree as XML (type as element name)
3. **ae tags**: Custom attributes via struct field tags

**Integration:**
```go
ui := ebitenui.UI{Container: root}
autoui.Register(&ui)  // Registers autoui.* commands
```

---

## Commands

### autoui.tree

Export full widget tree as XML.

**Usage:**
```bash
autoebiten custom autoui.tree
```

**Output:**
```xml
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

### autoui.find

Find widgets by attribute (AND logic for multiple criteria).

**Usage:**
```bash
autoebiten custom autoui.find --request "type=Button"
```

**Output:**
```xml
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

### autoui.exists

Check if widgets matching a query exist. Returns JSON for use with `wait-for`.

**Usage:**
```bash
autoebiten custom autoui.exists --request "type=Button"
```

**Output:**
```json
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

### autoui.call

Invoke method on widget.

**Usage:**
```bash
autoebiten custom autoui.call --request '{"target":"id=submit-btn","method":"Click","args":[]}'
```

**Output:**
```json
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
