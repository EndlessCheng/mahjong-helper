package util

import (
	"testing"
	"fmt"
	"github.com/EndlessCheng/mahjong-helper/util/model"
)

// TODO: assert slice
func Test_handInfo_calcFu(t *testing.T) {
	// 自摸
	for _, tiles := range []string{
		"33m 112233445566z",     // 七对自摸25符
		"345m 345s 334455p 44z", // 平和自摸20符
		"33345m 345s 345789p",   // 平和自摸20符
		"345m 222s 334455p 11z", // 20 + 中张暗刻4 + 连风4 + 自摸2 = 30符
		"234m 222s 334455p 11z", // 20 + 坎张2 + 中张暗刻4 + 连风4 + 自摸2 = 32符 = 进40符
		"22334455m 234s 234p",   // 20 符 / 30 符
		"678s 11m 123345p 666z",
	} {
		fmt.Print(tiles + " = ")
		tiles34 := MustStrToTiles34(tiles)
		results := DivideTiles34(tiles34)
		if len(results) == 0 {
			fmt.Println("[国士/未和牌]")
			continue
		}
		for _, result := range results {
			_hi := &_handInfo{
				PlayerInfo: &model.PlayerInfo{
					HandTiles34:   tiles34,
					IsTsumo:       true,
					WinTile:       MustStrToTile34("3m"),
					RoundWindTile: 27,
					SelfWindTile:  27,
				},
				divideResult: result,
			}
			fmt.Printf("%s %d 符, ", result.String(), _hi.calcFu(false))
		}
		fmt.Println()
	}

	fmt.Println()

	// 荣和
	for _, tiles := range []string{
		"345m 345s 334455p 44z", // 20 + 门清10符 = 30符
		"345m 222s 334455p 44z", // 20 + 门清10符 + 中张暗刻4 = 34符 = 进40符
	} {
		fmt.Print(tiles + " = ")
		tiles34 := MustStrToTiles34(tiles)
		results := DivideTiles34(tiles34)
		if len(results) == 0 {
			fmt.Println("[国士/未和牌]")
			continue
		}
		for _, result := range results {
			_hi := &_handInfo{
				PlayerInfo: &model.PlayerInfo{
					HandTiles34:   tiles34,
					IsTsumo:       false,
					WinTile:       MustStrToTile34("3m"),
					RoundWindTile: 27,
					SelfWindTile:  27,
				},
				divideResult: result,
			}
			fmt.Printf("%s %d 符, ", result.String(), _hi.calcFu(false))
		}
		fmt.Println()
	}

	fmt.Println()

	// 副露平和荣和
	for _, tiles := range []string{
		"345m 345s 345p 44z", // 副露平和30符
	} {
		fmt.Print(tiles + " = ")
		tiles34 := MustStrToTiles34(tiles)
		results := DivideTiles34(tiles34)
		if len(results) == 0 {
			fmt.Println("[国士/未和牌]")
			continue
		}
		for _, result := range results {
			_hi := &_handInfo{
				PlayerInfo: &model.PlayerInfo{
					HandTiles34:   tiles34,
					Melds:         []model.Meld{{MeldType: model.MeldTypeChi}},
					IsTsumo:       false,
					WinTile:       MustStrToTile34("3m"),
					RoundWindTile: 27,
					SelfWindTile:  27,
				},
				divideResult: result,
			}
			fmt.Printf("%s %d 符, ", result.String(), _hi.calcFu(true))
		}
		fmt.Println()
	}
}

func Benchmark_handInfo_calcFu(b *testing.B) {
	tiles34 := MustStrToTiles34("111234678m 11122z")
	results := DivideTiles34(tiles34)
	_hi := &_handInfo{
		PlayerInfo: &model.PlayerInfo{
			HandTiles34:   tiles34,
			Melds:         nil,
			IsTsumo:       true,
			WinTile:       MustStrToTile34("3m"),
			RoundWindTile: MustStrToTile34("2z"),
			SelfWindTile:  MustStrToTile34("2z"),
		},
		divideResult: results[0],
	}
	for i := 0; i < b.N; i++ {
		// 19.5 ns/op
		_hi.calcFu(false)
	}
}
