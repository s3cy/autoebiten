// Example demonstrating autoui integration with autoebiten.
// This example shows how to use autoui to enable CLI-based widget inspection
// and interaction for E2E testing.
//
// Runnable examples in docs/autoui.md reference this game.
//
// Build and run:
//
//	cd examples/autoui
//	go build -o autoui_demo
//	autoebiten launch -- ./autoui_demo &
//
// Then try the autoui commands:
//
//	autoebiten custom autoui.tree
//	autoebiten custom autoui.find --request "type=Button"
//	autoebiten custom autoui.call --request '{"target":"id=submit-btn","method":"Click"}'
package main

import (
	"image"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	ebitenuiImage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/s3cy/autoebiten"
	"github.com/s3cy/autoebiten/autoui"
)

// PlayerCard demonstrates ae tag usage for custom XML attributes.
type PlayerCard struct {
	PlayerID   string `ae:"player_id"`
	PlayerName string `ae:"player_name"`
	Level      int    `ae:"level"`
	ID         string `ae:"id"`
	Role       string `ae:"role"`
}

type Game struct {
	ui *ebitenui.UI
}

func main() {
	// Create simple button images (colored rectangles)
	buttonImage := &widget.ButtonImage{
		Idle:     createTestNineSlice(100, 40, color.RGBA{80, 80, 120, 255}),
		Pressed:  createTestNineSlice(100, 40, color.RGBA{60, 60, 100, 255}),
		Disabled: createTestNineSlice(100, 40, color.RGBA{120, 120, 120, 255}),
	}

	// Create root container
	root := widget.NewContainer()
	root.GetWidget().Rect = image.Rect(0, 0, 640, 480)

	// Submit button - matches doc example
	submitBtn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CustomData(&PlayerCard{
				PlayerID:   "p001",
				PlayerName: "Alice",
				Level:      42,
				ID:         "submit-btn",
				Role:       "primary",
			}),
		),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			log.Println("Submit button clicked!")
		}),
	)
	submitBtn.GetWidget().Rect = image.Rect(100, 50, 300, 90)

	// Cancel button - matches doc example
	cancelBtn := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			log.Println("Cancel button clicked!")
		}),
	)
	cancelBtn.GetWidget().Rect = image.Rect(100, 200, 300, 240)
	cancelBtn.GetWidget().CustomData = map[string]string{"id": "cancel-btn", "role": "secondary"}

	root.AddChild(submitBtn)
	root.AddChild(cancelBtn)

	// Create UI (minimal setup without theme)
	ui := &ebitenui.UI{Container: root}

	// Register autoui commands with autoebiten
	// Enables: autoui.tree, autoui.find, autoui.at, autoui.xpath, autoui.call, autoui.highlight
	autoui.Register(ui)

	g := &Game{ui: ui}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("autoui Demo")

	if err := ebiten.RunGameWithOptions(g, &ebiten.RunGameOptions{InitUnfocused: true}); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	// Required: Process RPC requests from CLI
	if !autoebiten.Update() {
		return ebiten.Termination
	}
	// Note: Skip g.ui.Update() to avoid validation that requires fonts
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 50, 255})
	g.ui.Draw(screen)

	// Required: Enable screenshot capture
	autoebiten.Capture(screen)

	// Optional: Draw highlights for visual debugging (on top of UI)
	autoui.DrawHighlights(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// createTestNineSlice creates a simple colored nine-slice for button rendering.
func createTestNineSlice(w, h int, c color.Color) *ebitenuiImage.NineSlice {
	img := ebiten.NewImage(w, h)
	img.Fill(c)
	return ebitenuiImage.NewNineSliceSimple(img, 0, 0)
}
