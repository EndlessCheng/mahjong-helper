package util

import (
	"fmt"
)

const agariState = -1

// 根据手牌计算向听数
// 3k+1 和 3k+2 张牌都行
func CalculateShanten(tiles34 []int) int {
	countOfTiles := CountOfTiles34(tiles34)
	if countOfTiles > 14 {
		panic(fmt.Sprintln("[CalculateShanten] 参数错误 >14", tiles34, countOfTiles))
	}

	_tiles34 := make([]int, 34)
	copy(_tiles34, tiles34)
	st := shanten{
		minShanten: 8, // 不考虑国士无双和七对子的最大向听
		tiles:      _tiles34,
	}

	if countOfTiles >= 13 {
		// 考虑七对子
		st.minShanten = st.scanChitoitsu()
	}

	st.removeCharacterTiles(countOfTiles)

	initMentsu := (14 - countOfTiles) / 3

	// TODO: 检查原理  对比 JS 检查正确性？
	st.scan(initMentsu)

	return st.minShanten
}

type shanten struct {
	numberMelds         int
	numberTatsu         int
	numberPairs         int
	numberJidahai       int
	numberCharacters    int
	numberIsolatedTiles int
	minShanten          int
	tiles               []int
}

func (st *shanten) removeCharacterTiles(countOfTiles int) {
	number := 0
	isolated := 0

	for i := 27; i < 34; i++ {
		if st.tiles[i] == 4 {
			st.numberMelds++
			st.numberJidahai++
			number |= 1 << uint(i-27)
			isolated |= 1 << uint(i-27)
		}

		if st.tiles[i] == 3 {
			st.numberMelds++
		}

		if st.tiles[i] == 2 {
			st.numberPairs++
		}

		if st.tiles[i] == 1 {
			isolated |= 1 << uint(i-27)
		}
	}

	if st.numberJidahai != 0 && (countOfTiles%3) == 2 {
		st.numberJidahai--
	}

	if isolated != 0 {
		st.numberIsolatedTiles |= 1 << 27
		if (number | isolated) == number {
			st.numberCharacters |= 1 << 27
		}
	}
}

func (st *shanten) scanChitoitsu() int {
	shanten := st.minShanten

	indices := []int{0, 8, 9, 17, 18, 26, 27, 28, 29, 30, 31, 32, 33}

	completedTerminals := 0
	for _, i := range indices {
		if st.tiles[i] >= 2 {
			completedTerminals++
		}
	}

	terminals := 0
	for _, i := range indices {
		if st.tiles[i] != 0 {
			terminals++
		}
	}

	indices = []int{1, 2, 3, 4, 5, 6, 7, 10, 11, 12, 13, 14, 15, 16, 19, 20, 21, 22, 23, 24, 25}

	completedPairs := completedTerminals
	for _, i := range indices {
		if st.tiles[i] >= 2 {
			completedPairs++
		}
	}

	pairs := terminals
	for _, i := range indices {
		if st.tiles[i] != 0 {
			pairs++
		}
	}

	retShanten := 6 - completedPairs
	if pairs < 7 && 7-pairs != 0 {
		retShanten++
	}
	if retShanten < shanten {
		shanten = retShanten
	}

	//retShanten = 13 - terminals
	//if completedTerminals != 0 {
	//	retShanten--
	//}
	//if retShanten < shanten {
	//	shanten = retShanten
	//}

	return shanten
}

