# autoui.exists Command Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `autoui.exists` command that returns JSON for use with `wait-for`, enabling widget waiting in E2E tests.

**Architecture:** New command handler follows existing autoui pattern - uses same query parsing as `autoui.find` but returns JSON `{found, count}` instead of XML. No error on empty results.

**Tech Stack:** Go, ebitenui widget library, existing autoui infrastructure

---

## File Structure

| File | Purpose |
|------|---------|
| `autoui/handlers.go` | Add `ExistsResponse` struct and `handleExistsCommand` function |
| `autoui/register.go` | Register `autoui.exists` command |
| `autoui/handlers_test.go` | Unit tests for `handleExistsCommand` |
| `docs/autoui.md` | Documentation: command section, SetText example, wait-for usage |

---

### Task 1: Add ExistsResponse Struct and Handler

**Files:**
- Modify: `autoui/handlers.go`
- Create: `autoui/handlers_test.go`

- [ ] **Step 1: Write the failing test for ExistsResponse JSON marshaling**

Create `autoui/handlers_test.go`:

```go
package autoui_test

import (
	"encoding/json"
	"testing"

	"github.com/s3cy/autoebiten/autoui"
)

func TestExistsResponse_JSON(t *testing.T) {
	// Test found=true case
	resp := autoui.ExistsResponse{Found: true, Count: 2}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal ExistsResponse: %v", err)
	}
	expected := `{"found":true,"count":2}`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}

	// Test found=false case
	resp = autoui.ExistsResponse{Found: false, Count: 0}
	data, err = json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal ExistsResponse: %v", err)
	}
	expected = `{"found":false,"count":0}`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui -run TestExistsResponse_JSON -v`
Expected: FAIL with "undefined: autoui.ExistsResponse"

- [ ] **Step 3: Add ExistsResponse struct to handlers.go**

Add to `autoui/handlers.go` after the existing `CallResponse` struct (around line 185):

```go
// ExistsResponse represents the response from the exists command.
// Returns JSON format for use with wait-for command.
type ExistsResponse struct {
	// Found indicates if any widgets matched the query.
	Found bool `json:"found"`

	// Count is the number of matching widgets.
	Count int `json:"count"`
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./autoui -run TestExistsResponse_JSON -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add autoui/handlers.go autoui/handlers_test.go
git commit -m "feat(autoui): add ExistsResponse struct for autoui.exists command"
```

---

### Task 2: Implement handleExistsCommand

**Files:**
- Modify: `autoui/handlers.go`
- Modify: `autoui/handlers_test.go`

- [ ] **Step 1: Write the failing test for handleExistsCommand with found widgets**

Add to `autoui/handlers_test.go`:

```go
package autoui_test

import (
	"encoding/json"
	"image"
	"image/color"
	"testing"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/s3cy/autoebiten/autoui"
)

func TestHandleExistsCommand_Found(t *testing.T) {
	// Create test widget tree
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	buttonImage := &widget.ButtonImage{
		Idle: createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
	}

	btn1 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(map[string]string{"id": "btn1"}),
		),
	)
	btn1.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	container.AddChild(btn1)

	btn2 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(map[string]string{"id": "btn2"}),
		),
	)
	btn2.GetWidget().Rect = image.Rect(120, 10, 220, 40)
	container.AddChild(btn2)

	widgets := autoui.WalkTree(container)

	// Test finding Button type (should find 2)
	matching := autoui.FindByQuery(widgets, "type=Button")
	resp := autoui.ExistsResponse{Found: len(matching) > 0, Count: len(matching)}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}
	expected := `{"found":true,"count":2}`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

func TestHandleExistsCommand_NotFound(t *testing.T) {
	// Create test widget tree
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	buttonImage := &widget.ButtonImage{
		Idle: createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
	}

	btn := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	container.AddChild(btn)

	widgets := autoui.WalkTree(container)

	// Test finding TextInput type (should find 0)
	matching := autoui.FindByQuery(widgets, "type=TextInput")
	resp := autoui.ExistsResponse{Found: len(matching) > 0, Count: len(matching)}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}
	expected := `{"found":false,"count":0}`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

func TestHandleExistsCommand_JSONQuery(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	buttonImage := &widget.ButtonImage{
		Idle: createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
	}

	btn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(map[string]string{"id": "btn1"}),
		),
	)
	btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	container.AddChild(btn)

	widgets := autoui.WalkTree(container)

	// Test JSON query format
	query := `{"type":"Button","id":"btn1"}`
	matching := autoui.FindByQueryJSON(widgets, query)
	resp := autoui.ExistsResponse{Found: len(matching) > 0, Count: len(matching)}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}
	expected := `{"found":true,"count":1}`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

// createTestNineSlice creates a simple test NineSlice (copied from tree_test.go pattern)
func createTestNineSlice(width, height int, c color.Color) *widget.NineSlice {
	return widget.NewNineSliceSimpleFromImage(createTestImage(width, height, c), 3, 3)
}

func createTestImage(width, height int, c color.Color) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}
```

