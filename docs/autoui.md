# autoui Package Documentation

EbitenUI automation helper for autoebiten. Enables CLI-based widget tree inspection, search, method invocation, and visual debugging for LLM-assisted E2E testing.

## Overview

`autoui` bridges the gap between autoebiten's RPC-based automation and ebitenui's widget hierarchy. It provides commands for:

- **Tree inspection**: Export full widget hierarchy as XML
- **Widget search**: Find widgets by attributes or XPath queries
- **Method invocation**: Call widget methods via reflection
- **Visual debugging**: Highlight widgets on screen

## Installation

```bash
go get github.com/s3cy/autoebiten
```

## Quick Start

```go
package main

import (
    "github.com/ebitenui/ebitenui"
    "github.com/ebitenui/ebitenui/widget"
    "github.com/s3cy/autoebiten/autoui"
)

func main() {
    // Create your UI
    root := widget.NewContainer(
        widget.ContainerOpts.Layout(widget.NewRowLayout()),
    )

    btn := widget.NewButton(
        widget.ButtonOpts.Text("Start Game", face, colors),
        widget.WidgetOpts.CustomData(struct {
            ID string `ae:"id"`
        }{ID: "start-btn"}),
    )
    root.AddChild(btn)

    ui := ebitenui.UI{Container: root}

    // Register autoui commands
    autoui.Register(&ui)

    // Run your game...
}
```

## Commands

All commands are accessed via `autoebiten custom <command>`:

### autoui.tree

Returns the full widget tree as XML.

```bash
autoebiten custom autoui.tree
```

**Output:**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<UI>
  <Container _addr="0x14000123450" x="0" y="0" width="800" height="600" visible="true" disabled="false">
    <Button _addr="0x14000123480" x="100" y="50" width="200" height="40" visible="true" disabled="false" text="Start Game" state="idle" id="start-btn"/>
  </Container>
</UI>
```

### autoui.find

Find widgets by simple attribute matching.

```bash
# Find by ID (uses ae tag from CustomData)
autoebiten custom autoui.find --request "id=start-btn"

# Find all buttons
autoebiten custom autoui.find --request "type=Button"

# Find visible buttons
autoebiten custom autoui.find --request "visible=true"

# JSON for multiple criteria (AND logic)
autoebiten custom autoui.find --request '{"type":"Button","state":"idle"}'
```

**Output format:** `autoui.find` returns matched widgets directly without `<UI>` wrapper or hierarchy reconstruction. For full hierarchy, use `autoui.tree`.

```xml
<Button _addr="0x14000123480" x="100" y="50" width="200" height="40" visible="true" disabled="false" text="Start Game" state="idle" id="start-btn"/>
<Button _addr="0x14000123500" x="100" y="200" width="200" height="40" visible="true" disabled="false" text="Cancel" state="idle" id="cancel-btn"/>
```

### autoui.xpath

Find widgets using XPath 1.0 expressions.

```bash
# Find by ID
autoebiten custom autoui.xpath --request "//Button[@id='start-btn']"

# Find all visible buttons
autoebiten custom autoui.xpath --request "//Button[@visible='true']"

# Find focused widget
autoebiten custom autoui.xpath --request "//*[@focused='true']"

# Find widget at position
autoebiten custom autoui.xpath --request "//*[@x<100 and @y<100]"
```

**Output format:** Same flat format as `autoui.find` - matched widgets directly without `<UI>` wrapper.

### autoui.at

Get widget at specific screen coordinates.

```bash
autoebiten custom autoui.at --request "100,200"
# or
autoebiten custom autoui.at --request '{"x":100,"y":200}'
```

### autoui.call

Invoke methods on widgets via reflection.

```bash
# Click a button
autoebiten custom autoui.call --request '{"target":"id=start-btn","method":"Click"}'

# Focus a text input
autoebiten custom autoui.call --request '{"target":"id=username","method":"Focus","args":[true]}'

# Set slider value
autoebiten custom autoui.call --request '{"target":"id=volume","method":"SetCurrentValue","args":[0.5]}'
```

### autoui.highlight

Visually highlight widgets for debugging.

```bash
# Highlight by ID
autoebiten custom autoui.highlight --request "id=start-btn"

# Highlight at coordinates
autoebiten custom autoui.highlight --request "100,200"

# Clear all highlights
autoebiten custom autoui.highlight --request "clear"
```

## XML Format

Widgets are serialized to XML with the following structure:

```xml
<WidgetType _addr="ptr" x="0" y="0" width="100" height="30" visible="true" disabled="false" .../>
```

### Standard Attributes

All widgets include these attributes:

| Attribute | Description |
|-----------|-------------|
| `_addr` | Pointer address for exact identification |
| `x`, `y` | Screen position (top-left corner) |
| `width`, `height` | Dimensions in pixels |
| `visible` | Visibility state (`true`/`false`) |
| `disabled` | Disabled state (`true`/`false`) |

### Widget-specific attributes

| Widget Type | Extra Attributes |
|-------------|------------------|
| Button | text (requires font), state, toggle, focused |
| TextInput | text, cursor, selection_start, selection_end, focused |
| Checkbox | checked, state |
| Slider | value, min, max |
| ProgressBar | value, min, max |
| Label | text |
| Container | (layout info if available) |

**Note:** Button `text` attribute requires font setup via `ButtonOpts.Text()`. For buttons without fonts, use `id` from CustomData for queries:

```go
btn.GetWidget().CustomData = map[string]string{"id": "submit-btn"}
```

```bash
autoebiten custom autoui.find --request "id=submit-btn"
```

## Custom Data (ae tags)

Attach custom attributes to widgets via `WidgetOpts.CustomData()` using the `ae` tag for custom naming:

### Struct with ae tags

```go
type WidgetMeta struct {
    ID      string `ae:"id"`
    Section string `ae:"section"`
}

