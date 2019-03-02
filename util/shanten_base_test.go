package util

import (
	"testing"
)

func TestCalculateShanten(t *testing.T) {
	t.Log(CalculateShanten(MustStrToTiles34("13579m 12357s 135p"), false) == 3)
	t.Log(CalculateShanten(MustStrToTiles34("123456789m 147s 14m"), false) == 1)
	t.Log(CalculateShanten(MustStrToTiles34("123456789m 147s 1m"), false) == 2)
	t.Log(CalculateShanten(MustStrToTiles34("258m 258s 258p 12345z"), true))
}

func BenchmarkCalculateShanten(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MustStrToTiles34("147m 147s 147p 12345z")
	}
}
