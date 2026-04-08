package autoui

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/antchfx/xmlquery"
)

// QueryXPath executes an XPath 1.0 expression against the widget tree.
// It converts widgets to XML, executes the XPath query, and maps results back to WidgetInfo.
// Returns an error if the XPath expression is invalid or XML parsing fails.
func QueryXPath(widgets []WidgetInfo, xpathExpr string) ([]WidgetInfo, error) {
	if len(widgets) == 0 {
		return nil, nil
	}

	if xpathExpr == "" {
		return nil, nil
	}

	// Generate XML from widgets
	xmlData, err := MarshalWidgetTreeXML(widgets)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal widget tree to XML: %w", err)
	}

	// Parse XML
	doc, err := xmlquery.Parse(bytes.NewReader(xmlData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	// Execute XPath query
	nodes, err := xmlquery.QueryAll(doc, xpathExpr)
	if err != nil {
		return nil, fmt.Errorf("XPath query failed: %w", err)
	}

	// Map XML nodes back to WidgetInfo
	return xmlNodesToWidgetInfo(nodes, widgets), nil
}

// xmlNodesToWidgetInfo maps XML nodes back to their corresponding WidgetInfo.
// Matching is done by widget type and position (x, y coordinates).
// This handles the case where multiple widgets might share the same type.
func xmlNodesToWidgetInfo(nodes []*xmlquery.Node, widgets []WidgetInfo) []WidgetInfo {
	if len(nodes) == 0 || len(widgets) == 0 {
		return nil
	}

	result := make([]WidgetInfo, 0, len(nodes))
	for _, node := range nodes {
		// Skip the UI root element
		if node.Data == "UI" {
			// UI is a wrapper, not a real widget
			continue
		}

		// Get type from element name
		widgetType := node.Data

		// Get position attributes
		xAttr := node.SelectAttr("x")
		yAttr := node.SelectAttr("y")

		if xAttr == "" || yAttr == "" {
			// Cannot match without position
			continue
		}

		// Parse position values
		xVal, err := strconv.Atoi(xAttr)
		if err != nil {
			continue
		}
		yVal, err := strconv.Atoi(yAttr)
		if err != nil {
			continue
		}

		// Find matching widget
		for _, w := range widgets {
			if w.Type == widgetType &&
				w.Rect.Min.X == xVal &&
				w.Rect.Min.Y == yVal {
				result = append(result, w)
				break // Only add first match
			}
		}
	}

	return result
}