package game

import (
	. "TowerDefense/internal/vars"
	"container/heap"
	"math"
)

// 為了在 Go 實作 Dijkstra 尋路，我們需要一個優先佇列 (Priority Queue)
type PathNode struct {
	X, Y int
	Cost int // 累計到城堡的通行成本
}

type PriorityQueue []*PathNode

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].Cost < pq[j].Cost } // 成本小的優先
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*PathNode)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// 建立一個地圖導航網格，存放每一格走到城堡的方向
type NavigationMap struct {
	Directions [GridSize][GridSize]int // 存 dx, dy 的索引 (0:上, 1:下, 2:左, 3:右, -1:無法抵達)
	TotalCosts [GridSize][GridSize]int // 每一格到城堡的總花費
}

// 計算從全地圖到城堡的最短權重路徑 (Dijkstra)
func CalculateNavigation(gm *GameMap, ta *TowerArray) *NavigationMap {
	nav := &NavigationMap{}
	centerX := GridSize / 2
	centerY := GridSize / 2

	// 初始化成本為無限大，方向為 -1
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			nav.TotalCosts[y][x] = math.MaxInt32
			nav.Directions[y][x] = -1
		}
	}

	// 建立優先佇列並將城堡（起點）放進去
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	nav.TotalCosts[centerY][centerX] = 0
	heap.Push(&pq, &PathNode{X: centerX, Y: centerY, Cost: 0})

	for pq.Len() > 0 {
		curr := heap.Pop(&pq).(*PathNode)

		// 如果找到的成本已經比記錄的大，跳過
		if curr.Cost > nav.TotalCosts[curr.Y][curr.X] {
			continue
		}

		// 檢查四個鄰居 (使用你宣告的 dx, dy)
		for i := 0; i < 4; i++ {
			nx := curr.X + Dx[i]
			ny := curr.Y + Dy[i]

			// 邊界檢查
			if nx < 0 || nx >= GridSize || ny < 0 || ny >= GridSize {
				continue
			}

			tile := gm.Grid[ny][nx]
			cost := TerrainCost[tile.Type]

			// 檢查該位置是否有塔
			hasTower := false
			if ta != nil {
				for k := 0; k < ta.Count; k++ {
					if int(ta.Towers[k].X) == nx && int(ta.Towers[k].Y) == ny {
						hasTower = true
						break
					}
				}
			}

			// 如果是山脈、湖泊，或是上面已經蓋了玩家的塔，則無法通行 (-1)
			if cost == -1 || hasTower {
				continue
			}

			// 計算走到這格的總成本
			nextCost := curr.Cost + cost

			// 如果這條路徑比之前找到的更省時
			if nextCost < nav.TotalCosts[ny][nx] {
				nav.TotalCosts[ny][nx] = nextCost
				// 記錄方向：因為我們是從城堡「反向」找出來的，
				// 所以鄰居要走到當前這格，方向必須反轉（使用你宣告的 reverseDir）
				nav.Directions[ny][nx] = ReverseDir[i]

				heap.Push(&pq, &PathNode{X: nx, Y: ny, Cost: nextCost})
			}
		}
	}
	return nav
}
