package util

import (
	"testing"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"fmt"
)

func Test_search13(t *testing.T) {
	humanTiles := "5555m"
	tiles34 := MustStrToTiles34(humanTiles)
	pi := model.NewSimplePlayerInfo(tiles34, nil)
	shanten := CalculateShanten(tiles34)
	fmt.Println(NumberToChineseShanten(shanten))
	fmt.Print(_search13(shanten, pi, -1))
}

func Test_searchShanten14(t *testing.T) {
	humanTiles := "12688m 33579p 24s 56z"
	tiles34 := MustStrToTiles34(humanTiles)
	pi := model.NewSimplePlayerInfo(tiles34, nil)
	shanten := CalculateShanten(tiles34)
	fmt.Println(NumberToChineseShanten(shanten))
	fmt.Print(searchShanten14(shanten, pi, shanten-1))
	fmt.Println("倒退回" + NumberToChineseShanten(shanten+1))
	fmt.Print(searchShanten14(shanten+1, pi, shanten))
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
	tiles34 := MustStrToTiles34("55678m 3467p 24668s")
	shanten := CalculateShanten(tiles34)
	pi := model.NewSimplePlayerInfo(tiles34, nil)
	for i := 0; i < b.N; i++ {
		// 19,343,607 ns/op
		searchShanten14(shanten, pi, -1)
	}
}

func BenchmarkSearchShanten3(b *testing.B) {
	tiles34 := MustStrToTiles34("12688m 33579p 24s 56z")
	shanten := CalculateShanten(tiles34)
	pi := model.NewSimplePlayerInfo(tiles34, nil)
	for i := 0; i < b.N; i++ {
		// 92,369,360 ns/op
		searchShanten14(shanten, pi, -1)
	}
}
