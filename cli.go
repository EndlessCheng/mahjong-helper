package main

import (
	"fmt"
	"strings"
	"github.com/fatih/color"
	"sort"
	"github.com/EndlessCheng/mahjong-helper/util"
)

type handsRisk struct {
	tile int
	risk float64
}

// 34 种牌的危险度
type riskTable util.RiskTiles34

func (t riskTable) printWithHands(counts []int) {
	const tab = "   "

	// 打印现物/NC且剩余数=0
	fmt.Printf(tab)
	for i, c := range counts {
		if c > 0 && t[i] == 0 {
			color.New(color.FgHiBlue).Printf(" " + mahjongZH[i])
		}
	}
	fmt.Println()

	// 打印危险牌，按照铳率排序&高亮
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

// 对手的各自危险度
type riskTables []riskTable

func (ts riskTables) printWithHands(counts []int, leftCounts []int) {
	// 打印安牌，危险牌
	printed := false
	names := []string{"", "下家", "对家", "上家"}
	for i := len(ts) - 1; i >= 1; i-- {
		riskTable := ts[i]
		if len(riskTable) > 0 {
			printed = true
			fmt.Println(names[i] + "安牌:")
			riskTable.printWithHands(counts)
		}
	}

	// 打印 NC OC
	if printed {
		ncSafeTileList := util.CalcNCSafeTiles(leftCounts).FilterWithHands(counts)
		ocSafeTileList := util.CalcOCSafeTiles(leftCounts).FilterWithHands(counts)
		if len(ncSafeTileList) > 0 {
			fmt.Printf("NC:")
			for _, safeTile := range ncSafeTileList {
				fmt.Printf(" " + mahjongZH[safeTile.Tile34])
			}
			fmt.Println()
		}
		if len(ocSafeTileList) > 0 {
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
