package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"strconv"
)

func TestCalculateRiskTiles34(t *testing.T) {
	t.Skip()
	safeTiles34 := make([]bool, 34)
	for _, tile := range MustStrToTiles("4m 5p 3s") {
		safeTiles34[tile] = true
	}
	leftTiles34 := InitLeftTiles34WithTiles34(MustStrToTiles34("11112222p 7777s"))
	risk34 := CalculateRiskTiles34(9, safeTiles34, leftTiles34, nil, 27, 28)
	for i, risk := range risk34 {
		t.Log(Mahjong[i], risk)
	}
}

func TestCalculateRiskTiles34_19(t *testing.T) {
	t.Skip()
	safeTiles34 := make([]bool, 34)
	for _, tile := range MustStrToTiles("4m") {
		safeTiles34[tile] = true
	}
	leftTiles34 := InitLeftTiles34WithTiles34(MustStrToTiles34("11117777m"))
	risk34 := CalculateRiskTiles34(9, safeTiles34, leftTiles34, nil, 27, 28)
	for i, risk := range risk34 {
		t.Log(Mahjong[i], risk)
	}
}

// 寻找绝对安牌（不考虑国士）
func TestCalculateRiskTiles34_SafeTile(t *testing.T) {
	assert := assert.New(t)

	safeTiles34 := make([]bool, 34)
	leftTiles34 := InitLeftTiles34WithTiles34(MustStrToTiles34("2222333377779999m 22228888p 333355557777s 4444z"))
	risk34 := CalculateRiskTiles34(9, safeTiles34, leftTiles34, nil, 27, 28)
	//for i, risk := range risk34 {
	//	t.Log(Mahjong[i], risk)
	//}
	for i, risk := range risk34 {
		if InInts(i, MustStrToTiles("29m 4z")) {
			assert.Zero(risk, strconv.Itoa(i))
		} else {
			assert.NotZero(risk, strconv.Itoa(i))
		}
	}

	safeTiles34 = make([]bool, 34)
	for _, tile := range MustStrToTiles("4m 5p 6s") {
		safeTiles34[tile] = true
	}
	leftTiles34 = InitLeftTiles34WithTiles34(MustStrToTiles34("111177778888m 1111222288889999p 222233339999s"))
	risk34 = CalculateRiskTiles34(9, safeTiles34, leftTiles34, nil, 27, 28)
	//for i, risk := range risk34 {
	//	t.Log(Mahjong[i], risk)
	//}
	for i, risk := range risk34 {
		if InInts(i, MustStrToTiles("1478m 12589p 2369s")) {
			assert.Zero(risk, strconv.Itoa(i))
		} else {
			assert.NotZero(risk, strconv.Itoa(i))
		}
	}
}

func TestCalculateLeftNoSujiTiles(t *testing.T) {
	assert := assert.New(t)

	safeTiles34 := make([]bool, 34)

	// 初始（18 无筋）
	leftTiles34 := InitLeftTiles34()
	assert.Equal("123789m 123789p 123789s", TilesToStr(CalculateLeftNoSujiTiles(safeTiles34, leftTiles34)))

	// 断幺壁
	leftTiles34 = InitLeftTiles34WithTiles34(MustStrToTiles34("33337777m 22228888p 5555s"))
	assert.Equal("37p 1289s", TilesToStr(CalculateLeftNoSujiTiles(safeTiles34, leftTiles34)))

	// 幺九壁
	leftTiles34 = InitLeftTiles34WithTiles34(MustStrToTiles34("11119999m 11119999p 11119999s"))
	assert.Equal("2378m 2378p 2378s", TilesToStr(CalculateLeftNoSujiTiles(safeTiles34, leftTiles34)))

	// 现物（0 无筋）
	discardTiles34 := MustStrToTiles34("123789m 123789p 123789s")
	leftTiles34 = InitLeftTiles34()
	for i, c := range discardTiles34 {
		safeTiles34[i] = c > 0
	}
	assert.Equal("", TilesToStr(CalculateLeftNoSujiTiles(safeTiles34, leftTiles34)))

	// 筋（0 无筋）
	discardTiles34 = MustStrToTiles34("456m 456p 456s")
	leftTiles34 = InitLeftTiles34()
	for i, c := range discardTiles34 {
		safeTiles34[i] = c > 0
	}
	assert.Equal("", TilesToStr(CalculateLeftNoSujiTiles(safeTiles34, leftTiles34)))
}
