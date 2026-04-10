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

**Runnable Example:**

```bash
# Build and run the demo
cd examples/autoui
go build -o autoui_demo
autoebiten launch -- ./autoui_demo &
```


```bash
# Get widget tree
autoebiten custom autoui.tree
```

**Output:**
```xml

OK: <UI>
  <Container _addr="<ADDR>" disabled="false" height="480" visible="true" width="640" x="0" y="0">
    <Button _addr="<ADDR>" disabled="false" height="40" id="submit-btn" role="primary" state="unchecked" visible="true" width="200" x="100" y="50"></Button>
    <Button _addr="<ADDR>" disabled="false" height="40" id="cancel-btn" role="secondary" state="unchecked" visible="true" width="200" x="100" y="200"></Button>
  </Container>
</UI>

```

---

### autoui.at

Find widget at coordinates.

**Usage:**
```bash
autoebiten custom autoui.at --request "100,200"
autoebiten custom autoui.at --request '{"x":100,"y":200}'
```

**Runnable Example:**

```bash
# Build and run the demo
cd examples/autoui
go build -o autoui_demo
autoebiten launch -- ./autoui_demo &
```

```bash
# Find widget at position (150, 70) - the Submit button
autoebiten custom autoui.at --request "150,70"
```

**Output:**
```xml

OK: <Button _addr="<ADDR>" disabled="false" height="40" id="submit-btn" role="primary" state="unchecked" visible="true" width="200" x="100" y="50"></Button>

```

**Error:** `no widget found at coordinates`

---

### autoui.find

Find widgets by attribute (AND logic for multiple criteria).

**Usage:**
```bash
autoebiten custom autoui.find --request "type=Button"
autoebiten custom autoui.find --request '{"type":"Button","text":"Submit"}'
```

**Runnable Example:**

```bash
# Build and run the demo
cd examples/autoui
go build -o autoui_demo
autoebiten launch -- ./autoui_demo &
```

```bash
# Find all buttons
autoebiten custom autoui.find --request "type=Button"
```

**Output:**
```xml

OK: <Button _addr="<ADDR>" disabled="false" height="40" id="submit-btn" role="primary" state="unchecked" visible="true" width="200" x="100" y="50"></Button><Button _addr="<ADDR>" disabled="false" height="40" id="cancel-btn" role="secondary" state="unchecked" visible="true" width="200" x="100" y="200"></Button>

```

```bash
# Find submit button specifically
autoebiten custom autoui.find --request "id=submit-btn"
```

**Output:**
```xml

OK: <Button _addr="<ADDR>" disabled="false" height="40" id="submit-btn" role="primary" state="unchecked" visible="true" width="200" x="100" y="50"></Button>

```

**Note:** `autoui.find` returns matched widgets directly without `<UI>` wrapper or hierarchy reconstruction. For full hierarchy, use `autoui.tree`.

**Error:** `no widgets found matching query`

---

### autoui.xpath

XPath 1.0 query on widget tree.

**Usage:**
```bash
autoebiten custom autoui.xpath --request "//Button[@id='submit-btn']"
autoebiten custom autoui.xpath --request "//Button[contains(@id,'submit')]"
```

**Runnable Example:**

```bash
# Build and run the demo
cd examples/autoui
go build -o autoui_demo
autoebiten launch -- ./autoui_demo &
```

```bash
# Find all buttons
autoebiten custom autoui.xpath --request "//Button"
```

**Output:**
```xml

OK: <Button _addr="<ADDR>" disabled="false" height="40" id="submit-btn" role="primary" state="unchecked" visible="true" width="200" x="100" y="50"></Button><Button _addr="<ADDR>" disabled="false" height="40" id="cancel-btn" role="secondary" state="unchecked" visible="true" width="200" x="100" y="200"></Button>

```

**Note:** `autoui.xpath` returns matched widgets directly without `<UI>` wrapper. For full hierarchy, use `autoui.tree`.

**XPath examples:**
```
//Button                      → All Button widgets
//Button[@disabled='true']    → Disabled buttons
//Button[@id='submit-btn']    → Button with specific id
//*[contains(@id,'submit')]   → Any widget with 'submit' in id
```

