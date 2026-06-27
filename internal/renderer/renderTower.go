package renderer

import (
	"image"
	"image/color"

	. "TowerDefense/internal/vars"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var TowerTileRect = [TowerTotal]image.Rectangle{
	TowerNone:   image.Rect(0*TileSize, 0*TileSize, (0+1)*TileSize, (0+1)*TileSize),
	TowerBasic:  image.Rect(1*TileSize, 10*TileSize, (1+1)*TileSize, (10+1)*TileSize),
	TowerSniper: image.Rect(1*TileSize, 11*TileSize, (1+1)*TileSize, (11+1)*TileSize),
	TowerRapid:  image.Rect(2*TileSize, 10*TileSize, (2+1)*TileSize, (10+1)*TileSize),
}

// TowerData 塔防資料結構 (用於繪製)
type TowerDrawData struct {
	X, Y   int
	Typ    int
	Range  float32
	Active bool
}
type TowerDataArr struct {
	Towers [MaxTowers]TowerDrawData
	Count  int
}

func (r *Renderer) DrawTowers(screen *ebiten.Image, tda *TowerDataArr) {
	screenW := float64(screen.Bounds().Dx())
	screenH := float64(screen.Bounds().Dy())

	var camMatrix ebiten.GeoM
	camMatrix.Translate(-r.Camera.X, -r.Camera.Y)
	camMatrix.Scale(r.Camera.Zoom, r.Camera.Zoom)
	camMatrix.Translate(screenW/2, screenH/2)

	for i := 0; i < tda.Count; i++ {
		t := tda.Towers[i] // 塔防資料
		if !t.Active {
			continue
		}
		// 世界座標 (以網格中心為基準)
		wx, wy := float64(t.X*TileSize+TileSize/2), float64(t.Y*TileSize+TileSize/2)
		sx, sy := camMatrix.Apply(wx, wy)

		// 繪製塔
		radius := float32(TileSize/2) * float32(r.Camera.Zoom)
		var c color.Color
		switch t.Typ {
		case 1: // TowerBasic
			c = color.RGBA{0, 0, 255, 255} // 基礎塔是藍色
		case 2: // TowerSniper
			c = color.RGBA{128, 0, 128, 255} // 狙擊塔是紫色
		case 3: // TowerRapid
			c = color.RGBA{0, 255, 255, 255} // 速射塔是青色
		default:
			c = color.White
		}
		vector.FillCircle(screen, float32(sx), float32(sy), radius, c, true)

		// 繪製攻擊範圍 (半透明)
		rangeRadius := t.Range * float32(TileSize) * float32(r.Camera.Zoom)
		vector.StrokeCircle(screen, float32(sx), float32(sy), rangeRadius, 1, color.RGBA{255, 255, 255, 50}, true)
	}
}
