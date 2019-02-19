package main

import (
	"sort"
	"fmt"
)

// map[mahjong下标]数量
type needTiles map[int]int

func (nt needTiles) allCount() (count int) {
	for _, cnt := range nt {
		count += cnt
	}
	return count
}

func (nt needTiles) indexes() []int {
	if len(nt) == 0 {
		return nil
	}

	tileIndexes := make([]int, 0, len(nt))
	for idx := range nt {
		tileIndexes = append(tileIndexes, idx)
	}
	sort.Ints(tileIndexes)

	return tileIndexes
}

func (nt needTiles) parseIndex() (allCount int, indexes []int) {
	return nt.allCount(), nt.indexes()
}

func (nt needTiles) _parse(template [34]string) (allCount int, tiles []string) {
	if len(nt) == 0 {
		return 0, nil
	}

	tileIndexes := make([]int, 0, len(nt))
	for idx, cnt := range nt {
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

func (nt needTiles) parse() (allCount int, tiles []string) {
	return nt._parse(mahjong)
}

func (nt needTiles) parseZH() (allCount int, tilesZH []string) {
	return nt._parse(mahjongZH)
}

func (nt needTiles) tilesZH() []string {
	_, tiles := nt.parseZH()
	return tiles
}

func (nt needTiles) String() string {
	return fmt.Sprint(nt.tilesZH())
}

func (nt needTiles) containAllIndexes(anotherNeeds needTiles) bool {
	for k := range anotherNeeds {
		if _, ok := nt[k]; !ok {
			return false
		}
	}
	return true
}

// 是否包含字牌
func (nt needTiles) containHonors() bool {
	indexes := nt.indexes()
	if len(indexes) == 0 {
		return false
	}
	return indexes[len(indexes)-1] >= 27
}

func (nt needTiles) fixCountsWithLeftCounts(leftCounts []int) {
	for k := range nt {
		nt[k] = leftCounts[k]
	}
}
