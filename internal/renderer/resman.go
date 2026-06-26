package renderer

import (
	. "TowerDefense/internal/vars"
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type ResourceManager struct {
	TerrainLayer *ebiten.Image
	CastleLayer  *ebiten.Image
	TileImage    *ebiten.Image
	EnemyImage   *ebiten.Image
}

func loadResources(fn string) *ebiten.Image {
	f, err := os.Open(fn)
	if err != nil {
		fmt.Printf("Error opening %s: %s\n", fn, err)
		return nil
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil
	}
	return ebiten.NewImageFromImage(img)
}

func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		TerrainLayer: ebiten.NewImage(ScreenWidth, ScreenHeight),
		CastleLayer:  ebiten.NewImage(ScreenWidth, ScreenHeight),
		TileImage:    loadResources("assets/greentile.png"),
		EnemyImage:   loadResources("assets/enemys.png"),
	}
}
