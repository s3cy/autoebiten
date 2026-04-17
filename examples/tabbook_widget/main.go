// Example demonstrating TabBook widget automation with autoui.
// This example shows how to use proxy methods to interact with TabBook widgets.
//
// Runnable example referenced in docs/autoui-reference.md.
//
// Build and run:
//
//	cd examples/tabbook_widget
//	go build -o tabbook_demo
//	autoebiten launch -- ./tabbook_demo &
//
// Then try the autoui commands:
//
//	autoebiten custom autoui.call --request '{"target":"type=TabBook","method":"Tabs"}'
//	autoebiten custom autoui.call --request '{"target":"type=TabBook","method":"SetTabByIndex","args":[1]}'
//	autoebiten custom autoui.call --request '{"target":"type=TabBook","method":"TabIndex"}'
package main

import (
	"image"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/s3cy/autoebiten"
	"github.com/s3cy/autoebiten/autoui"
)

type Game struct {
	ui *ebitenui.UI
}

func main() {
	// Create root container
	root := widget.NewContainer()
	root.GetWidget().Rect = image.Rect(0, 0, 640, 480)

	// Create TabBook
	tabBook := widget.NewTabBook()

	// Position the TabBook
	tabBook.GetWidget().Rect = image.Rect(50, 50, 590, 430)

	// Add TabBook to container
	root.AddChild(tabBook)

	// Create UI (minimal setup without theme)
	ui := &ebitenui.UI{Container: root}

	// Register autoui commands with autoebiten
	autoui.Register(ui)

	g := &Game{ui: ui}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("TabBook Widget Demo")

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
