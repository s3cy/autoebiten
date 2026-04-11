package autoui_test

import (
	"testing"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/s3cy/autoebiten/autoui"
)

func TestRegisterRadioGroup(t *testing.T) {
	btn1 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn2 := widget.NewButton(widget.ButtonOpts.ToggleMode())

	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2))

	autoui.RegisterRadioGroup("test-group", rg)

	// Verify retrieval
	retrieved := autoui.GetRadioGroup("test-group")
	if retrieved != rg {
		t.Error("retrieved RadioGroup does not match registered")
	}

	// Clean up
	autoui.UnregisterRadioGroup("test-group")
	retrieved = autoui.GetRadioGroup("test-group")
	if retrieved != nil {
		t.Error("RadioGroup still exists after unregister")
	}
}

func TestGetRadioGroupNotFound(t *testing.T) {
	retrieved := autoui.GetRadioGroup("nonexistent")
	if retrieved != nil {
		t.Error("expected nil for nonexistent RadioGroup")
	}
}

func TestRegisterRadioGroupReplace(t *testing.T) {
	btn1 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	rg1 := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1))
	rg2 := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1))

	autoui.RegisterRadioGroup("replace-group", rg1)
	autoui.RegisterRadioGroup("replace-group", rg2)

	retrieved := autoui.GetRadioGroup("replace-group")
	if retrieved != rg2 {
		t.Error("RadioGroup was not replaced")
	}

	autoui.UnregisterRadioGroup("replace-group")
}
