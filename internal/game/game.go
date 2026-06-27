package game

import (
	"TowerDefense/internal/renderer"
	. "TowerDefense/internal/vars"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Map      *GameMap
	Renderer *renderer.Renderer
}

// convertGridToRendererFormat 將 GameMap 的 Grid 轉換為 renderer 需要的格式
func convertGridToRendererFormat(gm *GameMap) [GridSize][GridSize]struct{ Type TerrainType } {
	var result [GridSize][GridSize]struct{ Type TerrainType }
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			result[y][x].Type = gm.Grid[y][x].Type
		}
	}
	return result
}

func NewGame() *Game {
	game := &Game{}
	game.Map = NewGameMap(GenerateMapScale) // Initialize the game map
	game.Renderer = renderer.NewRenderer()  // Initialize the renderer
	game.Renderer.DrawTerrainLayer(convertGridToRendererFormat(game.Map))
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
