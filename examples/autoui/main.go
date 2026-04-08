// Example demonstrating autoui integration with autoebiten.
//
// This example shows how to:
//   1. Create a UI with identifiable widgets using CustomData
//   2. Register autoui commands for CLI access
//   3. Use autoebiten for RPC support
//   4. Draw visual highlights from autoui.highlight command
//
// Run this example and use the CLI to inspect widgets:
//
//	# View the complete widget tree
//	autoebiten custom autoui.tree
//
//	# Find a widget by its ID
//	autoebiten custom autoui.find --request "id=start-btn"
//
//	# Find widgets using XPath expressions
//	autoebiten custom autoui.xpath --request "//Button[@visible='true']"
//
//	# Highlight all buttons (visual debugging)
//	autoebiten custom autoui.highlight --request "type=Button"
//
//	# Click a button programmatically
//	autoebiten custom autoui.call --request '{"target":"id=start-btn","method":"Click","args":[]}'
//
//	# Focus a text input
//	autoebiten custom autoui.call --request '{"target":"id=name-input","method":"Focus","args":[true]}'
//
// Build and run:
//
//	go build -o autoui-example ./examples/autoui/ && ./autoui-example
//
// Then use the autoebiten CLI to interact with the running game.
package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	ebitenui "github.com/ebitenui/ebitenui"
	ebitenuiImage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/s3cy/autoebiten"
	"github.com/s3cy/autoebiten/autoui"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

// customData is a struct for widget identification.
// The XML tags control how attributes appear in autoui.tree output.
type customData struct {
	ID   string `xml:"id,attr"`
	Role string `xml:"role,attr"`
}

// Game demonstrates autoui integration with autoebiten.
type Game struct {
	ui         *ebitenui.UI
	gameState  string
	clickCount int
	inputText  string
}

// NewGame creates a new game instance with UI setup.
func NewGame() *Game {
	g := &Game{
		gameState:  "idle",
		clickCount: 0,
		inputText:  "",
	}

	// Create the root container
	// CustomData is added directly to the container widget
	root := widget.NewContainer()
	root.GetWidget().Rect = image.Rect(0, 0, screenWidth, screenHeight)
	root.GetWidget().CustomData = customData{
		ID:   "root-container",
		Role: "main",
	}

	// Create button images for styling
	buttonImage := createButtonImage()

	// Button text color
	buttonColor := &widget.ButtonTextColor{
		Idle:     color.White,
		Disabled: color.Gray{128},
	}

	// Create Start button with CustomData ID
	// Note: CustomData must be passed through ButtonOpts.WidgetOpts wrapper
	startBtn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Start Game", nil, buttonColor),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(customData{
				ID:   "start-btn",
				Role: "primary",
			}),
		),
	)
	startBtn.GetWidget().Rect = image.Rect(50, 50, 200, 90)
	root.AddChild(startBtn)

	// Create Stop button
	stopBtn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Stop Game", nil, buttonColor),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(customData{
				ID:   "stop-btn",
				Role: "secondary",
			}),
		),
	)
	stopBtn.GetWidget().Rect = image.Rect(50, 100, 200, 140)
	root.AddChild(stopBtn)

	// Create Reset button
	resetBtn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Reset", nil, buttonColor),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(customData{
				ID:   "reset-btn",
				Role: "tertiary",
			}),
		),
	)
	resetBtn.GetWidget().Rect = image.Rect(50, 150, 200, 190)
	root.AddChild(resetBtn)

	// Create TextInput with CustomData ID
	// Note: CustomData must be passed through TextInputOpts.WidgetOpts wrapper
	textInput := widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(300, 30),
			widget.WidgetOpts.CustomData(customData{
				ID: "name-input",
			}),
		),
	)
	textInput.GetWidget().Rect = image.Rect(50, 200, 350, 230)
	root.AddChild(textInput)

	// Create the UI
	ui := &ebitenui.UI{
		Container: root,
	}

	// Register autoui commands for CLI access
	// This enables commands: autoui.tree, autoui.at, autoui.find, autoui.xpath,
	// autoui.call, autoui.highlight
	autoui.Register(ui)

	g.ui = ui
	return g
}

// Update is called every frame.
func (g *Game) Update() error {
	// autoebiten.Update() handles input injection and game loop control
	// Note: If using the Ebiten patch integration, this call would panic
	// (the patch handles updates automatically)
	if !autoebiten.Update() {
		return fmt.Errorf("exit requested")
	}

	// Update the UI
	g.ui.Update()

	return nil
}

// Draw is called every frame.
func (g *Game) Draw(screen *ebiten.Image) {
	// Fill background
	screen.Fill(color.RGBA{0x20, 0x20, 0x40, 0xff})

	// Draw the UI
	g.ui.Draw(screen)

	// Draw any active highlights from autoui.highlight command
	// This renders red rectangles around widgets that were highlighted via CLI
	autoui.DrawHighlights(screen)

	// Draw status info
	msg := fmt.Sprintf("State: %s | Clicks: %d | Input: %s\n\nCLI Commands:\n  autoui.tree\n  autoui.find id=start-btn\n  autoui.highlight type=Button\n  autoui.call target=id=start-btn method=Click",
		g.gameState, g.clickCount, g.inputText)
	ebitenutil.DebugPrintAt(screen, msg, 400, 50)

	// Capture screenshot for CLI requests
	// Note: If using the Ebiten patch integration, this call would panic
	// (the patch handles captures automatically)
	autoebiten.Capture(screen)
}

// Layout returns the screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// createButtonImage creates styled button images for the UI.
func createButtonImage() *widget.ButtonImage {
	return &widget.ButtonImage{
		Idle:     createNineSlice(150, 40, color.RGBA{80, 80, 120, 255}),
		Pressed:  createNineSlice(150, 40, color.RGBA{60, 60, 100, 255}),
		Disabled: createNineSlice(150, 40, color.RGBA{120, 120, 120, 255}),
	}
}

// createNineSlice creates a simple NineSlice for button styling.
func createNineSlice(w, h int, c color.Color) *ebitenuiImage.NineSlice {
	img := ebiten.NewImage(w, h)
	img.Fill(c)
	return ebitenuiImage.NewNineSliceSimple(img, 3, 3)
}

func main() {
	// Set window properties
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("autoui Demo - autoebiten + ebitenui")

	// Create game
	game := NewGame()

	// Run the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal("Failed to run game:", err)
	}
}