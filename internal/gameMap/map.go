package gameMap

import (
	. "TowerDefense/internal/vars"
	"math"
	"math/rand/v2"
)

type Tile struct {
	Type     TerrainType
	HasTower bool
}

type GameMap struct {
	Grid [GridSize][GridSize]Tile
}

const SafeRadius = 4 // 安全城堡半徑
// 根據雜訊產生多種地形，並保留中央城堡區
func NewGameMap(scale float32) *GameMap {
	gm := &GameMap{}
	centerX := GridSize / 2
	centerY := GridSize / 2

	// 加入隨機偏移
	offsetX := rand.Float32() * 1000
	offsetY := rand.Float32() * 1000

	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			// 1. 檢查中央城堡安全區
			dist := math.Sqrt(float64((x-centerX)*(x-centerX) + (y-centerY)*(y-centerY)))
			if dist < SafeRadius {
				if x == centerX && y == centerY {
					gm.Grid[y][x] = Tile{Type: TerrainCastle}
				} else {
					gm.Grid[y][x] = Tile{Type: TerrainGrass} // 城堡周圍預設草地
				}
				continue
			}

			// 2. 保護區外利用 Simplex Noise 分層 (映射到 0.0 ~ 1.0)
			nx := float32(x)*scale + offsetX
			ny := float32(y)*scale + offsetY
			elevation := (Snoise2(nx, ny) + 1.0) / 2.0

			var t TerrainType
			// 依據高度細分 6 種自然地形
			switch {
			case elevation < 0.25:
				t = TerrainLake // 低窪變湖泊 (-1 無法通行)
			case elevation < 0.35:
				t = TerrainSand // 水邊沙地 (5 慢速)
			case elevation < 0.50:
				t = TerrainGrass // 普通草地 (2 正常)
			case elevation < 0.65:
				t = TerrainStone // 乾淨石地 (1 快速路徑)
			case elevation < 0.75:
				t = TerrainHill // 丘陵地 (3 微慢)
			case elevation < 0.85:
				t = TerrainMud // 泥濘地 (8 極慢)
			default:
				t = TerrainMountain // 高山 (-1 無法通行)
			}

			gm.Grid[y][x] = Tile{Type: t}
		}
	}
	return gm
}
