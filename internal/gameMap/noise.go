package gameMap

import (
	"math/rand"
	"time"
)

func fastFloor(x float32) int {
	if float32(int(x)) <= x {
		return int(x)
	}
	return int(x) - 1
}

var perm [512]uint8

func init() {
	// 使用目前時間作為隨機種子
	rand.Seed(time.Now().UnixNano())

	// 初始化置換表
	p := make([]uint8, 256)
	for i := range p {
		p[i] = uint8(i)
	}

	// 隨機打亂
	rand.Shuffle(len(p), func(i, j int) {
		p[i], p[j] = p[j], p[i]
	})

	// 複製到 512 長度的數組中以避免取模運算
	for i := 0; i < 512; i++ {
		perm[i] = p[i&255]
	}
}

func grad2(hash uint8, x, y float32) float32 {
	h := hash & 7
	u := x
	v := y
	if h >= 4 {
		u = y
		v = x
	}
	if h&1 != 0 {
		u = -u
	}
	if h&2 != 0 {
		v = -v
	}
	return u + v
}

func Snoise2(x, y float32) float32 {
	const F2 float32 = 0.366025403 // 0.5*(sqrt(3.0)-1.0)
	const G2 float32 = 0.211324865 // (3.0-sqrt(3.0))/6.0

	var n0, n1, n2 float32

	s := (x + y) * F2
	xs := x + s
	ys := y + s
	i := fastFloor(xs)
	j := fastFloor(ys)

	t := float32(i+j) * G2
	X0 := float32(i) - t
	Y0 := float32(j) - t
	x0 := x - X0
	y0 := y - Y0

	var i1, j1 uint8
	if x0 > y0 {
		i1 = 1
		j1 = 0
	} else {
		i1 = 0
		j1 = 1
	}

	x1 := x0 - float32(i1) + G2
	y1 := y0 - float32(j1) + G2
	x2 := x0 - 1.0 + 2.0*G2
	y2 := y0 - 1.0 + 2.0*G2

	ii := uint8(i & 255)
	jj := uint8(j & 255)

	t0 := 0.5 - x0*x0 - y0*y0
	if t0 >= 0.0 {
		t0 *= t0
		// 修正：修正 perm 雙重查表的索引方式
		n0 = t0 * t0 * grad2(perm[uint16(ii)+uint16(perm[jj])], x0, y0)
	}

	t1 := 0.5 - x1*x1 - y1*y1
	if t1 >= 0.0 {
		t1 *= t1
		n1 = t1 * t1 * grad2(perm[uint16(ii)+uint16(i1)+uint16(perm[uint16(jj)+uint16(j1)])], x1, y1)
	}

	t2 := 0.5 - x2*x2 - y2*y2
	if t2 >= 0.0 {
		t2 *= t2
		n2 = t2 * t2 * grad2(perm[uint16(ii)+1+uint16(perm[uint16(jj)+1])], x2, y2)
	}
	// 關鍵修正：Simplex Noise 需要乘上一個放大係數（約 70），否則數值會太小
	return 70.0 * (n0 + n1 + n2)
}
