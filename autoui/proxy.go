package autoui

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
)

// ProxyHandler handles a named operation on a widget, bypassing the whitelist.
type ProxyHandler func(w widget.PreferredSizeLocateableWidget, args []any) (any, error)

// proxyHandlers maps proxy method names to their handlers.
var proxyHandlers = map[string]ProxyHandler{
	"SelectEntryByIndex": handleSelectEntryByIndex,
	"SelectedEntryIndex": handleSelectedEntryIndex,
}

// GetProxyHandler returns the handler for a proxy method name, or nil if not found.
func GetProxyHandler(name string) ProxyHandler {
	return proxyHandlers[name]
}

// handleSelectEntryByIndex selects a list entry by its position index.
func handleSelectEntryByIndex(w widget.PreferredSizeLocateableWidget, args []any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("SelectEntryByIndex requires 1 argument (index)")
	}

	index, ok := args[0].(float64)
	if !ok {
		return nil, fmt.Errorf("index must be a number, got %T", args[0])
	}
	idx := int(index)

	if list, ok := w.(*widget.List); ok {
		entries := list.Entries()
		if idx < 0 || idx >= len(entries) {
			return nil, fmt.Errorf("index %d out of range (0-%d)", idx, len(entries)-1)
		}
		list.SetSelectedEntry(entries[idx])
		return nil, nil
	}

	return nil, fmt.Errorf("SelectEntryByIndex requires widget of type *List, got %T", w)
}

// handleSelectedEntryIndex returns the index of the currently selected entry.
func handleSelectedEntryIndex(w widget.PreferredSizeLocateableWidget, args []any) (any, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("SelectedEntryIndex requires no arguments")
	}

	if list, ok := w.(*widget.List); ok {
		selected := list.SelectedEntry()
		entries := list.Entries()
		for i, e := range entries {
			if e == selected {
				return int64(i), nil
			}
		}
		return int64(-1), nil
	}

	return nil, fmt.Errorf("SelectedEntryIndex requires widget of type *List, got %T", w)
}