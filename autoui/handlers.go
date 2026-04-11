package autoui

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ebitenui/ebitenui"
	"github.com/s3cy/autoebiten"
)

// CoordinateRequest represents a coordinate request in JSON format.
type CoordinateRequest struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// radioGroupPrefix is the target prefix for RadioGroup method calls.
const radioGroupPrefix = "radiogroup="

// handleTreeCommand handles the "tree" command which returns the full widget tree as XML.
func handleTreeCommand(ui *ebitenui.UI) func(ctx autoebiten.CommandContext) {
	return func(ctx autoebiten.CommandContext) {
		if ui == nil {
			ctx.Respond("error: UI not registered")
			return
		}

		// Walk the widget tree
		widgets := SnapshotTree(ui)

		// Marshal to XML
		xmlData, err := MarshalWidgetTreeXML(widgets)
		if err != nil {
			ctx.Respond("error: failed to marshal widget tree: " + err.Error())
			return
		}

		ctx.Respond(string(xmlData))
	}
}

// handleAtCommand handles the "at" command which returns the widget at given coordinates.
// Supports two formats: "x,y" or JSON `{"x":100,"y":200}`.
func handleAtCommand(ui *ebitenui.UI) func(ctx autoebiten.CommandContext) {
	return func(ctx autoebiten.CommandContext) {
		if ui == nil {
			ctx.Respond("error: UI not registered")
			return
		}

		request := ctx.Request()
		if request == "" {
			ctx.Respond("error: missing coordinates")
			return
		}

		// Parse coordinates
		x, y, err := parseCoordinates(request)
		if err != nil {
			ctx.Respond("error: " + err.Error())
			return
		}

		// Walk the widget tree
		widgets := SnapshotTree(ui)

		// Find widget at coordinates
		widget := FindAt(widgets, x, y)
		if widget == nil {
			ctx.Respond("error: no widget found at coordinates")
			return
		}

		// Marshal to XML
		xmlData, err := MarshalWidgetXML(*widget)
		if err != nil {
			ctx.Respond("error: failed to marshal widget: " + err.Error())
			return
		}

		ctx.Respond(string(xmlData))
	}
}

// handleFindCommand handles the "find" command which returns widgets matching a query.
// Supports two formats: "key=value" or JSON `{"key":"value",...}`.
func handleFindCommand(ui *ebitenui.UI) func(ctx autoebiten.CommandContext) {
	return func(ctx autoebiten.CommandContext) {
		if ui == nil {
			ctx.Respond("error: UI not registered")
			return
		}

		request := ctx.Request()

		// Walk the widget tree
		widgets := SnapshotTree(ui)

		// Determine query format and find widgets
		var matching []WidgetInfo
		if strings.HasPrefix(request, "{") {
			// JSON format
			matching = FindByQueryJSON(widgets, request)
		} else {
			// Simple key=value format
			matching = FindByQuery(widgets, request)
		}

		if len(matching) == 0 {
			ctx.Respond("error: no widgets found matching query")
			return
		}

		// Marshal to XML
		xmlData, err := MarshalWidgetsXML(matching)
		if err != nil {
			ctx.Respond("error: failed to marshal widgets: " + err.Error())
			return
		}

		ctx.Respond(string(xmlData))
	}
}

// handleXPathCommand handles the "xpath" command which returns widgets matching an XPath expression.
func handleXPathCommand(ui *ebitenui.UI) func(ctx autoebiten.CommandContext) {
	return func(ctx autoebiten.CommandContext) {
		if ui == nil {
			ctx.Respond("error: UI not registered")
			return
		}

		xpathExpr := ctx.Request()
		if xpathExpr == "" {
			ctx.Respond("error: missing XPath expression")
			return
		}

		// Walk the widget tree
		widgets := SnapshotTree(ui)

		// Execute XPath query
		matching, err := QueryXPath(widgets, xpathExpr)
		if err != nil {
			ctx.Respond("error: " + err.Error())
			return
		}

		if len(matching) == 0 {
			ctx.Respond("error: no widgets found matching XPath")
			return
		}

		// Marshal to XML
		xmlData, err := MarshalWidgetsXML(matching)
		if err != nil {
			ctx.Respond("error: failed to marshal widgets: " + err.Error())
			return
		}

		ctx.Respond(string(xmlData))
	}
}

// CallRequest represents a method invocation request.
type CallRequest struct {
	// Target is the query string to find the target widget.
	// Supports "key=value" or JSON format.
	Target string `json:"target"`

	// Method is the name of the method to invoke.
	Method string `json:"method"`

	// Args are the arguments to pass to the method.
	Args []any `json:"args"`
}

