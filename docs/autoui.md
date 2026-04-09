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
└─ Need full tree? → autoui.tree command
```

**Acting on widgets:**
```
├─ Need to click/interact? → autoui.call command
└─ Need visual debugging? → autoui.highlight command
```

**Testing:**
```
└─ Need autoui in tests? → Use RunCustom with autoui commands
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
<UI>
  <Container x="0" y="0" width="640" height="480" visible="true">
    <Button x="100" y="50" width="200" height="40" id="submit-btn" disabled="false"/>
    <Button x="100" y="200" width="200" height="40" id="cancel-btn" disabled="false"/>
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
<Button x="100" y="50" width="200" height="40" id="submit-btn" disabled="false" visible="true"/>
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
<UI>
  <Button x="100" y="50" width="200" height="40" id="submit-btn" disabled="false"/>
  <Button x="100" y="200" width="200" height="40" id="cancel-btn" disabled="false"/>
</UI>
```

```bash
# Find submit button specifically
autoebiten custom autoui.find --request "id=submit-btn"
```

**Output:**
```xml
<UI>
  <Button x="100" y="50" width="200" height="40" id="submit-btn" disabled="false"/>
</UI>
```

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
<UI>
  <Button x="100" y="50" width="200" height="40" id="submit-btn" disabled="false"/>
  <Button x="100" y="200" width="200" height="40" id="cancel-btn" disabled="false"/>
</UI>
```

**XPath examples:**
```
//Button                      → All Button widgets
//Button[@disabled='true']    → Disabled buttons
//Button[@id='submit-btn']    → Button with specific id
//*[contains(@id,'submit')]   → Any widget with 'submit' in id
```

**Error:** `no widgets found matching XPath`

---

### autoui.call

Invoke method on widget.

**Usage:**
```bash
autoebiten custom autoui.call --request '{"target":"id=submit-btn","method":"Click","args":[]}'
```

**Runnable Example:**

```bash
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
{"success":true}
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
ok: highlighted 2 widgets
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
ok: highlights cleared
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
{"success":true}
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
| Button | text, state, toggle, focused |
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

---

### Custom Data (ae tags)

Add custom attributes via struct field tags:

```go
type PlayerCard struct {
    PlayerID   string `ae:"player_id"`
    PlayerName string `ae:"player_name"`
    Level      int    `ae:"level"`
    Hidden     string `ae:"-"`  // Skip this field
}

// In widget creation:
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
<UI>
  <Button x="100" y="50" width="200" height="40" id="submit-btn"/>
  <Button x="100" y="200" width="200" height="40" id="cancel-btn"/>
</UI>
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
<UI>
  <Container x="0" y="0" width="640" height="480">
    <Button x="100" y="50" id="submit-btn" disabled="false"/>
    <Button x="100" y="200" id="cancel-btn" disabled="false"/>
  </Container>
</UI>
```

```bash
# Step 2: Click Submit button
autoebiten custom autoui.call --request '{"target":"id=submit-btn","method":"Click"}'
```

**Output:**
```json
{"success":true}
```

---

### Example 2: Debug Widget State

**Scenario:** Check why a button isn't responding.

```bash
cd examples/autoui
go build -o autoui_demo
autoebiten launch -- ./autoui_demo &
```

```bash
# Highlight all buttons to see positions
autoebiten custom autoui.highlight --request "type=Button"
```

```bash
# Check if submit button is disabled
autoebiten custom autoui.xpath --request "//Button[@id='submit-btn']"
```

**Output:**
```xml
<Button x="100" y="50" id="submit-btn" disabled="false"/>
```

```bash
# Clear highlights
autoebiten custom autoui.highlight --request "clear"
```

---

### Example 3: Find by Custom Data

**Scenario:** Locate widgets by game-specific attributes.

```bash
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
<Button id="submit-btn" player_id="p001" player_name="Alice" level="42"/>
```

```bash
# XPath with custom attribute
autoebiten custom autoui.xpath --request "//*[number(@level) > 40]"
```

**Output:**
```xml
<Button id="submit-btn" level="42"/>
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

**Output:**
```
=== RUN   TestButtonClickFlow
    --- PASS: TestButtonClickFlow (2.34s)
PASS
```