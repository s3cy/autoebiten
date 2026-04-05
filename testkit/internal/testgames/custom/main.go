package main

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
	"sync/atomic"

	"github.com/s3cy/autoebiten"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

// CustomGame is a test game with custom commands.
type CustomGame struct {
	counter atomic.Int32
}

// Update is called every frame.
func (g *CustomGame) Update() error {
	if !autoebiten.Update() {
		return fmt.Errorf("exit requested")
	}
	return nil
}

// Draw is called every frame.
func (g *CustomGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x00, 0x66, 0x00, 0xFF})
	autoebiten.Capture(screen)
}

// Layout is called to get the screen size.
func (g *CustomGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// RegisterCommands registers custom commands.
func (g *CustomGame) RegisterCommands() {
	// Echo command - returns the request back
	autoebiten.Register("echo", func(ctx autoebiten.CommandContext) {
		ctx.Respond(ctx.Request())
	})

	// Counter command - increments and returns count
	autoebiten.Register("counter", func(ctx autoebiten.CommandContext) {
		newVal := g.counter.Add(1)
		ctx.Respond(strconv.Itoa(int(newVal)))
	})

	// GetCounter command - returns current count without incrementing
	autoebiten.Register("getCounter", func(ctx autoebiten.CommandContext) {
		val := g.counter.Load()
		ctx.Respond(strconv.Itoa(int(val)))
	})

	// ResetCounter command - resets count to 0
	autoebiten.Register("resetCounter", func(ctx autoebiten.CommandContext) {
		g.counter.Store(0)
		ctx.Respond("0")
	})
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("testkit custom command test game")

	g := &CustomGame{}
	g.RegisterCommands()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
