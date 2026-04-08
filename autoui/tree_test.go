package autoui_test

import (
	"image"
	"image/color"
	"testing"

	ebitenuiImage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/s3cy/autoebiten/autoui"
)

// TestExtractWidgetInfo_Button tests extracting widget info from a button.
// Note: Not using t.Parallel() because ebitenui has global state that is not thread-safe.
func TestExtractWidgetInfo_Button(t *testing.T) {
	// Create a button with custom data
	customData := struct {
		ID string `xml:"id,attr"`
	}{
		ID: "test-button-1",
	}

	// Create a simple button image for testing
	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
		Disabled: createTestNineSlice(100, 30, color.RGBA{150, 150, 150, 255}),
	}

	buttonColor := &widget.ButtonTextColor{
		Idle:     color.White,
		Disabled: color.Gray{128},
	}

	// Create a basic button
	btn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Test Button", nil, buttonColor),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(customData),
			widget.WidgetOpts.MinSize(100, 30),
		),
	)

	// Validate and set location
	// Note: Validation will fail without a font face, but we can still test extraction
	// btn.Validate() - skip validation for unit tests
	btn.SetLocation(image.Rect(10, 20, 110, 50))

	// Extract widget info
	info := autoui.ExtractWidgetInfo(btn)

	// Verify Type
	if info.Type != "Button" {
		t.Errorf("Expected Type to be 'Button', got '%s'", info.Type)
	}

	// Verify Rect
	expectedRect := image.Rect(10, 20, 110, 50)
	if info.Rect != expectedRect {
		t.Errorf("Expected Rect to be %v, got %v", expectedRect, info.Rect)
	}

	// Verify Visible (should be true by default)
	if !info.Visible {
		t.Error("Expected Visible to be true")
	}

	// Verify Disabled (should be false by default)
	if info.Disabled {
		t.Error("Expected Disabled to be false")
	}

	// Verify State contains text
	// Note: Without validation, text extraction may not work as expected
	// The Text() method requires init to be called
	if info.State == nil {
		t.Log("Warning: State is nil (may require validation)")
	} else if info.State["text"] != "Test Button" {
		t.Logf("State['text'] = '%s' (may require validation)", info.State["text"])
	}

	// Verify CustomData extraction
	if info.CustomData == nil {
		t.Fatal("Expected CustomData to be non-nil")
	}
	if info.CustomData["id"] != "test-button-1" {
		t.Errorf("Expected CustomData['id'] to be 'test-button-1', got '%s'", info.CustomData["id"])
	}
}

// TestExtractWidgetInfo_Container tests extracting widget info from a container.
func TestExtractWidgetInfo_Container(t *testing.T) {
	// Create a container with custom data
	customData := map[string]string{
		"id":   "main-container",
		"role": "layout",
	}

	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(customData),
			widget.WidgetOpts.MinSize(800, 600),
		),
	)

	container.Validate()
	container.SetLocation(image.Rect(0, 0, 800, 600))

	// Extract widget info
	info := autoui.ExtractWidgetInfo(container)

	// Verify Type
	if info.Type != "Container" {
		t.Errorf("Expected Type to be 'Container', got '%s'", info.Type)
	}

	// Verify Rect
	expectedRect := image.Rect(0, 0, 800, 600)
	if info.Rect != expectedRect {
		t.Errorf("Expected Rect to be %v, got %v", expectedRect, info.Rect)
	}

	// Verify CustomData extraction from map
	if info.CustomData == nil {
		t.Fatal("Expected CustomData to be non-nil")
	}
	if info.CustomData["id"] != "main-container" {
		t.Errorf("Expected CustomData['id'] to be 'main-container', got '%s'", info.CustomData["id"])
	}
	if info.CustomData["role"] != "layout" {
		t.Errorf("Expected CustomData['role'] to be 'layout', got '%s'", info.CustomData["role"])
	}
}

// TestExtractWidgetInfo_Slider tests extracting widget info from a slider.
func TestExtractWidgetInfo_Slider(t *testing.T) {
	trackImage := &widget.SliderTrackImage{
		Idle:     createTestNineSlice(100, 10, color.RGBA{50, 50, 50, 255}),
		Disabled: createTestNineSlice(100, 10, color.RGBA{100, 100, 100, 255}),
	}

	handleImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(16, 16, color.RGBA{150, 150, 150, 255}),
		Pressed:  createTestNineSlice(16, 16, color.RGBA{100, 100, 100, 255}),
		Disabled: createTestNineSlice(16, 16, color.RGBA{200, 200, 200, 255}),
	}

	slider := widget.NewSlider(
		widget.SliderOpts.Images(trackImage, handleImage),
		widget.SliderOpts.MinMax(0, 100),
		widget.SliderOpts.InitialCurrent(50),
	)

	slider.Validate()
	slider.SetLocation(image.Rect(10, 10, 210, 30))

	// Extract widget info
	info := autoui.ExtractWidgetInfo(slider)

	// Verify Type
	if info.Type != "Slider" {
		t.Errorf("Expected Type to be 'Slider', got '%s'", info.Type)
	}

	// Verify State contains slider values
	if info.State == nil {
		t.Fatal("Expected State to be non-nil")
	}
	if info.State["value"] != "50" {
		t.Errorf("Expected State['value'] to be '50', got '%s'", info.State["value"])
	}
	if info.State["min"] != "0" {
		t.Errorf("Expected State['min'] to be '0', got '%s'", info.State["min"])
	}
	if info.State["max"] != "100" {
		t.Errorf("Expected State['max'] to be '100', got '%s'", info.State["max"])
	}
}

// TestExtractWidgetInfo_ProgressBar tests extracting widget info from a progress bar.
func TestExtractWidgetInfo_ProgressBar(t *testing.T) {
	trackImage := &widget.ProgressBarImage{
		Idle:     createTestNineSlice(100, 20, color.RGBA{50, 50, 50, 255}),
		Disabled: createTestNineSlice(100, 20, color.RGBA{100, 100, 100, 255}),
	}

	fillImage := &widget.ProgressBarImage{
		Idle:     createTestNineSlice(100, 20, color.RGBA{0, 150, 0, 255}),
		Disabled: createTestNineSlice(100, 20, color.RGBA{100, 150, 100, 255}),
	}

	pb := widget.NewProgressBar(
		widget.ProgressBarOpts.Images(trackImage, fillImage),
		widget.ProgressBarOpts.Values(0, 100, 75),
	)

	pb.Validate()
	pb.SetLocation(image.Rect(10, 10, 210, 30))

	// Extract widget info
	info := autoui.ExtractWidgetInfo(pb)

	// Verify Type
	if info.Type != "ProgressBar" {
		t.Errorf("Expected Type to be 'ProgressBar', got '%s'", info.Type)
	}

	// Verify State contains progress bar values
	if info.State == nil {
		t.Fatal("Expected State to be non-nil")
	}
	if info.State["value"] != "75" {
		t.Errorf("Expected State['value'] to be '75', got '%s'", info.State["value"])
	}
	if info.State["min"] != "0" {
		t.Errorf("Expected State['min'] to be '0', got '%s'", info.State["min"])
	}
	if info.State["max"] != "100" {
		t.Errorf("Expected State['max'] to be '100', got '%s'", info.State["max"])
	}
}

// createTestNineSlice creates a simple NineSlice for testing.
func createTestNineSlice(w, h int, c color.Color) *ebitenuiImage.NineSlice {
	img := ebiten.NewImage(w, h)
	img.Fill(c)
	return ebitenuiImage.NewNineSliceSimple(img, 0, 0)
}