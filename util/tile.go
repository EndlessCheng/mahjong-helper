package util

import (
	"strings"
	"errors"
	"fmt"
	"sort"
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

// map[进张牌]剩余数
type Waits map[int]int

func (w Waits) AllCount() (count int) {
	for _, cnt := range w {
		count += cnt
	}
	return count
}

func (w Waits) indexes() []int {
	if len(w) == 0 {
		return nil
	}

	tileIndexes := make([]int, 0, len(w))
	for idx := range w {
		tileIndexes = append(tileIndexes, idx)
	}
	sort.Ints(tileIndexes)

	return tileIndexes
}

func (w Waits) ParseIndex() (allCount int, indexes []int) {
	return w.AllCount(), w.indexes()
}

func (w Waits) _parse(template [34]string) (allCount int, tiles []string) {
	if len(w) == 0 {
		return 0, nil
	}

	tileIndexes := make([]int, 0, len(w))
	for idx, cnt := range w {
		tileIndexes = append(tileIndexes, idx)
		allCount += cnt
	}
	sort.Ints(tileIndexes)

	tiles = make([]string, len(tileIndexes))
	for i, idx := range tileIndexes {
		tiles[i] = template[idx]
	}

	return allCount, tiles
}

func (w Waits) parse() (allCount int, tiles []string) {
	return w._parse(mahjong)
}

func (w Waits) parseZH() (allCount int, tilesZH []string) {
	return w._parse(mahjongZH)
}

func (w Waits) tilesZH() []string {
	_, tiles := w.parseZH()
	return tiles
}

func (w Waits) String() string {
	return fmt.Sprintf("%d 进张 %s", w.AllCount(), TilesToMergedStrWithBracket(w.indexes()))
}

func (w Waits) containAllIndexes(anotherNeeds Waits) bool {
	for k := range anotherNeeds {
		if _, ok := w[k]; !ok {
			return false
		}
	}
	return true
}

// 是否包含字牌
func (w Waits) containHonors() bool {
	indexes := w.indexes()
	if len(indexes) == 0 {
		return false
	}
	return indexes[len(indexes)-1] >= 27
}

func (w Waits) FixCountsWithLeftCounts(leftCounts []int) {
	if len(leftCounts) != 34 {
		return
	}
	for k := range w {
		w[k] = leftCounts[k]
	}
}

//

// TODO: 相关 1z<->27的转换代码，手牌解析等

func CountOfTiles(tiles []int) (count int) {
	for _, c := range tiles {
		count += c
	}
	return
}

func StrToTile34(tile string) (tile34 int, err error) {
	idx := byteAtStr(tile[1], "mpsz")
	if idx == -1 {
		return -1, fmt.Errorf("[StrToTile34] 参数错误: %s", tile)
	}
	i := tile[0]
	if i == '0' {
		i = '5'
	}
	return 9*idx + int(i-'1'), nil
}

// e.g. "22m 24p" => (4, [0, 2, 0, 0, ...,0, 10, 12])
func StrToTiles34(tiles string) (num int, tiles34 []int, err error) {
	tiles = strings.TrimSpace(tiles)
	if tiles == "" {
		return 0, nil, errors.New("参数错误: 处理的手牌不能为空")
	}

	var hands []int
	for _, split := range strings.Split(tiles, " ") {
		split = strings.TrimSpace(split)
		if len(split) <= 1 {
			return 0, nil, errors.New("参数错误: " + split)
		}
		for i := range split[:len(split)-1] {
			single := split[i:i+1] + split[len(split)-1:]
			tile, err := StrToTile34(single)
			if err != nil {
				return -1, nil, err
			}
			hands = append(hands, tile)
		}
	}

	tiles34 = make([]int, 34)
	for _, index := range hands {
		tiles34[index]++
		if tiles34[index] > 4 {
			return 0, nil, errors.New("参数错误: 超过4张一样的牌！")
		}
	}

	return len(hands), tiles34, nil
}

func MustStrToTiles34(tiles string) []int {
	_, tiles34, err := StrToTiles34(tiles)
	if err != nil {
		panic(err)
	}
	return tiles34
}

func MustStrToTile34(tile string) int {
	tile34, err := StrToTile34(tile)
	if err != nil {
		panic(err)
	}
	return tile34
}

// [0, 2, 9] => "13m 1p"
func TilesToMergedStr(tiles []int) (res string) {
	sort.Ints(tiles)
	merge := func(lowerIndex, upperIndex int, endsWith string) {
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
	merge(0, 9, "m ")
	merge(9, 18, "p ")
	merge(18, 27, "s ")
	merge(27, 34, "z")
	return strings.TrimSpace(res)
}

func Tile34ToMergedStr(tile34 int) (res string) {
	return Tiles34ToMergedStr([]int{tile34})
}

// [0, 2, 9] => "[13m 1p]"
func TilesToMergedStrWithBracket(tiles []int) (res string) {
	return "[" + TilesToMergedStr(tiles) + "]"
}

func Tiles34ToMergedStr(tiles34 []int) (res string) {
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

func Tiles34ToMergedStrWithBracket(tiles34 []int) (res string) {
	return "[" + Tiles34ToMergedStr(tiles34) + "]"
}
