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
		ID string `ae:"id"`
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
				ID string `ae:"id"`
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

// TestMarshalXML_DeepChain tests XML generation for a deep chain structure.
// Structure: Container -> Container -> Container (parent, child, grandchild).
// Expected XML: <Container id="root"><Container id="child"><Container id="grandchild"/></Container></Container>
func TestMarshalXML_DeepChain(t *testing.T) {
	// Create root container
	root := widget.NewContainer()
	root.GetWidget().Rect = image.Rect(0, 0, 800, 600)
	root.GetWidget().CustomData = map[string]string{"id": "root"}

	// Create child container (nested inside root)
	child := widget.NewContainer()
	child.GetWidget().Rect = image.Rect(0, 50, 800, 550)
	child.GetWidget().CustomData = map[string]string{"id": "child"}
	root.AddChild(child)

	// Create grandchild container (nested inside child)
	grandchild := widget.NewContainer()
	grandchild.GetWidget().Rect = image.Rect(0, 100, 800, 500)
	grandchild.GetWidget().CustomData = map[string]string{"id": "grandchild"}
	child.AddChild(grandchild)

	infoList := autoui.WalkTree(root)
	xmlData, err := autoui.MarshalWidgetTreeXML(infoList)
	if err != nil {
		t.Fatalf("MarshalWidgetTreeXML failed: %v", err)
	}

	xmlStr := string(xmlData)

	// Parse the XML to verify structure
	var node autoui.WidgetNode
	if err := xml.Unmarshal(xmlData, &node); err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	// Verify root has one child (child container)
	if len(node.Children) != 1 {
		t.Fatalf("Expected UI to have 1 child (root), got %d", len(node.Children))
	}
	rootNode := node.Children[0]
	if rootNode.Attrs["id"] != "root" {
		t.Errorf("Expected root id='root', got '%s'", rootNode.Attrs["id"])
	}

	// Verify root contains child container
	if len(rootNode.Children) != 1 {
		t.Fatalf("Expected root to have 1 child, got %d: %s", len(rootNode.Children), xmlStr)
	}
	childNode := rootNode.Children[0]
	if childNode.Attrs["id"] != "child" {
		t.Errorf("Expected child id='child', got '%s'", childNode.Attrs["id"])
	}

	// Verify child contains grandchild container
	if len(childNode.Children) != 1 {
		t.Fatalf("Expected child to have 1 child (grandchild), got %d: %s", len(childNode.Children), xmlStr)
	}
	grandchildNode := childNode.Children[0]
	if grandchildNode.Attrs["id"] != "grandchild" {
		t.Errorf("Expected grandchild id='grandchild', got '%s'", grandchildNode.Attrs["id"])
	}

	// Grandchild should have no children
	if len(grandchildNode.Children) != 0 {
		t.Errorf("Expected grandchild to have 0 children, got %d", len(grandchildNode.Children))
	}
}

// TestMarshalXML_WideTree tests XML generation for a wide tree structure.
// Structure: Container -> [Container, Container] (parent with two sibling children).
// Expected XML: <Container id="root"><Container id="child1"/><Container id="child2"/></Container>
// BUG: Currently this produces nested output instead of siblings due to buildNodeTree's stack approach.
func TestMarshalXML_WideTree(t *testing.T) {
	// Create root container
	root := widget.NewContainer()
	root.GetWidget().Rect = image.Rect(0, 0, 800, 600)
	root.GetWidget().CustomData = map[string]string{"id": "root"}

	// Create first child container
	child1 := widget.NewContainer()
	child1.GetWidget().Rect = image.Rect(0, 0, 400, 600)
	child1.GetWidget().CustomData = map[string]string{"id": "child1"}
	root.AddChild(child1)

	// Create second child container (sibling of child1, NOT nested)
	child2 := widget.NewContainer()
	child2.GetWidget().Rect = image.Rect(400, 0, 800, 600)
	child2.GetWidget().CustomData = map[string]string{"id": "child2"}
	root.AddChild(child2)

	infoList := autoui.WalkTree(root)
	xmlData, err := autoui.MarshalWidgetTreeXML(infoList)
	if err != nil {
		t.Fatalf("MarshalWidgetTreeXML failed: %v", err)
	}

	xmlStr := string(xmlData)

	// Parse the XML to verify structure
	var node autoui.WidgetNode
	if err := xml.Unmarshal(xmlData, &node); err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	// Verify root has one child (root container)
	if len(node.Children) != 1 {
		t.Fatalf("Expected UI to have 1 child (root), got %d", len(node.Children))
	}
	rootNode := node.Children[0]
	if rootNode.Attrs["id"] != "root" {
		t.Errorf("Expected root id='root', got '%s'", rootNode.Attrs["id"])
	}

	// CRITICAL: Verify root contains BOTH child1 AND child2 as SIBLINGS
	// Current implementation BUG: child2 becomes nested inside child1 instead
	if len(rootNode.Children) != 2 {
		t.Errorf("Expected root to have 2 children (child1 and child2 as siblings), got %d\nXML:\n%s",
			len(rootNode.Children), xmlStr)
	}

	// If we have 2 children, verify they are the correct siblings
	if len(rootNode.Children) == 2 {
		if rootNode.Children[0].Attrs["id"] != "child1" {
			t.Errorf("Expected first child id='child1', got '%s'", rootNode.Children[0].Attrs["id"])
		}
		if rootNode.Children[1].Attrs["id"] != "child2" {
			t.Errorf("Expected second child id='child2', got '%s'", rootNode.Children[1].Attrs["id"])
		}
		// Both siblings should have no children (they're leaf containers)
		if len(rootNode.Children[0].Children) != 0 {
			t.Errorf("Expected child1 to have 0 children, got %d", len(rootNode.Children[0].Children))
		}
		if len(rootNode.Children[1].Children) != 0 {
			t.Errorf("Expected child2 to have 0 children, got %d", len(rootNode.Children[1].Children))
		}
	}
}

