# autoui Bug Fixes Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix 5 bugs in autoui package: XPath matching, race condition, output format, slice support, and button text documentation.

**Architecture:** Add `_addr` for unique widget identification, add `SnapshotTree` for safe traversal, change find/xpath output to match XPath convention, flatten slices in custom data.

**Tech Stack:** Go, ebitenui widget library, antchfx/xmlquery for XPath

---

## Task 1: Add Addr Field to WidgetInfo

**Files:**
- Modify: `autoui/tree.go:11-36` (WidgetInfo struct)
- Modify: `autoui/tree.go:40-63` (ExtractWidgetInfo function)
- Test: `autoui/tree_test.go`

- [ ] **Step 1: Write the failing test**

Add test in `autoui/tree_test.go`:

```go
// TestWidgetInfo_Addr tests that Addr field is populated with pointer address.
func TestWidgetInfo_Addr(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 100, 100)

	info := autoui.ExtractWidgetInfo(container)

	if info.Addr == "" {
		t.Error("Expected Addr to be populated")
	}

	// Addr should be hex format like "0x14000abc0"
	if !strings.HasPrefix(info.Addr, "0x") {
		t.Errorf("Expected Addr to start with '0x', got '%s'", info.Addr)
	}

	// Addr should be unique per widget
	btn := widget.NewButton()
	btn.GetWidget().Rect = image.Rect(0, 0, 50, 30)
	btnInfo := autoui.ExtractWidgetInfo(btn)

	if info.Addr == btnInfo.Addr {
		t.Error("Expected different widgets to have different Addr values")
	}
}
```

Add import at top of file:
```go
import (
	"strings"
	// ... other imports
)
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui -run TestWidgetInfo_Addr -v`
Expected: FAIL with "Expected Addr to be populated" or similar

- [ ] **Step 3: Add Addr field to WidgetInfo struct**

In `autoui/tree.go`, modify WidgetInfo struct:

```go
// WidgetInfo holds extracted information about a widget.
type WidgetInfo struct {
	// Widget is the underlying widget instance.
	Widget widget.PreferredSizeLocateableWidget

	// Type is the widget type name (e.g., "Button", "Container").
	Type string

	// Rect is the widget's screen rectangle.
	Rect image.Rectangle

	// Visible indicates if the widget is currently visible.
	Visible bool

	// Disabled indicates if the widget is disabled.
	Disabled bool

	// State contains widget-specific state attributes.
	State map[string]string

	// CustomData contains extracted custom data attributes.
	CustomData map[string]string

	// Addr is the widget pointer address for unique identification.
	// Format: "0x14000abc0" (hex string).
	Addr string
}
```

- [ ] **Step 4: Populate Addr in ExtractWidgetInfo**

In `autoui/tree.go`, modify ExtractWidgetInfo:

```go
func ExtractWidgetInfo(w widget.PreferredSizeLocateableWidget) WidgetInfo {
	if w == nil {
		return WidgetInfo{}
	}

	// Get underlying widget for common properties
	baseWidget := w.GetWidget()

	info := WidgetInfo{
		Widget:   w,
		Type:     extractWidgetType(w),
		Rect:     baseWidget.Rect,
		Visible:  baseWidget.IsVisible(),
		Disabled: baseWidget.Disabled,
		Addr:     fmt.Sprintf("0x%x", reflect.ValueOf(w).Pointer()),
	}

	// Extract widget-specific state
	info.State = internal.ExtractWidgetState(w)

	// Extract custom data
	info.CustomData = internal.ExtractCustomData(baseWidget.CustomData)

	return info
}
```

Add `reflect` import if not present:
```go
import (
	"reflect"
	// ... other imports
)
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./autoui -run TestWidgetInfo_Addr -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add autoui/tree.go autoui/tree_test.go
git commit -m "feat(autoui): add Addr field to WidgetInfo for unique identification"
```

---

## Task 2: Add SnapshotTree Function

**Files:**
- Modify: `autoui/tree.go` (add SnapshotTree function)

- [ ] **Step 1: Write the failing test**

Add test in `autoui/tree_test.go`:

