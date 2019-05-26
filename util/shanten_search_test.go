package util

import (
	"testing"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"fmt"
)

func Test_search13(t *testing.T) {
	humanTiles := "5555m"
	tiles34 := MustStrToTiles34(humanTiles)
	shanten := CalculateShanten(tiles34)
	fmt.Print(_search13(shanten, model.NewSimplePlayerInfo(tiles34, nil)))
}

func TestSearchShanten(t *testing.T) {
	humanTiles := "1122334455667z 1m"
	tiles34 := MustStrToTiles34(humanTiles)
	shanten := CalculateShanten(tiles34)
	fmt.Print(searchShanten14(shanten, model.NewSimplePlayerInfo(tiles34, nil)))
}

func BenchmarkSearchShanten0(b *testing.B) {
	tiles34 := MustStrToTiles34("234788m 234567s 33z")
	shanten := CalculateShanten(tiles34)
	pi := model.NewSimplePlayerInfo(tiles34, nil)
	for i := 0; i < b.N; i++ {
		// 11,536 ns/op
		searchShanten14(shanten, pi)
	}
}

func BenchmarkSearchShanten1(b *testing.B) {
	tiles34 := MustStrToTiles34("33455m 668p 345667s")
	shanten := CalculateShanten(tiles34)
	pi := model.NewSimplePlayerInfo(tiles34, nil)
	for i := 0; i < b.N; i++ {
		// 361,680 ns/op
		searchShanten14(shanten, pi)
	}
}

func BenchmarkSearchShanten2(b *testing.B) {
	tiles34 := MustStrToTiles34("55678m 3467p 24668s")
	shanten := CalculateShanten(tiles34)
	pi := model.NewSimplePlayerInfo(tiles34, nil)
	for i := 0; i < b.N; i++ {
		// 19,343,607 ns/op
		searchShanten14(shanten, pi)
	}
}

func BenchmarkSearchShanten3(b *testing.B) {
	tiles34 := MustStrToTiles34("12688m 33579p 24s 56z")
	shanten := CalculateShanten(tiles34)
	pi := model.NewSimplePlayerInfo(tiles34, nil)
	for i := 0; i < b.N; i++ {
		// 92,369,360 ns/op
		searchShanten14(shanten, pi)
	}
}
