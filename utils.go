package main

import (
	"fmt"
	"os"
	"github.com/fatih/color"
	"bufio"
)

func errorExit(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	fmt.Println("按任意键退出...")
	bufio.NewReader(os.Stdin).ReadByte()
	os.Exit(1)
}

//

// 进张数优劣
func getWaitsCountColor(shanten int, waitsCount float64) color.Attribute {
	_getWaitsCountColor := func(fixedWaitsCount float64) color.Attribute {
		switch {
		case fixedWaitsCount < 13: // 4.3*3
			return color.FgHiCyan  // FgHiBlue FgHiCyan
		case fixedWaitsCount <= 18: // 6*3
			return color.FgHiYellow
		default: // >6*3
			return color.FgHiRed
		}
	}

	if shanten == 0 {
		return _getWaitsCountColor(waitsCount * 3)
	}
	weight := 1
	for i := 1; i < shanten; i++ {
		weight *= 2
	}
	return _getWaitsCountColor(waitsCount / float64(weight))
}

// 他家中张舍牌提示
func getOtherDiscardAlertColor(index int) color.Attribute {
	if index >= 27 {
		return color.FgWhite
	}
	switch index%9 + 1 {
	case 1, 2, 8, 9:
		return color.FgWhite
	case 3, 7:
		return color.FgHiYellow
	case 4, 5, 6:
		return color.FgHiRed
	default:
		panic(fmt.Errorf("[getOtherDiscardAlertColor] 代码有误: index = %d", index))
	}
}

// 铳率高低
func getNumRiskColor(risk float64) color.Attribute {
	switch {
	//case risk < 3:
	//	return color.FgHiBlue
	case risk < 5:
		return color.FgHiCyan
		//case risk < 7.5:
		//	return color.FgYellow
	case risk < 10:
		return color.FgHiYellow
	case risk < 15:
		return color.FgHiRed
	default:
		return color.FgRed
	}
}
