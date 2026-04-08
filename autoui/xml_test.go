package autoui_test

import (
	"encoding/xml"
	"image"
	"image/color"
	"strings"
	"testing"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/s3cy/autoebiten/autoui"
)

// TestMarshalXML_Button tests marshaling a single button to XML.
func TestMarshalXML_Button(t *testing.T) {
	// Create a button with custom data
	customData := struct {
		ID string `xml:"id,attr"`
	}{
		ID: "btn-submit",
	}

	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
		Disabled: createTestNineSlice(100, 30, color.RGBA{150, 150, 150, 255}),
	}

	buttonColor := &widget.ButtonTextColor{
		Idle:     color.White,
		Disabled: color.Gray{128},
	}

	btn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Click Me", nil, buttonColor),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(customData),
			widget.WidgetOpts.MinSize(100, 30),
		),
	)

	btn.GetWidget().Rect = image.Rect(10, 20, 110, 50)
	btn.GetWidget().Disabled = true

	// Extract widget info
	info := autoui.ExtractWidgetInfo(btn)

	// Marshal to XML
	xmlData, err := autoui.MarshalWidgetXML(info)
	if err != nil {
		t.Fatalf("MarshalWidgetXML failed: %v", err)
	}

	// Parse and verify XML structure
	var node autoui.WidgetNode
	if err := xml.Unmarshal(xmlData, &node); err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	// Verify element name is the widget type
	if node.XMLName.Local != "Button" {
		t.Errorf("Expected XMLName.Local to be 'Button', got '%s'", node.XMLName.Local)
	}

	// Verify position attributes
	if node.Attrs["x"] != "10" {
		t.Errorf("Expected x='10', got '%s'", node.Attrs["x"])
	}
	if node.Attrs["y"] != "20" {
		t.Errorf("Expected y='20', got '%s'", node.Attrs["y"])
	}
	if node.Attrs["width"] != "100" {
		t.Errorf("Expected width='100', got '%s'", node.Attrs["width"])
	}
	if node.Attrs["height"] != "30" {
		t.Errorf("Expected height='30', got '%s'", node.Attrs["height"])
	}

	// Verify state attributes
	if node.Attrs["visible"] != "true" {
		t.Errorf("Expected visible='true', got '%s'", node.Attrs["visible"])
	}
	if node.Attrs["disabled"] != "true" {
		t.Errorf("Expected disabled='true', got '%s'", node.Attrs["disabled"])
	}

	// Verify custom data
	if node.Attrs["id"] != "btn-submit" {
		t.Errorf("Expected id='btn-submit', got '%s'", node.Attrs["id"])
	}

	// Verify it's a self-closing tag (no children)
	xmlStr := string(xmlData)
	if !strings.Contains(xmlStr, "/>") && !strings.Contains(xmlStr, "</Button>") {
		t.Error("Expected self-closing or end tag for Button element")
	}
}

// TestMarshalXML_ContainerWithChildren tests marshaling a container tree to XML.
func TestMarshalXML_ContainerWithChildren(t *testing.T) {
	// Create root container with custom data
	rootData := map[string]string{
		"id":   "root-container",
		"role": "main",
	}
	root := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(rootData),
		),
	)
	root.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	// Create child button
	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
		Disabled: createTestNineSlice(100, 30, color.RGBA{150, 150, 150, 255}),
	}

	buttonColor := &widget.ButtonTextColor{
		Idle:     color.White,
		Disabled: color.Gray{128},
	}

	btn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Submit", nil, buttonColor),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(struct {
				ID string `xml:"id,attr"`
			}{ID: "btn-submit"}),
		),
	)
	btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	root.AddChild(btn)

	// Walk tree and marshal
	infoList := autoui.WalkTree(root)
	xmlData, err := autoui.MarshalWidgetTreeXML(infoList)
	if err != nil {
		t.Fatalf("MarshalWidgetTreeXML failed: %v", err)
	}

	// Verify XML structure
	xmlStr := string(xmlData)

	// Should have UI root element
	if !strings.Contains(xmlStr, "<UI") {
		t.Error("Expected <UI> root element")
	}

	// Should have Container with children
	if !strings.Contains(xmlStr, "<Container") {
		t.Error("Expected <Container> element")
	}

	// Should have nested Button
	if !strings.Contains(xmlStr, "<Button") {
		t.Error("Expected <Button> element inside Container")
	}

	// Verify attributes are present
	if !strings.Contains(xmlStr, "id=\"root-container\"") {
		t.Error("Expected id='root-container' attribute on Container")
	}
	if !strings.Contains(xmlStr, "id=\"btn-submit\"") {
		t.Error("Expected id='btn-submit' attribute on Button")
	}
}

