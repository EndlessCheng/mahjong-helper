package main

import (
	"fmt"
	"sort"

	"github.com/EndlessCheng/mahjong-helper/util"
	"github.com/fatih/color"
)

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
			color.New(GetNumRiskColor(hr.risk * fixedRiskMulti)).Printf(" " + util.MahjongZH[hr.tile])
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
