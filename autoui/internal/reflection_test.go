package internal_test

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/s3cy/autoebiten/autoui/internal"
)

func TestGetRadioGroupElements(t *testing.T) {
	// Create buttons for radio group
	btn1 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn2 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn3 := widget.NewButton(widget.ButtonOpts.ToggleMode())

	// Create radio group
	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2, btn3))

	// Extract elements via reflection
	elements := internal.GetRadioGroupElements(rg)

	if len(elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(elements))
	}

	// Verify elements are the same buttons
	if elements[0] != btn1 {
		t.Error("element 0 is not btn1")
	}
	if elements[1] != btn2 {
		t.Error("element 1 is not btn2")
	}
	if elements[2] != btn3 {
		t.Error("element 2 is not btn3")
	}
}

func TestGetRadioGroupElementsEmpty(t *testing.T) {
	rg := widget.NewRadioGroup()
	elements := internal.GetRadioGroupElements(rg)

	if len(elements) != 0 {
		t.Errorf("expected 0 elements for empty RadioGroup, got %d", len(elements))
	}
}

func TestGetTabBookTabs(t *testing.T) {
	// Create TabBookTab instances (minimal setup)
	tab1 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("General"))
	tab2 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("Settings"))

	// Note: TabBook requires full validation which needs theme setup
	// For this test, we create a TabBook and manually set tabs field via reflection
	// This tests the getter function, not the full TabBook initialization

	tb := widget.NewTabBook()
	// Set tabs via reflection for testing
	setPrivateField(tb, "tabs", []*widget.TabBookTab{tab1, tab2})

	tabs := internal.GetTabBookTabs(tb)

	if len(tabs) != 2 {
		t.Fatalf("expected 2 tabs, got %d", len(tabs))
	}

	if tabs[0] != tab1 {
		t.Error("tab 0 is not tab1")
	}
	if tabs[1] != tab2 {
		t.Error("tab 1 is not tab2")
	}
}

func TestGetTabBookTabLabel(t *testing.T) {
	tab := widget.NewTabBookTab(widget.TabBookTabOpts.Label("Settings"))
	label := internal.GetTabBookTabLabel(tab)

	if label != "Settings" {
		t.Errorf("expected label 'Settings', got '%s'", label)
	}
}

// Helper for testing: setPrivateField
func setPrivateField(obj interface{}, fieldName string, value interface{}) {
	v := reflect.ValueOf(obj).Elem()
	field := v.FieldByName(fieldName)
	fieldPtr := unsafe.Pointer(field.UnsafeAddr())
	realField := reflect.NewAt(field.Type(), fieldPtr).Elem()
	realField.Set(reflect.ValueOf(value))
}
