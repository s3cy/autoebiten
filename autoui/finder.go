package autoui

import (
	"encoding/json"
	"image"
	"strings"
)

// FindByQuery finds widgets matching a simple key=value query.
// The query format is "key=value" for exact match.
// Empty query returns all widgets.
func FindByQuery(widgets []WidgetInfo, query string) []WidgetInfo {
	if query == "" {
		return widgets
	}

	// Parse key=value format
	parts := strings.SplitN(query, "=", 2)
	if len(parts) != 2 {
		// Invalid format, return empty
		return nil
	}

	key := parts[0]
	value := parts[1]

	criteria := map[string]string{key: value}
	return filterWidgets(widgets, criteria)
}

// FindByQueryJSON finds widgets matching JSON criteria with AND logic.
// The query format is `{"key":"value",...}` - all criteria must match.
// Invalid JSON returns empty results.
func FindByQueryJSON(widgets []WidgetInfo, query string) []WidgetInfo {
	if query == "" {
		return nil
	}

	var criteria map[string]string
	if err := json.Unmarshal([]byte(query), &criteria); err != nil {
		// Invalid JSON, return empty
		return nil
	}

	if len(criteria) == 0 {
		return nil
	}

	return filterWidgets(widgets, criteria)
}

// filterWidgets filters the widget list by matching all criteria (AND logic).
func filterWidgets(widgets []WidgetInfo, criteria map[string]string) []WidgetInfo {
	if len(widgets) == 0 || len(criteria) == 0 {
		return nil
	}

	result := make([]WidgetInfo, 0, len(widgets))
	for _, w := range widgets {
		if matchesCriteria(w, criteria) {
			result = append(result, w)
		}
	}

	return result
}

// matchesCriteria checks if a widget matches all criteria (AND logic).
func matchesCriteria(w WidgetInfo, criteria map[string]string) bool {
	for key, expectedValue := range criteria {
		actualValue := getWidgetAttributeValue(w, key)
		if actualValue != expectedValue {
			return false
		}
	}
	return true
}

// getWidgetAttributeValue retrieves an attribute value from WidgetInfo.
// Lookup order:
// 1. Special keys: type, x, y, width, height, visible, disabled
// 2. State map: text, checked, value, etc.
// 3. CustomData map: user-defined attributes
func getWidgetAttributeValue(w WidgetInfo, key string) string {
	// Check special keys first
	switch key {
	case "type":
		return w.Type
	case "x":
		return formatInt(w.Rect.Min.X)
	case "y":
		return formatInt(w.Rect.Min.Y)
	case "width":
		return formatInt(w.Rect.Dx())
	case "height":
		return formatInt(w.Rect.Dy())
	case "visible":
		return formatBool(w.Visible)
	case "disabled":
		return formatBool(w.Disabled)
	}

	// Check State map
	if w.State != nil {
		if val, ok := w.State[key]; ok {
			return val
		}
	}

	// Check CustomData map
	if w.CustomData != nil {
		if val, ok := w.CustomData[key]; ok {
			return val
		}
	}

	return ""
}

// FindAt finds the top-most widget at the given coordinates.
// It searches in reverse order (last widget first) to find widgets added later.
// Returns nil if no widget contains the point.
func FindAt(widgets []WidgetInfo, x, y int) *WidgetInfo {
	if len(widgets) == 0 {
		return nil
	}

	point := image.Point{X: x, Y: y}

	// Search in reverse order (top-most first)
	// WalkTree adds widgets in depth-first order, so later widgets are "on top"
	for i := len(widgets) - 1; i >= 0; i-- {
		w := widgets[i]
		if point.In(w.Rect) {
			return &w
		}
	}

	return nil
}