// CallResponse represents the response from a method invocation.
type CallResponse struct {
	// Success indicates if the invocation was successful.
	Success bool `json:"success"`

	// Error contains the error message if invocation failed.
	Error string `json:"error,omitempty"`

	// Result contains the return value from the method (if applicable).
	// Used for getter methods like ActiveIndex, TabIndex, etc.
	Result any `json:"result,omitempty"`
}

// ExistsResponse represents the response from the exists command.
// Returns JSON format for use with wait-for command.
type ExistsResponse struct {
	// Found indicates if any widgets matched the query.
	Found bool `json:"found"`

	// Count is the number of matching widgets.
	Count int `json:"count"`
}

// handleCallCommand handles the "call" command which invokes a method on a widget.
// Request format: `{"target":"query","method":"name","args":[...]}`.
// Special target format: "radiogroup=name" for RadioGroup method calls.
func handleCallCommand(ui *ebitenui.UI) func(ctx autoebiten.CommandContext) {
	return func(ctx autoebiten.CommandContext) {
		if ui == nil {
			ctx.Respond("error: UI not registered")
			return
		}

		request := ctx.Request()
		if request == "" {
			ctx.Respond("error: missing call request")
			return
		}

		// Parse the call request
		var callReq CallRequest
		if err := json.Unmarshal([]byte(request), &callReq); err != nil {
			ctx.Respond("error: invalid JSON format: " + err.Error())
			return
		}

		if callReq.Target == "" {
			ctx.Respond("error: missing target query")
			return
		}

		if callReq.Method == "" {
			ctx.Respond("error: missing method name")
			return
		}

		// Check for RadioGroup target prefix
		if strings.HasPrefix(callReq.Target, radioGroupPrefix) {
			name := strings.TrimPrefix(callReq.Target, radioGroupPrefix)
			rg := GetRadioGroup(name)
			if rg == nil {
				ctx.Respond(fmt.Sprintf("error: RadioGroup '%s' not registered. Did you call autoui.RegisterRadioGroup?", name))
				return
			}

			result, err := InvokeRadioGroupMethod(rg, callReq.Method, callReq.Args)

			response := CallResponse{Success: err == nil}
			if err != nil {
				response.Error = err.Error()
			}
			if result != nil {
				response.Result = result
			}

			respData, _ := json.Marshal(response)
			ctx.Respond(string(respData))
			return
		}

		// Walk the widget tree
		widgets := SnapshotTree(ui)

		// Find the target widget
		var targetWidget *WidgetInfo
		var matching []WidgetInfo

		// Use appropriate finder based on target format
		if len(callReq.Target) > 0 && callReq.Target[0] == '{' {
			matching = FindByQueryJSON(widgets, callReq.Target)
		} else {
			matching = FindByQuery(widgets, callReq.Target)
		}

		if len(matching) == 0 {
			ctx.Respond("error: no widget found matching target query")
			return
		}

		// Use the first matching widget
		targetWidget = &matching[0]

		// Check for proxy handler first (for methods that return values)
		var result any
		var err error

		if handler := GetProxyHandler(callReq.Method); handler != nil {
			// Use proxy handler for special methods (Tabs, TabIndex, etc.)
			result, err = handler(targetWidget.Widget, callReq.Args)
		} else {
			// Use regular InvokeMethod for standard methods
			err = InvokeMethod(targetWidget.Widget, callReq.Method, callReq.Args)
		}

		// Build response
		response := CallResponse{
			Success: err == nil,
		}
		if err != nil {
			response.Error = err.Error()
		}
		if result != nil {
			response.Result = result
		}

		// Marshal response to JSON
		respData, err := json.Marshal(response)
		if err != nil {
			ctx.Respond("error: failed to marshal response: " + err.Error())
			return
		}

		ctx.Respond(string(respData))
	}
}

// HighlightRequest represents a highlight request in JSON format.
type HighlightRequest struct {
	Query string `json:"query"` // Query to find widgets to highlight
	X     int    `json:"x"`     // Optional: X coordinate for direct highlight
	Y     int    `json:"y"`     // Optional: Y coordinate for direct highlight
}

