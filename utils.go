package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/fatih/color"
	"sort"
)

func _errorExit(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}

// e.g. "3m" => 2
func _convert(tile string) int {
	for i, m := range mahjong {
		if m == tile {
			return i
		}
	}
	_errorExit("参数错误:", tile)
	return -1
}

// e.g. "13m 24p" => (4, [0, 2, 10, 12])
func convert(tiles string) (num int, cnt []int) {
	cnt = make([]int, 34)

	tiles = strings.TrimSpace(tiles)
	splits := strings.Split(tiles, " ")

	var result []int
	for _, split := range splits {
		split = strings.TrimSpace(split)
		if split[0] >= '1' && split[0] <= '9' {
			for i := range split[:len(split)-1] {
				single := split[i:i+1] + split[len(split)-1:]
				result = append(result, _convert(single))
			}
		} else {
			result = append(result, _convert(split))
		}
	}

	for _, m := range result {
		cnt[m]++
		if cnt[m] > 4 {
			_errorExit("参数错误: 超过4张一样的牌！")
		}
	}

	return len(result), cnt
}

func in(a string, arr []string) bool {
	for _, _a := range arr {
		if _a == a {
			return true
		}
	}
	return false
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

//

func _getAttr(n float64) []color.Attribute {
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

func colorNumber(n int) {
	color.New(_getAttr(float64(n))...).Printf("%d", n)
}

func colorNumberF(n float64) {
	color.New(_getAttr(n)...).Printf("%.2f", n)
}
