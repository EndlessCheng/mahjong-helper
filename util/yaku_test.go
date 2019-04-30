package util

import (
	"testing"
	"fmt"
)

func TestFindNormalYaku(t *testing.T) {
	for _, tiles := range []string{
		"123456789m 12344s",     // [44s 123m 456m 789m 123s]
		"111234678m 11122z",     // [22z 111m 111z 234m 678m]
		"22334455m 234s 234p",   // [22m 345m 345m 234p 234s] [55m 234m 234m 234p 234s]
		"111222333m 234s 11z",   // [11z 111m 222m 333m 234s] [11z 123m 123m 123m 234s]
		"112233m 112233p 11z",   // 两杯口，不是七对子
		"11223344556677z",       // 七对子
		"119m 19p 19s 1234567z", // 自行判断
		"234m 234p 23466999s",
	} {
		fmt.Print(tiles + " = ")
		tiles34 := MustStrToTiles34(tiles)
		for _, result := range DivideTiles34(tiles34) {
			fmt.Printf("%s ", result.String())
			fmt.Printf("%v", FindYakuList(&HandInfo{
				Divide: result,
			}))
		}
		fmt.Println()
	}
}