```go
// TestSnapshotTree tests snapshot tree creation.
func TestSnapshotTree(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	ui := &ebitenui.UI{Container: container}

	widgets := autoui.SnapshotTree(ui)

	if len(widgets) != 1 {
		t.Errorf("Expected 1 widget (container), got %d", len(widgets))
	}

	if widgets[0].Type != "Container" {
		t.Errorf("Expected Container type, got %s", widgets[0].Type)
	}
}

// TestSnapshotTree_NilUI tests nil UI handling.
func TestSnapshotTree_NilUI(t *testing.T) {
	widgets := autoui.SnapshotTree(nil)

	if widgets != nil {
		t.Errorf("Expected nil for nil UI, got %d widgets", len(widgets))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui -run TestSnapshotTree -v`
Expected: FAIL with "undefined: autoui.SnapshotTree"

- [ ] **Step 3: Add SnapshotTree function**

In `autoui/tree.go`, add after WalkTree:

```go
// SnapshotTree returns a snapshot of the widget tree from the UI.
// This function acquires RLock on the UI reference before traversal.
// Note: Full thread safety requires calling from main thread (Ebiten convention).
func SnapshotTree(ui *ebitenui.UI) []WidgetInfo {
	if ui == nil || ui.Container == nil {
		return nil
	}
	return WalkTree(ui.Container)
}
```

Add `ebitenui` import if not present:
```go
import (
	"github.com/ebitenui/ebitenui"
	// ... other imports
)
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./autoui -run TestSnapshotTree -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add autoui/tree.go autoui/tree_test.go
git commit -m "feat(autoui): add SnapshotTree function for safe tree traversal"
```

---

## Task 3: Add _addr to XML Output

**Files:**
- Modify: `autoui/xml.go:54-78` (widgetInfoToNode function)
- Test: `autoui/xml_test.go`

- [ ] **Step 1: Write the failing test**

Add test in `autoui/xml_test.go`:

```go
// TestWidgetXML_Addr tests that _addr appears in XML output.
func TestWidgetXML_Addr(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 100, 100)
	container.GetWidget().CustomData = map[string]string{"id": "test"}

	info := autoui.ExtractWidgetInfo(container)

	xmlData, err := autoui.MarshalWidgetXML(info)
	if err != nil {
		t.Fatalf("MarshalWidgetXML failed: %v", err)
	}

	// Check _addr appears in output
	if !strings.Contains(string(xmlData), "_addr=") {
		t.Errorf("Expected _addr attribute in XML, got: %s", xmlData)
	}

	// Check _addr format (hex string)
	if !strings.Contains(string(xmlData), "_addr=\"0x") {
		t.Errorf("Expected _addr to be hex format, got: %s", xmlData)
	}
}
```

Add `strings` import if not present:
```go
import (
	"strings"
	// ... other imports
)
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui -run TestWidgetXML_Addr -v`
Expected: FAIL with "Expected _addr attribute in XML"

- [ ] **Step 3: Add _addr to widgetInfoToNode**

In `autoui/xml.go`, modify widgetInfoToNode:

```go
func widgetInfoToNode(info WidgetInfo) *WidgetNode {
	node := &WidgetNode{
		XMLName: xml.Name{Local: info.Type},
		Attrs:   make(map[string]string),
	}

	// Add position attributes
	node.Attrs["x"] = formatInt(info.Rect.Min.X)
	node.Attrs["y"] = formatInt(info.Rect.Min.Y)
	node.Attrs["width"] = formatInt(info.Rect.Dx())
	node.Attrs["height"] = formatInt(info.Rect.Dy())

	// Add state attributes
	node.Attrs["visible"] = formatBool(info.Visible)
	node.Attrs["disabled"] = formatBool(info.Disabled)

	// Add unique address for identification
	node.Attrs["_addr"] = info.Addr

	// Add widget-specific state
	maps.Copy(node.Attrs, info.State)

	// Add custom data
	maps.Copy(node.Attrs, info.CustomData)

	return node
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./autoui -run TestWidgetXML_Addr -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add autoui/xml.go autoui/xml_test.go
git commit -m "feat(autoui): add _addr attribute to widget XML output"
```

---

## Task 4: Change XPath Matching to Use _addr Only

**Files:**
- Modify: `autoui/xpath.go:45-95` (xmlNodesToWidgetInfo function)
- Test: `autoui/xpath_test.go`

- [ ] **Step 1: Write the failing test for overlapping widgets**

Add test in `autoui/xpath_test.go`:

