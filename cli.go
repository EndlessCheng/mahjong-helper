package main

import (
	"fmt"
	"strings"
	"github.com/fatih/color"
	"sort"
)

type ting0Improve struct {
	drawIndex    int
	discardIndex int
	needs        needTiles
}

type ting0ImproveList []ting0Improve

// 13 张牌，计算默听改良
func (l ting0ImproveList) calcGoodImprove(counts []int) needTiles {
	goodTiles := needTiles{}
	for _, improve := range l {
		if _, ok := goodTiles[improve.drawIndex]; ok {
			continue
		}

		if improve.needs.containHonors() {
			goodTiles[improve.drawIndex] = 4 - counts[improve.drawIndex]
		} else {
			if count := improve.needs.allCount(); count > 4 {
				goodTiles[improve.drawIndex] = 4 - counts[improve.drawIndex]
			}
		}
	}
	return goodTiles
}

func (l ting0ImproveList) print() {
	if len(l) == 0 {
		fmt.Println("没有合适的改良")
		return
	}

	for _, improve := range l {
		count, tiles := improve.needs.parseZH()
		text := fmt.Sprintf("摸 %s 切 %s，听 %v, %d 张", mahjongZH[improve.drawIndex], mahjongZH[improve.discardIndex], tiles, count)
		var ting0Color color.Attribute
		if improve.needs.containHonors() {
			// 听字牌算良型听牌
			ting0Color = color.FgHiRed
		} else {
			ting0Color = getTingCountColor(float64(count))
		}
		color.New(ting0Color).Println(text)
	}
}

//

type ting0Discard struct {
	discardIndex int
	needs        needTiles
	improveTiles needTiles
}

type ting0DiscardList []ting0Discard

func (l ting0DiscardList) print() {
	for _, discard := range l {
		count, tiles := discard.needs.parseZH()
		printer := color.New(getTingCountColor(float64(count)))

		improveCount := discard.improveTiles.allCount()
		printer.Printf(" 切 %s, 听 %v, %d 张 (%d 张默改, %.2f 改良比)", mahjongZH[discard.discardIndex], tiles, count, improveCount, float64(improveCount)/float64(count))

		if agariRate := calcAgariRate(discard.needs); agariRate > 0 {
			printer.Printf(" (%.2f%% 和了率)", agariRate)
		}

		fmt.Println()
	}
}

//

type ting1Detail struct {
	needs needTiles

	avgImproveTing1Count float64
	improveWayCount      int
	avgTingCount         float64
}

// TODO: 提醒切这张牌可以断幺
// TODO: 赤牌改良提醒
// TODO: 5万(赤)
func (r *ting1Detail) print() {
	if r.improveWayCount > 0 {
		if r.improveWayCount >= 100 { // 三位数
			fmt.Printf("%-6.2f[%3d改良]", r.avgImproveTing1Count, r.improveWayCount)
		} else {
			fmt.Printf("%-6.2f[%2d 改良]", r.avgImproveTing1Count, r.improveWayCount)
		}
	} else {
		fmt.Print(strings.Repeat(" ", 15))
	}

	fmt.Print(" ")
	color.New(getTingCountColor(r.avgTingCount)).Printf("%5.2f", r.avgTingCount)
	fmt.Print(" 听牌数")

	fmt.Println()
}

//

type ting1Discard struct {
	discardIndex int
	needs        needTiles
	ting1Detail  *ting1Detail
}

type ting1DiscardList []ting1Discard

func (l ting1DiscardList) maxAvgTing1Count() float64 {
	maxAvg := 0.0
	for _, discard := range l {
		if discard.ting1Detail.improveWayCount > 0 {
			maxAvg = maxFloat64(maxAvg, discard.ting1Detail.avgImproveTing1Count)
		} else {
			maxAvg = maxFloat64(maxAvg, float64(discard.needs.allCount()))
		}
	}
	return maxAvg
}

// 是否为完全一向听
func (l ting1DiscardList) isGood() bool {
	return l.maxAvgTing1Count() >= 16
}

func (l ting1DiscardList) print() {
	for _, discard := range l {
		count, indexes := discard.needs.parseIndex()

		// 过滤掉不算改良的向听倒退
		if inIntSlice(discard.discardIndex, indexes) && (count <= 24 || count > 100) {
			continue
		}

		// 8     切 3索 [2万, 7万]
		// 9.20  [20 改良]  4.00 听牌数

		fmt.Println()
		colorTing1Count(count)
		fmt.Print("切 ")
		color.New(getRiskColor(discard.discardIndex)).Print(mahjongZH[discard.discardIndex])
		fmt.Print(" [")
		color.New(getSafeColor(indexes[0])).Print(mahjongZH[indexes[0]])
		for _, index := range indexes[1:] {
			fmt.Print(", ")
			color.New(getSafeColor(index)).Print(mahjongZH[index])
		}
		fmt.Print("]")
		fmt.Println()
		discard.ting1Detail.print()
		//flushBuffer()
	}
}

//

// 交互模式下，两向听进张的最低值
//var _ting2MinCount = -1
//
//func setTing2MinCount(count int) {
//	_ting2MinCount = count
//}
//
//func resetTing2MinCount() {
//	_ting2MinCount = -1
//}

type ting2Discard struct {
	discardIndex int
	needs        needTiles
}

type ting2DiscardList []ting2Discard

func (l ting2DiscardList) maxTing2Count() int {
	maxTing2Count := 0
	for _, discard := range l {
		maxTing2Count = maxInt(maxTing2Count, discard.needs.allCount())
	}
	return maxTing2Count
}

// 按照 needs.allCount() 从大到小排序
func (l ting2DiscardList) sort() {
	sort.Slice(l, func(i, j int) bool {
		return l[i].needs.allCount() > l[j].needs.allCount()
	})
}

func (l ting2DiscardList) print() {
	const printLimitExceptMax = 5
	printCount := 0

	l.sort()
	maxTing2Count := l[0].needs.allCount()
	for _, discard := range l {
		if count, tiles := discard.needs.parse(); count == maxTing2Count || printCount < printLimitExceptMax {
			printCount++
			colorTing2Count(count)
			fmt.Printf("   切 %s %v\n", mahjongZH[discard.discardIndex], tiles)
		}
	}
}
