package util

import "fmt"

func _calcKey(tiles34 []int) (key int) {
	bitPos := -1

	// 数牌
	idx := -1
	for i := 0; i < 3; i++ {
		prevInHand := false // 上一张牌是否在手牌中
		for j := 0; j < 9; j++ {
			idx++
			if c := tiles34[idx]; c > 0 {
				prevInHand = true
				bitPos++
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

	return
}

// 3k+2 张牌，是否和牌（不检测国士无双）
func IsAgari(tiles34 []int) bool {
	key := _calcKey(tiles34)
	_, isAgari := winTable[key]
	return isAgari
}

//

// 3k+2 张牌的某种拆解结果
type DivideResult struct {
	PairTile          int   // 雀头牌
	KotsuTiles        []int // 刻子牌（注意 len(KotsuTiles) 为自摸时的暗刻数，荣和时的暗刻数需要另加逻辑判断）
	ShuntsuFirstTiles []int // 顺子牌的第一张（如 678s 的 6s）

	// 由于生成 winTable 的代码是不考虑具体是什么牌的，
	// 所以只能判断如七对子、九莲宝灯、一气通贯、两杯口、一杯口等和「形状」有关的役，
	// 像国士无双、断幺、全带、三色、绿一色等，和具体的牌/位置有关的役是判断不出的，需要另加逻辑判断
	IsChiitoi       bool // 七对子
	IsChuurenPoutou bool // 九莲宝灯
	IsIttsuu        bool // 一气通贯（注意：未考虑副露！）
	IsRyanpeikou    bool // 两杯口（IsRyanpeikou == true 时 IsIipeikou == false）
	IsIipeikou      bool // 一杯口
}

// 调试用
func (d *DivideResult) String() string {
	if d.IsChiitoi {
		return "[七对子]"
	}

	output := ""

	humanTilesList := []string{TilesToStr([]int{d.PairTile, d.PairTile})}
	for _, kotsuTile := range d.KotsuTiles {
		humanTilesList = append(humanTilesList, TilesToStr([]int{kotsuTile, kotsuTile, kotsuTile}))
	}
	for _, shuntsuFirstTile := range d.ShuntsuFirstTiles {
		humanTilesList = append(humanTilesList, TilesToStr([]int{shuntsuFirstTile, shuntsuFirstTile + 1, shuntsuFirstTile + 2}))
	}
	output += fmt.Sprint(humanTilesList)

	if d.IsChuurenPoutou {
		output += "[九莲宝灯]"
	}
	if d.IsIttsuu {
		output += "[一气通贯]"
	}
	if d.IsRyanpeikou {
		output += "[两杯口]"
	}
	if d.IsIipeikou {
		output += "[一杯口]"
	}

	return output
}

// 3k+2 张牌，返回所有可能的拆解，没有拆解表示未和牌（不检测国士无双）
// http://hp.vector.co.jp/authors/VA046927/mjscore/mjalgorism.html
// http://hp.vector.co.jp/authors/VA046927/mjscore/AgariIndex.java
func DivideTiles34(tiles34 []int) (divideResults []*DivideResult) {
	tiles14 := make([]int, 14)
	tiles14TailIndex := 0

	key := 0
	bitPos := -1

	// 数牌
	idx := -1
	for i := 0; i < 3; i++ {
		prevInHand := false // 上一张牌是否在手牌中
		for j := 0; j < 9; j++ {
			idx++
			if c := tiles34[idx]; c > 0 {
				tiles14[tiles14TailIndex] = idx
				tiles14TailIndex++

				prevInHand = true
				bitPos++
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
			tiles14[tiles14TailIndex] = i
			tiles14TailIndex++

			bitPos++
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
		return
	}

	// 3bit  0: 刻子数(0～4)
	// 3bit  3: 顺子数(0～4)
	// 4bit  6: 雀头位置(1～13)
	// 4bit 10: 面子位置1(0～13) 刻子在前，顺子在后
	// 4bit 14: 面子位置2(0～13)
	// 4bit 18: 面子位置3(0～13)
	// 4bit 22: 面子位置4(0～13)
	// 1bit 26: 七对子
	// 1bit 27: 九莲宝灯
	// 1bit 28: 一气通贯
	// 1bit 29: 两杯口
	// 1bit 30: 一杯口
	for _, r := range results {
		// 雀头
		pairTile := tiles14[(r>>6)&0xF]

		// 刻子
		numKotsu := r & 0x7
		kotsuTiles := make([]int, numKotsu)
		for i := range kotsuTiles {
			kotsuTiles[i] = tiles14[(r>>uint(10+i*4))&0xF]
		}

		// 顺子的第一张牌
		numShuntsu := (r >> 3) & 0x7
		shuntsuFirstTiles := make([]int, numShuntsu)
		for i := range shuntsuFirstTiles {
			shuntsuFirstTiles[i] = tiles14[(r>>uint(10+(numKotsu+i)*4))&0xF]
		}

		divideResults = append(divideResults, &DivideResult{
			PairTile:          pairTile,
			KotsuTiles:        kotsuTiles,
			ShuntsuFirstTiles: shuntsuFirstTiles,
			IsChiitoi:         r&(1<<26) != 0,
			IsChuurenPoutou:   r&(1<<27) != 0,
			IsIttsuu:          r&(1<<28) != 0,
			IsRyanpeikou:      r&(1<<29) != 0,
			IsIipeikou:        r&(1<<30) != 0,
		})
	}

	return
}