// handleHighlightCommand handles the "highlight" command which adds visual highlights.
// Supports three modes:
// 1. "clear" - clears all highlights
// 2. "x,y" - highlights widget at coordinates
// 3. JSON `{"query":"..."}` - highlights widgets matching query
func handleHighlightCommand(ui *ebitenui.UI) func(ctx autoebiten.CommandContext) {
	return func(ctx autoebiten.CommandContext) {
		if ui == nil {
			ctx.Respond("error: UI not registered")
			return
		}

		request := ctx.Request()

		// Handle clear mode
		if request == "clear" {
			ClearHighlights()
			ctx.Respond("ok: highlights cleared")
			return
		}

		// Walk the widget tree for coordinate/query modes
		widgets := SnapshotTree(ui)

		// Handle coordinate mode (x,y format)
		if len(request) > 0 && request[0] != '{' {
			x, y, err := parseCoordinates(request)
			if err == nil {
				// Coordinate mode
				widget := FindAt(widgets, x, y)
				if widget == nil {
					ctx.Respond("error: no widget found at coordinates")
					return
				}

				AddHighlight(widget.Rect)
				ctx.Respond("ok: highlighted widget at coordinates")
				return
			}
			// If coordinate parsing failed, try as simple query
			matching := FindByQuery(widgets, request)
			if len(matching) == 0 {
				ctx.Respond("error: no widgets found")
				return
			}

			for _, w := range matching {
				AddHighlight(w.Rect)
			}
			ctx.Respond(fmt.Sprintf("ok: highlighted %d widgets", len(matching)))
			return
		}

		// Handle JSON format
		var highlightReq HighlightRequest
		if err := json.Unmarshal([]byte(request), &highlightReq); err != nil {
			ctx.Respond("error: invalid format, expected 'clear', 'x,y', query, or JSON")
			return
		}

		// If coordinates provided, use them
		if highlightReq.X != 0 || highlightReq.Y != 0 {
			widget := FindAt(widgets, highlightReq.X, highlightReq.Y)
			if widget == nil {
				ctx.Respond("error: no widget found at coordinates")
				return
			}

			AddHighlight(widget.Rect)
			ctx.Respond("ok: highlighted widget at coordinates")
			return
		}

		// If query provided, use it
		if highlightReq.Query != "" {
			matching := FindByQuery(widgets, highlightReq.Query)
			if len(matching) == 0 {
				ctx.Respond("error: no widgets found matching query")
				return
			}

			for _, w := range matching {
				AddHighlight(w.Rect)
			}
			ctx.Respond(fmt.Sprintf("ok: highlighted %d widgets", len(matching)))
			return
		}

		ctx.Respond("error: missing coordinates or query in JSON request")
	}
}

// handleExistsCommand handles the "exists" command which returns JSON indicating
// if widgets matching a query exist. Returns JSON for use with wait-for.
// Unlike find, this never returns an error for empty results.
func handleExistsCommand(ui *ebitenui.UI) func(ctx autoebiten.CommandContext) {
	return func(ctx autoebiten.CommandContext) {
		if ui == nil {
			// Return JSON error, not plain text
			resp := ExistsResponse{Found: false, Count: 0}
			data, _ := json.Marshal(resp)
			ctx.Respond(string(data))
			return
		}

		request := ctx.Request()

		// Walk the widget tree
		widgets := WalkTree(ui.Container)

		// Determine query format and find widgets
		var matching []WidgetInfo
		if strings.HasPrefix(request, "{") {
			// JSON format
			matching = FindByQueryJSON(widgets, request)
		} else {
			// Simple key=value format (empty request matches all)
			matching = FindByQuery(widgets, request)
		}

		// Build response - always valid JSON, no error on empty
		resp := ExistsResponse{
			Found: len(matching) > 0,
			Count: len(matching),
		}

		data, err := json.Marshal(resp)
		if err != nil {
			ctx.Respond(`{"found":false,"count":0}`)
			return
		}

		ctx.Respond(string(data))
	}
}

// parseCoordinates parses coordinates from "x,y" format or JSON format.
func parseCoordinates(request string) (int, int, error) {
	// Try JSON format first
	if strings.HasPrefix(request, "{") {
		var coord CoordinateRequest
		if err := json.Unmarshal([]byte(request), &coord); err != nil {
			return 0, 0, fmt.Errorf("invalid JSON format: %w", err)
		}
		return coord.X, coord.Y, nil
	}

	// Parse "x,y" format
	parts := strings.Split(request, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid coordinate format, expected 'x,y' or JSON")
	}

	x, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid x coordinate: %w", err)
	}

	y, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid y coordinate: %w", err)
	}

	return x, y, nil
}
