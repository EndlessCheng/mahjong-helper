package util

import (
	"fmt"
)

const (
	shantenStateAgari  = -1
	shantenStateTenpai = 0
)

// 参考 http://ara.moo.jp/mjhmr/shanten.htm
// 七对子向听数 = 6-对子数+max(0,7-种类数)
func CalculateShantenOfChiitoi(tiles34 []int) int {
	shanten := 6
	numKind := 0
	for _, c := range tiles34 {
		if c == 0 {
			continue
		}
		if c >= 2 {
			shanten--
		}
		numKind++
	}
	shanten += MaxInt(0, 7-numKind)
	return shanten
}

type shanten struct {
	tiles         []int
	numberMelds   int
	numberTatsu   int
	numberPairs   int
	numberJidahai int // 13枚にしてから少なくとも打牌しなければならない字牌の数 -> これより向聴数は下がらない
	ankanTiles    int // 暗杠，28bit 位压缩：27bit数牌|1bit字牌
	isolatedTiles int // 孤张，28bit 位压缩：27bit数牌|1bit字牌
	minShanten    int
}

func (st *shanten) scanCharacterTiles(countOfTiles int) {
	ankanTiles := 0    // 暗杠，7bit 位压缩
	isolatedTiles := 0 // 孤张，7bit 位压缩

	for i, c := range st.tiles[27:] {
		if c == 0 {
			continue
		}
		switch c {
		case 1:
			isolatedTiles |= 1 << uint(i)
		case 2:
			st.numberPairs++
		case 3:
			st.numberMelds++
		case 4:
			st.numberMelds++
			st.numberJidahai++
			ankanTiles |= 1 << uint(i)
			isolatedTiles |= 1 << uint(i)
		}
	}

	if st.numberJidahai > 0 && countOfTiles%3 == 2 {
		st.numberJidahai--
	}

	if isolatedTiles > 0 {
		st.isolatedTiles |= 1 << 27
		if ankanTiles|isolatedTiles == ankanTiles {
			// 此孤张不能视作单骑做雀头的材料
			st.ankanTiles |= 1 << 27
		}
	}
}

// 计算一般型（非七对子和国士无双）的向听数
// 参考 http://ara.moo.jp/mjhmr/shanten.htm
func (st *shanten) calcNormalShanten() int {
	_shanten := 8 - 2*st.numberMelds - st.numberTatsu - st.numberPairs
	numMentsuKouho := st.numberMelds + st.numberTatsu
	if st.numberPairs > 0 {
		numMentsuKouho += st.numberPairs - 1 // 有雀头时面子候补-1
	} else if st.ankanTiles > 0 && st.isolatedTiles > 0 {
		if st.ankanTiles|st.isolatedTiles == st.ankanTiles { // 没有雀头，且除了暗杠外没有孤张，这连单骑都算不上
			// 比如 5555m 应该算作一向听
			_shanten++
		}
	}
	if numMentsuKouho > 4 { // 面子候补过多
		_shanten += numMentsuKouho - 4
	}
	if _shanten != shantenStateAgari && _shanten < st.numberJidahai {
		return st.numberJidahai
	}
	return _shanten
}

// 拆分出一个暗刻
func (st *shanten) increaseSet(k int) {
	st.tiles[k] -= 3
	st.numberMelds++
}

func (st *shanten) decreaseSet(k int) {
	st.tiles[k] += 3
	st.numberMelds--
}

// 拆分出一个雀头
func (st *shanten) increasePair(k int) {
	st.tiles[k] -= 2
	st.numberPairs++
}

func (st *shanten) decreasePair(k int) {
	st.tiles[k] += 2
	st.numberPairs--
}

// 拆分出一个顺子
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

// 拆分出一个两面/边张搭子
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

// 拆分出一个坎张搭子
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

// 拆分出一个孤张（浮牌）
func (st *shanten) increaseIsolatedTile(k int) {
	st.tiles[k]--
	st.isolatedTiles |= 1 << uint(k)
}
func (st *shanten) decreaseIsolatedTile(k int) {
	st.tiles[k]++
	st.isolatedTiles &^= 1 << uint(k)
}

