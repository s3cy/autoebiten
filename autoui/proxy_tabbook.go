package autoui

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/s3cy/autoebiten/autoui/internal"
)

// TabInfo represents information about a tab in a TabBook.
type TabInfo struct {
	Index    int    `json:"index"`
	Label    string `json:"label"`
	Disabled bool   `json:"disabled"`
}

// handleTabs returns a list of TabInfo for all tabs in a TabBook.
func handleTabs(w widget.PreferredSizeLocateableWidget, args []any) (any, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("Tabs requires no arguments")
	}

	tb, ok := w.(*widget.TabBook)
	if !ok {
		return nil, fmt.Errorf("Tabs requires widget of type *TabBook, got %T", w)
	}

	tabs := internal.GetTabBookTabs(tb)
	result := make([]TabInfo, len(tabs))

	for i, tab := range tabs {
		result[i] = TabInfo{
			Index:    i,
			Label:    internal.GetTabBookTabLabel(tab),
			Disabled: tab.Disabled,
		}
	}

	return result, nil
}

// handleTabIndex returns the index of the active tab in a TabBook.
// Returns -1 if no tab is active.
func handleTabIndex(w widget.PreferredSizeLocateableWidget, args []any) (any, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("TabIndex requires no arguments")
	}

	tb, ok := w.(*widget.TabBook)
	if !ok {
		return nil, fmt.Errorf("TabIndex requires widget of type *TabBook, got %T", w)
	}

	activeTab := tb.Tab()
	if activeTab == nil {
		return int64(-1), nil
	}

	tabs := internal.GetTabBookTabs(tb)
	for i, tab := range tabs {
		if tab == activeTab {
			return int64(i), nil
		}
	}

	return int64(-1), nil
}

// handleSetTabByIndex sets the active tab by its index.
func handleSetTabByIndex(w widget.PreferredSizeLocateableWidget, args []any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("SetTabByIndex requires 1 argument (index)")
	}

	// JSON unmarshals numbers to float64, so we expect float64 input
	index, ok := args[0].(float64)
	if !ok {
		return nil, fmt.Errorf("index must be a number, got %T", args[0])
	}
	idx := int(index)

	tb, ok := w.(*widget.TabBook)
	if !ok {
		return nil, fmt.Errorf("SetTabByIndex requires widget of type *TabBook, got %T", w)
	}

	tabs := internal.GetTabBookTabs(tb)
	if len(tabs) == 0 {
		return nil, fmt.Errorf("TabBook has no tabs")
	}
	if idx < 0 || idx >= len(tabs) {
		return nil, fmt.Errorf("index %d out of range (0-%d)", idx, len(tabs)-1)
	}

	tab := tabs[idx]
	if tab.Disabled {
		return nil, fmt.Errorf("tab at index %d is disabled", idx)
	}

	tb.SetTab(tab)
	return nil, nil
}

// handleTabLabel returns the label of the active tab in a TabBook.
// Returns empty string if no tab is active.
func handleTabLabel(w widget.PreferredSizeLocateableWidget, args []any) (any, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("TabLabel requires no arguments")
	}

	tb, ok := w.(*widget.TabBook)
	if !ok {
		return nil, fmt.Errorf("TabLabel requires widget of type *TabBook, got %T", w)
	}

	activeTab := tb.Tab()
	if activeTab == nil {
		return "", nil
	}

	return internal.GetTabBookTabLabel(activeTab), nil
}

// handleSetTabByLabel sets the active tab by its label.
func handleSetTabByLabel(w widget.PreferredSizeLocateableWidget, args []any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("SetTabByLabel requires 1 argument (label)")
	}

	label, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("label must be a string, got %T", args[0])
	}

	tb, ok := w.(*widget.TabBook)
	if !ok {
		return nil, fmt.Errorf("SetTabByLabel requires widget of type *TabBook, got %T", w)
	}

	tabs := internal.GetTabBookTabs(tb)
	for _, tab := range tabs {
		if internal.GetTabBookTabLabel(tab) == label {
			if tab.Disabled {
				return nil, fmt.Errorf("tab with label '%s' is disabled", label)
			}
			tb.SetTab(tab)
			return nil, nil
		}
	}

	return nil, fmt.Errorf("tab with label '%s' not found", label)
}