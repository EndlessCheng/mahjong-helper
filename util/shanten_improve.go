package util

import (
	"fmt"
	"sort"
	"math"
)

// map[改良牌]进张（选择进张数最大的）
type Improves map[int]Waits

// 1/4/7/10/13 张手牌的分析结果
type WaitsWithImproves13 struct {
	// 原手牌
	Tiles34 []int

	// 向听数
	Shanten int

	// 进张：摸到这张牌能让向听数前进
	Waits Waits

	// TODO: 鸣牌进张：他家打出这张牌，可以鸣牌，且能让向听数前进
	//MeldWaits Waits

	// map[进张牌]向听前进后的进张数（这里让向听前进的切牌选择的是使「向听前进后的进张数最大」的切牌）
	NextShantenWaitsCountMap map[int]int

	// 向听前进后的进张数的加权均值
	AvgNextShantenWaitsCount float64

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

	// TODO: 役种提醒
	// TODO: 赤牌改良提醒
	// TODO: 打点均值？
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
		mixedScore := r.AvgImproveWaitsCount * r.AvgNextShantenWaitsCount
		for i := 2; i <= r.Shanten; i++ {
			mixedScore /= 4
		}
		s += fmt.Sprintf(" %.2f %s进张（%.2f 综合分）",
			r.AvgNextShantenWaitsCount,
			NumberToChineseShanten(r.Shanten-1),
			mixedScore,
		)
	}
	if r.Shanten >= 0 && r.Shanten <= 1 {
		s += fmt.Sprintf("（%.2f%% 参考和率）", r.AvgAgariRate)
	}
	return s
}

// 1/4/7/10/13 张牌，计算向听数、进张（考虑了剩余枚数）
func CalculateShantenAndWaits13(tiles34 []int, leftTiles34 []int, isOpen bool) (shanten int, waits Waits) {
	if len(leftTiles34) == 0 {
		leftTiles34 = InitLeftTiles34WithTiles34(tiles34)
	}

	shanten = CalculateShanten(tiles34, isOpen)

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

		if leftTiles34[i] == 0 {
			// 无法摸到这张牌
			continue
		}

		// 摸牌
		tiles34[i]++
		if newShanten := CalculateShanten(tiles34, isOpen); newShanten < shanten {
			// 向听前进了，则换的这张牌为进张，进张数即剩余枚数
			waits[i] = leftTiles34[i]
		}
		tiles34[i]--
	}

	return
}

