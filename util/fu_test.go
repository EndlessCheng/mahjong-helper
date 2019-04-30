package util

import (
	"testing"
	"fmt"
)

func TestDivideResult_Fu(t *testing.T) {
	for _, tiles := range []string{
		"123456789m 12344s",     // 一气 [44s 123m 456m 789m 123s]
		"111234678m 11122z",     // [22z 111m 111z 234m 678m]
		"22334455m 234s 234p",   // [22m 345m 345m 234p 234s] [55m 234m 234m 234p 234s]
		"111222333m 234s 11z",   // [11z 111m 222m 333m 234s] [11z 123m 123m 123m 234s]
		"112233m 112233p 11z",   // 两杯口，不是七对子
		"11223344556677z",       // 七对子
		"119m 19p 19s 1234567z", // 国士无双自行判断
		"11m 345p",
		"1122m",
		"11m 112233p",
		"11m 123456789p",
		"11m 111p 111s",
		"111m 11p 111s",
		"111m 111p 11s",
	} {
		fmt.Print(tiles + " = ")
		tiles34 := MustStrToTiles34(tiles)
		results := DivideTiles34(tiles34)
		if len(results) == 0 {
			fmt.Println("[国士/未和牌]")
			continue
		}
		for _, result := range results {
			hi := &HandInfo{&_handInfo{},result}
			fmt.Printf("%s %#v %d 符, ", result.String(), result, hi.Fu())
		}
		fmt.Println()
	}
}
