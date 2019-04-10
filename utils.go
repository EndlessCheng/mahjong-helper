package main

import (
	"fmt"
	"os"
		"github.com/fatih/color"
				"bufio"
)

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
		colors = append(colors, color.FgHiBlue)
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
		return color.FgHiBlue
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
			return color.FgHiBlue
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
			return color.FgHiBlue
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
		return color.FgHiCyan
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
