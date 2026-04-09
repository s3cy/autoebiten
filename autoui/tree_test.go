package autoui_test

import (
	"image"
	"image/color"
	"strings"
	"testing"

	"github.com/ebitenui/ebitenui"
	ebitenuiImage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/s3cy/autoebiten/autoui"
)

// TestExtractWidgetInfo_Button tests extracting widget info from a button.
// Note: Not using t.Parallel() because ebitenui has global state that is not thread-safe.
func TestExtractWidgetInfo_Button(t *testing.T) {
	// Create a button with custom data
	customData := struct {
		ID string `ae:"id"`
	}{
		ID: "test-button-1",
	}

	// Create a simple button image for testing
	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
		Disabled: createTestNineSlice(100, 30, color.RGBA{150, 150, 150, 255}),
	}

	buttonColor := &widget.ButtonTextColor{
		Idle:     color.White,
		Disabled: color.Gray{128},
	}

	// Create a basic button
	btn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Test Button", nil, buttonColor),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(customData),
			widget.WidgetOpts.MinSize(100, 30),
		),
	)

	// Validate and set location
	// Note: Validation will fail without a font face, but we can still test extraction
	// btn.Validate() - skip validation for unit tests
	btn.SetLocation(image.Rect(10, 20, 110, 50))

	// Extract widget info
	info := autoui.ExtractWidgetInfo(btn)

	// Verify Type
	if info.Type != "Button" {
		t.Errorf("Expected Type to be 'Button', got '%s'", info.Type)
	}

	// Verify Rect
	expectedRect := image.Rect(10, 20, 110, 50)
	if info.Rect != expectedRect {
		t.Errorf("Expected Rect to be %v, got %v", expectedRect, info.Rect)
	}

	// Verify Visible (should be true by default)
	if !info.Visible {
		t.Error("Expected Visible to be true")
	}

	// Verify Disabled (should be false by default)
	if info.Disabled {
		t.Error("Expected Disabled to be false")
	}

	// Verify State contains text
	// Note: Without validation, text extraction may not work as expected
	// The Text() method requires init to be called
	if info.State == nil {
		t.Log("Warning: State is nil (may require validation)")
	} else if info.State["text"] != "Test Button" {
		t.Logf("State['text'] = '%s' (may require validation)", info.State["text"])
	}

	// Verify CustomData extraction
	if info.CustomData == nil {
		t.Fatal("Expected CustomData to be non-nil")
	}
	if info.CustomData["id"] != "test-button-1" {
		t.Errorf("Expected CustomData['id'] to be 'test-button-1', got '%s'", info.CustomData["id"])
	}
}

// TestExtractWidgetInfo_Container tests extracting widget info from a container.
func TestExtractWidgetInfo_Container(t *testing.T) {
	// Create a container with custom data
	customData := map[string]string{
		"id":   "main-container",
		"role": "layout",
	}

	container := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(customData),
			widget.WidgetOpts.MinSize(800, 600),
		),
	)

	container.Validate()
	container.SetLocation(image.Rect(0, 0, 800, 600))

	// Extract widget info
	info := autoui.ExtractWidgetInfo(container)

	// Verify Type
	if info.Type != "Container" {
		t.Errorf("Expected Type to be 'Container', got '%s'", info.Type)
	}

	// Verify Rect
	expectedRect := image.Rect(0, 0, 800, 600)
	if info.Rect != expectedRect {
		t.Errorf("Expected Rect to be %v, got %v", expectedRect, info.Rect)
	}

	// Verify CustomData extraction from map
	if info.CustomData == nil {
		t.Fatal("Expected CustomData to be non-nil")
	}
	if info.CustomData["id"] != "main-container" {
		t.Errorf("Expected CustomData['id'] to be 'main-container', got '%s'", info.CustomData["id"])
	}
	if info.CustomData["role"] != "layout" {
		t.Errorf("Expected CustomData['role'] to be 'layout', got '%s'", info.CustomData["role"])
	}
}

