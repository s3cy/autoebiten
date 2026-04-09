---
name: widget-state-extraction
description: Expand autoui widget state extraction to support all ebitenui widget types
type: project
---

# Widget State Extraction Expansion

## Overview

Expand autoui's `ExtractWidgetState` function to support all ebitenui widgets that expose meaningful state for E2E testing.

**Why:** Full widget coverage enables comprehensive game UI automation. Currently only 6 of ~21 widgets are supported.

**Scope:** Add state extraction for 9 additional widgets. Skip 6 widgets that have no useful public state.

## Widgets to Add

### Interactive Input Widgets

#### List
- **API:** `Entries() []any`, `SelectedEntry() any`, `IsFocused() bool`
- **Attributes:** `entries` (count), `selected` (string representation), `focused`

#### RadioGroup
- **API:** `Active() RadioGroupElement`
- **Attributes:** `active` (element type - typically Button or Checkbox)

#### TextArea
- **API:** `GetText() string`
- **Attributes:** `text` (content)

#### ComboButton
- **API:** `ContentVisible bool` (field), `Label() string`
- **Attributes:** `label`, `open` (dropdown state)

#### ListComboButton
- **API:** `SelectedEntry() any`, `Label() string`, `ContentVisible() bool`, `IsFocused() bool`
- **Attributes:** `label`, `selected`, `open`, `focused`

### Navigation/Container Widgets

#### TabBook
- **API:** `Tab() *TabBookTab`
- **Attributes:** `active_tab` (tab label from Tab.Label)

#### ScrollContainer
- **API:** `ScrollLeft float64`, `ScrollTop float64` (fields), `ContentRect() image.Rectangle`
- **Attributes:** `scroll_x`, `scroll_y`, `content_width`, `content_height`

### Display/Misc Widgets

#### Text
- **API:** `Label string` (field), `MaxWidth float64` (field)
- **Attributes:** `text`, `max_width`

#### Window
- **API:** `Modal bool`, `Draggable bool`, `Resizeable bool`, `FocusedWindow bool` (fields)
- **Attributes:** `modal`, `draggable`, `resizeable`, `focused_window`

## Widgets to Skip

These widgets have no useful public state for E2E testing:

| Widget | Reason |
|--------|--------|
| Container | `children` field is private |
| Graphic | Only holds image data, no testable state |
| ToolTip | Transient popup, internal state machine |
| Caret | Internal widget, private visibility state |
| DragAndDrop | Complex transient interaction state |
| FlipBook | No public state exposed |

## Implementation

### File: `autoui/internal/widgetstate.go`

Add 9 case statements to `ExtractWidgetState` switch:

```go
case *widget.List:
    extractListState(v, result)
case *widget.RadioGroup:
    extractRadioGroupState(v, result)
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
case *widget.Window:
    extractWindowState(v, result)
```

Add 9 extractor functions following existing patterns:

```go
func extractListState(list *widget.List, result map[string]string) {
    result["entries"] = fmt.Sprintf("%d", len(list.Entries()))
    if selected := list.SelectedEntry(); selected != nil {
        result["selected"] = fmt.Sprintf("%v", selected)
    }
    if list.IsFocused() {
        result["focused"] = "true"
    }
}
```

### No Changes Required

- `autoui/handlers.go` - commands already handle any widget type
- `autoui/caller.go` - reflection-based invocation works with any method
- `autoui/register.go` - no new commands needed

## Testing

### Unit Tests

Add tests to `autoui/internal/widgetstate_test.go` for each extractor:

```go
func TestExtractWidgetState_List(t *testing.T) {
    list := widget.NewList(
        widget.ListOpts.Entries([]any{"item1", "item2"}),
        widget.ListOpts.EntryLabelFunc(func(e any) string { return e.(string) }),
    )
    list.Validate()
    
    result := ExtractWidgetState(list)
    
    assert.Equal(t, "2", result["entries"])
}
```

### Test Coverage

- Each extractor function tested independently
- Edge cases: empty collections, nil selections, unfocused widgets
- String formatting verified

## Success Criteria

- All 9 new widgets have state extraction
- Unit tests pass with >80% coverage on new code
- Existing tests continue to pass
- `autoui.find` and `autoui.at` commands work with new widget types