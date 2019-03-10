package util

import (
	"testing"
)

func TestCalculateShantenClosed(t *testing.T) {
	t.Log(CalculateShanten(MustStrToTiles34("13579m 12357s 135p"), false) == 3)
	t.Log(CalculateShanten(MustStrToTiles34("123456789m 147s 14m"), false) == 1)
	t.Log(CalculateShanten(MustStrToTiles34("123456789m 147s 1m"), false) == 2)
	t.Log(CalculateShanten(MustStrToTiles34("258m 258s 258p 12345z"), true) == 8) // 不考虑国士无双和七对子的最大向听
	t.Log(CalculateShanten(MustStrToTiles34("123456789m 1134p"), false) == 0)
	t.Log(CalculateShanten(MustStrToTiles34("123456789m 11345p"), false) == -1)
}

func TestCalculateShantenOpen(t *testing.T) {
	//t.Log(CalculateShanten(MustStrToTiles34("123m"), true))
	t.Log(CalculateShanten(MustStrToTiles34("2247m"), true) == 1)
	t.Log(CalculateShanten(MustStrToTiles34("11234m"), true) == -1)
}

func BenchmarkCalculateShantenClosed(b *testing.B) {
	tiles34 := MustStrToTiles34("13579m 12357s 135p")
	for i := 0; i < b.N; i++ {
		// 1931 ns/op
		CalculateShanten(tiles34, false)
	}
}

func BenchmarkCalculateShantenOpen(b *testing.B) {
	tiles34 := MustStrToTiles34("2247m")
	for i := 0; i < b.N; i++ {
		// 150 ns/op
		CalculateShanten(tiles34, true)
	}
}
