// Package state_exporter provides a sample game with state exporting capabilities.
// It demonstrates how to use autoebiten's state query and custom command features.
package state_exporter

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/s3cy/autoebiten"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
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

// Enemy represents an enemy entity in the game.
type Enemy struct {
	Name   string
	Health int
	X, Y   float64
}

// Game implements ebiten.Game interface and exports its state.
type Game struct {
	State GameState
}

// NewGame creates a new game instance with initialized state.
func NewGame() *Game {
	g := &Game{
		State: GameState{
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
	autoebiten.RegisterStateExporter("gamestate", &g.State)

	// Register custom commands
	autoebiten.Register("heal", func(ctx autoebiten.CommandContext) {
		old := g.State.Player.Health
		g.State.Player.Health = min(g.State.Player.Health+20, 100)
		ctx.Respond(fmt.Sprintf("Healed from %d to %d", old, g.State.Player.Health))
	})

	return g
}

// Update handles game logic each tick.
func (g *Game) Update() error {
	if !autoebiten.Update() {
		return fmt.Errorf("exit requested")
	}

	// Movement
	speed := 2.0
	if autoebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.State.Player.X += speed
	}
	if autoebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.State.Player.X -= speed
	}
	if autoebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.State.Player.Y -= speed
	}
	if autoebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.State.Player.Y += speed
	}

	// Damage
	if autoebiten.IsKeyPressed(ebiten.KeyD) {
		g.State.Player.Health = max(g.State.Player.Health-5, 0)
	}

	return nil
}

// Draw renders the game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x10, 0x20, 0x40, 0xff})

	msg := "=== State Exporter Demo ===\n\n"
	msg += fmt.Sprintf("Player: (%.1f, %.1f)\n", g.State.Player.X, g.State.Player.Y)
	msg += fmt.Sprintf("Health: %d  Mana: %d\n", g.State.Player.Health, g.State.Player.Mana)
	msg += fmt.Sprintf("Score: %d\n", g.State.Score)
	msg += "\nEnemies:\n"
	for i, e := range g.State.Enemies {
		msg += fmt.Sprintf("  %d: %s (HP:%d)\n", i, e.Name, e.Health)
	}
	msg += "\nCLI Commands:\n"
	msg += "  state --name gamestate --path Player.Health\n"
	msg += "  state --name gamestate --path Enemies.0.Name\n"

	ebitenutil.DebugPrint(screen, msg)
	autoebiten.Capture(screen)
}

// Layout returns the game's logical screen size.
func (g *Game) Layout(_, _ int) (int, int) {
	return ScreenWidth, ScreenHeight
}