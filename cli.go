package main

import (
	"fmt"
	"github.com/EndlessCheng/mahjong-helper/util"
	"github.com/fatih/color"
	"math"
	"sort"
	"strings"
)

func printAccountInfo(accountID int) {
	fmt.Printf("您的账号 ID 为 ")
	color.New(color.FgHiGreen).Printf("%d", accountID)
	fmt.Printf("，该数字为雀魂服务器账号数据库中的 ID，该值越小表示您的注册时间越早\n")
}

//

func (p *playerInfo) printDiscards() {
	// TODO: 高亮不合理的舍牌或危险舍牌，如
	// - 一开始就切中张
	// - 开始切中张后，手切了幺九牌（也有可能是有人碰了牌，比如 133m 有人碰了 2m）
	// - 切了 dora，提醒一下
	// - 切了赤宝牌
	// - 有人立直的情况下，多次切出危险度高的牌（有可能是对方读准了牌，或者对方手里的牌与牌河加起来产生了安牌）
	// - 其余可以参考贴吧的《魔神之眼》翻译 https://tieba.baidu.com/p/3311909701
	//      举个简单的例子,如果出现手切了一个对子的情况的话那么基本上就不可能是七对子。
	//      如果对方早巡手切了一个两面搭子的话，那么就可以推理出他在做染手或者牌型是对子型，如果他立直或者鸣牌的话，也比较容易读出他的手牌。
	// https://tieba.baidu.com/p/3311909701
	//      鸣牌之后和终盘的手切牌要尽量记下来，别人手切之前的安牌应该先切掉
	// https://tieba.baidu.com/p/3372239806
	//      吃牌时候打出来的牌的颜色是危险的；碰之后全部的牌都是危险的

	fmt.Printf(p.name + ":")
	for i, disTile := range p.discardTiles {
		fmt.Printf(" ")
		// TODO: 显示 dora, 赤宝牌
		bgColor := color.BgBlack
		fgColor := color.FgWhite
		var tile string
		if disTile >= 0 { // 手切
			tile = util.Mahjong[disTile]
			if disTile >= 27 {
				tile = util.MahjongU[disTile] // 关注字牌的手切
			}
			if p.isNaki { // 副露
				fgColor = getOtherDiscardAlertColor(disTile) // 高亮中张手切
				if util.InInts(i, p.meldDiscardsAt) {
					bgColor = color.BgWhite // 鸣牌时切的那张牌要背景高亮
					fgColor = color.FgBlack
				}
			}
		} else { // 摸切
			disTile = ^disTile
			tile = util.Mahjong[disTile]
			fgColor = color.FgHiBlack // 暗色显示
		}
		color.New(bgColor, fgColor).Print(tile)
	}
	fmt.Println()
}

//

type handsRisk struct {
	tile int
	risk float64
}

// 34 种牌的危险度
type riskTable util.RiskTiles34

func (t riskTable) printWithHands(hands []int, fixedRiskMulti float64) (containLine bool) {
	// 打印铳率=0的牌（现物，或NC且剩余数=0）
	safeCount := 0
	for i, c := range hands {
		if c > 0 && t[i] == 0 {
			fmt.Printf(" " + util.MahjongZH[i])
			safeCount++
		}
	}

	// 打印危险牌，按照铳率排序&高亮
	handsRisks := []handsRisk{}
	for i, c := range hands {
		if c > 0 && t[i] > 0 {
			handsRisks = append(handsRisks, handsRisk{i, t[i]})
		}
	}
	sort.Slice(handsRisks, func(i, j int) bool {
		return handsRisks[i].risk < handsRisks[j].risk
	})
	if len(handsRisks) > 0 {
		if safeCount > 0 {
			fmt.Print(" |")
			containLine = true
		}
		for _, hr := range handsRisks {
			// 颜色考虑了听牌率
			color.New(getNumRiskColor(hr.risk * fixedRiskMulti)).Printf(" " + util.MahjongZH[hr.tile])
		}
	}

	return
}

func (t riskTable) getBestDefenceTile(tiles34 []int) (result int) {
	minRisk := 100.0
	maxRisk := 0.0
	for tile, c := range tiles34 {
		if c == 0 {
			continue
		}
		risk := t[tile]
		if risk < minRisk {
			minRisk = risk
			result = tile
		}
		if risk > maxRisk {
			maxRisk = risk
		}
	}
	if maxRisk == 0 {
		return -1
	}
	return result
}

