package autoui_test

import (
	"image"
	"strings"
	"testing"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/s3cy/autoebiten/autoui"
)

func TestProxy_SelectEntryByIndex(t *testing.T) {
	entries := []any{"Entry A", "Entry B", "Entry C"}
	list := widget.NewList(
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e any) string {
			return e.(string)
		}),
	)
	list.GetWidget().Rect = image.Rect(10, 10, 210, 200)

	info := autoui.ExtractWidgetInfo(list)

	result, err := autoui.InvokeMethodWithResult(info.Widget, "SelectEntryByIndex", []any{float64(1)})
	if err != nil {
		t.Errorf("SelectEntryByIndex failed: %v", err)
	}

	selected := list.SelectedEntry()
	if selected != entries[1] {
		t.Errorf("Expected 'Entry B', got %v", selected)
	}

	// result should be nil (setter, no return value)
	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestProxy_SelectedEntryIndex(t *testing.T) {
	entries := []any{"Entry A", "Entry B", "Entry C"}
	list := widget.NewList(
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e any) string {
			return e.(string)
		}),
	)
	list.GetWidget().Rect = image.Rect(10, 10, 210, 200)
	list.SetSelectedEntry(entries[1])

	info := autoui.ExtractWidgetInfo(list)

	result, err := autoui.InvokeMethodWithResult(info.Widget, "SelectedEntryIndex", nil)
	if err != nil {
		t.Errorf("SelectedEntryIndex failed: %v", err)
	}

	index, ok := result.(int64)
	if !ok {
		t.Errorf("Expected int64, got %T", result)
	}
	if index != 1 {
		t.Errorf("Expected index 1, got %d", index)
	}
}

func TestProxy_SelectedEntryIndex_NoSelection(t *testing.T) {
	entries := []any{"Entry A", "Entry B", "Entry C"}
	list := widget.NewList(
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e any) string {
			return e.(string)
		}),
	)
	list.GetWidget().Rect = image.Rect(10, 10, 210, 200)
	// No selection set

	info := autoui.ExtractWidgetInfo(list)

	result, err := autoui.InvokeMethodWithResult(info.Widget, "SelectedEntryIndex", nil)
	if err != nil {
		t.Errorf("SelectedEntryIndex failed: %v", err)
	}

	index, ok := result.(int64)
	if !ok {
		t.Errorf("Expected int64, got %T", result)
	}
	if index != -1 {
		t.Errorf("Expected index -1 for no selection, got %d", index)
	}
}

func TestProxy_SelectEntryByIndex_OutOfBounds(t *testing.T) {
	entries := []any{"A", "B"}
	list := widget.NewList(
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e any) string { return e.(string) }),
	)
	list.GetWidget().Rect = image.Rect(10, 10, 210, 200)

	info := autoui.ExtractWidgetInfo(list)

	_, err := autoui.InvokeMethodWithResult(info.Widget, "SelectEntryByIndex", []any{float64(5)})
	if err == nil {
		t.Error("Expected error for out of bounds index")
	}
	if !strings.Contains(err.Error(), "out of range") {
		t.Errorf("Expected 'out of range' error, got: %v", err)
	}
}

func TestProxy_SelectEntryByIndex_NegativeIndex(t *testing.T) {
	entries := []any{"A", "B"}
	list := widget.NewList(
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e any) string { return e.(string) }),
	)
	list.GetWidget().Rect = image.Rect(10, 10, 210, 200)

	info := autoui.ExtractWidgetInfo(list)

	_, err := autoui.InvokeMethodWithResult(info.Widget, "SelectEntryByIndex", []any{float64(-1)})
	if err == nil {
		t.Error("Expected error for negative index")
	}
	if !strings.Contains(err.Error(), "out of range") {
		t.Errorf("Expected 'out of range' error, got: %v", err)
	}
}

func TestProxy_SelectEntryByIndex_WrongWidgetType(t *testing.T) {
	btn := widget.NewButton()
	btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)

	info := autoui.ExtractWidgetInfo(btn)

	_, err := autoui.InvokeMethodWithResult(info.Widget, "SelectEntryByIndex", []any{float64(0)})
	if err == nil {
		t.Error("Expected error for wrong widget type")
	}
	if !strings.Contains(err.Error(), "requires widget of type *List") {
		t.Errorf("Expected type error, got: %v", err)
	}
}

