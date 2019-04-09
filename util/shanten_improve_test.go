package util

import (
	"testing"
)

func TestCalculateShantenWithImproveClosed(t *testing.T) {
	for _, tiles := range []string{
		"11357m 13579p 135s",
		"123456789m 1135s",
		"123456789m 1134s",
		"123456789m 1234z",
	} {
		tiles34 := MustStrToTiles34(tiles)
		if CountOfTiles34(tiles34) != 13 {
			t.Error(tiles, "不是13张牌")
			continue
		}
		shanten, waits, _ := CalculateShantenAndWaits13(tiles34, false)
		t.Log(tiles, "=", NumberToChineseShanten(shanten), waits)
	}
}

func TestCalculateShantenWithImproveOpen(t *testing.T) {
	for _, tiles := range []string{
		"1234p",
		"1234z",
		"5p",
		"66m 12334p 334s 777z",
		"66m 123334p 34s 777z",
	} {
		tiles34 := MustStrToTiles34(tiles)
		shanten, waits, meldWaits := CalculateShantenAndWaits13(tiles34, true)
		t.Log(tiles, "=", NumberToChineseShanten(shanten), waits, "+", meldWaits)
	}
}

func TestCalculateShantenWithImproves13Closed(t *testing.T) {
	for _, tiles := range []string{
		"11357m 13579p 135s",
		"123456789m 1135s",
		"123456789m 1134s",
		"123456789m 1234z",
		"3m 12668p 5678s 222z",
		"6m 12668p 5678s 222z",
	} {
		tiles34 := MustStrToTiles34(tiles)
		if CountOfTiles34(tiles34) != 13 {
			t.Error(tiles, "不是13张牌")
			continue
		}
		result := CalculateShantenWithImproves13(tiles34, false)
		t.Log(tiles, "=\n"+result.String())
	}
}

func TestCalculateShantenWithImproves13Open(t *testing.T) {
	for _, tiles := range []string{
		"1234m",
		"1135m",
		"5p",
	} {
		tiles34 := MustStrToTiles34(tiles)
		result := CalculateShantenWithImproves13(tiles34, true)
		t.Log(tiles, "=\n"+result.String())
	}
}

func TestCalculateShantenWithImproves14Closed(t *testing.T) {
	tiles := "124679m 3678p 2366s"
	tiles = "11379m 347p 277s 111z"
	tiles = "334578m 11468p 235s"
	tiles = "478m 33588p 457899s"
	tiles = "2233688m 1234p 379s"
	tiles = "1233347m 23699p 88s"
	tiles = "334m 22457p 23456s 1z"
	tiles = "334m 122478p 23456s"
	tiles = "1m 258p 258s 1234567z"
	tiles = "4567m 4579p 344588s"
	tiles = "2479999m 45667p 13s" // 切任何一张都不会向听倒退
	tiles = "25667m 27789p 37s 44z"
	tiles = "111444777m 11177s"
	tiles = "2468m 33578p 22356s"
	tiles = "57m 4455p 12345699s"
	tiles = "57m 3445667p 12399s"
	shanten, results, incShantenResults := CalculateShantenWithImproves14(MustStrToTiles34(tiles), false)
	t.Log(NumberToChineseShanten(shanten))
	for _, result := range results {
		t.Log(result)
	}
	if len(incShantenResults) > 0 {
		t.Log(NumberToChineseShanten(shanten + 1))
		for _, result := range incShantenResults {
			t.Log(result)
		}
	} else {
		t.Log("无向听倒退的切牌")
	}
}

func BenchmarkCalculateShantenWithImproves14Closed(b *testing.B) {
	tiles34 := MustStrToTiles34("124679m 3678p 2366s")
	for i := 0; i < b.N; i++ {
		// 剪枝前：0.28s
		// 剪枝后：0.22s
		CalculateShantenWithImproves14(tiles34, false)
	}
}

func TestCalculateShantenWithImproves14Open(t *testing.T) {
	tiles := "35m"
	shanten, results, incShantenResults := CalculateShantenWithImproves14(MustStrToTiles34(tiles), true)
	t.Log(NumberToChineseShanten(shanten))
	for _, result := range results {
		t.Log(result)
	}
	t.Log(NumberToChineseShanten(shanten + 1))
	for _, result := range incShantenResults {
		t.Log(result)
	}
}

func TestCalculateMeld(t *testing.T) {
	tiles := "1234m 112z"
	tiles = "23445667m 11z"
	tiles = "112356799m 1233z"
	tiles = "78m 12355p 789s" // ***
	result := CalculateShantenWithImproves13(MustStrToTiles34(tiles), true)
	t.Log("原手牌" + NumberToChineseShanten(result.Shanten))
	t.Log(result)

	tile := "1m"
	tile = "3m" // "1z"
	tile = "4m"
	tile = "4p"
	shanten, results, incShantenResults := CalculateMeld(MustStrToTiles34(tiles), MustStrToTile34(tile), true)
	t.Log("鸣牌后" + NumberToChineseShanten(shanten))
	for _, result := range results {
		t.Log(result)
	}
	t.Log("鸣牌后" + NumberToChineseShanten(shanten+1))
	for _, result := range incShantenResults {
		t.Log(result)
	}
}
