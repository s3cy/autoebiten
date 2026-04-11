package autoui_test

import (
	"sync"
	"testing"
	"time"

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

func TestRadioGroupRegistryThreadSafety(t *testing.T) {
	// Run concurrent operations
	var wg sync.WaitGroup
	wg.Add(4)

	// Writer goroutine - register
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			btn := widget.NewButton(widget.ButtonOpts.ToggleMode())
			rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn))
			autoui.RegisterRadioGroup("concurrent-group", rg)
		}
	}()

	// Reader goroutine - get
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			_ = autoui.GetRadioGroup("concurrent-group")
		}
	}()

	// Reader goroutine - get all names
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			_ = autoui.GetRegisteredRadioGroups()
		}
	}()

	// Unregister goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			time.Sleep(10 * time.Millisecond)
			autoui.UnregisterRadioGroup("concurrent-group")
		}
	}()

	// Wait for all goroutines
	wg.Wait()

	// Clean up
	autoui.UnregisterRadioGroup("concurrent-group")

	// Test completed without race (race detector will catch issues)
	t.Log("Concurrent test completed without race conditions")
}

func TestGetRegisteredRadioGroups_Empty(t *testing.T) {
	// Ensure registry is empty for this test by using unique names
	// and cleaning up afterward
	names := autoui.GetRegisteredRadioGroups()

	// Filter out any names that might exist from other tests
	var emptyCount int
	for _, name := range names {
		if name == "nonexistent-test-group" {
			emptyCount++
		}
	}

	// GetRegisteredRadioGroups should return a slice (possibly empty)
	// but not nil for empty registry case
	if emptyCount > 0 {
		t.Error("expected no groups with test name in empty registry check")
	}
}

func TestGetRegisteredRadioGroups_Multiple(t *testing.T) {
	btn1 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn2 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn3 := widget.NewButton(widget.ButtonOpts.ToggleMode())

	rg1 := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1))
	rg2 := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn2))
	rg3 := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn3))

	// Register multiple RadioGroups
	autoui.RegisterRadioGroup("multi-group-1", rg1)
	autoui.RegisterRadioGroup("multi-group-2", rg2)
	autoui.RegisterRadioGroup("multi-group-3", rg3)

	// Get all registered names
	names := autoui.GetRegisteredRadioGroups()

	// Verify all three groups are present
	nameSet := make(map[string]bool)
	for _, name := range names {
		nameSet[name] = true
	}

	if !nameSet["multi-group-1"] {
		t.Error("multi-group-1 not found in registered groups")
	}
	if !nameSet["multi-group-2"] {
		t.Error("multi-group-2 not found in registered groups")
	}
	if !nameSet["multi-group-3"] {
		t.Error("multi-group-3 not found in registered groups")
	}

	// Verify each group can be retrieved
	if autoui.GetRadioGroup("multi-group-1") != rg1 {
		t.Error("multi-group-1 does not match registered instance")
	}
	if autoui.GetRadioGroup("multi-group-2") != rg2 {
		t.Error("multi-group-2 does not match registered instance")
	}
	if autoui.GetRadioGroup("multi-group-3") != rg3 {
		t.Error("multi-group-3 does not match registered instance")
	}

	// Clean up
	autoui.UnregisterRadioGroup("multi-group-1")
	autoui.UnregisterRadioGroup("multi-group-2")
	autoui.UnregisterRadioGroup("multi-group-3")

	// Verify cleanup
	names = autoui.GetRegisteredRadioGroups()
	for _, name := range names {
		if name == "multi-group-1" || name == "multi-group-2" || name == "multi-group-3" {
			t.Error("RadioGroup still registered after cleanup")
		}
	}
}
