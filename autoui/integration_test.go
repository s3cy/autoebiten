package autoui_test

import (
	"encoding/json"
	"image"
	"image/color"
	"strings"
	"testing"

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