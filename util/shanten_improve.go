package util

import (
	"fmt"
	"sort"
	"math"
	"github.com/EndlessCheng/mahjong-helper/util/model"
)

// map[改良牌]进张（选择进张数最大的）
type Improves map[int]Waits

// 1/4/7/10/13 张手牌的分析结果
type WaitsWithImproves13 struct {
	// 场风
	RoundWindTile34 int

	// 自风
	SelfWindTile34 int

	// 原手牌
	Tiles34 []int

	// 剩余牌
	LeftTiles34 []int

	// 向听数
	Shanten int

	// 进张：摸到这张牌能让向听数前进
	// 考虑了剩余枚数
	// 若某个进张 4 枚都可见，则该进张的 value 值为 0
	Waits Waits

	// TODO: 鸣牌进张：他家打出这张牌，可以鸣牌，且能让向听数前进
	//MeldWaits Waits

	// map[进张牌]向听前进后的进张数（这里让向听前进的切牌选择的是使「向听前进后的进张数最大」的切牌）
	NextShantenWaitsCountMap map[int]int

	// 向听前进后的进张数的加权均值
	AvgNextShantenWaitsCount float64

	// 综合了进张与向听前进后进张的评分
	MixedWaitsScore float64

	// 改良：摸到这张牌虽不能让向听数前进，但可以让进张变多
	// len(Improves) 即为改良的牌的种数
	Improves Improves

	// 改良情况数，这里计算的是有多少种使进张增加的切牌方式
	ImproveWayCount int

	// 在没有摸到进张时的改良后进张数的加权均值（计算时，对于既不是进张也不是改良的牌，其进张数为 Waits.AllCount()）
	// 这里只考虑一巡的改良均值
	// TODO: 在考虑改良的情况下，如何计算向听前进所需要的摸牌次数的期望值？
	AvgImproveWaitsCount float64

	// 向听前进后，若听牌，其最大和率的加权均值
	// 若已听牌，则该值为当前手牌和率
	AvgAgariRate float64

	// 振听可能率（一向听和听牌时）
	FuritenRate float64

	// 役种
	YakuTypes []int

	// 宝牌个数（手牌+副露）
	DoraCount int

	// 荣和打点期望
	RonPoint float64

	// 自摸打点期望
	TsumoPoint float64

	// TODO: 赤牌改良提醒
}

// waitsCount := float64(r.Waits.AllCount())
// return (waitsCount + (1-waitsCount/leftCount)*r.AvgImproveWaitsCount) * 100 / leftCount

// 进张和向听前进后进张的评分
// 这里粗略地近似为向听前进两次的概率
func (r *WaitsWithImproves13) mixedWaitsScore() float64 {
	if r.Waits.AllCount() == 0 || r.AvgNextShantenWaitsCount == 0 {
		return 0
	}
	leftCount := float64(CountOfTiles34(r.LeftTiles34))
	p2 := float64(r.Waits.AllCount()) / leftCount
	//p2 := r.AvgImproveWaitsCount / leftCount
	p1 := r.AvgNextShantenWaitsCount / leftCount
	//if r.AvgAgariRate > 0 {
	//	p1 = r.AvgAgariRate / 100
	//}
	p2_, p1_ := 1-p2, 1-p1
	const leftTurns = 10.0 // math.Max(5.0, leftCount/4)
	sumP2 := p2_ * (1 - math.Pow(p2_, leftTurns)) / p2
	sumP1 := p1_ * (1 - math.Pow(p1_, leftTurns)) / p1
	result := p2 * p1 * (sumP2 - sumP1) / (p2_ - p1_)
	return result * 100
}

// 调试用
func (r *WaitsWithImproves13) String() string {
	s := fmt.Sprintf("%d 进张 %s\n%.2f 改良进张 [%d(%d) 种]",
		r.Waits.AllCount(),
		//r.Waits.AllCount()+r.MeldWaits.AllCount(),
		TilesToStrWithBracket(r.Waits.indexes()),
		r.AvgImproveWaitsCount,
		len(r.Improves),
		r.ImproveWayCount,
	)
	if r.Shanten >= 1 {
		mixedScore := r.MixedWaitsScore
		//for i := 2; i <= r.Shanten; i++ {
		//	mixedScore /= 4
		//}
		s += fmt.Sprintf(" %.2f %s进张（%.2f 综合分）",
			r.AvgNextShantenWaitsCount,
			NumberToChineseShanten(r.Shanten-1),
			mixedScore,
		)
	}
	if r.Shanten >= 0 && r.Shanten <= 1 {
		s += fmt.Sprintf("（%.2f%% 参考和率）", r.AvgAgariRate)
		if r.FuritenRate > 0 {
			if r.FuritenRate < 1 {
				s += "[可能振听]"
			} else {
				s += "[振听]"
			}
		}
		s += YakuTypesWithDoraToStr(r.YakuTypes, r.DoraCount)
	}
	if r.RonPoint > 0 {
		s += fmt.Sprintf("[%d荣和点数]", int(math.Round(r.RonPoint)))
	}
	if r.TsumoPoint > 0 {
		s += fmt.Sprintf("[%d自摸点数]", int(math.Round(r.TsumoPoint)))
	}
	return s
}

