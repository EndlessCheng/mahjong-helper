package util

import "github.com/EndlessCheng/mahjong-helper/util/model"

func roundUpFu(fu int) int {
	return ((fu-1)/10 + 1) * 10
}

// 根据手牌拆解结果，结合场况计算符数
func (d *DivideResult) Fu(winTile int, isTsumo bool, melds []model.Meld, roundWindTile int, selfWindTile int) int {
	// 特殊：七对子计 25 符
	if d.IsChiitoi {
		return 25
	}

	const baseFu = 20

	// 符底
	fu := baseFu

	// 暗刻加符
	for _, tile := range d.KotsuTiles {
		var _fu int
		// 荣和算明刻
		if !isTsumo && tile == winTile {
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
	for _, meld := range melds {
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
		if isYaochupai(meld.Tiles[0]) {
			_fu *= 2
		}
		fu += _fu
	}

	// 雀头加符（场风与自风重合时计 4 符）
	if d.PairTile == roundWindTile || d.PairTile == selfWindTile || d.PairTile >= 31 {
		fu += 2
		if roundWindTile == selfWindTile {
			fu += 2
		}
	}

	// 是否鸣牌
	isNaki := false
	for _, meld := range melds {
		if meld.MeldType != model.MeldTypeAnkan {
			isNaki = true
			break
		}
	}

	// 特殊：门清 + 两面自摸 + 平和型，计 20 符
	if !isNaki && isTsumo && fu == baseFu {
		// 考虑能否两面和牌
		for _, tile := range d.ShuntsuFirstTiles {
			if tile%9 < 6 && tile == winTile || tile%9 > 0 && tile+2 == winTile {
				return 20
			}
		}
	}

	// 门清加符
	if !isNaki {
		fu += 10
	}

	// 自摸加符
	if isTsumo {
		fu += 2
	}

	// 边张、坎张、单骑和牌加符
	// 考虑能否不为两面和牌
	if d.PairTile == winTile {
		fu += 2 // 单骑和牌加符
	} else {
		for _, tile := range d.ShuntsuFirstTiles {
			if tile+1 == winTile {
				fu += 2 // 坎张和牌加符
				break
			}
			if tile%9 == 0 && tile+2 == winTile || tile%9 == 6 && tile == winTile {
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
