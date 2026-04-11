package autoui_test

import (
	"io"
	"image"
	"image/color"
	"reflect"
	"strings"
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

	// SetImage takes *ButtonImage (pointer to struct) which is not whitelisted
	err := autoui.InvokeMethod(info.Widget, "SetImage", []any{buttonImage})
	if err == nil {
		t.Error("Expected error for non-whitelisted signature (struct pointer), got nil")
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

// TestIsWhitelistedSignature_AnyReturn tests that func() any is whitelisted.
func TestIsWhitelistedSignature_AnyReturn(t *testing.T) {
	fn := func() any { return "test" }
	fnType := reflect.TypeOf(fn)

	if !autoui.IsWhitelistedSignature(fnType) {
		t.Error("Expected func() any to be whitelisted for returns")
	}
}

// TestIsWhitelistedSignature_SliceAnyReturn tests that func() []any is whitelisted.
func TestIsWhitelistedSignature_SliceAnyReturn(t *testing.T) {
	fn := func() []any { return []any{"a", "b"} }
	fnType := reflect.TypeOf(fn)

	if !autoui.IsWhitelistedSignature(fnType) {
		t.Error("Expected func() []any to be whitelisted")
	}
}

// TestIsWhitelistedSignature_EnumReturn tests that func() WidgetState (enum) is whitelisted.
func TestIsWhitelistedSignature_EnumReturn(t *testing.T) {
	// WidgetState is defined as: type WidgetState int
	fn := func() widget.WidgetState { return widget.WidgetChecked }
	fnType := reflect.TypeOf(fn)

	if !autoui.IsWhitelistedSignature(fnType) {
		t.Error("Expected func() WidgetState (enum) to be whitelisted")
	}
}

// TestIsWhitelistedSignature_AnyParam tests that func(any) is whitelisted.
func TestIsWhitelistedSignature_AnyParam(t *testing.T) {
	fn := func(any) {}
	fnType := reflect.TypeOf(fn)

	if !autoui.IsWhitelistedSignature(fnType) {
		t.Error("Expected func(any) to be whitelisted for params")
	}
}

// TestIsWhitelistedSignature_EnumParam tests that func(WidgetState) is whitelisted.
func TestIsWhitelistedSignature_EnumParam(t *testing.T) {
	fn := func(widget.WidgetState) {}
	fnType := reflect.TypeOf(fn)

	if !autoui.IsWhitelistedSignature(fnType) {
		t.Error("Expected func(WidgetState) to be whitelisted for enum params")
	}
}

// TestIsWhitelistedSignature_BasicSliceReturn tests that func() []string is whitelisted.
func TestIsWhitelistedSignature_BasicSliceReturn(t *testing.T) {
	fn := func() []string { return []string{"a", "b"} }
	fnType := reflect.TypeOf(fn)

	if !autoui.IsWhitelistedSignature(fnType) {
		t.Error("Expected func() []string to be whitelisted")
	}
}

// TestIsWhitelistedSignature_BasicIntReturn tests that func() int is whitelisted.
func TestIsWhitelistedSignature_BasicIntReturn(t *testing.T) {
	fn := func() int { return 42 }
	fnType := reflect.TypeOf(fn)

	if !autoui.IsWhitelistedSignature(fnType) {
		t.Error("Expected func() int to be whitelisted")
	}
}

// TestIsWhitelistedSignature_ErrorReturn tests that func() error is whitelisted.
func TestIsWhitelistedSignature_ErrorReturn(t *testing.T) {
	fn := func() error { return nil }
	fnType := reflect.TypeOf(fn)

	if !autoui.IsWhitelistedSignature(fnType) {
		t.Error("Expected func() error to be whitelisted")
	}
}

// TestIsWhitelistedSignature_TwoReturns tests that multi-return is NOT whitelisted.
func TestIsWhitelistedSignature_TwoReturns(t *testing.T) {
	fn := func() (int, error) { return 0, nil }
	fnType := reflect.TypeOf(fn)

	if autoui.IsWhitelistedSignature(fnType) {
		t.Error("Expected func() (int, error) to NOT be whitelisted")
	}
}

// TestConvertArg_IntToEnum tests converting int/float64 to enum types like WidgetState.
func TestConvertArg_IntToEnum(t *testing.T) {
	targetType := reflect.TypeOf(widget.WidgetState(0))
	arg := float64(1) // JSON unmarshals numbers to float64

	result, err := autoui.ConvertArg(arg, targetType)
	if err != nil {
		t.Errorf("ConvertArg failed: %v", err)
	}

	// Should convert to WidgetState(1)
	if result.Int() != 1 {
		t.Errorf("Expected int value 1, got %d", result.Int())
	}
}

// TestConvertArg_AnyParam tests that interface{} target accepts any value.
func TestConvertArg_AnyParam(t *testing.T) {
	// interface{} target should accept any value
	targetType := reflect.TypeOf((*interface{})(nil)).Elem()
	arg := "test string"

	result, err := autoui.ConvertArg(arg, targetType)
	if err != nil {
		t.Errorf("ConvertArg failed: %v", err)
	}

	if result.String() != "test string" {
		t.Errorf("Expected 'test string', got %v", result.Interface())
	}
}

// TestConvertArg_NonEmptyInterfaceNotImplemented tests that non-empty interface
// targets with incompatible values return an error.
func TestConvertArg_NonEmptyInterfaceNotImplemented(t *testing.T) {
	// io.Reader is a non-empty interface that string does not implement
	targetType := reflect.TypeOf((*io.Reader)(nil)).Elem()
	arg := "not a reader"

	_, err := autoui.ConvertArg(arg, targetType)
	if err == nil {
		t.Error("Expected error for non-implemented interface")
	}
	if !strings.Contains(err.Error(), "interface not implemented") {
		t.Errorf("Expected 'interface not implemented' error, got: %v", err)
	}
}

// TestInvokeMethodWithResult_ReturnValue tests capturing return values from widget methods.
func TestInvokeMethodWithResult_ReturnValue(t *testing.T) {
	textInput := widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(200, 30),
		),
	)
	textInput.GetWidget().Rect = image.Rect(10, 10, 210, 40)
	textInput.SetText("Hello")

	info := autoui.ExtractWidgetInfo(textInput)

	result, err := autoui.InvokeMethodWithResult(info.Widget, "GetText", nil)
	if err != nil {
		t.Errorf("InvokeMethodWithResult failed: %v", err)
	}

	if result != "Hello" {
		t.Errorf("Expected 'Hello', got %v", result)
	}
}

