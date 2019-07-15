package util

import "math"

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func Lower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		c += 32
	}
	return c
}

func Upper(c byte) byte {
	if c >= 'a' && c <= 'z' {
		c -= 32
	}
	return c
}

func ByteAtStr(b byte, s string) int {
	for i, _b := range []byte(s) {
		if _b == b {
			return i
		}
	}
	return -1
}

func InInts(e int, arr []int) bool {
	for _, _e := range arr {
		if e == _e {
			return true
		}
	}
	return false
}

func InStrings(e string, arr []string) bool {
	for _, _e := range arr {
		if e == _e {
			return true
		}
	}
	return false
}

// 258m 258p 258s 12345z 在不考虑国士无双和七对子时为八向听
var chineseShanten = []string{"和了", "听牌", "一向听", "两向听", "三向听", "四向听", "五向听", "六向听", "七向听", "八向听"}

// -1=和了，0=和牌，1=一向听，……
func NumberToChineseShanten(num int) string {
	return chineseShanten[num+1]
}

func rateAboveOne(x, y int) float64 {
	return rateAboveOneFloat64(float64(x), float64(y))
}

func rateAboveOneFloat64(x, y float64) float64 {
	if x == y {
		return 1
	}
	if x == 0 || y == 0 {
		return math.MaxFloat64
	}
	if x > y {
		return x / y
	}
	return y / x
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func InDelta(a, b, delta float64) bool {
	return math.Abs(a-b) < delta
}

func Equal(a, b float64) bool {
	const eps = 1e-5
	return InDelta(a, b, eps)
}
