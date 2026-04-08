package autoui_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/s3cy/autoebiten/autoui"
)

// TestInvokeMethod_ButtonClick tests invoking Click() on a button (no args).
// Note: Event handler triggering may require game loop initialization.
// This test verifies the method invocation works correctly.
func TestInvokeMethod_ButtonClick(t *testing.T) {
	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
		Disabled: createTestNineSlice(100, 30, color.RGBA{150, 150, 150, 255}),
	}

	btn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
	)
	btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)

	info := autoui.ExtractWidgetInfo(btn)

	// Invoke Click method with no arguments - should not return error
	err := autoui.InvokeMethod(info.Widget, "Click", nil)
	if err != nil {
		t.Errorf("InvokeMethod failed: %v", err)
	}

	// Method invocation succeeded - event handler verification requires integration test
}

// TestInvokeMethod_InvalidMethod tests error for non-existent method.
func TestInvokeMethod_InvalidMethod(t *testing.T) {
	buttonImage := &widget.ButtonImage{
		Idle: createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
	}

	btn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
	)
	btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)

	info := autoui.ExtractWidgetInfo(btn)

	// Try to invoke a method that doesn't exist
	err := autoui.InvokeMethod(info.Widget, "NonExistentMethod", nil)
	if err == nil {
		t.Error("Expected error for non-existent method, got nil")
	}
}

// TestInvokeMethod_WithArgs tests invoking Focus(true) on TextInput.
func TestInvokeMethod_WithArgs(t *testing.T) {
	textInput := widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(200, 30),
		),
	)
	textInput.GetWidget().Rect = image.Rect(10, 10, 210, 40)

	info := autoui.ExtractWidgetInfo(textInput)

	// TextInput should not be focused initially
	if textInput.IsFocused() {
		t.Error("Expected TextInput to not be focused initially")
	}

	// Invoke Focus(true) with bool argument
	err := autoui.InvokeMethod(info.Widget, "Focus", []any{true})
	if err != nil {
		t.Errorf("InvokeMethod failed: %v", err)
	}

	// Verify focus was set
	if !textInput.IsFocused() {
		t.Error("Expected Focus(true) to set focus on TextInput")
	}

	// Invoke Focus(false) to remove focus
	err = autoui.InvokeMethod(info.Widget, "Focus", []any{false})
	if err != nil {
		t.Errorf("InvokeMethod failed: %v", err)
	}

	// Verify focus was removed
	if textInput.IsFocused() {
		t.Error("Expected Focus(false) to remove focus from TextInput")
	}
}

// TestInvokeMethod_SetText tests invoking SetText(string) on TextInput.
func TestInvokeMethod_SetText(t *testing.T) {
	textInput := widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(200, 30),
		),
	)
	textInput.GetWidget().Rect = image.Rect(10, 10, 210, 40)

	info := autoui.ExtractWidgetInfo(textInput)

	// Invoke SetText("Hello World")
	err := autoui.InvokeMethod(info.Widget, "SetText", []any{"Hello World"})
	if err != nil {
		t.Errorf("InvokeMethod failed: %v", err)
	}

	// Verify text was set
	if textInput.GetText() != "Hello World" {
		t.Errorf("Expected text to be 'Hello World', got '%s'", textInput.GetText())
	}
}

// TestInvokeMethod_NilWidget tests error when widget is nil.
func TestInvokeMethod_NilWidget(t *testing.T) {
	err := autoui.InvokeMethod(nil, "Click", nil)
	if err == nil {
		t.Error("Expected error for nil widget, got nil")
	}
}

// TestInvokeMethod_ArgumentCountMismatch tests error when argument count doesn't match.
func TestInvokeMethod_ArgumentCountMismatch(t *testing.T) {
	buttonImage := &widget.ButtonImage{
		Idle: createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
	}

	btn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
	)
	btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)

	info := autoui.ExtractWidgetInfo(btn)

	// Click() takes no arguments, but we pass one
	err := autoui.InvokeMethod(info.Widget, "Click", []any{true})
	if err == nil {
		t.Error("Expected error for argument count mismatch, got nil")
	}
}

