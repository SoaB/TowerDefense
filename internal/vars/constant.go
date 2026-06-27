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
	InitialMoney     = 100
)

var Dx = [4]int{0, 0, -1, 1}
var Dy = [4]int{-1, 1, 0, 0}
var ReverseDir = [4]int{1, 0, 3, 2}

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

var TerrainCost = [TerrainTotal]int{
	TerrainGrass:    2,
	TerrainHill:     3,
	TerrainSand:     5,
	TerrainMud:      8,
	TerrainStone:    1,
	TerrainMountain: -1,
	TerrainLake:     -1,
	TerrainCastle:   0,
}

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

type TowerType int8

const (
	TowerNone TowerType = iota
	TowerBasic
	TowerSniper
	TowerRapid
	TowerTotal
)

var TowerCost = [TowerTotal]int{
	TowerNone:   0,
	TowerBasic:  10,
	TowerSniper: 20,
	TowerRapid:  30,
}
var TowerName = [TowerTotal]string{
	TowerNone:   "None",
	TowerBasic:  "Basic Tower",
	TowerSniper: "Sniper Tower",
	TowerRapid:  "Rapid Tower",
}
