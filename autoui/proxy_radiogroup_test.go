package autoui

import (
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"github.com/ebitenui/ebitenui/widget"
)

// TestInvokeRadioGroupMethod tests the InvokeRadioGroupMethod function.
func TestInvokeRadioGroupMethod(t *testing.T) {
	// Create buttons for radio group
	btn1 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn2 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn3 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	setButtonText(btn1, "Option 1")
	setButtonText(btn2, "Option 2")
	setButtonText(btn3, "Option 3")

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2, btn3))

	// Test Elements method
	result, err := InvokeRadioGroupMethod(rg, "Elements", nil)
	if err != nil {
		t.Fatalf("InvokeRadioGroupMethod failed: %v", err)
	}

	elements, ok := result.([]RadioGroupElementInfo)
	if !ok {
		t.Fatalf("expected []RadioGroupElementInfo, got %T", result)
	}

	if len(elements) != 3 {
		t.Errorf("expected 3 elements, got %d", len(elements))
	}
}

// TestInvokeRadioGroupMethodUnknown tests InvokeRadioGroupMethod with unknown method.
func TestInvokeRadioGroupMethodUnknown(t *testing.T) {
	btn := widget.NewButton(widget.ButtonOpts.ToggleMode())
	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn))

	_, err := InvokeRadioGroupMethod(rg, "UnknownMethod", nil)
	if err == nil {
		t.Fatal("expected error for unknown method")
	}

	if !strings.Contains(err.Error(), "unknown RadioGroup method") {
		t.Errorf("expected error to contain 'unknown RadioGroup method', got: %s", err.Error())
	}
}

// TestHandleRadioGroupElements tests the Elements handler.
func TestHandleRadioGroupElements(t *testing.T) {
	// Create buttons with labels using reflection
	btn1 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn2 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	setButtonText(btn1, "Option A")
	setButtonText(btn2, "Option B")

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2))

	result, err := handleRadioGroupElements(rg, nil)
	if err != nil {
		t.Fatalf("handleRadioGroupElements failed: %v", err)
	}

	elements, ok := result.([]RadioGroupElementInfo)
	if !ok {
		t.Fatalf("expected []RadioGroupElementInfo, got %T", result)
	}

	if len(elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(elements))
	}

	// Verify element info
	if elements[0].Type != "Button" {
		t.Errorf("element 0: expected type 'Button', got '%s'", elements[0].Type)
	}
	if elements[0].Label != "Option A" {
		t.Errorf("element 0: expected label 'Option A', got '%s'", elements[0].Label)
	}

	if elements[1].Type != "Button" {
		t.Errorf("element 1: expected type 'Button', got '%s'", elements[1].Type)
	}
	if elements[1].Label != "Option B" {
		t.Errorf("element 1: expected label 'Option B', got '%s'", elements[1].Label)
	}
}

// TestHandleRadioGroupElementsCheckboxes tests Elements handler with Checkbox elements.
func TestHandleRadioGroupElementsCheckboxes(t *testing.T) {
	// Create checkboxes with labels using reflection
	cb1 := widget.NewCheckbox()
	cb2 := widget.NewCheckbox()
	setCheckboxText(cb1, "Check A")
	setCheckboxText(cb2, "Check B")

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(cb1, cb2))

	result, err := handleRadioGroupElements(rg, nil)
	if err != nil {
		t.Fatalf("handleRadioGroupElements failed: %v", err)
	}

	elements, ok := result.([]RadioGroupElementInfo)
	if !ok {
		t.Fatalf("expected []RadioGroupElementInfo, got %T", result)
	}

	if len(elements) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(elements))
	}

	// Verify element info
	if elements[0].Type != "Checkbox" {
		t.Errorf("element 0: expected type 'Checkbox', got '%s'", elements[0].Type)
	}
	if elements[0].Label != "Check A" {
		t.Errorf("element 0: expected label 'Check A', got '%s'", elements[0].Label)
	}

	if elements[1].Type != "Checkbox" {
		t.Errorf("element 1: expected type 'Checkbox', got '%s'", elements[1].Type)
	}
	if elements[1].Label != "Check B" {
		t.Errorf("element 1: expected label 'Check B', got '%s'", elements[1].Label)
	}
}

