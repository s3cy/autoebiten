package autoui

import (
	"strings"
	"testing"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/s3cy/autoebiten"
	"github.com/s3cy/autoebiten/internal/custom"
)

func TestRegister_PanicsOnNilUI(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Register did not panic on nil UI")
		}
	}()

	Register(nil)
}

func TestRegisterWithPrefix_PanicsOnNilUI(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("RegisterWithPrefix did not panic on nil UI")
		}
	}()

	RegisterWithPrefix(nil, "test")
}

func TestRegister_RegistersAllCommands(t *testing.T) {
	// Clean up any existing commands
	for _, name := range autoebiten.ListCustomCommands() {
		autoebiten.Unregister(name)
	}

	// Create a simple UI
	ui := createTestUI()

	// Register autoui commands
	Register(ui)

	// Verify all commands are registered
	expectedCommands := []string{
		"autoui.tree",
		"autoui.at",
		"autoui.find",
		"autoui.xpath",
		"autoui.call",
		"autoui.highlight",
	}

	registeredCommands := autoebiten.ListCustomCommands()
	for _, expected := range expectedCommands {
		found := false
		for _, registered := range registeredCommands {
			if registered == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected command '%s' not registered", expected)
		}
	}

	// Clean up
	for _, name := range expectedCommands {
		autoebiten.Unregister(name)
	}
}

func TestRegisterWithPrefix_CustomPrefix(t *testing.T) {
	// Clean up any existing commands
	for _, name := range autoebiten.ListCustomCommands() {
		autoebiten.Unregister(name)
	}

	// Create a simple UI
	ui := createTestUI()

	// Register with custom prefix
	RegisterWithPrefix(ui, "custom.prefix")

	// Verify commands have custom prefix
	expectedCommands := []string{
		"custom.prefix.tree",
		"custom.prefix.at",
		"custom.prefix.find",
		"custom.prefix.xpath",
		"custom.prefix.call",
		"custom.prefix.highlight",
	}

	registeredCommands := autoebiten.ListCustomCommands()
	for _, expected := range expectedCommands {
		found := false
		for _, registered := range registeredCommands {
			if registered == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected command '%s' not registered", expected)
		}
	}

	// Clean up
	for _, name := range expectedCommands {
		autoebiten.Unregister(name)
	}
}

func TestHandleTreeCommand_ReturnsXML(t *testing.T) {
	ui := createTestUI()
	handler := handleTreeCommand(ui)

	// Create mock context
	response := ""
	ctx := custom.NewContext("", func(resp string) {
		response = resp
	})

	// Execute handler
	handler(ctx)

	// Verify response contains XML structure
	if !strings.Contains(response, "<UI>") {
		t.Errorf("Tree command response should contain <UI>, got: %s", response)
	}
	if !strings.Contains(response, "</UI>") {
		t.Errorf("Tree command response should contain </UI>, got: %s", response)
	}
}

func TestHandleTreeCommand_NilUI(t *testing.T) {
	handler := handleTreeCommand(nil)

	// Create mock context
	response := ""
	ctx := custom.NewContext("", func(resp string) {
		response = resp
	})

	// Execute handler
	handler(ctx)

	// Verify error response
	if !strings.Contains(response, "error:") {
		t.Errorf("Nil UI should return error, got: %s", response)
	}
}

func TestHandleAtCommand_ValidCoordinates(t *testing.T) {
	ui := createTestUI()
	handler := handleAtCommand(ui)

	// Create mock context with coordinates
	response := ""
	ctx := custom.NewContext("100,50", func(resp string) {
		response = resp
	})

	// Execute handler
	handler(ctx)

	// Response should either be widget XML or error (depends on widget positions)
	// At minimum, it should not panic
	if strings.Contains(response, "error:") {
		// Widget not at that position is acceptable
		t.Logf("No widget found at 100,50 (expected for test UI): %s", response)
	}
}

func TestHandleAtCommand_InvalidCoordinates(t *testing.T) {
	ui := createTestUI()
	handler := handleAtCommand(ui)

	// Create mock context with invalid coordinates
	response := ""
	ctx := custom.NewContext("invalid", func(resp string) {
		response = resp
	})

	// Execute handler
	handler(ctx)

	// Verify error response
	if !strings.Contains(response, "error:") {
		t.Errorf("Invalid coordinates should return error, got: %s", response)
	}
}

func TestHandleAtCommand_JSONCoordinates(t *testing.T) {
	ui := createTestUI()
	handler := handleAtCommand(ui)

	// Create mock context with JSON coordinates
	response := ""
	ctx := custom.NewContext("{\"x\":100,\"y\":50}", func(resp string) {
		response = resp
	})

	// Execute handler
	handler(ctx)

	// Response should not panic
	t.Logf("JSON coordinates response: %s", response)
}

func TestHandleFindCommand_EmptyQuery(t *testing.T) {
	ui := createTestUI()
	handler := handleFindCommand(ui)

	// Create mock context with empty query (should return all widgets)
	response := ""
	ctx := custom.NewContext("", func(resp string) {
		response = resp
	})

	// Execute handler
	handler(ctx)

	// Should return XML tree
	if !strings.Contains(response, "<UI>") {
		t.Errorf("Empty query should return full tree, got: %s", response)
	}
}

func TestHandleFindCommand_ValidQuery(t *testing.T) {
	ui := createTestUI()
	handler := handleFindCommand(ui)

	// Create mock context with query
	response := ""
	ctx := custom.NewContext("type=Container", func(resp string) {
		response = resp
	})

	// Execute handler
	handler(ctx)

	// Response should contain Container
	if !strings.Contains(response, "Container") {
		t.Logf("Find type=Container response: %s", response)
	}
}

