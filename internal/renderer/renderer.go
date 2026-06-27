package renderer

import (
	. "TowerDefense/internal/vars"
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