// TestHandleRadioGroupElementsEmpty tests Elements handler with empty RadioGroup.
func TestHandleRadioGroupElementsEmpty(t *testing.T) {
	rg := widget.NewRadioGroup()

	result, err := handleRadioGroupElements(rg, nil)
	if err != nil {
		t.Fatalf("handleRadioGroupElements failed: %v", err)
	}

	elements, ok := result.([]RadioGroupElementInfo)
	if !ok {
		t.Fatalf("expected []RadioGroupElementInfo, got %T", result)
	}

	if len(elements) != 0 {
		t.Errorf("expected 0 elements for empty RadioGroup, got %d", len(elements))
	}
}

// TestHandleRadioGroupElementsActive tests that the Active field is set correctly.
func TestHandleRadioGroupElementsActive(t *testing.T) {
	btn1 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn2 := widget.NewButton(widget.ButtonOpts.ToggleMode())

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2))

	// Set active to btn2
	rg.SetActive(btn2)

	result, err := handleRadioGroupElements(rg, nil)
	if err != nil {
		t.Fatalf("handleRadioGroupElements failed: %v", err)
	}

	elements := result.([]RadioGroupElementInfo)

	if elements[0].Active {
		t.Error("element 0 should not be active")
	}
	if !elements[1].Active {
		t.Error("element 1 should be active")
	}
}

// TestHandleRadioGroupElementsWrongArgs tests Elements handler with wrong arguments.
func TestHandleRadioGroupElementsWrongArgs(t *testing.T) {
	btn := widget.NewButton(widget.ButtonOpts.ToggleMode())
	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn))

	_, err := handleRadioGroupElements(rg, []any{"unexpected"})
	if err == nil {
		t.Fatal("expected error for non-empty arguments")
	}
}

// TestHandleRadioGroupActiveIndex tests the ActiveIndex handler.
func TestHandleRadioGroupActiveIndex(t *testing.T) {
	btn1 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn2 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn3 := widget.NewButton(widget.ButtonOpts.ToggleMode())

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2, btn3))

	// Set active to btn2 (index 1)
	rg.SetActive(btn2)

	result, err := handleRadioGroupActiveIndex(rg, nil)
	if err != nil {
		t.Fatalf("handleRadioGroupActiveIndex failed: %v", err)
	}

	index, ok := result.(int64)
	if !ok {
		t.Fatalf("expected int64, got %T", result)
	}

	if index != 1 {
		t.Errorf("expected index 1, got %d", index)
	}
}

// TestHandleRadioGroupActiveIndexNoSelection tests ActiveIndex when no element is selected.
func TestHandleRadioGroupActiveIndexNoSelection(t *testing.T) {
	btn1 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn2 := widget.NewButton(widget.ButtonOpts.ToggleMode())

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2))
	// Don't set any active element

	result, err := handleRadioGroupActiveIndex(rg, nil)
	if err != nil {
		t.Fatalf("handleRadioGroupActiveIndex failed: %v", err)
	}

	index, ok := result.(int64)
	if !ok {
		t.Fatalf("expected int64, got %T", result)
	}

	if index != -1 {
		t.Errorf("expected index -1 for no selection, got %d", index)
	}
}

// TestHandleRadioGroupActiveIndexWrongArgs tests ActiveIndex with wrong arguments.
func TestHandleRadioGroupActiveIndexWrongArgs(t *testing.T) {
	btn := widget.NewButton(widget.ButtonOpts.ToggleMode())
	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn))

	_, err := handleRadioGroupActiveIndex(rg, []any{"unexpected"})
	if err == nil {
		t.Fatal("expected error for non-empty arguments")
	}
}

