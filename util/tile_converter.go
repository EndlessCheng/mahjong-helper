package util

import (
	"strings"
	"sort"
	"errors"
	"fmt"
)

// e.g. "3m" => 2
func StrToTile34(humanTile string) (tile34 int, err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New("[StrToTile34] 参数错误: " + humanTile)
		}
	}()

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
func StrToTiles34(humanTiles string) (tiles34 []int, err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New("[StrToTiles34] 参数错误: " + humanTiles)
		}
	}()

	// 在 mpsz 后面加上空格方便解析不含空格的 humanTiles
	humanTiles = strings.Replace(humanTiles, "m", "m ", -1)
	humanTiles = strings.Replace(humanTiles, "p", "p ", -1)
	humanTiles = strings.Replace(humanTiles, "s", "s ", -1)
	humanTiles = strings.Replace(humanTiles, "z", "z ", -1)

	humanTiles = strings.TrimSpace(humanTiles)
	if humanTiles == "" {
		return nil, errors.New("[StrToTiles34] 参数错误: 处理的手牌不能为空")
	}

	tiles34 = make([]int, 34)
	for _, split := range strings.Split(humanTiles, " ") {
		split = strings.TrimSpace(split)
		if split == "" {
			continue
		}
		if len(split) < 2 {
			return nil, errors.New("[StrToTiles34] 参数错误: " + humanTiles)
		}
		tileType := split[len(split)-1:]
		for _, c := range split[:len(split)-1] {
			tile := string(c) + tileType
			tile34, er := StrToTile34(tile)
			if er != nil {
				return nil, er
			}
			tiles34[tile34]++
			if tiles34[tile34] > 4 {
				return nil, fmt.Errorf("[StrToTiles34] 参数错误: %s 有超过 4 张一样的牌", humanTiles)
			}
		}
	}
	return
}

func MustStrToTiles34(humanTiles string) []int {
	tiles34, err := StrToTiles34(humanTiles)
	if err != nil {
		panic(err)
	}
	return tiles34
}

// e.g. "11122z" => [27, 27, 27, 28, 28]
func StrToTiles(humanTiles string) (tiles []int, err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New("[StrToTiles34] 参数错误: " + humanTiles)
		}
	}()

	humanTiles = strings.TrimSpace(humanTiles)
	if humanTiles == "" {
		return nil, errors.New("[StrToTiles34] 参数错误: 处理的手牌不能为空")
	}

	for _, split := range strings.Split(humanTiles, " ") {
		split = strings.TrimSpace(split)
		if len(split) < 2 {
			return nil, errors.New("[StrToTiles34] 参数错误: " + humanTiles)
		}
		tileType := split[len(split)-1:]
		for _, c := range split[:len(split)-1] {
			tile := string(c) + tileType
			tile34, er := StrToTile34(tile)
			if er != nil {
				return nil, er
			}
			tiles = append(tiles, tile34)
		}
	}
	return
}

func MustStrToTiles(humanTiles string) []int {
	tiles, err := StrToTiles(humanTiles)
	if err != nil {
		panic(err)
	}
	return tiles
}

//

// e.g. [9, 11, 27] => "13p 1z"
func TilesToStr(tiles []int) (humanTiles string) {
	sort.Ints(tiles)
	convert := func(lowerIndex, upperIndex int, endsWith string) {
		found := false
		for _, idx := range tiles {
			if idx >= lowerIndex && idx < upperIndex {
				found = true
				humanTiles += string('1' + idx - lowerIndex)
			}
		}
		if found {
			humanTiles += endsWith
		}
	}
	convert(0, 9, "m ")
	convert(9, 18, "p ")
	convert(18, 27, "s ")
	convert(27, 34, "z")
	return strings.TrimSpace(humanTiles)
}

func Tile34ToStr(tile34 int) string {
	return TilesToStr([]int{tile34})
}

func Tiles34ToStr(tiles34 []int) (humanTiles string) {
	merge := func(lowerIndex, upperIndex int, endsWith string) {
		found := false
		for i, c := range tiles34 {
			if i >= lowerIndex && i < upperIndex {
				for j := 0; j < c; j++ {
					found = true
					humanTiles += string('1' + i - lowerIndex)
				}
			}
		}
		if found {
			humanTiles += endsWith
		}
	}
	merge(0, 9, "m ")
	merge(9, 18, "p ")
	merge(18, 27, "s ")
	merge(27, 34, "z")
	return strings.TrimSpace(humanTiles)
}

// e.g. [9, 11, 27] => "[13p 1z]"
func TilesToStrWithBracket(tiles []int) string {
	return "[" + TilesToStr(tiles) + "]"
}

func Tiles34ToStrWithBracket(tiles34 []int) string {
	return "[" + Tiles34ToStr(tiles34) + "]"
}
