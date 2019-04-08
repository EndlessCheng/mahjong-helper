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
		if handsTiles34[safeTile.Tile34] > 0 {
			newSafeTiles34 = append(newSafeTiles34, safeTile)
		}
	}
	newSafeTiles34.sort()
	return newSafeTiles34
}

// 根据剩余牌 leftTiles34 中的某些牌是否为 0，来判断哪些牌非常安全（Double No Chance：只输单骑、双碰）
func CalcDNCSafeTiles(leftTiles34 []int) (dncSafeTiles []int) {
	const leftLimit = 0
	nc := func(idx int) bool {
		return leftTiles34[idx] == leftLimit
	}
	or := func(idx ...int) bool {
		for _, i := range idx {
			if nc(i) {
				return true
			}
		}
		return false
	}
	and := func(idx ...int) bool {
		for _, i := range idx {
			if !nc(i) {
				return false
			}
		}
		return true
	}

	for i := 0; i < 3; i++ {
		// 2/3断的1
		if or(9*i+1, 9*i+2) {
			dncSafeTiles = append(dncSafeTiles, 9*i)
		}
		// 3/14断的2
		if nc(9*i+2) || and(9*i, 9*i+3) {
			dncSafeTiles = append(dncSafeTiles, 9*i+1)
		}
		// 14/24/25断的3（4567同理）
		for j := 2; j <= 6; j++ {
			idx := 9*i + j
			if and(idx-2, idx+1) || and(idx-1, idx+1) || and(idx-1, idx+2) {
				dncSafeTiles = append(dncSafeTiles, idx)
			}
		}
		// 7/69断的8
		if nc(9*i+6) || and(9*i+5, 9*i+8) {
			dncSafeTiles = append(dncSafeTiles, 9*i+7)
		}
		// 7/8断的9
		if or(9*i+6, 9*i+7) {
			dncSafeTiles = append(dncSafeTiles, 9*i+8)
		}
	}
	return
}

// 根据剩余牌 leftTiles34 中的某些牌是否为 0，来判断哪些牌较为安全（No Chance：只输单骑、双碰、边张、坎张）
func CalcNCSafeTiles(leftTiles34 []int) (ncSafeTiles WallSafeTileList) {
	const leftLimit = 0
	nc := func(idx int) bool {
		return leftTiles34[idx] == leftLimit
	}
	or := func(idx ...int) bool {
		for _, i := range idx {
			if nc(i) {
				return true
			}
		}
		return false
	}

	const safeType = WallSafeTypeNC
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			idx := 9*i + j
			if or(idx+1, idx+2) {
				ncSafeTiles = append(ncSafeTiles, WallSafeTile{idx, safeType})
			}
		}
		for j := 3; j < 6; j++ {
			idx := 9*i + j
			if or(idx-2, idx-1) && or(idx+1, idx+2) {
				ncSafeTiles = append(ncSafeTiles, WallSafeTile{idx, safeType})
			}
		}
		for j := 6; j < 9; j++ {
			idx := 9*i + j
			if or(idx-2, idx-1) {
				ncSafeTiles = append(ncSafeTiles, WallSafeTile{idx, safeType})
			}
		}
	}
	ncSafeTiles.sort()
	return
}

// 根据剩余牌 leftTiles34 中的某些牌是否为 1，来判断哪些牌较为安全（One Chance：早巡大概率只输单骑、双碰、边张、坎张）
func CalcOCSafeTiles(leftTiles34 []int) (ocSafeTiles WallSafeTileList) {
	const leftLimit = 1
	nc := func(idx int) bool {
		return leftTiles34[idx] == leftLimit
	}
	or := func(idx ...int) bool {
		for _, i := range idx {
			if nc(i) {
				return true
			}
		}
		return false
	}
	and := func(idx ...int) bool {
		for _, i := range idx {
			if !nc(i) {
				return false
			}
		}
		return true
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			idx := 9*i + j
			if and(idx+1, idx+2) {
				ocSafeTiles = append(ocSafeTiles, WallSafeTile{idx, WallSafeTypeOC_Double})
			} else if or(idx+1, idx+2) {
				ocSafeTiles = append(ocSafeTiles, WallSafeTile{idx, WallSafeTypeOC_Single})
			}
		}
		for j := 3; j < 6; j++ {
			idx := 9*i + j
			if or(idx-2, idx-1) && or(idx+1, idx+2) {
				if and(idx-2, idx-1, idx+1, idx+2) {
					// 两边都是 double
					ocSafeTiles = append(ocSafeTiles, WallSafeTile{idx, WallSafeTypeOC_Double})
				} else if and(idx-2, idx-1) || and(idx+1, idx+2) {
					// 一半 double，一半不是 double
					ocSafeTiles = append(ocSafeTiles, WallSafeTile{idx, WallSafeTypeOC_Mix})
				} else {
					ocSafeTiles = append(ocSafeTiles, WallSafeTile{idx, WallSafeTypeOC_Single})
				}
			}
		}
		for j := 6; j < 9; j++ {
			idx := 9*i + j
			if and(idx-2, idx-1) {
				ocSafeTiles = append(ocSafeTiles, WallSafeTile{idx, WallSafeTypeOC_Double})
			} else if or(idx-2, idx-1) {
				ocSafeTiles = append(ocSafeTiles, WallSafeTile{idx, WallSafeTypeOC_Single})
			}
		}
	}
	ocSafeTiles.sort()
	return
}

func CalcWallTiles(leftTiles34 []int) (safeTiles WallSafeTileList) {
	safeTiles = append(safeTiles, CalcNCSafeTiles(leftTiles34)...)
	safeTiles = append(safeTiles, CalcOCSafeTiles(leftTiles34)...)
	safeTiles.sort()
	return
}
