package util

import "github.com/EndlessCheng/mahjong-helper/util/model"

// 是否鸣牌
func (hi *HandInfo) isNaki() bool {
	if hi._isNaki != nil {
		return *hi._isNaki
	}
	naki := false
	for _, meld := range hi.Melds {
		if meld.MeldType != model.MeldTypeAnkan {
			naki = true
			break
		}
	}
	hi._isNaki = &naki
	return naki
}

// 是否包含字牌
func (hi *HandInfo) containHonor() bool {
	if hi._containHonor != nil {
		return *hi._containHonor
	}
	ch := func() bool {
		if hi.Divide.PairTile >= 27 {
			return true
		}
		for _, tile := range hi.Divide.KotsuTiles {
			if tile >= 27 {
				return true
			}
		}
		for _, meld := range hi.Melds {
			if meld.MeldType != model.MeldTypeChi && meld.Tiles[0] >= 27 {
				return true
			}
		}
		return false
	}
	cont := ch()
	hi._containHonor = &cont
	return cont
}

// 是否为役牌，用于算役种（役牌、平和）、雀头加符
func (hi *HandInfo) isYakuTile(tile int) bool {
	return tile >= 31 || tile == hi.RoundWindTile || tile == hi.SelfWindTile
}

// 暗刻个数，用于算三暗刻、四暗刻、符数（如 456666 荣和 6，这里算一个暗刻）
func (hi *HandInfo) numAnkou() int {
	num := len(hi.Divide.KotsuTiles)
	if hi.IsTsumo {
		return num
	}
	// 荣和的牌在雀头里
	if hi.WinTile == hi.Divide.PairTile {
		return num
	}
	// 荣和的牌在顺子里
	for _, tile := range hi.Divide.ShuntsuFirstTiles {
		if hi.WinTile >= tile && hi.WinTile <= tile+2 {
			return num
		}
	}
	// 荣和的牌只在刻子里，该刻子算明刻
	return num - 1
}

// 杠子个数，用于算三杠子、四杠子
func (hi *HandInfo) numKantsu() int {
	cnt := 0
	for _, meld := range hi.Melds {
		if meld.MeldType == model.MeldTypeAnkan || meld.MeldType == model.MeldTypeMinkan || meld.MeldType == model.MeldTypeKakan {
			cnt++
		}
	}
	return cnt
}

//

// 门清限定
func (hi *HandInfo) chiitoi() bool {
	return hi.Divide.IsChiitoi
}

// 门清限定
func (hi *HandInfo) pinfu() bool {
	// 雀头不能是役牌，且不能是单骑和牌
	if hi.isYakuTile(hi.WinTile) || hi.WinTile == hi.Divide.PairTile {
		return false
	}
	drs := hi.Divide.ShuntsuFirstTiles
	// 不能有刻子
	if len(drs) < 4 {
		return false
	}
	for _, tile := range drs {
		// 可以两面和牌
		if tile%9 < 6 && tile == hi.WinTile || tile%9 > 0 && tile+2 == hi.WinTile {
			return true
		}
	}
	// 没有两面和牌
	return false
}

// 门清限定
func (hi *HandInfo) ryanpeikou() bool {
	return hi.Divide.IsRyanpeikou
}

// 门清限定
// 两杯口时不算一杯口
func (hi *HandInfo) iipeikou() bool {
	return hi.Divide.IsIipeikou
}

func (hi *HandInfo) sanshokuDoujun() bool {
	shuntsuFirstTiles := append([]int{}, hi.Divide.ShuntsuFirstTiles...)
	for _, meld := range hi.Melds {
		if meld.MeldType == model.MeldTypeChi {
			shuntsuFirstTiles = append(shuntsuFirstTiles, meld.Tiles[0])
		}
	}
	if len(shuntsuFirstTiles) < 3 {
		return false
	}
	var sMan, sPin, sSou []int
	for _, s := range shuntsuFirstTiles {
		if isMan(s) {
			sMan = append(sMan, s)
		} else if isPin(s) {
			sPin = append(sPin, s)
		} else { // isSou
			sSou = append(sSou, s)
		}
	}
	for _, man := range sMan {
		for _, pin := range sPin {
			for _, sou := range sSou {
				if man == pin-9 && man == sou-18 {
					return true
				}
			}
		}
	}
	return false
}

