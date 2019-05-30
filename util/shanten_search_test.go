package util

import (
	"testing"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func Test_search13(t *testing.T) {
	t.Skip()
	humanTiles := "5555m"
	humanTiles = "55678m 3467p 2466s"
	tiles34 := MustStrToTiles34(humanTiles)
	pi := model.NewSimplePlayerInfo(tiles34, nil)
	shanten := CalculateShanten(tiles34)
	fmt.Println(NumberToChineseShanten(shanten))
	fmt.Print(_search13(shanten, pi, shanten-1))
}

func Test_searchShanten14(t *testing.T) {
	t.Skip()
	humanTiles := "466m 234467p 77s"
	tiles34 := MustStrToTiles34(humanTiles)
	pi := model.NewSimplePlayerInfo(tiles34, nil)
	shanten := CalculateShanten(tiles34)
	fmt.Println(NumberToChineseShanten(shanten))
	fmt.Print(searchShanten14(shanten, pi, -1))
	fmt.Println("倒退回" + NumberToChineseShanten(shanten+1))
	fmt.Print(searchShanten14(shanten+1, pi, -1))
}

func TestCalculateShantenAndWaits13(t *testing.T) {
	assert := assert.New(t)

	toString := func(shanten int, waits Waits) string {
		return NumberToChineseShanten(shanten) + " " + waits.String()
	}

	// closed
	assert.Equal("听牌 3 进张 [7z]", toString(CalculateShantenAndWaits13(MustStrToTiles34("1122334455667z"), nil)))
	assert.Equal("听牌 4 进张 [4s]", toString(CalculateShantenAndWaits13(MustStrToTiles34("123456789m 1135s"), nil)))
	assert.Equal("听牌 8 进张 [25s]", toString(CalculateShantenAndWaits13(MustStrToTiles34("123456789m 1134s"), nil)))
	assert.Equal("一向听 61 进张 [12345678m 47p 12345678s]", toString(CalculateShantenAndWaits13(MustStrToTiles34("3456m 3456s 44456p"), nil)))
	assert.Equal("两向听 12 进张 [1234z]", toString(CalculateShantenAndWaits13(MustStrToTiles34("123456789m 1234z"), nil)))
	assert.Equal("三向听 32 进张 [46m 2468p 24s]", toString(CalculateShantenAndWaits13(MustStrToTiles34("11357m 13579p 135s"), nil)))

	// open
	assert.Equal("听牌 3 进张 [5p]", toString(CalculateShantenAndWaits13(MustStrToTiles34("5p"), nil)))
	assert.Equal("听牌 6 进张 [14p]", toString(CalculateShantenAndWaits13(MustStrToTiles34("1234p"), nil)))
	assert.Equal("一向听 132 进张 [12346789m 123456789p 123456789s 1234567z]", toString(CalculateShantenAndWaits13(MustStrToTiles34("5555m"), nil)))
	assert.Equal("两向听 12 进张 [1234z]", toString(CalculateShantenAndWaits13(MustStrToTiles34("1234z"), nil)))
}

func BenchmarkSearchShanten0(b *testing.B) {
	tiles34 := MustStrToTiles34("234788m 234567s 33z")
	shanten := CalculateShanten(tiles34)
	pi := model.NewSimplePlayerInfo(tiles34, nil)
	for i := 0; i < b.N; i++ {
		// 11,536 ns/op
		searchShanten14(shanten, pi, -1)
	}
}

func BenchmarkSearchShanten1(b *testing.B) {
	tiles34 := MustStrToTiles34("33455m 668p 345667s")
	shanten := CalculateShanten(tiles34)
	pi := model.NewSimplePlayerInfo(tiles34, nil)
	for i := 0; i < b.N; i++ {
		// 361,680 ns/op
		searchShanten14(shanten, pi, -1)
	}
}

func BenchmarkSearchShanten2(b *testing.B) {
	tiles34 := MustStrToTiles34("4888m 499p 134557s 4z")
	shanten := CalculateShanten(tiles34)
	pi := model.NewSimplePlayerInfo(tiles34, nil)
	for i := 0; i < b.N; i++ {
		// 4,781,838 ns/op
		searchShanten14(shanten, pi, -1)
	}
}

func BenchmarkSearchShanten3(b *testing.B) {
	tiles34 := MustStrToTiles34("488m 499p 134557s 56z")
	shanten := CalculateShanten(tiles34)
	pi := model.NewSimplePlayerInfo(tiles34, nil)
	for i := 0; i < b.N; i++ {
		// 437,529,400 ns/op
		searchShanten14(shanten, pi, -1)
	}
}

func BenchmarkSearchShanten4(b *testing.B) {
	tiles34 := MustStrToTiles34("488m 49p 134557s 456z")
	shanten := CalculateShanten(tiles34)
	pi := model.NewSimplePlayerInfo(tiles34, nil)
	for i := 0; i < b.N; i++ {
		// 48,888,144,400 ns/op
		searchShanten14(shanten, pi, -1)
	}
}
