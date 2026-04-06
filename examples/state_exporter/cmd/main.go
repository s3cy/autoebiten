package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/s3cy/autoebiten/examples/state_exporter"
)

func main() {
	ebiten.SetWindowSize(state_exporter.ScreenWidth, state_exporter.ScreenHeight)
	ebiten.SetWindowTitle("State Exporter Demo")

	if err := ebiten.RunGame(state_exporter.NewGame()); err != nil {
		log.Fatal(err)
	}
}