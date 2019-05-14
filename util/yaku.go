package util

import (
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"sort"
)

// 是否包含字牌
// cache 这货的话，其他地方都要 copy 了，目前项目采用引用的方式，不适合 cache
func (hi *_handInfo) containHonor() bool {
	// 门清时简化
	if len(hi.Melds) == 0 {
		for i := 27; i < 34; i++ {
			if hi.HandTiles34[i] > 0 {
				return true
			}
		}
		return false
	}
	if hi.divideResult.PairTile >= 27 {
		return true
	}
	for _, tile := range hi.divideResult.KotsuTiles {
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

// 是否为役牌，用于算役种（役牌、平和）、雀头加符
func (hi *_handInfo) isYakuTile(tile int) bool {
	return tile >= 31 || tile == hi.RoundWindTile || tile == hi.SelfWindTile
}

// 是否为连风牌
func (hi *_handInfo) isDoubleWindTile(tile int) bool {
	return hi.RoundWindTile == hi.SelfWindTile && tile == hi.RoundWindTile
}

// 暗刻个数，用于算三暗刻、四暗刻、符数（如 456666 荣和 6，这里算一个暗刻）
func (hi *_handInfo) numAnkou() int {
	num := len(hi.divideResult.KotsuTiles)
	if hi.IsTsumo {
		return num
	}
	// 荣和的牌在雀头里
	if hi.WinTile == hi.divideResult.PairTile {
		return num
	}
	// 荣和的牌在顺子里
	for _, tile := range hi.divideResult.ShuntsuFirstTiles {
		if hi.WinTile >= tile && hi.WinTile <= tile+2 {
			return num
		}
	}
	// 荣和的牌只在刻子里，该刻子算明刻
	return num - 1
}

// 杠子个数，用于算三杠子、四杠子
func (hi *_handInfo) numKantsu() int {
	cnt := 0
	for _, meld := range hi.Melds {
		if meld.IsKan() {
			cnt++
		}
	}
	return cnt
}

//

// 门清限定
func (hi *_handInfo) daburii() bool {
	return hi.IsDaburii
}

// 门清限定
func (hi *_handInfo) riichi() bool {
	return !hi.IsDaburii && hi.IsRiichi
}

// 门清限定
func (hi *_handInfo) tsumo() bool {
	return !hi.IsNaki() && hi.IsTsumo
}

// 门清限定
func (hi *_handInfo) chiitoi() bool {
	return hi.divideResult.IsChiitoi
}

// 门清限定
func (hi *_handInfo) pinfu() bool {
	// 不能是单骑和牌，雀头不能是役牌
	if hi.WinTile == hi.divideResult.PairTile || hi.isYakuTile(hi.divideResult.PairTile) {
		return false
	}
	drs := hi.divideResult.ShuntsuFirstTiles
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
func (hi *_handInfo) ryanpeikou() bool {
	return hi.divideResult.IsRyanpeikou
}

// 门清限定
// 两杯口时不算一杯口
func (hi *_handInfo) iipeikou() bool {
	return hi.divideResult.IsIipeikou
}

func (hi *_handInfo) sanshokuDoujun() bool {
	shuntsuFirstTiles := append([]int{}, hi.divideResult.ShuntsuFirstTiles...)
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

func (hi *_handInfo) ittsuu() bool {
	if !hi.IsNaki() {
		return hi.divideResult.IsIttsuu
	}
	shuntsuFirstTiles := append([]int{}, hi.divideResult.ShuntsuFirstTiles...)
	for _, meld := range hi.Melds {
		if meld.MeldType == model.MeldTypeChi {
			shuntsuFirstTiles = append(shuntsuFirstTiles, meld.Tiles[0])
		}
	}
	if len(shuntsuFirstTiles) < 3 {
		return false
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

func (hi *_handInfo) toitoi() bool {
	numKotsu := len(hi.divideResult.KotsuTiles)
	for _, meld := range hi.Melds {
		if meld.MeldType != model.MeldTypeChi {
			numKotsu++
		}
	}
	return numKotsu == 4
}

// 荣和的刻子是明刻
// 注意 456666 这样的荣和 6，算暗刻
func (hi *_handInfo) sanAnkou() bool {
	if len(hi.divideResult.KotsuTiles) < 3 {
		return false
	}
	return hi.numAnkou() == 3
}

func (hi *_handInfo) sanshokuDoukou() bool {
	kotsuTiles := append([]int{}, hi.divideResult.KotsuTiles...)
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

func (hi *_handInfo) sanKantsu() bool {
	if len(hi.Melds) < 3 {
		return false
	}
	return hi.numKantsu() == 3
}

func (hi *_handInfo) tanyao() bool {
	if len(hi.Melds) == 0 {
		// 门清时简单判断
		for _, tile := range YaochuTiles {
			if hi.HandTiles34[tile] > 0 {
				return false
			}
		}
		return true
	}

	dr := hi.divideResult
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

// 返回役牌个数，连风算两个
func (hi *_handInfo) numYakuhai() int {
	cnt := 0
	for _, tile := range hi.divideResult.KotsuTiles {
		if hi.isYakuTile(tile) {
			cnt++
			if hi.isDoubleWindTile(tile) {
				cnt++
			}
		}
	}
	for _, meld := range hi.Melds {
		tile := meld.Tiles[0]
		if meld.MeldType != model.MeldTypeChi && hi.isYakuTile(tile) {
			cnt++
			if hi.isDoubleWindTile(tile) {
				cnt++
			}
		}
	}
	return cnt
}

func (hi *_handInfo) _chantai() bool {
	dr := hi.divideResult
	// 必须有顺子
	if len(dr.ShuntsuFirstTiles) == 0 {
		hasShuntsu := false
		for _, meld := range hi.Melds {
			if meld.MeldType == model.MeldTypeChi {
				hasShuntsu = true
				break
			}
		}
		if !hasShuntsu {
			return false
		}
	}
	// 所有雀头和面子都要包含幺九牌
	if !isYaochupai(dr.PairTile) {
		return false
	}
	for _, tile := range dr.KotsuTiles {
		if !isYaochupai(tile) {
			return false
		}
	}
	for _, tile := range dr.ShuntsuFirstTiles {
		if !isYaochupai(tile) && !isYaochupai(tile + 2) {
			return false
		}
	}
	for _, meld := range hi.Melds {
		tiles := meld.Tiles
		if meld.MeldType == model.MeldTypeChi {
			if !isYaochupai(tiles[0]) && !isYaochupai(tiles[2]) {
				return false
			}
		} else {
			if !isYaochupai(tiles[0]) {
				return false
			}
		}
	}
	return true
}

func (hi *_handInfo) chanta() bool {
	return hi.containHonor() && hi._chantai()
}

func (hi *_handInfo) junchan() bool {
	return !hi.containHonor() && hi._chantai()
}

func (hi *_handInfo) honroutou() bool {
	if len(hi.Melds) == 0 {
		// 门清时简单判断
		cnt := 0
		for _, tile := range YaochuTiles {
			cnt += hi.HandTiles34[tile]
		}
		return cnt == 14
	}
	if !hi.containHonor() {
		return false
	}
	dr := hi.divideResult
	// 不能有顺子
	if len(dr.ShuntsuFirstTiles) > 0 {
		return false
	}
	for _, meld := range hi.Melds {
		if meld.MeldType == model.MeldTypeChi {
			return false
		}
	}
	// 所有雀头和刻子都要包含幺九牌
	if !isYaochupai(dr.PairTile) {
		return false
	}
	for _, tile := range dr.KotsuTiles {
		if !isYaochupai(tile) {
			return false
		}
	}
	for _, meld := range hi.Melds {
		if !isYaochupai(meld.Tiles[0]) {
			return false
		}
	}
	return true
}

func (hi *_handInfo) shousangen() bool {
	// 检查雀头
	dr := hi.divideResult
	if dr.PairTile < 31 {
		return false
	}
	// 检查三元牌刻子个数
	cnt := 0
	for _, tile := range dr.KotsuTiles {
		if tile >= 31 {
			cnt++
		}
	}
	for _, meld := range hi.Melds {
		if meld.MeldType != model.MeldTypeChi && meld.Tiles[0] >= 31 {
			cnt++
		}
	}
	return cnt == 2
}

func (hi *_handInfo) _numSuit() int {
	cntMan := 0
	cntPin := 0
	cntSou := 0
	cnt := func(tile int) {
		if isMan(tile) {
			cntMan++
		} else if isPin(tile) {
			cntPin++
		} else if isSou(tile) {
			cntSou++
		}
	}

	if hi.divideResult.IsChiitoi {
		// 七对子特殊判断
		for i, c := range hi.HandTiles34 {
			if c > 0 {
				cnt(i)
			}
		}
	} else {
		dr := hi.divideResult
		cnt(dr.PairTile)
		for _, tile := range dr.KotsuTiles {
			cnt(tile)
		}
		for _, tile := range dr.ShuntsuFirstTiles {
			cnt(tile)
		}
		for _, meld := range hi.Melds {
			cnt(meld.Tiles[0])
		}
	}

	numSuit := 0
	if cntMan > 0 {
		numSuit++
	}
	if cntPin > 0 {
		numSuit++
	}
	if cntSou > 0 {
		numSuit++
	}
	return numSuit
}

func (hi *_handInfo) honitsu() bool {
	return hi.containHonor() && hi._numSuit() == 1
}

func (hi *_handInfo) chinitsu() bool {
	return !hi.containHonor() && hi._numSuit() == 1
}

type yakuChecker func(*_handInfo) bool

var yakuCheckerMap = map[int]yakuChecker{
	YakuDaburii:        (*_handInfo).daburii,
	YakuRiichi:         (*_handInfo).riichi,
	YakuChiitoi:        (*_handInfo).chiitoi,
	YakuTsumo:          (*_handInfo).tsumo,
	YakuPinfu:          (*_handInfo).pinfu,
	YakuRyanpeikou:     (*_handInfo).ryanpeikou,
	YakuIipeikou:       (*_handInfo).iipeikou,
	YakuSanshokuDoujun: (*_handInfo).sanshokuDoujun,
	YakuIttsuu:         (*_handInfo).ittsuu,
	YakuToitoi:         (*_handInfo).toitoi,
	YakuSanAnkou:       (*_handInfo).sanAnkou,
	YakuSanshokuDoukou: (*_handInfo).sanshokuDoukou,
	YakuSanKantsu:      (*_handInfo).sanKantsu,
	YakuTanyao:         (*_handInfo).tanyao,
	YakuChanta:         (*_handInfo).chanta,
	YakuJunchan:        (*_handInfo).junchan,
	YakuHonroutou:      (*_handInfo).honroutou,
	YakuShousangen:     (*_handInfo).shousangen,
	YakuHonitsu:        (*_handInfo).honitsu,
	YakuChinitsu:       (*_handInfo).chinitsu,
}

func findYakuTypes(hi *_handInfo) (yakuTypes []int) {
	// TODO: 先检测是否有役满

	var yakuHanMap _yakuHanMap
	if !hi.IsNaki() {
		yakuHanMap = YakuHanMap
	} else {
		yakuHanMap = NakiYakuHanMap
	}

	for yakuType := range yakuHanMap {
		if checker, ok := yakuCheckerMap[yakuType]; ok {
			if checker(hi) {
				yakuTypes = append(yakuTypes, yakuType)
			}
		}
	}

	// 役牌单独算（连风算两个）
	numYakuhai := hi.numYakuhai()
	for i := 0; i < numYakuhai; i++ {
		yakuTypes = append(yakuTypes, YakuYakuhai)
	}

	sort.Ints(yakuTypes)
	return
}

// 寻找所有可能的役种
// 调用前请设置 WinTile
func FindAllYakuTypes(playerInfo *model.PlayerInfo) (yakuTypes []int) {
	canYaku := make([]bool, maxYakuType)
	for _, result := range DivideTiles34(playerInfo.HandTiles34) {
		_hi := &_handInfo{
			PlayerInfo:   playerInfo,
			divideResult: result,
		}
		yakuTypes := findYakuTypes(_hi)
		for _, t := range yakuTypes {
			canYaku[t] = true
		}
	}
	for yakuType, isYaku := range canYaku {
		if isYaku {
			yakuTypes = append(yakuTypes, yakuType)
		}
	}
	return
}
