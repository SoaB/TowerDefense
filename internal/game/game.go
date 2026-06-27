package game

import (
	"TowerDefense/internal/renderer"
	. "TowerDefense/internal/vars"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Map         *GameMap
	Nav         *NavigationMap
	Enemies     *EnemyArray
	Towers      *TowerArray
	Projectiles *ProjectileArray
	Ticks       int
	Money       int
	Renderer    *renderer.Renderer
}

// convertGridToRendererFormat 將 GameMap 的 Grid 轉換為 renderer 需要的格式
func convertGridToRendererFormat(gm *GameMap) *renderer.TerrainDataArr {
	tda := &renderer.TerrainDataArr{} // 初始化地形數據陣列
	for y := 0; y < GridSize; y++ {
		for x := 0; x < GridSize; x++ {
			tda.Grid[y][x].TerrainType = gm.Grid[y][x].Type
		}
	}
	return tda
}

// convertTowersToRendererFormat 將 TowerArray 轉換為 renderer 需要的格式
func convertTowersToRendererFormat(towers *TowerArray) *renderer.TowerDataArr {
	result := &renderer.TowerDataArr{}
	result.Count = towers.Count
	for i := 0; i < towers.Count; i++ {
		t := towers.Towers[i]
		result.Towers[i].X = int(t.X)
		result.Towers[i].Y = int(t.Y)
		result.Towers[i].Typ = int(t.Typ)
		result.Towers[i].Range = t.Range
		result.Towers[i].Active = t.Active
	}
	return result
}

// convertEnemiesToRendererFormat 將 EnemyArray 轉換為 renderer 需要的格式
func convertEnemiesToRendererFormat(enemies *EnemyArray) *renderer.EnemyDataArr {
	result := &renderer.EnemyDataArr{} // 初始化敵人數據陣列
	result.Count = enemies.Count
	for i := 0; i < enemies.Count; i++ {
		e := enemies.Enemies[i]
		status := 0
		if e.Status == EnemyHurt {
			status = 1
		}
		result.Enemies[i].X = float64(e.X)
		result.Enemies[i].Y = float64(e.Y)
		result.Enemies[i].Typ = int(e.Typ)
		result.Enemies[i].Health = e.Health
		result.Enemies[i].MaxHealth = e.MaxHealth
		result.Enemies[i].Status = status
	}
	return result
}

// convertProjectilesToRendererFormat 將 ProjectileArray 轉換為 renderer 需要的格式
func convertProjectilesToRendererFormat(projectiles *ProjectileArray) *renderer.ProjectileDataArr {
	result := &renderer.ProjectileDataArr{} // 初始化 projectiles 数据数组
	result.Count = projectiles.Count
	for i := 0; i < projectiles.Count; i++ {
		p := projectiles.Projectiles[i]
		result.Projectiles[i].X = float32(p.X)
		result.Projectiles[i].Y = float32(p.Y)
		result.Projectiles[i].Active = p.Active
	}
	return result
}

func NewGame() *Game {
	g := &Game{}
	g.Money = InitialMoney                       // 初始化金錢
	g.Enemies = NewEnemyArray()                  // 初始化敵人列表
	g.Towers = NewTowerArray()                   // 初始化塔列表
	g.Map = NewGameMap(GenerateMapScale)         // 初始化遊戲地圖
	g.Nav = CalculateNavigation(g.Map, g.Towers) // 計算初始導航地圖
	g.Renderer = renderer.NewRenderer()          // 初始化資源管理器
	g.Renderer.DrawTerrainLayer(convertGridToRendererFormat(g.Map))
	g.Renderer.DrawCastleLayer()
	g.Projectiles = NewProjectileArray()
	return g
}

func (g *Game) Update() error {
	g.Ticks++
	g.Renderer.Camera.Update()
	// 測試：每 60 幀生成一個敵人
	if g.Ticks%60 == 0 {
		// 在地圖邊緣隨機位置生成敵人，並確保有路徑
		rx, ry := GetRandomEdgePosition(g.Map, g.Nav)
		g.Enemies.TryAddEnemy(EnemyMinion, rx, ry)
	}
	if g.Ticks%300 == 0 {
		rx, ry := GetRandomEdgePosition(g.Map, g.Nav)
		g.Enemies.TryAddEnemy(EnemyBoss, rx, ry)
	}
	// 更新塔和攻擊，並獲得擊殺金錢
	g.Towers.Update(1.0/60.0, g.Enemies, g.Projectiles)
	g.Money += g.Projectiles.Update(g.Enemies)
	// 更新敵人移動和路徑
	g.Enemies.Update(g.Nav, g.Map)

	// 點擊左鍵放置基礎塔
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.Towers.Place(g, TowerBasic)
	}
	// 點擊右鍵放置狙擊塔
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		g.Towers.Place(g, TowerSniper)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Renderer.DrawLayer(screen)
	g.Renderer.DrawEnemies(screen, convertEnemiesToRendererFormat(g.Enemies))             // 更新敵人數據
	g.Renderer.DrawProjectiles(screen, convertProjectilesToRendererFormat(g.Projectiles)) // 更新 projectiles 数据
	g.Renderer.DrawTowers(screen, convertTowersToRendererFormat(g.Towers))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
