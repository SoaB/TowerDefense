package renderer

import (
	. "TowerDefense/internal/vars"

	"github.com/hajimehoshi/ebiten/v2"
)

type Renderer struct {
	ResMan *ResourceManager
	Camera *Camera
}

func NewRenderer() *Renderer {
	renderer := &Renderer{}
	renderer.ResMan = NewResourceManager()
	renderer.Camera = NewCamera(float64(ScreenWidth/2), float64(ScreenHeight/2), 1.0)
	return renderer
}
func (r *Renderer) Draw(screen *ebiten.Image) {
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