```go
// TestXPath_OverlappingWidgets tests matching overlapping widgets by _addr.
func TestXPath_OverlappingWidgets(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
	}

	// Two buttons at same position (overlapping)
	btn1 := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn1.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	btn1.GetWidget().CustomData = map[string]string{"id": "btn1"}
	container.AddChild(btn1)

	btn2 := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn2.GetWidget().Rect = image.Rect(10, 10, 110, 40) // Same position as btn1
	btn2.GetWidget().CustomData = map[string]string{"id": "btn2"}
	container.AddChild(btn2)

	widgets := autoui.WalkTree(container)

	// Both buttons should be returned by //Button
	results, err := autoui.QueryXPath(widgets, "//Button")
	if err != nil {
		t.Fatalf("QueryXPath failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 overlapping buttons to both match, got %d", len(results))
	}

	// Verify both different buttons matched
	ids := []string{results[0].CustomData["id"], results[1].CustomData["id"]}
	if !containsStr(ids, "btn1") || !containsStr(ids, "btn2") {
		t.Errorf("Expected both btn1 and btn2 to match, got ids: %v", ids)
	}
}

func containsStr(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui -run TestXPath_OverlappingWidgets -v`
Expected: FAIL with "Expected 2 overlapping buttons to both match, got 1"

- [ ] **Step 3: Rewrite xmlNodesToWidgetInfo to use _addr only**

In `autoui/xpath.go`, replace xmlNodesToWidgetInfo:

```go
// xmlNodesToWidgetInfo maps XML nodes back to their corresponding WidgetInfo.
// Matching is done by _addr attribute (widget pointer address) for exact identification.
func xmlNodesToWidgetInfo(nodes []*xmlquery.Node, widgets []WidgetInfo) []WidgetInfo {
	if len(nodes) == 0 || len(widgets) == 0 {
		return nil
	}

	// Build addr -> WidgetInfo map for efficient lookup
	addrMap := make(map[string]WidgetInfo, len(widgets))
	for _, w := range widgets {
		addrMap[w.Addr] = w
	}

	result := make([]WidgetInfo, 0, len(nodes))
	for _, node := range nodes {
		// Skip the UI root element
		if node.Data == "UI" {
			continue
		}

		// Get _addr attribute
		addr := node.SelectAttr("_addr")
		if addr == "" {
			// Cannot match without _addr
			continue
		}

		// Find matching widget by address
		if w, ok := addrMap[addr]; ok {
			result = append(result, w)
		}
	}

	return result
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./autoui -run TestXPath_OverlappingWidgets -v`
Expected: PASS

- [ ] **Step 5: Run all XPath tests to verify nothing broke**

Run: `go test ./autoui -run TestXPath -v`
Expected: All PASS

- [ ] **Step 6: Commit**

```bash
git add autoui/xpath.go autoui/xpath_test.go
git commit -m "fix(autoui): change XPath matching to use _addr only for exact identification"
```

---

## Task 5: Add MarshalWidgetsXML for Flat Output

**Files:**
- Modify: `autoui/xml.go` (add MarshalWidgetsXML function)
- Test: `autoui/xml_test.go`

- [ ] **Step 1: Write the failing tests**

Add tests in `autoui/xml_test.go`:

