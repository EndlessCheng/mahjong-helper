package util

import "github.com/EndlessCheng/mahjong-helper/util/model"

// 门清限定
func (hi *_handInfo) suuAnkou() bool {
	// 非单骑和牌
	return hi.WinTile != hi.divideResult.PairTile && hi.numAnkou() == 4
}

// 门清限定
func (hi *_handInfo) suuAnkouTanki() bool {
	// 单骑和牌
	return hi.WinTile == hi.divideResult.PairTile && hi.numAnkou() == 4
}

// 计算在指定牌中的刻子个数
func (hi *_handInfo) _countSpecialKotsu(specialTilesL, specialTilesLR int) (cnt int) {
	for _, tile := range hi.divideResult.KotsuTiles {
		if tile >= specialTilesL && tile <= specialTilesLR {
			cnt++
		}
	}
	for _, meld := range hi.Melds {
		if meld.MeldType != model.MeldTypeChi {
			if tile := meld.Tiles[0]; tile >= specialTilesL && tile <= specialTilesLR {
				cnt++
			}
		}
	}
	return
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
	if hi.divideResult.PairTile < 27 {
		return false
	}
	if len(hi.divideResult.ShuntsuFirstTiles) > 0 {
		return false
	}
	for _, tile := range hi.divideResult.KotsuTiles {
		if tile < 27 {
			return false
		}
	}
	for _, meld := range hi.Melds {
		if meld.MeldType == model.MeldTypeChi {
			return false
		}
		if meld.Tiles[0] < 27 {
			return false
		}
	}
	return true
}

func (hi *_handInfo) chinroutou() bool {
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
	if len(hi.divideResult.ShuntsuFirstTiles) > 0 {
		return false
	}
	for _, tile := range hi.divideResult.KotsuTiles {
		if !isValid(tile) {
			return false
		}
	}
	for _, meld := range hi.Melds {
		if meld.MeldType == model.MeldTypeChi {
			return false
		}
		if !isValid(meld.Tiles[0]) {
			return false
		}
	}
	return true
}

func (hi *_handInfo) ryuuiisou() bool {
	isValid := func(tile int) bool {
		return InInts(tile, []int{19, 20, 21, 23, 25, 32})
	}

	if !isValid(hi.divideResult.PairTile) {
		return false
	}
	for _, tile := range hi.divideResult.ShuntsuFirstTiles {
		if tile != 19 { // 只能是 234s
			return false
		}
	}
	for _, tile := range hi.divideResult.KotsuTiles {
		if !isValid(tile) {
			return false
		}
	}
	for _, meld := range hi.Melds {
		tile := meld.Tiles[0]
		if meld.MeldType == model.MeldTypeChi {
			if tile != 19 { // 只能是 234s
				return false
			}
		} else {
			if !isValid(tile) {
				return false
			}
		}
	}
	return true
}

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
	for ; idx < end; idx++ {
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

func findYakumanTypes(hi *_handInfo) (yakumanTypes []int) {
	for yakuman, checker := range yakumanCheckerMap {
		if checker(hi) {
			yakumanTypes = append(yakumanTypes, yakuman)
		}
	}
	return
}
