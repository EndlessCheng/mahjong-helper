package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCalculateShantenOfChiitoi(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(-1, CalculateShantenOfChiitoi(MustStrToTiles34("11223344556677z")))
	assert.Equal(0, CalculateShantenOfChiitoi(MustStrToTiles34("1223344556677z")))
	assert.Equal(0, CalculateShantenOfChiitoi(MustStrToTiles34("1m 1223344556677z")))
	assert.Equal(0, CalculateShantenOfChiitoi(MustStrToTiles34("1223344556677z")))
	assert.Equal(1, CalculateShantenOfChiitoi(MustStrToTiles34("12m 123344556677z")))
	assert.Equal(1, CalculateShantenOfChiitoi(MustStrToTiles34("1m 123344556677z")))
	assert.Equal(1, CalculateShantenOfChiitoi(MustStrToTiles34("11222233445566z")))
	assert.Equal(5, CalculateShantenOfChiitoi(MustStrToTiles34("11112222333344z")))
	assert.Equal(3, CalculateShantenOfChiitoi(MustStrToTiles34("33m 5555p 66s 556666z")))
	assert.Equal(2, CalculateShantenOfChiitoi(MustStrToTiles34("577m 23677p 245577s")))
}

func TestCalculateShantenOfNormal(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(8, CalculateShantenOfNormal(MustStrToTiles34("19m 19p 19s 1234567z"), 14))
	assert.Equal(7, CalculateShantenOfNormal(MustStrToTiles34("19m 199p 19s 1234567z"), 13))
	assert.Equal(3, CalculateShantenOfNormal(MustStrToTiles34("577m 23677p 245577s"), 14))
}

func TestCalculateShanten(t *testing.T) {
	assert := assert.New(t)

	// Closed
	assert.Equal(1, CalculateShanten(MustStrToTiles34("33m 5555p 66s 556666z")))
	assert.Equal(4, CalculateShanten(MustStrToTiles34("13579m 13579s 135p")))
	assert.Equal(3, CalculateShanten(MustStrToTiles34("13579m 12379s 135p")))
	assert.Equal(1, CalculateShanten(MustStrToTiles34("123456789m 147s 14m")))
	assert.Equal(2, CalculateShanten(MustStrToTiles34("123456789m 147s 1m")))
	assert.Equal(6, CalculateShanten(MustStrToTiles34("258m 258s 258p 12345z"))) // 和牌最远
	assert.Equal(0, CalculateShanten(MustStrToTiles34("123456789m 1134p")))
	assert.Equal(-1, CalculateShanten(MustStrToTiles34("123456789m 11345p")))

	// Open
	assert.Equal(0, CalculateShanten(MustStrToTiles34("1m")))
	assert.Equal(0, CalculateShanten(MustStrToTiles34("1555m")))
	assert.Equal(1, CalculateShanten(MustStrToTiles34("2247m")))
	assert.Equal(-1, CalculateShanten(MustStrToTiles34("11234m")))
	assert.Equal(1, CalculateShanten(MustStrToTiles34("5555m")))
	assert.Equal(1, CalculateShanten(MustStrToTiles34("5555z")))
}

func BenchmarkCalculateShantenClosed(b *testing.B) {
	tiles34 := MustStrToTiles34("13579m 12357s 135p")
	for i := 0; i < b.N; i++ {
		// 1758 ns/op
		CalculateShanten(tiles34)
	}
}

func BenchmarkCalculateShantenOpen(b *testing.B) {
	tiles34 := MustStrToTiles34("2247m")
	for i := 0; i < b.N; i++ {
		// 100 ns/op
		CalculateShanten(tiles34)
	}
}
