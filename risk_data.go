package main

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
var riskData = [][]float64{
	{},
	{5.7, 5.7, 5.8, 4.7, 3.4, 2.5, 2.5, 3.1, 5.6, 3.8, 1.8, -1, -1, 2.1, 1.2, 0.5, 2.4, 1.4, 1.2}, // 1
	{},
	{},
	{},
	{}, // 5
	{},
	{},
	{},
	{12.8, 13.1, 9.5, 8.6, 7.4, 7.4, 7.2, 7.9, 5.5, 3.9, 1.8, 2.2, 2.3, 4.6, 1.9, 0.3, 4.0, 1.8, 0.2}, // 9
	{},
	{},
	{},
	{},
	{},
	{}, // 15
	{},
	{},
	{}, // 18
	{},
}

var (
	// [需要判断危险度的牌号(0-8)][是否有对应的现物(0-1或0-3)]
	// 123789: 无现物，有现物
	// 4: 无17现物，无1有7，有1无7，有17
	// 56: 同上
	tileTypeTable = [][]tileType{
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
	ziTileType = [][]tileType{
		{tileTypeOtaHaiLeft1, tileTypeOtaHaiLeft2, tileTypeOtaHaiLeft3, tileTypeOtaHaiLeft3},
		{tileTypeYakuHaiLeft1, tileTypeYakuHaiLeft2, tileTypeYakuHaiLeft3, tileTypeOtaHaiLeft3},
	}
)

//  TODO  var noChanceMatrix = [
//            {'indices': [1], 'blocked_tiles': [0]},
//            {'indices': [2], 'blocked_tiles': [0, 1]},
//            {'indices': [3], 'blocked_tiles': [1, 2]},
//            {'indices': [4], 'blocked_tiles': [2, 6]},
//            {'indices': [5], 'blocked_tiles': [6, 7]},
//            {'indices': [6], 'blocked_tiles': [7, 8]},
//            {'indices': [7], 'blocked_tiles': [8]},
//            {'indices': [1, 5], 'blocked_tiles': [3]},
//            {'indices': [2, 6], 'blocked_tiles': [4]},
//            {'indices': [3, 7], 'blocked_tiles': [5]},
//            {'indices': [1, 4], 'blocked_tiles': [2, 3]},
//            {'indices': [2, 5], 'blocked_tiles': [3, 4]},
//            {'indices': [3, 6], 'blocked_tiles': [4, 5]},
//            {'indices': [4, 7], 'blocked_tiles': [5, 6]},
//        ]