**Error:** `no widgets found matching XPath`

---

### autoui.exists

Check if widgets matching a query exist. Returns JSON for use with `wait-for`.

**Usage:**
```bash
autoebiten custom autoui.exists --request "type=Dialog"
autoebiten custom autoui.exists --request '{"type":"TextInput","id":"name-input"}'
```

**Runnable Example:**

```bash
# Build and run the demo
cd examples/autoui
go build -o autoui_demo
autoebiten launch -- ./autoui_demo &
```

```bash
# Check if any Button widgets exist
autoebiten custom autoui.exists --request "type=Button"
```

**Output:**
```json

OK: {"found":true,"count":2}

```

**Key difference from autoui.find:** Returns JSON instead of XML. Never errors on empty results - useful with `wait-for`.

---

### autoui.call

Invoke method on widget.

**Usage:**
```bash
autoebiten custom autoui.call --request '{"target":"id=submit-btn","method":"Click","args":[]}'
```

**Runnable Example:**

```bash
# Build and run the demo
cd examples/autoui
go build -o autoui_demo
autoebiten launch -- ./autoui_demo &
```

```bash
# Click the Submit button
autoebiten custom autoui.call --request '{"target":"id=submit-btn","method":"Click","args":[]}'
```

**Output:**
```json

OK: {"success":true}

```

**Request format:**
```json
{
  "target": "type=Button",     // Query to find widget (key=value or JSON)
  "method": "Click",           // Method name to invoke
  "args": []                   // Arguments (optional)
}
```

**Whitelisted signatures:**
- `func()`
- `func(bool)`, `func(int)`, `func(float64)`, `func(string)`
- `func() error`, `func(bool) error`, etc.

**SetText Example:**

```bash
# Set text in a TextInput widget
autoebiten custom autoui.call --request '{"target":"id=name-input","method":"SetText","args":["Alice"]}'
```

**Output:**
```json

OK: error: no widget found matching target query

```

**Why use SetText:** TextInput's `SetText(string)` method fires `ChangedEvent` for game logic. Setting the field directly would bypass events.

**Error:** `{"success":false,"error":"method 'X' not found"}`

---

### autoui.highlight

Add visual highlight rectangles (red, 3-second duration).

**Usage:**
```bash
autoebiten custom autoui.highlight --request "clear"
autoebiten custom autoui.highlight --request "100,200"
autoebiten custom autoui.highlight --request "type=Button"
autoebiten custom autoui.highlight --request '{"query":"type=Button"}'
```

**Runnable Example:**

```bash
# Build and run the demo
cd examples/autoui
go build -o autoui_demo
autoebiten launch -- ./autoui_demo &
```

```bash
# Highlight all buttons
autoebiten custom autoui.highlight --request "type=Button"
```

**Output:**
```

OK: ok: highlighted 2 widgets

```

```bash
# Take screenshot to see the highlights
autoebiten screenshot --output highlighted.png
```

```bash
# Clear highlights
autoebiten custom autoui.highlight --request "clear"
```

**Output:**
```

OK: ok: highlights cleared

```

---

## Clicking Without Coordinates

Instead of mouse clicks with hard-coded x,y:

```bash
# OLD: Hard-coded coordinates (breaks if UI changes)
autoebiten mouse -x 150 -y 70 --button MouseButtonLeft
```

Use autoui to find and click by widget identity:

```bash
# NEW: Click by widget identity (resilient to layout changes)
autoebiten custom autoui.call --request '{"target":"id=submit-btn","method":"Click"}'
```

**Runnable Example:**

```bash
# Build and run the demo
cd examples/autoui
go build -o autoui_demo
autoebiten launch -- ./autoui_demo &
```

```bash
# Click Submit button without knowing its position
autoebiten custom autoui.call --request '{"target":"id=submit-btn","method":"Click"}'
```

**Output:**
```json

OK: {"success":true}

```

**Why this matters:**

1. Layout changes don't break tests
2. No screenshot measurement needed
3. Works across screen sizes
4. Same pattern for any widget method

