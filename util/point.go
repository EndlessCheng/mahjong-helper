package util

import "github.com/EndlessCheng/mahjong-helper/util/model"

// TODO: 考虑大三元和大四喜的包牌？

func roundUpPoint(point int) int {
	if point == 0 {
		return 0
	}
	return ((point-1)/100 + 1) * 100
}

func calcBasicPoint(han int, fu int, yakumanTimes int) (basicPoint int) {
	switch {
	case yakumanTimes > 0:
		basicPoint = 8000 * yakumanTimes
	case han >= 13: // 累计役满
		basicPoint = 8000
	case han >= 11: // 三倍满
		basicPoint = 6000
	case han >= 8: // 倍满
		basicPoint = 4000
	case han >= 6: // 跳满
		basicPoint = 3000
	default:
		basicPoint = fu * (1 << uint(2+han))
		if basicPoint > 2000 { // 满贯
			basicPoint = 2000
		}
	}
	return
}

// 番数 符数 役满倍数 是否为亲家
// 返回荣和点数
func CalcPointRon(han int, fu int, yakumanTimes int, isParent bool) (point int) {
	basicPoint := calcBasicPoint(han, fu, yakumanTimes)
	if isParent {
		point = 6 * basicPoint
	} else {
		point = 4 * basicPoint
	}
	return roundUpPoint(point)
}

// 番数 符数 役满倍数 是否为亲家
// 返回自摸时的子家支付点数和亲家支付点数
func CalcPointTsumo(han int, fu int, yakumanTimes int, isParent bool) (childPoint int, parentPoint int) {
	basicPoint := calcBasicPoint(han, fu, yakumanTimes)
	if isParent {
		childPoint = 2 * basicPoint
	} else {
		childPoint = basicPoint
		parentPoint = 2 * basicPoint
	}
	return roundUpPoint(childPoint), roundUpPoint(parentPoint)
}

// 番数 符数 役满倍数 是否为亲家
// 返回自摸时的点数
func CalcPointTsumoSum(han int, fu int, yakumanTimes int, isParent bool) int {
	childPoint, parentPoint := CalcPointTsumo(han, fu, yakumanTimes, isParent)
	if isParent {
		return 3 * childPoint
	}
	return 2*childPoint + parentPoint
}

// TODO: 振听只能自摸

// 计算荣和点数
// 调用前请设置 WinTile
// 无役时返回 0
func CalcRonPointWithHands(playerInfo *model.PlayerInfo) (ronPoint int) {
	for _, result := range DivideTiles34(playerInfo.HandTiles34) {
		_hi := &_handInfo{
			PlayerInfo:   playerInfo,
			divideResult: result,
		}
		yakuTypes := findYakuTypes(_hi)
		if len(yakuTypes) == 0 {
			continue
		}
		han := CalcYakuHan(yakuTypes, _hi.isNaki())
		han += _hi.DoraCount
		fu := _hi.calcFu()
		yakumanTimes := CalcYakumanTimes()
		point := CalcPointRon(han, fu, yakumanTimes, _hi.IsParent)
		// 高点法
		ronPoint = MaxInt(ronPoint, point)
	}
	return
}

// TODO: 计算自摸点数

// TODO: 考虑里宝时的荣和、自摸点数
