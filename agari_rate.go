package main

// 仅为无筋数据，未考虑场况等，仅供参考
var agariTable = [...][4]float64{
	{26.3, 41.6, 50.1, 55.0},
	{19.2, 31.7, 38.2, 42.0},
	{14.8, 25.5, 32.0, 36.8},
	{11.8, 20.3, 26.7, 31.0},
	{11.8, 20.3, 26.7, 31.0},
	{11.8, 20.3, 26.7, 31.0},
	{14.8, 25.5, 32.0, 36.8},
	{19.2, 31.7, 38.2, 42.0},
	{26.3, 41.6, 50.1, 55.0},
}

func calcAgariRate(needs needTiles) float64 {
	agariRate := 0.0
	for idx, num := range needs {
		if num == 0 {
			continue
		}
		if idx > 27 {
			// 字牌的和率暂不考虑，返回 -1
			return -1
		}
		idx %= 9
		rate := agariTable[idx][num-1]
		agariRate = agariRate + rate - agariRate*rate/100
	}
	return agariRate
}
