package autoui

import (
	"encoding/xml"
	"fmt"
	"maps"
	"sort"

	"github.com/ebitenui/ebitenui/widget"
)

// WidgetNode represents a widget as an XML node with dynamic element names.
// It implements xml.Marshaler for custom XML element name generation.
type WidgetNode struct {
	XMLName  xml.Name
	Attrs    map[string]string
	Children []*WidgetNode
}

// stringKeyValue is a helper struct for sorted map iteration.
type stringKeyValue struct {
	Key   string
	Value string
}

// MarshalWidgetXML converts a single WidgetInfo to XML bytes.
// The widget type becomes the XML element name.
func MarshalWidgetXML(info WidgetInfo) ([]byte, error) {
	node := widgetInfoToNode(info)
	return xml.Marshal(node)
}

// MarshalWidgetTreeXML converts a flat list of widgets to a tree XML structure.
// It reconstructs the parent-child relationships and wraps the result in a <UI> root element.
func MarshalWidgetTreeXML(widgets []WidgetInfo) ([]byte, error) {
	if len(widgets) == 0 {
		return []byte("<UI/>"), nil
	}

	// Build tree from flat list
	root := widgetInfoToNode(widgets[0])
	buildNodeTree(root, widgets)

	// Wrap in UI root element
	uiNode := &WidgetNode{
		XMLName:  xml.Name{Local: "UI"},
		Children: []*WidgetNode{root},
	}

	return xml.MarshalIndent(uiNode, "", "  ")
}

// widgetInfoToNode converts a WidgetInfo to a WidgetNode.
// It extracts position, state, and custom data attributes.
func widgetInfoToNode(info WidgetInfo) *WidgetNode {
	node := &WidgetNode{
		XMLName: xml.Name{Local: info.Type},
		Attrs:   make(map[string]string),
	}

	// Add position attributes
	node.Attrs["x"] = formatInt(info.Rect.Min.X)
	node.Attrs["y"] = formatInt(info.Rect.Min.Y)
	node.Attrs["width"] = formatInt(info.Rect.Dx())
	node.Attrs["height"] = formatInt(info.Rect.Dy())

	// Add state attributes
	node.Attrs["visible"] = formatBool(info.Visible)
	node.Attrs["disabled"] = formatBool(info.Disabled)

	// Add widget-specific state
	maps.Copy(node.Attrs, info.State)

	// Add custom data
	maps.Copy(node.Attrs, info.CustomData)

	return node
}

// MarshalXML implements xml.Marshaler for custom element name generation.
// It writes the element with its type as the name and all attributes sorted alphabetically.
func (n *WidgetNode) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// Use the XMLName for the element name
	start.Name = n.XMLName

	// Build sorted attributes
	sorted := sortedMap(n.Attrs)
	for _, kv := range sorted {
		attr := xml.Attr{
			Name:  xml.Name{Local: kv.Key},
			Value: kv.Value,
		}
		start.Attr = append(start.Attr, attr)
	}

	// Write start element
	if err := e.EncodeToken(start); err != nil {
		return fmt.Errorf("failed to encode start element: %w", err)
	}

	// Write children if present
	for _, child := range n.Children {
		if err := e.Encode(child); err != nil {
			return fmt.Errorf("failed to encode child: %w", err)
		}
	}

	// Write end element
	end := xml.EndElement{Name: n.XMLName}
	if err := e.EncodeToken(end); err != nil {
		return fmt.Errorf("failed to encode end element: %w", err)
	}

	return nil
}

// buildNodeTree reconstructs the tree structure from a flat list of widgets.
// It uses Widget.Parent() to determine actual parent-child relationships.
// For filtered results where parent containers may not be in the list, orphan widgets are
// added to the root node as siblings.
func buildNodeTree(root *WidgetNode, widgets []WidgetInfo) {
	if len(widgets) <= 1 {
		return
	}

	// Map base widgets to their corresponding nodes for parent lookup
	// Key is the Widget pointer (from GetWidget()), value is the XML node
	widgetToNode := map[*widget.Widget]*WidgetNode{}
	widgetToNode[widgets[0].Widget.GetWidget()] = root

	// Process remaining widgets (skip widgets[0] which is already the root)
	for _, info := range widgets[1:] {
		node := widgetInfoToNode(info)
		baseWidget := info.Widget.GetWidget()
		widgetToNode[baseWidget] = node

		// Find parent using Widget.Parent()
		parentWidget := baseWidget.Parent()
		if parentNode, ok := widgetToNode[parentWidget]; ok {
			parentNode.Children = append(parentNode.Children, node)
		} else {
			// Fallback: no parent found (filtered result without parent container)
			// Add to root node as a sibling
			root.Children = append(root.Children, node)
		}
	}
}

// sortedMap returns the map entries sorted by key for deterministic output.
func sortedMap(m map[string]string) []stringKeyValue {
	if m == nil {
		return nil
	}

	result := make([]stringKeyValue, 0, len(m))
	for key, value := range m {
		result = append(result, stringKeyValue{Key: key, Value: value})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Key < result[j].Key
	})

	return result
}

// formatInt converts an int to its string representation.
func formatInt(v int) string {
	return fmt.Sprintf("%d", v)
}

// formatBool converts a bool to its string representation.
func formatBool(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

// UnmarshalXML implements xml.Unmarshaler for WidgetNode.
// It reads the element and populates the Attrs map.
func (n *WidgetNode) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	n.XMLName = start.Name
	n.Attrs = make(map[string]string)

	// Extract attributes into map
	for _, attr := range start.Attr {
		n.Attrs[attr.Name.Local] = attr.Value
	}

	// Read content (children or end element)
	for {
		token, err := d.Token()
		if err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}

		switch t := token.(type) {
		case xml.StartElement:
			// Decode child element
			child := &WidgetNode{}
			if err := d.DecodeElement(child, &t); err != nil {
				return fmt.Errorf("failed to decode child: %w", err)
			}
			n.Children = append(n.Children, child)

		case xml.EndElement:
			// End of this element
			if t.Name == start.Name {
				return nil
			}
		}
	}
}