// 1/4/7/10/13 张牌，计算向听数、进张（考虑了剩余枚数）
func CalculateShantenAndWaits13(tiles34 []int, leftTiles34 []int) (shanten int, waits Waits) {
	if len(leftTiles34) == 0 {
		leftTiles34 = InitLeftTiles34WithTiles34(tiles34)
	}

	shanten = CalculateShanten(tiles34)

	// 剪枝：检测非浮牌，在不考虑国士无双的情况下，这种牌是不可能让向听数前进的（但有改良的可能，不过 CalculateShantenAndWaits13 函数不考虑这个）
	// 此处优化提升了约 30% 的性能
	//needCheck34 := make([]bool, 34)
	//idx := -1
	//for i := 0; i < 3; i++ {
	//	for j := 0; j < 9; j++ {
	//		idx++
	//		if tiles34[idx] == 0 {
	//			continue
	//		}
	//		if j == 0 {
	//			needCheck34[idx] = true
	//			needCheck34[idx+1] = true
	//			needCheck34[idx+2] = true
	//		} else if j == 1 {
	//			needCheck34[idx-1] = true
	//			needCheck34[idx] = true
	//			needCheck34[idx+1] = true
	//			needCheck34[idx+2] = true
	//		} else if j < 7 {
	//			needCheck34[idx-2] = true
	//			needCheck34[idx-1] = true
	//			needCheck34[idx] = true
	//			needCheck34[idx+1] = true
	//			needCheck34[idx+2] = true
	//		} else if j == 7 {
	//			needCheck34[idx-2] = true
	//			needCheck34[idx-1] = true
	//			needCheck34[idx] = true
	//			needCheck34[idx+1] = true
	//		} else {
	//			needCheck34[idx-2] = true
	//			needCheck34[idx-1] = true
	//			needCheck34[idx] = true
	//		}
	//	}
	//}
	//for i := 27; i < 34; i++ {
	//	if tiles34[i] > 0 {
	//		needCheck34[i] = true
	//	}
	//}

	waits = Waits{}
	for i := 0; i < 34; i++ {
		//if !needCheck34[i] {
		//	continue
		//}

		if tiles34[i] == 4 {
			// 无法摸到这张牌
			continue
		}

		// 摸牌
		tiles34[i]++
		if newShanten := CalculateShanten(tiles34); newShanten < shanten {
			// 向听前进了，则换的这张牌为进张，进张数即剩余枚数
			// 有可能为 0，但这对于判断振听是有帮助的，所以记录
			waits[i] = leftTiles34[i]
		}
		tiles34[i]--
	}

	return
}