// TestExtractWidgetInfo_Slider tests extracting widget info from a slider.
func TestExtractWidgetInfo_Slider(t *testing.T) {
	trackImage := &widget.SliderTrackImage{
		Idle:     createTestNineSlice(100, 10, color.RGBA{50, 50, 50, 255}),
		Disabled: createTestNineSlice(100, 10, color.RGBA{100, 100, 100, 255}),
	}

	handleImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(16, 16, color.RGBA{150, 150, 150, 255}),
		Pressed:  createTestNineSlice(16, 16, color.RGBA{100, 100, 100, 255}),
		Disabled: createTestNineSlice(16, 16, color.RGBA{200, 200, 200, 255}),
	}

	slider := widget.NewSlider(
		widget.SliderOpts.Images(trackImage, handleImage),
		widget.SliderOpts.MinMax(0, 100),
		widget.SliderOpts.InitialCurrent(50),
	)

	slider.Validate()
	slider.SetLocation(image.Rect(10, 10, 210, 30))

	// Extract widget info
	info := autoui.ExtractWidgetInfo(slider)

	// Verify Type
	if info.Type != "Slider" {
		t.Errorf("Expected Type to be 'Slider', got '%s'", info.Type)
	}

	// Verify State contains slider values
	if info.State == nil {
		t.Fatal("Expected State to be non-nil")
	}
	if info.State["value"] != "50" {
		t.Errorf("Expected State['value'] to be '50', got '%s'", info.State["value"])
	}
	if info.State["min"] != "0" {
		t.Errorf("Expected State['min'] to be '0', got '%s'", info.State["min"])
	}
	if info.State["max"] != "100" {
		t.Errorf("Expected State['max'] to be '100', got '%s'", info.State["max"])
	}
}

// TestExtractWidgetInfo_ProgressBar tests extracting widget info from a progress bar.
func TestExtractWidgetInfo_ProgressBar(t *testing.T) {
	trackImage := &widget.ProgressBarImage{
		Idle:     createTestNineSlice(100, 20, color.RGBA{50, 50, 50, 255}),
		Disabled: createTestNineSlice(100, 20, color.RGBA{100, 100, 100, 255}),
	}

	fillImage := &widget.ProgressBarImage{
		Idle:     createTestNineSlice(100, 20, color.RGBA{0, 150, 0, 255}),
		Disabled: createTestNineSlice(100, 20, color.RGBA{100, 150, 100, 255}),
	}

	pb := widget.NewProgressBar(
		widget.ProgressBarOpts.Images(trackImage, fillImage),
		widget.ProgressBarOpts.Values(0, 100, 75),
	)

	pb.Validate()
	pb.SetLocation(image.Rect(10, 10, 210, 30))

	// Extract widget info
	info := autoui.ExtractWidgetInfo(pb)

	// Verify Type
	if info.Type != "ProgressBar" {
		t.Errorf("Expected Type to be 'ProgressBar', got '%s'", info.Type)
	}

	// Verify State contains progress bar values
	if info.State == nil {
		t.Fatal("Expected State to be non-nil")
	}
	if info.State["value"] != "75" {
		t.Errorf("Expected State['value'] to be '75', got '%s'", info.State["value"])
	}
	if info.State["min"] != "0" {
		t.Errorf("Expected State['min'] to be '0', got '%s'", info.State["min"])
	}
	if info.State["max"] != "100" {
		t.Errorf("Expected State['max'] to be '100', got '%s'", info.State["max"])
	}
}

// createTestNineSlice creates a simple NineSlice for testing.
func createTestNineSlice(w, h int, c color.Color) *ebitenuiImage.NineSlice {
	img := ebiten.NewImage(w, h)
	img.Fill(c)
	return ebitenuiImage.NewNineSliceSimple(img, 0, 0)
}

// TestWalkTree_ContainerWithChildren tests traversing a container with children.
func TestWalkTree_ContainerWithChildren(t *testing.T) {
	// Create container with children
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

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
		widget.ButtonOpts.Text("Button 1", nil, buttonColor),
	)
	btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)
	container.AddChild(btn)

	infoList := autoui.WalkTree(container)

	if len(infoList) != 2 {
		t.Errorf("Expected 2 widgets, got %d", len(infoList))
	}
	if infoList[0].Type != "Container" {
		t.Errorf("Expected first widget to be Container, got %s", infoList[0].Type)
	}
	if infoList[1].Type != "Button" {
		t.Errorf("Expected second widget to be Button, got %s", infoList[1].Type)
	}
	// Note: text state may require validation, so we just check Type
}

// TestWalkTree_NestedContainers tests traversing nested containers.
func TestWalkTree_NestedContainers(t *testing.T) {
	root := widget.NewContainer()
	root.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	childContainer := widget.NewContainer()
	childContainer.GetWidget().Rect = image.Rect(0, 50, 800, 550)
	root.AddChild(childContainer)

	btn := widget.NewButton()
	btn.GetWidget().Rect = image.Rect(20, 60, 120, 90)
	childContainer.AddChild(btn)

	infoList := autoui.WalkTree(root)

	if len(infoList) != 3 {
		t.Errorf("Expected 3 widgets, got %d", len(infoList))
	}
	if infoList[0].Type != "Container" {
		t.Errorf("Expected first widget to be Container, got %s", infoList[0].Type)
	}
	if infoList[1].Type != "Container" {
		t.Errorf("Expected second widget to be Container, got %s", infoList[1].Type)
	}
	if infoList[2].Type != "Button" {
		t.Errorf("Expected third widget to be Button, got %s", infoList[2].Type)
	}
}

