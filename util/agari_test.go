package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"strings"
)

func TestIsAgari(t *testing.T) {
	for _, humanTiles := range []string{
		"123456789m 12344s",
		"11122345678999s",
		"111234678m 11122z",
		"22334455m 234s 234p",
		"111222333m 234s 11z",
		"112233m 112233p 11z",
		"11223344556677z",   // 七对子
		"1133556699m 1122s", // 七对子
		"11m 345p",
		"11m 112233p",
		"11m 123456789p",
		"11m 111p 111s",
		"111m 11p 111s",
		"111m 111p 11s",
	} {
		assert.True(t, IsAgari(MustStrToTiles34(humanTiles)), humanTiles)
	}

	for _, humanTiles := range []string{
		"119m 19p 19s 1234567z", // 国士无双自行判断
		"1133555599m 1122s",
		"1122m",
		"8888p",
		"7777z",
		"66778p 1122345s 77z",
	} {
		assert.False(t, IsAgari(MustStrToTiles34(humanTiles)), humanTiles)
	}
}

func TestDivideTiles34(t *testing.T) {
	assert := assert.New(t)

	const otherDivideResult = "国士 or 未和牌"
	divideTiles := func(humanTiles string) string {
		drs := DivideTiles34(MustStrToTiles34(humanTiles))
		if len(drs) == 0 {
			return otherDivideResult
		}
		results := []string{}
		for _, dr := range drs {
			results = append(results, dr.String())
		}
		return strings.Join(results, ", ")
	}

	assert.Equal("[七对子]", divideTiles("11223344556677z"))
	assert.Equal("[七对子]", divideTiles("223344m 11335577s"))
	assert.Equal("[99s 111s 123s 456s 789s][九莲宝灯][一气通贯]", divideTiles("11112345678999s"))
	assert.Equal("[22s 111s 999s 345s 678s][九莲宝灯]", divideTiles("11122345678999s"))
	assert.Equal("[11s 999s 123s 345s 678s][九莲宝灯]", divideTiles("11123345678999s"))
	assert.Equal("[99s 111s 234s 456s 789s][九莲宝灯]", divideTiles("11123445678999s"))
	assert.Equal("[55s 111s 999s 234s 678s][九莲宝灯]", divideTiles("11123455678999s"))
	assert.Equal("[11s 999s 123s 456s 789s][九莲宝灯][一气通贯]", divideTiles("11123456789999s"))
	assert.Equal("[44s 123m 456m 789m 123s][一气通贯]", divideTiles("123456789m 12344s"))
	assert.Equal("[11m 123p 456p 789p][一气通贯]", divideTiles("11m 123456789p"))
	assert.Equal("[11p 123p 456p 789p][一气通贯]", divideTiles("11123456789p"))
	assert.Equal("[11z 123m 123m 123p 123p][两杯口]", divideTiles("112233m 112233p 11z"))
	assert.Equal("[22m 345m 345m 234p 234s][一杯口], [55m 234m 234m 234p 234s][一杯口]", divideTiles("22334455m 234s 234p"))
	assert.Equal("[11z 111m 222m 333m 234s], [11z 123m 123m 123m 234s][一杯口]", divideTiles("111222333m 234s 11z"))
	assert.Equal("[11m 234m 234m], [44m 123m 123m]", divideTiles("11223344m"))
	assert.Equal("[22z 111m 111z 234m 678m]", divideTiles("111234678m 11122z"))
	assert.Equal("[11m 345p]", divideTiles("11m 345p"))
	assert.Equal("[55p]", divideTiles("55p"))
	assert.Equal("[11m 111p 111s]", divideTiles("11m 111p 111s"))
	assert.Equal("[11p 111m 111s]", divideTiles("111m 11p 111s"))
	assert.Equal("[11s 111m 111p]", divideTiles("111m 111p 11s"))

	assert.Equal(otherDivideResult, divideTiles("119m 19p 19s 1234567z"))

	assert.Equal(otherDivideResult, divideTiles("4888m 499p 134557s 4z"))
	assert.Equal(otherDivideResult, divideTiles("1122m"))
	assert.Equal(otherDivideResult, divideTiles("5m"))
}

func BenchmarkIsAgari(b *testing.B) {
	tiles34 := MustStrToTiles34("123456789m 12344s")
	for i := 0; i < b.N; i++ {
		// 83.9 ns/op
		IsAgari(tiles34)
	}
}

func BenchmarkDivideTiles34(b *testing.B) {
	tiles34 := MustStrToTiles34("123456789m 12344s")
	for i := 0; i < b.N; i++ {
		// 236 ns/op
		DivideTiles34(tiles34)
	}
}