// 1/4/7/10/13 张牌，计算向听数、进张、改良等（考虑了剩余枚数）
func CalculateShantenWithImproves13(playerInfo *model.PlayerInfo) (r *WaitsWithImproves13) {
	if len(playerInfo.LeftTiles34) == 0 {
		playerInfo.FillLeftTiles34()
	}

	tiles34 := playerInfo.HandTiles34
	leftTiles34 := playerInfo.LeftTiles34
	shanten13, waits := CalculateShantenAndWaits13(tiles34, leftTiles34)
	waitsCount := waits.AllCount()

	nextShantenWaitsCountMap := map[int]int{} // map[进张牌]听多少张牌
	improves := Improves{}
	improveWayCount := 0
	// 对于每张牌，摸到之后的手牌进张数（如果摸到的是 waits 中的牌，则进张数视作 waitsCount）
	improveWaitsCount34 := make([]int, 34)
	// 初始化成基本进张
	for i := 0; i < 34; i++ {
		improveWaitsCount34[i] = waitsCount
	}
	avgAgariRate := 0.0

	canYaku := make([]bool, maxYakuType)
	fillYakuTypes := func(_shanten13 int, _waits Waits) {
		if _shanten13 != 0 { // 只考虑听牌
			return
		}
		for tile, left := range _waits {
			if left == 0 {
				continue
			}
			tiles34[tile]++
			playerInfo.WinTile = tile
			_yakuTypes := FindAllYakuTypes(playerInfo)
			for _, t := range _yakuTypes {
				canYaku[t] = true
			}
			tiles34[tile]--
		}
	}

	// 若听牌，计算役种
	fillYakuTypes(shanten13, waits)

	for i := 0; i < 34; i++ {
		if leftTiles34[i] == 0 {
			// 无法摸到这张牌
			continue
		}
		// 从剩余牌中摸牌
		leftTiles34[i]--
		tiles34[i]++
		if _, ok := waits[i]; ok {
			// 摸到的是进张
			maxAgariRate := 0.0
			for j := 0; j < 34; j++ {
				if tiles34[j] == 0 || j == i {
					continue
				}
				// 切牌
				tiles34[j]--
				// 向听前进才是正确的切牌
				if newShanten13, newWaits := CalculateShantenAndWaits13(tiles34, leftTiles34); newShanten13 < shanten13 {
					// 切牌一般切进张最多的
					if waitsCount := newWaits.AllCount(); waitsCount > nextShantenWaitsCountMap[i] {
						nextShantenWaitsCountMap[i] = waitsCount
					}
					// 听牌一般切和率最高的，TODO: 除非打点更高，比如说听到 dora 上，或者有三色等
					if newShanten13 == 0 {
						maxAgariRate = math.Max(maxAgariRate, CalculateAgariRate(newWaits, playerInfo.DiscardTiles))
					}
					// 若前进后听牌（当前为一向听），计算役种
					fillYakuTypes(newShanten13, newWaits)
				}
				tiles34[j]++
			}
			// 加权：进张牌的剩余枚数*和率
			avgAgariRate += float64(leftTiles34[i]+1) * maxAgariRate
		} else {
			// 摸到的不是进张，但可能有改良
			for j := 0; j < 34; j++ {
				if tiles34[j] == 0 || j == i {
					continue
				}
				// 切牌
				tiles34[j]--
				// 正确的切牌
				if newShanten13, improveWaits := CalculateShantenAndWaits13(tiles34, leftTiles34); newShanten13 == shanten13 {
					// 若进张数变多，则为改良
					if improveWaitsCount := improveWaits.AllCount(); improveWaitsCount > waitsCount {
						improveWayCount++
						if improveWaitsCount > improveWaitsCount34[i] {
							improveWaitsCount34[i] = improveWaitsCount
							// improves 选的是进张数最大的改良
							improves[i] = improveWaits
						}
						//fmt.Println(fmt.Sprintf("    摸 %s 切 %s 改良:", MahjongZH[i], MahjongZH[j]), improveWaitsCount, TilesToStrWithBracket(improveWaits.indexes()))
					}
				}
				tiles34[j]++
			}
		}
		tiles34[i]--
		leftTiles34[i]++
	}
	avgAgariRate /= float64(waitsCount)
	if shanten13 == 0 {
		avgAgariRate = CalculateAgariRate(waits, playerInfo.DiscardTiles)
	}

	yakuTypes := []int{}
	for yakuType, can := range canYaku {
		if can {
			yakuTypes = append(yakuTypes, yakuType)
		}
	}

	_tiles34 := make([]int, 34)
	copy(_tiles34, tiles34)
	r = &WaitsWithImproves13{
		RoundWindTile34: playerInfo.RoundWindTile,
		SelfWindTile34:  playerInfo.SelfWindTile,
		Tiles34:         _tiles34,
		LeftTiles34:     leftTiles34,
		Shanten:         shanten13,
		Waits:           waits,
		NextShantenWaitsCountMap: nextShantenWaitsCountMap,
		Improves:                 improves,
		ImproveWayCount:          improveWayCount,
		AvgImproveWaitsCount:     float64(waitsCount),
		AvgAgariRate:             avgAgariRate,
		YakuTypes:                yakuTypes,
		DoraCount:                playerInfo.DoraCount,
	}

	// 对于听牌及一向听，判断是否有振听可能
	if shanten13 <= 1 {
		for _, discardTile := range playerInfo.DiscardTiles {
			if _, ok := waits[discardTile]; ok {
				r.FuritenRate = 0.5 // TODO: 待完善
				if shanten13 == 0 {
					// 听牌时，若听的牌在舍牌中，则构成振听
					r.FuritenRate = 1
					// 修正振听时的和率
					r.AvgAgariRate *= FuritenAgariMulti
				}
			}
		}
	}

	// 非振听且待牌有役时计算荣和点数
	// TODO: 立直时考虑中里的分数
	if r.FuritenRate == 0 && shanten13 == 0 {
		sum := 0
		w := 0
		for tile, left := range waits {
			if left == 0 {
				continue
			}
			tiles34[tile]++
			playerInfo.WinTile = tile
			ronPoint := CalcRonPointWithHands(playerInfo)
			// 不考虑无役（如后附，片听）
			if ronPoint > 0 {
				sum += ronPoint * left
				w += left
			}
			tiles34[tile]--
		}
		if w > 0 {
			r.RonPoint = float64(sum) / float64(w)
		}
	}

	// TODO: 自摸点数
	// TODO: 立直时考虑中里的分数

	// 分析
	if len(nextShantenWaitsCountMap) > 0 {
		nextShantenWaitsSum := 0
		weight := 0
		for tile, c := range nextShantenWaitsCountMap {
			w := leftTiles34[tile]
			nextShantenWaitsSum += w * c
			weight += w
		}
		r.AvgNextShantenWaitsCount = float64(nextShantenWaitsSum) / float64(weight)
	}
	if len(improves) > 0 {
		improveWaitsSum := 0
		weight := 0
		for i := 0; i < 34; i++ {
			w := leftTiles34[i]
			improveWaitsSum += w * improveWaitsCount34[i]
			weight += w
		}
		r.AvgImproveWaitsCount = float64(improveWaitsSum) / float64(weight)
	}
	r.MixedWaitsScore = r.mixedWaitsScore()

	return
}

