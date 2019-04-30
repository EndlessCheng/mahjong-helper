package util

type YakuType int

const (
	// https://en.wikipedia.org/wiki/Japanese_Mahjong_yaku
	// Special criteria
	YakuRiichi YakuType = iota
	YakuChiitoi

	// Yaku based on luck
	YakuTsumo
	YakuIppatsu
	YakuHaitei
	YakuHoutei
	YakuRinshan
	YakuChankan
	YakuDaburii

	// Yaku based on sequences
	YakuPinfu
	YakuRyanpeikou
	YakuIipeikou
	YakuSanshokuDoujun  // *
	YakuIttsuu          // *

	// Yaku based on triplets and/or quads
	YakuToitoi
	YakuSanAnkou  // TODO 注意荣和的刻子是明刻
	YakuSanshokuDoukou
	YakuSanKantsu

	// Yaku based on terminal or honor tiles
	YakuTanyao
	YakuYakuhai
	YakuChanta     // * 必须有顺子
	YakuJunchan    // * 必须有顺子
	YakuHonroutou  // 七对也算
	YakuShousangen

	// Yaku based on suits
	YakuHonitsu   // *
	YakuChinitsu  // *

	// TODO: 役满
)

var YakuHanMap = map[YakuType]int{
	YakuRiichi:  1,
	YakuChiitoi: 2,

	YakuTsumo:   1,
	YakuIppatsu: 1,
	YakuHaitei:  1,
	YakuHoutei:  1,
	YakuRinshan: 1,
	YakuChankan: 1,
	YakuDaburii: 2,

	YakuPinfu:          1,
	YakuRyanpeikou:     3,
	YakuIipeikou:       1,
	YakuSanshokuDoujun: 2,
	YakuIttsuu:         2,

	YakuToitoi:         2,
	YakuSanAnkou:       2,
	YakuSanshokuDoukou: 2,
	YakuSanKantsu:      2,

	YakuTanyao:     1,
	YakuYakuhai:    1,
	YakuChanta:     2,
	YakuJunchan:    3,
	YakuHonroutou:  2,
	YakuShousangen: 2,

	YakuHonitsu:  3,
	YakuChinitsu: 6,
}

var YakumanTimesMap = map[YakuType]int{

}

//type yakuCheckerFunc func(hi *HandInfo, dr DivideResult, melds [][]int) bool

func (hi *HandInfo) isYakuTile(tile int) bool {
	return tile >= 31 || tile == hi.RoundWindTile || tile == hi.SelfWindTile
}

func (hi *HandInfo) chiitoi() bool {
	return hi.Divide.IsChiitoi
}

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

func (hi *HandInfo) ryanpeikou() bool {
	return hi.Divide.IsRyanpeikou
}

// 两杯口时无一杯口
func (hi *HandInfo) iipeikou() bool {
	return hi.Divide.IsIipeikou
}

func (hi *HandInfo) sanshokuDoujun() bool {
	drs := hi.Divide.ShuntsuFirstTiles
	if len(drs) < 3 {
		return false
	}
	var sMan, sPin, sSou []int
	for _, s := range drs {
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
	drs := hi.Divide.ShuntsuFirstTiles
	if len(drs) < 3 {
		return false
	}
	uniqueS := []int{drs[0]}
	for _, s := range drs[1:] {
		if s != uniqueS[len(uniqueS)-1] {
			uniqueS = append(uniqueS, s)
		}
	}
	if len(uniqueS) < 3 {
		return false
	}
	return false
	//return drs[0]%9 == 0 && drs[0] == drs[1]-3 && drs[0] == drs[2]-6 ||
	//	len(uniqueS) == 4 && drs[1]%9 == 0 && drs[1] == drs[2]-3 && drs[1] == drs[3]-6
}

func (hi *HandInfo) toitoi() bool {
	return len(hi.Divide.KotsuTiles) == 4
}

//func (hi *HandInfo) sanAnkou() bool {
//	cntAnkou := 0
//	for
//	return len(hi.Divide.KotsuTiles) == 4
//}

func (hi *HandInfo) sanshokuDoukou() bool {
	drk := hi.Divide.KotsuTiles
	if len(drk) < 3 {
		return false
	}
	var kMan, kPin, kSou []int
	for _, tile := range drk {
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

func FindYakuList(hi *HandInfo) (yakuList []YakuType) {

	return
}
