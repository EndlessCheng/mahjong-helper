package util

import "testing"

func TestCalculateShantenWithImprove(t *testing.T) {
	t.Log(CalculateShantenAndWaits13(MustStrToTiles34("11357m 13579p 135s"), false))
}

func TestCalculateShantenWithImproves13(t *testing.T) {
	tiles := "334m 22478p 23456s"
	t.Log(CalculateShantenWithImproves13(MustStrToTiles34(tiles), false))
}

func TestCalculateShantenWithImproves14(t *testing.T) {
	tiles := "124679m 3678p 2366s"
	tiles = "11379m 347p 277s 111z"
	tiles = "334578m 11468p 235s"
	tiles = "478m 33588p 457899s"
	tiles = "2233688m 1234p 379s"
	tiles = "1233347m 23699p 88s"
	tiles = "334m 22457p 23456s 1z"
	tiles = "334m 122478p 23456s"
	tiles = "1m 258p 258s 1234567z"
	shanten, results, _ := CalculateShantenWithImproves14(MustStrToTiles34(tiles), false)
	t.Log(shanten)
	for _, result := range results {
		t.Log(result)
	}
}

func BenchmarkCalculateShantenWithImproves14(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 剪枝前：1.91s
		// 剪枝后：1.19s
		CalculateShantenWithImproves14(MustStrToTiles34("124679m 3678p 2366s"), false)
	}
}
