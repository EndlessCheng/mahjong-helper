package main

import (
	"github.com/EndlessCheng/mahjong-helper/util"
	"fmt"
	"strings"
	"github.com/fatih/color"
	"github.com/EndlessCheng/mahjong-helper/util/model"
)

func alertBackwardToShanten2(results util.Hand14AnalysisResultList, incShantenResults util.Hand14AnalysisResultList) {
	if len(results) == 0 || len(incShantenResults) == 0 {
		return
	}

	if results[0].Result13.Waits.AllCount() < 10 {
		if results[0].Result13.MixedWaitsScore < incShantenResults[0].Result13.MixedWaitsScore {
			color.HiGreen("建议向听倒退")
		}
	}
}

func _printIncShantenResults14(shanten int, incShantenResults14 util.Hand14AnalysisResultList, mixedRiskTable riskTable) {
	if len(incShantenResults14) == 0 {
		return
	}

	if len(incShantenResults14[0].OpenTiles) > 0 {
		fmt.Print("鸣牌后")
	}
	// "倒退回" +
	fmt.Println(util.NumberToChineseShanten(shanten+1) + "：")
	for _, result := range incShantenResults14 {
		printWaitsWithImproves13_oneRow(result.Result13, result.DiscardTile, result.OpenTiles, mixedRiskTable)
	}
}

func analysisTiles34(playerInfo *model.PlayerInfo, mixedRiskTable riskTable) error {
	humanTiles := util.Tiles34ToStr(playerInfo.HandTiles34)
	if len(playerInfo.Melds) > 0 {
		humanTiles += " &"
		for i := len(playerInfo.Melds) - 1; i >= 0; i-- {
			humanTiles += " " + util.TilesToStr(playerInfo.Melds[i].Tiles)
		}
	}
	fmt.Println(humanTiles)
	fmt.Println(strings.Repeat("=", len(humanTiles)))

	countOfTiles := util.CountOfTiles34(playerInfo.HandTiles34)
	switch countOfTiles % 3 {
	case 1:
		result := util.CalculateShantenWithImproves13(playerInfo)
		fmt.Println(util.NumberToChineseShanten(result.Shanten) + "：")
		printWaitsWithImproves13_oneRow(result, -1, nil, mixedRiskTable)
	case 2:
		shanten, results14, incShantenResults14 := util.CalculateShantenWithImproves14(playerInfo)

		if shanten == -1 {
			color.HiRed("【已胡牌】")
			break
		}

		if shanten == 0 {
			if len(results14) > 0 {
				r13 := results14[0].Result13
				if r13.RiichiPoint > 0 && r13.FuritenRate == 0 && r13.DamaPoint >= 5200 {
					color.HiGreen("默听打点充足：追求和率默听，追求打点立直")
				}
				// 局收支相近时，提示：局收支相近，追求和率打xx，追求打点打xx
			}
		} else if shanten == 1 {
			if len(playerInfo.DiscardTiles) < 9 {
				alertBackwardToShanten2(results14, incShantenResults14)
			}
		}

		if len(results14) > 0 {
			fmt.Println(util.NumberToChineseShanten(shanten) + "：")
			for _, result := range results14 {
				printWaitsWithImproves13_oneRow(result.Result13, result.DiscardTile, result.OpenTiles, mixedRiskTable)
			}
		}
		if len(incShantenResults14) > 0 {
			_printIncShantenResults14(shanten, incShantenResults14, mixedRiskTable)
		}
	default:
		return fmt.Errorf("参数错误: %d 张牌", countOfTiles)
	}

	fmt.Println()

	return nil
}

// 分析鸣牌
// playerInfo: 自家信息
// targetTile34: 他家舍牌
// isRedFive: 此舍牌是否为赤5
// allowChi: 是否能吃
// mixedRiskTable: 危险度表
func analysisMeld(playerInfo *model.PlayerInfo, targetTile34 int, isRedFive bool, allowChi bool, mixedRiskTable riskTable) {
	// 原始手牌分析
	result := util.CalculateShantenWithImproves13(playerInfo)

	// 副露分析
	shanten, results14, incShantenResults14 := util.CalculateMeld(playerInfo, targetTile34, isRedFive, allowChi)

	if len(results14) == 0 && len(incShantenResults14) == 0 {
		return
	}

	raw := util.Tiles34ToStr(playerInfo.HandTiles34) + " + " + util.Tile34ToStr(targetTile34) + "?"
	fmt.Println(raw)
	fmt.Println(strings.Repeat("=", len(raw)))

	fmt.Println("当前" + util.NumberToChineseShanten(result.Shanten) + "：")
	printWaitsWithImproves13_oneRow(result, -1, nil, mixedRiskTable)

	if shanten == -1 {
		color.HiRed("【已胡牌】")
		return
	}

	fmt.Print("鸣牌后")

	if shanten == 0 {
		// 局收支相近时，提示：局收支相近，追求和率打xx，追求打点打xx
	} else if shanten == 1 {
		//if len(playerInfo.DiscardTiles) < 9 {
		//	alertBackwardToShanten2(results14, incShantenResults14)
		//}
	}

	// 打印结果
	// FIXME: 选择很多时如何精简何切选项？
	const maxShown = 10

	if len(results14) > 0 {
		fmt.Println(util.NumberToChineseShanten(shanten) + "：")
		shownResults14 := results14
		if len(shownResults14) > maxShown {
			shownResults14 = shownResults14[:maxShown]
		}
		for _, result := range shownResults14 {
			printWaitsWithImproves13_oneRow(result.Result13, result.DiscardTile, result.OpenTiles, mixedRiskTable)
		}
	}

	if len(incShantenResults14) > 0 {
		shownIncResults14 := incShantenResults14
		if len(shownIncResults14) > maxShown {
			shownIncResults14 = shownIncResults14[:maxShown]
		}
		_printIncShantenResults14(shanten, shownIncResults14, mixedRiskTable)
	}
}

func analysisHumanTiles(humanTilesInfo *model.HumanTilesInfo) (tiles34 []int, err error) {
	humanTiles := humanTilesInfo.HumanTiles
	doraTiles := []int{}
	if humanTilesInfo.HumanDoraTiles != "" {
		doraTiles = util.MustStrToTiles(humanTilesInfo.HumanDoraTiles)
	}

	splits := strings.Split(humanTiles, "+")
	if len(splits) == 2 {
		tiles34, err = util.StrToTiles34(splits[0])
		if err != nil {
			return
		}

		rawTargetTile := strings.TrimSpace(splits[1])
		if len(rawTargetTile) > 2 {
			rawTargetTile = rawTargetTile[:2]
		}
		var targetTile34 int
		targetTile34, err = util.StrToTile34(rawTargetTile)
		if err != nil {
			return
		}

		var melds []model.Meld
		//melds = append(melds, model.Meld{MeldType: model.MeldTypePon, Tiles: util.MustStrToTiles("777z")})
		playerInfo := model.NewSimplePlayerInfo(tiles34, melds)
		playerInfo.DoraTiles = doraTiles
		isRedFive := false
		analysisMeld(playerInfo, targetTile34, isRedFive, true, nil)
		return
	}

	tiles34, err = util.StrToTiles34(humanTiles)
	if err != nil {
		return
	}

	playerInfo := model.NewSimplePlayerInfo(tiles34, nil)
	playerInfo.DoraTiles = doraTiles
	//playerInfo.IsTsumo = true
	err = analysisTiles34(playerInfo, nil)
	return
}
