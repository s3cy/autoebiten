package autoui

import (
	"fmt"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/s3cy/autoebiten/autoui/internal"
)

// RadioGroupElementInfo represents information about a RadioGroup element.
type RadioGroupElementInfo struct {
	Type   string `json:"type"`
	Active bool   `json:"active"`
	Label  string `json:"label,omitempty"`
}

// RadioGroupHandler handles a named operation on a RadioGroup.
type RadioGroupHandler func(rg *widget.RadioGroup, args []any) (any, error)

// radioGroupProxyHandlers maps proxy method names to their handlers.
// This is a read-only registry after initialization - safe for concurrent access.
var radioGroupProxyHandlers = map[string]RadioGroupHandler{
	"Elements":         handleRadioGroupElements,
	"ActiveIndex":      handleRadioGroupActiveIndex,
	"SetActiveByIndex": handleRadioGroupSetActiveByIndex,
	"ActiveLabel":      handleRadioGroupActiveLabel,
	"SetActiveByLabel": handleRadioGroupSetActiveByLabel,
}

// InvokeRadioGroupMethod invokes a named method on a RadioGroup.
// Returns the result of the method call or an error if the method is not found
// or the method execution fails.
func InvokeRadioGroupMethod(rg *widget.RadioGroup, method string, args []any) (any, error) {
	handler, ok := radioGroupProxyHandlers[method]
	if !ok {
		return nil, fmt.Errorf("unknown RadioGroup method: %s", method)
	}
	return handler(rg, args)
}

// handleRadioGroupElements returns a list of RadioGroupElementInfo for all elements.
func handleRadioGroupElements(rg *widget.RadioGroup, args []any) (any, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("Elements requires no arguments")
	}

	elements := internal.GetRadioGroupElements(rg)
	active := rg.Active()

	result := make([]RadioGroupElementInfo, len(elements))
	for i, elem := range elements {
		info := RadioGroupElementInfo{
			Active: elem == active,
		}

		// Determine type and extract label via type switch
		switch e := elem.(type) {
		case *widget.Button:
			info.Type = "Button"
			if text := e.Text(); text != nil {
				info.Label = text.Label
			}
		case *widget.Checkbox:
			info.Type = "Checkbox"
			if text := e.Text(); text != nil {
				info.Label = text.Label
			}
		default:
			info.Type = "Unknown"
		}

		result[i] = info
	}

	return result, nil
}

// handleRadioGroupActiveIndex returns the index of the active element, or -1 if none.
func handleRadioGroupActiveIndex(rg *widget.RadioGroup, args []any) (any, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("ActiveIndex requires no arguments")
	}

	active := rg.Active()
	if active == nil {
		return int64(-1), nil
	}

	elements := internal.GetRadioGroupElements(rg)
	for i, elem := range elements {
		if elem == active {
			return int64(i), nil
		}
	}

	return int64(-1), nil
}

// handleRadioGroupSetActiveByIndex sets the active element by its index.
func handleRadioGroupSetActiveByIndex(rg *widget.RadioGroup, args []any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("SetActiveByIndex requires 1 argument (index)")
	}

	// JSON unmarshals numbers to float64, so we expect float64 input
	index, ok := args[0].(float64)
	if !ok {
		return nil, fmt.Errorf("index must be a number, got %T", args[0])
	}
	idx := int(index)

	elements := internal.GetRadioGroupElements(rg)
	if len(elements) == 0 {
		return nil, fmt.Errorf("RadioGroup has no elements")
	}
	if idx < 0 || idx >= len(elements) {
		return nil, fmt.Errorf("index %d out of range (0-%d)", idx, len(elements)-1)
	}

	rg.SetActive(elements[idx])
	return nil, nil
}

// handleRadioGroupActiveLabel returns the label of the active element, or "" if none.
func handleRadioGroupActiveLabel(rg *widget.RadioGroup, args []any) (any, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("ActiveLabel requires no arguments")
	}

	active := rg.Active()
	if active == nil {
		return "", nil
	}

	// Extract label via type switch
	switch e := active.(type) {
	case *widget.Button:
		if text := e.Text(); text != nil {
			return text.Label, nil
		}
	case *widget.Checkbox:
		if text := e.Text(); text != nil {
			return text.Label, nil
		}
	}

	return "", nil
}

// handleRadioGroupSetActiveByLabel sets the active element by its label.
func handleRadioGroupSetActiveByLabel(rg *widget.RadioGroup, args []any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("SetActiveByLabel requires 1 argument (label)")
	}

	label, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("label must be a string, got %T", args[0])
	}

	elements := internal.GetRadioGroupElements(rg)
	for _, elem := range elements {
		var elemLabel string

		switch e := elem.(type) {
		case *widget.Button:
			if text := e.Text(); text != nil {
				elemLabel = text.Label
			}
		case *widget.Checkbox:
			if text := e.Text(); text != nil {
				elemLabel = text.Label
			}
		}

		if elemLabel == label {
			rg.SetActive(elem)
			return nil, nil
		}
	}

	return nil, fmt.Errorf("element with label '%s' not found", label)
}
