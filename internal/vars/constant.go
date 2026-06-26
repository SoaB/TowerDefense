package vars

// 版本號
const (
	TowerDefenseVersion = "1.0.0"
)

// 遊戲
const (
	// 游戏名稱
	GameName = "Tower Defense"
	// 游戏描述
	GameDescription = "A tower defense game."
	// 游戏作者
	GameAuthor = "Your Name"
	// 游戏版本
	GameVersion = TowerDefenseVersion
)

// 游戏配置
const (
	// 游戏宽度
	GameWidth = 960
	// 游戏高度
	GameHeight = 960
	// 游戏帧率
	GameFrameRate = 60
)

// 游戏资源路径
const (
	ResourcePath = "./assets"
)

const (
	GridSize         = 30 // 网格大小
	TileSize         = 16 // 块大小
	ScreenWidth      = GridSize * TileSize
	ScreenHeight     = GridSize * TileSize
	MaxEnemies       = 100
	MaxTowers        = 100
	MaxProjectiles   = 200
	GenerateMapScale = 0.08
)

type TerrainType int

const (
	TerrainGrass TerrainType = iota
	TerrainHill
	TerrainSand
	TerrainMud
	TerrainStone
	TerrainMountain
	TerrainLake
	TerrainCastle
	TerrainTotal
)

type EnemyType int8

const (
	EnemyNone EnemyType = iota
	EnemyMinion
	EnemyBoss
	EnemyTotal
)

type EnemyStatus int8

const (
	EnemyHealthy EnemyStatus = iota
	EnemyHurt
	EnemyDead
)
