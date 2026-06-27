package game

import (
	. "TowerDefense/internal/vars"
	"math"
)

type Projectile struct {
	X, Y     float32
	TargetID int
	Speed    float32
	Damage   float32
	Active   bool
}

type ProjectileArray struct {
	Projectiles [MaxProjectiles]Projectile
	Count       int
}

func NewProjectileArray() *ProjectileArray {
	return &ProjectileArray{}
}

func (ps *ProjectileArray) Add(x, y float32, targetId int, damage float32) {
	if ps.Count >= len(ps.Projectiles) {
		return
	}

	ps.Projectiles[ps.Count] = Projectile{
		X:        x,
		Y:        y,
		TargetID: targetId,
		Speed:    4.0,
		Damage:   damage,
		Active:   true,
	}
	ps.Count++
}

func (ps *ProjectileArray) Update(enemies *EnemyArray) int {
	totalReward := 0
	for i := ps.Count - 1; i >= 0; i-- {
		p := &ps.Projectiles[i]
		if !p.Active {
			continue
		}
		target := enemies.FindByID(p.TargetID)
		if target == nil {
			ps.Remove(i)
			continue
		}
		dx := target.X - p.X
		dy := target.Y - p.Y
		// 計算距離的平方
		distSq := dx*dx + dy*dy
		// 計算速度的平方
		speedSq := p.Speed * p.Speed
		// 如果距離平方小於速度平方，代表這一個 frame 就會撞到
		if distSq < speedSq {
			totalReward += enemies.TakeDamageByID(p.TargetID, p.Damage)
			ps.Remove(i)
			continue
		}
		// 只有在需要移動時才開根號，減少不必要的運算
		dist := float32(math.Sqrt(float64(distSq)))
		p.X += dx / dist * p.Speed
		p.Y += dy / dist * p.Speed
	}
	return totalReward
}

func (ps *ProjectileArray) Remove(idx int) {
	last := ps.Count - 1
	if idx != last {
		ps.Projectiles[idx] = ps.Projectiles[last]
	}
	ps.Projectiles[last].Active = false
	ps.Count--
}
