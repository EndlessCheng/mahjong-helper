package util

import "testing"

func TestOutsideTiles(t *testing.T) {
	for i := 0; i < 34; i++ {
		t.Log(TilesToStr(OutsideTiles(i)))
	}
}
