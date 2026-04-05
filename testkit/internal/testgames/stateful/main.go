package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/s3cy/autoebiten"
	"github.com/s3cy/autoebiten/testkit"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 240
	moveSpeed    = 2
)

// Player represents the player character.
type Player struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Health int     `json:"health"`
}

// InventoryItem represents an item in the inventory.
type InventoryItem struct {
	Name string `json:"name"`
	Qty  int    `json:"qty"`
}

// StatefulGame is a test game with player state and movement.
type StatefulGame struct {
	Player    Player          `json:"player"`
	Inventory []InventoryItem `json:"inventory"`
	Skills    map[string]int  `json:"skills"`
	TickCount int64           `json:"tickCount"`
}

// NewStatefulGame creates a new stateful game with initial state.
func NewStatefulGame() *StatefulGame {
	return &StatefulGame{
		Player: Player{
			X:      100,
			Y:      100,
			Health: 100,
		},
		Inventory: []InventoryItem{
			{Name: "Sword", Qty: 1},
			{Name: "Shield", Qty: 1},
		},
		Skills: map[string]int{
			"Sword":   10,
			"Shield":  5,
			"Magic":   3,
		},
		TickCount: 0,
	}
}

// Update is called every frame.
func (g *StatefulGame) Update() error {
	if !autoebiten.Update() {
		return fmt.Errorf("exit requested")
	}

	g.TickCount++

	// Movement on key press
	if autoebiten.IsKeyPressed(ebiten.KeyArrowRight) || autoebiten.IsKeyPressed(ebiten.KeyD) {
		g.Player.X += moveSpeed
	}
	if autoebiten.IsKeyPressed(ebiten.KeyArrowLeft) || autoebiten.IsKeyPressed(ebiten.KeyA) {
		g.Player.X -= moveSpeed
	}
	if autoebiten.IsKeyPressed(ebiten.KeyArrowUp) || autoebiten.IsKeyPressed(ebiten.KeyW) {
		g.Player.Y -= moveSpeed
	}
	if autoebiten.IsKeyPressed(ebiten.KeyArrowDown) || autoebiten.IsKeyPressed(ebiten.KeyS) {
		g.Player.Y += moveSpeed
	}

	return nil
}

// Draw is called every frame.
func (g *StatefulGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x33, 0x33, 0x33, 0xFF})
	autoebiten.Capture(screen)
}

// Layout is called to get the screen size.
func (g *StatefulGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// RegisterCommands registers custom commands.
func (g *StatefulGame) RegisterCommands() {
	autoebiten.Register("testkit.state", testkit.StateExporter(g))
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("testkit stateful test game")

	g := NewStatefulGame()
	g.RegisterCommands()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
