package game

import (
	. "TowerDefense/internal/vars"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tower struct {
	Typ      TowerType
	X, Y     float32 // 網格座標 (grid units)，例如 (0, 0) 代表左上角第一個格子
	Level    int
	Range    float32 // 攻擊範圍（網格單位）
	Damage   float32 // 攻擊傷害
	Cooldown float32 // 當前冷卻時間計時器 (秒)，0 代表可以攻擊
	MaxCD    float32 // 攻擊間隔/冷卻上限 (秒)
	Active   bool    // 該防禦塔是否啟用/存在
}

type TowerArray struct {
	Towers [MaxTowers]Tower
	Count  int
}

func NewTowerArray() *TowerArray {
	return &TowerArray{
		Towers: [MaxTowers]Tower{},
		Count:  0,
	}
}

// TryAddTower 嘗試在網格上放置一座建築物
func (ta *TowerArray) TryAddTower(typ TowerType, gridX, gridY int, gm *GameMap) bool {
	// 檢查防禦塔數量是否已達上限
	if ta.Count >= MaxTowers {
		return false
	}
	// 檢查放置邊界，防止陣列越界 (Index Out of Bounds)
	if gridX < 0 || gridX >= GridSize || gridY < 0 || gridY >= GridSize {
		return false
	}
	// 檢查地形：如果是不可通行區域（如水、山脈，代價為 -1）或是城堡（主基地），則無法建造
	tile := gm.Grid[gridY][gridX]
	if TerrainCost[tile.Type] == -1 || tile.Type == TerrainCastle {
		return false
	}
	// 檢查該網格上是否已經存在其他防禦塔，避免重疊建造
	for i := 0; i < ta.Count; i++ {
		if int(ta.Towers[i].X) == gridX && int(ta.Towers[i].Y) == gridY {
			return false
		}
	}
	// 初始化基礎屬性（預設為 TowerBasic 基礎塔）
	rangeVal := float32(5.0)
	damage := float32(20.0)
	maxCD := float32(1.0)
	// 根據防禦塔類型微調數值
	switch typ {
	case TowerSniper:
		rangeVal = 10.0
		damage = 500.0
		maxCD = 3.0 // 攻速慢、單發傷害極高、範圍廣
	case TowerRapid:
		rangeVal = 5.0
		damage = 5.0 // 修正：原本程式碼誤填為 500.0，機槍塔傷害通常較低但攻速極快
		maxCD = 0.2  // 攻速極快
	}
	// 將新塔寫入陣列，並遞增計數器
	ta.Towers[ta.Count] = Tower{
		Typ:      typ,
		X:        float32(gridX),
		Y:        float32(gridY),
		Level:    1,
		Range:    rangeVal,
		Damage:   damage,
		Cooldown: 0,
		MaxCD:    maxCD,
		Active:   true,
	}
	ta.Count++
	return true
}

// Update 每幀更新防禦塔狀態，包含 CD 計算、尋找目標與發射子彈
func (ta *TowerArray) Update(dt float32, enemies *EnemyArray, projectiles *ProjectileArray) {
	// 將 tileSize 轉為 float32，避免在迴圈內反覆轉型影響效能
	fTileSize := float32(TileSize)

	for i := 0; i < ta.Count; i++ {
		t := &ta.Towers[i]
		if !t.Active {
			continue
		}
		// 處理冷卻時間倒數
		if t.Cooldown > 0 {
			t.Cooldown -= dt
		}
		// 冷卻完畢，開始索敵與攻擊
		if t.Cooldown <= 0 {
			targetIdx := -1
			// 預設最大搜尋距離的平方（只考慮範圍內的敵人）
			minDist := t.Range * t.Range
			// 遍歷所有敵人，尋找最靠近塔的目標
			for j := 0; j < enemies.Count; j++ {
				e := &enemies.Enemies[j]
				if !e.Active {
					continue
				}
				// Bug 修正：將敵人的像素座標轉換為網格座標（加入 float32 轉型防止編譯錯誤）
				ex := e.X / fTileSize
				ey := e.Y / fTileSize
				// 計算塔與敵人之間的平方距離
				dist := (ex-t.X)*(ex-t.X) + (ey-t.Y)*(ey-t.Y)
				// 邏輯優化：使用 '<' 確保穩定鎖定最靠近的單一目標，減少蓋格誤差
				if dist < minDist {
					minDist = dist
					targetIdx = j
				}
			}
			// 如果範圍內有找到有效的敵人
			if targetIdx != -1 {
				target := &enemies.Enemies[targetIdx]
				// Bug 修正：計算子彈生成的像素座標（塔的中心點），將 tileSize 轉為 float32
				startX := t.X*fTileSize + fTileSize/2
				startY := t.Y*fTileSize + fTileSize/2
				// 生成子彈並導向該敵人的 ID
				projectiles.Add(startX, startY, target.Id, t.Damage)
				// 重設冷卻時間計時器
				t.Cooldown = t.MaxCD
			}
		}
	}
}

// hasPathToCastle 檢查地圖邊緣是否有路徑可以到達城堡
func hasPathToCastle(nav *NavigationMap) bool {
	for x := 0; x < GridSize; x++ {
		if nav.Directions[0][x] != -1 || nav.Directions[GridSize-1][x] != -1 {
			return true
		}
	}
	for y := 0; y < GridSize; y++ {
		if nav.Directions[y][0] != -1 || nav.Directions[y][GridSize-1] != -1 {
			return true
		}
	}
	return false
}

func (ta *TowerArray) Place(g *Game, towerType TowerType) {
	mx, my := ebiten.CursorPosition()
	screenW := float64(ScreenWidth)
	screenH := float64(ScreenHeight)
	wx := (float64(mx)-screenW/2)/g.Renderer.Camera.Zoom + g.Renderer.Camera.X
	wy := (float64(my)-screenH/2)/g.Renderer.Camera.Zoom + g.Renderer.Camera.Y
	gridX := int(wx / TileSize)
	gridY := int(wy / TileSize)

	if g.Money >= TowerCost[towerType] {
		if g.Towers.TryAddTower(towerType, gridX, gridY, g.Map) {
			newNav := CalculateNavigation(g.Map, g.Towers)
			if hasPathToCastle(newNav) {
				g.Nav = newNav
				g.Money -= TowerCost[towerType]
				fmt.Printf("Placed %s ! Money left: %d\n", TowerName[towerType], g.Money)
			} else {
				g.Towers.Count--
				fmt.Println("Cannot place tower here: path to castle would be blocked!")
			}
		}
	}
}