// TestHandleRadioGroupSetActiveByIndex tests setting active element by index.
func TestHandleRadioGroupSetActiveByIndex(t *testing.T) {
	btn1 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn2 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn3 := widget.NewButton(widget.ButtonOpts.ToggleMode())

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2, btn3))

	// Set active to index 2
	_, err := handleRadioGroupSetActiveByIndex(rg, []any{float64(2)})
	if err != nil {
		t.Fatalf("handleRadioGroupSetActiveByIndex failed: %v", err)
	}

	// Verify the active element
	active := rg.Active()
	if active != btn3 {
		t.Error("active element should be btn3")
	}
}

// TestHandleRadioGroupSetActiveByIndexOutOfRange tests setting active with invalid index.
func TestHandleRadioGroupSetActiveByIndexOutOfRange(t *testing.T) {
	btn := widget.NewButton(widget.ButtonOpts.ToggleMode())
	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn))

	// Try to set index 5 (out of range)
	_, err := handleRadioGroupSetActiveByIndex(rg, []any{float64(5)})
	if err == nil {
		t.Fatal("expected error for out of range index")
	}

	if !strings.Contains(err.Error(), "out of range") {
		t.Errorf("expected error to contain 'out of range', got: %s", err.Error())
	}
}

// TestHandleRadioGroupSetActiveByIndexNegative tests setting active with negative index.
func TestHandleRadioGroupSetActiveByIndexNegative(t *testing.T) {
	btn := widget.NewButton(widget.ButtonOpts.ToggleMode())
	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn))

	// Try to set index -1 (negative)
	_, err := handleRadioGroupSetActiveByIndex(rg, []any{float64(-1)})
	if err == nil {
		t.Fatal("expected error for negative index")
	}

	if !strings.Contains(err.Error(), "out of range") {
		t.Errorf("expected error to contain 'out of range', got: %s", err.Error())
	}
}

// TestHandleRadioGroupSetActiveByIndexEmpty tests setting active on empty RadioGroup.
func TestHandleRadioGroupSetActiveByIndexEmpty(t *testing.T) {
	rg := widget.NewRadioGroup()

	_, err := handleRadioGroupSetActiveByIndex(rg, []any{float64(0)})
	if err == nil {
		t.Fatal("expected error for empty RadioGroup")
	}

	if !strings.Contains(err.Error(), "no elements") {
		t.Errorf("expected error to contain 'no elements', got: %s", err.Error())
	}
}

// TestHandleRadioGroupSetActiveByIndexWrongArgs tests SetActiveByIndex with wrong arguments.
func TestHandleRadioGroupSetActiveByIndexWrongArgs(t *testing.T) {
	btn := widget.NewButton(widget.ButtonOpts.ToggleMode())
	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn))

	// Test with no arguments
	_, err := handleRadioGroupSetActiveByIndex(rg, nil)
	if err == nil {
		t.Fatal("expected error for no arguments")
	}

	// Test with wrong type
	_, err = handleRadioGroupSetActiveByIndex(rg, []any{"not a number"})
	if err == nil {
		t.Fatal("expected error for wrong argument type")
	}
}

// TestHandleRadioGroupActiveLabel tests getting the label of the active element.
func TestHandleRadioGroupActiveLabel(t *testing.T) {
	btn1 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn2 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	setButtonText(btn1, "First Option")
	setButtonText(btn2, "Second Option")

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2))

	// Set active to btn2
	rg.SetActive(btn2)

	result, err := handleRadioGroupActiveLabel(rg, nil)
	if err != nil {
		t.Fatalf("handleRadioGroupActiveLabel failed: %v", err)
	}

	label, ok := result.(string)
	if !ok {
		t.Fatalf("expected string, got %T", result)
	}

	if label != "Second Option" {
		t.Errorf("expected label 'Second Option', got '%s'", label)
	}
}

// TestHandleRadioGroupActiveLabelCheckbox tests ActiveLabel with Checkbox elements.
func TestHandleRadioGroupActiveLabelCheckbox(t *testing.T) {
	cb1 := widget.NewCheckbox()
	cb2 := widget.NewCheckbox()
	setCheckboxText(cb1, "Check One")
	setCheckboxText(cb2, "Check Two")

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(cb1, cb2))

	// Set active to cb1
	rg.SetActive(cb1)

	result, err := handleRadioGroupActiveLabel(rg, nil)
	if err != nil {
		t.Fatalf("handleRadioGroupActiveLabel failed: %v", err)
	}

	label, ok := result.(string)
	if !ok {
		t.Fatalf("expected string, got %T", result)
	}

	if label != "Check One" {
		t.Errorf("expected label 'Check One', got '%s'", label)
	}
}

