package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/fatih/color"
	"sort"
	"strconv"
	"errors"
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
	os.Exit(1)
}

// e.g. "3m" => 2
func _convert(tile string) (int, error) {
	for i, m := range mahjong {
		if m == tile {
			return i, nil
		}
	}
	return -1, fmt.Errorf("参数错误: %s", tile)
}

// e.g. "22m 24p" => (4, [0, 2, 0, 0, ...,0, 10, 12])
func convert(tiles string) (num int, counts []int, err error) {
	tiles = strings.TrimSpace(tiles)
	if tiles == "" {
		return 0, nil, errors.New("参数错误: 处理的手牌不能为空")
	}

	var indexes []int
	for _, split := range strings.Split(tiles, " ") {
		split = strings.TrimSpace(split)
		if len(split) <= 1 {
			return 0, nil, errors.New("参数错误: " + split)
		}
		for i := range split[:len(split)-1] {
			single := split[i:i+1] + split[len(split)-1:]
			tile, err := _convert(single)
			if err != nil {
				return -1, nil, err
			}
			indexes = append(indexes, tile)
		}
	}

	counts = make([]int, len(mahjong))
	for _, index := range indexes {
		counts[index]++
		if counts[index] > 4 {
			return 0, nil, errors.New("参数错误: 超过4张一样的牌！")
		}
	}

	return len(indexes), counts, nil
}

func countsToString(counts []int) (string, error) {
	if len(counts) != len(mahjong) {
		return "", fmt.Errorf("counts 长度必须为 %d", len(mahjong))
	}

	sb := strings.Builder{}
	for i, type_ := range [...]string{"m", "p", "s", "z"} {
		wrote := false
		for j := 0; j < 9; j++ {
			idx := 9*i + j
			if idx >= len(mahjong) {
				break
			}
			for k := 0; k < counts[idx]; k++ {
				sb.WriteString(strconv.Itoa(j + 1))
				wrote = true
			}
		}
		if wrote {
			sb.WriteString(type_ + " ")
		}
	}
	return strings.TrimSpace(sb.String()), nil
}

func inStrSlice(a string, arr []string) bool {
	for _, _a := range arr {
		if _a == a {
			return true
		}
	}
	return false
}

func inIntSlice(a int, arr []int) bool {
	for _, _a := range arr {
		if _a == a {
			return true
		}
	}
	return false
}

func intAtSlice(a int, arr []int) int {
	for i, _a := range arr {
		if _a == a {
			return i
		}
	}
	return -1
}

func uniqueStrings(strings []string) []string {
	u := make([]string, 0, len(strings))
	mp := make(map[string]struct{}, len(strings))

	for _, val := range strings {
		if _, ok := mp[val]; !ok {
			mp[val] = struct{}{}
			u = append(u, val)
		}
	}

	sort.Strings(u)
	return u
}

func countPairs(counts []int) int {
	pairs := 0
	for i := 0; i < len(mahjong); i++ {
		if counts[i] >= 2 {
			pairs++
		}
	}
	return pairs
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

func colorTing1Count(n int) {
	color.New(_getTingCountAttr(float64(n))...).Printf("%-6d", n)
}

func colorTing2Count(n int) {
	color.New(_getTingCountAttr(float64(n) / 2)...).Printf("%2d", n)
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

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func inInts(e int, arr []int) bool {
	for _, _e := range arr {
		if e == _e {
			return true
		}
	}
	return false
}

func stringToSlice(str string) []string {
	arr := make([]string, len(str))
	for i, c := range str {
		arr[i] = string(c)
	}
	return arr
}
