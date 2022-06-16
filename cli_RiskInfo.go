package main


import (
	"fmt"

	"github.com/fatih/color"
	"github.com/EndlessCheng/mahjong-helper/util"
)

type RiskInfo struct {
	// 三麻为 3，四麻为 4
	playerNumber int

	// 该玩家的听牌率（立直时为 100.0）
	tenpaiRate float64

	// 该玩家的安牌
	// 若该玩家有杠操作，把杠的那张牌也算作安牌，这有助于判断筋壁危险度
	safeTiles34 []bool

	// 各种牌的铳率表
	riskTable riskTable

	// 剩余无筋 123789
	// 总计 18 种。剩余无筋牌数量越少，该无筋牌越危险
	leftNoSujiTiles []int

	// 是否摸切立直
	isTsumogiriRiichi bool

	// 荣和点数
	// 仅调试用
	_ronPoint float64
}

type RiskInfoList []*RiskInfo

// 考虑了听牌率的综合危险度
func (l RiskInfoList) mixedRiskTable() riskTable {
	mixedRiskTable := make(riskTable, 34)
	for i := range mixedRiskTable {
		mixedRisk := 0.0
		for _, ri := range l[1:] {
			if ri.tenpaiRate <= 15 {
				continue
			}
			_risk := ri.riskTable[i] * ri.tenpaiRate / 100
			mixedRisk = mixedRisk + _risk - mixedRisk*_risk/100
		}
		mixedRiskTable[i] = mixedRisk
	}
	return mixedRiskTable
}

func (l RiskInfoList) printWithHands(hands []int, leftCounts []int) {
	// 听牌率超过一定值就打印铳率
	const (
		minShownTenpaiRate4 = 50.0
		minShownTenpaiRate3 = 20.0
	)

	minShownTenpaiRate := minShownTenpaiRate4
	if l[0].playerNumber == 3 {
		minShownTenpaiRate = minShownTenpaiRate3
	}

	dangerousPlayerCount := 0
	// 打印安牌，危险牌
	names := []string{"", "下家", "對家", "上家"}
	for i := len(l) - 1; i >= 1; i-- {
		tenpaiRate := l[i].tenpaiRate
		if len(l[i].riskTable) > 0 && (DebugMode || tenpaiRate > minShownTenpaiRate) {
			dangerousPlayerCount++
			fmt.Print(names[i] + "安牌:")
			//if debugMode {
			//fmt.Printf("(%d*%2.2f%%听牌率)", int(l[i]._ronPoint), l[i].tenpaiRate)
			//}
			containLine := l[i].riskTable.printWithHands(hands, tenpaiRate/100)

			// 打印听牌率
			fmt.Print(" ")
			if !containLine {
				fmt.Print("  ")
			}
			fmt.Print("[")
			if tenpaiRate == 100 {
				fmt.Print("100.%")
			} else {
				fmt.Printf("%4.1f%%", tenpaiRate)
			}
			fmt.Print("听牌率]")

			// 打印无筋数量
			fmt.Print(" ")
			const badMachiLimit = 3
			noSujiInfo := ""
			if l[i].isTsumogiriRiichi {
				noSujiInfo = "摸切立直"
			} else if len(l[i].leftNoSujiTiles) == 0 {
				noSujiInfo = "愚形听牌/振听"
			} else if len(l[i].leftNoSujiTiles) <= badMachiLimit {
				noSujiInfo = "可能愚形听牌/振听"
			}
			if noSujiInfo != "" {
				fmt.Printf("[%d无筋: ", len(l[i].leftNoSujiTiles))
				color.New(color.FgHiYellow).Printf("%s", noSujiInfo)
				fmt.Print("]")
			} else {
				fmt.Printf("[%d无筋]", len(l[i].leftNoSujiTiles))
			}

			fmt.Println()
		}
	}

	// 若不止一个玩家立直/副露，打印加权综合铳率（考虑了听牌率）
	mixedPlayers := 0
	for _, ri := range l[1:] {
		if ri.tenpaiRate > 0 {
			mixedPlayers++
		}
	}
	if dangerousPlayerCount > 0 && mixedPlayers > 1 {
		fmt.Print("综合安牌:")
		mixedRiskTable := l.mixedRiskTable()
		mixedRiskTable.printWithHands(hands, 1)
		fmt.Println()
	}

	// 打印因 NC OC 产生的安牌
	// TODO: 重构至其他函数
	if dangerousPlayerCount > 0 {
		ncSafeTileList := util.CalcNCSafeTiles(leftCounts).FilterWithHands(hands)
		ocSafeTileList := util.CalcOCSafeTiles(leftCounts).FilterWithHands(hands)
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
