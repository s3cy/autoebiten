# Widget State Extraction Expansion Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add state extraction for 9 additional ebitenui widgets to enable full UI automation coverage.

**Architecture:** Extend the existing type switch in `ExtractWidgetState` with 9 new cases, each delegating to a dedicated extractor function following the established pattern.

**Tech Stack:** Go 1.x, ebitenui widget library, standard testing package

---

## File Structure

| File | Purpose |
|------|---------|
| `autoui/internal/widgetstate.go` | Add 9 case statements + 9 extractor functions |
| `autoui/internal/widgetstate_test.go` | Add 9 test functions for new extractors |

---

### Task 1: Add List Widget State Extraction

**Files:**
- Modify: `autoui/internal/widgetstate.go:15-30` (switch statement)
- Modify: `autoui/internal/widgetstate.go:120-134` (add extractor function)
- Test: `autoui/internal/widgetstate_test.go`

- [ ] **Step 1: Write the failing test**

```go
// TestExtractWidgetState_List tests list state extraction.
func TestExtractWidgetState_List(t *testing.T) {
	scrollImage := &widget.ScrollContainerImage{
		Idle: createTestNineSlice(100, 100, color.RGBA{50, 50, 50, 255}),
		Mask: createTestNineSlice(100, 100, color.RGBA{255, 255, 255, 255}),
	}

	list := widget.NewList(
		widget.ListOpts.Entries([]any{"item1", "item2", "item3"}),
		widget.ListOpts.EntryLabelFunc(func(e any) string { return e.(string) }),
		widget.ListOpts.ScrollContainerImage(scrollImage),
	)
	list.Validate()
	list.SetLocation(image.Rect(0, 0, 100, 100))

	state := internal.ExtractWidgetState(list)
	if state == nil {
		t.Fatal("Expected non-nil state")
	}

	if state["entries"] != "3" {
		t.Errorf("Expected entries='3', got '%s'", state["entries"])
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui/internal/... -v -run TestExtractWidgetState_List`
Expected: FAIL - List not handled in switch

- [ ] **Step 3: Add case statement to ExtractWidgetState**

In `widgetstate.go`, add to the switch statement after the existing cases:

```go
case *widget.List:
	extractListState(v, result)
```

- [ ] **Step 4: Add extractListState function**

Add after `extractProgressBarState`:

