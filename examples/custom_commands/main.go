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

// Game demonstrates custom commands.
type Game struct {
	playerHealth int
	playerMana   int
	lastMessage  string
	deferredCtx  autoebiten.CommandContext
	deferredTick int
}

// NewGame creates a new game instance.
func NewGame() *Game {
	g := &Game{
		playerHealth: 100,
		playerMana:   50,
		lastMessage:  "Waiting for commands...",
	}

	// Register custom commands
	autoebiten.Register("getPlayerInfo", func(ctx autoebiten.CommandContext) {
		info := fmt.Sprintf("Health: %d, Mana: %d", g.playerHealth, g.playerMana)
		ctx.Respond(info)
		g.lastMessage = "Sent player info"
	})

	autoebiten.Register("heal", func(ctx autoebiten.CommandContext) {
		oldHealth := g.playerHealth
		g.playerHealth = min(g.playerHealth+20, 100)
		ctx.Respond(fmt.Sprintf("Healed from %d to %d", oldHealth, g.playerHealth))
		g.lastMessage = "Player healed"
	})

	autoebiten.Register("damage", func(ctx autoebiten.CommandContext) {
		oldHealth := g.playerHealth
		g.playerHealth = max(g.playerHealth-10, 0)
		ctx.Respond(fmt.Sprintf("Damaged from %d to %d", oldHealth, g.playerHealth))
		g.lastMessage = "Player damaged"
	})

	autoebiten.Register("echo", func(ctx autoebiten.CommandContext) {
		// Echo back the request
		ctx.Respond(fmt.Sprintf("Echo: %s", ctx.Request()))
		g.lastMessage = "Echo command received"
	})

	autoebiten.Register("deferred", func(ctx autoebiten.CommandContext) {
		// Store the context and respond later
		g.lastMessage = "Deferred command will respond in 60 ticks..."
		g.deferredCtx = ctx
		g.deferredTick = 60
	})

	return g
}

// Update is called every frame.
func (g *Game) Update() error {
	if !autoebiten.Update() {
		return fmt.Errorf("exit requested")
	}

	// Handle deferred response
	if g.deferredTick > 0 {
		g.deferredTick--
		if g.deferredTick == 0 && g.deferredCtx != nil {
			g.deferredCtx.Respond("Deferred response after 60 ticks!")
			g.deferredCtx = nil
			g.lastMessage = "Deferred response sent"
		}
	}

	// Check input
	if autoebiten.IsKeyPressed(ebiten.KeyH) {
		g.playerHealth = min(g.playerHealth+5, 100)
	}
	if autoebiten.IsKeyPressed(ebiten.KeyD) {
		g.playerHealth = max(g.playerHealth-5, 0)
	}

	return nil
}

// Draw is called every frame.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x00, 0x00, 0x66, 0xff})

	// Draw UI
	msg := "=== Custom Commands Demo ===\n\n"
	msg += fmt.Sprintf("Health: %d | Mana: %d\n\n", g.playerHealth, g.playerMana)
	msg += "Available Commands:\n"
	msg += "  getPlayerInfo - Get player stats\n"
	msg += "  heal          - Heal player (+20)\n"
	msg += "  damage        - Damage player (-10)\n"
	msg += "  echo <text>   - Echo back text\n"
	msg += "  deferred      - Response after 60 ticks\n\n"
	msg += fmt.Sprintf("Last: %s", g.lastMessage)

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
	ebiten.SetWindowTitle("autoebiten Custom Commands Demo")

	g := NewGame()

	if err := ebiten.RunGameWithOptions(g, &ebiten.RunGameOptions{InitUnfocused: true}); err != nil {
		log.Fatal("Failed to run game:", err)
	}
}