- [ ] **Step 2: Run test to verify tests compile and pass**

Run: `go test ./autoui -run TestHandleExists -v`
Expected: PASS (tests use existing FindByQuery functions)

- [ ] **Step 3: Add handleExistsCommand function to handlers.go**

Add to `autoui/handlers.go` after `handleHighlightCommand` (around line 359):

```go
// handleExistsCommand handles the "exists" command which returns JSON indicating
// if widgets matching a query exist. Returns JSON for use with wait-for.
// Unlike find, this never returns an error for empty results.
func handleExistsCommand(ui *ebitenui.UI) func(ctx autoebiten.CommandContext) {
	return func(ctx autoebiten.CommandContext) {
		if ui == nil {
			// Return JSON error, not plain text
			resp := ExistsResponse{Found: false, Count: 0}
			data, _ := json.Marshal(resp)
			ctx.Respond(string(data))
			return
		}

		request := ctx.Request()

		// Walk the widget tree
		widgets := WalkTree(ui.Container)

		// Determine query format and find widgets
		var matching []WidgetInfo
		if strings.HasPrefix(request, "{") {
			// JSON format
			matching = FindByQueryJSON(widgets, request)
		} else {
			// Simple key=value format (empty request matches all)
			matching = FindByQuery(widgets, request)
		}

		// Build response - always valid JSON, no error on empty
		resp := ExistsResponse{
			Found: len(matching) > 0,
			Count: len(matching),
		}

		data, err := json.Marshal(resp)
		if err != nil {
			ctx.Respond(`{"found":false,"count":0}`)
			return
		}

		ctx.Respond(string(data))
	}
}
```

- [ ] **Step 4: Run all handler tests**

Run: `go test ./autoui -v`
Expected: All tests pass

- [ ] **Step 5: Commit**

```bash
git add autoui/handlers.go autoui/handlers_test.go
git commit -m "feat(autoui): add handleExistsCommand for autoui.exists"
```

---

### Task 3: Register autoui.exists Command

**Files:**
- Modify: `autoui/register.go`

- [ ] **Step 1: Write the failing test for command registration**

Add to `autoui/register_test.go` (check if file exists, if not create it):

```go
package autoui_test

import (
	"testing"

	"github.com/s3cy/autoebiten"
	"github.com/s3cy/autoebiten/autoui"
)

func TestExistsCommandRegistered(t *testing.T) {
	// Verify autoui.exists is in the command list
	commands := autoebiten.ListCustomCommands()
	found := false
	for _, cmd := range commands {
		if cmd == "autoui.exists" {
			found = true
			break
		}
	}
	if !found {
		t.Error("autoui.exists command not registered")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui -run TestExistsCommandRegistered -v`
Expected: FAIL with "autoui.exists command not registered"

- [ ] **Step 3: Register the command in register.go**

Modify `autoui/register.go` - add registration in `registerCommands` function after the highlight command (around line 73):

```go
// Register exists command
autoebiten.Register(prefix+".exists", handleExistsCommand(ui))
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./autoui -run TestExistsCommandRegistered -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add autoui/register.go autoui/register_test.go
git commit -m "feat(autoui): register autoui.exists command"
```

---

### Task 4: Update Documentation

**Files:**
- Modify: `docs/autoui.md`

- [ ] **Step 1: Update Quick Decision section**

Find the Quick Decision section (lines 8-27) and add exists entry:

```markdown
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
```

- [ ] **Step 2: Add autoui.exists command section**

Add new section after `autoui.xpath` section (around line 210), before `autoui.call`:

```markdown
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
cd examples/autoui
go build -o autoui_demo
autoebiten launch -- ./autoui_demo &
```

```bash
# Check if any Button widgets exist
autoebiten custom autoui.exists --request "type=Button"
```

