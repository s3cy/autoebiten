package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/s3cy/autoebiten"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

// SimpleGame is a minimal test game that renders a colored screen.
type SimpleGame struct{}

// Update is called every frame.
func (g *SimpleGame) Update() error {
	if !autoebiten.Update() {
		return fmt.Errorf("exit requested")
	}
	return nil
}

// Draw is called every frame.
func (g *SimpleGame) Draw(screen *ebiten.Image) {
	// Fill with a distinctive color for screenshot verification
	screen.Fill(color.RGBA{0x42, 0x86, 0xF4, 0xFF}) // Google blue
	autoebiten.Capture(screen)
}

// Layout is called to get the screen size.
func (g *SimpleGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("testkit simple test game")

	g := &SimpleGame{}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
