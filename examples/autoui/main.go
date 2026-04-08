// Example demonstrating autoui integration with autoebiten.
// This example shows how to use autoui to enable CLI-based widget inspection
// and interaction for E2E testing.
package main

import (
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
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout()),
	)

	// Create buttons
	startBtn := widget.NewButton(
		widget.ButtonOpts.Text("Start Game", nil, nil),
	)
	// Set custom data for identification (must be set after creation)
	startBtn.GetWidget().CustomData = map[string]string{
		"id":   "start-btn",
		"role": "primary",
	}

	quitBtn := widget.NewButton(
		widget.ButtonOpts.Text("Quit", nil, nil),
	)
	quitBtn.GetWidget().CustomData = map[string]string{
		"id":   "quit-btn",
		"role": "secondary",
	}

	root.AddChild(startBtn)
	root.AddChild(quitBtn)

	// Create UI
	ui := &ebitenui.UI{Container: root}

	// Register autoui commands with autoebiten
	// This enables CLI commands: autoui.tree, autoui.find, autoui.at, etc.
	autoui.Register(ui)

	game := &Game{ui: ui}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("autoui Demo")

	// Run the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	// Required: Call autoebiten.Update() to process RPC requests
	if !autoebiten.Update() {
		return ebiten.Termination
	}
	g.ui.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)

	// Required: Capture screen for screenshot requests
	autoebiten.Capture(screen)

	// Optional: Draw any active highlights (for visual debugging)
	// Call this after UI.Draw to render highlights on top
	autoui.DrawHighlights(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