func TestProxy_SelectedEntryIndex_WrongWidgetType(t *testing.T) {
	btn := widget.NewButton()
	btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)

	info := autoui.ExtractWidgetInfo(btn)

	_, err := autoui.InvokeMethodWithResult(info.Widget, "SelectedEntryIndex", nil)
	if err == nil {
		t.Error("Expected error for wrong widget type")
	}
	if !strings.Contains(err.Error(), "requires widget of type *List") {
		t.Errorf("Expected type error, got: %v", err)
	}
}

func TestProxy_SelectEntryByIndex_WrongArgCount(t *testing.T) {
	entries := []any{"A", "B"}
	list := widget.NewList(
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e any) string { return e.(string) }),
	)
	list.GetWidget().Rect = image.Rect(10, 10, 210, 200)

	info := autoui.ExtractWidgetInfo(list)

	_, err := autoui.InvokeMethodWithResult(info.Widget, "SelectEntryByIndex", nil)
	if err == nil {
		t.Error("Expected error for missing argument")
	}
	if !strings.Contains(err.Error(), "requires 1 argument") {
		t.Errorf("Expected argument count error, got: %v", err)
	}
}

func TestProxy_SelectEntryByIndex_WrongArgType(t *testing.T) {
	entries := []any{"A", "B"}
	list := widget.NewList(
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e any) string { return e.(string) }),
	)
	list.GetWidget().Rect = image.Rect(10, 10, 210, 200)

	info := autoui.ExtractWidgetInfo(list)

	_, err := autoui.InvokeMethodWithResult(info.Widget, "SelectEntryByIndex", []any{"not a number"})
	if err == nil {
		t.Error("Expected error for wrong argument type")
	}
	if !strings.Contains(err.Error(), "must be a number") {
		t.Errorf("Expected type error, got: %v", err)
	}
}

func TestProxy_SelectedEntryIndex_WithArgs(t *testing.T) {
	entries := []any{"A", "B"}
	list := widget.NewList(
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e any) string { return e.(string) }),
	)
	list.GetWidget().Rect = image.Rect(10, 10, 210, 200)

	info := autoui.ExtractWidgetInfo(list)

	_, err := autoui.InvokeMethodWithResult(info.Widget, "SelectedEntryIndex", []any{"unexpected"})
	if err == nil {
		t.Error("Expected error for unexpected argument")
	}
	if !strings.Contains(err.Error(), "requires no arguments") {
		t.Errorf("Expected argument count error, got: %v", err)
	}
}

func TestProxy_SelectEntryByIndex_FirstEntry(t *testing.T) {
	entries := []any{"First", "Second", "Third"}
	list := widget.NewList(
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e any) string { return e.(string) }),
	)
	list.GetWidget().Rect = image.Rect(10, 10, 210, 200)

	info := autoui.ExtractWidgetInfo(list)

	_, err := autoui.InvokeMethodWithResult(info.Widget, "SelectEntryByIndex", []any{float64(0)})
	if err != nil {
		t.Errorf("SelectEntryByIndex failed: %v", err)
	}

	selected := list.SelectedEntry()
	if selected != entries[0] {
		t.Errorf("Expected 'First', got %v", selected)
	}
}

func TestProxy_SelectEntryByIndex_LastEntry(t *testing.T) {
	entries := []any{"First", "Second", "Third"}
	list := widget.NewList(
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e any) string { return e.(string) }),
	)
	list.GetWidget().Rect = image.Rect(10, 10, 210, 200)

	info := autoui.ExtractWidgetInfo(list)

	_, err := autoui.InvokeMethodWithResult(info.Widget, "SelectEntryByIndex", []any{float64(2)})
	if err != nil {
		t.Errorf("SelectEntryByIndex failed: %v", err)
	}

	selected := list.SelectedEntry()
	if selected != entries[2] {
		t.Errorf("Expected 'Third', got %v", selected)
	}
}

// TestProxy_SelectEntryByIndex_EmptyList tests selecting from an empty list.
func TestProxy_SelectEntryByIndex_EmptyList(t *testing.T) {
	list := widget.NewList(
		widget.ListOpts.Entries([]any{}), // Empty list
		widget.ListOpts.EntryLabelFunc(func(e any) string { return e.(string) }),
	)
	list.GetWidget().Rect = image.Rect(10, 10, 210, 200)

	info := autoui.ExtractWidgetInfo(list)

	_, err := autoui.InvokeMethodWithResult(info.Widget, "SelectEntryByIndex", []any{float64(0)})
	if err == nil {
		t.Error("Expected error for empty list")
	}
	if !strings.Contains(err.Error(), "has no entries") {
		t.Errorf("Expected 'has no entries' error, got: %v", err)
	}
}