// TestHandleRadioGroupActiveLabelNoSelection tests ActiveLabel when no element is selected.
func TestHandleRadioGroupActiveLabelNoSelection(t *testing.T) {
	btn := widget.NewButton(widget.ButtonOpts.ToggleMode())
	setButtonText(btn, "Option")

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn))
	// Don't set any active element

	result, err := handleRadioGroupActiveLabel(rg, nil)
	if err != nil {
		t.Fatalf("handleRadioGroupActiveLabel failed: %v", err)
	}

	label, ok := result.(string)
	if !ok {
		t.Fatalf("expected string, got %T", result)
	}

	if label != "" {
		t.Errorf("expected empty label for no selection, got '%s'", label)
	}
}

// TestHandleRadioGroupActiveLabelWrongArgs tests ActiveLabel with wrong arguments.
func TestHandleRadioGroupActiveLabelWrongArgs(t *testing.T) {
	btn := widget.NewButton(widget.ButtonOpts.ToggleMode())
	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn))

	_, err := handleRadioGroupActiveLabel(rg, []any{"unexpected"})
	if err == nil {
		t.Fatal("expected error for non-empty arguments")
	}
}

// TestHandleRadioGroupSetActiveByLabel tests setting active element by label.
func TestHandleRadioGroupSetActiveByLabel(t *testing.T) {
	btn1 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn2 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn3 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	setButtonText(btn1, "Alpha")
	setButtonText(btn2, "Beta")
	setButtonText(btn3, "Gamma")

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2, btn3))

	// Set active to "Beta"
	_, err := handleRadioGroupSetActiveByLabel(rg, []any{"Beta"})
	if err != nil {
		t.Fatalf("handleRadioGroupSetActiveByLabel failed: %v", err)
	}

	// Verify the active element
	active := rg.Active()
	if active != btn2 {
		t.Error("active element should be btn2")
	}
}

// TestHandleRadioGroupSetActiveByLabelCheckbox tests SetActiveByLabel with Checkbox elements.
func TestHandleRadioGroupSetActiveByLabelCheckbox(t *testing.T) {
	cb1 := widget.NewCheckbox()
	cb2 := widget.NewCheckbox()
	setCheckboxText(cb1, "Option X")
	setCheckboxText(cb2, "Option Y")

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(cb1, cb2))

	// Set active to "Option Y"
	_, err := handleRadioGroupSetActiveByLabel(rg, []any{"Option Y"})
	if err != nil {
		t.Fatalf("handleRadioGroupSetActiveByLabel failed: %v", err)
	}

	// Verify the active element
	active := rg.Active()
	if active != cb2 {
		t.Error("active element should be cb2")
	}
}

// TestHandleRadioGroupSetActiveByLabelNotFound tests setting active with non-existent label.
func TestHandleRadioGroupSetActiveByLabelNotFound(t *testing.T) {
	btn := widget.NewButton(widget.ButtonOpts.ToggleMode())
	setButtonText(btn, "Existing")

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn))

	// Try to set to non-existent label
	_, err := handleRadioGroupSetActiveByLabel(rg, []any{"NonExistent"})
	if err == nil {
		t.Fatal("expected error for non-existent label")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected error to contain 'not found', got: %s", err.Error())
	}
}

// TestHandleRadioGroupSetActiveByLabelWrongArgs tests SetActiveByLabel with wrong arguments.
func TestHandleRadioGroupSetActiveByLabelWrongArgs(t *testing.T) {
	btn := widget.NewButton(widget.ButtonOpts.ToggleMode())
	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn))

	// Test with no arguments
	_, err := handleRadioGroupSetActiveByLabel(rg, nil)
	if err == nil {
		t.Fatal("expected error for no arguments")
	}

	// Test with wrong type
	_, err = handleRadioGroupSetActiveByLabel(rg, []any{123})
	if err == nil {
		t.Fatal("expected error for wrong argument type")
	}
}

