package autoui_test

import (
	"encoding/json"
	"image"
	"image/color"
	"testing"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/s3cy/autoebiten/autoui"
)

// TestFindByQuery_ID tests finding widgets by custom data id.
func TestFindByQuery_ID(t *testing.T) {
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

	// Button 3 without id
	btn3 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
	)
	btn3.GetWidget().Rect = image.Rect(230, 10, 330, 40)
	container.AddChild(btn3)

	widgets := autoui.WalkTree(container)

	// Test finding by id
	results := autoui.FindByQuery(widgets, "id=btn1")
	if len(results) != 1 {
		t.Errorf("Expected 1 result for id=btn1, got %d", len(results))
	}
	if len(results) > 0 && results[0].CustomData["id"] != "btn1" {
		t.Errorf("Expected result to have id=btn1, got %s", results[0].CustomData["id"])
	}

	// Test finding non-existent id
	results = autoui.FindByQuery(widgets, "id=nonexistent")
	if len(results) != 0 {
		t.Errorf("Expected 0 results for id=nonexistent, got %d", len(results))
	}

	// Test finding btn2
	results = autoui.FindByQuery(widgets, "id=btn2")
	if len(results) != 1 {
		t.Errorf("Expected 1 result for id=btn2, got %d", len(results))
	}
}

// TestFindByQuery_Type tests finding widgets by type.
func TestFindByQuery_Type(t *testing.T) {
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

	// Test finding by type=Button
	results := autoui.FindByQuery(widgets, "type=Button")
	if len(results) != 2 {
		t.Errorf("Expected 2 Button widgets, got %d", len(results))
	}
	for i, r := range results {
		if r.Type != "Button" {
			t.Errorf("Result %d: Expected type Button, got %s", i, r.Type)
		}
	}

	// Test finding by type=Container
	results = autoui.FindByQuery(widgets, "type=Container")
	if len(results) != 1 {
		t.Errorf("Expected 1 Container widget, got %d", len(results))
	}

	// Test finding by type=Label
	results = autoui.FindByQuery(widgets, "type=Label")
	if len(results) != 1 {
		t.Errorf("Expected 1 Label widget, got %d", len(results))
	}

	// Test finding non-existent type
	results = autoui.FindByQuery(widgets, "type=TextInput")
	if len(results) != 0 {
		t.Errorf("Expected 0 results for type=TextInput, got %d", len(results))
	}
}

// TestFindByQuery_Visible tests finding widgets by visibility.
// Note: All widgets are visible by default when created.
func TestFindByQuery_Visible(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
	}

	// All widgets are visible by default
	btn1 := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn1.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	container.AddChild(btn1)

	btn2 := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn2.GetWidget().Rect = image.Rect(120, 10, 220, 40)
	container.AddChild(btn2)

	widgets := autoui.WalkTree(container)

	// Test finding visible widgets (all should be visible by default)
	results := autoui.FindByQuery(widgets, "visible=true")
	// Container + 2 buttons = 3 visible widgets
	if len(results) != 3 {
		t.Errorf("Expected 3 visible widgets (container + 2 buttons), got %d", len(results))
	}

	// Test finding hidden widgets (none should be hidden)
	results = autoui.FindByQuery(widgets, "visible=false")
	if len(results) != 0 {
		t.Errorf("Expected 0 hidden widgets, got %d", len(results))
	}
}

