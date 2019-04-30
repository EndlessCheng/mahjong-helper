package util

const (
	// https://en.wikipedia.org/wiki/Japanese_Mahjong_yaku
	YakuChiitoi = iota

	YakuPinfu
	YakuIipeikou
	YakuRyanpeikou
	YakuSanshokuDoujun
	YakuIttsuu

	YakuToitoi
	YakuSanAnkou
	YakuSanshokuDoukou
	YakuSanKantsu  // TODO

	YakuTanyao
	YakuYakuhai
	YakuChanta
	YakuJunchan
	YakuHonroutou
	YakuShousangen

	YakuHonitsu
	YakuChinitsu

	// TODO: 役满
)

type handInfo struct {
	divideResult  DivideResult
	melds         [][]int
	winTile       int
	roundWindTile int
	selfWindTile  int
}

type yakuCheckerFunc func(hi *handInfo, dr DivideResult, melds [][]int) bool

func (hi *handInfo) isYakuTile(tile int) bool {
	return tile >= 31 || tile == hi.roundWindTile || tile == hi.selfWindTile
}

func (hi *handInfo) chiitoi() bool {
	return hi.divideResult.IsChiitoitsu()
}

func (hi *handInfo) pinfu() bool {
	// 雀头不能是役牌，且不能是单骑和牌
	dr := hi.divideResult
	if hi.isYakuTile(hi.winTile) || hi.winTile == dr.pairTile {
		return false
	}
	drs := hi.divideResult.ShuntsuFirstTiles
	if len(drs) < 4 {
		return false
	}
	for _, s := range drs {
		if hi.winTile == s {
			// 不能和 89 边张
			return s%9 <= 5
		}
		if hi.winTile == s+2 {
			// 不能和 12 边张
			return s%9 >= 1
		}
	}
	return false
}

func (hi *handInfo) ryanpeikou() bool {
	drs := hi.divideResult.ShuntsuFirstTiles
	if len(drs) < 4 {
		return false
	}
	return drs[0] == drs[1] && drs[2] == drs[3]
}

// 需要先判断是否为两杯口
func (hi *handInfo) iipeikou() bool {
	drs := hi.divideResult.ShuntsuFirstTiles
	if len(drs) < 2 {
		return false
	}
	for i := range drs[:len(drs)-1] {
		if drs[i] == drs[i+1] {
			return true
		}
	}
	return false
}

func (hi *handInfo) sanshokuDoujun() bool {
	drs := hi.divideResult.ShuntsuFirstTiles
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

func (hi *handInfo) ittsuu() bool {
	drs := hi.divideResult.ShuntsuFirstTiles
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

func (hi *handInfo) toitoi() bool {
	return len(hi.divideResult.KotsuTiles) == 4
}

//
//func (hi *handInfo) sanAnkou() bool {
//	cntAnkou := 0
//	for
//	return len(hi.divideResult.KotsuTiles) == 4
//}

//

//var yakuCheckerMap = map[int]yakuCheckerFunc{
//	//1: (*handInfo).Chiitoi,
//}
//
//var YakuMap = map[int]string{
//
//}
//

// 先找个三色看看~
func FindNormalYaku(hi *handInfo) bool {
	return hi.sanshokuDoujun()
}

func FindNormalYakuSimple(tiles34 []int) bool {
	for _, result := range DivideTiles34(tiles34) {
		if FindNormalYaku(&handInfo{divideResult: result}) {
			return true
		}
	}
	return false
}