type WaitsWithImproves14 struct {
	// 切牌后的手牌分析结果
	Result13 *WaitsWithImproves13
	// 需要切的牌
	DiscardTile int
	// 切掉这张牌后的向听数
	Shanten int
	// 副露信息（没有副露就是 nil）
	// 比如用 23m 吃了牌，OpenTiles 就是 [1,2]
	OpenTiles []int
}

func (r *WaitsWithImproves14) String() string {
	meldInfo := ""
	if len(r.OpenTiles) > 0 {
		meldType := "吃"
		if r.OpenTiles[0] == r.OpenTiles[1] {
			meldType = "碰"
		}
		meldInfo = fmt.Sprintf("用 %s%s %s，", string([]rune(MahjongZH[r.OpenTiles[0]])[:1]), MahjongZH[r.OpenTiles[1]], meldType)
	}
	return meldInfo + fmt.Sprintf("切 %s: %s", MahjongZH[r.DiscardTile], r.Result13.String())
}

type WaitsWithImproves14List []*WaitsWithImproves14

// 排序，若 needImprove 为 true，则优先按照 AvgImproveWaitsCount 排序
func (l WaitsWithImproves14List) Sort(needImprove bool) {
	sort.Slice(l, func(i, j int) bool {
		ri, rj := l[i].Result13, l[j].Result13

		// 听牌的话和率优先
		// TODO: 考虑打点
		if l[0].Shanten == 0 {
			if !Equal(ri.AvgAgariRate, rj.AvgAgariRate) {
				return ri.AvgAgariRate > rj.AvgAgariRate
			}
		}

		if needImprove {
			if !Equal(ri.AvgImproveWaitsCount, rj.AvgImproveWaitsCount) {
				return ri.AvgImproveWaitsCount > rj.AvgImproveWaitsCount
			}
		}

		// 排序规则：综合评分 - 进张 - 前进后的进张 - 和率 - 改良 - 好牌先走
		// 必须注意到的一点是，随着游戏的进行，进张会被他家打出，所以进张是有减少的趋势的
		// 对于一向听，考虑到未听牌之前要听的牌会被他家打出而造成听牌时的枚数降低，所以听牌枚数比和率更重要
		// 对比当前进张与前进后的进张，在二者乘积相近的情况下（注意这个前提），由于进张越大听牌速度越快，听牌时的进张数也就越接近预期进张数，所以进张越多越好（再次强调是在二者乘积相近的情况下）

		if !Equal(ri.MixedWaitsScore, rj.MixedWaitsScore) {
			return ri.MixedWaitsScore > rj.MixedWaitsScore
		}

		riWaitsCount, rjWaitsCount := ri.Waits.AllCount(), rj.Waits.AllCount()
		if riWaitsCount != rjWaitsCount {
			return riWaitsCount > rjWaitsCount
		}

		if !Equal(ri.AvgNextShantenWaitsCount, rj.AvgNextShantenWaitsCount) {
			return ri.AvgNextShantenWaitsCount > rj.AvgNextShantenWaitsCount
		}

		// shanten == 1
		if !Equal(ri.AvgAgariRate, rj.AvgAgariRate) {
			return ri.AvgAgariRate > rj.AvgAgariRate
		}

		if !Equal(ri.AvgImproveWaitsCount, rj.AvgImproveWaitsCount) {
			return ri.AvgImproveWaitsCount > rj.AvgImproveWaitsCount
		}

		idxI, idxJ := l[i].DiscardTile, l[j].DiscardTile
		if idxI >= 27 && idxJ >= 27 {
			// TODO 场风不为自风时：下家风 > 对家风 > 上家风 > 场风
		}

		// 好牌先走
		if idxI < 27 && idxJ < 27 {
			idxI %= 9
			if idxI > 4 {
				idxI = 8 - idxI
			}
			idxJ %= 9
			if idxJ > 4 {
				idxJ = 8 - idxJ
			}
			return idxI > idxJ
		}
		return idxI < idxJ

		//// 改良种类、方式多的优先
		//if len(ri.Improves) != len(rj.Improves) {
		//	return len(ri.Improves) > len(rj.Improves)
		//}
		//if ri.ImproveWayCount != rj.ImproveWayCount {
		//	return ri.ImproveWayCount > rj.ImproveWayCount
		//}
	})
}