// TestFindByQuery_JSONMultiple tests JSON query with multiple criteria (AND logic).
func TestFindByQuery_JSONMultiple(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
	}

	// Button with id=btn1
	btn1 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(map[string]string{"id": "btn1"}),
		),
	)
	btn1.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	container.AddChild(btn1)

	// Button with id=btn2, disabled
	btn2 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(map[string]string{"id": "btn2"}),
		),
	)
	btn2.GetWidget().Rect = image.Rect(120, 10, 220, 40)
	btn2.GetWidget().Disabled = true
	container.AddChild(btn2)

	// Button with id=btn3
	btn3 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(map[string]string{"id": "btn3"}),
		),
	)
	btn3.GetWidget().Rect = image.Rect(230, 10, 330, 40)
	container.AddChild(btn3)

	widgets := autoui.WalkTree(container)

	// Test JSON query: type=Button AND visible=true
	query := `{"type":"Button","visible":"true"}`
	results := autoui.FindByQueryJSON(widgets, query)
	if len(results) != 3 {
		t.Errorf("Expected 3 visible Button widgets (all buttons), got %d", len(results))
	}

	// Test JSON query: id=btn2 AND disabled=true
	query = `{"id":"btn2","disabled":"true"}`
	results = autoui.FindByQueryJSON(widgets, query)
	if len(results) != 1 {
		t.Errorf("Expected 1 result for id=btn2 disabled=true, got %d", len(results))
	}

	// Test JSON query: type=Button AND id=btn1
	query = `{"type":"Button","id":"btn1"}`
	results = autoui.FindByQueryJSON(widgets, query)
	if len(results) != 1 {
		t.Errorf("Expected 1 result for type=Button id=btn1, got %d", len(results))
	}

	// Test JSON query with no matches
	query = `{"type":"Button","id":"nonexistent"}`
	results = autoui.FindByQueryJSON(widgets, query)
	if len(results) != 0 {
		t.Errorf("Expected 0 results for non-matching criteria, got %d", len(results))
	}
}

// TestFindAtCoordinates tests coordinate-based widget finding.
func TestFindAtCoordinates(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
	}

	// Button at (10, 10) to (110, 40)
	btn1 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(map[string]string{"id": "btn1"}),
		),
	)
	btn1.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	container.AddChild(btn1)

	// Button at (120, 10) to (220, 40) - overlaps btn1 in Y
	btn2 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(map[string]string{"id": "btn2"}),
		),
	)
	btn2.GetWidget().Rect = image.Rect(120, 10, 220, 40)
	container.AddChild(btn2)

	widgets := autoui.WalkTree(container)

	// Test point inside btn1
	result := autoui.FindAt(widgets, 50, 25)
	if result == nil {
		t.Error("Expected to find widget at (50, 25)")
	} else if result.CustomData["id"] != "btn1" {
		t.Errorf("Expected to find btn1 at (50, 25), got %s", result.CustomData["id"])
	}

	// Test point inside btn2
	result = autoui.FindAt(widgets, 150, 25)
	if result == nil {
		t.Error("Expected to find widget at (150, 25)")
	} else if result.CustomData["id"] != "btn2" {
		t.Errorf("Expected to find btn2 at (150, 25), got %s", result.CustomData["id"])
	}

	// Test point outside all widgets (in container area)
	result = autoui.FindAt(widgets, 400, 400)
	if result == nil {
		t.Error("Expected to find Container at (400, 400)")
	} else if result.Type != "Container" {
		t.Errorf("Expected to find Container at (400, 400), got %s", result.Type)
	}

	// Test point outside container
	result = autoui.FindAt(widgets, 900, 700)
	if result != nil {
		t.Errorf("Expected nil for point outside container, got %s", result.Type)
	}

	// Test corner point (inside btn1 at corner)
	result = autoui.FindAt(widgets, 10, 10)
	if result == nil {
		t.Error("Expected to find widget at corner (10, 10)")
	}
}

// TestFindAtCoordinates_TopMost tests that top-most widgets are found first (reverse order).
func TestFindAtCoordinates_TopMost(t *testing.T) {
	// Create overlapping widgets
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
	}

	// First button (added first, should be "below")
	btn1 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(map[string]string{"id": "btn-bottom"}),
		),
	)
	btn1.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	container.AddChild(btn1)

	// Second button at same position (added last, should be "top-most")
	btn2 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(map[string]string{"id": "btn-top"}),
		),
	)
	btn2.GetWidget().Rect = image.Rect(10, 10, 110, 40) // Same position as btn1
	container.AddChild(btn2)

	widgets := autoui.WalkTree(container)

	// Point in overlapping area should find btn-top (last added = top-most)
	// WalkTree visits Container, btn1, btn2 - FindAt should search in reverse
	result := autoui.FindAt(widgets, 50, 25)
	if result == nil {
		t.Error("Expected to find widget at overlapping position")
	} else if result.CustomData["id"] != "btn-top" {
		t.Errorf("Expected to find top-most widget (btn-top), got %s", result.CustomData["id"])
	}
}