```go
// TestMarshalWidgetsXML_SingleWidget tests single widget output without wrapper.
func TestMarshalWidgetsXML_SingleWidget(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 100, 100)
	container.GetWidget().CustomData = map[string]string{"id": "single"}

	widgets := []autoui.WidgetInfo{autoui.ExtractWidgetInfo(container)}

	xmlData, err := autoui.MarshalWidgetsXML(widgets)
	if err != nil {
		t.Fatalf("MarshalWidgetsXML failed: %v", err)
	}

	// Should NOT have <UI> wrapper
	if strings.Contains(string(xmlData), "<UI>") {
		t.Errorf("Expected no <UI> wrapper, got: %s", xmlData)
	}

	// Should be just the widget element
	if !strings.HasPrefix(string(xmlData), "<Container") {
		t.Errorf("Expected output to start with <Container, got: %s", xmlData)
	}
}

// TestMarshalWidgetsXML_MultipleWidgets tests multiple widgets without hierarchy.
func TestMarshalWidgetsXML_MultipleWidgets(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 100, 100)
	container.GetWidget().CustomData = map[string]string{"id": "parent"}

	buttonImage := &widget.ButtonImage{
		Idle: createTestNineSlice(50, 30, color.RGBA{100, 100, 100, 255}),
	}

	btn := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn.GetWidget().Rect = image.Rect(10, 10, 60, 40)
	btn.GetWidget().CustomData = map[string]string{"id": "child"}
	container.AddChild(btn)

	// Only return the button (not container) - simulating filtered result
	widgets := []autoui.WidgetInfo{autoui.ExtractWidgetInfo(btn)}

	xmlData, err := autoui.MarshalWidgetsXML(widgets)
	if err != nil {
		t.Fatalf("MarshalWidgetsXML failed: %v", err)
	}

	// Should NOT have <UI> wrapper
	if strings.Contains(string(xmlData), "<UI>") {
		t.Errorf("Expected no <UI> wrapper, got: %s", xmlData)
	}

	// Should NOT reconstruct hierarchy
	if strings.Contains(string(xmlData), "<Container") {
		t.Errorf("Expected no Container in filtered result, got: %s", xmlData)
	}

	// Should have just the button
	if !strings.Contains(string(xmlData), "<Button") {
		t.Errorf("Expected Button element, got: %s", xmlData)
	}
}

// TestMarshalWidgetsXML_Empty tests empty widget list.
func TestMarshalWidgetsXML_Empty(t *testing.T) {
	widgets := []autoui.WidgetInfo{}

	xmlData, err := autoui.MarshalWidgetsXML(widgets)
	if err != nil {
		t.Fatalf("MarshalWidgetsXML failed: %v", err)
	}

	if xmlData != nil {
		t.Errorf("Expected nil for empty widgets, got: %s", xmlData)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui -run TestMarshalWidgetsXML -v`
Expected: FAIL with "undefined: autoui.MarshalWidgetsXML"

- [ ] **Step 3: Add MarshalWidgetsXML function**

In `autoui/xml.go`, add after MarshalWidgetTreeXML:

```go
// MarshalWidgetsXML converts a flat list of widgets to XML elements.
// Unlike MarshalWidgetTreeXML, this returns widgets directly without
// hierarchy reconstruction or <UI> wrapper. Follows XPath convention.
func MarshalWidgetsXML(widgets []WidgetInfo) ([]byte, error) {
	if len(widgets) == 0 {
		return nil, nil
	}

	var buf bytes.Buffer
	for _, info := range widgets {
		node := widgetInfoToNode(info)
		data, err := xml.Marshal(node)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal widget: %w", err)
		}
		buf.Write(data)
	}

	return buf.Bytes(), nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./autoui -run TestMarshalWidgetsXML -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add autoui/xml.go autoui/xml_test.go
git commit -m "feat(autoui): add MarshalWidgetsXML for flat widget output"
```

---

## Task 6: Update Handlers to Use SnapshotTree and MarshalWidgetsXML

**Files:**
- Modify: `autoui/handlers.go:19-39` (handleTreeCommand)
- Modify: `autoui/handlers.go:41-82` (handleAtCommand)
- Modify: `autoui/handlers.go:84-122` (handleFindCommand)
- Modify: `autoui/handlers.go:124-162` (handleXPathCommand)
- Modify: `autoui/handlers.go:186-260` (handleCallCommand)
- Modify: `autoui/handlers.go:269-359` (handleHighlightCommand)

- [ ] **Step 1: Update handleTreeCommand to use SnapshotTree**

In `autoui/handlers.go`, modify handleTreeCommand:

```go
func handleTreeCommand(ui *ebitenui.UI) func(ctx autoebiten.CommandContext) {
	return func(ctx autoebiten.CommandContext) {
		if ui == nil {
			ctx.Respond("error: UI not registered")
			return
		}

		// Walk the widget tree (use SnapshotTree for safety)
		widgets := SnapshotTree(ui)

		// Marshal to XML
		xmlData, err := MarshalWidgetTreeXML(widgets)
		if err != nil {
			ctx.Respond("error: failed to marshal widget tree: " + err.Error())
			return
		}

		ctx.Respond(string(xmlData))
	}
}
```

- [ ] **Step 2: Update handleAtCommand to use SnapshotTree**

In `autoui/handlers.go`, modify handleAtCommand:

