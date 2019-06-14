package util

import (
	"testing"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"github.com/stretchr/testify/assert"
)

func Test_handInfo_calcFu(t *testing.T) {
	assert := assert.New(t)

	fuList := func(humanTiles string, humanWinTile string, isTsumo bool) (l []int) {
		l = []int{}
		playerInfo := MustParseHumanTilesWithMelds(humanTiles)
		playerInfo.WinTile = MustStrToTile34(humanWinTile)
		playerInfo.IsTsumo = isTsumo
		results := DivideTiles34(playerInfo.HandTiles34)
		if len(results) == 0 {
			return
		}
		isNaki := playerInfo.IsNaki()
		for _, result := range results {
			_hi := &_handInfo{
				PlayerInfo:   playerInfo,
				divideResult: result,
			}
			l = append(l, _hi.calcFu(isNaki))
		}
		return
	}

	// 七对
	assert.EqualValues([]int{25}, fuList("33m 112233445566z", "3m", true))
	assert.EqualValues([]int{25}, fuList("33m 112233445566z", "3m", false))

	// 平和
	assert.EqualValues([]int{20}, fuList("345m 345s 334455p 44z", "3m", true))
	assert.EqualValues([]int{20}, fuList("33345m 345s 345789p", "3m", true))
	assert.EqualValues([]int{30}, fuList("345m 345s 334455p 44z", "3m", false))
	assert.EqualValues([]int{30}, fuList("33345m 345s 345789p", "3m", false))

	// 20 + 中张暗刻4 + 连风4 + 自摸2
	assert.EqualValues([]int{30}, fuList("345m 222s 334455p 11z", "3m", true))

	// 20 + 中张暗刻4 + 连风4 + 自摸2 + 坎张2（刚好跳符）
	assert.EqualValues([]int{40}, fuList("234m 222s 334455p 11z", "3m", true))

	// 20 + 中张暗刻4 + 连风4 + 自摸2 + 边张2（刚好跳符）
	assert.EqualValues([]int{40}, fuList("123m 222s 334455p 11z", "3m", true))

	// 平和 / 坎张
	assert.EqualValues([]int{20, 30}, fuList("22334455m 234s 234p", "3m", true))
	assert.EqualValues([]int{30, 40}, fuList("22334455m 234s 234p", "3m", false))

	// 副露平和型
	assert.EqualValues([]int{30}, fuList("345m 345s 345p 44z # 567m", "3m", true))
	assert.EqualValues([]int{30}, fuList("345m 345s 345p 44z # 567m", "3m", false))

	// 副露坎张
	assert.EqualValues([]int{30}, fuList("345m 345s 345p 44z # 567m", "3m", true))
	assert.EqualValues([]int{30}, fuList("345m 345s 345p 44z # 567m", "3m", false))

	// 暗杠
	assert.EqualValues([]int{70}, fuList("123456m 678s 11p # 1111S", "6s", false))
	assert.EqualValues([]int{60}, fuList("123m 678s 11p # 456m 1111S", "6s", false))

	// 明杠 + 明刻 + 暗刻
	assert.EqualValues([]int{50}, fuList("123p 66677z # 8888p 999m", "3p", true))

	// (1番)110 符
	assert.EqualValues([]int{110}, fuList("234s 11777z # 1111M 9999P", "7z", false))

	// 青天井规则下的 170 符（四杠子 四暗刻单骑）
	assert.EqualValues([]int{170}, fuList("77z # 1111S 9999P 1111M 9999M", "7z", false))
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
		// 17.5 ns/op
		_hi.calcFu(false)
	}
}
