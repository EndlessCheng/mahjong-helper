package util

// 门清限定
func (hi *_handInfo) suuAnkou() bool {
	if hi.WinTile == hi.divideResult.PairTile {
		return false
	}
	// 非单骑和牌
	n, _ := hi.numAnkou()
	return n == 4
}

// 门清限定
func (hi *_handInfo) suuAnkouTanki() bool {
	if hi.WinTile != hi.divideResult.PairTile {
		return false
	}
	// 单骑和牌
	n, _ := hi.numAnkou()
	return n == 4
}

func (hi *_handInfo) daisangen() bool {
	return hi._countSpecialKotsu(31, 33) == 3
}

func (hi *_handInfo) shousuushii() bool {
	return hi.divideResult.PairTile >= 27 && hi.divideResult.PairTile <= 30 && hi._countSpecialKotsu(27, 30) == 3
}

func (hi *_handInfo) daisuushii() bool {
	return hi._countSpecialKotsu(27, 30) == 4
}

func (hi *_handInfo) tsuuiisou() bool {
	if hi.divideResult.IsChiitoi {
		// 大七星
		for _, c := range hi.HandTiles34[27:] {
			if c == 0 {
				return false
			}
		}
		return true
	}

	if hi.divideResult.PairTile < 27 {
		return false
	}
	if len(hi.allShuntsuFirstTiles) > 0 {
		return false
	}
	for _, tile := range hi.allKotsuTiles {
		if tile < 27 {
			return false
		}
	}
	return true
}

func (hi *_handInfo) chinroutou() bool {
	if hi.divideResult.IsChiitoi {
		return false
	}

	isValid := func(tile int) bool {
		if tile >= 27 {
			return false
		}
		t9 := tile % 9
		return t9 == 0 || t9 == 8
	}

	if !isValid(hi.divideResult.PairTile) {
		return false
	}
	if len(hi.allShuntsuFirstTiles) > 0 {
		return false
	}
	for _, tile := range hi.allKotsuTiles {
		if !isValid(tile) {
			return false
		}
	}
	return true
}

var _ryuuTiles = []int{19, 20, 21, 23, 25, 32}

func (hi *_handInfo) ryuuiisou() bool {
	if hi.divideResult.IsChiitoi {
		return false
	}

	for _, tile := range hi.allShuntsuFirstTiles {
		if tile != 19 { // 只能是 234s
			return false
		}
	}
	if !InInts(hi.divideResult.PairTile, _ryuuTiles) {
		return false
	}
	for _, tile := range hi.allKotsuTiles {
		if !InInts(tile, _ryuuTiles) {
			return false
		}
	}
	return true
}

// 调用前已经不是七对了
func (hi *_handInfo) _isChuuren9() bool {
	// 去掉 WinTile 后，剩余的牌必须是 1112345678999
	// 也就是说，hi.HandTiles34[hi.WinTile] 多出的那一枚必须正好是 WinTile
	tileType := hi.WinTile / 9
	tiles34 := hi.HandTiles34
	idx := 9 * tileType
	if tiles34[idx] == 4 {
		return hi.WinTile == idx
	}
	end := 9*tileType + 8
	for ; idx < end; idx++ { // 2~8
		if tiles34[idx] == 2 {
			return hi.WinTile == idx
		}
	}
	if tiles34[idx] == 4 {
		return hi.WinTile == idx
	}
	return false
}

// 门清限定
func (hi *_handInfo) chuuren() bool {
	return hi.divideResult.IsChuurenPoutou && !hi._isChuuren9()
}

// 门清限定
func (hi *_handInfo) chuuren9() bool {
	return hi.divideResult.IsChuurenPoutou && hi._isChuuren9()
}

func (hi *_handInfo) suuKantsu() bool {
	return hi.numKantsu() == 4
}

var yakumanCheckerMap = map[int]yakuChecker{
	YakuSuuAnkou:      (*_handInfo).suuAnkou,
	YakuSuuAnkouTanki: (*_handInfo).suuAnkouTanki,
	YakuDaisangen:     (*_handInfo).daisangen,
	YakuShousuushii:   (*_handInfo).shousuushii,
	YakuDaisuushii:    (*_handInfo).daisuushii,
	YakuTsuuiisou:     (*_handInfo).tsuuiisou,
	YakuChinroutou:    (*_handInfo).chinroutou,
	YakuRyuuiisou:     (*_handInfo).ryuuiisou,
	YakuChuuren:       (*_handInfo).chuuren,
	YakuChuuren9:      (*_handInfo).chuuren9,
	YakuSuuKantsu:     (*_handInfo).suuKantsu,
}

//

func (*_handInfo) checkAllPairs(tiles []int) bool {
	for _, c := range tiles {
		if c != 2 {
			return false
		}
	}
	return true
}

// 门清限定
func (hi *_handInfo) daisuurin() bool {
	return hi.checkAllPairs(hi.HandTiles34[1:8])
}

// 门清限定
func (hi *_handInfo) daisharin() bool {
	return hi.checkAllPairs(hi.HandTiles34[9+1 : 9+8])
}

// 门清限定
func (hi *_handInfo) daichikurin() bool {
	return hi.checkAllPairs(hi.HandTiles34[18+1 : 18+8])
}

// 门清限定
func (hi *_handInfo) daichisei() bool {
	return hi.checkAllPairs(hi.HandTiles34[27:])
}

var oldYakumanCheckerMap = map[int]yakuChecker{
	YakuDaisuurin:   (*_handInfo).daisuurin,
	YakuDaisharin:   (*_handInfo).daisharin,
	YakuDaichikurin: (*_handInfo).daichikurin,
	YakuDaichisei:   (*_handInfo).daichisei,
}

// 检测役满
// 结果未排序
// *计算前必须设置顺子牌和刻子牌
func findYakumanTypes(hi *_handInfo, isNaki bool) (yakumanTypes []int) {
	var yakumanTimesMap _yakumanTimesMap
	if !isNaki {
		yakumanTimesMap = YakumanTimesMap
	} else {
		yakumanTimesMap = NakiYakumanTimesMap
	}

	for yakuman := range yakumanTimesMap {
		if checker, ok := yakumanCheckerMap[yakuman]; ok {
			if checker(hi) {
				yakumanTypes = append(yakumanTypes, yakuman)
			}
		}
	}

	if considerOldYaku && !isNaki {
		for yakuman := range OldYakumanTimesMap {
			if checker, ok := oldYakumanCheckerMap[yakuman]; ok {
				if checker(hi) {
					yakumanTypes = append(yakumanTypes, yakuman)
				}
			}
		}
	}

	return
}