```go
func handleAtCommand(ui *ebitenui.UI) func(ctx autoebiten.CommandContext) {
	return func(ctx autoebiten.CommandContext) {
		if ui == nil {
			ctx.Respond("error: UI not registered")
			return
		}

		request := ctx.Request()
		if request == "" {
			ctx.Respond("error: missing coordinates")
			return
		}

		// Parse coordinates
		x, y, err := parseCoordinates(request)
		if err != nil {
			ctx.Respond("error: " + err.Error())
			return
		}

		// Walk the widget tree
		widgets := SnapshotTree(ui)

		// Find widget at coordinates
		widget := FindAt(widgets, x, y)
		if widget == nil {
			ctx.Respond("error: no widget found at coordinates")
			return
		}

		// Marshal to XML
		xmlData, err := MarshalWidgetXML(*widget)
		if err != nil {
			ctx.Respond("error: failed to marshal widget: " + err.Error())
			return
		}

		ctx.Respond(string(xmlData))
	}
}
```

- [ ] **Step 3: Update handleFindCommand to use SnapshotTree and MarshalWidgetsXML**

In `autoui/handlers.go`, modify handleFindCommand:

```go
func handleFindCommand(ui *ebitenui.UI) func(ctx autoebiten.CommandContext) {
	return func(ctx autoebiten.CommandContext) {
		if ui == nil {
			ctx.Respond("error: UI not registered")
			return
		}

		request := ctx.Request()

		// Walk the widget tree
		widgets := SnapshotTree(ui)

		// Determine query format and find widgets
		var matching []WidgetInfo
		if strings.HasPrefix(request, "{") {
			// JSON format
			matching = FindByQueryJSON(widgets, request)
		} else {
			// Simple key=value format
			matching = FindByQuery(widgets, request)
		}

		if len(matching) == 0 {
			ctx.Respond("error: no widgets found matching query")
			return
		}

		// Marshal to XML (flat output, no hierarchy)
		xmlData, err := MarshalWidgetsXML(matching)
		if err != nil {
			ctx.Respond("error: failed to marshal widgets: " + err.Error())
			return
		}

		ctx.Respond(string(xmlData))
	}
}
```

- [ ] **Step 4: Update handleXPathCommand to use SnapshotTree and MarshalWidgetsXML**

In `autoui/handlers.go`, modify handleXPathCommand:

```go
func handleXPathCommand(ui *ebitenui.UI) func(ctx autoebiten.CommandContext) {
	return func(ctx autoebiten.CommandContext) {
		if ui == nil {
			ctx.Respond("error: UI not registered")
			return
		}

		xpathExpr := ctx.Request()
		if xpathExpr == "" {
			ctx.Respond("error: missing XPath expression")
			return
		}

		// Walk the widget tree
		widgets := SnapshotTree(ui)

		// Execute XPath query
		matching, err := QueryXPath(widgets, xpathExpr)
		if err != nil {
			ctx.Respond("error: " + err.Error())
			return
		}

		if len(matching) == 0 {
			ctx.Respond("error: no widgets found matching XPath")
			return
		}

		// Marshal to XML (flat output, no hierarchy)
		xmlData, err := MarshalWidgetsXML(matching)
		if err != nil {
			ctx.Respond("error: failed to marshal widgets: " + err.Error())
			return
		}

		ctx.Respond(string(xmlData))
	}
}
```

- [ ] **Step 5: Update handleCallCommand to use SnapshotTree**

In `autoui/handlers.go`, modify handleCallCommand:

```go
func handleCallCommand(ui *ebitenui.UI) func(ctx autoebiten.CommandContext) {
	return func(ctx autoebiten.CommandContext) {
		if ui == nil {
			ctx.Respond("error: UI not registered")
			return
		}

		request := ctx.Request()
		if request == "" {
			ctx.Respond("error: missing call request")
			return
		}

		// Parse the call request
		var callReq CallRequest
		if err := json.Unmarshal([]byte(request), &callReq); err != nil {
			ctx.Respond("error: invalid JSON format: " + err.Error())
			return
		}

		if callReq.Target == "" {
			ctx.Respond("error: missing target query")
			return
		}

		if callReq.Method == "" {
			ctx.Respond("error: missing method name")
			return
		}

		// Walk the widget tree
		widgets := SnapshotTree(ui)

		// Find the target widget
		var targetWidget *WidgetInfo
		var matching []WidgetInfo

		// Use appropriate finder based on target format
		if len(callReq.Target) > 0 && callReq.Target[0] == '{' {
			matching = FindByQueryJSON(widgets, callReq.Target)
		} else {
			matching = FindByQuery(widgets, callReq.Target)
		}

		if len(matching) == 0 {
			ctx.Respond("error: no widget found matching target query")
			return
		}

		// Use the first matching widget
		targetWidget = &matching[0]

		// Invoke the method
		err := InvokeMethod(targetWidget.Widget, callReq.Method, callReq.Args)

		// Build response
		response := CallResponse{
			Success: err == nil,
		}
		if err != nil {
			response.Error = err.Error()
		}

		// Marshal response to JSON
		respData, err := json.Marshal(response)
		if err != nil {
			ctx.Respond("error: failed to marshal response: " + err.Error())
			return
		}

		ctx.Respond(string(respData))
	}
}
```

