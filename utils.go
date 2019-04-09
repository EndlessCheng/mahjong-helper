package main

import (
	"fmt"
	"os"
		"github.com/fatih/color"
				"bufio"
)

var mahjong = [...]string{
	"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m",
	"1p", "2p", "3p", "4p", "5p", "6p", "7p", "8p", "9p",
	"1s", "2s", "3s", "4s", "5s", "6s", "7s", "8s", "9s",
	"1z", "2z", "3z", "4z", "5z", "6z", "7z",
}

var mahjongU = [...]string{
	"1M", "2M", "3M", "4M", "5M", "6M", "7M", "8M", "9M",
	"1p", "2p", "3p", "4p", "5p", "6p", "7p", "8p", "9p",
	"1S", "2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S",
	"1Z", "2Z", "3Z", "4Z", "5Z", "6Z", "7Z",
}

var mahjongZH = [...]string{
	"1万", "2万", "3万", "4万", "5万", "6万", "7万", "8万", "9万",
	"1饼", "2饼", "3饼", "4饼", "5饼", "6饼", "7饼", "8饼", "9饼",
	"1索", "2索", "3索", "4索", "5索", "6索", "7索", "8索", "9索",
	"东", "南", "西", "北", "白", "发", "中",
}

func _errorExit(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	fmt.Println("按任意键退出...")
	bufio.NewReader(os.Stdin).ReadByte()
	os.Exit(1)
}

//

func _getTingCountAttr(n float64) []color.Attribute {
	var colors []color.Attribute
	switch {
	case n < 14:
		colors = append(colors, color.FgBlue)
	case n <= 18:
		colors = append(colors, color.FgYellow)
	case n < 24:
		colors = append(colors, color.FgHiRed)
	default:
		colors = append(colors, color.FgRed)
	}
	return colors
}

// *
func getShantenWaitsCountColors(shanten int, waitsCount int) []color.Attribute {
	if shanten == 0 {
		return _getTingCountAttr(float64(waitsCount * 3))
	}
	div := 1
	for i := 1; i < shanten; i++ {
		div *= 2
	}
	return _getTingCountAttr(float64(waitsCount) / float64(div))
}

func getTingCountColor(count float64) color.Attribute {
	switch {
	case count < 5:
		return color.FgBlue
	case count <= 6:
		return color.FgYellow
	case count <= 8:
		return color.FgHiRed
	default:
		return color.FgRed
	}
}

func getNextShantenWaitsCountColor(shanten int, avgNextShantenWaitsCount float64) color.Attribute {
	div := 1
	for i := 1; i < shanten; i++ {
		div *= 3
	}
	return getTingCountColor(avgNextShantenWaitsCount / float64(div))
}

func getSimpleRiskColor(index int) color.Attribute {
	if index >= 27 {
		return color.FgHiBlue
	} else {
		_i := index%9 + 1
		switch _i {
		case 1, 9:
			return color.FgBlue
		case 2, 8:
			return color.FgYellow
		case 3, 7:
			return color.FgHiYellow
		case 4, 5, 6:
			return color.FgRed
		default:
			_errorExit("代码有误: _i = ", _i)
		}
	}
	return -1
}

func getDiscardAlertColor(index int) color.Attribute {
	if index >= 27 {
		return color.FgWhite
	} else {
		_i := index%9 + 1
		switch _i {
		case 1, 9:
			return color.FgWhite
		case 2, 8:
			return color.FgYellow
		case 3, 7:
			return color.FgHiYellow
		case 4, 5, 6:
			return color.FgRed
		default:
			_errorExit("代码有误: _i = ", _i)
		}
	}
	return -1
}

func getSafeColor(index int) color.Attribute {
	if index >= 27 {
		return color.FgRed
	} else {
		_i := index%9 + 1
		switch _i {
		case 1, 9:
			return color.FgRed
		case 2, 8:
			return color.FgHiYellow
		case 3, 7:
			return color.FgYellow
		case 4, 5, 6:
			return color.FgBlue
		default:
			_errorExit("代码有误: _i = ", _i)
		}
	}
	return -1
}

func getNumRiskColor(risk float64) color.Attribute {
	switch {
	case risk < 3:
		return color.FgHiBlue
	case risk < 5:
		return color.FgBlue
	case risk < 7.5:
		return color.FgYellow
	case risk < 10:
		return color.FgHiYellow
	case risk < 15:
		return color.FgHiRed
	default:
		return color.FgRed
	}
}

//

func lower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		c += 32
	}
	return c
}

func upper(c byte) byte {
	if c >= 'a' && c <= 'z' {
		c -= 32
	}
	return c
}
