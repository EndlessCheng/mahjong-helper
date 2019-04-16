package util

import "fmt"

// 14 张牌的某种拆解结果
// 若该拆解没有刻子和顺子则为七对子/国士
type DivideResult struct {
	// 雀头牌
	pairTile int
	// 刻子牌
	KotsuTiles []int
	// 顺子牌的第一张（如 678s 的 6s）
	ShuntsuFirstTiles []int
}

// 调试用
func (d *DivideResult) String() string {
	humanTiles := []string{TilesToStr([]int{d.pairTile, d.pairTile})}
	for _, kotsuTile := range d.KotsuTiles {
		humanTiles = append(humanTiles, TilesToStr([]int{kotsuTile, kotsuTile, kotsuTile}))
	}
	for _, shuntsuFirstTile := range d.ShuntsuFirstTiles {
		humanTiles = append(humanTiles, TilesToStr([]int{shuntsuFirstTile, shuntsuFirstTile + 1, shuntsuFirstTile + 2}))
	}
	return fmt.Sprint(humanTiles)
}

// 判断是否为特殊牌型（七对子/国士）
func (d *DivideResult) IsSpecial() bool {
	return len(d.KotsuTiles) == 0 && len(d.ShuntsuFirstTiles) == 0
}

// 14张牌，返回所有可能的拆解，没有拆解表示未和牌
// http://hp.vector.co.jp/authors/VA046927/mjscore/mjalgorism.html
// http://hp.vector.co.jp/authors/VA046927/mjscore/AgariIndex.java
func DivideTiles34(tiles34 []int) (divideResults []DivideResult) {
	key := 0
	bitPos := -1

	pos34InHand14 := make([]int, 14)
	handPos := 0

	// 数牌
	for i := 0; i < 3; i++ {
		prevInHand := false // 上一张牌是否在手牌中
		for j := 0; j < 9; j++ {
			idx := i*9 + j
			if c := tiles34[idx]; c > 0 {
				prevInHand = true
				bitPos++
				pos34InHand14[handPos] = idx
				handPos++
				switch c {
				case 2:
					key |= 0x3 << uint(bitPos)
					bitPos += 2
				case 3:
					key |= 0xF << uint(bitPos)
					bitPos += 4
				case 4:
					key |= 0x3F << uint(bitPos)
					bitPos += 6
				}
			} else {
				if prevInHand {
					prevInHand = false
					key |= 0x1 << uint(bitPos)
					bitPos++
				}
			}
		}
		if prevInHand {
			key |= 0x1 << uint(bitPos)
			bitPos++
		}
	}

	// 字牌
	for i := 27; i < 34; i++ {
		if c := tiles34[i]; c > 0 {
			bitPos++
			pos34InHand14[handPos] = i
			handPos++
			switch c {
			case 2:
				key |= 0x3 << uint(bitPos)
				bitPos += 2
			case 3:
				key |= 0xF << uint(bitPos)
				bitPos += 4
			case 4:
				key |= 0x3F << uint(bitPos)
				bitPos += 6
			}
			key |= 0x1 << uint(bitPos)
			bitPos++
		}
	}

	results, ok := winTable[key]
	if !ok {
		return nil
	}

	for _, r := range results {
		// 雀头
		pairTile := pos34InHand14[(r>>6)&0xF]

		// 刻子
		numKotsu := r & 0x7
		kotsuTiles := make([]int, numKotsu)
		for i := range kotsuTiles {
			kotsuTiles[i] = pos34InHand14[(r>>uint(10+i*4))&0xF]
		}

		// 顺子的第一张牌
		numShuntsu := (r >> 3) & 0x7
		shuntsuFirstTiles := make([]int, numShuntsu)
		for i := range shuntsuFirstTiles {
			shuntsuFirstTiles[i] = pos34InHand14[(r>>uint(10+(numKotsu+i)*4))&0xF]
		}

		divideResults = append(divideResults, DivideResult{
			pairTile:          pairTile,
			KotsuTiles:        kotsuTiles,
			ShuntsuFirstTiles: shuntsuFirstTiles,
		})
	}

	return
}