//

type riskInfo struct {
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

type riskInfoList []*riskInfo

// 考虑了听牌率的综合危险度
func (l riskInfoList) mixedRiskTable() riskTable {
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

func (l riskInfoList) printWithHands(hands []int, leftCounts []int) {
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
	names := []string{"", "下家", "对家", "上家"}
	for i := len(l) - 1; i >= 1; i-- {
		tenpaiRate := l[i].tenpaiRate
		if len(l[i].riskTable) > 0 && (debugMode || tenpaiRate > minShownTenpaiRate) {
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

//

func alertBackwardToShanten2(results util.Hand14AnalysisResultList, incShantenResults util.Hand14AnalysisResultList) {
	if len(results) == 0 || len(incShantenResults) == 0 {
		return
	}

	if results[0].Result13.Waits.AllCount() < 9 {
		if results[0].Result13.MixedWaitsScore < incShantenResults[0].Result13.MixedWaitsScore {
			color.HiGreen("向听倒退？")
		}
	}
}

// 需要提醒的役种
var yakuTypesToAlert = []int{
	//util.YakuKokushi,
	//util.YakuKokushi13,
	util.YakuSuuAnkou,
	util.YakuSuuAnkouTanki,
	util.YakuDaisangen,
	util.YakuShousuushii,
	util.YakuDaisuushii,
	util.YakuTsuuiisou,
	util.YakuChinroutou,
	util.YakuRyuuiisou,
	util.YakuChuuren,
	util.YakuChuuren9,
	util.YakuSuuKantsu,
	//util.YakuTenhou,
	//util.YakuChiihou,

	util.YakuChiitoi,
	util.YakuPinfu,
	util.YakuRyanpeikou,
	util.YakuIipeikou,
	util.YakuSanshokuDoujun,
	util.YakuIttsuu,
	util.YakuToitoi,
	util.YakuSanAnkou,
	util.YakuSanshokuDoukou,
	util.YakuSanKantsu,
	util.YakuTanyao,
	util.YakuChanta,
	util.YakuJunchan,
	util.YakuHonroutou,
	util.YakuShousangen,
	util.YakuHonitsu,
	util.YakuChinitsu,

	util.YakuShiiaruraotai,
	util.YakuUumensai,
	util.YakuSanrenkou,
	util.YakuIsshokusanjun,
}

/*

8     切 3索 听[2万, 7万]
9.20  [20 改良]  4.00 听牌数

4     听 [2万, 7万]
4.50  [ 4 改良]  55.36% 参考和率

8     45万吃，切 4万 听[2万, 7万]
9.20  [20 改良]  4.00 听牌数

*/
// 打印何切分析结果（双行）
func printWaitsWithImproves13_twoRows(result13 *util.Hand13AnalysisResult, discardTile34 int, openTiles34 []int) {
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
		fmt.Print(util.MahjongZH[discardTile34])
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

	fmt.Println()
}

type analysisResult struct {
	discardTile34     int
	isDiscardTileDora bool
	openTiles34       []int
	result13          *util.Hand13AnalysisResult

	mixedRiskTable riskTable

	highlightAvgImproveWaitsCount bool
	highlightMixedScore           bool
}

/*
4[ 4.56] 切 8饼 => 44.50% 参考和率[ 4 改良] [7p 7s] [默听2000] [三色] [振听]

4[ 4.56] 切 8饼 => 0.00% 参考和率[ 4 改良] [7p 7s] [无役]

31[33.58] 切7索 =>  5.23听牌数 [19.21速度] [16改良] [6789p 56789s] [局收支3120] [可能振听]

48[50.64] 切5饼 => 24.25一向听 [12改良] [123456789p 56789s]

31[33.62] 77索碰,切5饼 => 5.48听牌数 [15 改良] [123456789p]

*/
// 打印何切分析结果（单行）
func (r *analysisResult) printWaitsWithImproves13_oneRow() {
	discardTile34 := r.discardTile34
	openTiles34 := r.openTiles34
	result13 := r.result13

	shanten := result13.Shanten

	// 进张数
	waitsCount := result13.Waits.AllCount()
	c := getWaitsCountColor(shanten, float64(waitsCount))
	color.New(c).Printf("%2d", waitsCount)
	// 改良进张均值
	if len(result13.Improves) > 0 {
		if r.highlightAvgImproveWaitsCount {
			color.New(color.FgHiWhite).Printf("[%5.2f]", result13.AvgImproveWaitsCount)
		} else {
			fmt.Printf("[%5.2f]", result13.AvgImproveWaitsCount)
		}
	} else {
		fmt.Print(strings.Repeat(" ", 7))
	}

	fmt.Print(" ")

	// 是否为3k+2张牌的何切分析
	if discardTile34 != -1 {
		// 鸣牌分析
		if len(openTiles34) > 0 {
			meldType := "吃"
			if openTiles34[0] == openTiles34[1] {
				meldType = "碰"
			}
			color.New(color.FgHiWhite).Printf("%s%s", string([]rune(util.MahjongZH[openTiles34[0]])[:1]), util.MahjongZH[openTiles34[1]])
			fmt.Printf("%s,", meldType)
		}
		// 舍牌
		if r.isDiscardTileDora {
			color.New(color.FgHiWhite).Print("ド")
		} else {
			fmt.Print("切")
		}
		tileZH := util.MahjongZH[discardTile34]
		if discardTile34 >= 27 {
			tileZH = " " + tileZH
		}
		if r.mixedRiskTable != nil {
			// 若有实际危险度，则根据实际危险度来显示舍牌危险度
			risk := r.mixedRiskTable[discardTile34]
			if risk == 0 {
				fmt.Print(tileZH)
			} else {
				color.New(getNumRiskColor(risk)).Print(tileZH)
			}
		} else {
			fmt.Print(tileZH)
		}
	}

	fmt.Print(" => ")

	if shanten >= 1 {
		// 前进后的进张数均值
		incShanten := shanten - 1
		c := getWaitsCountColor(incShanten, result13.AvgNextShantenWaitsCount)
		color.New(c).Printf("%5.2f", result13.AvgNextShantenWaitsCount)
		fmt.Printf("%s", util.NumberToChineseShanten(incShanten))
		if incShanten >= 1 {
			//fmt.Printf("进张")
		} else { // incShanten == 0
			fmt.Printf("数")
			//if showAgariAboveShanten1 {
			//	fmt.Printf("（%.2f%% 参考和率）", result13.AvgAgariRate)
			//}
		}
	} else { // shanten == 0
		// 前进后的和率
		// 若振听或片听，则标红
		if result13.FuritenRate == 1 || result13.IsPartWait {
			color.New(color.FgHiRed).Printf("%5.2f%% 参考和率", result13.AvgAgariRate)
		} else {
			fmt.Printf("%5.2f%% 参考和率", result13.AvgAgariRate)
		}
	}

	// 手牌速度，用于快速过庄
	if result13.MixedWaitsScore > 0 && shanten >= 1 && shanten <= 2 {
		fmt.Print(" ")
		if r.highlightMixedScore {
			color.New(color.FgHiWhite).Printf("[%5.2f速度]", result13.MixedWaitsScore)
		} else {
			fmt.Printf("[%5.2f速度]", result13.MixedWaitsScore)
		}
	}

	// 局收支
	if showScore && result13.MixedRoundPoint != 0.0 {
		fmt.Print(" ")
		color.New(color.FgHiGreen).Printf("[局收支%4d]", int(math.Round(result13.MixedRoundPoint)))
	}

	// (默听)荣和点数
	if result13.DamaPoint > 0 {
		fmt.Print(" ")
		ronType := "荣和"
		if !result13.IsNaki {
			ronType = "默听"
		}
		color.New(color.FgHiGreen).Printf("[%s%d]", ronType, int(math.Round(result13.DamaPoint)))
	}

	// 立直点数，考虑了自摸、一发、里宝
	if result13.RiichiPoint > 0 {
		fmt.Print(" ")
		color.New(color.FgHiGreen).Printf("[立直%d]", int(math.Round(result13.RiichiPoint)))
	}

	if len(result13.YakuTypes) > 0 {
		// 役种（两向听以内开启显示）
		if result13.Shanten <= 2 {
			if !showAllYakuTypes && !debugMode {
				shownYakuTypes := []int{}
				for yakuType := range result13.YakuTypes {
					for _, yt := range yakuTypesToAlert {
						if yakuType == yt {
							shownYakuTypes = append(shownYakuTypes, yakuType)
						}
					}
				}
				if len(shownYakuTypes) > 0 {
					sort.Ints(shownYakuTypes)
					fmt.Print(" ")
					color.New(color.FgHiGreen).Printf(util.YakuTypesToStr(shownYakuTypes))
				}
			} else {
				// debug
				fmt.Print(" ")
				color.New(color.FgHiGreen).Printf(util.YakuTypesWithDoraToStr(result13.YakuTypes, result13.DoraCount))
			}
			// 片听
			if result13.IsPartWait {
				fmt.Print(" ")
				color.New(color.FgHiRed).Printf("[片听]")
			}
		}
	} else if result13.IsNaki && shanten >= 0 && shanten <= 2 {
		// 鸣牌时的无役提示（从听牌到两向听）
		fmt.Print(" ")
		color.New(color.FgHiRed).Printf("[无役]")
	}

	// 振听提示
	if result13.FuritenRate > 0 {
		fmt.Print(" ")
		if result13.FuritenRate < 1 {
			color.New(color.FgHiYellow).Printf("[可能振听]")
		} else {
			color.New(color.FgHiRed).Printf("[振听]")
		}
	}

	// 改良数
	if showScore {
		fmt.Print(" ")
		if len(result13.Improves) > 0 {
			fmt.Printf("[%2d改良]", len(result13.Improves))
		} else {
			fmt.Print(strings.Repeat(" ", 4))
			fmt.Print(strings.Repeat("　", 2)) // 全角空格
		}
	}

	// 进张类型
	fmt.Print(" ")
	waitTiles := result13.Waits.AvailableTiles()
	fmt.Print(util.TilesToStrWithBracket(waitTiles))

	//

	fmt.Println()

	if showImproveDetail {
		for tile, waits := range result13.Improves {
			fmt.Printf("摸 %s 改良成 %s\n", util.Mahjong[tile], waits.String())
		}
	}
}

func printResults14WithRisk(results14 util.Hand14AnalysisResultList, mixedRiskTable riskTable) {
	if len(results14) == 0 {
		return
	}

	maxMixedScore := -1.0
	maxAvgImproveWaitsCount := -1.0
	for _, result := range results14 {
		if result.Result13.MixedWaitsScore > maxMixedScore {
			maxMixedScore = result.Result13.MixedWaitsScore
		}
		if result.Result13.AvgImproveWaitsCount > maxAvgImproveWaitsCount {
			maxAvgImproveWaitsCount = result.Result13.AvgImproveWaitsCount
		}
	}

	if len(results14[0].OpenTiles) > 0 {
		fmt.Print("鸣牌后")
	}
	fmt.Println(util.NumberToChineseShanten(results14[0].Result13.Shanten) + "：")

	if results14[0].Result13.Shanten == 0 {
		// 检查听牌是否一样，但是打点不一样
		isDiffPoint := false
		baseWaits := results14[0].Result13.Waits
		baseDamaPoint := results14[0].Result13.DamaPoint
		baseRiichiPoint := results14[0].Result13.RiichiPoint
		for _, result14 := range results14[1:] {
			if baseWaits.Equals(result14.Result13.Waits) && (baseDamaPoint != result14.Result13.DamaPoint || baseRiichiPoint != result14.Result13.RiichiPoint) {
				isDiffPoint = true
				break
			}
		}

		if isDiffPoint {
			color.HiGreen("注意切牌选择：打点")
		}
	}

	// FIXME: 选择很多时如何精简何切选项？
	//const maxShown = 10
	//if len(results14) > maxShown { // 限制输出数量
	//	results14 = results14[:maxShown]
	//}
	for _, result := range results14 {
		r := &analysisResult{
			result.DiscardTile,
			result.IsDiscardDoraTile,
			result.OpenTiles,
			result.Result13,
			mixedRiskTable,
			result.Result13.AvgImproveWaitsCount == maxAvgImproveWaitsCount,
			result.Result13.MixedWaitsScore == maxMixedScore,
		}
		r.printWaitsWithImproves13_oneRow()
	}
}
