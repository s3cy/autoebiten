package internal

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
)

// ExtractWidgetState extracts widget-specific state as a string map.
// Each widget type has different state attributes that are extracted.
func ExtractWidgetState(w widget.PreferredSizeLocateableWidget) map[string]string {
	result := make(map[string]string)

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
	case *widget.List:
		extractListState(v, result)
	case *widget.TextArea:
		extractTextAreaState(v, result)
	case *widget.ComboButton:
		extractComboButtonState(v, result)
	case *widget.ListComboButton:
		extractListComboButtonState(v, result)
	case *widget.TabBook:
		extractTabBookState(v, result)
	case *widget.ScrollContainer:
		extractScrollContainerState(v, result)
	case *widget.Text:
		extractTextState(v, result)
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

// extractListState extracts state from a List widget.
// Attributes: entries, selected, focused
func extractListState(list *widget.List, result map[string]string) {
	result["entries"] = fmt.Sprintf("%d", len(list.Entries()))

	if selected := list.SelectedEntry(); selected != nil {
		result["selected"] = fmt.Sprintf("%v", selected)
	}

	if list.IsFocused() {
		result["focused"] = "true"
	}
}

// extractTextAreaState extracts state from a TextArea widget.
// Attributes: text
func extractTextAreaState(ta *widget.TextArea, result map[string]string) {
	result["text"] = ta.GetText()
}

// extractComboButtonState extracts state from a ComboButton widget.
// Attributes: label, open
func extractComboButtonState(cb *widget.ComboButton, result map[string]string) {
	// Set open state first (public field, won't panic)
	if cb.ContentVisible {
		result["open"] = "true"
	} else {
		result["open"] = "false"
	}

	// Label() may panic if widget not fully initialized
	defer func() {
		if r := recover(); r != nil {
			// Widget not fully initialized, skip label extraction
		}
	}()

	result["label"] = cb.Label()
}

// extractListComboButtonState extracts state from a ListComboButton widget.
// Attributes: label, selected, open, focused
func extractListComboButtonState(lcb *widget.ListComboButton, result map[string]string) {
	// Set open state first
	if lcb.ContentVisible() {
		result["open"] = "true"
	} else {
		result["open"] = "false"
	}

	// Label() may panic if widget not fully initialized
	defer func() {
		if r := recover(); r != nil {
			// Widget not fully initialized, skip label extraction
		}
	}()

	result["label"] = lcb.Label()

	if selected := lcb.SelectedEntry(); selected != nil {
		result["selected"] = fmt.Sprintf("%v", selected)
	}

	if lcb.IsFocused() {
		result["focused"] = "true"
	}
}

// extractTabBookState extracts state from a TabBook widget.
// Attributes: active_tab (label of active tab)
func extractTabBookState(tb *widget.TabBook, result map[string]string) {
	if tab := tb.Tab(); tab != nil {
		result["active_tab"] = GetTabBookTabLabel(tab)
	}
}

// extractScrollContainerState extracts state from a ScrollContainer widget.
// Attributes: scroll_x, scroll_y, content_width, content_height
func extractScrollContainerState(sc *widget.ScrollContainer, result map[string]string) {
	result["scroll_x"] = fmt.Sprintf("%.2f", sc.ScrollLeft)
	result["scroll_y"] = fmt.Sprintf("%.2f", sc.ScrollTop)

	contentRect := sc.ContentRect()
	result["content_width"] = fmt.Sprintf("%d", contentRect.Dx())
	result["content_height"] = fmt.Sprintf("%d", contentRect.Dy())
}

// extractTextState extracts state from a Text widget.
// Attributes: text, max_width
func extractTextState(txt *widget.Text, result map[string]string) {
	result["text"] = txt.Label

	if txt.MaxWidth > 0 {
		result["max_width"] = fmt.Sprintf("%.0f", txt.MaxWidth)
	}
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
