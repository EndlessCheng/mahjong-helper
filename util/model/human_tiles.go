package model

// 用于命令行分析
type HumanTilesInfo struct {
	HumanTiles     string
	HumanDoraTiles string
}

func NewSimpleHumanTilesInfo(humanTiles string) *HumanTilesInfo {
	return &HumanTilesInfo{
		HumanTiles: humanTiles,
	}
}
