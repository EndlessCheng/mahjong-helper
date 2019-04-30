package util

// TODO: 考虑大三元和大四喜的包牌？

func roundUpPoint(point int) int {
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