func (st *shanten) scanChitoitsuAndKokushi() int {
	shanten := st.minShanten

	indices := []int{0, 8, 9, 17, 18, 26, 27, 28, 29, 30, 31, 32, 33}

	completedTerminals := 0
	for _, i := range indices {
		if st.tiles[i] >= 2 {
			completedTerminals++
		}
	}

	terminals := 0
	for _, i := range indices {
		if st.tiles[i] != 0 {
			terminals++
		}
	}

	indices = []int{1, 2, 3, 4, 5, 6, 7, 10, 11, 12, 13, 14, 15, 16, 19, 20, 21, 22, 23, 24, 25}

	completedPairs := completedTerminals
	for _, i := range indices {
		if st.tiles[i] >= 2 {
			completedPairs++
		}
	}

	pairs := terminals
	for _, i := range indices {
		if st.tiles[i] != 0 {
			pairs++
		}
	}

	retShanten := 6 - completedPairs
	if pairs < 7 && 7-pairs != 0 {
		retShanten++
	}
	if retShanten < shanten {
		shanten = retShanten
	}

	retShanten = 13 - terminals
	if completedTerminals != 0 {
		retShanten--
	}
	if retShanten < shanten {
		shanten = retShanten
	}

	return shanten
}

func (st *shanten) scan(initMentsu int) {
	st.numberCharacters = 0
	for i := 0; i < 27; i++ {
		st.numberCharacters |= boolToInt(st.tiles[i] == 4) << uint(i)
	}
	st.numberMelds += initMentsu
	st.run(0)
}

func (st *shanten) run(depth int) {
	if st.minShanten == agariState {
		return
	}

	for st.tiles[depth] == 0 {
		depth++

		if depth >= 27 {
			break
		}
	}

	if depth >= 27 {
		st.updateResult()
		return
	}

	i := depth
	if i > 8 {
		i -= 9
	}
	if i > 8 {
		i -= 9
	}

	if st.tiles[depth] == 4 {
		st.increaseSet(depth)
		if i < 7 && st.tiles[depth+2] != 0 {
			if st.tiles[depth+1] != 0 {
				st.increaseSyuntsu(depth)
				st.run(depth + 1)
				st.decreaseSyuntsu(depth)
			}
			st.increaseTatsuSecond(depth)
			st.run(depth + 1)
			st.decreaseTatsuSecond(depth)
		}

		if i < 8 && st.tiles[depth+1] != 0 {
			st.increaseTatsuFirst(depth)
			st.run(depth + 1)
			st.decreaseTatsuFirst(depth)
		}

		st.increaseIsolatedTile(depth)
		st.run(depth + 1)
		st.decreaseIsolatedTile(depth)
		st.decreaseSet(depth)
		st.increasePair(depth)

		if i < 7 && st.tiles[depth+2] != 0 {
			if st.tiles[depth+1] != 0 {
				st.increaseSyuntsu(depth)
				st.run(depth)
				st.decreaseSyuntsu(depth)
			}
			st.increaseTatsuSecond(depth)
			st.run(depth + 1)
			st.decreaseTatsuSecond(depth)
		}

		if i < 8 && st.tiles[depth+1] != 0 {
			st.increaseTatsuFirst(depth)
			st.run(depth + 1)
			st.decreaseTatsuFirst(depth)
		}

		st.decreasePair(depth)
	}

	if st.tiles[depth] == 3 {
		st.increaseSet(depth)
		st.run(depth + 1)
		st.decreaseSet(depth)
		st.increasePair(depth)

		if i < 7 && st.tiles[depth+1] != 0 && st.tiles[depth+2] != 0 {
			st.increaseSyuntsu(depth)
			st.run(depth + 1)
			st.decreaseSyuntsu(depth)
		} else {
			if i < 7 && st.tiles[depth+2] != 0 {
				st.increaseTatsuSecond(depth)
				st.run(depth + 1)
				st.decreaseTatsuSecond(depth)
			}

			if i < 8 && st.tiles[depth+1] != 0 {
				st.increaseTatsuFirst(depth)
				st.run(depth + 1)
				st.decreaseTatsuFirst(depth)
			}
		}

		st.decreasePair(depth)

		if i < 7 && st.tiles[depth+2] >= 2 && st.tiles[depth+1] >= 2 {
			st.increaseSyuntsu(depth)
			st.increaseSyuntsu(depth)
			st.run(depth)
			st.decreaseSyuntsu(depth)
			st.decreaseSyuntsu(depth)
		}

	}

	if st.tiles[depth] == 2 {
		st.increasePair(depth)
		st.run(depth + 1)
		st.decreasePair(depth)
		if i < 7 && st.tiles[depth+2] != 0 && st.tiles[depth+1] != 0 {
			st.increaseSyuntsu(depth)
			st.run(depth)
			st.decreaseSyuntsu(depth)
		}
	}

	if st.tiles[depth] == 1 {
		if i < 6 && st.tiles[depth+1] == 1 && st.tiles[depth+2] != 0 && st.tiles[depth+3] != 4 {
			st.increaseSyuntsu(depth)
			st.run(depth + 2)
			st.decreaseSyuntsu(depth)
		} else {
			st.increaseIsolatedTile(depth)
			st.run(depth + 1)
			st.decreaseIsolatedTile(depth)

			if i < 7 && st.tiles[depth+2] != 0 {
				if st.tiles[depth+1] != 0 {
					st.increaseSyuntsu(depth)
					st.run(depth + 1)
					st.decreaseSyuntsu(depth)
				}
				st.increaseTatsuSecond(depth)
				st.run(depth + 1)
				st.decreaseTatsuSecond(depth)
			}

			if i < 8 && st.tiles[depth+1] != 0 {
				st.increaseTatsuFirst(depth)
				st.run(depth + 1)
				st.decreaseTatsuFirst(depth)
			}
		}
	}
}

