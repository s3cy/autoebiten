package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/s3cy/autoebiten"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

// Game is a simple demo game.
type Game struct{}

// Update is called every frame.
func (g *Game) Update() error {
	if !autoebiten.Update() {
		return fmt.Errorf("exit requested")
	}

	// Example: check if a key is pressed (real or injected)
	if autoebiten.IsKeyPressed(ebiten.KeySpace) {
		fmt.Println("Space is pressed!")
	}

	return nil
}

// Draw is called every frame.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x00, 0x00, 0x66, 0xff})

	// Draw some text
	msg := "autoebiten demo game\n"
	msg += "Press keys to see them registered\n"
	msg += "Use CLI to inject inputs"

	ebitenutil.DebugPrint(screen, msg)

	// Capture screenshot for CLI requests
	autoebiten.Capture(screen)
}

// Layout is called to get the screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("autoebiten Demo")

	g := &Game{}

	if err := ebiten.RunGameWithOptions(g, &ebiten.RunGameOptions{InitUnfocused: true}); err != nil {
		log.Fatal("Failed to run game:", err)
	}
}
