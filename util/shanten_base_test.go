package util

import (
	"testing"
)

func TestCalculateShantenClosed(t *testing.T) {
	t.Log(CalculateShanten(MustStrToTiles34("13579m 12357s 135p")) == 3)
	t.Log(CalculateShanten(MustStrToTiles34("123456789m 147s 14m")) == 1)
	t.Log(CalculateShanten(MustStrToTiles34("123456789m 147s 1m")) == 2)
	t.Log(CalculateShanten(MustStrToTiles34("258m 258s 258p 12345z")) == 6)
	t.Log(CalculateShanten(MustStrToTiles34("123456789m 1134p")) == 0)
	t.Log(CalculateShanten(MustStrToTiles34("123456789m 11345p")) == -1)
}

func TestCalculateShantenOpen(t *testing.T) {
	t.Log(CalculateShanten(MustStrToTiles34("2247m")) == 1)
	t.Log(CalculateShanten(MustStrToTiles34("11234m")) == -1)
}

func BenchmarkCalculateShantenClosed(b *testing.B) {
	tiles34 := MustStrToTiles34("13579m 12357s 135p")
	for i := 0; i < b.N; i++ {
		// 1806 ns/op
		CalculateShanten(tiles34)
	}
}

func BenchmarkCalculateShantenOpen(b *testing.B) {
	tiles34 := MustStrToTiles34("2247m")
	for i := 0; i < b.N; i++ {
		// 146 ns/op
		CalculateShanten(tiles34)
	}
}
