package main

import (
	"TowerDefense/internal/game"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := game.NewGame()
	ebiten.SetWindowTitle("Tower Defense")
	ebiten.SetWindowSize(960, 960)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