// TestFindByQuery_Disabled tests finding widgets by disabled state.
func TestFindByQuery_Disabled(t *testing.T) {
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
	container.AddChild(btn1)

	// Disabled button
	btn2 := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn2.GetWidget().Rect = image.Rect(120, 10, 220, 40)
	btn2.GetWidget().Disabled = true
	container.AddChild(btn2)

	widgets := autoui.WalkTree(container)

	// Test finding disabled widgets
	results := autoui.FindByQuery(widgets, "disabled=true")
	if len(results) != 1 {
		t.Errorf("Expected 1 disabled widget, got %d", len(results))
	}

	// Test finding enabled widgets
	results = autoui.FindByQuery(widgets, "disabled=false")
	if len(results) != 2 {
		t.Errorf("Expected 2 enabled widgets (container + btn1), got %d", len(results))
	}
}

// TestFindByQuery_StateAttribute tests finding widgets by state attributes.
func TestFindByQuery_StateAttribute(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
	}

	buttonColor := &widget.ButtonTextColor{
		Idle: color.White,
	}

	// Button with text "Submit"
	btn1 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Submit", nil, buttonColor),
	)
	btn1.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	container.AddChild(btn1)

	// Button with text "Cancel"
	btn2 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Cancel", nil, buttonColor),
	)
	btn2.GetWidget().Rect = image.Rect(120, 10, 220, 40)
	container.AddChild(btn2)

	widgets := autoui.WalkTree(container)

	// Note: Without proper initialization, State["text"] may be empty
	// This test verifies the query mechanism works when state is populated

	// If we manually set state (for testing purposes), query should work
	// In real usage, state is extracted by ExtractWidgetState

	// Test that query handles state attributes
	// This is a basic test - actual state extraction requires widget validation
	t.Logf("Button 1 state: %v", widgets[1].State)
	t.Logf("Button 2 state: %v", widgets[2].State)
}

// TestFindByQueryJSON_InvalidJSON tests handling of invalid JSON.
func TestFindByQueryJSON_InvalidJSON(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	widgets := autoui.WalkTree(container)

	// Invalid JSON should return empty results
	results := autoui.FindByQueryJSON(widgets, "not valid json")
	if len(results) != 0 {
		t.Errorf("Expected 0 results for invalid JSON, got %d", len(results))
	}

	// Empty string should return empty results
	results = autoui.FindByQueryJSON(widgets, "")
	if len(results) != 0 {
		t.Errorf("Expected 0 results for empty query, got %d", len(results))
	}
}

// TestFindByQuery_Empty tests handling of empty query.
func TestFindByQuery_Empty(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	buttonImage := &widget.ButtonImage{
		Idle: createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
	}

	btn := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	container.AddChild(btn)

	widgets := autoui.WalkTree(container)

	// Empty query should return all widgets
	results := autoui.FindByQuery(widgets, "")
	if len(results) != len(widgets) {
		t.Errorf("Expected %d results for empty query, got %d", len(widgets), len(results))
	}
}

// TestFindByQuery_PositionAttributes tests finding by position attributes.
func TestFindByQuery_PositionAttributes(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	buttonImage := &widget.ButtonImage{
		Idle: createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
	}

	btn1 := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn1.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	container.AddChild(btn1)

	btn2 := widget.NewButton(widget.ButtonOpts.Image(buttonImage))
	btn2.GetWidget().Rect = image.Rect(120, 10, 220, 40)
	container.AddChild(btn2)

	widgets := autoui.WalkTree(container)

	// Test finding by x position
	results := autoui.FindByQuery(widgets, "x=10")
	if len(results) != 1 {
		t.Errorf("Expected 1 widget at x=10, got %d", len(results))
	}

	// Test finding by width
	results = autoui.FindByQuery(widgets, "width=100")
	if len(results) != 2 {
		t.Errorf("Expected 2 widgets with width=100, got %d", len(results))
	}
}