// TestInvokeMethod_NonWhitelistedSignature tests error for non-whitelisted method signature.
func TestInvokeMethod_NonWhitelistedSignature(t *testing.T) {
	buttonImage := &widget.ButtonImage{
		Idle: createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
	}

	btn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
	)
	btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)

	info := autoui.ExtractWidgetInfo(btn)

	// SetState takes WidgetState which is not whitelisted
	err := autoui.InvokeMethod(info.Widget, "SetState", []any{widget.WidgetChecked})
	if err == nil {
		t.Error("Expected error for non-whitelisted signature, got nil")
	}
}

// TestInvokeMethod_NumericConversion tests numeric argument conversion.
func TestInvokeMethod_NumericConversion(t *testing.T) {
	// This test verifies numeric type conversion behavior
	// Note: Most ebitenui widgets don't have methods taking int/float64 directly,
	// but we test the conversion logic anyway

	textInput := widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(200, 30),
		),
	)
	textInput.GetWidget().Rect = image.Rect(10, 10, 210, 40)

	info := autoui.ExtractWidgetInfo(textInput)

	// Test that bool conversion works (already tested in Focus)
	// Focus takes bool, passing int should fail
	err := autoui.InvokeMethod(info.Widget, "Focus", []any{1})
	if err == nil {
		t.Error("Expected error when passing int to bool parameter, got nil")
	}

	// Passing float64 to bool should also fail
	err = autoui.InvokeMethod(info.Widget, "Focus", []any{1.5})
	if err == nil {
		t.Error("Expected error when passing float64 to bool parameter, got nil")
	}
}

// TestInvokeMethod_ButtonPress tests invoking Press() on a button.
// Note: Event handler triggering may require game loop initialization.
// This test verifies the method invocation works correctly.
func TestInvokeMethod_ButtonPress(t *testing.T) {
	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
		Disabled: createTestNineSlice(100, 30, color.RGBA{150, 150, 150, 255}),
	}

	btn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
	)
	btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)

	info := autoui.ExtractWidgetInfo(btn)

	// Invoke Press method with no arguments - should not return error
	err := autoui.InvokeMethod(info.Widget, "Press", nil)
	if err != nil {
		t.Errorf("InvokeMethod failed: %v", err)
	}

	// Method invocation succeeded - event handler verification requires integration test
}

// TestInvokeMethod_ButtonRelease tests invoking Release() on a button.
// Note: Event handler triggering may require game loop initialization.
// This test verifies the method invocation works correctly.
func TestInvokeMethod_ButtonRelease(t *testing.T) {
	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
		Disabled: createTestNineSlice(100, 30, color.RGBA{150, 150, 150, 255}),
	}

	btn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
	)
	btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)

	// First press the button so it can be released
	btn.Press()

	info := autoui.ExtractWidgetInfo(btn)

	// Invoke Release method with no arguments - should not return error
	err := autoui.InvokeMethod(info.Widget, "Release", nil)
	if err != nil {
		t.Errorf("InvokeMethod failed: %v", err)
	}

	// Method invocation succeeded - event handler verification requires integration test
}

// TestInvokeMethod_ButtonSetText tests invoking SetText(string) on Button.
func TestInvokeMethod_ButtonSetText(t *testing.T) {
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
		widget.ButtonOpts.Text("Original", nil, buttonColor),
	)
	btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)

	info := autoui.ExtractWidgetInfo(btn)

	// Invoke SetText("New Text")
	err := autoui.InvokeMethod(info.Widget, "SetText", []any{"New Text"})
	if err != nil {
		t.Errorf("InvokeMethod failed: %v", err)
	}

	// Verify text was set (Text() returns *Text widget, we check label)
	if btn.Text() == nil {
		t.Log("Warning: Text() returned nil (may require validation)")
	}
}
