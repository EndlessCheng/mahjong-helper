package util

import (
	"testing"
	"fmt"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"github.com/stretchr/testify/assert"
)

func TestCalcPointRon(t *testing.T) {
	assert.Equal(t, 5200, CalcPointRon(3, 40, 0, false))
	assert.Equal(t, 7700, CalcPointRon(3, 40, 0, true))
	assert.Equal(t, 8000, CalcPointRon(3, 70, 0, false))
	assert.Equal(t, 12000, CalcPointRon(4, 40, 0, true))
	assert.Equal(t, 32000, CalcPointRon(0, 0, 1, false))
	assert.Equal(t, 64000, CalcPointRon(0, 0, 2, false))
	assert.Equal(t, 96000, CalcPointRon(0, 0, 3, false))
	assert.Equal(t, 128000, CalcPointRon(0, 0, 4, false))
}

func TestCalcPointTsumoSum(t *testing.T) {
	assert.Equal(t, 5200, CalcPointTsumoSum(3, 40, 0, false))
	assert.Equal(t, 7800, CalcPointTsumoSum(3, 40, 0, true))
	assert.Equal(t, 12000, CalcPointTsumoSum(4, 40, 0, true))
}

func TestCalcRonPointWithHands(t *testing.T) {
	// 子家默听荣和
	newPI := func(humanTiles string, winHumanTile string) *model.PlayerInfo {
		return &model.PlayerInfo{
			HandTiles34:   MustStrToTiles34(humanTiles),
			WinTile:       MustStrToTile34(winHumanTile),
			RoundWindTile: MustStrToTile34("2z"),
			SelfWindTile:  MustStrToTile34("2z"),
		}
	}
	assert.Equal(t, 12000, CalcRonPoint(newPI("11m 112233445566z", "1m")))    // [七对 混老头 混一色]
	assert.Equal(t, 7700, CalcRonPoint(newPI("345m 345s 334455p 44z", "3m"))) // [平和 一杯口 三色]
	assert.Equal(t, 2600, CalcRonPoint(newPI("333m 333s 333345p 11z", "3m"))) // [三色同刻]
	assert.Equal(t, 8000, CalcRonPoint(newPI("22334455m 234s 234p", "3m")))   // 高点法取[一杯口 三色 断幺]
	assert.Equal(t, 12000, CalcRonPoint(newPI("234m 333p 55666777z", "3m")))  // [三暗刻 役牌 役牌 小三元]
	assert.Equal(t, 12000, CalcRonPoint(newPI("123445566789m 11z", "3m")))    // [一杯口 一气 混一色]
	assert.Equal(t, 3200, CalcRonPoint(newPI("123m 123999s 11155z", "3m")))   // [混全]
	assert.Equal(t, 5200, CalcRonPoint(newPI("334455m 667788s 77z", "3m")))   // [两杯口]
	assert.Equal(t, 7700, CalcRonPoint(newPI("334455m 667788s 44z", "3m")))   // [平和 两杯口]
	assert.Equal(t, 5200, CalcRonPoint(newPI("123m 123999s 11789p", "3m")))   // [纯全]
	assert.Equal(t, 2600, CalcRonPoint(newPI("345m 12355789s 222z", "3m")))   // [役牌 役牌]

	fmt.Println()

	// 子家立直荣和
	newPIWithRiichi := func(humanTiles string, winHumanTile string) *model.PlayerInfo {
		return &model.PlayerInfo{
			HandTiles34:   MustStrToTiles34(humanTiles),
			WinTile:       MustStrToTile34(winHumanTile),
			RoundWindTile: MustStrToTile34("2z"),
			SelfWindTile:  MustStrToTile34("3z"),
			IsRiichi:      true,
		}
	}
	assert.Equal(t, 1300, CalcRonPoint(newPIWithRiichi("345m 222789p 333s 66z", "3m"))) // [立直]

	// 子家立直荣和，带宝牌
	ronPoints := []int{}
	for doraCount := 0; doraCount < 13; doraCount++ {
		ronPoint := CalcRonPoint(&model.PlayerInfo{
			NumRedFives:   []int{doraCount, 0, 0}, // 方便算番
			HandTiles34:   MustStrToTiles34("345m 222789p 333s 66z"),
			WinTile:       MustStrToTile34("3m"),
			RoundWindTile: MustStrToTile34("2z"),
			SelfWindTile:  MustStrToTile34("3z"),
			IsRiichi:      true,
		})
		ronPoints = append(ronPoints, ronPoint)
	}
	assert.Equal(t, ronPoints, []int{1300, 2600, 5200, 8000, 8000, 12000, 12000, 16000, 16000, 16000, 24000, 24000, 32000})

	// 亲家立直荣和，带宝牌
	ronPoints = []int{}
	for doraCount := 0; doraCount < 13; doraCount++ {
		ronPoint := CalcRonPoint(&model.PlayerInfo{
			NumRedFives:   []int{doraCount, 0, 0}, // 方便算番
			HandTiles34:   MustStrToTiles34("345m 222789p 333s 66z"),
			WinTile:       MustStrToTile34("3m"),
			RoundWindTile: MustStrToTile34("2z"),
			SelfWindTile:  MustStrToTile34("3z"),
			IsParent:      true,
			IsRiichi:      true,
		})
		ronPoints = append(ronPoints, ronPoint)
	}
	assert.Equal(t, ronPoints, []int{2000, 3900, 7700, 12000, 12000, 18000, 18000, 24000, 24000, 24000, 36000, 36000, 48000})
}
