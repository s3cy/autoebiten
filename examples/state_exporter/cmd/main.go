package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/s3cy/autoebiten/examples/state_exporter"
)

func main() {
	ebiten.SetWindowSize(state_exporter.ScreenWidth, state_exporter.ScreenHeight)
	ebiten.SetWindowTitle("State Exporter Demo")

	g := state_exporter.NewGame()
	if err := ebiten.RunGameWithOptions(g, &ebiten.RunGameOptions{InitUnfocused: true}); err != nil {
		log.Fatal(err)
	}
}
