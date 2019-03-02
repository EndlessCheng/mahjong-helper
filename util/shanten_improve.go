package util

import (
	"fmt"
	"sort"
)

// map[改良牌]进张
type Improves map[int]Waits

// 13 张手牌的分析结果
type WaitsWithImproves13 struct {
	// 手牌
	Tiles34 []int

	// 向听数
	Shanten int

	// 进张：摸到这张牌可以让向听数前进
	Waits Waits

	// map[进张牌]向听前进后的进张数（这里让向听前进的切牌是最优切牌，即让向听前进后的进张数最大的切牌）
	NextShantenWaitsCountMap map[int]int

	// 改良：摸到这张牌虽不能让向听数前进，但可以让进张变多
	Improves Improves

	// 改良情况数
	ImproveWayCount int

	// 对于每张牌，摸到之后的手牌进张数（如果摸到的是 Waits 中的牌，则进张数视作摸到之前的进张数）
	ImproveWaitsCount34 []int

	// 在没有摸到进张时的改良进张数的加权均值
	AvgImproveWaitsCount float64

	// 向听前进后的进张数的加权均值
	AvgNextShantenWaitsCount float64
}

// avgImproveWaitsCount: 在没有摸到进张时的改良进张数的加权均值
func (r *WaitsWithImproves13) analysis() (avgImproveWaitsCount float64, avgNextShantenWaitsCount float64) {
	const leftTile = 4

	if len(r.Improves) > 0 {
		improveScore := 0
		weight := 0
		for i := 0; i < 34; i++ {
			w := leftTile - r.Tiles34[i]
			improveScore += w * r.ImproveWaitsCount34[i]
			weight += w
		}
		avgImproveWaitsCount = float64(improveScore) / float64(weight)
		r.AvgImproveWaitsCount = avgImproveWaitsCount
	} else {
		r.AvgImproveWaitsCount = float64(r.Waits.allCount())
	}

	nextShantenWaitsSum := 0
	weight := 0
	for tile, c := range r.NextShantenWaitsCountMap {
		w := leftTile - r.Tiles34[tile]
		nextShantenWaitsSum += w * c
		weight += w
	}
	avgNextShantenWaitsCount = float64(nextShantenWaitsSum) / float64(weight)
	r.AvgNextShantenWaitsCount = avgNextShantenWaitsCount

	return
}

// 调试用
func (r *WaitsWithImproves13) String() string {
	return fmt.Sprintf("%s\n%.2f [%d 改良] %.2f %s进张",
		r.Waits.String(),
		r.AvgImproveWaitsCount,
		r.ImproveWayCount,
		r.AvgNextShantenWaitsCount,
		NumberToChineseShanten(r.Shanten-1),
	)
}

// 13 张牌，计算向听数和进张
func CalculateShantenAndWaits13(tiles34 []int, isOpen bool) (shanten int, waits Waits) {
	shanten = CalculateShanten(tiles34, isOpen)

	const leftTile = 4

	// 剪枝：可以先计算出非相邻的牌，在不考虑国士无双的情况下，这种牌是不可能让向听数前进的（但有改良的可能，不过 CalculateShantenAndWaits13 函数不考虑这个）
	// 此处剪枝提高了约 60% 的性能！
	needChecks := make([]int, 0, 34)
	idx := -1
	for i := 0; i < 3; i++ {
		for j := 0; j < 9; j++ {
			idx++
			if tiles34[idx] == leftTile {
				continue
			}
			if j == 0 || j == 9 {
				if tiles34[idx] > 0 || tiles34[idx+1] > 0 {
					needChecks = append(needChecks, idx)
				}
			} else {
				if tiles34[idx-1] > 0 || tiles34[idx] > 0 || tiles34[idx+1] > 0 {
					needChecks = append(needChecks, idx)
				}
			}
		}
	}
	for i := 27; i < 34; i++ {
		if tiles34[i] == leftTile {
			continue
		}
		if tiles34[i] > 0 {
			needChecks = append(needChecks, i)
		}
	}

	waits = Waits{}
	for i := 0; i < 34; i++ {
		if tiles34[i] == 0 {
			continue
		}
		// 切掉其中一张牌
		tiles34[i]--
		for _, j := range needChecks {
			if j == i || tiles34[j] == leftTile {
				continue
			}
			if _, ok := waits[j]; ok {
				// 已经是进张的就不用考虑了
				continue
			}
			// 换成其他牌
			tiles34[j]++
			if newShanten := CalculateShanten(tiles34, isOpen); newShanten < shanten {
				// 向听前进了，则换的这张牌为进张
				waits[j] = leftTile - (tiles34[j] - 1)
			}
			tiles34[j]--
		}
		tiles34[i]++
	}
	return
}

