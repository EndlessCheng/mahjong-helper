package util

type tileType int8

const (
	tileTypeNoSuji5 tileType = iota
	tileTypeNoSuji46
	tileTypeNoSuji37
	tileTypeNoSuji28
	tileTypeNoSuji19
	tileTypeHalfSuji5
	tileTypeHalfSuji46A  // 现物为19
	tileTypeHalfSuji46B  // 现物为73
	tileTypeSuji37
	tileTypeSuji28
	tileTypeSuji19
	tileTypeDoubleSuji5
	tileTypeDoubleSuji46
	tileTypeYakuHaiLeft3  // 字牌-役牌-生牌
	tileTypeYakuHaiLeft2  // 字牌-役牌-现1枚
	tileTypeYakuHaiLeft1  // 字牌-役牌-现2枚
	tileTypeOtaHaiLeft3   // 字牌-客风-生牌
	tileTypeOtaHaiLeft2   // 字牌-客风-现1枚
	tileTypeOtaHaiLeft1   // 字牌-客风-现2枚
)

// [巡目][类型]
var RiskData = [][]float64{
	{},
	{5.7, 5.7, 5.8, 4.7, 3.4, 2.5, 2.5, 3.1, 5.6, 3.8, 1.8, -1, -1, 2.1, 1.2, 0.5, 2.4, 1.4, 1.2}, // 1
	{6.6, 6.9, 6.3, 5.2, 4.0, 3.5, 3.5, 4.1, 5.3, 3.5, 1.9, 0.8, 2.6, 2.3, 1.2, 0.5, 2.7, 1.3, 0.4},
	{7.7, 8.0, 6.7, 5.8, 4.6, 4.3, 4.1, 4.9, 5.2, 3.6, 1.8, 1.6, 2.0, 2.4, 1.2, 0.3, 2.6, 1.2, 0.3},
	{8.5, 8.9, 7.1, 6.2, 5.1, 4.8, 4.7, 5.6, 5.2, 3.8, 1.7, 1.6, 2.0, 2.6, 1.1, 0.2, 2.6, 1.2, 0.2},
	{9.4, 9.7, 7.5, 6.7, 5.5, 5.3, 5.1, 6.0, 5.3, 3.7, 1.7, 1.7, 2.0, 2.9, 1.2, 0.2, 2.8, 1.2, 0.2}, // 5
	{10.2, 10.5, 7.9, 7.1, 5.9, 5.8, 5.6, 6.4, 5.2, 3.7, 1.7, 1.8, 2.0, 3.2, 1.3, 0.2, 2.9, 1.3, 0.2},
	{11.0, 11.3, 8.4, 7.5, 6.3, 6.3, 6.1, 6.8, 5.3, 3.7, 1.7, 2.0, 2.1, 3.6, 1.4, 0.2, 3.2, 1.4, 0.2},
	{11.9, 12.2, 8.9, 8.0, 6.8, 6.9, 6.6, 7.4, 5.3, 3.8, 1.7, 2.1, 2.2, 4.0, 1.6, 0.2, 3.5, 1.6, 0.2},
	{12.8, 13.1, 9.5, 8.6, 7.4, 7.4, 7.2, 7.9, 5.5, 3.9, 1.8, 2.2, 2.3, 4.6, 1.9, 0.3, 4.0, 1.8, 0.2}, // 9
	{13.8, 14.1, 10.1, 9.2, 8.0, 8.0, 7.8, 8.5, 5.6, 4.0, 1.9, 2.4, 2.4, 5.3, 2.2, 0.3, 4.6, 2.1, 0.3},
	{14.9, 15.1, 10.8, 9.9, 8.7, 8.7, 8.5, 9.2, 5.7, 4.2, 2.0, 2.5, 2.6, 6.0, 2.6, 0.4, 5.1, 2.5, 0.3},
	{16.0, 16.3, 11.6, 10.6, 9.4, 9.4, 9.2, 9.9, 6.0, 4.4, 2.2, 2.7, 2.7, 6.8, 3.1, 0.4, 5.1, 2.5, 0.3},
	{17.2, 17.5, 12.4, 11.4, 10.2, 10.2, 10.0, 10.6, 6.2, 4.6, 2.4, 3.0, 3.0, 7.8, 3.7, 0.5, 6.6, 3.7, 0.5},
	{18.5, 18.8, 13.3, 12.3, 11.1, 11.0, 10.9, 11.4, 6.6, 4.9, 2.7, 3.2, 3.1, 8.8, 4.4, 0.7, 7.4, 4.4, 0.6},
	{19.9, 20.1, 14.3, 13.3, 12.0, 11.9, 11.8, 12.3, 7.0, 5.3, 3.0, 3.4, 3.4, 9.9, 5.2, 0.8, 8.4, 5.3, 0.8}, // 15
	{21.3, 21.7, 15.4, 14.3, 13.1, 12.9, 12.8, 13.3, 7.4, 5.7, 3.3, 3.7, 3.6, 11.2, 6.2, 1.0, 9.4, 6.5, 0.9},
	{22.9, 23.2, 16.6, 15.4, 14.2, 14.0, 13.8, 14.4, 8.0, 6.1, 3.6, 3.9, 3.9, 12.4, 7.3, 1.3, 10.5, 7.7, 1.2},
	{24.7, 24.9, 17.9, 16.7, 15.4, 15.2, 15.0, 15.6, 8.5, 6.6, 4.0, 4.3, 4.2, 13.9, 8.5, 1.7, 11.8, 9.4, 1.6}, // 18
	{27.5, 27.8, 20.4, 19.1, 17.8, 17.5, 17.5, 17.5, 9.8, 7.4, 5.0, 5.1, 5.1, 18.1, 12.1, 2.8, 14.7, 12.6, 2.1},
}
var MaxTurns = len(RiskData) - 1

var (
	// [需要判断危险度的牌号(0-8)][是否有对应的现物(0-1或0-3)]
	// 123789: 无现物，有现物
	// 4: 无17现物，无1有7，有1无7，有17
	// 56: 同上
	TileTypeTable = [][]tileType{
		{tileTypeNoSuji19, tileTypeSuji19},
		{tileTypeNoSuji28, tileTypeSuji28},
		{tileTypeNoSuji37, tileTypeSuji37},
		{tileTypeNoSuji46, tileTypeHalfSuji46B, tileTypeHalfSuji46A, tileTypeDoubleSuji46},
		{tileTypeNoSuji5, tileTypeHalfSuji5, tileTypeHalfSuji5, tileTypeDoubleSuji5},
		{tileTypeNoSuji46, tileTypeHalfSuji46A, tileTypeHalfSuji46B, tileTypeDoubleSuji46},
		{tileTypeNoSuji37, tileTypeSuji37},
		{tileTypeNoSuji28, tileTypeSuji28},
		{tileTypeNoSuji19, tileTypeSuji19},
	}
	// [是否为役牌(0-1)][剩余数-1]
	HonorTileType = [][]tileType{
		{tileTypeOtaHaiLeft1, tileTypeOtaHaiLeft2, tileTypeOtaHaiLeft3, tileTypeOtaHaiLeft3},
		{tileTypeYakuHaiLeft1, tileTypeYakuHaiLeft2, tileTypeYakuHaiLeft3, tileTypeOtaHaiLeft3},
	}
)
