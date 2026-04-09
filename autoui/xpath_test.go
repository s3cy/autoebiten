package autoui_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/s3cy/autoebiten/autoui"
)

// TestXPath_FindByID tests finding a button by id attribute using XPath.
func TestXPath_FindByID(t *testing.T) {
	// Create container with multiple buttons
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)
	container.GetWidget().CustomData = map[string]string{"id": "root"}

	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
		Disabled: createTestNineSlice(100, 30, color.RGBA{150, 150, 150, 255}),
	}

	// Button 1 with id "btn1"
	btn1 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(map[string]string{"id": "btn1"}),
		),
	)
	btn1.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	container.AddChild(btn1)

	// Button 2 with id "btn2"
	btn2 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(map[string]string{"id": "btn2"}),
		),
	)
	btn2.GetWidget().Rect = image.Rect(120, 10, 220, 40)
	container.AddChild(btn2)

	widgets := autoui.WalkTree(container)

	// Test XPath: //Button[@id='btn1']
	results, err := autoui.QueryXPath(widgets, "//Button[@id='btn1']")
	if err != nil {
		t.Fatalf("QueryXPath failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result for //Button[@id='btn1'], got %d", len(results))
	}

	if len(results) > 0 {
		if results[0].Type != "Button" {
			t.Errorf("Expected result type to be Button, got %s", results[0].Type)
		}
		if results[0].CustomData["id"] != "btn1" {
			t.Errorf("Expected result to have id='btn1', got %s", results[0].CustomData["id"])
		}
	}

	// Test XPath: //Button[@id='btn2']
	results, err = autoui.QueryXPath(widgets, "//Button[@id='btn2']")
	if err != nil {
		t.Fatalf("QueryXPath failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result for //Button[@id='btn2'], got %d", len(results))
	}

	// Test XPath: //Button[@id='nonexistent'] - should return empty
	results, err = autoui.QueryXPath(widgets, "//Button[@id='nonexistent']")
	if err != nil {
		t.Fatalf("QueryXPath failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results for non-existent id, got %d", len(results))
	}
}

// TestXPath_FindVisibleButtons tests finding buttons with visible='true' using XPath.
func TestXPath_FindVisibleButtons(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
	}

	// All buttons visible by default
	btn1 := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn1.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	btn1.GetWidget().CustomData = map[string]string{"id": "btn1"}
	container.AddChild(btn1)

	btn2 := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn2.GetWidget().Rect = image.Rect(120, 10, 220, 40)
	btn2.GetWidget().CustomData = map[string]string{"id": "btn2"}
	container.AddChild(btn2)

	btn3 := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn3.GetWidget().Rect = image.Rect(230, 10, 330, 40)
	btn3.GetWidget().CustomData = map[string]string{"id": "btn3"}
	container.AddChild(btn3)

	widgets := autoui.WalkTree(container)

	// Test XPath: //Button[@visible='true']
	results, err := autoui.QueryXPath(widgets, "//Button[@visible='true']")
	if err != nil {
		t.Fatalf("QueryXPath failed: %v", err)
	}

	// All 3 buttons should be visible
	if len(results) != 3 {
		t.Errorf("Expected 3 visible buttons, got %d", len(results))
	}

	// Verify all results are Buttons and visible
	for i, r := range results {
		if r.Type != "Button" {
			t.Errorf("Result %d: Expected type Button, got %s", i, r.Type)
		}
		if !r.Visible {
			t.Errorf("Result %d: Expected visible=true", i)
		}
	}

	// Test XPath: //Button[@visible='false'] - should return empty
	results, err = autoui.QueryXPath(widgets, "//Button[@visible='false']")
	if err != nil {
		t.Fatalf("QueryXPath failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 hidden buttons, got %d", len(results))
	}

	// Test XPath: Find all visible widgets (any type)
	results, err = autoui.QueryXPath(widgets, "//*[@visible='true']")
	if err != nil {
		t.Fatalf("QueryXPath failed: %v", err)
	}

	// Container + 3 buttons = 4 visible widgets
	if len(results) != 4 {
		t.Errorf("Expected 4 visible widgets (container + 3 buttons), got %d", len(results))
	}
}

