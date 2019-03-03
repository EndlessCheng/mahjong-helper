package util

import (
	"testing"
)

func TestCheckWin(t *testing.T) {
	for _, tiles := range []string{
		"123456789m 12344s",
		"111234678m 11122z",
	} {
		tiles34 := MustStrToTiles34(tiles)
		if !CheckWin(tiles34) {
			t.Error("CheckWin failed at", tiles)
		}
	}
}

func BenchmarkCheckWin(b *testing.B) {
	// 92.7 ns/op
	tiles34 := MustStrToTiles34("123456789m 12344s")
	for i := 0; i < b.N; i++ {
		CheckWin(tiles34)
	}
}
