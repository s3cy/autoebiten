package internal

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
)

// ExtractWidgetState extracts widget-specific state as a string map.
// Each widget type has different state attributes that are extracted.
func ExtractWidgetState(w widget.PreferredSizeLocateableWidget) map[string]string {
	result := make(map[string]string)

	// Get the underlying widget for common properties
	baseWidget := w.GetWidget()

	// Type switch for widget-specific extraction
	switch v := w.(type) {
	case *widget.Button:
		extractButtonState(v, result)
	case *widget.TextInput:
		extractTextInputState(v, result)
	case *widget.Checkbox:
		extractCheckboxState(v, result)
	case *widget.Slider:
		extractSliderState(v, result)
	case *widget.Label:
		extractLabelState(v, result)
	case *widget.ProgressBar:
		extractProgressBarState(v, result)
	}

	// Add base widget state if applicable
	if baseWidget.Disabled {
		result["disabled"] = "true"
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

// extractButtonState extracts state from a Button widget.
// Attributes: text, state, toggle
func extractButtonState(btn *widget.Button, result map[string]string) {
	// Extract text if available
	if text := btn.Text(); text != nil {
		result["text"] = text.Label
	}

	// Extract button state
	result["state"] = widgetStateToString(btn.State())

	// Check if toggle mode
	if btn.ToggleMode {
		result["toggle"] = "true"
	}

	// Check if focused
	if btn.IsFocused() {
		result["focused"] = "true"
	}
}

// extractTextInputState extracts state from a TextInput widget.
// Attributes: text, cursor, focused
func extractTextInputState(input *widget.TextInput, result map[string]string) {
	// Extract text content
	result["text"] = input.GetText()

	// Note: Cursor position is internal and not exposed by TextInput
	// We can only determine if it's focused
	if input.IsFocused() {
		result["focused"] = "true"
	}
}

// extractCheckboxState extracts state from a Checkbox widget.
// Attributes: checked, triState
func extractCheckboxState(cb *widget.Checkbox, result map[string]string) {
	// Extract checkbox state
	result["state"] = widgetStateToString(cb.State())

	// Extract label text if available
	if text := cb.Text(); text != nil {
		result["text"] = text.Label
	}

	// Determine checked status from state
	if cb.State() == widget.WidgetChecked {
		result["checked"] = "true"
	} else if cb.State() == widget.WidgetGreyed {
		result["checked"] = "greyed"
	} else {
		result["checked"] = "false"
	}

	// Check if focused
	if cb.IsFocused() {
		result["focused"] = "true"
	}
}

// extractSliderState extracts state from a Slider widget.
// Attributes: value, min, max
func extractSliderState(slider *widget.Slider, result map[string]string) {
	result["value"] = fmt.Sprintf("%d", slider.Current)
	result["min"] = fmt.Sprintf("%d", slider.Min)
	result["max"] = fmt.Sprintf("%d", slider.Max)

	// Check if focused
	if slider.IsFocused() {
		result["focused"] = "true"
	}
}

// extractLabelState extracts state from a Label widget.
// Attributes: text
func extractLabelState(label *widget.Label, result map[string]string) {
	result["text"] = label.Label
}

// extractProgressBarState extracts state from a ProgressBar widget.
// Attributes: value, min, max
func extractProgressBarState(pb *widget.ProgressBar, result map[string]string) {
	result["value"] = fmt.Sprintf("%d", pb.GetCurrent())
	result["min"] = fmt.Sprintf("%d", pb.Min)
	result["max"] = fmt.Sprintf("%d", pb.Max)
}

// widgetStateToString converts a WidgetState enum to string.
func widgetStateToString(state widget.WidgetState) string {
	switch state {
	case widget.WidgetUnchecked:
		return "unchecked"
	case widget.WidgetChecked:
		return "checked"
	case widget.WidgetGreyed:
		return "greyed"
	default:
		return fmt.Sprintf("unknown(%d)", state)
	}
}