package util

import "github.com/EndlessCheng/mahjong-helper/util/model"

// 四副露大吊车，不能有暗杠
func (hi *_handInfo) shiiaruraotai() bool {
	if len(hi.Melds) < 4 {
		return false
	}
	for _, meld := range hi.Melds {
		if meld.MeldType == model.MeldTypeAnkan {
			return false
		}
	}
	return true
}

func (hi *_handInfo) uumensai() bool {
	suits := map[int]bool{}
	addSuit := func(tile int) {
		if tile < 27 {
			suits[tile/9] = true
		} else if tile < 31 {
			suits[3] = true
		} else {
			suits[4] = true
		}
	}
	addSuit(hi.divideResult.PairTile)
	for _, tile := range hi.allShuntsuFirstTiles {
		addSuit(tile)
	}
	for _, tile := range hi.allKotsuTiles {
		addSuit(tile)
	}
	return len(suits) == 5
}

func (hi *_handInfo) sanrenkou() bool {
	if len(hi.allKotsuTiles) < 3 {
		return false
	}
	if hi.allKotsuTiles[0] < 27 && hi.allKotsuTiles[0]%9 < 7 && hi.allKotsuTiles[0]+1 == hi.allKotsuTiles[1] && hi.allKotsuTiles[1]+1 == hi.allKotsuTiles[2] {
		return true
	}
	if len(hi.allKotsuTiles) == 4 {
		if hi.allKotsuTiles[1] < 27 && hi.allKotsuTiles[1]%9 < 7 && hi.allKotsuTiles[1]+1 == hi.allKotsuTiles[2] && hi.allKotsuTiles[2]+1 == hi.allKotsuTiles[3] {
			return true
		}
	}
	return false
}

func (hi *_handInfo) isshokusanjun() bool {
	if len(hi.allShuntsuFirstTiles) < 3 {
		return false
	}
	if hi.allShuntsuFirstTiles[0] == hi.allShuntsuFirstTiles[1] && hi.allShuntsuFirstTiles[1] == hi.allShuntsuFirstTiles[2] {
		return true
	}
	if len(hi.allShuntsuFirstTiles) == 4 {
		if hi.allShuntsuFirstTiles[1] == hi.allShuntsuFirstTiles[2] && hi.allShuntsuFirstTiles[2] == hi.allShuntsuFirstTiles[3] {
			return true
		}
	}
	return false
}

var oldYakuCheckerMap = map[int]yakuChecker{
	YakuShiiaruraotai: (*_handInfo).shiiaruraotai,
	YakuUumensai:      (*_handInfo).uumensai,
	YakuSanrenkou:     (*_handInfo).sanrenkou,
	YakuIsshokusanjun: (*_handInfo).isshokusanjun,
}