func (l *WaitsWithImproves14List) filterOutDiscard(cantDiscardTile int) {
	newResults := WaitsWithImproves14List{}
	for _, r := range *l {
		if r.DiscardTile != cantDiscardTile {
			newResults = append(newResults, r)
		}
	}
	*l = newResults
}

func (l WaitsWithImproves14List) addOpenTile(openTiles []int) {
	for _, r := range l {
		r.OpenTiles = openTiles
	}
}

// 2/5/8/11/14 张牌，计算向听数、进张、改良、向听倒退等
func CalculateShantenWithImproves14(playerInfo *model.PlayerInfo) (shanten int, waitsWithImproves WaitsWithImproves14List, incShantenResults WaitsWithImproves14List) {
	if len(playerInfo.LeftTiles34) == 0 {
		playerInfo.FillLeftTiles34()
	}

	tiles34 := playerInfo.HandTiles34
	shanten = CalculateShanten(tiles34)

	for i := 0; i < 34; i++ {
		if tiles34[i] == 0 {
			continue
		}
		tiles34[i]-- // 切牌
		playerInfo.DiscardTiles = append(playerInfo.DiscardTiles, i)
		result13 := CalculateShantenWithImproves13(playerInfo)
		playerInfo.DiscardTiles = playerInfo.DiscardTiles[:len(playerInfo.DiscardTiles)-1]
		r := &WaitsWithImproves14{
			Result13:    result13,
			DiscardTile: i,
			Shanten:     result13.Shanten,
		}
		if result13.Shanten == shanten {
			waitsWithImproves = append(waitsWithImproves, r)
		} else {
			// 向听倒退
			incShantenResults = append(incShantenResults, r)
		}
		tiles34[i]++
	}

	needImprove := func(l []*WaitsWithImproves14) bool {
		if len(l) == 0 {
			return false
		}

		shanten := l[0].Shanten
		// 一向听及以下进张优先，改良其次
		if shanten <= 1 {
			return false
		}

		maxWaitsCount := 0
		for _, r := range waitsWithImproves {
			maxWaitsCount = MaxInt(maxWaitsCount, r.Result13.Waits.AllCount())
		}

		// 两向听及以上的七对子考虑改良
		return maxWaitsCount <= 6*shanten+3
	}

	ni := needImprove(waitsWithImproves)
	waitsWithImproves.Sort(ni)
	ni = needImprove(incShantenResults)
	incShantenResults.Sort(ni)
	return
}

