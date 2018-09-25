package main

import (
	"fmt"
	"os"
	"strings"
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