// TestXPath_FindByPosition tests finding widgets by position coordinates using XPath.
func TestXPath_FindByPosition(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)
	container.GetWidget().CustomData = map[string]string{"id": "root"}

	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
	}

	// Button at x=10, y=10
	btn1 := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn1.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	btn1.GetWidget().CustomData = map[string]string{"id": "btn-at-10-10"}
	container.AddChild(btn1)

	// Button at x=120, y=10
	btn2 := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn2.GetWidget().Rect = image.Rect(120, 10, 220, 40)
	btn2.GetWidget().CustomData = map[string]string{"id": "btn-at-120-10"}
	container.AddChild(btn2)

	// Button at x=10, y=50
	btn3 := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn3.GetWidget().Rect = image.Rect(10, 50, 110, 80)
	btn3.GetWidget().CustomData = map[string]string{"id": "btn-at-10-50"}
	container.AddChild(btn3)

	widgets := autoui.WalkTree(container)

	// Test XPath: Find widgets with x=10
	results, err := autoui.QueryXPath(widgets, "//*[@x='10']")
	if err != nil {
		t.Fatalf("QueryXPath failed: %v", err)
	}

	// btn1 and btn3 have x=10
	if len(results) != 2 {
		t.Errorf("Expected 2 widgets with x=10, got %d", len(results))
	}

	// Test XPath: Find widgets with y=10
	results, err = autoui.QueryXPath(widgets, "//*[@y='10']")
	if err != nil {
		t.Fatalf("QueryXPath failed: %v", err)
	}

	// btn1 and btn2 have y=10
	if len(results) != 2 {
		t.Errorf("Expected 2 widgets with y=10, got %d", len(results))
	}

	// Test XPath: Find widget with specific x and y
	results, err = autoui.QueryXPath(widgets, "//*[@x='10' and @y='10']")
	if err != nil {
		t.Fatalf("QueryXPath failed: %v", err)
	}

	// Only btn1 has x=10 AND y=10
	if len(results) != 1 {
		t.Errorf("Expected 1 widget at x=10 and y=10, got %d", len(results))
	}

	if len(results) > 0 && results[0].CustomData["id"] != "btn-at-10-10" {
		t.Errorf("Expected btn-at-10-10, got %s", results[0].CustomData["id"])
	}

	// Test XPath: Find widgets in upper-left (x<100 and y<100)
	results, err = autoui.QueryXPath(widgets, "//*[@x<100 and @y<100]")
	if err != nil {
		t.Fatalf("QueryXPath failed: %v", err)
	}

	// btn1 (x=10, y=10), btn3 (x=10, y=50) - both have x<100 and y<100
	// Note: XPath number comparison works on attribute values
	if len(results) < 2 {
		t.Errorf("Expected at least 2 widgets with x<100 and y<100, got %d", len(results))
	}
}

// TestXPath_FindByType tests finding widgets by type using XPath.
func TestXPath_FindByType(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
	}

	btn1 := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn1.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	container.AddChild(btn1)

	btn2 := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn2.GetWidget().Rect = image.Rect(120, 10, 220, 40)
	container.AddChild(btn2)

	// Add a label
	label := widget.NewLabel(widget.LabelOpts.Text("Label", nil, nil))
	label.GetWidget().Rect = image.Rect(10, 50, 200, 70)
	container.AddChild(label)

	widgets := autoui.WalkTree(container)

	// Test XPath: //Button - find all Buttons
	results, err := autoui.QueryXPath(widgets, "//Button")
	if err != nil {
		t.Fatalf("QueryXPath failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 Button widgets, got %d", len(results))
	}

	for i, r := range results {
		if r.Type != "Button" {
			t.Errorf("Result %d: Expected type Button, got %s", i, r.Type)
		}
	}

	// Test XPath: //Container - find all Containers
	results, err = autoui.QueryXPath(widgets, "//Container")
	if err != nil {
		t.Fatalf("QueryXPath failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 Container widget, got %d", len(results))
	}

	// Test XPath: //Label - find all Labels
	results, err = autoui.QueryXPath(widgets, "//Label")
	if err != nil {
		t.Fatalf("QueryXPath failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 Label widget, got %d", len(results))
	}
}