func (st *shanten) run(depth int) {
	if st.minShanten == shantenStateAgari {
		return
	}

	// skip
	for ; depth < 27 && st.tiles[depth] == 0; depth++ {
	}

	if depth >= 27 {
		_shanten := st.calcNormalShanten()
		st.minShanten = MinInt(st.minShanten, _shanten)
		return
	}

	// i := depth % 9
	// 快速取模
	i := depth
	if i > 8 {
		i -= 9
	}
	if i > 8 {
		i -= 9
	}

	// 手牌拆解
	switch st.tiles[depth] {
	case 1:
		// 孤立牌は２つ以上取る必要は無い -> 雀头のほうが向聴数は下がる -> ３枚 -> 雀头＋孤立は雀头から取る
		// 孤立牌は合計８枚以上取る必要は無い
		if i < 6 && st.tiles[depth+1] == 1 && st.tiles[depth+2] > 0 && st.tiles[depth+3] < 4 {
			// 延べ単
			// 顺子
			st.increaseSyuntsu(depth)
			st.run(depth + 2)
			st.decreaseSyuntsu(depth)
		} else {
			// 浮牌
			st.increaseIsolatedTile(depth)
			st.run(depth + 1)
			st.decreaseIsolatedTile(depth)

			if i < 7 && st.tiles[depth+2] > 0 {
				if st.tiles[depth+1] != 0 {
					// 顺子
					st.increaseSyuntsu(depth)
					st.run(depth + 1)
					st.decreaseSyuntsu(depth)
				}
				// 坎张搭子
				st.increaseTatsuSecond(depth)
				st.run(depth + 1)
				st.decreaseTatsuSecond(depth)
			}
			if i < 8 && st.tiles[depth+1] > 0 {
				// 两面/边张搭子
				st.increaseTatsuFirst(depth)
				st.run(depth + 1)
				st.decreaseTatsuFirst(depth)
			}
		}
	case 2:
		// 雀头
		st.increasePair(depth)
		st.run(depth + 1)
		st.decreasePair(depth)

		if i < 7 && st.tiles[depth+1] > 0 && st.tiles[depth+2] > 0 {
			// 顺子
			st.increaseSyuntsu(depth)
			st.run(depth)
			st.decreaseSyuntsu(depth)
		}
	case 3:
		// 暗刻
		st.increaseSet(depth)
		st.run(depth + 1)
		st.decreaseSet(depth)

		st.increasePair(depth)
		if i < 7 && st.tiles[depth+1] > 0 && st.tiles[depth+2] > 0 {
			// 雀头+顺子
			st.increaseSyuntsu(depth)
			st.run(depth + 1)
			st.decreaseSyuntsu(depth)
		} else {
			if i < 7 && st.tiles[depth+2] > 0 {
				// 雀头+坎张搭子
				st.increaseTatsuSecond(depth)
				st.run(depth + 1)
				st.decreaseTatsuSecond(depth)
			}
			if i < 8 && st.tiles[depth+1] > 0 {
				// 雀头+两面/边张搭子
				st.increaseTatsuFirst(depth)
				st.run(depth + 1)
				st.decreaseTatsuFirst(depth)
			}
		}
		st.decreasePair(depth)

		if i < 7 && st.tiles[depth+1] >= 2 && st.tiles[depth+2] >= 2 {
			// 一杯口
			st.increaseSyuntsu(depth)
			st.increaseSyuntsu(depth)
			st.run(depth)
			st.decreaseSyuntsu(depth)
			st.decreaseSyuntsu(depth)
		}
	case 4:
		st.increaseSet(depth)
		if i < 7 && st.tiles[depth+2] > 0 {
			if st.tiles[depth+1] > 0 {
				// 暗刻+顺子
				st.increaseSyuntsu(depth)
				st.run(depth + 1)
				st.decreaseSyuntsu(depth)
			}
			// 暗刻+坎张搭子
			st.increaseTatsuSecond(depth)
			st.run(depth + 1)
			st.decreaseTatsuSecond(depth)
		}
		if i < 8 && st.tiles[depth+1] > 0 {
			// 暗刻+两面/边张搭子
			st.increaseTatsuFirst(depth)
			st.run(depth + 1)
			st.decreaseTatsuFirst(depth)
		}
		// 暗刻+孤张
		st.increaseIsolatedTile(depth)
		st.run(depth + 1)
		st.decreaseIsolatedTile(depth)
		st.decreaseSet(depth)

		st.increasePair(depth)
		if i < 7 && st.tiles[depth+2] > 0 {
			if st.tiles[depth+1] > 0 {
				// 雀头+顺子
				st.increaseSyuntsu(depth)
				st.run(depth)
				st.decreaseSyuntsu(depth)
			}
			// 雀头+坎张搭子
			st.increaseTatsuSecond(depth)
			st.run(depth + 1)
			st.decreaseTatsuSecond(depth)
		}
		if i < 8 && st.tiles[depth+1] > 0 {
			// 雀头+两面/边张搭子
			st.increaseTatsuFirst(depth)
			st.run(depth + 1)
			st.decreaseTatsuFirst(depth)
		}
		st.decreasePair(depth)
	}
}

// 根据手牌计算一般型（不考虑七对国士）的向听数
// 3k+1 和 3k+2 张牌都行
func CalculateShantenOfNormal(tiles34 []int, countOfTiles int) int {
	st := shanten{
		numberMelds: (14 - countOfTiles) / 3,
		minShanten:  8, // 不考虑国士无双和七对子的最大向听
		tiles:       tiles34,
	}

	st.scanCharacterTiles(countOfTiles)

	for i, c := range st.tiles[:27] {
		if c == 4 {
			st.ankanTiles |= 1 << uint(i)
		}
	}

	st.run(0)

	return st.minShanten
}

// 根据手牌计算向听数（不考虑国士）
// 3k+1 和 3k+2 张牌都行
func CalculateShanten(tiles34 []int) int {
	countOfTiles := CountOfTiles34(tiles34) // 若入参带 countOfTiles，能节省约 5% 的时间
	if countOfTiles > 14 {
		panic(fmt.Sprintln("[CalculateShanten] 参数错误 >14", tiles34, countOfTiles))
	}
	minShanten := CalculateShantenOfNormal(tiles34, countOfTiles)
	if countOfTiles >= 13 { // 考虑七对子
		minShanten = MinInt(minShanten, CalculateShantenOfChiitoi(tiles34))
	}
	return minShanten
}
