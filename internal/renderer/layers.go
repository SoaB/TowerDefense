package renderer

import (
	"image"
	"image/color"

	. "TowerDefense/internal/vars"

	"TowerDefense/internal/gameMap"

	"github.com/hajimehoshi/ebiten/v2"
)

var TerrainTileRect = [TerrainTotal]image.Rectangle{
	TerrainGrass:    image.Rect(2*TileSize, 1*TileSize, (2+1)*TileSize, (1+1)*TileSize),
	TerrainHill:     image.Rect(4*TileSize, 8*TileSize, (4+1)*TileSize, (8+1)*TileSize),
	TerrainSand:     image.Rect(3*TileSize, 0*TileSize, (3+1)*TileSize, (0+1)*TileSize),
	TerrainMud:      image.Rect(10*TileSize, 1*TileSize, (10+1)*TileSize, (1+1)*TileSize),
	TerrainStone:    image.Rect(4*TileSize, 9*TileSize, (4+1)*TileSize, (9+1)*TileSize),
	TerrainMountain: image.Rect(4*TileSize, 11*TileSize, (4+1)*TileSize, (11+1)*TileSize),
	TerrainLake:     image.Rect(10*TileSize, 2*TileSize, (10+1)*TileSize, (2+1)*TileSize),
	TerrainCastle:   image.Rect(2*TileSize, 3*TileSize, (2+1)*TileSize, (3+1)*TileSize),
}

func (r *Renderer) DrawTerrainLayer(gm *gameMap.GameMap) {
	r.ResMan.TerrainLayer.Fill(color.Transparent)

	op := &ebiten.DrawImageOptions{}

	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			tileType := gm.Grid[y][x].Type
			if tileType == TerrainCastle { // TerrainCastle
				continue
			}
			dX, dY := float64(x*TileSize), float64(y*TileSize)
			op.GeoM.Reset()
			op.GeoM.Translate(dX, dY)
			tr := TerrainTileRect[tileType]
			r.ResMan.TerrainLayer.DrawImage(r.ResMan.TileImage.SubImage(tr).(*ebiten.Image), op)
		}
	}
}

func (r *Renderer) DrawCastleLayer() {
	r.ResMan.CastleLayer.Fill(color.Transparent)
	cx, cy := float64(15), float64(15) // 30/2 = 15
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(cx*TileSize, cy*TileSize)
	tr := TerrainTileRect[TerrainCastle]
	r.ResMan.CastleLayer.DrawImage(r.ResMan.TileImage.SubImage(tr).(*ebiten.Image), op) // TerrainCastle
}
