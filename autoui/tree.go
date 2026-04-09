package autoui

import (
	"fmt"
	"image"
	"reflect"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/s3cy/autoebiten/autoui/internal"
)

// WidgetInfo holds extracted information about a widget.
// This is the internal representation used before XML conversion.
type WidgetInfo struct {
	// Widget is the underlying widget instance.
	Widget widget.PreferredSizeLocateableWidget

	// Type is the widget type name (e.g., "Button", "Container").
	Type string

	// Rect is the widget's screen rectangle.
	Rect image.Rectangle

	// Visible indicates if the widget is currently visible.
	Visible bool

	// Disabled indicates if the widget is disabled.
	Disabled bool

	// State contains widget-specific state attributes.
	// Keys depend on widget type (e.g., "text", "checked", "value").
	State map[string]string

	// CustomData contains extracted custom data attributes.
	// Flattened from the widget's CustomData field.
	CustomData map[string]string

	// Addr is the widget pointer address for unique identification.
	// Format: "0x14000abc0" (hex string).
	Addr string
}

// ExtractWidgetInfo extracts widget information from a widget instance.
// It captures the widget's type, geometry, visibility, state, and custom data.
func ExtractWidgetInfo(w widget.PreferredSizeLocateableWidget) WidgetInfo {
	if w == nil {
		return WidgetInfo{}
	}

	// Get underlying widget for common properties
	baseWidget := w.GetWidget()

	info := WidgetInfo{
		Widget:   w,
		Type:     extractWidgetType(w),
		Rect:     baseWidget.Rect,
		Visible:  baseWidget.IsVisible(),
		Disabled: baseWidget.Disabled,
		Addr:     fmt.Sprintf("0x%x", reflect.ValueOf(w).Pointer()),
	}

	// Extract widget-specific state
	info.State = internal.ExtractWidgetState(w)

	// Extract custom data
	info.CustomData = internal.ExtractCustomData(baseWidget.CustomData)

	return info
}

// extractWidgetType returns the widget type name using reflection.
// It strips the "*" prefix for pointer types and returns just the type name.
func extractWidgetType(w widget.PreferredSizeLocateableWidget) string {
	t := reflect.TypeOf(w)

	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}

// WalkTree traverses the widget hierarchy depth-first and returns a flat list of all widgets.
// It starts from the root widget and recursively visits all children in Container widgets.
func WalkTree(root widget.PreferredSizeLocateableWidget) []WidgetInfo {
	if root == nil {
		return nil
	}

	var result []WidgetInfo
	walkTreeRecursive(root, &result)
	return result
}

// walkTreeRecursive recursively traverses the widget tree and appends widget info to the result.
// It performs depth-first traversal, processing the current widget before its children.
func walkTreeRecursive(w widget.PreferredSizeLocateableWidget, result *[]WidgetInfo) {
	// Extract info for current widget
	info := ExtractWidgetInfo(w)
	*result = append(*result, info)

	// If this is a Container, traverse children
	if container, ok := w.(widget.Containerer); ok {
		children := container.Children()
		for _, child := range children {
			walkTreeRecursive(child, result)
		}
	}
}

// SnapshotTree returns a snapshot of the widget tree from the UI.
// If RWLock was provided via RegisterOptions, it acquires RLock before traversal.
func SnapshotTree(ui *ebitenui.UI) []WidgetInfo {
	if ui == nil || ui.Container == nil {
		return nil
	}

	// Acquire read lock if user provided one for thread safety
	if treeRWLock != nil {
		treeRWLock.RLock()
		defer treeRWLock.RUnlock()
	}

	return WalkTree(ui.Container)
}