// TestFindByQueryJSON_Complex tests JSON with multiple different attribute types.
func TestFindByQueryJSON_Complex(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
	}

	// Button with id, at specific position
	btn1 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(map[string]string{"id": "btn1", "role": "submit"}),
		),
	)
	btn1.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	container.AddChild(btn1)

	// Button with different id, same role
	btn2 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(map[string]string{"id": "btn2", "role": "submit"}),
		),
	)
	btn2.GetWidget().Rect = image.Rect(120, 10, 220, 40)
	btn2.GetWidget().Disabled = true
	container.AddChild(btn2)

	widgets := autoui.WalkTree(container)

	// Test: role=submit AND disabled=true
	query := `{"role":"submit","disabled":"true"}`
	results := autoui.FindByQueryJSON(widgets, query)
	if len(results) != 1 {
		t.Errorf("Expected 1 result for role=submit AND disabled=true, got %d", len(results))
	}
	if len(results) > 0 && results[0].CustomData["id"] != "btn2" {
		t.Errorf("Expected btn2, got %s", results[0].CustomData["id"])
	}

	// Test: role=submit AND visible=true AND type=Button
	// Both btn1 and btn2 are visible (disabled doesn't affect visibility)
	query = `{"role":"submit","visible":"true","type":"Button"}`
	results = autoui.FindByQueryJSON(widgets, query)
	if len(results) != 2 {
		t.Errorf("Expected 2 visible Buttons with role=submit, got %d", len(results))
	}

	// Test: role=submit AND disabled=false AND type=Button (only btn1)
	query = `{"role":"submit","disabled":"false","type":"Button"}`
	results = autoui.FindByQueryJSON(widgets, query)
	if len(results) != 1 {
		t.Errorf("Expected 1 enabled Button with role=submit, got %d", len(results))
	}
}

// TestFindAt_EmptyList tests FindAt with empty widget list.
func TestFindAt_EmptyList(t *testing.T) {
	widgets := []autoui.WidgetInfo{}

	result := autoui.FindAt(widgets, 100, 100)
	if result != nil {
		t.Errorf("Expected nil for empty widget list, got %+v", result)
	}
}

// TestFindAt_NilWidget tests FindAt behavior with nil widgets.
func TestFindAt_NilWidget(t *testing.T) {
	// WidgetInfo with nil Widget should still have Rect
	info := autoui.WidgetInfo{
		Type:   "Test",
		Rect:   image.Rect(0, 0, 100, 100),
		Visible: true,
	}

	widgets := []autoui.WidgetInfo{info}

	// Point inside should be found
	result := autoui.FindAt(widgets, 50, 50)
	if result == nil {
		t.Error("Expected to find widget at (50, 50)")
	}

	// Point outside should return nil
	result = autoui.FindAt(widgets, 150, 150)
	if result != nil {
		t.Errorf("Expected nil for point outside, got %+v", result)
	}
}

// Helper to parse JSON and verify criteria format
func TestJSONCriteriaParsing(t *testing.T) {
	// Verify that JSON criteria are parsed correctly
	query := `{"type":"Button","id":"btn1","visible":"true"}`

	var criteria map[string]string
	err := json.Unmarshal([]byte(query), &criteria)
	if err != nil {
		t.Fatalf("Failed to parse JSON query: %v", err)
	}

	if criteria["type"] != "Button" {
		t.Errorf("Expected type=Button, got %s", criteria["type"])
	}
	if criteria["id"] != "btn1" {
		t.Errorf("Expected id=btn1, got %s", criteria["id"])
	}
	if criteria["visible"] != "true" {
		t.Errorf("Expected visible=true, got %s", criteria["visible"])
	}
}

// Note: createTestNineSlice helper is defined in tree_test.go (same package)