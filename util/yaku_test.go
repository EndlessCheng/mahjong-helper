package util

import (
	"testing"
	"fmt"
	"github.com/EndlessCheng/mahjong-helper/util/model"
)

func Test_findYakuTypes(t *testing.T) {
	for _, tiles := range []string{
		"11m 112233445566z",     // [七对 混老头 混一色]
		"345m 345s 334455p 44z", // [平和 一杯口 三色]
		"333m 333s 333345p 11z", // [三暗刻 三色同刻]
		"22334455m 234s 234p",   // [22m 345m 345m 234p 234s][平和 一杯口 断幺], [55m 234m 234m 234p 234s][一杯口 三色 断幺]
		"234m 333p 55666777z",   // [三暗刻 役牌 役牌 小三元]
		"123445566789m 11z",     // [一杯口 一气 混一色]
		"111222333444m 11z",     // [11z 111m 222m 333m 444m][四暗刻 混一色], [11z 444m 123m 123m 123m][一杯口 混一色]
		"11122345678999m",       // [九莲宝灯]
		"123m 123999s 11155z",   // [役牌 混全]
		"334455m 667788s 77z",   // [两杯口]
		"334455m 667788s 44z",   // [平和 两杯口]
		"123m 123999s 11789p",   // [纯全]
		//"111999m 111p 11122z",   // [四暗刻]([对对 役牌 混老头])
	} {
		fmt.Print(tiles + " = ")
		tiles34 := MustStrToTiles34(tiles)
		for _, result := range DivideTiles34(tiles34) {
			yakuTypes := findYakuTypes(&_handInfo{
				PlayerInfo: &model.PlayerInfo{
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

	fmt.Println()

	// 三暗刻的荣和判定
	for _, tiles := range []string{
		"333m 333p 333567s 11z", // [三色同刻]
		"333345m 333p 333s 11z", // [三暗刻 三色同刻]
	} {
		fmt.Print(tiles + " = ")
		tiles34 := MustStrToTiles34(tiles)
		for _, result := range DivideTiles34(tiles34) {
			yakuTypes := findYakuTypes(&_handInfo{
				PlayerInfo: &model.PlayerInfo{
					HandTiles34:   tiles34,
					IsTsumo:       false,
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

	fmt.Println()

	for _, tiles := range []string{
		"333m 77z", // [对对 三杠子 役牌 混一色]
	} {
		fmt.Print(tiles + " = ")
		tiles34 := MustStrToTiles34(tiles)
		for _, result := range DivideTiles34(tiles34) {
			yakuTypes := findYakuTypes(&_handInfo{
				PlayerInfo: &model.PlayerInfo{
					HandTiles34: tiles34,
					Melds: []model.Meld{
						{MeldType: model.MeldTypeMinkan, Tiles: MustStrToTiles("1111z")},
						{MeldType: model.MeldTypeMinkan, Tiles: MustStrToTiles("2222z")},
						{MeldType: model.MeldTypeMinkan, Tiles: MustStrToTiles("3333z")},
					},
					IsTsumo:       false,
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

	fmt.Println()

	// 无役检测
	for _, tiles := range []string{
		"333m 123s 123p 77z", // [无役]
	} {
		fmt.Print(tiles + " = ")
		tiles34 := MustStrToTiles34(tiles)
		for _, result := range DivideTiles34(tiles34) {
			yakuTypes := findYakuTypes(&_handInfo{
				PlayerInfo: &model.PlayerInfo{
					HandTiles34: tiles34,
					Melds: []model.Meld{
						{MeldType: model.MeldTypeChi, Tiles: MustStrToTiles("789p")},
					},
					IsTsumo:       false,
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
