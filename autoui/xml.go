package autoui

import (
	"encoding/xml"
	"fmt"
	"sort"
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
	buildNodeTree(root, widgets[1:])

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
	for key, value := range info.State {
		// Don't override position/state attributes
		if !isReservedAttr(key) {
			node.Attrs[key] = value
		}
	}

	// Add custom data
	for key, value := range info.CustomData {
		// Don't override position/state attributes
		if !isReservedAttr(key) {
			node.Attrs[key] = value
		}
	}

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
// It uses the Container type to determine parent-child relationships.
func buildNodeTree(root *WidgetNode, remaining []WidgetInfo) {
	if len(remaining) == 0 {
		return
	}

	// Track container contexts - first element is the root
	containerStack := []*WidgetNode{root}

	for _, info := range remaining {
		node := widgetInfoToNode(info)

		// If this is a Container, add to current context and push to stack
		if info.Type == "Container" {
			// Add to current container's children
			if len(containerStack) > 0 {
				containerStack[len(containerStack)-1].Children = append(
					containerStack[len(containerStack)-1].Children,
					node,
				)
			}
			// Push to stack (it becomes the new context)
			containerStack = append(containerStack, node)
		} else {
			// Non-container: add to current container context
			if len(containerStack) > 0 {
				containerStack[len(containerStack)-1].Children = append(
					containerStack[len(containerStack)-1].Children,
					node,
				)
			}
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

// isReservedAttr checks if an attribute name is reserved for position/state attributes.
func isReservedAttr(key string) bool {
	reserved := []string{"x", "y", "width", "height", "visible", "disabled"}
	for _, r := range reserved {
		if key == r {
			return true
		}
	}
	return false
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