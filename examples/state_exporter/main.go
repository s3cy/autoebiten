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

// GameState holds all game data for state queries.
type GameState struct {
	Player struct {
		X      float64
		Y      float64
		Health int
		Mana   int
	}
	Enemies []Enemy
	Score   int
}

type Enemy struct {
	Name   string
	Health int
	X, Y   float64
}

// Game implements ebiten.Game interface.
type Game struct {
	state GameState
}

func NewGame() *Game {
	g := &Game{
		state: GameState{
			Player: struct {
				X      float64
				Y      float64
				Health int
				Mana   int
			}{X: 100, Y: 100, Health: 100, Mana: 50},
			Enemies: []Enemy{
				{Name: "Goblin", Health: 30, X: 300, Y: 200},
				{Name: "Orc", Health: 50, X: 400, Y: 300},
			},
			Score: 0,
		},
	}

	// Register state exporter for StateQuery
	autoebiten.RegisterStateExporter("gamestate", &g.state)

	// Register custom commands
	autoebiten.Register("heal", func(ctx autoebiten.CommandContext) {
		old := g.state.Player.Health
		g.state.Player.Health = min(g.state.Player.Health+20, 100)
		ctx.Respond(fmt.Sprintf("Healed from %d to %d", old, g.state.Player.Health))
	})

	return g
}

func (g *Game) Update() error {
	if !autoebiten.Update() {
		return fmt.Errorf("exit requested")
	}

	// Movement
	speed := 2.0
	if autoebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.state.Player.X += speed
	}
	if autoebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.state.Player.X -= speed
	}
	if autoebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.state.Player.Y -= speed
	}
	if autoebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.state.Player.Y += speed
	}

	// Damage
	if autoebiten.IsKeyPressed(ebiten.KeyD) {
		g.state.Player.Health = max(g.state.Player.Health-5, 0)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x10, 0x20, 0x40, 0xff})

	msg := "=== State Exporter Demo ===\n\n"
	msg += fmt.Sprintf("Player: (%.1f, %.1f)\n", g.state.Player.X, g.state.Player.Y)
	msg += fmt.Sprintf("Health: %d  Mana: %d\n", g.state.Player.Health, g.state.Player.Mana)
	msg += fmt.Sprintf("Score: %d\n", g.state.Score)
	msg += "\nEnemies:\n"
	for i, e := range g.state.Enemies {
		msg += fmt.Sprintf("  %d: %s (HP:%d)\n", i, e.Name, e.Health)
	}
	msg += "\nCLI Commands:\n"
	msg += "  state --name gamestate --path Player.Health\n"
	msg += "  state --name gamestate --path Enemies.0.Name\n"

	ebitenutil.DebugPrint(screen, msg)
	autoebiten.Capture(screen)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("State Exporter Demo")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}