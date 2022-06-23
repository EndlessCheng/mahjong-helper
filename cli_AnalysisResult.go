package main

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/EndlessCheng/mahjong-helper/util"
	"github.com/fatih/color"
)

// define struct
type AnalysisResult struct {
	discardTile34     int
	isDiscardTileDora bool
	openTiles34       []int
	result13          *util.Hand13AnalysisResult

	MixedRiskTable riskTable

	highlightAvgImproveWaitsCount bool
	highlightMixedScore           bool
}

var (
	discardTile34 int
	openTiles34   []int
	result13      *util.Hand13AnalysisResult
	shanten       int
)

/*
4[ 4.56] 切 8饼 => 44.50% 参考和率[ 4 改良] [7p 7s] [默听2000] [三色] [振听]

4[ 4.56] 切 8饼 => 0.00% 参考和率[ 4 改良] [7p 7s] [无役]

31[33.58] 切7索 =>  5.23听牌数 [19.21速度] [16改良] [6789p 56789s] [局收支3120] [可能振听]

48[50.64] 切5饼 => 24.25一向听 [12改良] [123456789p 56789s]

31[33.62] 77索碰,切5饼 => 5.48听牌数 [15 改良] [123456789p]

*/
// 打印何切分析结果（单行）
func (r *AnalysisResult) PrintWaitsWithImproves13_oneRow() {
	discardTile34 = r.discardTile34
	openTiles34 = r.openTiles34
	result13 = r.result13

	shanten = result13.Shanten

	// 进张数
	countDrawingAUsefulTile()

	// 改良进张均值
	averageOfChangingDAUT(r)

	// 是否为3k+2张牌的何切分析
	wODOf_M3P2_Analysis(r)

	// 向听前进后的进张数的加权均值
	waOfDAUTAfterShanTenAdvance()

	// 手牌速度，用于快速过庄
	dAUTEfficiency(r)

	// 局收支；(默听)荣和点数；立直点数，考虑了自摸、一发、里宝
	expectedScore()

	// 牌型
	yaKuType()

	// 振听提示
	fuRiTenHint()

	// 改良数
	countChangeDAUT()

	// 进张类型
	typeOfDAUT()

	// 是否可以改變牌型
	theHandCanImproveOrNot()
}

// count drawing a useful tile funcction
func countDrawingAUsefulTile() {
	// 进张数
	waitsCount := result13.Waits.AllCount()
	c := getWaitsCountColor(shanten, float64(waitsCount))
	color.New(c).Printf("%2d ", waitsCount)
}

// average of changing drawing a useful tile(DAUT)
func averageOfChangingDAUT(r *AnalysisResult) {
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
}

// is "Whitch one to Discard(WOD)" of (meld * 3 + 2) type(M3P2) analysis
func wODOf_M3P2_Analysis(r *AnalysisResult) {
	if discardTile34 != -1 {
		// 鳴牌分析
		if len(openTiles34) > 0 {
			meldType := "吃"
			if openTiles34[0] == openTiles34[1] {
				meldType = "碰"
			}
			color.New(color.FgHiWhite).Printf("%s%s", string([]rune(util.MahjongZH[openTiles34[0]])[:1]),
				util.MahjongZH[openTiles34[1]])
			fmt.Printf("%s,", meldType)
		}
		// 舍牌
		tileZH := util.MahjongZH[discardTile34] + " "
		// if it's dora print red text
		if r.isDiscardTileDora {
			color.New(color.FgHiRed).Printf("切")
		} else {
			fmt.Print("切")
		}
		// if tiel kind is z append space before tileZH
		if discardTile34 >= 27 {
			tileZH = " " + tileZH
		}
		// if the discard tile is not safe
		if r.MixedRiskTable != nil {
			// 若有实际危险度，则根据实际危险度来显示舍牌危险度
			risk := r.MixedRiskTable[discardTile34]
			if risk == 0 {
				fmt.Print(tileZH)
			} else {
				color.New(GetNumRiskColor(risk)).Print(tileZH)
			}
		} else {
			fmt.Print(tileZH)
		}
	}

	fmt.Print(" => ")
}

// weighted-average(wa) of drawing a useful tile after shanten advance
func waOfDAUTAfterShanTenAdvance() {
	if shanten >= 1 {
		// 前进后的进张数均值
		incShanten := shanten - 1
		c := getWaitsCountColor(incShanten, result13.AvgNextShantenWaitsCount)
		color.New(c).Printf("%5.2f ", result13.AvgNextShantenWaitsCount)
		fmt.Printf("%s ", util.NumberToChineseShanten(incShanten))
		if incShanten >= 1 {
			//fmt.Printf("进张")
		} else { // incShanten == 0
			fmt.Printf("數")
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

}

// drawing a useful tile(DAUT) efficiency
func dAUTEfficiency(r *AnalysisResult) {
	if result13.MixedWaitsScore > 0 && shanten >= 1 /*&& shanten <= 2*/ {
		// fmt.Print(" ")
		if r.highlightMixedScore {
			color.New(color.FgHiWhite).Printf("[%5.2f速度]", result13.MixedWaitsScore)
		} else {
			fmt.Printf("[%5.2f速度]", result13.MixedWaitsScore)
		}
	}

}

// expected score
func expectedScore() {
	// 局收支
	if ShowScore && result13.MixedRoundPoint != 0.0 {
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
}

// yaKuType
func yaKuType() {
	if len(result13.YakuTypes) > 0 {
		// 役种（两向听以内开启显示）
		if result13.Shanten <= 2 {
			if !ShowAllYakuTypes && !DebugMode {
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
				color.New(color.FgHiRed).Printf("[片聽]")
			}
		}
	} else if result13.IsNaki && shanten >= 0 && shanten <= 2 {
		// 鳴牌时的无役提示（从听牌到两向听）
		fmt.Print(" ")
		color.New(color.FgHiRed).Printf("[無役]")
	}
}

// furiten hint
func fuRiTenHint() {
	if result13.FuritenRate > 0 {
		fmt.Print(" ")
		if result13.FuritenRate < 1 {
			color.New(color.FgHiYellow).Printf("[可能振聽]")
		} else {
			color.New(color.FgHiRed).Printf("[振聽]")
		}
	}
}

// count change drawing a useful tile(DAUT)
func countChangeDAUT() {
	if ShowScore {
		fmt.Print(" ")
		if len(result13.Improves) > 0 {
			fmt.Printf("[%2d改良]", len(result13.Improves))
		} else {
			fmt.Print(strings.Repeat(" ", 4))
			fmt.Print(strings.Repeat("　", 2)) // 全角空格
		}
	}

}

// type of drawing a useful til(DAUT)
func typeOfDAUT() {
	fmt.Printf(" %s \n", util.TilesToStrWithBracket(result13.Waits.AvailableTiles()))
}

// the hand can improve or not
func theHandCanImproveOrNot() {
	if ShowImproveDetail {
		for tile, waits := range result13.Improves {
			fmt.Printf("摸 %s 改良成 %s\n", util.Mahjong[tile], waits.String())
		}
	}
}
