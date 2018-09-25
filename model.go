package main

import "sort"

type needTiles map[int]int

func (t needTiles) parse() (allCount int, tiles []string) {
	if len(t) == 0 {
		return 0, nil
	}

	tileIndexes := make([]int, 0, len(t))
	for idx, cnt := range t {
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
