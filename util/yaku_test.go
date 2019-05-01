package util

import (
	"testing"
	"fmt"
)

func TestFindNormalYaku(t *testing.T) {
	for _, tiles := range []string{
		"11m 112233445566z",     // [七对 混老头 混一色]
		"345m 345s 334455p 44z", // [平和 一杯口 三色]
		"222m 222s 222345p 11z", // [三暗刻 三色同刻]
		"22334455m 234s 234p",   // [22m 345m 345m 234p 234s][平和 一杯口 断幺], [55m 234m 234m 234p 234s][一杯口 三色 断幺]
	} {
		fmt.Print(tiles + " = ")
		tiles34 := MustStrToTiles34(tiles)
		for _, result := range DivideTiles34(tiles34) {
			yakuTypes := findYakuTypes(&_handInfo{
				HandInfo: &HandInfo{
					HandTiles34:   tiles34,
					IsTsumo:       true,
					WinTile:       MustStrToTile34("3m"),
					RoundWindTile: 27,
					SelfWindTile:  27,
				},
				divideResult: result,
			})
			fmt.Printf("%s %v, ", result.String(), YakuTypesToStr(yakuTypes))
		}
		fmt.Println()
	}
}
