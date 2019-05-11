package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCalculateShantenClosed(t *testing.T) {
	assert.Equal(t, 4, CalculateShanten(MustStrToTiles34("13579m 13579s 135p")))
	assert.Equal(t, 3, CalculateShanten(MustStrToTiles34("13579m 12379s 135p")))
	assert.Equal(t, 1, CalculateShanten(MustStrToTiles34("123456789m 147s 14m")))
	assert.Equal(t, 2, CalculateShanten(MustStrToTiles34("123456789m 147s 1m")))
	assert.Equal(t, 6, CalculateShanten(MustStrToTiles34("258m 258s 258p 12345z"))) // 和牌最远
	assert.Equal(t, 0, CalculateShanten(MustStrToTiles34("123456789m 1134p")))
	assert.Equal(t, -1, CalculateShanten(MustStrToTiles34("123456789m 11345p")))
}

func TestCalculateShantenOpen(t *testing.T) {
	assert.Equal(t, 0, CalculateShanten(MustStrToTiles34("1m")))
	assert.Equal(t, 1, CalculateShanten(MustStrToTiles34("2247m")))
	assert.Equal(t, -1, CalculateShanten(MustStrToTiles34("11234m")))
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