// 计算最小向听数，鸣牌方式
func calculateMeldShanten(tiles34 []int, calledTile int, allowChi bool) (minShanten int, combinations []model.Meld) {
	// 是否能碰
	if tiles34[calledTile] >= 2 {
		combinations = append(combinations, model.Meld{
			MeldType:   model.MeldTypePon,
			Tiles:      []int{calledTile, calledTile, calledTile},
			SelfTiles:  []int{calledTile, calledTile},
			CalledTile: calledTile,
		})
	}
	// 是否能吃
	if allowChi && calledTile < 27 {
		checkChi := func(tileA, tileB int) {
			if tiles34[tileA] > 0 && tiles34[tileB] > 0 {
				_tiles := []int{tileA, tileB, calledTile}
				sort.Ints(_tiles)
				combinations = append(combinations, model.Meld{
					MeldType:   model.MeldTypeChi,
					Tiles:      _tiles,
					SelfTiles:  []int{tileA, tileB},
					CalledTile: calledTile,
				})
			}
		}
		t9 := calledTile % 9
		if t9 >= 2 {
			checkChi(calledTile-2, calledTile-1)
		}
		if t9 >= 1 && t9 <= 7 {
			checkChi(calledTile-1, calledTile+1)
		}
		if t9 <= 6 {
			checkChi(calledTile+1, calledTile+2)
		}
	}

	// 计算所有鸣牌下的最小向听数
	minShanten = 99
	for _, c := range combinations {
		tiles34[c.SelfTiles[0]]--
		tiles34[c.SelfTiles[1]]--
		shanten := CalculateShanten(tiles34)
		minShanten = MinInt(minShanten, shanten)
		tiles34[c.SelfTiles[0]]++
		tiles34[c.SelfTiles[1]]++
	}

	return
}

// TODO 鸣牌的情况判断（待重构）
// 编程时注意他家切掉的这张牌是否算到剩余数中
//if isOpen {
//if newShanten, combinations, shantens := calculateMeldShanten(tiles34, i, true); newShanten < shanten {
//	// 向听前进了，说明鸣牌成功，则换的这张牌为鸣牌进张
//	// 计算进张数：若能碰则 =剩余数*3，否则 =剩余数
//	meldWaits[i] = leftTile - tiles34[i]
//	for i, comb := range combinations {
//		if comb[0] == comb[1] && shantens[i] == newShanten {
//			meldWaits[i] *= 3
//			break
//		}
//	}
//}
//}

func CalculateMeld(playerInfo *model.PlayerInfo, calledTile int, allowChi bool) (shanten int, waitsWithImproves WaitsWithImproves14List, incShantenResults WaitsWithImproves14List) {
	if len(playerInfo.LeftTiles34) == 0 {
		playerInfo.FillLeftTiles34()
	}

	tiles34 := playerInfo.HandTiles34
	shanten, combinations := calculateMeldShanten(tiles34, calledTile, allowChi)

	for _, c := range combinations {
		tiles34[c.SelfTiles[0]]--
		tiles34[c.SelfTiles[1]]--
		playerInfo.Melds = append(playerInfo.Melds, c)
		_shanten, _waitsWithImproves, _incShantenResults := CalculateShantenWithImproves14(playerInfo)
		playerInfo.Melds = playerInfo.Melds[:len(playerInfo.Melds)-1]
		tiles34[c.SelfTiles[0]]++
		tiles34[c.SelfTiles[1]]++

		// 去掉现物食替的情况
		_waitsWithImproves.filterOutDiscard(calledTile)
		_incShantenResults.filterOutDiscard(calledTile)

		// 去掉筋食替的情况
		cantDiscardTile := -1
		if c.SelfTiles[0] < calledTile && c.SelfTiles[1] < calledTile && calledTile >= 3 {
			cantDiscardTile = calledTile - 3
		} else if c.SelfTiles[0] > calledTile && c.SelfTiles[1] > calledTile && calledTile <= 5 {
			cantDiscardTile = calledTile + 3
		}
		if cantDiscardTile != -1 {
			_waitsWithImproves.filterOutDiscard(cantDiscardTile)
			_incShantenResults.filterOutDiscard(cantDiscardTile)
		}

		_waitsWithImproves.addOpenTile(c.SelfTiles)
		_incShantenResults.addOpenTile(c.SelfTiles)

		if _shanten == shanten {
			waitsWithImproves = append(waitsWithImproves, _waitsWithImproves...)
			incShantenResults = append(incShantenResults, _incShantenResults...)
		} else if _shanten == shanten+1 {
			incShantenResults = append(incShantenResults, _waitsWithImproves...)
		}
	}

	waitsWithImproves.Sort(false)
	incShantenResults.Sort(false)

	return
}
