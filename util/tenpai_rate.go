package util

import "github.com/EndlessCheng/mahjong-helper/util/model"

// 没有立直时，根据玩家的副露、手切来判断其听牌率 (0-100)
// TODO: 传入 *model.PlayerInfo
func CalcTenpaiRate(melds []*model.Meld, discardTiles []int, meldDiscardsAt []int) float64 {
	isNaki := false
	for _, meld := range melds {
		if meld.MeldType != model.MeldTypeAnkan {
			isNaki = true
		}
	}

	if !isNaki {
		// 默听听牌率近似为巡目数
		turn := len(discardTiles)
		return float64(turn)
	}

	if len(melds) == 4 {
		return 100
	}

	_tenpaiRate := tenpaiRate[len(melds)]

	turn := MinInt(len(discardTiles), len(_tenpaiRate)-1)
	_tenpaiRateWithTurn := _tenpaiRate[turn]

	// 计算上一次副露后的手切数
	// 注意连续开杠时，副露数 len(melds) 是不等于副露时的切牌数 len(meldDiscardsAt) 的
	countTedashi := 0
	if len(meldDiscardsAt) > 0 {
		latestDiscardAt := meldDiscardsAt[len(meldDiscardsAt)-1]
		if len(discardTiles) > latestDiscardAt {
			for _, disTile := range discardTiles[latestDiscardAt+1:] {
				if disTile >= 0 {
					countTedashi++
				}
			}
		}
	}
	countTedashi = MinInt(countTedashi, len(_tenpaiRateWithTurn)-1)

	return _tenpaiRateWithTurn[countTedashi]
}