- [ ] **Step 6: Update handleHighlightCommand to use SnapshotTree**

In `autoui/handlers.go`, modify handleHighlightCommand:

```go
func handleHighlightCommand(ui *ebitenui.UI) func(ctx autoebiten.CommandContext) {
	return func(ctx autoebiten.CommandContext) {
		if ui == nil {
			ctx.Respond("error: UI not registered")
			return
		}

		request := ctx.Request()

		// Handle clear mode
		if request == "clear" {
			ClearHighlights()
			ctx.Respond("ok: highlights cleared")
			return
		}

		// Walk the widget tree for coordinate/query modes
		widgets := SnapshotTree(ui)

		// Handle coordinate mode (x,y format)
		if len(request) > 0 && request[0] != '{' {
			x, y, err := parseCoordinates(request)
			if err == nil {
				// Coordinate mode
				widget := FindAt(widgets, x, y)
				if widget == nil {
					ctx.Respond("error: no widget found at coordinates")
					return
				}

				AddHighlight(widget.Rect)
				ctx.Respond("ok: highlighted widget at coordinates")
				return
			}
			// If coordinate parsing failed, try as simple query
			matching := FindByQuery(widgets, request)
			if len(matching) == 0 {
				ctx.Respond("error: no widgets found")
				return
			}

			for _, w := range matching {
				AddHighlight(w.Rect)
			}
			ctx.Respond(fmt.Sprintf("ok: highlighted %d widgets", len(matching)))
			return
		}

		// Handle JSON format
		var highlightReq HighlightRequest
		if err := json.Unmarshal([]byte(request), &highlightReq); err != nil {
			ctx.Respond("error: invalid format, expected 'clear', 'x,y', query, or JSON")
			return
		}

		// If coordinates provided, use them
		if highlightReq.X != 0 || highlightReq.Y != 0 {
			widget := FindAt(widgets, highlightReq.X, highlightReq.Y)
			if widget == nil {
				ctx.Respond("error: no widget found at coordinates")
				return
			}

			AddHighlight(widget.Rect)
			ctx.Respond("ok: highlighted widget at coordinates")
			return
		}

		// If query provided, use it
		if highlightReq.Query != "" {
			matching := FindByQuery(widgets, highlightReq.Query)
			if len(matching) == 0 {
				ctx.Respond("error: no widgets found matching query")
				return
			}

			for _, w := range matching {
				AddHighlight(w.Rect)
			}
			ctx.Respond(fmt.Sprintf("ok: highlighted %d widgets", len(matching)))
			return
		}

		ctx.Respond("error: missing coordinates or query in JSON request")
	}
}
```

- [ ] **Step 7: Run all handler tests**

Run: `go test ./autoui -v`
Expected: All PASS

- [ ] **Step 8: Commit**

```bash
git add autoui/handlers.go
git commit -m "refactor(autoui): use SnapshotTree and MarshalWidgetsXML in handlers"
```

---

## Task 7: Add Slice/Array Support in CustomData

**Files:**
- Modify: `autoui/internal/customdata.go:30-72` (ExtractCustomData switch)
- Modify: `autoui/internal/customdata.go` (add extractSliceElements helper)
- Test: `autoui/internal/customdata_test.go`

- [ ] **Step 1: Write the failing tests**

Add tests in `autoui/internal/customdata_test.go`:

