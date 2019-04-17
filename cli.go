package main

import (
	"fmt"
	"strings"
	"github.com/fatih/color"
	"sort"
	"github.com/EndlessCheng/mahjong-helper/util"
)

func printAccountInfo(accountID int) {
	fmt.Printf("您的账号 ID 为 ")
	color.New(color.FgHiGreen).Printf("%d", accountID)
	fmt.Printf("，该数字为雀魂服务器账号数据库中的 ID，该值越小表示您的注册时间越早\n")
}

//

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
			color.New(color.FgHiBlue).Printf(" " + util.MahjongZH[i])
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
		color.New(getNumRiskColor(hr.risk)).Printf(" " + util.MahjongZH[hr.tile])
	}
}

type riskInfo struct {
	// 各种牌的铳率表
	riskTable riskTable

	// 剩余无筋 123789
	// 总计 18 种。剩余无筋牌数量越少，该无筋牌越危险
	leftNoSujiTiles []int
}

type riskInfoList []riskInfo

func (ri riskInfoList) printWithHands(counts []int, leftCounts []int) {
	// 打印安牌，危险牌
	printed := false
	names := []string{"", "下家", "对家", "上家"}
	for i := len(ri) - 1; i >= 1; i-- {
		riskTable := ri[i].riskTable
		if len(riskTable) > 0 {
			printed = true
			fmt.Println(names[i] + "安牌:")
			riskTable.printWithHands(counts)

			// 打印无筋数量和种类
			noSujiInfo := util.TilesToStr(ri[i].leftNoSujiTiles)
			if len(ri[i].leftNoSujiTiles) == 0 {
				noSujiInfo = "非好型听牌/振听"
			}
			fmt.Printf(" [%d无筋: %s]", len(ri[i].leftNoSujiTiles), noSujiInfo)

			fmt.Println()
		}
	}

	// 打印因 NC OC 产生的安牌
	// TODO: 重构至其他函数
	if printed {
		ncSafeTileList := util.CalcNCSafeTiles(leftCounts).FilterWithHands(counts)
		ocSafeTileList := util.CalcOCSafeTiles(leftCounts).FilterWithHands(counts)
		if len(ncSafeTileList) > 0 {
			fmt.Printf("NC:")
			for _, safeTile := range ncSafeTileList {
				fmt.Printf(" " + util.MahjongZH[safeTile.Tile34])
			}
			fmt.Println()
		}
		if len(ocSafeTileList) > 0 {
			fmt.Printf("OC:")
			for _, safeTile := range ocSafeTileList {
				fmt.Printf(" " + util.MahjongZH[safeTile.Tile34])
			}
			fmt.Println()
		}

		// 下面这个是另一种显示方式：显示壁牌
		//printedNC := false
		//for i, c := range leftCounts[:27] {
		//	if c != 0 || i%9 == 0 || i%9 == 8 {
		//		continue
		//	}
		//	if !printedNC {
		//		printedNC = true
		//		fmt.Printf("NC:")
		//	}
		//	fmt.Printf(" " + util.MahjongZH[i])
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
		//	fmt.Printf(" " + util.MahjongZH[i])
		//}
		//if printedOC {
		//	fmt.Println()
		//}
		fmt.Println()
	}
}

//

/*

8     切 3索 听[2万, 7万]
9.20  [20 改良]  4.00 听牌数

4     听 [2万, 7万]
4.50  [ 4 改良]  55.36% 参考和率

8     45万吃，切 4万 听[2万, 7万]
9.20  [20 改良]  4.00 听牌数

*/
// 打印何切分析结果（双行）
func printWaitsWithImproves13_twoRows(result13 *util.WaitsWithImproves13, discardTile34 int, openTiles34 []int) {
	shanten := result13.Shanten
	waits := result13.Waits

	waitsCount, waitTiles := waits.ParseIndex()
	c := getWaitsCountColor(shanten, float64(waitsCount))
	color.New(c).Printf("%-6d", waitsCount)
	if discardTile34 != -1 {
		if len(openTiles34) > 0 {
			meldType := "吃"
			if openTiles34[0] == openTiles34[1] {
				meldType = "碰"
			}
			color.New(color.FgHiWhite).Printf("%s%s", string([]rune(util.MahjongZH[openTiles34[0]])[:1]), util.MahjongZH[openTiles34[1]])
			fmt.Printf("%s，", meldType)
		}
		fmt.Print("切 ")
		if shanten <= 1 {
			color.New(getSelfDiscardRiskColor(discardTile34)).Print(util.MahjongZH[discardTile34])
		} else {
			fmt.Print(util.MahjongZH[discardTile34])
		}
		fmt.Print(" ")
	}
	//fmt.Print("等")
	//if shanten <= 1 {
	//	fmt.Print("[")
	//	if len(waitTiles) > 0 {
	//		fmt.Print(util.MahjongZH[waitTiles[0]])
	//		for _, idx := range waitTiles[1:] {
	//			fmt.Print(", " + util.MahjongZH[idx])
	//		}
	//	}
	//	fmt.Println("]")
	//} else {
	fmt.Println(util.TilesToStrWithBracket(waitTiles))
	//}

	if len(result13.Improves) > 0 {
		fmt.Printf("%-6.2f[%2d 改良]", result13.AvgImproveWaitsCount, len(result13.Improves))
	} else {
		fmt.Print(strings.Repeat(" ", 15))
	}

	fmt.Print(" ")

	if shanten >= 1 {
		c := getWaitsCountColor(shanten-1, result13.AvgNextShantenWaitsCount)
		color.New(c).Printf("%5.2f", result13.AvgNextShantenWaitsCount)
		fmt.Printf(" %s", util.NumberToChineseShanten(shanten-1))
		if shanten >= 2 {
			fmt.Printf("进张")
		} else { // shanten == 1
			fmt.Printf("数")
			if showAgariAboveShanten1 {
				fmt.Printf("（%.2f%% 参考和率）", result13.AvgAgariRate)
			}
		}
		if showScore {
			mixedScore := result13.MixedWaitsScore
			//for i := 2; i <= shanten; i++ {
			//	mixedScore /= 4
			//}
			fmt.Printf("（%.2f 综合分）", mixedScore)
		}
	} else { // shanten == 0
		fmt.Printf("%5.2f%% 参考和率", result13.AvgAgariRate)
	}

	//if dangerous {
	//	// TODO: 提示危险度！
	//}

	fmt.Println()
}