// TestXPath_FindDisabled tests finding disabled widgets using XPath.
func TestXPath_FindDisabled(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
		Disabled: createTestNineSlice(100, 30, color.RGBA{150, 150, 150, 255}),
	}

	// Enabled button
	btn1 := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn1.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	btn1.GetWidget().Disabled = false
	btn1.GetWidget().CustomData = map[string]string{"id": "btn-enabled"}
	container.AddChild(btn1)

	// Disabled button
	btn2 := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn2.GetWidget().Rect = image.Rect(120, 10, 220, 40)
	btn2.GetWidget().Disabled = true
	btn2.GetWidget().CustomData = map[string]string{"id": "btn-disabled"}
	container.AddChild(btn2)

	widgets := autoui.WalkTree(container)

	// Test XPath: //Button[@disabled='true']
	results, err := autoui.QueryXPath(widgets, "//Button[@disabled='true']")
	if err != nil {
		t.Fatalf("QueryXPath failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 disabled button, got %d", len(results))
	}

	if len(results) > 0 {
		if !results[0].Disabled {
			t.Error("Expected result to be disabled")
		}
		if results[0].CustomData["id"] != "btn-disabled" {
			t.Errorf("Expected btn-disabled, got %s", results[0].CustomData["id"])
		}
	}

	// Test XPath: //Button[@disabled='false']
	results, err = autoui.QueryXPath(widgets, "//Button[@disabled='false']")
	if err != nil {
		t.Fatalf("QueryXPath failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 enabled button, got %d", len(results))
	}

	if len(results) > 0 && results[0].CustomData["id"] != "btn-enabled" {
		t.Errorf("Expected btn-enabled, got %s", results[0].CustomData["id"])
	}
}

// TestXPath_NestedContainer tests finding widgets inside nested containers.
func TestXPath_NestedContainer(t *testing.T) {
	root := widget.NewContainer()
	root.GetWidget().Rect = image.Rect(0, 0, 800, 600)
	root.GetWidget().CustomData = map[string]string{"id": "root"}

	// Nested container
	child := widget.NewContainer()
	child.GetWidget().Rect = image.Rect(0, 50, 800, 550)
	child.GetWidget().CustomData = map[string]string{"id": "child"}
	root.AddChild(child)

	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
	}

	// Button inside child container
	btn := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn.GetWidget().Rect = image.Rect(20, 60, 120, 90)
	btn.GetWidget().CustomData = map[string]string{"id": "btn-nested"}
	child.AddChild(btn)

	widgets := autoui.WalkTree(root)

	// Test XPath: //Button[@id='btn-nested']
	results, err := autoui.QueryXPath(widgets, "//Button[@id='btn-nested']")
	if err != nil {
		t.Fatalf("QueryXPath failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 nested button, got %d", len(results))
	}

	// Test XPath: //Container/Button - find buttons inside containers
	results, err = autoui.QueryXPath(widgets, "//Container/Button")
	if err != nil {
		t.Fatalf("QueryXPath failed: %v", err)
	}

	// The button inside child container
	if len(results) != 1 {
		t.Errorf("Expected 1 Button inside Container, got %d", len(results))
	}
}

// TestXPath_EmptyWidgetList tests XPath with empty widget list.
func TestXPath_EmptyWidgetList(t *testing.T) {
	widgets := []autoui.WidgetInfo{}

	results, err := autoui.QueryXPath(widgets, "//Button")
	if err != nil {
		t.Fatalf("QueryXPath failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results for empty widget list, got %d", len(results))
	}
}

// TestXPath_OverlappingWidgets tests matching overlapping widgets by _addr.
func TestXPath_OverlappingWidgets(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	buttonImage := &widget.ButtonImage{
		Idle:    createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed: createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
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

// containsStr checks if a string is in a slice.
func containsStr(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// TestXPath_InvalidExpression tests handling of invalid XPath expressions.
func TestXPath_InvalidExpression(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	buttonImage := &widget.ButtonImage{
		Idle: createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
	}

	btn := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	container.AddChild(btn)

	widgets := autoui.WalkTree(container)

	// Test invalid XPath expression
	results, err := autoui.QueryXPath(widgets, "invalid[xpath")
	if err == nil {
		t.Error("Expected error for invalid XPath expression")
	}
	if len(results) != 0 {
		t.Errorf("Expected 0 results for invalid XPath, got %d", len(results))
	}
}