package util

import "testing"

func TestCalcWallTiles34(t *testing.T) {
	for _, tiles := range []string{
		"2222777888m",
		"33337777m",
		"333777m",
		"333444777m",
	} {
		t.Log(CalcWallTiles34(invert(MustStrToTiles34(tiles))))
	}
}