/*
4[ 4.56] 切 8饼 => 44.50% 参考和率[ 4 改良] [7p 7s]

31[33.58] 切7索 => 5.23听牌数 [16改良] [6789p 56789s]

48[50.64] 切5饼 => 24.25一向听进张 [12改良] [123456789p 56789s]

31[33.62] 77索碰,切5饼 => 5.48听牌数 [15 改良] [123456789p]

*/
// 打印何切分析结果（单行）
func printWaitsWithImproves13_oneRow(result13 *util.WaitsWithImproves13, discardTile34 int, openTiles34 []int) {
	shanten := result13.Shanten

	// 打印进张数
	waitsCount, waitTiles := result13.Waits.ParseIndex()
	c := getWaitsCountColor(shanten, float64(waitsCount))
	color.New(c).Printf("%2d", waitsCount)
	// 打印改良进张均值
	if len(result13.Improves) > 0 {
		fmt.Printf("[%5.2f]", result13.AvgImproveWaitsCount)
	} else {
		fmt.Print(strings.Repeat(" ", 7))
	}

	fmt.Print(" ")

	// 是否为3k+2张牌的何切分析
	if discardTile34 != -1 {
		// 打印鸣牌分析
		if len(openTiles34) > 0 {
			meldType := "吃"
			if openTiles34[0] == openTiles34[1] {
				meldType = "碰"
			}
			color.New(color.FgHiWhite).Printf("%s%s", string([]rune(util.MahjongZH[openTiles34[0]])[:1]), util.MahjongZH[openTiles34[1]])
			fmt.Printf("%s,", meldType)
		}
		// 打印舍牌
		fmt.Print("切")
		tileZH := util.MahjongZH[discardTile34]
		if discardTile34 >= 27 {
			tileZH = " " + tileZH
		}
		if shanten <= 1 {
			color.New(getSelfDiscardRiskColor(discardTile34)).Print(tileZH)
		} else {
			fmt.Print(tileZH)
		}
	}

	fmt.Print(" => ")

	if shanten >= 1 {
		// 打印前进后的进张数均值
		c := getWaitsCountColor(shanten-1, result13.AvgNextShantenWaitsCount)
		color.New(c).Printf("%5.2f", result13.AvgNextShantenWaitsCount)
		fmt.Printf("%s", util.NumberToChineseShanten(shanten-1))
		if shanten >= 2 {
			//fmt.Printf("进张")
		} else { // shanten == 1
			fmt.Printf("数")
			if showAgariAboveShanten1 {
				fmt.Printf("（%.2f%% 参考和率）", result13.AvgAgariRate)
			}
		}
		if showScore {
			mixedScore := result13.MixedWaitsScore
			//for i := 2; i <= shanten; i++ {
			//	mixedScore /= 4
			//}
			fmt.Printf("（%.2f 综合分）", mixedScore)
		}
	} else { // shanten == 0
		fmt.Printf("%5.2f%% 参考和率", result13.AvgAgariRate)
	}

	//if dangerous {
	//	// TODO: 提示危险度！
	//}

	fmt.Print(" ")

	// 打印改良数
	if len(result13.Improves) > 0 {
		fmt.Printf("[%2d改良]", len(result13.Improves))
	} else {
		fmt.Print(strings.Repeat(" ", 4))
		fmt.Print(strings.Repeat("　", 2)) // 全角空格
	}

	fmt.Print(" ")

	// 打印进张类型
	fmt.Print(util.TilesToStrWithBracket(waitTiles))

	fmt.Println()
}
