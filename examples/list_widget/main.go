// Example demonstrating List widget automation with autoui.
// This example shows how to use proxy methods to interact with List widgets.
//
// Runnable example referenced in docs/autoui-reference.md.
//
// Build and run:
//
//	cd examples/list_widget
//	go build -o list_demo
//	autoebiten launch -- ./list_demo &
//
// Then try the autoui commands:
//
//	autoebiten custom autoui.call --request '{"target":"type=List","method":"Entries"}'
//	autoebiten custom autoui.call --request '{"target":"type=List","method":"SelectEntryByIndex","args":[1]}'
//	autoebiten custom autoui.call --request '{"target":"type=List","method":"SelectedEntryIndex"}'
package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/s3cy/autoebiten"
	"github.com/s3cy/autoebiten/autoui"
)

//go:embed assets/fonts/notosans-regular.ttf
var fontData []byte

type Game struct {
	ui *ebitenui.UI
}

func main() {
	// Load font
	face, err := loadFont(20)
	if err != nil {
		log.Fatal(err)
	}

	// Create root container
	root := widget.NewContainer()
	root.GetWidget().Rect = image.Rect(0, 0, 640, 480)

	// Create a List widget with sample entries
	entries := []any{"Option A", "Option B", "Option C", "Option D", "Option E"}
	list := widget.NewList(
		widget.ListOpts.Entries(entries),
		widget.ListOpts.EntryLabelFunc(func(e any) string {
			return e.(string)
		}),
		widget.ListOpts.EntryFontFace(face),
	)
	list.GetWidget().Rect = image.Rect(100, 50, 300, 250)
	list.GetWidget().CustomData = map[string]string{"id": "main-list"}

	// Add to container
	root.AddChild(list)

	// Create UI
	ui := &ebitenui.UI{Container: root}

	// Register autoui
	autoui.Register(ui)

	g := &Game{ui: ui}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("List Widget Demo")

	if err := ebiten.RunGameWithOptions(g, &ebiten.RunGameOptions{InitUnfocused: true}); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	if !autoebiten.Update() {
		return ebiten.Termination
	}
	g.ui.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 50, 255})
	g.ui.Draw(screen)
	autoebiten.Capture(screen)
	autoui.DrawHighlights(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func loadFont(size float64) (*text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		return nil, err
	}

	var face text.Face = &text.GoTextFace{
		Source: s,
		Size:   size,
	}
	return &face, nil
}
