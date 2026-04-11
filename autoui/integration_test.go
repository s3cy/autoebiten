package autoui_test

import (
	"encoding/json"
	"image"
	"image/color"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/s3cy/autoebiten"
	"github.com/s3cy/autoebiten/autoui"
	"github.com/s3cy/autoebiten/internal/custom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegration_FullWorkflow tests the complete registration and command execution flow.
// This test verifies that all commands work end-to-end with real ebitenui widgets.
func TestIntegration_FullWorkflow(t *testing.T) {
	// Clean up any existing commands from previous tests
	for _, name := range autoebiten.ListCustomCommands() {
		autoebiten.Unregister(name)
	}

	// Reset UI reference
	autoui.GetUI() // This is safe to call even if nil

	// 1. Build real UI tree with multiple widgets
	root := widget.NewContainer()
	root.GetWidget().Rect = image.Rect(0, 0, 800, 600)

	// Create button images for testing
	buttonImage := createButtonImage()

	buttonColor := &widget.ButtonTextColor{
		Idle:     color.White,
		Disabled: color.Gray{128},
	}

	// Create first button with custom data
	btn1 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Button One", nil, buttonColor),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(struct {
				ID   string `ae:"id"`
				Role string `ae:"role"`
			}{
				ID:   "btn-1",
				Role: "primary",
			}),
		),
	)
	btn1.GetWidget().Rect = image.Rect(50, 50, 150, 80)
	root.AddChild(btn1)

	// Create second button
	btn2 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Button Two", nil, buttonColor),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(struct {
				ID   string `ae:"id"`
				Role string `ae:"role"`
			}{
				ID:   "btn-2",
				Role: "secondary",
			}),
		),
	)
	btn2.GetWidget().Rect = image.Rect(50, 100, 150, 130)
	root.AddChild(btn2)

	// Create a TextInput for testing method invocation
	textInput := widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(200, 30),
			widget.WidgetOpts.CustomData(struct {
				ID string `ae:"id"`
			}{
				ID: "input-1",
			}),
		),
	)
	textInput.GetWidget().Rect = image.Rect(50, 200, 250, 230)
	root.AddChild(textInput)

	// Create the UI
	ui := &ebitenui.UI{Container: root}

	// 2. Register autoui commands
	autoui.Register(ui)

	// Verify all commands are registered
	registeredCommands := autoebiten.ListCustomCommands()
	expectedCommands := []string{
		"autoui.tree",
		"autoui.at",
		"autoui.find",
		"autoui.xpath",
		"autoui.call",
		"autoui.highlight",
	}

	for _, expected := range expectedCommands {
		found := false
		for _, registered := range registeredCommands {
			if registered == expected {
				found = true
				break
			}
		}
		require.True(t, found, "Expected command '%s' to be registered", expected)
	}

	// 3. Test tree command
	t.Run("tree command", func(t *testing.T) {
		result := executeCommand("autoui.tree", "")
		assert.Contains(t, result, "<UI>")
		assert.Contains(t, result, "</UI>")
		assert.Contains(t, result, "<Container")
		assert.Contains(t, result, "<Button")
		assert.Contains(t, result, "<TextInput")
		// Verify structure contains expected widget types
		buttonCount := strings.Count(result, "<Button")
		assert.Equal(t, 2, buttonCount, "Expected 2 Button elements in tree")
		textInputCount := strings.Count(result, "<TextInput")
		assert.Equal(t, 1, textInputCount, "Expected 1 TextInput element in tree")
	})

	// 4. Test at command
	t.Run("at command with coordinates", func(t *testing.T) {
		// Find button at coordinates (100, 65) - should be btn1
		result := executeCommand("autoui.at", "100,65")
		assert.NotContains(t, result, "error:")
		assert.Contains(t, result, "<Button")
		assert.Contains(t, result, "x=\"50\"")
		assert.Contains(t, result, "y=\"50\"")
	})

	t.Run("at command with JSON coordinates", func(t *testing.T) {
		result := executeCommand("autoui.at", "{\"x\":100,\"y\":115}")
		assert.NotContains(t, result, "error:")
		assert.Contains(t, result, "<Button")
		// Should find btn2 at (100, 115)
		assert.Contains(t, result, "y=\"100\"")
	})

	t.Run("at command outside container", func(t *testing.T) {
		// Use coordinates outside the root container (900, 900)
		result := executeCommand("autoui.at", "900,900")
		assert.Contains(t, result, "error: no widget found")
	})

	// 5. Test find command
	t.Run("find command by type", func(t *testing.T) {
		result := executeCommand("autoui.find", "type=Button")
		assert.NotContains(t, result, "error:")
		assert.Contains(t, result, "<Button")
		buttonCount := strings.Count(result, "<Button")
		assert.Equal(t, 2, buttonCount, "Expected 2 buttons found")
	})

	t.Run("find command TextInput", func(t *testing.T) {
		result := executeCommand("autoui.find", "type=TextInput")
		assert.NotContains(t, result, "error:")
		assert.Contains(t, result, "<TextInput")
	})

	t.Run("find command by custom data", func(t *testing.T) {
		result := executeCommand("autoui.find", "id=btn-1")
		assert.NotContains(t, result, "error:")
		assert.Contains(t, result, "<Button")
		assert.Contains(t, result, "id=\"btn-1\"")
		buttonCount := strings.Count(result, "<Button")
		assert.Equal(t, 1, buttonCount, "Expected 1 button with id=btn-1")
	})

	t.Run("find command by role", func(t *testing.T) {
		result := executeCommand("autoui.find", "role=primary")
		assert.NotContains(t, result, "error:")
		assert.Contains(t, result, "<Button")
		assert.Contains(t, result, "role=\"primary\"")
	})

	t.Run("find command with JSON", func(t *testing.T) {
		result := executeCommand("autoui.find", "{\"type\":\"Container\"}")
		assert.NotContains(t, result, "error:")
		assert.Contains(t, result, "<Container")
	})

	t.Run("find command no match", func(t *testing.T) {
		result := executeCommand("autoui.find", "id=nonexistent")
		assert.Contains(t, result, "error: no widgets found")
	})

	// 6. Test xpath command
	t.Run("xpath command find all buttons", func(t *testing.T) {
		result := executeCommand("autoui.xpath", "//Button")
		assert.NotContains(t, result, "error:")
		assert.Contains(t, result, "<Button")
		buttonCount := strings.Count(result, "<Button")
		assert.Equal(t, 2, buttonCount, "Expected 2 buttons from XPath")
	})

	t.Run("xpath command find by attribute", func(t *testing.T) {
		result := executeCommand("autoui.xpath", "//Button[@id='btn-1']")
		assert.NotContains(t, result, "error:")
		assert.Contains(t, result, "<Button")
		assert.Contains(t, result, "id=\"btn-1\"")
		buttonCount := strings.Count(result, "<Button")
		assert.Equal(t, 1, buttonCount, "Expected 1 button with id='btn-1'")
	})

	t.Run("xpath command find TextInput", func(t *testing.T) {
		result := executeCommand("autoui.xpath", "//TextInput")
		assert.NotContains(t, result, "error:")
		assert.Contains(t, result, "<TextInput")
	})

	t.Run("xpath command invalid expression", func(t *testing.T) {
		result := executeCommand("autoui.xpath", "invalid/[xpath")
		assert.Contains(t, result, "error:")
	})

	// 7. Test highlight command
	t.Run("highlight command clear", func(t *testing.T) {
		result := executeCommand("autoui.highlight", "clear")
		assert.Contains(t, result, "ok: highlights cleared")
	})

	t.Run("highlight command by coordinates", func(t *testing.T) {
		result := executeCommand("autoui.highlight", "100,65")
		assert.Contains(t, result, "ok: highlighted widget")
		// Clear after test
		executeCommand("autoui.highlight", "clear")
	})

	t.Run("highlight command by query", func(t *testing.T) {
		result := executeCommand("autoui.highlight", "type=Button")
		assert.Contains(t, result, "ok: highlighted")
		assert.Contains(t, result, "widgets")
		// Clear after test
		executeCommand("autoui.highlight", "clear")
	})

	// 8. Test call command with TextInput (has Focus method)
	t.Run("call command focus widget", func(t *testing.T) {
		// TextInput has Focus(bool) method which is whitelisted
		result := executeCommand("autoui.call", "{\"target\":\"id=input-1\",\"method\":\"Focus\",\"args\":[true]}")
		assert.Contains(t, result, "\"success\":true")

		// Unfocus for cleanup
		executeCommand("autoui.call", "{\"target\":\"id=input-1\",\"method\":\"Focus\",\"args\":[false]}")
	})

	t.Run("call command click button", func(t *testing.T) {
		// Button has Click() method with no args
		result := executeCommand("autoui.call", "{\"target\":\"id=btn-1\",\"method\":\"Click\",\"args\":[]}")
		assert.Contains(t, result, "\"success\":true")
	})

	t.Run("call command invalid method", func(t *testing.T) {
		result := executeCommand("autoui.call", "{\"target\":\"id=btn-1\",\"method\":\"NonExistentMethod\",\"args\":[]}")
		// Error is returned in JSON format with success=false
		assert.Contains(t, result, "\"success\":false")
	})

	t.Run("call command missing target", func(t *testing.T) {
		result := executeCommand("autoui.call", "{\"method\":\"Click\",\"args\":[]}")
		assert.Contains(t, result, "error: missing target query")
	})

	t.Run("call command invalid JSON", func(t *testing.T) {
		result := executeCommand("autoui.call", "not json")
		assert.Contains(t, result, "error:")
	})

	// 9. Cleanup - unregister all commands
	for _, cmd := range expectedCommands {
		autoebiten.Unregister(cmd)
	}

	// Verify all commands are unregistered
	registeredCommands = autoebiten.ListCustomCommands()
	for _, cmd := range expectedCommands {
		assert.NotContains(t, registeredCommands, cmd, "Command '%s' should be unregistered", cmd)
	}
}

