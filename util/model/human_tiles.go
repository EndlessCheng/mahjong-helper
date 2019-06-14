package model

import (
	"strings"
	"fmt"
)

// 用于命令行分析
type HumanTilesInfo struct {
	// 手牌 & 副露(暗杠用大写表示) + 要鸣的牌(可以吃)
	HumanTiles     string // 24688m 34s # 6666P 234p + 3m
	HumanDoraTiles string // 13m6p 不能有空格
	IsTsumo        bool

	HumanMelds      []string // 从 HumanTiles 解析出来的副露
	HumanTargetTile string   // 从 HumanTiles 解析出来的被鸣的牌
}

func NewSimpleHumanTilesInfo(humanTiles string) *HumanTilesInfo {
	return &HumanTilesInfo{
		HumanTiles: humanTiles,
	}
}

const (
	SepMeld       = "#"
	SepTargetTile = "+"
)

// 简单地处理 HumanTiles，拆分成一些子字符串
func (i *HumanTilesInfo) SelfParse() error {
	raw := strings.TrimSpace(i.HumanTiles)

	splits := strings.Split(raw, SepTargetTile)
	if len(splits) >= 2 {
		raw = strings.TrimSpace(splits[0])
		tile := strings.TrimSpace(splits[1])
		if len(tile) < 2 {
			return fmt.Errorf("输入错误: %s", i.HumanTiles)
		}
		i.HumanTargetTile = tile[:2]
	}

	splits = strings.Split(raw, SepMeld)
	if len(splits) >= 2 {
		raw = strings.TrimSpace(splits[0])
		humanMelds := strings.TrimSpace(splits[1])
		// 在 mpsz 后面加上空格方便解析不含空格的 humanTiles
		for _, tileType := range []string{"m", "p", "s", "z"} {
			humanMelds = strings.Replace(humanMelds, tileType, tileType+" ", -1)
			tileType = strings.ToUpper(tileType) // 暗杠
			humanMelds = strings.Replace(humanMelds, tileType, tileType+" ", -1)
		}
		humanMelds = strings.TrimSpace(humanMelds)
		for _, humanMeld := range strings.Split(humanMelds, " ") {
			if humanMeld != "" {
				i.HumanMelds = append(i.HumanMelds, humanMeld)
			}
		}
	}

	i.HumanTiles = raw
	return nil
}
