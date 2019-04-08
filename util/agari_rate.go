package util

var (
	// 仅为无筋数据，未考虑筋牌、早外、NC、是否为宝牌、其他场况等，仅供参考
	// 6~10巡目 [牌0-4][剩余数1-4]
	// from: 勝つための現代麻雀技術論
	agariTable = [...][4]float64{
		{26.3, 41.6, 50.1, 55.0},
		{19.2, 31.7, 38.2, 42.0},
		{14.8, 25.5, 32.0, 36.8},
		{11.8, 20.3, 26.7, 31.0},
		{11.8, 20.3, 26.7, 31.0},
	}

	// 8巡目 [剩余数1-3]
	// from:「統計学」のマージャン戦術
	// FIXME: 这条仅适用于单骑，双碰不适用
	honorTileAgariTable = [3]float64{47.5, 58.0, 49.5}
)

// selfDiscards: 自家舍牌，用于分析骗筋时的和率
func CalculateAgariRate(waits Waits, selfDiscards []int) float64 {
	agariRate := 0.0
	for tile, left := range waits {
		if left == 0 {
			continue
		}
		var rate float64
		if tile < 27 {
			tile %= 9
			if tile > 4 {
				tile = 8 - tile
			}
			rate = agariTable[tile][left-1]
		} else {
			rate = honorTileAgariTable[left-1]
		}
		agariRate = agariRate + rate - agariRate*rate/100
	}
	return agariRate
}
