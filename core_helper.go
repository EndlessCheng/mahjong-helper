package main

var debugMode = false

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
