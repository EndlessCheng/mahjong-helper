package main

import (
	"fmt"
	"strings"

	"github.com/EndlessCheng/mahjong-helper/util"
	"github.com/fatih/color"
)

func printAccountInfo(accountID int) {
	fmt.Printf("您的账号 ID 为 ")
	color.New(color.FgHiGreen).Printf("%d", accountID)
	fmt.Printf("，该数字为雀魂服务器账号数据库中的 ID，该值越小表示您的注册时间越早\n")
}

type handsRisk struct {
	tile int
	risk float64
}

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
	util.YakuSuuAnkou,      // 四暗刻
	util.YakuSuuAnkouTanki, // 四暗刻單騎
	util.YakuDaisangen,     // 大三元
	util.YakuShousuushii,   // 小四喜
	util.YakuDaisuushii,    // 大四喜
	util.YakuTsuuiisou,     // 字一色
	util.YakuChinroutou,    // 清老頭
	util.YakuRyuuiisou,     // 綠一色
	util.YakuChuuren,       // 九蓮寶燈
	util.YakuChuuren9,      // 純正九蓮寶燈
	util.YakuSuuKantsu,
	// util.YakuTenhou,  // 天胡
	// util.YakuChiihou, // 地胡

	util.YakuChiitoi,        // 七對子
	util.YakuPinfu,          // 平胡
	util.YakuRyanpeikou,     // 二盃口
	util.YakuIipeikou,       // 一盃口
	util.YakuSanshokuDoujun, // 三色同順
	util.YakuIttsuu,         // 一氣貫通
	util.YakuToitoi,         // 對對胡
	util.YakuSanAnkou,       // 三暗刻
	util.YakuSanshokuDoukou, // 三色同刻
	util.YakuSanKantsu,      // 三槓子
	util.YakuTanyao,         // 斷么九
	util.YakuChanta,         // 混全帶么九
	util.YakuJunchan,        // 純全帶么九
	util.YakuHonroutou,      // 混老頭
	util.YakuShousangen,     // 小三元
	util.YakuHonitsu,        // 混一色
	util.YakuChinitsu,       // 清一色

	util.YakuShiiaruraotai, // 十二落台
	util.YakuUumensai,      // 五門齊
	util.YakuSanrenkou,     // 三連刻
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
			fmt.Printf("進張")
		} else { // shanten == 1
			fmt.Printf("數")
			if ShowAgariAboveShanten1 {
				fmt.Printf("（%.2f%% 参考和率）", result13.AvgAgariRate)
			}
		}
		if ShowScore {
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

func PrintResults14WithRisk(results14 util.Hand14AnalysisResultList, mixedRiskTable riskTable) {
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
		fmt.Print("鳴牌後")
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
			color.HiGreen("注意切牌選擇：打點")
		}
	}

	// FIXME: 选择很多时如何精简何切选项？
	//const maxShown = 10
	//if len(results14) > maxShown { // 限制输出数量
	//	results14 = results14[:maxShown]
	//}
	for _, result := range results14 {
		r := &AnalysisResult{
			result.DiscardTile,
			result.IsDiscardDoraTile,
			result.OpenTiles,
			result.Result13,
			mixedRiskTable,
			result.Result13.AvgImproveWaitsCount == maxAvgImproveWaitsCount,
			result.Result13.MixedWaitsScore == maxMixedScore,
		}
		r.PrintWaitsWithImproves13_oneRow()
	}
}