func (st *shanten) updateResult() {
	retShanten := 8 - st.numberMelds*2 - st.numberTatsu - st.numberPairs
	nMentsuKouho := st.numberMelds + st.numberTatsu
	if st.numberPairs != 0 {
		nMentsuKouho += st.numberPairs - 1
	} else if st.numberCharacters != 0 && st.numberIsolatedTiles != 0 {
		if (st.numberCharacters | st.numberIsolatedTiles) == st.numberCharacters {
			retShanten++
		}
	}

	if nMentsuKouho > 4 {
		retShanten += nMentsuKouho - 4
	}

	if retShanten != agariState && retShanten < st.numberJidahai {
		retShanten = st.numberJidahai
	}

	if retShanten < st.minShanten {
		st.minShanten = retShanten
	}
}

func (st *shanten) increaseSet(k int) {
	st.tiles[k] -= 3
	st.numberMelds++
}

func (st *shanten) decreaseSet(k int) {
	st.tiles[k] += 3
	st.numberMelds--
}

func (st *shanten) increasePair(k int) {
	st.tiles[k] -= 2
	st.numberPairs++
}

func (st *shanten) decreasePair(k int) {
	st.tiles[k] += 2
	st.numberPairs--
}

func (st *shanten) increaseSyuntsu(k int) {
	st.tiles[k]--
	st.tiles[k+1]--
	st.tiles[k+2]--
	st.numberMelds++
}

func (st *shanten) decreaseSyuntsu(k int) {
	st.tiles[k]++
	st.tiles[k+1]++
	st.tiles[k+2]++
	st.numberMelds--
}

func (st *shanten) increaseTatsuFirst(k int) {
	st.tiles[k]--
	st.tiles[k+1]--
	st.numberTatsu++
}

func (st *shanten) decreaseTatsuFirst(k int) {
	st.tiles[k]++
	st.tiles[k+1]++
	st.numberTatsu--
}

func (st *shanten) increaseTatsuSecond(k int) {
	st.tiles[k]--
	st.tiles[k+2]--
	st.numberTatsu++
}

func (st *shanten) decreaseTatsuSecond(k int) {
	st.tiles[k]++
	st.tiles[k+2]++
	st.numberTatsu--
}

func (st *shanten) increaseIsolatedTile(k int) {
	st.tiles[k]--
	st.numberIsolatedTiles |= 1 << uint(k)
}

func (st *shanten) decreaseIsolatedTile(k int) {
	st.tiles[k]++
	st.numberIsolatedTiles &= ^(1 << uint(k))
}