**Workflow for LLM:**
1. Run `autoui.tree` to discover widgets
2. Identify target by type, text, or custom attributes
3. Call method: `autoui.call '{"target":"...","method":"Click"}'`

---

## Waiting for Widgets

Use `autoui.exists` with `wait-for` to block until a widget appears.

**CLI:**

```bash
autoebiten wait-for --condition "custom:autoui.exists:type=Dialog.found == true" --timeout 5s
```


**testkit:**
```go
func TestDialogAppears(t *testing.T) {
    game := testkit.Launch(t, "./mygame")
    defer game.Shutdown()

    // Wait for dialog to appear
    // Use .found to extract the boolean from the JSON response
    game.WaitFor(
        "custom:autoui.exists:type=Dialog.found == true",
        "5s",
        "",
    )

    // Now interact with the dialog
    game.RunCustom("autoui.call", `{"target":"type=Dialog","method":"Close"}`)
}
```

**Response path syntax:** Use `.found` or `.count` after the query to extract a specific field from the JSON response `{"found":true,"count":1}`.

---

## XML Format

### Structure

All autoui commands return XML (except autoui.call which returns JSON).

**Element names:** Widget type (Button, Container, TextInput, etc.)

**Standard attributes:**
```
x, y          → Position (pixels)
width, height → Size (pixels)
visible       → "true" or "false"
disabled      → "true" or "false"
```

**Widget-specific attributes:**

| Widget | Attributes |
|--------|------------|
| Button | text (requires font), state, toggle, focused |
| TextInput | text, focused |
| Checkbox | checked, state, text, focused |
| Slider | value, min, max, focused |
| Label | text |
| ProgressBar | value, min, max |
| List | entries, selected, focused |
| TextArea | text |
| ComboButton | label, open |
| ListComboButton | label, selected, open, focused |
| TabBook | active_tab |
| ScrollContainer | scroll_x, scroll_y, content_width, content_height |
| Text | text, max_width |

**Note:** Button `text` attribute requires font setup via `ButtonOpts.Text()`. For buttons without fonts, use `id` from CustomData for queries.

---

### Widget Identification (_addr)

Every widget includes a `_addr` attribute containing its pointer address:

```xml
<Button _addr="0x14000abc0" x="100" y="50" id="submit-btn"/>
```

`_addr` is used internally for exact widget identification. It changes between runs, so don't hard-code it in tests. Most queries should use `id` or other stable attributes instead.

---

### Custom Data (ae tags)

Add custom attributes via struct field tags:

package main

// PlayerCard demonstrates ae tag usage for custom XML attributes.
type PlayerCard struct {
	PlayerID   string `ae:"player_id"`
	PlayerName string `ae:"player_name"`
	Level      int    `ae:"level"`
}


**In widget creation:**
```go
btn := widget.NewButton(
    widget.ButtonOpts.Image(buttonImage),
    widget.ButtonOpts.WidgetOpts(
        widget.WidgetOpts.CustomData(&PlayerCard{
            PlayerID:   "p001",
            PlayerName: "Alice",
            Level:      42,
        }),
    ),
)
```

**Output:**
```xml
<Button id="submit-btn" player_id="p001" player_name="Alice" level="42"/>
```

**Rules:**
- `ae:"name"` → Use specified name as attribute
- `ae:"-"` → Skip field
- No ae tag → Use field name as-is

Nested structs flatten with dot notation:
```go
type Data struct {
    Player struct {
        Name string `ae:"name"`
    } `ae:"player"`
}
// Output: player.name="Alice"
```

**Slice flattening:**

Slices and arrays are flattened with indexed keys:

```go
type PlayerData struct {
    Name string   `ae:"name"`
    Tags []string `ae:"tags"`
}

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
autoebiten custom autoui.xpath --request "//Button[contains(@player_id,'p001')]"
```


Nested slices use multi-level indexing:
```xml
<!-- [][]string: [["a","b"],["c"]] -->
<Widget data.0.0="a" data.0.1="b" data.1.0="c"/>
```

---

