// Example demonstrating RadioGroup widget automation with autoui.
// This example shows how to use proxy methods to interact with RadioGroup widgets.
//
// Runnable example referenced in docs/autoui-reference.md.
//
// Build and run:
//
//	cd examples/radiogroup_widget
//	make run
//
// Then try the autoui commands:
//
//	autoebiten custom autoui.call --request '{"target":"radiogroup=settings-group","method":"Elements"}'
//	autoebiten custom autoui.call --request '{"target":"radiogroup=settings-group","method":"SetActiveByIndex","args":[1]}'
//	autoebiten custom autoui.call --request '{"target":"radiogroup=settings-group","method":"ActiveLabel"}'
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

	// Create toggle buttons for radio group
	btn1 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn2 := widget.NewButton(widget.ButtonOpts.ToggleMode())
	btn3 := widget.NewButton(widget.ButtonOpts.ToggleMode())

	// Position buttons
	btn1.GetWidget().Rect = image.Rect(100, 50, 200, 90)
	btn2.GetWidget().Rect = image.Rect(100, 100, 200, 140)
	btn3.GetWidget().Rect = image.Rect(100, 150, 200, 190)

	// Create RadioGroup with the buttons
	radioGroup := widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(btn1, btn2, btn3),
	)

	// Register the RadioGroup with autoui (required for RadioGroup access)
	autoui.RegisterRadioGroup("settings-group", radioGroup)

	// Add buttons to container
	root.AddChild(btn1)
	root.AddChild(btn2)
	root.AddChild(btn3)

	// Create UI (minimal setup without theme)
	ui := &ebitenui.UI{Container: root}

	// Register autoui commands with autoebiten
	autoui.Register(ui)
	// Register radio group again after autoui.Register to ensure it's in the registry
	autoui.RegisterRadioGroup("settings-group", radioGroup)

	g := &Game{ui: ui}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("RadioGroup Widget Demo")

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