func (hi *HandInfo) ittsuu() bool {
	if !hi.isNaki() {
		return hi.Divide.IsIttsuu
	}
	shuntsuFirstTiles := append([]int{}, hi.Divide.ShuntsuFirstTiles...)
	for _, meld := range hi.Melds {
		if meld.MeldType == model.MeldTypeChi {
			shuntsuFirstTiles = append(shuntsuFirstTiles, meld.Tiles[0])
		}
	}
	// 若有 123，找是否有同色的 456 和 789
	for _, tile := range shuntsuFirstTiles {
		if tile%9 == 0 {
			has456 := false
			has789 := false
			for _, otherTile := range shuntsuFirstTiles {
				if otherTile == tile+3 {
					has456 = true
				} else if otherTile == tile+6 {
					has789 = true
				}
			}
			if has456 && has789 {
				return true
			}
		}
	}
	return false
}

func (hi *HandInfo) toitoi() bool {
	numKotsu := len(hi.Divide.KotsuTiles)
	for _, meld := range hi.Melds {
		if meld.MeldType != model.MeldTypeChi {
			numKotsu++
		}
	}
	return numKotsu == 4
}

// 荣和的刻子是明刻
// 注意 456666 这样的荣和 6，算暗刻
func (hi *HandInfo) sanAnkou() bool {
	if len(hi.Divide.KotsuTiles) < 3 {
		return false
	}
	return hi.numAnkou() == 3
}

func (hi *HandInfo) sanshokuDoukou() bool {
	kotsuTiles := append([]int{}, hi.Divide.KotsuTiles...)
	for _, meld := range hi.Melds {
		if meld.MeldType != model.MeldTypeChi {
			kotsuTiles = append(kotsuTiles, meld.Tiles[0])
		}
	}
	if len(kotsuTiles) < 3 {
		return false
	}
	var kMan, kPin, kSou []int
	for _, tile := range kotsuTiles {
		if isMan(tile) {
			kMan = append(kMan, tile)
		} else if isPin(tile) {
			kPin = append(kPin, tile)
		} else if isSou(tile) {
			kSou = append(kSou, tile)
		}
	}
	for _, man := range kMan {
		for _, pin := range kPin {
			for _, sou := range kSou {
				if man == pin-9 && man == sou-18 {
					return true
				}
			}
		}
	}
	return false
}

func (hi *HandInfo) sanKantsu() bool {
	if len(hi.Melds) < 3 {
		return false
	}
	return hi.numKantsu() == 3
}

func (hi *HandInfo) tanyao() bool {
	dr := hi.Divide
	if isYaochupai(dr.PairTile) {
		return false
	}
	for _, tile := range dr.KotsuTiles {
		if isYaochupai(tile) {
			return false
		}
	}
	for _, tile := range dr.ShuntsuFirstTiles {
		if isYaochupai(tile) || isYaochupai(tile+2) {
			return false
		}
	}
	for _, meld := range hi.Melds {
		tiles := meld.Tiles
		if meld.MeldType == model.MeldTypeChi {
			if isYaochupai(tiles[0]) || isYaochupai(tiles[2]) {
				return false
			}
		} else {
			if isYaochupai(tiles[0]) {
				return false
			}
		}
	}
	return true
}

// 返回役牌个数
func (hi *HandInfo) yakuhai() int {
	cnt := 0
	for _, tile := range hi.Divide.KotsuTiles {
		if hi.isYakuTile(tile) {
			cnt++
		}
	}
	for _, meld := range hi.Melds {
		if meld.MeldType != model.MeldTypeChi && hi.isYakuTile(meld.Tiles[0]) {
			cnt++
		}
	}
	return cnt
}

func (hi *HandInfo) chanta() bool {
	if !hi.containHonor() {
		return false
	}
	// TODO
	return false
}

func (hi *HandInfo) junchan() bool {
	if hi.containHonor() {
		return false
	}
	// TODO
	return false
}

func (hi *HandInfo) honroutou() bool {
	if !hi.containHonor() {
		return false
	}
	// TODO
	return false
}

func (hi *HandInfo) shousangen() bool {
	// TODO
	return false
}

func (hi *HandInfo) honitsu() bool {
	if !hi.containHonor() {
		return false
	}
	// TODO
	return false
}

func (hi *HandInfo) chinitsu() bool {
	if hi.containHonor() {
		return false
	}
	// TODO
	return false
}

func FindYakuList(hi *HandInfo) (yakuList []YakuType) {
	// TODO
	return
}
