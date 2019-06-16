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
	// 暂时不防默听玩家
	if !isNaki {
		return 0
	}

	if len(melds) == 4 {
		return 100
	}

	_tenpaiRate := tenpaiRate[len(melds)]

	turn := MinInt(len(discardTiles), len(_tenpaiRate)-1)
	_tenpaiRateWithTurn := _tenpaiRate[turn]

	// 计算上一次副露后的手切数
	countTedashi := 0
	if len(meldDiscardsAt) > 0 { // FIXME 实际上这恒为 true，只不过天凤偶尔会有先收到自家摸牌再收到上家摸牌的问题
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