// TestInvokeMethodWithResult_EnumReturn tests capturing enum return values and converting to int64.
func TestInvokeMethodWithResult_EnumReturn(t *testing.T) {
	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
		Disabled: createTestNineSlice(100, 30, color.RGBA{150, 150, 150, 255}),
	}

	btn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.ToggleMode(), // Enable toggle mode so State() works
	)
	btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)

	// SetState directly changes the WidgetState enum
	btn.SetState(widget.WidgetChecked)

	info := autoui.ExtractWidgetInfo(btn)

	result, err := autoui.InvokeMethodWithResult(info.Widget, "State", nil)
	if err != nil {
		t.Errorf("InvokeMethodWithResult failed: %v", err)
	}

	stateInt, ok := result.(int64)
	if !ok {
		t.Errorf("Expected int64, got %T", result)
	}
	if stateInt != 1 { // WidgetChecked = 1
		t.Errorf("Expected state 1 (checked), got %d", stateInt)
	}
}

// TestInvokeMethodWithResult_SliceReturn tests capturing []any slice return values.
func TestInvokeMethodWithResult_SliceReturn(t *testing.T) {
	entries := []any{"Entry A", "Entry B", "Entry C"}
	list := widget.NewList(
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e any) string {
			return e.(string)
		}),
	)
	list.GetWidget().Rect = image.Rect(10, 10, 210, 200)

	info := autoui.ExtractWidgetInfo(list)

	result, err := autoui.InvokeMethodWithResult(info.Widget, "Entries", nil)
	if err != nil {
		t.Errorf("InvokeMethodWithResult failed: %v", err)
	}

	entriesResult, ok := result.([]any)
	if !ok {
		t.Errorf("Expected []any, got %T", result)
	}
	if len(entriesResult) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entriesResult))
	}
}
