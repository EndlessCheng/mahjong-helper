package util

import (
	"strings"
	"sort"
	"errors"
	"fmt"
)

// e.g. "3m" => 2
func StrToTile34(humanTile string) (int, error) {
	humanTile = strings.TrimSpace(humanTile)
	if len(humanTile) != 2 {
		return -1, errors.New("[StrToTile34] 参数错误: " + humanTile)
	}

	idx := byteAtStr(humanTile[1], "mpsz")
	if idx == -1 {
		return -1, errors.New("[StrToTile34] 参数错误: " + humanTile)
	}
	i := humanTile[0]
	if i == '0' {
		i = '5'
	}
	return 9*idx + int(i-'1'), nil
}

func MustStrToTile34(humanTile string) int {
	tile34, err := StrToTile34(humanTile)
	if err != nil {
		panic(err)
	}
	return tile34
}

// e.g. "224m 24p" => [0, 2, 0, 1, 0, ..., 1, 0, 1, ...]
func StrToTiles34(humanTiles string) ([]int, error) {
	humanTiles = strings.TrimSpace(humanTiles)
	if humanTiles == "" {
		return nil, errors.New("[StrToTiles34] 参数错误: 处理的手牌不能为空")
	}

	tiles34 := make([]int, 34)
	for _, split := range strings.Split(humanTiles, " ") {
		split = strings.TrimSpace(split)
		if len(split) < 2 {
			return nil, errors.New("[StrToTiles34] 参数错误: " + humanTiles)
		}
		tileType := split[len(split)-1:]
		for _, c := range split[:len(split)-1] {
			tile := string(c) + tileType
			tile34, err := StrToTile34(tile)
			if err != nil {
				return nil, err
			}
			tiles34[tile34]++
			if tiles34[tile34] > 4 {
				return nil, fmt.Errorf("[StrToTiles34] 参数错误: %s 有超过 4 张一样的牌", humanTiles)
			}
		}
	}
	return tiles34, nil
}

func MustStrToTiles34(tiles string) []int {
	tiles34, err := StrToTiles34(tiles)
	if err != nil {
		panic(err)
	}
	return tiles34
}

//

// e.g. [9, 11, 27] => "13p 1z"
func TilesToStr(tiles []int) (res string) {
	sort.Ints(tiles)
	convert := func(lowerIndex, upperIndex int, endsWith string) {
		found := false
		for _, idx := range tiles {
			if idx >= lowerIndex && idx < upperIndex {
				found = true
				res += string('1' + idx - lowerIndex)
			}
		}
		if found {
			res += endsWith
		}
	}
	convert(0, 9, "m ")
	convert(9, 18, "p ")
	convert(18, 27, "s ")
	convert(27, 34, "z")
	return strings.TrimSpace(res)
}

func Tile34ToStr(tile34 int) string {
	return TilesToStr([]int{tile34})
}

func Tiles34ToStr(tiles34 []int) (res string) {
	merge := func(lowerIndex, upperIndex int, endsWith string) {
		found := false
		for i, c := range tiles34 {
			if i >= lowerIndex && i < upperIndex {
				for j := 0; j < c; j++ {
					found = true
					res += string('1' + i - lowerIndex)
				}
			}
		}
		if found {
			res += endsWith
		}
	}
	merge(0, 9, "m ")
	merge(9, 18, "p ")
	merge(18, 27, "s ")
	merge(27, 34, "z")
	return strings.TrimSpace(res)
}

// e.g. [9, 11, 27] => "[13p 1z]"
func TilesToStrWithBracket(tiles []int) (res string) {
	return "[" + TilesToStr(tiles) + "]"
}

func Tiles34ToStrWithBracket(tiles34 []int) (res string) {
	return "[" + Tiles34ToStr(tiles34) + "]"
}
