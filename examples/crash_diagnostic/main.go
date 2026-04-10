package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/s3cy/autoebiten"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

// Game is a crash diagnostic demo game.
type Game struct {
	crashAfterRPC bool
	counter       int
}

// Update is called every frame.
func (g *Game) Update() error {
	if !autoebiten.Update() {
		return fmt.Errorf("exit requested")
	}

	// Simulate some work and log output
	g.counter++
	if g.counter%60 == 0 {
		fmt.Printf("Game running... tick %d\n", g.counter)
	}

	// Check if we should crash after RPC connection
	if g.crashAfterRPC && g.counter >= 180 { // Crash after ~3 seconds
		panic("intentional crash after RPC connection")
	}

	return nil
}

// Draw is called every frame.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x00, 0x00, 0x66, 0xff})

	msg := "Crash Diagnostic Demo\n"
	msg += "=====================\n"
	if g.crashAfterRPC {
		msg += "Mode: Will crash AFTER RPC connection\n"
		msg += fmt.Sprintf("Counter: %d (crashes at 180)\n", g.counter)
	} else {
		msg += "Mode: Normal operation (no crash)\n"
	}
	msg += "\nUse CLI to interact with this game:\n"
	msg += "  autoebiten ping\n"
	msg += "  autoebiten screenshot\n"
	msg += "  autoebiten exit"

	ebitenutil.DebugPrint(screen, msg)

	// Capture screenshot for CLI requests
	autoebiten.Capture(screen)
}

// Layout is called to get the screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Parse command-line arguments
	crashBeforeRPC := flag.Bool("crash-before-rpc", false, "Crash before RPC connection (in main)")
	crashAfterRPC := flag.Bool("crash-after-rpc", false, "Crash after RPC connection (in Update)")
	flag.Parse()

	fmt.Println("Starting crash diagnostic demo...")
	fmt.Printf("Flags: crash-before-rpc=%v, crash-after-rpc=%v\n", *crashBeforeRPC, *crashAfterRPC)

	// Simulate some initialization time
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Initialization complete")

	// Crash before RPC connection if requested
	if *crashBeforeRPC {
		fmt.Println("About to crash before RPC connection!")
		time.Sleep(100 * time.Millisecond)
		panic("intentional crash before RPC connection")
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("autoebiten Crash Diagnostic Demo")

	g := &Game{
		crashAfterRPC: *crashAfterRPC,
	}

	if err := ebiten.RunGameWithOptions(g, &ebiten.RunGameOptions{InitUnfocused: true}); err != nil {
		log.Fatal("Failed to run game:", err)
	}
}
