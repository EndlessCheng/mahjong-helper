package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTilesToMergedStr(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("13m 1p", TilesToStr([]int{0, 2, 9}))
	assert.Equal("134m", TilesToStr([]int{0, 2, 3}))
	assert.Equal("1m 1p 1s 17z", TilesToStr([]int{0, 9, 18, 27, 33}))
	assert.Equal("19m 19p 19s 17z", TilesToStr([]int{0, 8, 9, 17, 18, 26, 27, 33}))
}

func TestConvert(t *testing.T) {
	for _, tiles := range []string{
		//"123456789m123456789p123456789s1234567z",
		"123456789m 123456789p 123456789s 1234567z",
		//"",
		"1114569m 1456999p 1456669s 14567777z",
		"6m",
		"7p",
		"6s",
		"3z",
		"7z",
		"7p 7s",
		"45s",
		//"0m 0p 0s",
		//"307m 307p 304s 1234567z",
	} {
		assert.Equal(t, tiles, Tiles34ToStr(MustStrToTiles34(tiles)))
	}
}
