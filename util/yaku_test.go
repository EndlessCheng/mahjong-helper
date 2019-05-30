package util

import (
	"testing"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"strings"
	"github.com/stretchr/testify/assert"
)

func Test_findYakuTypes(t *testing.T) {
	assert := assert.New(t)

	calcStrYaku := func(humanTiles string, humanWinTile string, isTsumo bool, melds ...model.Meld) string {
		output := ""
		pi := &model.PlayerInfo{
			HandTiles34:   MustStrToTiles34(humanTiles),
			Melds:         melds,
			IsTsumo:       isTsumo,
			WinTile:       MustStrToTile34(humanWinTile),
			RoundWindTile: 27,
			SelfWindTile:  27,
		}
		for _, result := range DivideTiles34(pi.HandTiles34) {
			yakuTypes := findYakuTypes(&_handInfo{
				PlayerInfo:   pi,
				divideResult: result,
			})
			output += YakuTypesToStr(yakuTypes) + " "
		}
		return strings.TrimSpace(output)
	}

	assert.Equal("[七对 混老头 混一色]", calcStrYaku("99s 112233445566z", "9s", false))
	assert.Equal("[七对 混一色]", calcStrYaku("22m 112233445566z", "2m", false))
	assert.Equal("[平和 一杯口 三色]", calcStrYaku("345m 345s 334455p 44z", "3m", false))
	assert.Equal("[三色同刻]", calcStrYaku("333m 333s 333345p 11z", "3m", false))
	assert.Equal("[平和 一杯口 断幺] [一杯口 三色 断幺]", calcStrYaku("22334455m 234s 234p", "3m", false))
	assert.Equal("[三暗刻 役牌 役牌 小三元]", calcStrYaku("234m 333p 55666777z", "3m", false))
	assert.Equal("[一杯口 一通 混一色]", calcStrYaku("123445566789m 11z", "3m", false))
	assert.Equal("[对对 三暗刻 混一色] [一杯口 混一色]", calcStrYaku("111222333444m 11z", "3m", false))
	assert.Equal("[四暗刻] [自摸 一杯口 混一色]", calcStrYaku("111222333444m 11z", "3m", true))
	assert.Equal("[役牌 役牌 混全]", calcStrYaku("123m 123999s 11155z", "3m", false))
	assert.Equal("[两杯口]", calcStrYaku("334455m 667788s 77z", "3m", false))
	assert.Equal("[平和 两杯口]", calcStrYaku("334455m 667788s 44z", "3m", false))
	assert.Equal("[纯全]", calcStrYaku("123m 123999s 11789p", "3m", false))

	// 役满
	assert.Equal("[九莲]", calcStrYaku("11122345678999m", "3m", false))
	assert.Equal("[纯正九莲]", calcStrYaku("11123345678999m", "3m", false))
	assert.Equal("[绿一色]", calcStrYaku("22334466688s 666z", "6z", false))
	assert.Equal("[四暗刻]", calcStrYaku("111999m 111p 11122z", "1z", true))
	assert.Equal("[小四喜 字一色]", calcStrYaku("11122233344555z", "1z", false))
	assert.Equal("[字一色]", calcStrYaku("11223344556677z", "1z", false))
	assert.Equal("[四暗刻单骑 大四喜 字一色]", calcStrYaku("11122233344455z", "5z", false))
	assert.Equal("[大三元]", calcStrYaku("12333m 555666777z", "1m", false))
	assert.Equal("[清老头]", calcStrYaku("111999m 111999s 11p", "1m", false))

	// 三暗刻判定
	assert.Equal("[三色同刻]", calcStrYaku("333m 333p 333567s 11z", "3m", false))
	assert.Equal("[三暗刻 三色同刻]", calcStrYaku("333345m 333p 333s 11z", "3m", false))

	// 副露相关
	assert.Equal("[一通 役牌 役牌 混一色]", calcStrYaku("123p 11177z", "3p", false,
		model.Meld{MeldType: model.MeldTypeChi, Tiles: MustStrToTiles("456p")},
		model.Meld{MeldType: model.MeldTypeChi, Tiles: MustStrToTiles("789p")},
	))
	assert.Equal("[对对 役牌 役牌 混老头]", calcStrYaku("111p 11177z", "1p", false,
		model.Meld{MeldType: model.MeldTypePon, Tiles: MustStrToTiles("999p")},
		model.Meld{MeldType: model.MeldTypePon, Tiles: MustStrToTiles("111s")},
	))
	assert.Equal("[对对 三杠子 混一色]", calcStrYaku("333m 77z", "3m", false,
		model.Meld{MeldType: model.MeldTypeMinkan, Tiles: MustStrToTiles("4444z")},
		model.Meld{MeldType: model.MeldTypeMinkan, Tiles: MustStrToTiles("2222z")},
		model.Meld{MeldType: model.MeldTypeMinkan, Tiles: MustStrToTiles("3333z")},
	))
	assert.Equal("[对对 三杠子 断幺]", calcStrYaku("333m 77s", "3m", false,
		model.Meld{MeldType: model.MeldTypeMinkan, Tiles: MustStrToTiles("4444s")},
		model.Meld{MeldType: model.MeldTypeMinkan, Tiles: MustStrToTiles("2222s")},
		model.Meld{MeldType: model.MeldTypeMinkan, Tiles: MustStrToTiles("3333s")},
	))
	assert.Equal("[四杠子]", calcStrYaku("77z", "7z", false,
		model.Meld{MeldType: model.MeldTypeMinkan, Tiles: MustStrToTiles("1111z")},
		model.Meld{MeldType: model.MeldTypeAnkan, Tiles: MustStrToTiles("1111p")},
		model.Meld{MeldType: model.MeldTypeKakan, Tiles: MustStrToTiles("2222z")},
		model.Meld{MeldType: model.MeldTypeMinkan, Tiles: MustStrToTiles("3333z")},
	))
	assert.Equal("[四暗刻单骑 大四喜 字一色 四杠子]", calcStrYaku("77z", "7z", false,
		model.Meld{MeldType: model.MeldTypeAnkan, Tiles: MustStrToTiles("1111z")},
		model.Meld{MeldType: model.MeldTypeAnkan, Tiles: MustStrToTiles("2222z")},
		model.Meld{MeldType: model.MeldTypeAnkan, Tiles: MustStrToTiles("3333z")},
		model.Meld{MeldType: model.MeldTypeAnkan, Tiles: MustStrToTiles("4444z")},
	))

	// 无役
	assert.Equal("[无役]", calcStrYaku("333m 123s 123p 77z", "3m", false,
		model.Meld{MeldType: model.MeldTypeChi, Tiles: MustStrToTiles("789p")},
	))
}

// 寻找所有可能的役种
// 调用前请设置 WinTile
func _findAllYakuTypes(playerInfo *model.PlayerInfo) (yakuTypes []int) {
	canYaku := make([]bool, maxYakuType)
	for _, result := range DivideTiles34(playerInfo.HandTiles34) {
		_hi := &_handInfo{
			PlayerInfo:   playerInfo,
			divideResult: result,
		}
		yakuTypes := findYakuTypes(_hi)
		for _, t := range yakuTypes {
			canYaku[t] = true
		}
	}
	for yakuType, isYaku := range canYaku {
		if isYaku {
			yakuTypes = append(yakuTypes, yakuType)
		}
	}
	return
}

func BenchmarkFindAllYakuTypes(b *testing.B) {
	pi := &model.PlayerInfo{
		HandTiles34:   MustStrToTiles34("345m 345789p 34555s"),
		IsTsumo:       false,
		WinTile:       MustStrToTile34("5s"),
		RoundWindTile: 27,
		SelfWindTile:  27,
	}
	for i := 0; i < b.N; i++ {
		// 1746 ns/op
		_findAllYakuTypes(pi)
	}
}
