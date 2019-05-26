package util

const (
	// 参考:「統計学」のマージャン戦術 & 知るだけで強くなる麻雀の2択
	furitenBaseAgariRate = 5.9
	honorDoraAgariMulti  = 35.0 / 48.0
	numberDoraAgariMulti = 26.0 / 38.0
	ryanmenAgariMulti    = 0.91
)

var (
	// TODO: 基于巡目的和了率数据
	// TODO: 考虑读山的和了率？
	// TODO: 早外、NC、其他场况（其他家不要的牌）
	// https://github.com/EndlessCheng/mahjong-helper/issues/46

	// 数牌和率
	// 6~10巡目 [牌0-4][剩余数]
	// 参考: 勝つための現代麻雀技術論
	agariMap = map[tileType][5]float64{
		tileTypeNoSuji19:     {0, 26.3, 41.6, 50.1, 54.0},
		tileTypeNoSuji28:     {0, 19.2, 31.7, 38.2, 42.0},
		tileTypeNoSuji37:     {0, 14.8, 25.5, 32.0, 36.8},
		tileTypeNoSuji46:     {0, 11.8, 20.3, 26.7, 31.0},
		tileTypeNoSuji5:      {0, 11.8, 20.3, 26.7, 31.0},
		tileTypeSuji19:       {0, 36.1, 60.0, 67.9, 0},
		tileTypeSuji28:       {0, 24.9, 42.7, 51.2, 56.5},
		tileTypeSuji37:       {0, 17.2, 33.1, 43.5, 48.9},
		tileTypeDoubleSuji46: {0, 16.5, 35.5, 45.4, 50.0},
		tileTypeDoubleSuji5:  {0, 16.5, 35.5, 45.4, 50.0},
		tileTypeHalfSuji46A:  {0, 12.9, 24.7, 30.9, 35.4},
		tileTypeHalfSuji46B:  {0, 12.9, 24.7, 30.9, 35.4},
		tileTypeHalfSuji5:    {0, 12.9, 24.7, 30.9, 35.4},
	}

	// 字牌非单骑和率
	// 6~10巡目 [剩余数]
	// 参考: 勝つための現代麻雀技術論
	honorTileNonDankiAgariTable = [...]float64{0, 34.0, 52.0}

	// 字牌单骑和率
	// 8巡目 [剩余数]
	// 参考:「統計学」のマージャン戦術
	honorTileDankiAgariTable = [...]float64{0, 47.5, 58.0, 49.5}
)
