package util

import (
	"github.com/EndlessCheng/mahjong-helper/util/model"
)

// TODO: 考虑大三元和大四喜的包牌？

func roundUpPoint(point int) int {
	if point == 0 {
		return 0
	}
	return ((point-1)/100 + 1) * 100
}

func calcBasicPoint(han int, fu int, yakumanTimes int) (basicPoint int) {
	switch {
	case yakumanTimes > 0: // (x倍)役满
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

//

type PointResult struct {
	han          int
	fu           int
	yakumanTimes int
	isParent     bool
	Point        int
}

// 已和牌，计算自摸或荣和时的点数（不考虑里宝、一发等情况）
// 无役时返回的 Point == 0
// 调用前请设置 IsTsumo WinTile
func CalcPoint(playerInfo *model.PlayerInfo) (result *PointResult) {
	result = &PointResult{}
	for _, divideResult := range DivideTiles34(playerInfo.HandTiles34) {
		_hi := &_handInfo{
			PlayerInfo:   playerInfo,
			divideResult: divideResult,
		}
		yakuTypes := findYakuTypes(_hi)
		if len(yakuTypes) == 0 {
			// 此手牌拆解下无役
			continue
		}
		han := CalcYakuHan(yakuTypes, _hi.IsNaki())
		han += _hi.CountDora()
		fu := _hi.calcFu()
		yakumanTimes := CalcYakumanTimes(yakuTypes)
		var pt int
		if _hi.IsTsumo {
			pt = CalcPointTsumoSum(han, fu, yakumanTimes, _hi.IsParent)
		} else {
			pt = CalcPointRon(han, fu, yakumanTimes, _hi.IsParent)
		}
		_result := &PointResult{
			han,
			fu,
			yakumanTimes,
			_hi.IsParent,
			pt,
		}
		// 高点法
		if pt > result.Point {
			result = _result
		} else if pt == result.Point {
			if han > result.han {
				result = _result
			}
		}
	}
	return
}

// 已听牌，根据 playerInfo 提供的信息计算加权和率后的平均点数
// 无役时返回 0
func CalcAvgPoint(playerInfo model.PlayerInfo, waits Waits) (avgPoint float64) {
	isFuriten := playerInfo.IsFuriten(waits)
	if isFuriten {
		// 振听只能自摸
		playerInfo.IsTsumo = true
	}

	tileAgariRate := CalculateAgariRateOfEachTile(waits, &playerInfo)
	sum := 0.0
	weight := 0.0
	for tile, left := range waits {
		if left == 0 {
			continue
		}
		playerInfo.HandTiles34[tile]++
		playerInfo.WinTile = tile
		result := CalcPoint(&playerInfo)
		playerInfo.HandTiles34[tile]--
		if result.Point == 0 {
			// 不考虑部分无役（如后附、片听）
			continue
		}
		var w float64
		if playerInfo.IsTsumo {
			w = float64(left) // 如果是自摸的话，只看枚数
		} else {
			w = tileAgariRate[tile] // 荣和考虑各个牌的和率
		}
		pt := float64(result.Point)
		if playerInfo.IsRiichi {
			// 如果立直了，需要考虑一发和里宝
			pt = result.fixedRiichiPoint(isFuriten)
		}
		sum += pt * w
		weight += w
	}
	if weight > 0 {
		avgPoint = sum / weight
	}
	return
}

// 计算立直时的平均点数（考虑一发和里宝）
// 已鸣牌时返回 0
// TODO: 剩余不到 4 张无法立直
// TODO: 分数不足 1000 无法立直
func CalcAvgRiichiPoint(playerInfo model.PlayerInfo, waits Waits) (avgRiichiPoint float64) {
	if playerInfo.IsNaki() {
		return 0
	}
	playerInfo.IsRiichi = true
	return CalcAvgPoint(playerInfo, waits)
}
