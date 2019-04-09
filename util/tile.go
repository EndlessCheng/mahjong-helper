package util

import (
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
	return fmt.Sprintf("%d 进张 %s", w.AllCount(), TilesToStrWithBracket(w.indexes()))
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

func CountOfTiles(tiles []int) (count int) {
	for _, c := range tiles {
		count += c
	}
	return
}