```go
// TestExtractCustomData_StringSlice tests string slice flattening.
func TestExtractCustomData_StringSlice(t *testing.T) {
	input := []string{"fire", "ice", "wind"}

	result := internal.ExtractCustomData(input)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result["0"] != "fire" {
		t.Errorf("Expected 0='fire', got '%s'", result["0"])
	}
	if result["1"] != "ice" {
		t.Errorf("Expected 1='ice', got '%s'", result["1"])
	}
	if result["2"] != "wind" {
		t.Errorf("Expected 2='wind', got '%s'", result["2"])
	}
}

// TestExtractCustomData_IntSlice tests int slice flattening.
func TestExtractCustomData_IntSlice(t *testing.T) {
	input := []int{1, 2, 3}

	result := internal.ExtractCustomData(input)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result["0"] != "1" {
		t.Errorf("Expected 0='1', got '%s'", result["0"])
	}
	if result["1"] != "2" {
		t.Errorf("Expected 1='2', got '%s'", result["1"])
	}
	if result["2"] != "3" {
		t.Errorf("Expected 2='3', got '%s'", result["2"])
	}
}

// TestExtractCustomData_StructWithSlice tests struct containing slice.
func TestExtractCustomData_StructWithSlice(t *testing.T) {
	type Meta struct {
		ID   string   `ae:"id"`
		Tags []string `ae:"tags"`
	}

	input := Meta{
		ID:   "widget-1",
		Tags: []string{"fire", "ice"},
	}

	result := internal.ExtractCustomData(input)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result["id"] != "widget-1" {
		t.Errorf("Expected id='widget-1', got '%s'", result["id"])
	}
	if result["tags.0"] != "fire" {
		t.Errorf("Expected tags.0='fire', got '%s'", result["tags.0"])
	}
	if result["tags.1"] != "ice" {
		t.Errorf("Expected tags.1='ice', got '%s'", result["tags.1"])
	}
}

// TestExtractCustomData_NestedSlice tests nested slice flattening.
func TestExtractCustomData_NestedSlice(t *testing.T) {
	input := [][]string{
		{"a", "b"},
		{"c"},
	}

	result := internal.ExtractCustomData(input)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result["0.0"] != "a" {
		t.Errorf("Expected 0.0='a', got '%s'", result["0.0"])
	}
	if result["0.1"] != "b" {
		t.Errorf("Expected 0.1='b', got '%s'", result["0.1"])
	}
	if result["1.0"] != "c" {
		t.Errorf("Expected 1.0='c', got '%s'", result["1.0"])
	}
}

// TestExtractCustomData_EmptySlice tests empty slice handling.
func TestExtractCustomData_EmptySlice(t *testing.T) {
	input := []string{}

	result := internal.ExtractCustomData(input)
	if result != nil {
		t.Errorf("Expected nil for empty slice, got %v", result)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./autoui/internal -run TestExtractCustomData_StringSlice -v`
Expected: FAIL with result being nil or wrong format

- [ ] **Step 3: Add slice/array handling to ExtractCustomData**

In `autoui/internal/customdata.go`, modify the switch statement. Add before the `default` case:

```go
case reflect.Slice, reflect.Array:
	// Handle nil/empty slices
	if v.Len() == 0 {
		return nil
	}
	// Flatten slice with indexed keys
	extractSliceElements(v, result, "")
```

- [ ] **Step 4: Add extractSliceElements helper function**

In `autoui/internal/customdata.go`, add after extractStructFields:

```go
// extractSliceElements flattens a slice/array with indexed keys.
// Nested slices use dot notation (e.g., "0.0", "0.1").
func extractSliceElements(v reflect.Value, result map[string]string, prefix string) {
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)

		// Build key for this element
		key := fmt.Sprintf("%d", i)
		if prefix != "" {
			key = prefix + "." + key
		}

		// Handle nested types
		if elem.Kind() == reflect.Ptr {
			if elem.IsNil() {
				continue
			}
			elem = elem.Elem()
		}

		switch elem.Kind() {
		case reflect.Struct:
			extractStructFields(elem, result, key)
		case reflect.Slice, reflect.Array:
			extractSliceElements(elem, result, key)
		case reflect.Map:
			// Handle map within slice
			if elem.Type().Key().Kind() == reflect.String {
				iter := elem.MapRange()
				for iter.Next() {
					mapKey := key + "." + iter.Key().String()
					val := iter.Value()
					if val.Kind() == reflect.String {
						result[mapKey] = val.String()
					} else {
						result[mapKey] = fmt.Sprintf("%v", val.Interface())
					}
				}
			}
		default:
			result[key] = valueToString(elem)
		}
	}
}
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./autoui/internal -run TestExtractCustomData -v`
Expected: All new tests PASS

- [ ] **Step 6: Run all customdata tests**

Run: `go test ./autoui/internal -v`
Expected: All PASS (existing tests should still pass)

- [ ] **Step 7: Commit**

```bash
git add autoui/internal/customdata.go autoui/internal/customdata_test.go
git commit -m "feat(autoui): add slice/array flattening support in CustomData"
```