```go
// extractListState extracts state from a List widget.
// Attributes: entries, selected, focused
func extractListState(list *widget.List, result map[string]string) {
	result["entries"] = fmt.Sprintf("%d", len(list.Entries()))

	if selected := list.SelectedEntry(); selected != nil {
		result["selected"] = fmt.Sprintf("%v", selected)
	}

	if list.IsFocused() {
		result["focused"] = "true"
	}
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./autoui/internal/... -v -run TestExtractWidgetState_List`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add autoui/internal/widgetstate.go autoui/internal/widgetstate_test.go
git commit -m "feat(autoui): add List widget state extraction"
```

---

### Task 2: Add RadioGroup Widget State Extraction

**Files:**
- Modify: `autoui/internal/widgetstate.go`
- Test: `autoui/internal/widgetstate_test.go`

- [ ] **Step 1: Write the failing test**

```go
// TestExtractWidgetState_RadioGroup tests radio group state extraction.
func TestExtractWidgetState_RadioGroup(t *testing.T) {
	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
	}

	buttonColor := &widget.ButtonTextColor{
		Idle: color.White,
	}

	btn1 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Option 1", nil, buttonColor),
	)
	btn1.Validate()

	btn2 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Option 2", nil, buttonColor),
	)
	btn2.Validate()

	rg := widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(btn1, btn2),
	)

	state := internal.ExtractWidgetState(rg)
	if state == nil {
		t.Fatal("Expected non-nil state")
	}

	// Active should be set to one of the buttons
	if state["active"] == "" {
		t.Error("Expected active to be set")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui/internal/... -v -run TestExtractWidgetState_RadioGroup`
Expected: FAIL

- [ ] **Step 3: Add case statement to ExtractWidgetState**

```go
case *widget.RadioGroup:
	extractRadioGroupState(v, result)
```

- [ ] **Step 4: Add extractRadioGroupState function**

```go
// extractRadioGroupState extracts state from a RadioGroup widget.
// Attributes: active
func extractRadioGroupState(rg *widget.RadioGroup, result map[string]string) {
	if active := rg.Active(); active != nil {
		result["active"] = fmt.Sprintf("%T", active)
	}
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./autoui/internal/... -v -run TestExtractWidgetState_RadioGroup`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add autoui/internal/widgetstate.go autoui/internal/widgetstate_test.go
git commit -m "feat(autoui): add RadioGroup widget state extraction"
```

---

### Task 3: Add TextArea Widget State Extraction

**Files:**
- Modify: `autoui/internal/widgetstate.go`
- Test: `autoui/internal/widgetstate_test.go`

- [ ] **Step 1: Write the failing test**

```go
// TestExtractWidgetState_TextArea tests text area state extraction.
func TestExtractWidgetState_TextArea(t *testing.T) {
	scrollImage := &widget.ScrollContainerImage{
		Idle: createTestNineSlice(100, 100, color.RGBA{50, 50, 50, 255}),
		Mask: createTestNineSlice(100, 100, color.RGBA{255, 255, 255, 255}),
	}

	ta := widget.NewTextArea(
		widget.TextAreaOpts.Text("Hello World"),
		widget.TextAreaOpts.FontColor(color.White),
		widget.TextAreaOpts.ScrollContainerImage(scrollImage),
		widget.TextAreaOpts.ShowVerticalScrollbar(),
	)
	ta.Validate()
	ta.SetLocation(image.Rect(0, 0, 100, 100))

	state := internal.ExtractWidgetState(ta)
	if state == nil {
		t.Fatal("Expected non-nil state")
	}

	if state["text"] != "Hello World" {
		t.Errorf("Expected text='Hello World', got '%s'", state["text"])
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui/internal/... -v -run TestExtractWidgetState_TextArea`
Expected: FAIL

- [ ] **Step 3: Add case statement to ExtractWidgetState**

```go
case *widget.TextArea:
	extractTextAreaState(v, result)
```

- [ ] **Step 4: Add extractTextAreaState function**

```go
// extractTextAreaState extracts state from a TextArea widget.
// Attributes: text
func extractTextAreaState(ta *widget.TextArea, result map[string]string) {
	result["text"] = ta.GetText()
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./autoui/internal/... -v -run TestExtractWidgetState_TextArea`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add autoui/internal/widgetstate.go autoui/internal/widgetstate_test.go
git commit -m "feat(autoui): add TextArea widget state extraction"
```

---

### Task 4: Add ComboButton Widget State Extraction

**Files:**
- Modify: `autoui/internal/widgetstate.go`
- Test: `autoui/internal/widgetstate_test.go`

- [ ] **Step 1: Write the failing test**

```go
// TestExtractWidgetState_ComboButton tests combo button state extraction.
func TestExtractWidgetState_ComboButton(t *testing.T) {
	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
	}

	buttonColor := &widget.ButtonTextColor{
		Idle: color.White,
	}

	content := widget.NewContainer()

	cb := widget.NewComboButton(
		widget.ComboButtonOpts.ButtonOpts(
			widget.ButtonOpts.Image(buttonImage),
			widget.ButtonOpts.Text("Select", nil, buttonColor),
		),
		widget.ComboButtonOpts.Content(content),
	)
	cb.Validate()
	cb.SetLocation(image.Rect(0, 0, 100, 30))

	state := internal.ExtractWidgetState(cb)
	if state == nil {
		t.Fatal("Expected non-nil state")
	}

	if state["label"] != "Select" {
		t.Errorf("Expected label='Select', got '%s'", state["label"])
	}

	if state["open"] != "false" {
		t.Errorf("Expected open='false', got '%s'", state["open"])
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui/internal/... -v -run TestExtractWidgetState_ComboButton`
Expected: FAIL

- [ ] **Step 3: Add case statement to ExtractWidgetState**

```go
case *widget.ComboButton:
	extractComboButtonState(v, result)
```

- [ ] **Step 4: Add extractComboButtonState function**

```go
// extractComboButtonState extracts state from a ComboButton widget.
// Attributes: label, open
func extractComboButtonState(cb *widget.ComboButton, result map[string]string) {
	result["label"] = cb.Label()

	if cb.ContentVisible {
		result["open"] = "true"
	} else {
		result["open"] = "false"
	}
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./autoui/internal/... -v -run TestExtractWidgetState_ComboButton`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add autoui/internal/widgetstate.go autoui/internal/widgetstate_test.go
git commit -m "feat(autoui): add ComboButton widget state extraction"
```

---

### Task 5: Add ListComboButton Widget State Extraction

**Files:**
- Modify: `autoui/internal/widgetstate.go`
- Test: `autoui/internal/widgetstate_test.go`

- [ ] **Step 1: Write the failing test**

```go
// TestExtractWidgetState_ListComboButton tests list combo button state extraction.
func TestExtractWidgetState_ListComboButton(t *testing.T) {
	scrollImage := &widget.ScrollContainerImage{
		Idle: createTestNineSlice(100, 100, color.RGBA{50, 50, 50, 255}),
		Mask: createTestNineSlice(100, 100, color.RGBA{255, 255, 255, 255}),
	}

	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
	}

	buttonColor := &widget.ButtonTextColor{
		Idle: color.White,
	}

	lcb := widget.NewListComboButton(
		widget.ListComboButtonOpts.Entries([]any{"opt1", "opt2"}),
		widget.ListComboButtonOpts.EntryLabelFunc(
			func(e any) string { return e.(string) },
			func(e any) string { return e.(string) },
		),
		widget.ListComboButtonOpts.Text(nil, nil, buttonColor),
		widget.ListComboButtonOpts.ButtonParams(&widget.ButtonParams{
			Image: buttonImage,
		}),
		widget.ListComboButtonOpts.ListParams(&widget.ListParams{
			ScrollContainerImage: scrollImage,
		}),
	)
	lcb.Validate()
	lcb.SetLocation(image.Rect(0, 0, 100, 30))

	state := internal.ExtractWidgetState(lcb)
	if state == nil {
		t.Fatal("Expected non-nil state")
	}

	// Label should be set (first entry by default)
	if state["label"] == "" {
		t.Error("Expected label to be set")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui/internal/... -v -run TestExtractWidgetState_ListComboButton`
Expected: FAIL

- [ ] **Step 3: Add case statement to ExtractWidgetState**

```go
case *widget.ListComboButton:
	extractListComboButtonState(v, result)
```

- [ ] **Step 4: Add extractListComboButtonState function**

```go
// extractListComboButtonState extracts state from a ListComboButton widget.
// Attributes: label, selected, open, focused
func extractListComboButtonState(lcb *widget.ListComboButton, result map[string]string) {
	result["label"] = lcb.Label()

	if selected := lcb.SelectedEntry(); selected != nil {
		result["selected"] = fmt.Sprintf("%v", selected)
	}

	if lcb.ContentVisible() {
		result["open"] = "true"
	} else {
		result["open"] = "false"
	}

	if lcb.IsFocused() {
		result["focused"] = "true"
	}
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./autoui/internal/... -v -run TestExtractWidgetState_ListComboButton`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add autoui/internal/widgetstate.go autoui/internal/widgetstate_test.go
git commit -m "feat(autoui): add ListComboButton widget state extraction"
```

---

### Task 6: Add TabBook Widget State Extraction

**Files:**
- Modify: `autoui/internal/widgetstate.go`
- Test: `autoui/internal/widgetstate_test.go`

- [ ] **Step 1: Write the failing test**

```go
// TestExtractWidgetState_TabBook tests tab book state extraction.
func TestExtractWidgetState_TabBook(t *testing.T) {
	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(80, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(80, 30, color.RGBA{80, 80, 80, 255}),
	}

	buttonColor := &widget.ButtonTextColor{
		Idle: color.White,
	}

	tab1 := widget.NewTabBookTab(
		widget.TabBookTabOpts.Label("Tab 1"),
	)
	tab1.Validate()

	tab2 := widget.NewTabBookTab(
		widget.TabBookTabOpts.Label("Tab 2"),
	)
	tab2.Validate()

	tb := widget.NewTabBook(
		widget.TabBookOpts.Tabs(tab1, tab2),
		widget.TabBookOpts.TabButtonImage(buttonImage),
		widget.TabBookOpts.TabButtonText(nil, buttonColor),
	)
	tb.Validate()
	tb.SetLocation(image.Rect(0, 0, 200, 150))

	state := internal.ExtractWidgetState(tb)
	if state == nil {
		t.Fatal("Expected non-nil state")
	}

	// active_tab should indicate the current tab
	if state["active_tab"] == "" {
		t.Error("Expected active_tab to be set")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui/internal/... -v -run TestExtractWidgetState_TabBook`
Expected: FAIL

- [ ] **Step 3: Add case statement to ExtractWidgetState**

```go
case *widget.TabBook:
	extractTabBookState(v, result)
```

- [ ] **Step 4: Add extractTabBookState function**

```go
// extractTabBookState extracts state from a TabBook widget.
// Attributes: active_tab
func extractTabBookState(tb *widget.TabBook, result map[string]string) {
	if tab := tb.Tab(); tab != nil {
		// Tab label is stored in a private field, use widget type
		result["active_tab"] = fmt.Sprintf("%T", tab)
	}
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./autoui/internal/... -v -run TestExtractWidgetState_TabBook`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add autoui/internal/widgetstate.go autoui/internal/widgetstate_test.go
git commit -m "feat(autoui): add TabBook widget state extraction"
```

---

### Task 7: Add ScrollContainer Widget State Extraction

**Files:**
- Modify: `autoui/internal/widgetstate.go`
- Test: `autoui/internal/widgetstate_test.go`

- [ ] **Step 1: Write the failing test**

```go
// TestExtractWidgetState_ScrollContainer tests scroll container state extraction.
func TestExtractWidgetState_ScrollContainer(t *testing.T) {
	scrollImage := &widget.ScrollContainerImage{
		Idle: createTestNineSlice(100, 100, color.RGBA{50, 50, 50, 255}),
		Mask: createTestNineSlice(100, 100, color.RGBA{255, 255, 255, 255}),
	}

	content := widget.NewContainer()
	content.SetLocation(image.Rect(0, 0, 200, 200))

	sc := widget.NewScrollContainer(
		widget.ScrollContainerOpts.Content(content),
		widget.ScrollContainerOpts.Image(scrollImage),
	)
	sc.Validate()
	sc.SetLocation(image.Rect(0, 0, 100, 100))

	state := internal.ExtractWidgetState(sc)
	if state == nil {
		t.Fatal("Expected non-nil state")
	}

	if state["scroll_x"] != "0.00" {
		t.Errorf("Expected scroll_x='0.00', got '%s'", state["scroll_x"])
	}

	if state["scroll_y"] != "0.00" {
		t.Errorf("Expected scroll_y='0.00', got '%s'", state["scroll_y"])
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui/internal/... -v -run TestExtractWidgetState_ScrollContainer`
Expected: FAIL

- [ ] **Step 3: Add case statement to ExtractWidgetState**

```go
case *widget.ScrollContainer:
	extractScrollContainerState(v, result)
```

- [ ] **Step 4: Add extractScrollContainerState function**

```go
// extractScrollContainerState extracts state from a ScrollContainer widget.
// Attributes: scroll_x, scroll_y, content_width, content_height
func extractScrollContainerState(sc *widget.ScrollContainer, result map[string]string) {
	result["scroll_x"] = fmt.Sprintf("%.2f", sc.ScrollLeft)
	result["scroll_y"] = fmt.Sprintf("%.2f", sc.ScrollTop)

	contentRect := sc.ContentRect()
	result["content_width"] = fmt.Sprintf("%d", contentRect.Dx())
	result["content_height"] = fmt.Sprintf("%d", contentRect.Dy())
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./autoui/internal/... -v -run TestExtractWidgetState_ScrollContainer`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add autoui/internal/widgetstate.go autoui/internal/widgetstate_test.go
git commit -m "feat(autoui): add ScrollContainer widget state extraction"
```

---

### Task 8: Add Text Widget State Extraction

**Files:**
- Modify: `autoui/internal/widgetstate.go`
- Test: `autoui/internal/widgetstate_test.go`

- [ ] **Step 1: Write the failing test**

```go
// TestExtractWidgetState_Text tests text widget state extraction.
func TestExtractWidgetState_Text(t *testing.T) {
	txt := widget.NewText(
		widget.TextOpts.Text("Sample Text", nil, color.White),
	)
	txt.Validate()
	txt.SetLocation(image.Rect(0, 0, 100, 20))

	state := internal.ExtractWidgetState(txt)
	if state == nil {
		t.Fatal("Expected non-nil state")
	}

	if state["text"] != "Sample Text" {
		t.Errorf("Expected text='Sample Text', got '%s'", state["text"])
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui/internal/... -v -run TestExtractWidgetState_Text`
Expected: FAIL

- [ ] **Step 3: Add case statement to ExtractWidgetState**

```go
case *widget.Text:
	extractTextState(v, result)
```

- [ ] **Step 4: Add extractTextState function**

```go
// extractTextState extracts state from a Text widget.
// Attributes: text, max_width
func extractTextState(txt *widget.Text, result map[string]string) {
	result["text"] = txt.Label

	if txt.MaxWidth > 0 {
		result["max_width"] = fmt.Sprintf("%.0f", txt.MaxWidth)
	}
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./autoui/internal/... -v -run TestExtractWidgetState_Text`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add autoui/internal/widgetstate.go autoui/internal/widgetstate_test.go
git commit -m "feat(autoui): add Text widget state extraction"
```

---

### Task 9: Add Window Widget State Extraction

**Files:**
- Modify: `autoui/internal/widgetstate.go`
- Test: `autoui/internal/widgetstate_test.go`

- [ ] **Step 1: Write the failing test**

```go
// TestExtractWidgetState_Window tests window state extraction.
func TestExtractWidgetState_Window(t *testing.T) {
	contents := widget.NewContainer()

	win := widget.NewWindow(
		widget.WindowOpts.Contents(contents),
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Draggable(),
	)

	state := internal.ExtractWidgetState(win)
	if state == nil {
		t.Fatal("Expected non-nil state")
	}

	if state["modal"] != "true" {
		t.Errorf("Expected modal='true', got '%s'", state["modal"])
	}

	if state["draggable"] != "true" {
		t.Errorf("Expected draggable='true', got '%s'", state["draggable"])
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui/internal/... -v -run TestExtractWidgetState_Window`
Expected: FAIL

- [ ] **Step 3: Add case statement to ExtractWidgetState**

```go
case *widget.Window:
	extractWindowState(v, result)
```

- [ ] **Step 4: Add extractWindowState function**

```go
// extractWindowState extracts state from a Window widget.
// Attributes: modal, draggable, resizeable, focused_window
func extractWindowState(win *widget.Window, result map[string]string) {
	if win.Modal {
		result["modal"] = "true"
	}

	if win.Draggable {
		result["draggable"] = "true"
	}

	if win.Resizeable {
		result["resizeable"] = "true"
	}

	if win.FocusedWindow {
		result["focused_window"] = "true"
	}
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./autoui/internal/... -v -run TestExtractWidgetState_Window`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add autoui/internal/widgetstate.go autoui/internal/widgetstate_test.go
git commit -m "feat(autoui): add Window widget state extraction"
```

---

### Task 10: Run Full Test Suite and Verify Coverage

**Files:**
- None (verification only)

- [ ] **Step 1: Run all tests**

Run: `go test ./autoui/... -v`
Expected: All tests PASS

- [ ] **Step 2: Check test coverage**

Run: `go test ./autoui/internal/... -cover`
Expected: Coverage >80% on widgetstate.go

- [ ] **Step 3: Final commit if needed**

If any fixes were required:
```bash
git add autoui/internal/
git commit -m "fix(autoui): resolve test issues for widget state extraction"
```

---

## Success Criteria

- [ ] All 9 new widgets have state extraction functions
- [ ] All tests pass
- [ ] Test coverage >80% on new code
- [ ] Existing tests continue to pass
- [ ] 9 commits created (one per widget type)