// 1/4/7/10/13 张牌，计算向听数、进张、改良等（考虑了剩余枚数）
func CalculateShantenWithImproves13(tiles34 []int, leftTiles34 []int, isOpen bool) (r *WaitsWithImproves13) {
	if len(leftTiles34) == 0 {
		leftTiles34 = InitLeftTiles34WithTiles34(tiles34)
	}

	shanten13, waits := CalculateShantenAndWaits13(tiles34, leftTiles34, isOpen)
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
				if newShanten13, newWaits := CalculateShantenAndWaits13(tiles34, leftTiles34, isOpen); newShanten13 < shanten13 {
					// 切牌一般切进张最多的
					if waitsCount := newWaits.AllCount(); waitsCount > nextShantenWaitsCountMap[i] {
						nextShantenWaitsCountMap[i] = waitsCount
					}
					// 听牌一般切和率最高的，TODO: 除非打点更高
					// TODO: add selfDiscards
					if newShanten13 == 0 {
						maxAgariRate = math.Max(maxAgariRate, CalculateAgariRate(newWaits, nil))
					}
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
				if newShanten13, improveWaits := CalculateShantenAndWaits13(tiles34, leftTiles34, isOpen); newShanten13 == shanten13 {
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
		// TODO: add selfDiscards
		avgAgariRate = CalculateAgariRate(waits, nil)
	}

	_tiles34 := make([]int, 34)
	copy(_tiles34, tiles34)
	r = &WaitsWithImproves13{
		Tiles34:                  _tiles34,
		Shanten:                  shanten13,
		Waits:                    waits,
		NextShantenWaitsCountMap: nextShantenWaitsCountMap,
		Improves:                 improves,
		ImproveWayCount:          improveWayCount,
		AvgImproveWaitsCount:     float64(waitsCount),
		AvgAgariRate:             avgAgariRate,
	}

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

func (l WaitsWithImproves14List) Sort() {
	sort.Slice(l, func(i, j int) bool {
		ri, rj := l[i].Result13, l[j].Result13

		// 听牌的话直接按照和率排序，TODO: 考虑打点
		if l[0].Shanten == 0 {
			return ri.AvgAgariRate > rj.AvgAgariRate
		}

		// 改良*前进后的进张 - 改良 - 前进后的进张 - 进张 - 和率

		// 对于一向听，考虑到未听牌之前要听的牌会被他家打出而造成听牌时的枚数降低，所以听牌枚数比和率更重要

		riM, rjM := ri.AvgImproveWaitsCount*ri.AvgNextShantenWaitsCount, rj.AvgImproveWaitsCount*rj.AvgNextShantenWaitsCount
		if riM != rjM {
			return riM > rjM
		}

		if ri.AvgImproveWaitsCount != rj.AvgImproveWaitsCount {
			return ri.AvgImproveWaitsCount > rj.AvgImproveWaitsCount
		}

		if ri.AvgNextShantenWaitsCount != rj.AvgNextShantenWaitsCount {
			return ri.AvgNextShantenWaitsCount > rj.AvgNextShantenWaitsCount
		}

		riWaitsCount, rjWaitsCount := ri.Waits.AllCount(), rj.Waits.AllCount()
		if riWaitsCount != rjWaitsCount {
			return riWaitsCount > rjWaitsCount
		}

		if ri.AvgAgariRate != rj.AvgAgariRate {
			return ri.AvgAgariRate > rj.AvgAgariRate
		}

		// 改良种类、方式多的优先
		if len(ri.Improves) != len(rj.Improves) {
			return len(ri.Improves) > len(rj.Improves)
		}
		if ri.ImproveWayCount != rj.ImproveWayCount {
			return ri.ImproveWayCount > rj.ImproveWayCount
		}

		return l[i].DiscardTile > l[j].DiscardTile
	})
}

//func (l WaitsWithImproves14List) FilterWithLeftTiles34(leftTiles34 []int) {
//	for _, r := range l {
//		r.Result13.Waits.FixCountsWithLeftCounts(leftTiles34)
//	}
//	l.Sort()
//}

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
func CalculateShantenWithImproves14(tiles34 []int, leftTiles34 []int, isOpen bool) (shanten int, waitsWithImproves WaitsWithImproves14List, incShantenResults WaitsWithImproves14List) {
	if len(leftTiles34) == 0 {
		leftTiles34 = InitLeftTiles34WithTiles34(tiles34)
	}

	shanten = CalculateShanten(tiles34, isOpen)

	for i := 0; i < 34; i++ {
		if tiles34[i] == 0 {
			continue
		}
		tiles34[i]-- // 切牌
		result13 := CalculateShantenWithImproves13(tiles34, leftTiles34, isOpen)
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
	waitsWithImproves.Sort()
	incShantenResults.Sort()
	return
}

// 计算最小向听数，鸣牌方式，该鸣牌方式下的向听数
func calculateMeldShanten(tiles34 []int, tile int, allowChi bool) (minShanten int, combinations [][]int, shantens []int) {
	// 是否能碰
	if tiles34[tile] >= 2 {
		combinations = append(combinations, []int{tile, tile})
	}
	// 是否能吃
	if allowChi && tile < 27 {
		checkChi := func(tileA, tileB int) {
			if tiles34[tileA] > 0 && tiles34[tileB] > 0 {
				combinations = append(combinations, []int{tileA, tileB})
			}
		}
		t9 := tile % 9
		if t9 >= 2 {
			checkChi(tile-2, tile-1)
		}
		if t9 >= 1 && t9 <= 7 {
			checkChi(tile-1, tile+1)
		}
		if t9 <= 6 {
			checkChi(tile+1, tile+2)
		}
	}

	// 计算所有副露情况下的最小向听数
	minShanten = 99
	for _, c := range combinations {
		tiles34[c[0]]--
		tiles34[c[1]]--
		shanten := CalculateShanten(tiles34, true)
		minShanten = MinInt(minShanten, shanten)
		shantens = append(shantens, shanten)
		tiles34[c[0]]++
		tiles34[c[1]]++
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

func CalculateMeld(tiles34 []int, tile int, allowChi bool, leftTiles34 []int) (shanten int, waitsWithImproves WaitsWithImproves14List, incShantenResults WaitsWithImproves14List) {
	if len(leftTiles34) == 0 {
		leftTiles34 = InitLeftTiles34WithTiles34(tiles34)
	}

	shanten, combinations, _ := calculateMeldShanten(tiles34, tile, allowChi)

	for _, c := range combinations {
		tiles34[c[0]]--
		tiles34[c[1]]--
		_shanten, _waitsWithImproves, _incShantenResults := CalculateShantenWithImproves14(tiles34, leftTiles34, true)
		tiles34[c[0]]++
		tiles34[c[1]]++

		// 去掉现物食替的情况
		_waitsWithImproves.filterOutDiscard(tile)
		_incShantenResults.filterOutDiscard(tile)

		// 去掉筋食替的情况
		cantDiscardTile := -1
		if c[0] < tile && c[1] < tile && tile >= 3 {
			cantDiscardTile = tile - 3
		} else if c[0] > tile && c[1] > tile && tile <= 5 {
			cantDiscardTile = tile + 3
		}
		if cantDiscardTile != -1 {
			_waitsWithImproves.filterOutDiscard(cantDiscardTile)
			_incShantenResults.filterOutDiscard(cantDiscardTile)
		}

		_waitsWithImproves.addOpenTile(c[:])
		_incShantenResults.addOpenTile(c[:])

		if _shanten == shanten {
			waitsWithImproves = append(waitsWithImproves, _waitsWithImproves...)
			incShantenResults = append(incShantenResults, _incShantenResults...)
		} else if _shanten == shanten+1 {
			incShantenResults = append(incShantenResults, _waitsWithImproves...)
		}
	}

	waitsWithImproves.Sort()
	incShantenResults.Sort()

	return
}