widget.WidgetOpts.CustomData(WidgetMeta{
    ID:      "start-btn",
    Section: "main-menu",
})
```

**Output:**
```xml
<Button id="start-btn" section="main-menu" .../>
```

### Map[string]string

```go
widget.WidgetOpts.CustomData(map[string]string{
    "id": "btn1",
    "test-id": "start-button",
})
```

**Output:**
```xml
<Button id="btn1" test-id="start-button" .../>
```

### Nested structs

```go
type PlayerData struct {
    Name string `ae:"name"`
    Pos  struct {
        X int `ae:"x"`
        Y int `ae:"y"`
    } `ae:"pos"`
}
```

**Output:**
```xml
<Button name="Alice" pos.x="100" pos.y="200" .../>
```

### Slice flattening

```go
type PlayerData struct {
    Name string   `ae:"name"`
    Tags []string `ae:"tags"`
}

// In widget creation:
widget.WidgetOpts.CustomData(&PlayerData{
    Name: "Alice",
    Tags: []string{"fire", "ice", "wind"},
})
```

**Output:**
```xml
<Button name="Alice" tags.0="fire" tags.1="ice" tags.2="wind"/>
```

**XPath query:**
```bash
autoebiten custom autoui.xpath --request "//Button[@tags.0='fire']"
```

Nested slices use multi-level indexing:
```xml
<!-- [][]string{{"a","b"},{"c"}} -->
<data.0.0="a" data.0.1="b" data.1.0="c"/>
```

## Attribute Lookup Order

When searching with `autoui.find` or `autoui.xpath`, attributes are resolved in this order:

1. **Widget type** - `type` returns "Button", "Container", etc.
2. **Position** - `x`, `y`, `width`, `height` from widget rect
3. **State** - `visible`, `disabled`
4. **Widget-specific** - e.g., Button's `text`, `state`
5. **CustomData** - flattened attributes from `ae` tags or map

## Widget Identification (_addr)

Every widget includes an `_addr` attribute containing its pointer address:

```xml
<Button _addr="0x14000abc0" x="100" y="50" id="submit-btn"/>
```

**Note:** `_addr` is an internal attribute for exact widget identification. It changes between runs, so do not hard-code it in tests. Instead, read `_addr` dynamically from `autoui.tree` output if you need it for subsequent queries.

**Usage pattern:**
```bash
# Get widget tree, note the _addr
autoebiten custom autoui.tree

# Use _addr for exact query (if needed)
autoebiten custom autoui.xpath --request "//Button[@_addr='0x14000abc0']"
```

Most queries should use `id` or other stable attributes instead of `_addr`.

## Method Invocation

The `autoui.call` command uses reflection with a whitelist for safety:

**Supported signatures:**
- `func()` - No args
- `func(bool)` - Boolean
- `func(int)` - Integer
- `func(float64)` - Float
- `func(string)` - String

**Common methods:**

| Widget | Method | Args | Description |
|--------|--------|------|-------------|
| Button | `Click` | none | Simulate click |
| TextInput | `Focus` | `bool` | Set focus state |
| Slider | `SetCurrentValue` | `float64` | Set value |
| Checkbox | `Click` | none | Toggle state |

## Visual Debugging

Draw highlight rectangles in your game's Draw method:

```go
func (g *Game) Draw(screen *ebiten.Image) {
    g.ui.Draw(screen)
    autoui.DrawHighlights(screen) // Draw active highlights
}
```

Configure highlight duration:

```go
autoui.SetHighlightDuration(5 * time.Second) // Default: 3 seconds
```

## Testing

Use with testkit for E2E tests:

```go
func TestMainMenu(t *testing.T) {
    game := testkit.Launch(t, "./mygame")
    defer game.Shutdown()

    // Get UI tree
    tree := game.Custom("autoui.tree", "")
    t.Logf("UI Tree:\n%s", tree)

    // Find and click button
    result := game.Custom("autoui.find", "id=start-btn")
    require.Contains(t, result, "start-btn")

    game.Custom("autoui.call", `{"target":"id=start-btn","method":"Click"}`)
}
```

## API Reference

### Functions

```go
// Register all autoui commands
func Register(ui *ebitenui.UI)

// Register with custom prefix
func RegisterWithPrefix(ui *ebitenui.UI, prefix string)

// Draw highlight overlays (call in Draw)
func DrawHighlights(screen *ebiten.Image)

// Configure highlight duration
func SetHighlightDuration(d time.Duration)
```

## Examples

See [examples/autoui](../examples/autoui/main.go) for a complete working example.

## Troubleshooting

### Widget not found

- Check that `autoui.Register(&ui)` was called after UI construction
- Verify CustomData attributes are being extracted (check `autoui.tree` output)
- Ensure widget is visible (invisible widgets are still in tree but may be skipped by some queries)

### Method invocation fails

- Verify method name matches exactly (case-sensitive)
- Check that method signature is in the whitelist
- Ensure target widget exists before calling

### XPath queries not working

- Use single quotes for string values: `[@id='btn']` not `[@id="btn"]`
- Remember XPath is case-sensitive for element names
- Test query with `autoui.tree` output first
