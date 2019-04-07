package main

import (
	"fmt"
	"strings"
	"github.com/fatih/color"
	"sort"
	"github.com/EndlessCheng/mahjong-helper/util"
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
	agariRate    float64
}

type ting0DiscardList []ting0Discard

func (l ting0DiscardList) sort() {
	sort.Slice(l, func(i, j int) bool {
		return l[i].agariRate > l[j].agariRate
	})
}

func (l ting0DiscardList) print() {
	l.sort()
	for _, discard := range l {
		count, tiles := discard.needs.parseZH()
		//printer := color.New(getTingCountColor(float64(count)))

		improveCount := discard.improveTiles.allCount()
		fmt.Printf(" 切 %s, 听 %v, %d 张 (%d 张默改, %.2f 改良比)", mahjongZH[discard.discardIndex], tiles, count, improveCount, float64(improveCount)/float64(count+1))
		fmt.Printf(" (%.2f%% 和了率)\n", discard.agariRate)
	}
}

func (l ting0DiscardList) printWithLeftCounts(leftCounts []int) {
	if leftCounts != nil {
		for _, discard := range l {
			discard.needs.fixCountsWithLeftCounts(leftCounts)
			discard.improveTiles.fixCountsWithLeftCounts(leftCounts)
		}
	}
	l.print()
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

// http://blog.sina.com.cn/s/blog_721350d40101dto4.html
// http://blog.sina.com.cn/s/blog_721350d40102uwbe.html
// 按照
// 1. needs.allCount() 降序
// 2. ting1Detail.avgTingCount 降序
// 3. ting1Detail.avgImproveTing1Count 降序
// 4. TODO: needs.indexes 的和牌容易度
// 5. TODO: 切牌的危险度（好牌先打，或者先打安全点的牌）
func (l ting1DiscardList) sort() {
	sort.Slice(l, func(i, j int) bool {
		li, lj := l[i], l[j]

		if li.needs.allCount() != lj.needs.allCount() {
			return li.needs.allCount() > lj.needs.allCount()
		}

		if li.ting1Detail.avgTingCount != lj.ting1Detail.avgTingCount {
			return li.ting1Detail.avgTingCount > lj.ting1Detail.avgTingCount
		}

		// avgTingCount 为零视作没有改良，排同等级最末
		return li.ting1Detail.avgImproveTing1Count > lj.ting1Detail.avgImproveTing1Count
	})
}

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
	l.sort()

	for _, discard := range l {
		count, indexes := discard.needs.parseIndex()

		// 过滤掉不算改良的向听倒退
		if inIntSlice(discard.discardIndex, indexes) && (count <= 24 || count > 100) {
			continue
		}

		// 8     切 3索 [2万, 7万]
		// 9.20  [20 改良]  4.00 听牌数

		colorTing1Count(count)
		fmt.Print("切 ")
		color.New(getSimpleRiskColor(discard.discardIndex)).Print(mahjongZH[discard.discardIndex])
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

func (l ting1DiscardList) printWithLeftCounts(leftCounts []int) {
	if leftCounts != nil {
		for _, discard := range l {
			discard.needs.fixCountsWithLeftCounts(leftCounts)
			// TODO: discard.ting1Detail
			// TODO: 也就是说处理数据的过程移到此处
		}
	}
	l.print()
}

//

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
	if len(l) == 0 {
		fmt.Println("无有效两向听")
		return
	}

	const printLimitExceptMax = 5
	printCount := 0

	l.sort()
	maxTing2Count := l[0].needs.allCount()
	for _, discard := range l {
		if count := discard.needs.allCount(); count == maxTing2Count || printCount < printLimitExceptMax {
			printCount++
			colorTing2Count(count)
			fmt.Printf("   切 %s %v\n", mahjongZH[discard.discardIndex], util.TilesToMergedStrWithBracket(discard.needs.indexes()))
		}
	}
}

func (l ting2DiscardList) printWithLeftCounts(leftCounts []int) {
	if leftCounts != nil {
		for _, discard := range l {
			discard.needs.fixCountsWithLeftCounts(leftCounts)
		}
	}
	l.print()
}

//

type handsRisk struct {
	tile int
	risk float64
}

type riskTable []float64

func (t riskTable) printWithHands(counts []int) {
	tab := "   "
	fmt.Printf(tab)
	for i, c := range counts {
		if c > 0 && t[i] == 0 {
			color.New(color.FgHiBlue).Printf(" " + mahjongZH[i])
		}
	}
	fmt.Println()

	handsRisks := []handsRisk{}
	for i, c := range counts {
		if c > 0 && t[i] > 0 {
			handsRisks = append(handsRisks, handsRisk{i, t[i]})
		}
	}
	sort.Slice(handsRisks, func(i, j int) bool {
		return handsRisks[i].risk < handsRisks[j].risk
	})
	fmt.Printf(tab)
	for _, hr := range handsRisks {
		color.New(getNumRiskColor(hr.risk)).Printf(" " + mahjongZH[hr.tile])
	}
	fmt.Println()
}

type riskTables []riskTable

func (ts riskTables) printWithHands(counts []int, leftCounts []int) {
	printed := false
	names := []string{"下家", "对家", "上家"}
	for i, table := range ts {
		if len(table) > 0 {
			printed = true
			fmt.Println(names[i] + "安牌:")
			table.printWithHands(counts)
		}
	}

	// NC OC
	if printed {
		ncSafeTileList := util.CalcNCSafeTiles34(leftCounts).FilterWithHands(counts)
		if len(ncSafeTileList) != 0 {
			fmt.Printf("NC:")
			for _, safeTile := range ncSafeTileList {
				fmt.Printf(" " + mahjongZH[safeTile.Tile34])
			}
			fmt.Println()
		}
		ocSafeTileList := util.CalcOCSafeTiles34(leftCounts).FilterWithHands(counts)
		if len(ocSafeTileList) != 0 {
			fmt.Printf("OC:")
			for _, safeTile := range ocSafeTileList {
				fmt.Printf(" " + mahjongZH[safeTile.Tile34])
			}
			fmt.Println()
		}
		//printedNC := false
		//for i, c := range leftCounts[:27] {
		//	if c != 0 || i%9 == 0 || i%9 == 8 {
		//		continue
		//	}
		//	if !printedNC {
		//		printedNC = true
		//		fmt.Printf("NC:")
		//	}
		//	fmt.Printf(" " + mahjongZH[i])
		//}
		//if printedNC {
		//	fmt.Println()
		//}
		//printedOC := false
		//for i, c := range leftCounts[:27] {
		//	if c != 1 || i%9 == 0 || i%9 == 8 {
		//		continue
		//	}
		//	if !printedOC {
		//		printedOC = true
		//		fmt.Printf("OC:")
		//	}
		//	fmt.Printf(" " + mahjongZH[i])
		//}
		//if printedOC {
		//	fmt.Println()
		//}
		fmt.Println()
	}
}

//

func printAccountInfo(accountID int) {
	fmt.Printf("您的账号 ID 为 ")
	color.New(color.FgMagenta).Printf("%d", accountID)
	fmt.Printf("，该数字为雀魂服务器账号数据库中的 ID，该值越小表示您的注册时间越早\n")
}

/*

8     切 3索 听[2万, 7万]
9.20  [20 改良]  4.00 听牌数

4     听 [2万, 7万]
4.50  [ 4 改良]  55.36% 和了率


8     45万吃，切 4万 听[2万, 7万]
9.20  [20 改良]  4.00 听牌数

*/
// TODO: 按照和率排序？
func printWaitsWithImproves13(result13 *util.WaitsWithImproves13, discardTile34 int, openTiles34 []int) {
	shanten := result13.Shanten
	waits := result13.Waits

	waitsCount, waitTiles := waits.ParseIndex()
	colors := getShantenWaitsCountColors(shanten, waitsCount)
	color.New(colors...).Printf("%-6d", waitsCount)
	if discardTile34 != -1 {
		if len(openTiles34) > 0 {
			meldType := "吃"
			if openTiles34[0] == openTiles34[1] {
				meldType = "碰"
			}
			color.New(color.FgHiWhite).Printf("%s%s", string([]rune(mahjongZH[openTiles34[0]])[:1]), mahjongZH[openTiles34[1]])
			fmt.Printf("%s，", meldType)
		}
		fmt.Print("切 ")
		if shanten <= 1 {
			color.New(getSimpleRiskColor(discardTile34)).Print(mahjongZH[discardTile34])
		} else {
			fmt.Print(mahjongZH[discardTile34])
		}
		fmt.Print(" ")
	}
	//fmt.Print("等")
	if shanten <= 1 {
		fmt.Print("[")
		color.New(getSafeColor(waitTiles[0])).Print(mahjongZH[waitTiles[0]])
		for _, index := range waitTiles[1:] {
			fmt.Print(", ")
			color.New(getSafeColor(index)).Print(mahjongZH[index])
		}
		fmt.Print("]")
		fmt.Println()
	} else {
		fmt.Println(util.TilesToMergedStrWithBracket(waitTiles))
	}

	if result13.ImproveWayCount > 0 {
		if result13.ImproveWayCount >= 100 { // 三位数
			fmt.Printf("%-6.2f[%3d改良]", result13.AvgImproveWaitsCount, result13.ImproveWayCount)
		} else {
			fmt.Printf("%-6.2f[%2d 改良]", result13.AvgImproveWaitsCount, result13.ImproveWayCount)
		}
	} else {
		fmt.Print(strings.Repeat(" ", 15))
	}

	fmt.Print(" ")

	if shanten >= 1 {
		_color := getNextShantenWaitsCountColor(shanten, result13.AvgNextShantenWaitsCount)
		color.New(_color).Printf("%5.2f", result13.AvgNextShantenWaitsCount)
		fmt.Printf(" %s", util.NumberToChineseShanten(shanten-1))
		if shanten >= 2 {
			fmt.Printf("进张")
		} else {
			fmt.Printf("数")
			//fmt.Printf("（%.2f%% 参考和率）", result13.AvgAgariRate)
		}
		mixedScore := float64(waitsCount) * result13.AvgNextShantenWaitsCount
		for i := 2; i <= shanten; i++ {
			mixedScore /= 4
		}
		fmt.Printf("（%.2f 综合分）", mixedScore)
	} else { // shanten == 0
		fmt.Printf("%5.2f%% 参考和率", result13.AvgAgariRate)
	}

	//if dangerous {
	//	// TODO: 提示危险度！
	//}

	fmt.Println()
}
