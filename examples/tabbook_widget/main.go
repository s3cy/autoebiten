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
	"bytes"
	_ "embed"
	"image"
	"image/color"
	"log"

	ebitenuiImage "github.com/ebitenui/ebitenui/image"
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

	// Create button images for tabs
	buttonImage := createButtonImage()

	// Create tabs
	tab1 := widget.NewTabBookTab(
		widget.TabBookTabOpts.Label("Tab 1"),
		widget.TabBookTabOpts.ContainerOpts(
			widget.ContainerOpts.BackgroundImage(ebitenuiImage.NewNineSliceColor(color.NRGBA{100, 100, 150, 255})),
		),
	)
	tab2 := widget.NewTabBookTab(
		widget.TabBookTabOpts.Label("Tab 2"),
		widget.TabBookTabOpts.ContainerOpts(
			widget.ContainerOpts.BackgroundImage(ebitenuiImage.NewNineSliceColor(color.NRGBA{100, 150, 100, 255})),
		),
	)

	// Create TabBook
	tabBook := widget.NewTabBook(
		widget.TabBookOpts.TabButtonImage(buttonImage),
		widget.TabBookOpts.TabButtonText(&face, &widget.ButtonTextColor{Idle: color.White, Disabled: color.NRGBA{128, 128, 128, 255}}),
		widget.TabBookOpts.Tabs(tab1, tab2),
		widget.TabBookOpts.InitialTab(tab1),
	)

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
	// Note: Skip g.ui.Update() to avoid validation that requires full theme
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 50, 255})
	// Skip ui.Draw() - TabBook requires full theme setup for rendering
	// The autoui commands work at the widget tree level without rendering
	autoebiten.Capture(screen)
	autoui.DrawHighlights(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func loadFont(size float64) (text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}

func createButtonImage() *widget.ButtonImage {
	return &widget.ButtonImage{
		Idle:     ebitenuiImage.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255}),
		Hover:    ebitenuiImage.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255}),
		Pressed:  ebitenuiImage.NewNineSliceColor(color.NRGBA{R: 255, G: 100, B: 120, A: 255}),
		Disabled: ebitenuiImage.NewNineSliceColor(color.NRGBA{R: 120, G: 120, B: 120, A: 255}),
	}
}
