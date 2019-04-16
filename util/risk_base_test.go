package util

import "testing"

func TestCalculateRiskTiles34(t *testing.T) {

}

func TestCalculateLeftNoSujiTiles(t *testing.T) {
	discardTiles34 := MustStrToTiles34("124689m 1346p 38s")
	leftTiles34 := InitLeftTiles34WithTiles34(MustStrToTiles34("22225555s"))

	safeTiles34 := make([]bool, 34)
	for i, c := range discardTiles34 {
		if c >=1 {
			safeTiles34[i] = true
		}
	}
	leftNoSujiTiles := CalculateLeftNoSujiTiles(safeTiles34, leftTiles34)
	t.Log(len(leftNoSujiTiles), "无筋:", TilesToStr(leftNoSujiTiles))
}
