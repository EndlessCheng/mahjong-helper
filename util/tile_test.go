package util

import "testing"

func TestTilesToMergedStr(t *testing.T) {
	t.Log(TilesToMergedStr([]int{0, 2, 9}) == "13m 1p")
	t.Log(TilesToMergedStr([]int{0, 2, 3}) == "134m")
	t.Log(TilesToMergedStr([]int{0, 9, 18, 27, 33}) == "1m 1p 1s 17z")
	t.Log(TilesToMergedStr([]int{0, 8, 9, 17, 18, 26, 27, 33}) == "19m 19p 19s 17z")
}
