package autoui

import (
	"bytes"
	"fmt"

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
// Matching is done by _addr attribute (widget pointer address) for exact identification.
func xmlNodesToWidgetInfo(nodes []*xmlquery.Node, widgets []WidgetInfo) []WidgetInfo {
	if len(nodes) == 0 || len(widgets) == 0 {
		return nil
	}

	// Build addr -> WidgetInfo map for efficient lookup
	addrMap := make(map[string]WidgetInfo, len(widgets))
	for _, w := range widgets {
		addrMap[w.Addr] = w
	}

	result := make([]WidgetInfo, 0, len(nodes))
	for _, node := range nodes {
		// Skip the UI root element
		if node.Data == "UI" {
			continue
		}

		// Get _addr attribute
		addr := node.SelectAttr("_addr")
		if addr == "" {
			// Cannot match without _addr
			continue
		}

		// Find matching widget by address
		if w, ok := addrMap[addr]; ok {
			result = append(result, w)
		}
	}

	return result
}