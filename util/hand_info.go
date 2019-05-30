package util

import "github.com/EndlessCheng/mahjong-helper/util/model"

type _handInfo struct {
	*model.PlayerInfo
	divideResult *DivideResult // 手牌解析结果

	// *在计算役种前，缓存自己的顺子牌和刻子牌，这样能减少大量重复计算
	allShuntsuFirstTiles []int
	allKotsuTiles        []int
}

// 未排序。用于算一通、三色
func (hi *_handInfo) getAllShuntsuFirstTiles() []int {
	shuntsuFirstTiles := append([]int{}, hi.divideResult.ShuntsuFirstTiles...)
	for _, meld := range hi.Melds {
		if meld.MeldType == model.MeldTypeChi {
			shuntsuFirstTiles = append(shuntsuFirstTiles, meld.Tiles[0])
		}
	}
	return shuntsuFirstTiles
}

// 未排序。用于算对对、三色同刻
func (hi *_handInfo) getAllKotsuTiles() []int {
	kotsuTiles := append([]int{}, hi.divideResult.KotsuTiles...)
	for _, meld := range hi.Melds {
		if meld.MeldType != model.MeldTypeChi {
			kotsuTiles = append(kotsuTiles, meld.Tiles[0])
		}
	}
	return kotsuTiles
}

// 是否包含字牌（调用前需要设置刻子牌）
func (hi *_handInfo) containHonor() bool {
	// 七对子特殊处理
	if hi.divideResult.IsChiitoi {
		for _, c := range hi.HandTiles34[27:] {
			if c > 0 {
				return true
			}
		}
		return false
	}

	if hi.divideResult.PairTile >= 27 {
		return true
	}
	for _, tile := range hi.allKotsuTiles {
		if tile >= 27 {
			return true
		}
	}
	return false
}

// 是否为役牌，用于算役种（役牌、平和）、雀头加符
func (hi *_handInfo) isYakuTile(tile int) bool {
	return tile >= 31 || tile == hi.RoundWindTile || tile == hi.SelfWindTile
}

// 是否为连风牌
func (hi *_handInfo) isDoubleWindTile(tile int) bool {
	return hi.RoundWindTile == hi.SelfWindTile && tile == hi.RoundWindTile
}

// 暗杠个数，用于算三暗刻、四暗刻
func (hi *_handInfo) numAnkan() (cnt int) {
	for _, meld := range hi.Melds {
		if meld.MeldType == model.MeldTypeAnkan {
			cnt++
		}
	}
	return
}

// 杠子个数，用于算三杠子、四杠子
func (hi *_handInfo) numKantsu() (cnt int) {
	for _, meld := range hi.Melds {
		if meld.IsKan() {
			cnt++
		}
	}
	return
}

// 暗刻个数，用于算三暗刻、四暗刻、符数（如 456666 荣和 6，这里算一个暗刻）
// 即手牌暗刻和暗杠
func (hi *_handInfo) numAnkou() (cnt int, isMinkou bool) {
	num := len(hi.divideResult.KotsuTiles) + hi.numAnkan()
	// 自摸直接返回，无需讨论是否荣和了刻子
	if hi.IsTsumo {
		return num, false
	}
	// 荣和的牌在雀头里
	if hi.WinTile == hi.divideResult.PairTile {
		return num, false
	}
	// 荣和的牌在顺子里
	for _, tile := range hi.divideResult.ShuntsuFirstTiles {
		if hi.WinTile >= tile && hi.WinTile <= tile+2 {
			return num, false
		}
	}
	// 荣和的牌只在刻子里，该刻子算明刻
	return num - 1, true
}

// 计算在指定牌中的刻子个数
func (hi *_handInfo) _countSpecialKotsu(specialTilesL, specialTilesLR int) (cnt int) {
	for _, tile := range hi.allKotsuTiles {
		if tile >= specialTilesL && tile <= specialTilesLR {
			cnt++
		}
	}
	return
}
