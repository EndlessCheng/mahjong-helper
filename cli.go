package main

import (
	"fmt"
	"strings"
	"github.com/fatih/color"
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
		improveCount := discard.improveTiles.allCount()
		color.New(getTingCountColor(float64(count))).
			Printf(" 切 %s, 听 %v, %d 张 (%d 张默改, %.2f 改良率)\n", mahjongZH[discard.discardIndex], tiles, count, improveCount, float64(improveCount)/float64(count))
	}
}

//

type ting1Detail struct {
	avgImproveNum   float64
	improveWayCount int
	avgTingCount    float64
}

// TODO: 提醒切这张牌可以断幺！
// TODO: 赤牌改良提醒！！
// TODO: 5万(赤)
func (r *ting1Detail) print() {
	if r.improveWayCount > 0 {
		if r.improveWayCount >= 100 {
			fmt.Printf("%-6.2f[%3d改良]", r.avgImproveNum, r.improveWayCount)
		} else {
			fmt.Printf("%-6.2f[%2d 改良]", r.avgImproveNum, r.improveWayCount)
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

func (l ting1DiscardList) print() {
	for _, discard := range l {
		count, indexes := discard.needs.parseIndex()
		if inIntSlice(discard.discardIndex, indexes) {
			continue
		}

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
var _ting2MinCount = -1

func setTing2MinCount(count int) {
	_ting2MinCount = count
}

func resetTing2MinCount() {
	_ting2MinCount = -1
}

type ting2Discard struct {
	discardIndex int
	needs        needTiles
}

type ting2DiscardList []ting2Discard

func (l ting2DiscardList) print() {
	for _, discard := range l {
		if count, tiles := discard.needs.parse(); count >= _ting2MinCount {
			colorTing2Count(count)
			fmt.Printf("   切 %s %v\n", mahjongZH[discard.discardIndex], tiles)
		}
	}
}