// TestIntegration_CommandsAfterUnregister tests that commands are properly cleaned up.
func TestIntegration_CommandsAfterUnregister(t *testing.T) {
	// Clean up any existing commands
	for _, name := range autoebiten.ListCustomCommands() {
		autoebiten.Unregister(name)
	}

	// Create a simple UI
	root := widget.NewContainer()
	root.GetWidget().Rect = image.Rect(0, 0, 100, 100)

	ui := &ebitenui.UI{Container: root}

	// Register autoui
	autoui.Register(ui)

	// Verify commands work
	result := executeCommand("autoui.tree", "")
	assert.Contains(t, result, "<UI>")

	// Unregister all autoui commands
	for _, name := range autoebiten.ListCustomCommands() {
		if strings.HasPrefix(name, "autoui.") {
			autoebiten.Unregister(name)
		}
	}

	// Verify commands no longer work
	result = executeCommand("autoui.tree", "")
	assert.Contains(t, result, "error: command not found")

	result = executeCommand("autoui.find", "type=Container")
	assert.Contains(t, result, "error: command not found")

	// Re-register should work
	autoui.Register(ui)
	result = executeCommand("autoui.tree", "")
	assert.Contains(t, result, "<UI>")

	// Final cleanup
	for _, name := range autoebiten.ListCustomCommands() {
		autoebiten.Unregister(name)
	}
}