### Attribute Lookup Order

When querying with `autoui.find` or XPath:

1. **Built-in:** type, x, y, width, height, visible, disabled
2. **Widget state:** text, checked, value, etc.
3. **Custom data:** ae tag attributes

All attributes are strings. Use XPath for numeric comparison:
```bash
# String comparison (lexicographic)
//Slider[@value > '50']

# Numeric comparison
//Slider[number(@value) > 50]
```

---

## XPath Queries

XPath 1.0 syntax for widget tree queries.

**Runnable Example:**

```bash
# Build and run the demo
cd examples/autoui
go build -o autoui_demo
autoebiten launch -- ./autoui_demo &
```

```bash
# All buttons
autoebiten custom autoui.xpath --request "//Button"
```

**Output:**
```xml

OK: <Button _addr="<ADDR>" disabled="false" height="40" id="submit-btn" role="primary" state="unchecked" visible="true" width="200" x="100" y="50"></Button><Button _addr="<ADDR>" disabled="false" height="40" id="cancel-btn" role="secondary" state="unchecked" visible="true" width="200" x="100" y="200"></Button>

```

---

### Common Patterns

**Find by type:**
```
//Button               → All Button widgets
//Container            → All Container widgets
//*                    → All widgets
```

**Find by attribute:**
```
//Button[@id='submit-btn']           → Exact match by id
//Button[@disabled='true']           → Disabled buttons
//*[contains(@id,'submit')]          → Id contains 'submit'
//Button[@role='primary']            → By custom attribute
```

**Numeric comparisons:**
```
//Slider[number(@value) > 50]        → Value > 50
//*[number(@width) > 100]            → Width > 100px
```

**Position-based:**
```
//Button[1]                          → First button
//Button[last()]                     → Last button
```

**Parent-child:**
```
//Container/Button                   → Buttons inside Container
```

**Multiple conditions:**
```
//Button[@disabled='false' and contains(@text,'Click')]
```

**Union:**
```
//Button | //Label                   → Buttons and Labels
```

---

### XPath Functions

**String:**
```
contains(@attr, 'text')             → Attribute contains text
starts-with(@attr, 'pre')           → Starts with prefix
```

**Numeric:**
```
number(@attr)                       → Convert for comparison
```

**Position:**
```
position()                          → Current position
last()                              → Last position
```

**Boolean:**
```
not(@disabled='true')               → Negation
```

---

## testkit Integration

Use autoui commands in tests via `RunCustom`:

```go
func TestButtonExists(t *testing.T) {
    game := testkit.Launch(t, "./mygame")
    defer game.Shutdown()

    // Find button
    output, err := game.RunCustom("autoui.find", "type=Button")
    require.NoError(t, err)
    assert.NotContains(t, output, "no widgets found")
}
```

---

### Test Patterns

**Verify widget exists:**
```go
output, _ := game.RunCustom("autoui.find", `{"type":"Button","id":"submit-btn"}`)
assert.NotContains(t, output, "no widgets found")
```

**Verify widget state:**
```go
output, _ := game.RunCustom("autoui.xpath", "//Button[@disabled='true']")
assert.NotContains(t, output, "no widgets found")
```

**Click a button:**
```go
_, err := game.RunCustom("autoui.call", `{"target":"id=submit-btn","method":"Click"}`)
require.NoError(t, err)
```

**Highlight for debugging:**
```go
game.RunCustom("autoui.highlight", "type=Button")
game.ScreenshotToFile("debug.png")
game.RunCustom("autoui.highlight", "clear")
```

---

## Examples

### Example 1: Discover and Interact

**Scenario:** Find a button and click it.

```bash
# Build and run the demo
cd examples/autoui
go build -o autoui_demo
autoebiten launch -- ./autoui_demo &
```

```bash
# Step 1: See what widgets exist
autoebiten custom autoui.tree
```

