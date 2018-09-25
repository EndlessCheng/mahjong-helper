package main

import (
	"sort"
	"fmt"
)

// map[mahjong下标]数量
type needTiles map[int]int

func (nt needTiles) parse() (allCount int, tiles []string) {
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
		tiles[i] = mahjong[idx]
	}

	return allCount, tiles
}

func (nt needTiles) String() string {
	tiles := make([]string, 0, len(nt))
	for k := range nt {
		tiles = append(tiles, mahjong[k])
	}
	return fmt.Sprint(tiles)
}
