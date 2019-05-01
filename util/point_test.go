package util

import (
	"testing"
	"fmt"
)

func TestCalcPointRon(t *testing.T) {
	t.Log(CalcPointRon(3, 40, 0, false) == 5200)
	t.Log(CalcPointRon(3, 40, 0, true) == 7700)
	t.Log(CalcPointRon(4, 40, 0, true) == 12000)
}

func TestCalcPointTsumoSum(t *testing.T) {
	t.Log(CalcPointTsumoSum(3, 40, 0, false) == 5200)
	t.Log(CalcPointTsumoSum(3, 40, 0, true) == 7800)
	t.Log(CalcPointTsumoSum(4, 40, 0, true) == 12000)
}

func TestCalcRonPointWithHands(t *testing.T) {
	// 子家默听荣和
	for _, tiles := range []string{
		//                          荣和点数
		"11m 112233445566z",     // 12000 [七对 混老头 混一色]
		"345m 345s 334455p 44z", // 7700 [平和 一杯口 三色]
		"333m 333s 333345p 11z", // 2600 [三色同刻]
		"22334455m 234s 234p",   // 8000：高点法取[一杯口 三色 断幺]  [22m 345m 345m 234p 234s][平和 一杯口 断幺], [55m 234m 234m 234p 234s][一杯口 三色 断幺]
		"234m 333p 55666777z",   // 12000 [三暗刻 役牌 役牌 小三元]
		"123445566789m 11z",     // 12000 [一杯口 一气 混一色]
		"123m 123999s 11155z",   // 3200 [混全]
		"334455m 667788s 77z",   // 5200 [两杯口]
		"334455m 667788s 44z",   // 7700 [平和 两杯口]
		"123m 123999s 11789p",   // 5200 [纯全]
	} {
		fmt.Print(tiles + " = ")
		tiles34 := MustStrToTiles34(tiles)
		ronPoint := CalcRonPointWithHands(&HandInfo{
			HandTiles34:   tiles34,
			IsTsumo:       false,
			WinTile:       MustStrToTile34("3m"),
			RoundWindTile: 28,
			SelfWindTile:  29,
		})
		fmt.Print(ronPoint)
		fmt.Println()
	}

	fmt.Println()

	// 子家立直荣和
	for _, tiles := range []string{
		//                          荣和点数
		"345m 222789p 333s 66z", // 1300 [立直]
	} {
		fmt.Print(tiles + " = ")
		tiles34 := MustStrToTiles34(tiles)
		ronPoint := CalcRonPointWithHands(&HandInfo{
			HandTiles34:   tiles34,
			IsTsumo:       false,
			WinTile:       MustStrToTile34("3m"),
			RoundWindTile: 28,
			SelfWindTile:  29,
			IsRiichi:      true,
		})
		fmt.Print(ronPoint)
		fmt.Println()
	}

	// 子家立直荣和，带宝牌
	for _, tiles := range []string{
		//                          荣和点数
		"345m 222789p 333s 66z", // 1300 2600 5200 8000 8000 12000 12000 16000 16000 16000 24000 24000 32000
	} {
		fmt.Print(tiles + " = ")
		tiles34 := MustStrToTiles34(tiles)
		for doraCount := 0; doraCount < 13; doraCount++ {
			ronPoint := CalcRonPointWithHands(&HandInfo{
				HandTiles34:   tiles34,
				IsTsumo:       false,
				WinTile:       MustStrToTile34("3m"),
				RoundWindTile: 28,
				SelfWindTile:  29,
				DoraCount:     doraCount,
				IsRiichi:      true,
			})
			fmt.Print(ronPoint)
			fmt.Print(" ")
		}
		fmt.Println()
	}

	// 亲家立直荣和，带宝牌
	for _, tiles := range []string{
		//                          荣和点数
		"345m 222789p 333s 66z", // 2000 3900 7700 12000 12000 18000 18000 24000 24000 24000 36000 36000 48000
	} {
		fmt.Print(tiles + " = ")
		tiles34 := MustStrToTiles34(tiles)
		for doraCount := 0; doraCount < 13; doraCount++ {
			ronPoint := CalcRonPointWithHands(&HandInfo{
				HandTiles34:   tiles34,
				IsTsumo:       false,
				WinTile:       MustStrToTile34("3m"),
				RoundWindTile: 28,
				SelfWindTile:  29,
				DoraCount:     doraCount,
				IsParent:      true,
				IsRiichi:      true,
			})
			fmt.Print(ronPoint)
			fmt.Print(" ")
		}
		fmt.Println()
	}
}