**Output:**
```xml

OK: <UI>
  <Container _addr="<ADDR>" disabled="false" height="480" visible="true" width="640" x="0" y="0">
    <Button _addr="<ADDR>" disabled="false" height="40" id="submit-btn" role="primary" state="unchecked" visible="true" width="200" x="100" y="50"></Button>
    <Button _addr="<ADDR>" disabled="false" height="40" id="cancel-btn" role="secondary" state="unchecked" visible="true" width="200" x="100" y="200"></Button>
  </Container>
</UI>

```

```bash
# Step 2: Click Submit button
autoebiten custom autoui.call --request '{"target":"id=submit-btn","method":"Click"}'
```

**Output:**
```json

OK: {"success":true}

```

---

### Example 2: Debug Widget State

**Scenario:** Check why a button isn't responding.

```bash
# Build and run the demo
cd examples/autoui
go build -o autoui_demo
autoebiten launch -- ./autoui_demo &
```

```bash
# Highlight all buttons to see positions
autoebiten custom autoui.highlight --request "type=Button"
```

**Output:**
```

OK: ok: highlighted 2 widgets

```

```bash
# Check if submit button is disabled
autoebiten custom autoui.xpath --request "//Button[@id='submit-btn']"
```

**Output:**
```xml

OK: <Button _addr="<ADDR>" disabled="false" height="40" id="submit-btn" role="primary" state="unchecked" visible="true" width="200" x="100" y="50"></Button>

```

```bash
# Clear highlights
autoebiten custom autoui.highlight --request "clear"
```

**Output:**
```

OK: ok: highlights cleared

```

---

### Example 3: Find by Custom Data

**Scenario:** Locate widgets by game-specific attributes.

```bash
# Build and run the demo
cd examples/autoui
go build -o autoui_demo
autoebiten launch -- ./autoui_demo &
```

```bash
# Find by custom player_id attribute
autoebiten custom autoui.find --request "player_id=p001"
```

**Output:**
```xml

OK: error: no widgets found matching query

```

```bash
# XPath with custom attribute
autoebiten custom autoui.xpath --request "//*[number(@level) > 40]"
```

**Output:**
```xml

OK: error: no widgets found matching XPath

```

---

### Example 4: E2E Test

**Scenario:** Complete user flow in black-box test.

**Test file:** `e2e/autoui_test.go`

```go
package e2e_test

import (
    "testing"
    "time"

    "github.com/s3cy/autoebiten/testkit"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestButtonClickFlow(t *testing.T) {
    game := testkit.Launch(t, "./examples/autoui/autoui_demo")
    defer game.Shutdown()

    require.True(t, game.WaitFor(func() bool {
        return game.Ping() == nil
    }, 5*time.Second))

    // Verify Submit button exists
    output, err := game.RunCustom("autoui.find", `{"type":"Button","id":"submit-btn"}`)
    require.NoError(t, err)
    assert.NotContains(t, output, "no widgets found")

    // Click Submit
    output, err = game.RunCustom("autoui.call", `{"target":"id=submit-btn","method":"Click"}`)
    require.NoError(t, err)
    assert.Contains(t, output, `"success":true`)
}
```

**Run:**
```bash
cd examples/autoui && go build -o autoui_demo && cd ../..
go test -v ./e2e -run TestButtonClickFlow
```

---

## API Reference

### Functions

```go
// Register all autoui commands (default "autoui." prefix)
func Register(ui *ebitenui.UI)

// Register with custom prefix
func RegisterWithPrefix(ui *ebitenui.UI, prefix string)

// Register with RWLock for concurrent UI modifications
func RegisterWithOptions(ui *ebitenui.UI, prefix string, opts *RegisterOptions)

// Draw highlight overlays (call in Draw)
func DrawHighlights(screen *ebiten.Image)

// Configure highlight duration
func SetHighlightDuration(d time.Duration)
```

### Thread Safety

If you modify the UI tree from goroutines (`AddChild`, `RemoveChild`), provide an `RWLock`:

```go
var uiLock sync.RWMutex

autoui.RegisterWithOptions(ui, "autoui", &autoui.RegisterOptions{
    RWLock: &uiLock,
})

// In goroutines modifying UI:
uiLock.Lock()
container.AddChild(newWidget)
uiLock.Unlock()
```

---

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