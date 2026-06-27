package renderer

import (
	. "TowerDefense/internal/vars"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ProjectileData 子彈資料結構 (用於繪製)
type ProjectileData struct {
	X, Y   float32
	Active bool
}

type ProjectileDataArr struct {
	Projectiles [MaxProjectiles]ProjectileData
	Count       int
}

func (r *Renderer) DrawProjectiles(screen *ebiten.Image, projectiles *ProjectileDataArr) {
	for i := 0; i < projectiles.Count; i++ {
		p := projectiles.Projectiles[i]
		if !p.Active {
			continue
		}
		vector.FillCircle(
			screen,
			p.X, p.Y,
			3,
			color.White,
			false,
		)
	}
}