// TestWidgetXML_Addr tests that _addr appears in XML output.
func TestWidgetXML_Addr(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 100, 100)
	container.GetWidget().CustomData = map[string]string{"id": "test"}

	info := autoui.ExtractWidgetInfo(container)

	xmlData, err := autoui.MarshalWidgetXML(info)
	if err != nil {
		t.Fatalf("MarshalWidgetXML failed: %v", err)
	}

	// Check _addr appears in output
	if !strings.Contains(string(xmlData), "_addr=") {
		t.Errorf("Expected _addr attribute in XML, got: %s", xmlData)
	}

	// Check _addr format (hex string with 0x prefix)
	if !strings.Contains(string(xmlData), "_addr=\"0x") {
		t.Errorf("Expected _addr to be hex format (0x...), got: %s", xmlData)
	}
}

// TestMarshalXML_DistinguishesTreeStructures demonstrates that DeepChain and WideTree
// produce DIFFERENT XML structures. If they produce the same XML, there's a bug.
func TestMarshalXML_DistinguishesTreeStructures(t *testing.T) {
	// Build deep chain: root -> child -> grandchild
	deepRoot := widget.NewContainer()
	deepRoot.GetWidget().Rect = image.Rect(0, 0, 800, 600)
	deepRoot.GetWidget().CustomData = map[string]string{"id": "root"}

	deepChild := widget.NewContainer()
	deepChild.GetWidget().Rect = image.Rect(0, 50, 800, 550)
	deepChild.GetWidget().CustomData = map[string]string{"id": "child"}
	deepRoot.AddChild(deepChild)

	deepGrandchild := widget.NewContainer()
	deepGrandchild.GetWidget().Rect = image.Rect(0, 100, 800, 500)
	deepGrandchild.GetWidget().CustomData = map[string]string{"id": "grandchild"}
	deepChild.AddChild(deepGrandchild)

	// Build wide tree: root -> [child1, child2]
	wideRoot := widget.NewContainer()
	wideRoot.GetWidget().Rect = image.Rect(0, 0, 800, 600)
	wideRoot.GetWidget().CustomData = map[string]string{"id": "root"}

	wideChild1 := widget.NewContainer()
	wideChild1.GetWidget().Rect = image.Rect(0, 0, 400, 600)
	wideChild1.GetWidget().CustomData = map[string]string{"id": "child1"}
	wideRoot.AddChild(wideChild1)

	wideChild2 := widget.NewContainer()
	wideChild2.GetWidget().Rect = image.Rect(400, 0, 800, 600)
	wideChild2.GetWidget().CustomData = map[string]string{"id": "child2"}
	wideRoot.AddChild(wideChild2)

	// Generate XML for both
	deepXML, err := autoui.MarshalWidgetTreeXML(autoui.WalkTree(deepRoot))
	if err != nil {
		t.Fatalf("Deep chain marshal failed: %v", err)
	}

	wideXML, err := autoui.MarshalWidgetTreeXML(autoui.WalkTree(wideRoot))
	if err != nil {
		t.Fatalf("Wide tree marshal failed: %v", err)
	}

	deepStr := string(deepXML)
	wideStr := string(wideXML)

	// The two structures should produce DIFFERENT XML
	// Deep chain: nested 3 levels
	// Wide tree: nested 2 levels with 2 siblings
	if deepStr == wideStr {
		t.Errorf("BUG: Deep chain and wide tree produce identical XML!\nDeep:\n%s\nWide:\n%s",
			deepStr, wideStr)
	}
}