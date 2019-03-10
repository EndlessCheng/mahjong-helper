package util

import "sort"

const (
	WallSafeTypeNC = iota
	WallSafeTypeOC_Double
	WallSafeTypeOC_Mix
	WallSafeTypeOC_Single
)

type WallSafeTile struct {
	Tile34   int
	SafeType int
}

type WallSafeTileList []WallSafeTile

func (l WallSafeTileList) sort() {
	normalIndex := func(tile34 int) int {
		idx := tile34 % 9
		if idx >= 5 {
			// 5678 -> 3210
			idx = 8 - idx
		}
		return idx
	}

	sort.Slice(l, func(i, j int) bool {
		li, lj := l[i], l[j]

		liIndex := normalIndex(li.Tile34)
		ljIndex := normalIndex(lj.Tile34)
		// 先判断345
		if liIndex > 2 && ljIndex <= 2 || liIndex <= 2 && ljIndex > 2 {
			return liIndex < ljIndex
		}

		if li.SafeType != lj.SafeType {
			return li.SafeType < lj.SafeType
		}

		return liIndex < ljIndex
	})
}

func (l WallSafeTileList) FilterWithHands(handsTiles34 []int) WallSafeTileList {
	newSafeTiles34 := WallSafeTileList{}
	for _, safeTile := range l {
		if inInts(safeTile.Tile34, handsTiles34) {
			newSafeTiles34 = append(newSafeTiles34, safeTile)
		}
	}
	newSafeTiles34.sort()
	return newSafeTiles34
}

// 根据剩余牌 leftTiles34 中的某些牌是否为 0，来判断哪些牌较为安全（只输双碰、单骑、边张、坎张）
func CalcNCSafeTiles34(leftTiles34 []int) (ncSafeTiles34 WallSafeTileList) {
	const leftLimit = 0
	const safeType = WallSafeTypeNC
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			idx := 9*i + j
			if leftTiles34[idx+1] == leftLimit || leftTiles34[idx+2] == leftLimit {
				ncSafeTiles34 = append(ncSafeTiles34, WallSafeTile{idx, safeType})
			}
		}
		for j := 3; j < 6; j++ {
			idx := 9*i + j
			if (leftTiles34[idx-2] == leftLimit || leftTiles34[idx-1] == leftLimit) && (leftTiles34[idx+1] == leftLimit || leftTiles34[idx+2] == leftLimit) {
				ncSafeTiles34 = append(ncSafeTiles34, WallSafeTile{idx, safeType})
			}
		}
		for j := 6; j < 9; j++ {
			idx := 9*i + j
			if leftTiles34[idx-2] == leftLimit || leftTiles34[idx-1] == leftLimit {
				ncSafeTiles34 = append(ncSafeTiles34, WallSafeTile{idx, safeType})
			}
		}
	}
	ncSafeTiles34.sort()
	return
}

func CalcOCSafeTiles34(leftTiles34 []int) (ocSafeTiles34 WallSafeTileList) {
	const leftLimit = 1
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			idx := 9*i + j
			if leftTiles34[idx+1] == leftLimit && leftTiles34[idx+2] == leftLimit {
				ocSafeTiles34 = append(ocSafeTiles34, WallSafeTile{idx, WallSafeTypeOC_Double})
			} else if leftTiles34[idx+1] == leftLimit || leftTiles34[idx+2] == leftLimit {
				ocSafeTiles34 = append(ocSafeTiles34, WallSafeTile{idx, WallSafeTypeOC_Single})
			}
		}
		for j := 3; j < 6; j++ {
			idx := 9*i + j
			if (leftTiles34[idx-2] == leftLimit || leftTiles34[idx-1] == leftLimit) && (leftTiles34[idx+1] == leftLimit || leftTiles34[idx+2] == leftLimit) {
				if leftTiles34[idx-2] == leftLimit && leftTiles34[idx-1] == leftLimit && leftTiles34[idx+1] == leftLimit && leftTiles34[idx+2] == leftLimit {
					// 两边都是 double
					ocSafeTiles34 = append(ocSafeTiles34, WallSafeTile{idx, WallSafeTypeOC_Double})
				} else if leftTiles34[idx-2] == leftLimit && leftTiles34[idx-1] == leftLimit || leftTiles34[idx+1] == leftLimit && leftTiles34[idx+2] == leftLimit {
					// 一半 double，一半不是 double
					ocSafeTiles34 = append(ocSafeTiles34, WallSafeTile{idx, WallSafeTypeOC_Mix})
				} else {
					ocSafeTiles34 = append(ocSafeTiles34, WallSafeTile{idx, WallSafeTypeOC_Single})
				}
			}
		}
		for j := 6; j < 9; j++ {
			idx := 9*i + j
			if leftTiles34[idx-2] == leftLimit && leftTiles34[idx-1] == leftLimit {
				ocSafeTiles34 = append(ocSafeTiles34, WallSafeTile{idx, WallSafeTypeOC_Double})
			} else if leftTiles34[idx-2] == leftLimit || leftTiles34[idx-1] == leftLimit {
				ocSafeTiles34 = append(ocSafeTiles34, WallSafeTile{idx, WallSafeTypeOC_Single})
			}
		}
	}
	ocSafeTiles34.sort()
	return
}

func CalcWallTiles34(leftTiles34 []int) (safeTiles34 WallSafeTileList) {
	safeTiles34 = append(safeTiles34, CalcNCSafeTiles34(leftTiles34)...)
	safeTiles34 = append(safeTiles34, CalcOCSafeTiles34(leftTiles34)...)
	safeTiles34.sort()
	return
}