// TestRadioGroupElementInfoStruct tests that RadioGroupElementInfo struct has correct JSON tags.
func TestRadioGroupElementInfoStruct(t *testing.T) {
	info := RadioGroupElementInfo{
		Type:   "Button",
		Active: true,
		Label:  "Test Label",
	}

	// Verify struct fields
	if info.Type != "Button" {
		t.Errorf("expected Type 'Button', got '%s'", info.Type)
	}
	if !info.Active {
		t.Error("expected Active true")
	}
	if info.Label != "Test Label" {
		t.Errorf("expected Label 'Test Label', got '%s'", info.Label)
	}
}

// TestRadioGroupProxyHandlersRegistered tests that all RadioGroup handlers are registered.
func TestRadioGroupProxyHandlersRegistered(t *testing.T) {
	expectedHandlers := []string{
		"Elements",
		"ActiveIndex",
		"SetActiveByIndex",
		"ActiveLabel",
		"SetActiveByLabel",
	}

	for _, name := range expectedHandlers {
		handler, ok := radioGroupProxyHandlers[name]
		if !ok {
			t.Errorf("RadioGroup handler '%s' not registered", name)
		}
		if handler == nil {
			t.Errorf("RadioGroup handler '%s' is nil", name)
		}
	}
}

// setButtonText is a test helper to set the text field on a Button widget using reflection.
// This is needed because ButtonOpts.Text requires a font face which isn't available in tests.
func setButtonText(btn *widget.Button, label string) {
	txt := widget.NewText(widget.TextOpts.Text(label, nil, nil))
	v := reflect.ValueOf(btn).Elem()
	field := v.FieldByName("text")
	if !field.IsValid() {
		panic("setButtonText: 'text' field not found on Button")
	}
	fieldPtr := unsafe.Pointer(field.UnsafeAddr())
	realField := reflect.NewAt(field.Type(), fieldPtr).Elem()
	realField.Set(reflect.ValueOf(txt))
}

// setCheckboxText is a test helper to set the text field on a Checkbox widget using reflection.
// This is needed because CheckboxOpts.Text requires a font face which isn't available in tests.
// Checkbox uses a Label widget internally, which contains a Text widget.
func setCheckboxText(cb *widget.Checkbox, label string) {
	// Create a Text widget
	txt := widget.NewText(widget.TextOpts.Text(label, nil, nil))

	// Create a Label widget and set its text fields
	lbl := widget.NewLabel()
	setPrivateFieldOnCheckboxLabel(lbl, "Label", label)
	setPrivateFieldOnCheckboxLabel(lbl, "text", txt)

	// Set the Label on the Checkbox
	setPrivateFieldOnCheckbox(cb, "label", lbl)
	setPrivateFieldOnCheckbox(cb, "labelString", label)
}

// setPrivateFieldOnCheckbox is a test helper to set private fields on Checkbox.
func setPrivateFieldOnCheckbox(cb *widget.Checkbox, fieldName string, value interface{}) {
	v := reflect.ValueOf(cb).Elem()
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		panic("setPrivateFieldOnCheckbox: field '" + fieldName + "' not found on Checkbox")
	}
	fieldPtr := unsafe.Pointer(field.UnsafeAddr())
	realField := reflect.NewAt(field.Type(), fieldPtr).Elem()
	realField.Set(reflect.ValueOf(value))
}

// setPrivateFieldOnCheckboxLabel is a test helper to set private fields on Label (used by Checkbox).
func setPrivateFieldOnCheckboxLabel(lbl *widget.Label, fieldName string, value interface{}) {
	v := reflect.ValueOf(lbl).Elem()
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		panic("setPrivateFieldOnCheckboxLabel: field '" + fieldName + "' not found on Label")
	}
	fieldPtr := unsafe.Pointer(field.UnsafeAddr())
	realField := reflect.NewAt(field.Type(), fieldPtr).Elem()
	realField.Set(reflect.ValueOf(value))
}
