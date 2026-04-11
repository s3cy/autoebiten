package autoui_test

import (
	"encoding/json"
	"image"
	"image/color"
	"testing"

	"github.com/ebitenui/ebitenui/widget"
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

// createTestImage creates a simple test image for use in tests
func createTestImage(width, height int, c color.Color) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

// mockCommandContext is a mock implementation of autoebiten.CommandContext for testing.
type mockCommandContext struct {
	request  string
	response string
}

func (m *mockCommandContext) Request() string {
	return m.request
}

func (m *mockCommandContext) Respond(response string) {
	m.response = response
}

func (m *mockCommandContext) GetResponse() string {
	return m.response
}

// TestCallResponse_ResultField tests that CallResponse includes Result field.
func TestCallResponse_ResultField(t *testing.T) {
	resp := autoui.CallResponse{
		Success: true,
		Result:  42,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal CallResponse: %v", err)
	}

	// Should include result field
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if parsed["success"] != true {
		t.Errorf("Expected success=true, got %v", parsed["success"])
	}
	if parsed["result"] != 42.0 { // JSON numbers are float64
		t.Errorf("Expected result=42, got %v", parsed["result"])
	}
}

// TestCallResponse_NoResult tests CallResponse without result field.
func TestCallResponse_NoResult(t *testing.T) {
	resp := autoui.CallResponse{
		Success: true,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal CallResponse: %v", err)
	}

	expected := `{"success":true}`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}
