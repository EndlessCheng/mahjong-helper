package util

import "testing"

func TestTilesToMergedStr(t *testing.T) {
	t.Log(TilesToStr([]int{0, 2, 9}) == "13m 1p")
	t.Log(TilesToStr([]int{0, 2, 3}) == "134m")
	t.Log(TilesToStr([]int{0, 9, 18, 27, 33}) == "1m 1p 1s 17z")
	t.Log(TilesToStr([]int{0, 8, 9, 17, 18, 26, 27, 33}) == "19m 19p 19s 17z")
}

func TestTiles34ToMergedStr(t *testing.T) {
	tiles := "1119m 1999p 19s 17z"
	tiles34 := MustStrToTiles34(tiles)
	t.Log(Tiles34ToStr(tiles34) == tiles)
}