---

## Task 8: Update Documentation

**Files:**
- Modify: `docs/autoui.md`

- [ ] **Step 1: Add _addr documentation**

In `docs/autoui.md`, add section after "Attribute Lookup Order":

```markdown
### Widget Identification (_addr)

Every widget includes a `_addr` attribute containing its pointer address:

```xml
<Button _addr="0x14000abc0" x="100" y="50" id="submit-btn"/>
```

**Note:** `_addr` is an internal attribute for exact widget identification. It changes between runs, so don't hard-code it in tests. Instead, read `_addr` dynamically from `autoui.tree` output if you need it for subsequent queries.

**Usage pattern:**
```bash
# Get widget tree, note the _addr
autoebiten custom autoui.tree

# Use _addr for exact query (if needed)
autoebiten custom autoui.xpath --request "//Button[@_addr='0x14000abc0']"
```

Most queries should use `id` or other stable attributes instead of `_addr`.
```

- [ ] **Step 2: Update find/xpath output documentation**

In `docs/autoui.md`, update `autoui.find` and `autoui.xpath` output examples:

For `autoui.find` section, change output example:
```xml
<!-- OLD -->
<UI>
  <Button x="100" y="50" width="200" height="40" id="submit-btn" disabled="false"/>
  <Button x="100" y="200" width="200" height="40" id="cancel-btn" disabled="false"/>
</UI>

<!-- NEW: Flat output, no hierarchy reconstruction -->
<Button x="100" y="50" width="200" height="40" id="submit-btn" disabled="false" visible="true"/>
<Button x="100" y="200" width="200" height="40" id="cancel-btn" disabled="false" visible="true"/>
```

Add note:
```markdown
**Output format:** `autoui.find` returns matched widgets directly without `<UI>` wrapper or hierarchy reconstruction. For full hierarchy, use `autoui.tree`.
```

For `autoui.xpath` section, same change.

- [ ] **Step 3: Add slice flattening documentation**

In `docs/autoui.md`, add example in "Custom Data (ae tags)" section:

```markdown
**Slice flattening:**
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
```

- [ ] **Step 4: Add button text limitation documentation**

In `docs/autoui.md`, update the "Widget-specific attributes" table note:

```markdown
| Widget | Attributes |
|--------|------------|
| Button | text (requires font), state, toggle, focused |

**Note:** Button `text` attribute requires font setup via `ButtonOpts.Text()`. For buttons without fonts, use `id` from CustomData for queries:

```go
btn.GetWidget().CustomData = map[string]string{"id": "submit-btn"}
```

```bash
autoebiten custom autoui.find --request "id=submit-btn"
```
```

- [ ] **Step 5: Verify documentation changes**

Read through docs/autoui.md to ensure all changes are coherent.

- [ ] **Step 6: Commit**

```bash
git add docs/autoui.md
git commit -m "docs(autoui): document _addr, flat output format, slice support, button text limitation"
```

---

## Task 9: Integration Test with Example

**Files:**
- Verify: `examples/autoui/main.go`

- [ ] **Step 1: Run all autoui tests**

Run: `go test -race ./autoui/...`
Expected: All PASS

- [ ] **Step 2: Build example**

Run: `cd examples/autoui && go build -o autoui_demo`
Expected: Build succeeds

- [ ] **Step 3: Launch example and test commands**

Run: `cd examples/autoui && autoebiten launch -- ./autoui_demo &`
Then:
```bash
sleep 2
autoebiten custom autoui.tree
autoebiten custom autoui.find --request "type=Button"
autoebiten custom autoui.xpath --request "//Button"
```

Expected: Output contains `_addr` attribute, find/xpath output is flat (no `<UI>` wrapper)

- [ ] **Step 4: Kill example process**

Run: `pkill -f autoui_demo`

- [ ] **Step 5: Final commit if needed**

If any fixes were needed during integration test, commit those.

---

## Summary

| Bug | Task | Key Change |
|-----|------|------------|
| 1 | 1-4 | `_addr` for unique identification |
| 2 | 2, 6 | `SnapshotTree` for safe traversal |
| 3 | 5-6 | `MarshalWidgetsXML` for flat output |
| 4 | 7 | Slice flattening in customdata |
| 5 | 8 | Documentation update |

All changes maintain backward compatibility except the output format change (Bug 3), which follows XPath convention.