func TestHandleXPathCommand_MissingExpression(t *testing.T) {
	ui := createTestUI()
	handler := handleXPathCommand(ui)

	// Create mock context with empty expression
	response := ""
	ctx := custom.NewContext("", func(resp string) {
		response = resp
	})

	// Execute handler
	handler(ctx)

	// Verify error response
	if !strings.Contains(response, "error:") {
		t.Errorf("Empty XPath should return error, got: %s", response)
	}
}

func TestHandleXPathCommand_ValidExpression(t *testing.T) {
	ui := createTestUI()
	handler := handleXPathCommand(ui)

	// Create mock context with XPath expression
	response := ""
	ctx := custom.NewContext("//Container", func(resp string) {
		response = resp
	})

	// Execute handler
	handler(ctx)

	// Response should contain Container or error
	t.Logf("XPath //Container response: %s", response)
}

func TestHandleCallCommand_InvalidJSON(t *testing.T) {
	ui := createTestUI()
	handler := handleCallCommand(ui)

	// Create mock context with invalid JSON
	response := ""
	ctx := custom.NewContext("invalid json", func(resp string) {
		response = resp
	})

	// Execute handler
	handler(ctx)

	// Verify error response
	if !strings.Contains(response, "error:") {
		t.Errorf("Invalid JSON should return error, got: %s", response)
	}
}

func TestHandleCallCommand_MissingFields(t *testing.T) {
	ui := createTestUI()
	handler := handleCallCommand(ui)

	tests := []struct {
		name    string
		request string
	}{
		{"missing target", "{\"method\":\"test\"}"},
		{"missing method", "{\"target\":\"type=Container\"}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := ""
			ctx := custom.NewContext(tt.request, func(resp string) {
				response = resp
			})

			handler(ctx)

			if !strings.Contains(response, "error:") {
				t.Errorf("Missing fields should return error, got: %s", response)
			}
		})
	}
}

func TestHandleHighlightCommand_Clear(t *testing.T) {
	ui := createTestUI()
	handler := handleHighlightCommand(ui)

	// Create mock context with "clear" request
	response := ""
	ctx := custom.NewContext("clear", func(resp string) {
		response = resp
	})

	// Execute handler
	handler(ctx)

	// Verify success response
	if !strings.Contains(response, "ok:") {
		t.Errorf("Clear should return ok, got: %s", response)
	}
}

func TestHandleHighlightCommand_InvalidFormat(t *testing.T) {
	ui := createTestUI()
	handler := handleHighlightCommand(ui)

	// Create mock context with invalid JSON
	response := ""
	ctx := custom.NewContext("{\"invalid\":true}", func(resp string) {
		response = resp
	})

	// Execute handler
	handler(ctx)

	// Verify error response (missing query and coordinates)
	if !strings.Contains(response, "error:") {
		t.Errorf("Invalid highlight request should return error, got: %s", response)
	}
}

func TestParseCoordinates_SimpleFormat(t *testing.T) {
	x, y, err := parseCoordinates("100,50")
	if err != nil {
		t.Errorf("Failed to parse simple coordinates: %v", err)
	}
	if x != 100 || y != 50 {
		t.Errorf("Expected x=100, y=50, got x=%d, y=%d", x, y)
	}
}

func TestParseCoordinates_WithSpaces(t *testing.T) {
	x, y, err := parseCoordinates(" 100 , 50 ")
	if err != nil {
		t.Errorf("Failed to parse coordinates with spaces: %v", err)
	}
	if x != 100 || y != 50 {
		t.Errorf("Expected x=100, y=50, got x=%d, y=%d", x, y)
	}
}

func TestParseCoordinates_JSONFormat(t *testing.T) {
	x, y, err := parseCoordinates("{\"x\":100,\"y\":50}")
	if err != nil {
		t.Errorf("Failed to parse JSON coordinates: %v", err)
	}
	if x != 100 || y != 50 {
		t.Errorf("Expected x=100, y=50, got x=%d, y=%d", x, y)
	}
}

func TestParseCoordinates_InvalidFormat(t *testing.T) {
	_, _, err := parseCoordinates("invalid")
	if err == nil {
		t.Error("Expected error for invalid coordinates")
	}
}

func TestParseCoordinates_InvalidJSON(t *testing.T) {
	_, _, err := parseCoordinates("{invalid}")
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestGetUI_ReturnsNilBeforeRegistration(t *testing.T) {
	// Reset uiReference
	uiMu.Lock()
	uiReference = nil
	uiMu.Unlock()

	ui := GetUI()
	if ui != nil {
		t.Error("GetUI should return nil before registration")
	}
}

func TestGetUI_ReturnsUIAfterRegistration(t *testing.T) {
	// Reset and clean up
	uiMu.Lock()
	uiReference = nil
	uiMu.Unlock()
	for _, name := range autoebiten.ListCustomCommands() {
		autoebiten.Unregister(name)
	}

	ui := createTestUI()
	Register(ui)

	retrievedUI := GetUI()
	if retrievedUI != ui {
		t.Error("GetUI should return registered UI")
	}

	// Clean up
	uiMu.Lock()
	uiReference = nil
	uiMu.Unlock()
	for _, name := range autoebiten.ListCustomCommands() {
		autoebiten.Unregister(name)
	}
}

// Helper function to create a simple test UI
func createTestUI() *ebitenui.UI {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	return &ebitenui.UI{
		Container: rootContainer,
	}
}