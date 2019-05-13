package util

import "github.com/EndlessCheng/mahjong-helper/util/model"

func roundUpFu(fu int) int {
	return ((fu-1)/10 + 1) * 10
}

// 根据手牌拆解结果，结合场况计算符数
func (hi *_handInfo) calcFu() int {
	divideResult := hi.divideResult

	// 特殊：七对子计 25 符
	if divideResult.IsChiitoi {
		return 25
	}

	const baseFu = 20

	// 符底 20 符
	fu := baseFu

	// 暗刻加符
	// 若刻子数不等于暗刻数，则荣和的牌被算到了刻子中
	ronKotsu := len(divideResult.KotsuTiles) != hi.numAnkou()
	for _, tile := range divideResult.KotsuTiles {
		var _fu int
		// 荣和算明刻
		if ronKotsu && tile == hi.WinTile {
			_fu = 2
		} else {
			_fu = 4
		}
		if isYaochupai(tile) {
			_fu *= 2
		}
		fu += _fu
	}

	// 明刻、明杠、暗杠加符
	for _, meld := range hi.Melds {
		_fu := 0
		switch meld.MeldType {
		case model.MeldTypePon:
			_fu = 2
		case model.MeldTypeMinkan:
		case model.MeldTypeKakan:
			_fu = 8
		case model.MeldTypeAnkan:
			_fu = 16
		}
		if _fu > 0 {
			if isYaochupai(meld.Tiles[0]) {
				_fu *= 2
			}
			fu += _fu
		}
	}

	// 雀头加符（连风雀头计 4 符）
	if hi.isYakuTile(divideResult.PairTile) {
		fu += 2
		if hi.isDoubleYakuTile(divideResult.PairTile) {
			fu += 2
		}
	}

	// 是否鸣牌
	isNaki := hi.IsNaki()

	// 特殊：门清 + 自摸 + 平和型，计 20 符
	if !isNaki && hi.IsTsumo && fu == baseFu {
		// 考虑能否两面和牌
		for _, tile := range divideResult.ShuntsuFirstTiles {
			if tile%9 < 6 && tile == hi.WinTile || tile%9 > 0 && tile+2 == hi.WinTile {
				return 20
			}
		}
	}

	// 门清荣和加符
	if !isNaki && !hi.IsTsumo {
		fu += 10
	}

	// 自摸加符
	if hi.IsTsumo {
		fu += 2
	}

	// 边张、坎张、单骑和牌加符
	// 考虑能否不为两面和牌
	if divideResult.PairTile == hi.WinTile {
		fu += 2 // 单骑和牌加符
	} else {
		for _, tile := range divideResult.ShuntsuFirstTiles {
			if tile+1 == hi.WinTile {
				fu += 2 // 坎张和牌加符
				break
			}
			if tile%9 == 0 && tile+2 == hi.WinTile || tile%9 == 6 && tile == hi.WinTile {
				fu += 2 // 边张和牌加符（123 和 3，789 和 7）
				break
			}
		}
	}

	// 特殊：若仍然为 20 符（副露荣和平和型）视作 30 符
	if fu == baseFu {
		return 30
	}

	// 进位
	return roundUpFu(fu)
}