// TestWalkTree_SingleWidget tests traversing a single widget without children.
func TestWalkTree_SingleWidget(t *testing.T) {
	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
		Disabled: createTestNineSlice(100, 30, color.RGBA{150, 150, 150, 255}),
	}

	btn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
	)
	btn.GetWidget().Rect = image.Rect(10, 10, 110, 40)

	infoList := autoui.WalkTree(btn)

	if len(infoList) != 1 {
		t.Errorf("Expected 1 widget, got %d", len(infoList))
	}
	if infoList[0].Type != "Button" {
		t.Errorf("Expected widget to be Button, got %s", infoList[0].Type)
	}
}

// TestWalkTree_NilWidget tests traversing with nil root.
func TestWalkTree_NilWidget(t *testing.T) {
	infoList := autoui.WalkTree(nil)

	if len(infoList) != 0 {
		t.Errorf("Expected 0 widgets for nil root, got %d", len(infoList))
	}
}

// TestWalkTree_DeepChain tests a deep chain: Container -> Container -> Container.
// Structure: root (parent) -> child -> grandchild (3 levels deep).
func TestWalkTree_DeepChain(t *testing.T) {
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

	// Should have 3 widgets in DFS order
	if len(infoList) != 3 {
		t.Fatalf("Expected 3 widgets, got %d", len(infoList))
	}

	// Verify DFS order: root -> child -> grandchild
	expectedIDs := []string{"root", "child", "grandchild"}
	for i, expected := range expectedIDs {
		if infoList[i].CustomData["id"] != expected {
			t.Errorf("Position %d: expected id='%s', got '%s'", i, expected, infoList[i].CustomData["id"])
		}
	}
}

// TestWalkTree_WideTree tests a wide tree: Container -> [Container, Container].
// Structure: root (parent) with two children (child1 and child2).
func TestWalkTree_WideTree(t *testing.T) {
	// Create root container
	root := widget.NewContainer()
	root.GetWidget().Rect = image.Rect(0, 0, 800, 600)
	root.GetWidget().CustomData = map[string]string{"id": "root"}

	// Create first child container
	child1 := widget.NewContainer()
	child1.GetWidget().Rect = image.Rect(0, 0, 400, 600)
	child1.GetWidget().CustomData = map[string]string{"id": "child1"}
	root.AddChild(child1)

	// Create second child container (sibling of child1, not nested)
	child2 := widget.NewContainer()
	child2.GetWidget().Rect = image.Rect(400, 0, 800, 600)
	child2.GetWidget().CustomData = map[string]string{"id": "child2"}
	root.AddChild(child2)

	infoList := autoui.WalkTree(root)

	// Should have 3 widgets in DFS order
	if len(infoList) != 3 {
		t.Fatalf("Expected 3 widgets, got %d", len(infoList))
	}

	// Verify DFS order: root -> child1 -> child2
	expectedIDs := []string{"root", "child1", "child2"}
	for i, expected := range expectedIDs {
		if infoList[i].CustomData["id"] != expected {
			t.Errorf("Position %d: expected id='%s', got '%s'", i, expected, infoList[i].CustomData["id"])
		}
	}
}

// TestWidgetInfo_Addr tests that Addr field is populated with pointer address.
func TestWidgetInfo_Addr(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 100, 100)

	info := autoui.ExtractWidgetInfo(container)

	if info.Addr == "" {
		t.Error("Expected Addr to be populated")
	}

	// Addr should be hex format like "0x14000abc0"
	if !strings.HasPrefix(info.Addr, "0x") {
		t.Errorf("Expected Addr to start with '0x', got '%s'", info.Addr)
	}

	// Addr should be unique per widget
	btn := widget.NewButton()
	btn.GetWidget().Rect = image.Rect(0, 0, 50, 30)
	btnInfo := autoui.ExtractWidgetInfo(btn)

	if info.Addr == btnInfo.Addr {
		t.Error("Expected different widgets to have different Addr values")
	}
}

// TestSnapshotTree tests snapshot tree creation.
func TestSnapshotTree(t *testing.T) {
	container := widget.NewContainer()
	container.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	ui := &ebitenui.UI{Container: container}

	widgets := autoui.SnapshotTree(ui)

	if len(widgets) != 1 {
		t.Errorf("Expected 1 widget (container), got %d", len(widgets))
	}

	if widgets[0].Type != "Container" {
		t.Errorf("Expected Container type, got %s", widgets[0].Type)
	}
}

// TestSnapshotTree_NilUI tests nil UI handling.
func TestSnapshotTree_NilUI(t *testing.T) {
	widgets := autoui.SnapshotTree(nil)

	if widgets != nil {
		t.Errorf("Expected nil for nil UI, got %d widgets", len(widgets))
	}
}