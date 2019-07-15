package model

const (
	TileTypeMan = 0
	TileTypePin = 1
	TileTypeSou = 2
)

// TODO: 其他的也移过来
func InitLeftTiles34WithTiles34(tiles34 []int) []int {
	leftTiles34 := make([]int, 34)
	for i, count := range tiles34 {
		leftTiles34[i] = 4 - count
	}
	return leftTiles34
}

// 根据宝牌指示牌计算出宝牌
// isSannin: 是否为三麻
func DoraTile(doraIndicator int, isSannin bool) (dora int) {
	if doraIndicator < 27 { // mps
		if doraIndicator%9 < 8 {
			if isSannin && doraIndicator == 0 {
				// 三麻的1m->9m
				return 8
			}
			return doraIndicator + 1
		}
		return doraIndicator - 8
	}
	if doraIndicator < 31 { // 东南西北
		if doraIndicator < 30 {
			return doraIndicator + 1
		}
		return 27
	}
	if doraIndicator < 33 { // 白发中
		return doraIndicator + 1
	}
	return 31
}

// 根据宝牌指示牌计算出宝牌
// isSannin: 是否为三麻
func DoraList(doraIndicators []int, isSannin bool) (doraList []int) {
	for _, doraIndicator := range doraIndicators {
		doraList = append(doraList, DoraTile(doraIndicator, isSannin))
	}
	return
}
