package game

import (
	"TowerDefense/internal/gameMap"
	"TowerDefense/internal/renderer"
	. "TowerDefense/internal/vars"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Map      *gameMap.GameMap
	Renderer *renderer.Renderer
}

func NewGame() *Game {
	game := &Game{}
	game.Map = gameMap.NewGameMap(GenerateMapScale) // Initialize the game map
	game.Renderer = renderer.NewRenderer()          // Initialize the renderer
	game.Renderer.DrawTerrainLayer(game.Map)
	game.Renderer.DrawCastleLayer()
	return game
}
func (g *Game) Update() error {
	g.Renderer.Camera.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Renderer.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
