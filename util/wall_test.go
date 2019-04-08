package util

import "testing"

func TestCalcWallTiles34(t *testing.T) {
	for _, tiles := range []string{
		"2222777888m",
		"33337777m",
		"333777m",
		"333444777m",
		"8888m",
	} {
		t.Log(CalcWallTiles(invert(MustStrToTiles34(tiles))))
	}
}

func TestCalcNCSafeTiles34(t *testing.T) {
	for _, tiles := range []string{
		"8888m",
	} {
		leftTiles34 := invert(MustStrToTiles34(tiles))
		t.Log(CalcNCSafeTiles(leftTiles34).FilterWithHands(MustStrToTiles34("9m")))
	}
}

func TestCalcDNCSafeTiles(t *testing.T) {
	for _, tiles := range []string{
		"8888m",
		"33336666m",
		"33335555m",
		"33337777m",
		"333355557777m",
	} {
		leftTiles34 := invert(MustStrToTiles34(tiles))
		t.Log(CalcDNCSafeTiles(leftTiles34))
	}
}