// TestIntegration_CustomPrefix tests registering with a custom prefix.
func TestIntegration_CustomPrefix(t *testing.T) {
	// Clean up any existing commands
	for _, name := range autoebiten.ListCustomCommands() {
		autoebiten.Unregister(name)
	}

	// Create a simple UI
	root := widget.NewContainer()
	root.GetWidget().Rect = image.Rect(0, 0, 100, 100)
	ui := &ebitenui.UI{Container: root}

	// Register with custom prefix
	autoui.RegisterWithPrefix(ui, "custom.prefix")

	// Verify commands with custom prefix work
	result := executeCommand("custom.prefix.tree", "")
	assert.Contains(t, result, "<UI>")

	result = executeCommand("custom.prefix.find", "type=Container")
	assert.NotContains(t, result, "error:")
	assert.Contains(t, result, "<Container")

	// Verify default prefix doesn't work
	result = executeCommand("autoui.tree", "")
	assert.Contains(t, result, "error: command not found")

	// Cleanup
	for _, name := range autoebiten.ListCustomCommands() {
		autoebiten.Unregister(name)
	}
}

// executeCommand is a helper function that executes a command and returns the result.
func executeCommand(name, request string) string {
	handler := autoebiten.GetCustomCommand(name)
	if handler == nil {
		return "error: command not found"
	}

	var result string
	ctx := custom.NewContext(request, func(s string) {
		result = s
	})
	handler(ctx)
	return result
}

