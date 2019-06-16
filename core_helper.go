package main

var debugMode = false

type gameMode int

const (
	// TODO: 感觉有点杂乱需要重构
	gameModeMatch       gameMode = iota // 对战 - IsInit
	gameModeRecord                      // 解析牌谱
	gameModeRecordCache                 // 解析牌谱 - runMajsoulRecordAnalysisTask
	gameModeLive                        // 解析观战
)

const (
	dataSourceTypeTenhou = iota
	dataSourceTypeMajsoul
)

const (
	meldTypeChi    = iota // 吃
	meldTypePon           // 碰
	meldTypeAnkan         // 暗杠
	meldTypeMinkan        // 大明杠
	meldTypeKakan         // 加杠
)

// 负数变正数
func normalDiscardTiles(discardTiles []int) []int {
	newD := make([]int, len(discardTiles))
	copy(newD, discardTiles)
	for i, discardTile := range newD {
		if discardTile < 0 {
			newD[i] = ^discardTile
		}
	}
	return newD
}
