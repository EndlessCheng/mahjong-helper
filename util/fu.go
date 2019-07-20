package util

import (
	"github.com/EndlessCheng/mahjong-helper/util/model"
)

func roundUpFu(fu int) int {
	return ((fu-1)/10 + 1) * 10
}

// 根据手牌拆解结果，结合场况计算符数
func (hi *_handInfo) calcFu(isNaki bool) int {
	divideResult := hi.divideResult

	// 特殊：七对子计 25 符
	if divideResult.IsChiitoi {
		return 25
	}

	const baseFu = 20

	// 符底 20 符
	fu := baseFu

	// 暗刻加符
	_, ronKotsu := hi.numAnkou()
	for _, tile := range divideResult.KotsuTiles {
		var _fu int
		// 荣和刻子算明刻
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
		case model.MeldTypeMinkan, model.MeldTypeKakan:
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
		if hi.isDoubleWindTile(divideResult.PairTile) {
			fu += 2
		}
	}

	if fu == baseFu {
		// 手牌全是顺子，且雀头不是役牌
		if isNaki {
			// 无论怎样都不可能超过 30 符，直接返回
			return 30
		}
		// 门清状态下需要检测能否平和
		// 若没有平和则一定是坎张、边张、单骑和牌
		isPinfu := false
		for _, tile := range divideResult.ShuntsuFirstTiles {
			t9 := tile % 9
			if t9 < 6 && tile == hi.WinTile || t9 > 0 && tile+2 == hi.WinTile {
				isPinfu = true
				break
			}
		}
		if hi.IsTsumo {
			if isPinfu {
				// 门清自摸平和 20 符
				return 20
			}
			// 坎张、边张、单骑自摸，30 符
			return 30
		} else {
			// 荣和
			if isPinfu {
				// 门清平和荣和 30 符
				return 30
			}
			// 坎张、边张、单骑荣和，40 符
			return 40
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
				fu += 2 // 边张和牌加符
				break
			}
		}
	}

	// 进位
	return roundUpFu(fu)
}