// 13 张牌，计算向听数、进张、改良、向听倒退等
func CalculateShantenWithImproves13(tiles34 []int, isOpen bool) (waitsWithImproves *WaitsWithImproves13) {
	shanten := CalculateShanten(tiles34, isOpen)
	//fmt.Println(shanten)

	const leftTile = 4

	waits := Waits{}                          // 进张
	rawImprovesMap := map[int]Improves{}      // map[摸到idx]改良情况
	nextShantenWaitsCountMap := map[int]int{} // map[摸到idx]听多少张牌

	for i := 0; i < 34; i++ {
		if tiles34[i] == 0 {
			continue
		}
		improves := Improves{}
		// 切掉其中一张牌
		tiles34[i]--
		for j := 0; j < 34; j++ {
			if j == i || tiles34[j] == leftTile {
				continue
			}
			// 换成其他牌
			tiles34[j]++
			newShanten, _waits := CalculateShantenAndWaits13(tiles34, isOpen)
			if newShanten < shanten {
				// 向听前进了，则换的这张牌为进张
				if _, ok := waits[j]; !ok {
					waits[j] = leftTile - (tiles34[j] - 1)
				}
				if waitsCount := _waits.allCount(); waitsCount > nextShantenWaitsCountMap[j] {
					// 切牌一般切进张最多的
					nextShantenWaitsCountMap[j] = waitsCount
				}
			} else if newShanten == shanten {
				// 向听数没变，但可能是改良型，记录一下
				improves[j] = _waits
			} else {
				// TODO: 向听倒退
			}
			tiles34[j]--
		}
		tiles34[i]++
		rawImprovesMap[i] = improves
	}

	improves := Improves{}
	improveWayCount := 0

	baseWaitsCount, waitTiles := waits.ParseIndex()
	if baseWaitsCount == 0 {
		// TODO ?
	}
	improveWaitsCount34 := make([]int, 34)
	// 初始化成基本进张
	for i := 0; i < 34; i++ {
		improveWaitsCount34[i] = baseWaitsCount
	}
	for discardTile, improve := range rawImprovesMap {
		for drawTile, improveWaits := range improve {
			if inInts(drawTile, waitTiles) {
				// 跳过改良牌就是进张的情况
				continue
			}
			if improveWaitsCount := improveWaits.allCount(); improveWaitsCount > baseWaitsCount {
				// 进张数变多，是改良
				if improveWaitsCount > improveWaitsCount34[drawTile] {
					improveWaitsCount34[drawTile] = improveWaitsCount
					improves[drawTile] = improveWaits
				}
				improveWayCount++
				_ = discardTile
				//fmt.Println(fmt.Sprintf("    摸 %s 切 %s 改良:", mahjongZH[drawTile], mahjongZH[discardTile]), improveWaitsCount, TilesToMergedStrWithBracket(improveWaits.indexes()))
			}
		}
	}

	_tiles34 := make([]int, 34)
	copy(_tiles34, tiles34)
	waitsWithImproves = &WaitsWithImproves13{
		Tiles34: _tiles34,
		Shanten: shanten,
		Waits:   waits,
		NextShantenWaitsCountMap: nextShantenWaitsCountMap,
		Improves:                 improves,
		ImproveWayCount:          improveWayCount,
		ImproveWaitsCount34:      improveWaitsCount34,
	}
	waitsWithImproves.analysis()

	return
}

type WaitsWithImproves14 struct {
	Result13 *WaitsWithImproves13
	// 需要切的牌
	DiscardTile int
	// 切掉这张牌后的向听数
	Shanten int
}

func (r *WaitsWithImproves14) String() string {
	return fmt.Sprintf("切 %s: %s", mahjongZH[r.DiscardTile], r.Result13.String())
}

type WaitsWithImproves14List []*WaitsWithImproves14

func (l WaitsWithImproves14List) Sort() {
	sort.Slice(l, func(i, j int) bool {
		ri, rj := l[i].Result13, l[j].Result13

		if ri.Waits.allCount() != rj.Waits.allCount() {
			return ri.Waits.allCount() > rj.Waits.allCount()
		}

		if ri.AvgNextShantenWaitsCount != rj.AvgNextShantenWaitsCount {
			return ri.AvgNextShantenWaitsCount > rj.AvgNextShantenWaitsCount
		}

		if ri.AvgImproveWaitsCount != rj.AvgImproveWaitsCount {
			return ri.AvgImproveWaitsCount > rj.AvgImproveWaitsCount
		}

		if ri.ImproveWayCount != rj.ImproveWayCount {
			return ri.ImproveWayCount > rj.ImproveWayCount
		}

		return l[i].DiscardTile > l[j].DiscardTile
	})
}

// 14 张牌，计算向听数、进张、改良、向听倒退等
func CalculateShantenWithImproves14(tiles34 []int, isOpen bool) (shanten int, waitsWithImproves WaitsWithImproves14List, incShantenResults WaitsWithImproves14List) {
	shanten = CalculateShanten(tiles34, isOpen)
	//fmt.Println(shanten)

	for i := 0; i < 34; i++ {
		if tiles34[i] == 0 {
			continue
		}
		tiles34[i]-- // 切牌
		result13 := CalculateShantenWithImproves13(tiles34, isOpen)
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
