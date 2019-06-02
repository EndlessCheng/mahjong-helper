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
	Point      int
	FixedPoint float64 // 和牌时的期望点数

	han          int
	fu           int
	yakumanTimes int
	isParent     bool

	divideResult *DivideResult
	winTile      int
	yakuTypes    []int
	agariRate    float64 // 无役时的和率为 0
}

// 已和牌，计算自摸或荣和时的点数（不考虑里宝、一发等情况）
// 无役时返回的点数为 0（和率也为 0）
// 调用前请设置 IsTsumo WinTile
func CalcPoint(playerInfo *model.PlayerInfo) (result *PointResult) {
	result = &PointResult{}
	isNaki := playerInfo.IsNaki()
	var han, fu int
	numDora := playerInfo.CountDora()
	for _, divideResult := range DivideTiles34(playerInfo.HandTiles34) {
		_hi := &_handInfo{
			PlayerInfo:   playerInfo,
			divideResult: divideResult,
		}
		yakuTypes := findYakuTypes(_hi, isNaki)
		if len(yakuTypes) == 0 {
			// 此手牌拆解下无役
			continue
		}
		yakumanTimes := CalcYakumanTimes(yakuTypes, isNaki)
		if yakumanTimes == 0 {
			han = CalcYakuHan(yakuTypes, isNaki)
			han += numDora
			fu = _hi.calcFu(isNaki)
		}
		var pt int
		if _hi.IsTsumo {
			pt = CalcPointTsumoSum(han, fu, yakumanTimes, _hi.IsParent)
		} else {
			pt = CalcPointRon(han, fu, yakumanTimes, _hi.IsParent)
		}
		_result := &PointResult{
			pt,
			float64(pt),
			han,
			fu,
			yakumanTimes,
			_hi.IsParent,
			divideResult,
			_hi.WinTile,
			yakuTypes,
			0.0, // 后面会补上
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
// 有役时返回平均点数（立直时考虑自摸、一发和里宝）和各种侍牌下的对应点数
func CalcAvgPoint(playerInfo model.PlayerInfo, waits Waits) (avgPoint float64, pointResults []*PointResult) {
	isFuriten := playerInfo.IsFuriten(waits)
	if isFuriten {
		// 振听只能自摸，但是振听立直时考虑了这一点，所以只在默听或鸣牌时考虑
		if !playerInfo.IsRiichi {
			playerInfo.IsTsumo = true
		}
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
		result := CalcPoint(&playerInfo) // 非振听时，这里算出的是荣和的点数
		playerInfo.HandTiles34[tile]--
		if result.Point == 0 {
			// 不考虑部分无役（如后附、片听）
			continue
		}
		pt := float64(result.Point)
		if playerInfo.IsRiichi {
			// 如果立直了，需要考虑自摸、一发和里宝
			pt = result.fixedRiichiPoint(isFuriten)
			result.FixedPoint = pt
		}
		w := tileAgariRate[tile]
		sum += pt * w
		weight += w
		result.agariRate = w
		pointResults = append(pointResults, result)
	}
	if weight > 0 {
		avgPoint = sum / weight
	}
	return
}

// 计算立直时的平均点数（考虑自摸、一发和里宝）和各种侍牌下的对应点数
// 已鸣牌时返回 0
// TODO: 剩余不到 4 张无法立直
// TODO: 不足 1000 点无法立直
func CalcAvgRiichiPoint(playerInfo model.PlayerInfo, waits Waits) (avgRiichiPoint float64, pointResults []*PointResult) {
	if playerInfo.IsNaki() {
		return 0, nil
	}
	playerInfo.IsRiichi = true
	return CalcAvgPoint(playerInfo, waits)
}