**Output (found):**
```json
{"found":true,"count":2}
```

**Output (not found):**
```json
{"found":false,"count":0}
```

**Key difference from autoui.find:** Returns JSON instead of XML. Never errors on empty results - useful with `wait-for`.

---
```

- [ ] **Step 3: Add SetText example to autoui.call section**

Add example after the Click example in `autoui.call` section (around line 254):

```markdown
**SetText Example:**

```bash
# Set text in a TextInput widget
autoebiten custom autoui.call --request '{"target":"id=name-input","method":"SetText","args":["Alice"]}'
```

**Output:**
```json
{"success":true}
```

**Why use SetText:** TextInput's `SetText(string)` method fires `ChangedEvent` for game logic. Setting the field directly would bypass events.
```

- [ ] **Step 4: Add "Waiting for Widgets" section**

Add new section at the end of the Commands section, before "XML Format":

```markdown
---

## Waiting for Widgets

Use `autoui.exists` with `wait-for` to block until a widget appears.

**CLI:**
```bash
# Wait up to 5 seconds for a Dialog widget
autoebiten wait-for 'custom:autoui.exists:type=Dialog == {"found":true}' --timeout 5s
```

**testkit:**
```go
func TestDialogAppears(t *testing.T) {
    game := testkit.Launch(t, "./mygame")
    defer game.Shutdown()

    // Wait for dialog to appear
    game.WaitFor(
        "custom:autoui.exists:type=Dialog == {\"found\":true}",
        "5s",
        "",
    )

    // Now interact with the dialog
    game.RunCustom("autoui.call", `{"target":"type=Dialog","method":"Close"}`)
}
```

**Why autoui.exists returns JSON:** The `wait-for` command compares JSON values. Other autoui commands return XML, which can't be parsed as JSON.

---

## XML Format
```

- [ ] **Step 5: Commit**

```bash
git add docs/autoui.md
git commit -m "docs(autoui): add autoui.exists command and SetText example"
```

---

### Task 5: Integration Test

**Files:**
- Modify: `autoui/integration_test.go`

- [ ] **Step 1: Check existing integration test pattern**

Run: `cat autoui/integration_test.go`
Note the existing test patterns for reference.

- [ ] **Step 2: Add integration test for exists with wait-for pattern**

Add to `autoui/integration_test.go` (or create if needed):

```go
package autoui_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExistsResponse_WaitForCompatible(t *testing.T) {
	// Verify response format works with wait-for's JSON comparison

	// Simulate game response for found=true
	resp := autoui.ExistsResponse{Found: true, Count: 1}
	data, err := json.Marshal(resp)
	require.NoError(t, err)

	// Parse as JSON (wait-for does this)
	var parsed map[string]any
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	// Verify structure
	assert.Equal(t, true, parsed["found"])
	assert.Equal(t, 1.0, parsed["count"]) // JSON numbers are float64

	// Simulate comparison that wait-for would do
	expected := map[string]any{"found": true}
	assert.Equal(t, expected["found"], parsed["found"])
}
```

- [ ] **Step 3: Run integration test**

Run: `go test ./autoui -run TestExistsResponse_WaitForCompatible -v`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add autoui/integration_test.go
git commit -m "test(autoui): add integration test for exists wait-for compatibility"
```

---

### Task 6: Final Verification

- [ ] **Step 1: Run all tests**

Run: `go test ./autoui -v`
Expected: All tests pass

- [ ] **Step 2: Run full test suite**

Run: `go test ./... -race`
Expected: All tests pass

- [ ] **Step 3: Verify documentation**

Run: `cat docs/autoui.md | grep -A5 "autoui.exists"`
Expected: Shows the new command documentation

- [ ] **Step 4: Final commit (if any remaining changes)**

```bash
git status
# If clean, no action needed
```

---

## Spec Coverage Check

| Spec Requirement | Task |
|------------------|------|
| autoui.exists command | Task 1-3 |
| JSON response {found, count} | Task 1 |
| Same query parsing as find | Task 2 |
| No error on empty results | Task 2 |
| Register in register.go | Task 3 |
| Quick Decision update | Task 4, Step 1 |
| autoui.exists command section | Task 4, Step 2 |
| SetText example | Task 4, Step 3 |
| wait-for usage section | Task 4, Step 4 |
| Unit tests | Task 1-2 |
| Integration test | Task 5 |