// createButtonImage creates a test button image for integration tests.
func createButtonImage() *widget.ButtonImage {
	return &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 30, color.RGBA{100, 100, 100, 255}),
		Pressed:  createTestNineSlice(100, 30, color.RGBA{80, 80, 80, 255}),
		Disabled: createTestNineSlice(100, 30, color.RGBA{150, 150, 150, 255}),
	}
}

func TestExistsResponse_WaitForCompatible(t *testing.T) {
	// Verify response format works with wait-for's JSON comparison

	// Simulate game response for found=true
	resp := autoui.ExistsResponse{Found: true, Count: 1}
	data, err := json.Marshal(resp)
	require.NoError(t, err)

	// Parse as JSON (wait-for does this)
	var parsed map[string]any
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	// Verify structure
	assert.Equal(t, true, parsed["found"])
	assert.Equal(t, 1.0, parsed["count"]) // JSON numbers are float64

	// Simulate comparison that wait-for would do
	expected := map[string]any{"found": true}
	assert.Equal(t, expected["found"], parsed["found"])
}

// TestIntegration_RadioGroupCallHandler tests the radiogroup=name target syntax in call command.
func TestIntegration_RadioGroupCallHandler(t *testing.T) {
	// Clean up any existing commands from previous tests
	for _, name := range autoebiten.ListCustomCommands() {
		autoebiten.Unregister(name)
	}

	// Create buttons for RadioGroup elements
	buttonImage := createButtonImage()
	buttonColor := &widget.ButtonTextColor{
		Idle:     color.White,
		Disabled: color.Gray{128},
	}

	btn1 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Option A", nil, buttonColor),
	)
	btn1.GetWidget().Rect = image.Rect(10, 10, 110, 40)

	btn2 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Option B", nil, buttonColor),
	)
	btn2.GetWidget().Rect = image.Rect(10, 50, 110, 80)

	// Create RadioGroup with elements
	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2))

	// Create UI with root container
	root := widget.NewContainer()
	root.GetWidget().Rect = image.Rect(0, 0, 200, 100)
	ui := &ebitenui.UI{Container: root}

	// Register autoui commands
	autoui.Register(ui)

	// Register RadioGroup
	autoui.RegisterRadioGroup("test-call-group", rg)

	// Test ActiveIndex method via radiogroup=name target
	request := `{"target":"radiogroup=test-call-group","method":"ActiveIndex","args":[]}`
	result := executeCommand("autoui.call", request)

	// Cleanup
	autoui.UnregisterRadioGroup("test-call-group")
	for _, name := range autoebiten.ListCustomCommands() {
		autoebiten.Unregister(name)
	}

	// Verify response is successful
	assert.Contains(t, result, `"success":true`)

	// Verify result contains active_index (should be 0 for first element or -1 if none selected)
	var resp autoui.CallResponse
	err := json.Unmarshal([]byte(result), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)

	// Result should be present for getter methods
	if resp.Result != nil {
		// Result should be an integer (index)
		assert.Contains(t, result, `"result"`)
	}
}

// TestIntegration_RadioGroupCallHandler_NotFound tests error handling for unregistered RadioGroup.
func TestIntegration_RadioGroupCallHandler_NotFound(t *testing.T) {
	// Clean up any existing commands from previous tests
	for _, name := range autoebiten.ListCustomCommands() {
		autoebiten.Unregister(name)
	}

	// Create simple UI
	root := widget.NewContainer()
	root.GetWidget().Rect = image.Rect(0, 0, 100, 100)
	ui := &ebitenui.UI{Container: root}

	// Register autoui commands
	autoui.Register(ui)

	// Try to call a RadioGroup that is NOT registered
	request := `{"target":"radiogroup=nonexistent","method":"ActiveIndex","args":[]}`
	result := executeCommand("autoui.call", request)

	// Cleanup
	for _, name := range autoebiten.ListCustomCommands() {
		autoebiten.Unregister(name)
	}

	// Verify error response
	assert.Contains(t, result, "error:")
	assert.Contains(t, result, "nonexistent")
	assert.Contains(t, result, "not registered")
}