// TestMarshalXML_SortedAttributes tests that attributes are sorted for deterministic output.
func TestMarshalXML_SortedAttributes(t *testing.T) {
	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
		Disabled: createTestNineSlice(100, 30, color.RGBA{150, 150, 150, 255}),
	}

	buttonColor := &widget.ButtonTextColor{
		Idle:     color.White,
		Disabled: color.Gray{128},
	}

	btn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Test", nil, buttonColor),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(map[string]string{
				"zAttr": "z-value",
				"aAttr": "a-value",
				"mAttr": "m-value",
			}),
		),
	)
	btn.GetWidget().Rect = image.Rect(0, 0, 100, 30)

	info := autoui.ExtractWidgetInfo(btn)
	xmlData, err := autoui.MarshalWidgetXML(info)
	if err != nil {
		t.Fatalf("MarshalWidgetXML failed: %v", err)
	}

	// Attributes should be sorted alphabetically in the output
	xmlStr := string(xmlData)

	// Find positions of custom attributes
	aPos := strings.Index(xmlStr, "aAttr=")
	mPos := strings.Index(xmlStr, "mAttr=")
	zPos := strings.Index(xmlStr, "zAttr=")

	if aPos == -1 || mPos == -1 || zPos == -1 {
		t.Fatalf("Missing custom attributes in output: %s", xmlStr)
	}

	// Verify sorted order: aAttr should come before mAttr, which should come before zAttr
	if !(aPos < mPos && mPos < zPos) {
		t.Errorf("Attributes not sorted alphabetically:\n%s", xmlStr)
	}
}

// TestMarshalXML_EmptyWidget tests handling of empty/zero widget info.
func TestMarshalXML_EmptyWidget(t *testing.T) {
	info := autoui.WidgetInfo{
		Type:       "Widget",
		Rect:       image.Rect(0, 0, 0, 0),
		State:      nil,
		CustomData: nil,
	}

	xmlData, err := autoui.MarshalWidgetXML(info)
	if err != nil {
		t.Fatalf("MarshalWidgetXML failed: %v", err)
	}

	var node autoui.WidgetNode
	if err := xml.Unmarshal(xmlData, &node); err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	if node.XMLName.Local != "Widget" {
		t.Errorf("Expected XMLName.Local to be 'Widget', got '%s'", node.XMLName.Local)
	}
}

// TestMarshalXML_NestedContainers tests XML generation with nested containers.
func TestMarshalXML_NestedContainers(t *testing.T) {
	// Create root container
	root := widget.NewContainer()
	root.GetWidget().Rect = image.Rect(0, 0, 800, 600)
	root.GetWidget().CustomData = map[string]string{"id": "root"}

	// Create child container
	child := widget.NewContainer()
	child.GetWidget().Rect = image.Rect(0, 50, 800, 550)
	child.GetWidget().CustomData = map[string]string{"id": "child"}
	root.AddChild(child)

	// Create button inside child container
	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
		Disabled: createTestNineSlice(100, 30, color.RGBA{150, 150, 150, 255}),
	}
	btn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
	)
	btn.GetWidget().Rect = image.Rect(20, 60, 120, 90)
	btn.GetWidget().CustomData = map[string]string{"id": "btn"}
	child.AddChild(btn)

	infoList := autoui.WalkTree(root)
	xmlData, err := autoui.MarshalWidgetTreeXML(infoList)
	if err != nil {
		t.Fatalf("MarshalWidgetTreeXML failed: %v", err)
	}

	xmlStr := string(xmlData)

	// Verify nested structure: root contains child, child contains btn
	rootIdx := strings.Index(xmlStr, "id=\"root\"")
	childIdx := strings.Index(xmlStr, "id=\"child\"")
	btnIdx := strings.Index(xmlStr, "id=\"btn\"")

	// Root should appear first
	if rootIdx == -1 || childIdx == -1 || btnIdx == -1 {
		t.Fatalf("Missing expected elements in XML:\n%s", xmlStr)
	}

	// Verify nesting order (simplified check)
	// The structure should be: <Container id="root">...<Container id="child">...<Button id="btn"/>
	if !(rootIdx < childIdx) {
		t.Error("Root container should appear before child container")
	}
	if !(childIdx < btnIdx) {
		t.Error("Child container should appear before button")
	}
}