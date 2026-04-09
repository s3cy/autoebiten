package internal_test

import (
	"image"
	"image/color"
	"testing"

	ebitenuiImage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/s3cy/autoebiten/autoui/internal"
)

// TestExtractWidgetState_Button tests button state extraction.
func TestExtractWidgetState_Button(t *testing.T) {
	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
		Disabled: createTestNineSlice(100, 30, color.RGBA{150, 150, 150, 255}),
	}

	buttonColor := &widget.ButtonTextColor{
		Idle:     color.White,
		Disabled: color.Gray{128},
	}

	btn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Click Me", nil, buttonColor),
		widget.ButtonOpts.ToggleMode(),
	)
	btn.SetLocation(image.Rect(0, 0, 100, 30))

	state := internal.ExtractWidgetState(btn)
	if state == nil {
		t.Fatal("Expected non-nil state")
	}

	if state["toggle"] != "true" {
		t.Errorf("Expected toggle='true', got '%s'", state["toggle"])
	}
}

// TestExtractWidgetState_ProgressBar tests progress bar state extraction.
func TestExtractWidgetState_ProgressBar(t *testing.T) {
	trackImage := &widget.ProgressBarImage{
		Idle: createTestNineSlice(100, 20, color.RGBA{50, 50, 50, 255}),
	}

	fillImage := &widget.ProgressBarImage{
		Idle: createTestNineSlice(100, 20, color.RGBA{0, 150, 0, 255}),
	}

	pb := widget.NewProgressBar(
		widget.ProgressBarOpts.Images(trackImage, fillImage),
		widget.ProgressBarOpts.Values(0, 200, 150),
	)
	pb.Validate()
	pb.SetLocation(image.Rect(0, 0, 100, 20))

	state := internal.ExtractWidgetState(pb)
	if state == nil {
		t.Fatal("Expected non-nil state")
	}

	if state["value"] != "150" {
		t.Errorf("Expected value='150', got '%s'", state["value"])
	}
	if state["min"] != "0" {
		t.Errorf("Expected min='0', got '%s'", state["min"])
	}
	if state["max"] != "200" {
		t.Errorf("Expected max='200', got '%s'", state["max"])
	}
}

// createTestNineSlice creates a simple NineSlice for testing.
func createTestNineSlice(w, h int, c color.Color) *ebitenuiImage.NineSlice {
	img := ebiten.NewImage(w, h)
	img.Fill(c)
	return ebitenuiImage.NewNineSliceSimple(img, 0, 0)
}

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
		widget.ListOpts.HideVerticalSlider(),
		widget.ListOpts.HideHorizontalSlider(),
	)
	// Note: List.Validate() requires EntryFontFace which needs font loading
	// Testing extraction on unvalidated widget - Entries() should still work
	list.SetLocation(image.Rect(0, 0, 100, 100))

	state := internal.ExtractWidgetState(list)
	if state == nil {
		t.Fatal("Expected non-nil state")
	}

	if state["entries"] != "3" {
		t.Errorf("Expected entries='3', got '%s'", state["entries"])
	}
}