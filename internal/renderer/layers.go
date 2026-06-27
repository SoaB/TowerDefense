package renderer

import (
	"image"
	"image/color"

	. "TowerDefense/internal/vars"

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

type TerrainDrawData struct {
	TerrainType TerrainType
}

type TerrainDataArr struct {
	Grid [GridSize][GridSize]TerrainDrawData
}

func (r *Renderer) DrawTerrainLayer(tda *TerrainDataArr) {
	r.ResMan.TerrainLayer.Fill(color.Transparent)

	op := &ebiten.DrawImageOptions{}

	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			tileType := tda.Grid[y][x].TerrainType
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

func (r *Renderer) DrawLayer(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screenW := float64(screen.Bounds().Dx())
	screenH := float64(screen.Bounds().Dy())
	op.GeoM.Reset()
	op.GeoM.Translate(-r.Camera.X, -r.Camera.Y)
	op.GeoM.Scale(r.Camera.Zoom, r.Camera.Zoom)
	op.GeoM.Translate(screenW/2, screenH/2)
	screen.DrawImage(r.ResMan.TerrainLayer, op)
	screen.DrawImage(r.ResMan.CastleLayer, op)
}
