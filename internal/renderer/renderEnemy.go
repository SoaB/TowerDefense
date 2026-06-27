package renderer

import (
	"image"
	"image/color"

	. "TowerDefense/internal/vars"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var EnemyTileRect = [EnemyTotal]image.Rectangle{
	EnemyNone:   image.Rect(0*TileSize, 0*TileSize, (0+1)*TileSize, (0+1)*TileSize),
	EnemyMinion: image.Rect(0*TileSize, 7*TileSize, (0+1)*TileSize, (7+1)*TileSize),
	EnemyBoss:   image.Rect(1*TileSize, 0*TileSize, (1+1)*TileSize, (0+1)*TileSize),
}

// EnemyData 敵人資料結構 (用於繪製)
type EnemyData struct {
	X, Y      float64
	Typ       int
	Health    float32
	MaxHealth float32
	Status    int // 0=normal, 1=hurt
}
type EnemyDataArr struct {
	Enemies [MaxEnemies]EnemyData
	Count   int
}

func (r *Renderer) DrawEnemies(screen *ebiten.Image, eda *EnemyDataArr) {
	screenW := float64(screen.Bounds().Dx())
	screenH := float64(screen.Bounds().Dy())

	// 建立與 Draw 相同的轉換矩陣
	var camMatrix ebiten.GeoM
	camMatrix.Translate(-r.Camera.X, -r.Camera.Y)
	camMatrix.Scale(r.Camera.Zoom, r.Camera.Zoom)
	camMatrix.Translate(screenW/2, screenH/2)

	for i := 0; i < eda.Count; i++ {
		e := eda.Enemies[i] // 敵人資料
		// 將世界座標轉換為螢幕座標
		wx, wy := float64(e.X), float64(e.Y)
		sx, sy := camMatrix.Apply(wx, wy)
		// 繪製敵人圖片
		tr := EnemyTileRect[e.Typ]
		op := &ebiten.DrawImageOptions{}
		// 將圖片中心對齊座標
		op.GeoM.Translate(-float64(TileSize)/2, -float64(TileSize)/2)
		op.GeoM.Scale(r.Camera.Zoom, r.Camera.Zoom)
		op.GeoM.Translate(sx, sy)
		// 如果受傷，暫時變色 (紅色閃爍)
		if e.Status == 1 {
			op.ColorScale.Scale(1, 0.5, 0.5, 1)
		}
		screen.DrawImage(r.ResMan.EnemyImage.SubImage(tr).(*ebiten.Image), op)

		// 繪製血條
		radius := float32(TileSize/2) * float32(r.Camera.Zoom)
		healthBarWidth := float32(10) * float32(r.Camera.Zoom)
		healthBarHeight := float32(2) * float32(r.Camera.Zoom)
		healthRatio := e.Health / e.MaxHealth
		vector.FillRect(screen,
			float32(sx)-healthBarWidth/2,
			float32(sy)-radius-healthBarHeight-2,
			healthBarWidth, healthBarHeight,
			color.RGBA{100, 100, 100, 255}, true)
		vector.FillRect(screen,
			float32(sx)-healthBarWidth/2,
			float32(sy)-radius-healthBarHeight-2,
			healthBarWidth*healthRatio, healthBarHeight,
			color.RGBA{0, 255, 0, 255}, true)
	}
}
