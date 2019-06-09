package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCalculateRiskTiles34(t *testing.T) {
	safeTiles34 := make([]bool, 34)
	for _, tile := range MustStrToTiles("4m 5p 3s") {
		safeTiles34[tile] = true
	}
	leftTiles34 := InitLeftTiles34WithTiles34(MustStrToTiles34("11112222p 7777s"))
	risk34 := CalculateRiskTiles34(8, safeTiles34, leftTiles34, nil, 27, 28)
	for i, risk := range risk34 {
		t.Log(Mahjong[i], risk)
	}
}

func TestCalculateLeftNoSujiTiles(t *testing.T) {
	discardTiles34 := MustStrToTiles34("124689m 1346p 38s")
	leftTiles34 := InitLeftTiles34WithTiles34(MustStrToTiles34("22225555s"))

	safeTiles34 := make([]bool, 34)
	for i, c := range discardTiles34 {
		if c >= 1 {
			safeTiles34[i] = true
		}
	}
	assert.Equal(t, "28p 9s", TilesToStr(CalculateLeftNoSujiTiles(safeTiles34, leftTiles34)))
}