// TestIntegration_RadioGroupFullWorkflow tests full RadioGroup automation workflow.
func TestIntegration_RadioGroupFullWorkflow(t *testing.T) {
	// Clean up any existing commands from previous tests
	for _, name := range autoebiten.ListCustomCommands() {
		autoebiten.Unregister(name)
	}

	// Setup: Create buttons for RadioGroup elements
	buttonImage := createButtonImage()
	buttonColor := &widget.ButtonTextColor{
		Idle:     color.White,
		Disabled: color.Gray{128},
	}

	btn1 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Option A", nil, buttonColor),
	)
	btn1.GetWidget().Rect = image.Rect(10, 10, 110, 40)

	btn2 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Option B", nil, buttonColor),
	)
	btn2.GetWidget().Rect = image.Rect(10, 50, 110, 80)

	btn3 := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Option C", nil, buttonColor),
	)
	btn3.GetWidget().Rect = image.Rect(10, 90, 110, 120)

	// Create RadioGroup with 3 elements
	rg := widget.NewRadioGroup(widget.RadioGroupOpts.Elements(btn1, btn2, btn3))

	// Create UI with root container
	root := widget.NewContainer()
	root.GetWidget().Rect = image.Rect(0, 0, 200, 150)
	ui := &ebitenui.UI{Container: root}

	// Register autoui commands
	autoui.Register(ui)

	// Register RadioGroup
	autoui.RegisterRadioGroup("options", rg)

	// Test 1: Get elements - returns type info (label may be empty without validation)
	t.Run("Elements", func(t *testing.T) {
		request := `{"target":"radiogroup=options","method":"Elements","args":[]}`
		result := executeCommand("autoui.call", request)
		assert.Contains(t, result, `"success":true`)
		// Elements returns list of element info with type
		assert.Contains(t, result, `"type":"Button"`)
	})

	// Test 2: Get active index - no selection initially (returns -1)
	t.Run("ActiveIndex_NoSelection", func(t *testing.T) {
		request := `{"target":"radiogroup=options","method":"ActiveIndex","args":[]}`
		result := executeCommand("autoui.call", request)
		assert.Contains(t, result, `"success":true`)
		// No selection initially, returns -1
		assert.Contains(t, result, `"result":-1`)
	})

	// Test 3: Set active by index to first element
	t.Run("SetActiveByIndex", func(t *testing.T) {
		request := `{"target":"radiogroup=options","method":"SetActiveByIndex","args":[0]}`
		result := executeCommand("autoui.call", request)
		assert.Contains(t, result, `"success":true`)
	})

	// Test 4: Verify active index after setting
	t.Run("ActiveIndex_AfterSet", func(t *testing.T) {
		request := `{"target":"radiogroup=options","method":"ActiveIndex","args":[]}`
		result := executeCommand("autoui.call", request)
		assert.Contains(t, result, `"success":true`)
		assert.Contains(t, result, `"result":0`)
	})

	// Test 5: Tree output includes RadioGroup
	t.Run("TreeOutput", func(t *testing.T) {
		result := executeCommand("autoui.tree", "")
		assert.Contains(t, result, "<RadioGroup")
		assert.Contains(t, result, "name=\"options\"")
	})

	// Cleanup
	autoui.UnregisterRadioGroup("options")
	for _, name := range autoebiten.ListCustomCommands() {
		autoebiten.Unregister(name)
	}
}
// TestIntegration_TabBookFullWorkflow tests full TabBook automation workflow.
func TestIntegration_TabBookFullWorkflow(t *testing.T) {
	// Clean up any existing commands from previous tests
	for _, name := range autoebiten.ListCustomCommands() {
		autoebiten.Unregister(name)
	}

	// Setup: Create TabBookTab with content
	tab1 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("General"))
	tab1.GetWidget().Rect = image.Rect(0, 50, 400, 250)

	// Add a button to tab1
	buttonImage := createButtonImage()
	buttonColor := &widget.ButtonTextColor{
		Idle:     color.White,
		Disabled: color.Gray{128},
	}
	saveBtn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Save", nil, buttonColor),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(map[string]string{"id": "save-btn"}),
		),
	)
	saveBtn.GetWidget().Rect = image.Rect(10, 60, 110, 90)
	tab1.AddChild(saveBtn)

	// Create second tab (disabled)
	tab2 := widget.NewTabBookTab(widget.TabBookTabOpts.Label("Settings"))
	tab2.GetWidget().Rect = image.Rect(0, 50, 400, 250)
	tab2.Disabled = true

	// Create TabBook and set tabs via reflection (same pattern as tree tests)
	tb := widget.NewTabBook()
	tb.GetWidget().Rect = image.Rect(0, 0, 400, 300)
	setPrivateFieldOnTabBookIntegration(tb, "tabs", []*widget.TabBookTab{tab1, tab2})
	setPrivateFieldOnTabBookIntegration(tb, "tab", tab1) // Set active tab

	// Create root container with TabBook
	root := widget.NewContainer()
	root.GetWidget().Rect = image.Rect(0, 0, 500, 350)
	root.AddChild(tb)

	// Create UI
	ui := &ebitenui.UI{Container: root}

	// Register autoui commands
	autoui.Register(ui)

	// Test 1: Get tabs
	t.Run("Tabs", func(t *testing.T) {
		request := `{"target":"type=TabBook","method":"Tabs","args":[]}`
		result := executeCommand("autoui.call", request)
		assert.Contains(t, result, `"success":true`)
		assert.Contains(t, result, "General")
		assert.Contains(t, result, "Settings")
		assert.Contains(t, result, `"disabled":true`)
	})

	// Test 2: Get active tab label
	t.Run("TabLabel", func(t *testing.T) {
		request := `{"target":"type=TabBook","method":"TabLabel","args":[]}`
		result := executeCommand("autoui.call", request)
		assert.Contains(t, result, `"success":true`)
		assert.Contains(t, result, "General")
	})

	// Test 3: Get tab index
	t.Run("TabIndex", func(t *testing.T) {
		request := `{"target":"type=TabBook","method":"TabIndex","args":[]}`
		result := executeCommand("autoui.call", request)
		assert.Contains(t, result, `"success":true`)
		// First tab is active, index should be 0
		assert.Contains(t, result, `"result":0`)
	})

	// Test 4: Tree output includes TabBookTab
	t.Run("TreeOutput", func(t *testing.T) {
		result := executeCommand("autoui.tree", "")
		assert.Contains(t, result, "<TabBook")
		assert.Contains(t, result, "<TabBookTab")
		assert.Contains(t, result, "label=\"General\"")
		assert.Contains(t, result, "label=\"Settings\"")
		assert.Contains(t, result, "disabled=\"true\"")
	})

	// Test 5: Set tab by index (try to set to disabled tab - should fail)
	t.Run("SetTabByIndexDisabled", func(t *testing.T) {
		request := `{"target":"type=TabBook","method":"SetTabByIndex","args":[1]}`
		result := executeCommand("autoui.call", request)
		// Should fail because tab 1 is disabled
		assert.Contains(t, result, `"success":false`)
		assert.Contains(t, result, "disabled")
	})

	// Cleanup
	for _, name := range autoebiten.ListCustomCommands() {
		autoebiten.Unregister(name)
	}
}

// setPrivateFieldOnTabBookIntegration is a test helper to set private fields on TabBook.
func setPrivateFieldOnTabBookIntegration(tb *widget.TabBook, fieldName string, value any) {
	v := reflect.ValueOf(tb).Elem()
	field := v.FieldByName(fieldName)
	fieldPtr := unsafe.Pointer(field.UnsafeAddr())
	realField := reflect.NewAt(field.Type(), fieldPtr).Elem()
	realField.Set(reflect.ValueOf(value))
}
