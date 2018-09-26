package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/fatih/color"
	"sort"
	"strconv"
)

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

// e.g. "13m 24p" => (4, [0, 2, 10, 12])
func convert(tiles string) (num int, cnt []int, err error) {
	cnt = make([]int, 34)

	tiles = strings.TrimSpace(tiles)
	splits := strings.Split(tiles, " ")

	var result []int
	for _, split := range splits {
		split = strings.TrimSpace(split)
		if split[0] >= '1' && split[0] <= '9' {
			for i := range split[:len(split)-1] {
				single := split[i:i+1] + split[len(split)-1:]
				tile, err := _convert(single)
				if err != nil {
					return -1, nil, err
				}
				result = append(result, tile)
			}
		} else {
			tile, err := _convert(split)
			if err != nil {
				return -1, nil, err
			}
			result = append(result, tile)
		}
	}

	for _, m := range result {
		cnt[m]++
		if cnt[m] > 4 {
			_errorExit("参数错误: 超过4张一样的牌！")
		}
	}

	return len(result), cnt, nil
}

func countToString(cnt []int) string {
	sb := strings.Builder{}
	for i, type_ := range [...]string{"m", "p", "s"} {
		wrote := false
		for j := 0; j < 9; j++ {
			for k := 0; k < cnt[9*i+j]; k++ {
				sb.WriteString(strconv.Itoa(j + 1))
				wrote = true
			}
		}
		if wrote {
			sb.WriteString(type_ + " ")
		}
	}
	for i := 27; i < len(mahjong); i++ {
		for k := 0; k < cnt[i]; k++ {
			sb.WriteString(mahjong[i] + " ")
		}
	}
	return strings.TrimSpace(sb.String())
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

func colorNumber1(n int) {
	color.New(_getAttr(float64(n))...).Printf("%d", n)
}

func colorNumber2(n int) {
	color.New(_getAttr(float64(n) / 2)...).Printf("%d", n)
}
