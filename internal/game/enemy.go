package game

import (
	. "TowerDefense/internal/vars"
	"fmt"
	"math/rand"
)

type Enemy struct {
	Typ       EnemyType   // 敵人的種類
	Status    EnemyStatus // 敵人的狀態
	X, Y      float32     // 敵人的精確浮點數座標 (用於平滑移動)
	Speed     float32     // 敵人的基礎移動速度
	Health    float32     // 敵人的生命值
	MaxHealth float32     // 敵人的最大生命值
	Damage    float32     // 敵人的攻擊力
	Radius    float32     // 敵人的視野半徑
	Id        int         // 敵人的生成代數
	Active    bool        // 是否活著且在場上
	Reword    int         // 敵人的獎勵值
}

type EnemyArray struct {
	Enemies [MaxEnemies]Enemy
	Count   int // 目前在場上的敵人數量
	NextId  int // 下一個敵人的生成代數
}

func NewEnemyArray() *EnemyArray {
	return &EnemyArray{
		Enemies: [MaxEnemies]Enemy{},
		Count:   0,
		NextId:  0,
	}
}

func (es *EnemyArray) TryAddEnemy(typ EnemyType, x, y float32) bool {
	// 因為陣列是緊湊的，活著的敵人一定在 0 ~ es.Count-1
	// 所以只要 Count 還沒滿，直接加在 es.Count 的位置即可！不需要用迴圈找空位。
	if es.Count >= MaxEnemies {
		return false
	}
	health := float32(10.0)
	speed := float32(1.0)
	if typ == EnemyBoss {
		health = 50.0
		speed = 0.5
	}
	es.Enemies[es.Count] = Enemy{
		Typ:       typ,
		Status:    EnemyHealthy,
		X:         x,
		Y:         y,
		Speed:     speed,
		Health:    health,
		MaxHealth: health,
		Damage:    10.0,
		Radius:    10.0,
		Id:        es.NextId,
		Active:    true,
		Reword:    10, // 初始化獎勵值為 0
	}
	es.Count++
	es.NextId++
	return true
}

// Update 修正：改為倒序走訪，完美避開陣列移除時的索引衝突
func (es *EnemyArray) Update(nav *NavigationMap, gm *GameMap) {
	for i := es.Count - 1; i >= 0; i-- {
		e := &es.Enemies[i]
		if !e.Active {
			continue
		}
		// 算出敵人目前在哪個格子
		gridX := int(e.X / TileSize)
		gridY := int(e.Y / TileSize)
		// 防止越界
		if gridX < 0 || gridX >= GridSize || gridY < 0 || gridY >= GridSize {
			e.Active = false
			es.Remove(i)
			continue
		}
		// 檢查是否到達城堡 (中心點)
		castleX := GridSize / 2
		castleY := GridSize / 2
		if gridX == castleX && gridY == castleY {
			e.Status = EnemyDead
			e.Active = false
			es.Remove(i)
			continue
		}
		// 獲取該格子的推薦移動方向
		dirIdx := nav.Directions[gridY][gridX]
		if dirIdx == -1 {
			continue
		}
		// 根據方向向量移動
		dxVal := float32(Dx[dirIdx])
		dyVal := float32(Dy[dirIdx])
		e.X += dxVal * e.Speed
		e.Y += dyVal * e.Speed
	}
}
func (es *EnemyArray) FindByID(id int) *Enemy {
	for i := 0; i < es.Count; i++ {
		if es.Enemies[i].Active && es.Enemies[i].Id == id {
			return &es.Enemies[i]
		}
	}
	return nil
}

// TakeDamage 修正：健康的敵人也要能承受傷害
func (es *EnemyArray) TakeDamage(index int, amount float32) int {
	if index < 0 || index >= es.Count {
		return 0
	}
	reward := 0
	e := &es.Enemies[index]
	if e.Active && e.Status != EnemyDead {
		e.Health -= amount
		if e.Health < e.MaxHealth {
			e.Status = EnemyHurt
		}
		if e.Health <= 0 {
			reward = e.Reword
			e.Status = EnemyDead
			e.Active = false
			es.Remove(index)
		}
	}
	return reward
}
func (es *EnemyArray) TakeDamageByID(id int, damage float32) int {
	for i := 0; i < es.Count; i++ {
		e := &es.Enemies[i]
		if e.Id != id || !e.Active {
			continue
		}
		e.Health -= damage

		if e.Health <= 0 {
			es.Remove(i)
			return e.Reword
		}
		return 0
	}
	return 0
}

// Remove 修正：修正了混亂的大括號邏輯與 Active 狀態檢查
func (es *EnemyArray) Remove(idx int) {
	if idx >= es.Count || idx < 0 {
		fmt.Printf("Warning: Attempt to remove enemy with invalid index %d.\n", idx)
		return
	}
	// 無論 e.Active 是真是假，既然要 Remove，我們就執行覆蓋
	// 把最後一個活著的敵人 (es.Count-1) 移到當前被刪除的空位 (idx)
	lastIdx := es.Count - 1
	if idx != lastIdx {
		es.Enemies[idx] = es.Enemies[lastIdx]
	}
	// 清空最後一個位置的資料，釋放記憶體或重設狀態
	es.Enemies[lastIdx] = Enemy{
		Typ:    EnemyNone,
		Status: EnemyHealthy,
		Active: false,
	}
	es.Count--
	fmt.Printf("Enemy removed from index %d. New Count: %d\n", idx, es.Count)
}

// GetRandomEdgePosition 返回地圖邊緣的一個隨機座標，並確保該位置有路徑通往城堡
func GetRandomEdgePosition(gm *GameMap, nav *NavigationMap) (float32, float32) {
	for attempt := 0; attempt < 100; attempt++ {
		edge := rand.Intn(4) // 0:上, 1:下, 2:左, 3:右
		var x, y float32
		switch edge {
		case 0: // 上邊緣
			x = rand.Float32() * ScreenWidth
			y = 0
		case 1: // 下邊緣
			x = rand.Float32() * ScreenWidth
			y = ScreenHeight - 1
		case 2: // 左邊緣
			x = 0
			y = rand.Float32() * ScreenHeight
		case 3: // 右邊緣
			x = ScreenWidth - 1
			y = rand.Float32() * ScreenHeight
		}
		// 檢查地形是否可通行
		gridX := int(x / TileSize)
		gridY := int(y / TileSize)
		// 邊界安全檢查
		if gridX < 0 {
			gridX = 0
		} else if gridX >= GridSize {
			gridX = GridSize - 1
		}
		if gridY < 0 {
			gridY = 0
		} else if gridY >= GridSize {
			gridY = GridSize - 1
		}
		// 關鍵修正：除了檢查地形成本，還要檢查導航地圖中該點是否有方向引導（即是否有路通往城堡）
		if nav.Directions[gridY][gridX] != -1 {
			return x, y
		}
	}
	// 如果嘗試多次都沒找到，回傳地圖中心（城堡位置），作為保險
	return float32(GridSize / 2 * TileSize), float32(GridSize / 2 * TileSize)
}
