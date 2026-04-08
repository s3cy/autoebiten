package autoui

import (
	"image"
	"reflect"

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