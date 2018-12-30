package main

import (
	"strings"
	"fmt"
	"os"
	"time"
	"github.com/fatih/color"
)

var detailFlag = false
var interactFlag = false // 交互模式

// 13 张牌，检查是否听牌，返回听牌的种类和枚数
// （不考虑国士无双）
func checkTing0(counts []int) needTiles {
	needs := needTiles{}

	// 剪枝：检测浮牌
	// 此处优化提升了 3-5 倍的性能
	for i := 0; i < 3; i++ {
		cnt0 := 2
		for j := 0; j < 9; j++ {
			idx := 9*i + j
			c := counts[idx]
			switch {
			case c == 0:
				cnt0++
			case c == 1:
				if cnt0 < 2 {
					cnt0 = 0
				} else {
					if j+1 < 9 && counts[idx+1] > 0 {
						j++
						cnt0 = 0
					} else if j+2 < 9 && counts[idx+2] > 0 {
						j += 2
						cnt0 = 0
					} else {
						// 这是一张浮牌，要想和牌只能单骑这一张
						counts[idx]++ // 摸牌
						if checkWin(counts) { // 单骑和牌
							needs[idx] = 4 - (counts[idx] - 1)
						}
						counts[idx]--
						return needs
					}
				}
			case c >= 2:
				cnt0 = 0
			}
		}
	}
	for i := 27; i < len(mahjong); i++ {
		if counts[i] == 1 {
			// 这是一张浮牌，要想和牌只能单骑这一张
			counts[i]++ // 摸牌
			if checkWin(counts) { // 单骑和牌
				needs[i] = 4 - (counts[i] - 1)
			}
			counts[i]--
			return needs
		}
	}

	// 剪枝：下面这段可以加快 35%~40%
	// 只计算可能胡的牌，对于 345m ... 这样的牌来说，摸到 1m 是肯定不能胡牌的，这种情况就不用去调用 checkWin 了
	needChecks := make([]int, 0, len(mahjong))
	idx := -1
	for i := 0; i < 3; i++ {
		for j := 0; j < 9; j++ {
			idx++
			if counts[idx] == 4 {
				continue
			}
			if j == 0 || j == 9 {
				if counts[idx] > 0 || counts[idx+1] > 0 {
					needChecks = append(needChecks, idx)
				}
			} else {
				if counts[idx-1] > 0 || counts[idx] > 0 || counts[idx+1] > 0 {
					needChecks = append(needChecks, idx)
				}
			}
		}
	}
	for i := 27; i < 34; i++ {
		if counts[i] == 4 {
			continue
		}
		if counts[i] > 0 {
			needChecks = append(needChecks, i)
		}
	}
	for _, checkIndex := range needChecks {
		counts[checkIndex]++ // 摸牌
		if checkWin(counts) { // 和牌
			needs[checkIndex] = 4 - (counts[checkIndex] - 1)
		}
		counts[checkIndex]--
	}

	//for i := range mahjong {
	//	if counts[i] == 4 {
	//		continue
	//	}
	//	counts[i]++ // 摸牌
	//	if checkWin(counts) { // 和牌
	//		needs[i] = 4 - (counts[i] - 1)
	//	}
	//	counts[i]--
	//}

	return needs
}

// 13 张牌，计算默听时的改良情况
func checkTing0Improve(counts []int, tings needTiles) ting0ImproveList {
	ting0ImproveList := ting0ImproveList{}
	for i := range mahjong {
		if counts[i] == 4 {
			continue
		}
		if _, ok := tings[i]; ok {
			continue
		}
		counts[i]++ // 摸牌

		for j := range mahjong {
			if counts[j] == 0 || j == i {
				continue
			}
			counts[j]-- // 切牌
			if needs := checkTing0(counts); len(needs) > 0 && !tings.containAllIndexes(needs) {
				ting0ImproveList = append(ting0ImproveList, ting0Improve{i, j, needs})
			}
			counts[j]++
		}

		counts[i]--
	}
	return ting0ImproveList
}

// 14 张牌，检查切掉某张牌后是否听牌
func checkTing0Discard(counts []int) ting0DiscardList {
	ting0DiscardList := ting0DiscardList{}
	for i := range mahjong {
		if counts[i] >= 1 {
			counts[i]-- // 切牌
			if needs := checkTing0(counts); len(needs) > 0 {
				// 切掉这张后的默听改良率
				ting0ImproveList := checkTing0Improve(counts, needs)
				goodTiles := ting0ImproveList.calcGoodImprove(counts)
				ting0DiscardList = append(ting0DiscardList, ting0Discard{i, needs, goodTiles})
			}
			counts[i]++
		}
	}
	return ting0DiscardList
}

//var (
//	detailBuffer = strings.Builder{}
//)

//func flushBuffer() {
//	if detailFlag {
//		fmt.Print(detailBuffer.String())
//	}
//	detailBuffer.Reset()
//}

// 13张牌，检查一向听
func checkTing1(counts []int, recur bool) (needTiles, *ting1Detail) {
	needs := needTiles{}
	betterNeedsMap := map[int]map[int]needTiles{}
	tingCntMap := map[int]int{} // map[摸到idx]听多少张牌

	for i := range mahjong {
		if counts[i] >= 1 {
			tmpNeedsMap := map[int]needTiles{}
			counts[i]-- // 切掉其中一张牌
			for j := range mahjong {
				if j == i {
					continue
				}
				if counts[j] == 4 {
					continue
				}
				counts[j]++ // 换成其他牌
				if nd := checkTing0(counts); len(nd) > 0 {
					// 若能听牌，则换的这张牌为一向听的进张
					if _, ok := needs[j]; !ok {
						needs[j] = 4 - (counts[j] - 1)
					} else {
						// 比如说 57m22566s，切 5s/6s 来 8m 都听牌，但是听牌的数量有区别
					}
					if recur {
						if tingCnt := nd.allCount(); tingCnt > tingCntMap[j] {
							// 听牌一般听数量最多的
							tingCntMap[j] = tingCnt
						}
					}
				} else if recur {
					if betterNeeds, _ := checkTing1(counts, false); len(betterNeeds) > 0 {
						// 换成这张牌也是一向听，可能是改良型，记录一下
						tmpNeedsMap[j] = betterNeeds
					}
				}
				counts[j]--
			}
			counts[i]++
			betterNeedsMap[i] = tmpNeedsMap
		}
	}

	if !recur {
		return needs, nil
	}

	ting1Detail := ting1Detail{}
	//detailBuffer.Reset()

	// TODO: 振听?
	if allCount, tiles := needs.parse(); allCount > 0 {
		improveCount := make([]int, len(mahjong))
		for i := range mahjong {
			improveCount[i] = allCount
		}
		for _, tmpNeedsMap := range betterNeedsMap {
			for drawIdx, betterNeeds := range tmpNeedsMap {
				if inStrSlice(mahjong[drawIdx], tiles) {
					// 跳过改良牌就是一向听的进张的情况
					continue
				}
				if betterAllCount := betterNeeds.allCount(); betterAllCount > allCount {
					// 进张数变多，则为一向听的改良
					ting1Detail.improveWayCount++
					if betterAllCount > improveCount[drawIdx] {
						improveCount[drawIdx] = betterAllCount
					}
					//detailBuffer.WriteString(fmt.Sprintln(fmt.Sprintf("    摸 %s 切 %s 改良:", mahjongZH[drawIdx], mahjongZH[discardIdx]), betterAllCount, betterTiles))
				}
			}
		}

		if ting1Detail.improveWayCount > 0 {
			improveScore := 0
			weight := 0
			for i := range mahjong {
				w := 4 - counts[i]
				improveScore += w * improveCount[i]
				weight += w
			}
			ting1Detail.avgImproveNum = float64(improveScore) / float64(weight)
		}

		avgTingSum := 0
		weight := 0
		for idx, c := range tingCntMap {
			w := 4 - counts[idx]
			avgTingSum += w * c
			weight += w
		}
		// TODO: 根据1-9的牌来计算综合和牌率
		ting1Detail.avgTingCount = float64(avgTingSum) / float64(weight)
	}

	return needs, &ting1Detail
}

// 14张牌，可以一向听，何切
// 检查能一向听的切牌，对比：
// 1. 进张数
// 2. 改良之后的（加权）平均进张数
// 3. 听牌后的（加权）平均听牌数
// 4. 听牌后所听牌的名称（就是一向听的进张名称）（一般来说 14m 优于 25m。不过还是要根据场况来判断）
// // TODO: 赤牌改良提醒！！
func checkTing1Discard(counts []int) ting1DiscardList {
	ting1DiscardList := ting1DiscardList{}
	for i := range mahjong {
		if counts[i] >= 1 {
			counts[i]-- // 切牌
			needs, ting1Detail := checkTing1(counts, true)
			if len(needs) > 0 {
				ting1DiscardList = append(ting1DiscardList, ting1Discard{i, needs, ting1Detail})
			}
			counts[i]++
		}
	}
	return ting1DiscardList
}

// 13张牌，检查一向听（简化版）
func _simpleCheckTing1(counts []int) needTiles {
	needs := needTiles{}
	for i := range mahjong {
		if counts[i] >= 1 {
			counts[i]-- // 切掉其中一张牌
			for j := range mahjong {
				if j == i {
					continue
				}
				if counts[j] == 4 {
					continue
				}
				counts[j]++ // 换成其他牌
				if nd := checkTing0(counts); len(nd) > 0 {
					// 若能听牌，则换的这张牌为一向听的进张
					if _, ok := needs[j]; !ok {
						needs[j] = 4 - (counts[j] - 1)
					} else {
						// 比如说 57m22566s，切 5s/6s 来 8m 都听牌
					}
				}
				counts[j]--
			}
			counts[i]++
		}
	}
	return needs
}

// 13张牌，检查两向听
// TODO: 两向听时计算一向听的平均进张
func checkTing2(counts []int) needTiles {
	needs := needTiles{}
	for i := range mahjong {
		if counts[i] >= 1 {
			counts[i]-- // 切掉其中一张牌
			for j := range mahjong {
				if j == i {
					continue
				}
				if counts[j] == 4 {
					continue
				}
				counts[j]++ // 换成其他牌
				if nd := _simpleCheckTing1(counts); len(nd) > 0 {
					// 若能一向听，则换的这张牌为两向听的进张
					if _, ok := needs[j]; !ok {
						needs[j] = 4 - (counts[j] - 1)
					}
				}
				counts[j]--
			}
			counts[i]++
		}
	}
	return needs
}

// 14张牌，可以两向听，何切
func checkTing2Discard(counts []int) ting2DiscardList {
	ting2DiscardList := ting2DiscardList{}
	for i := range mahjong {
		if counts[i] >= 1 {
			counts[i]-- // 切牌
			if needs := checkTing2(counts); len(needs) > 0 {
				ting2DiscardList = append(ting2DiscardList, ting2Discard{i, needs})
			}
			counts[i]++
		}
	}
	return ting2DiscardList
}

func analysis(raw string) (num int, counts []int, err error) {
	fmt.Println(raw)
	fmt.Println(strings.Repeat("=", len(raw)))

	num, counts, err = convert(raw)
	if err != nil {
		return
	}

	if countPairs(counts) >= 4 {
		color.Yellow("对子手可能")
		fmt.Println()
	}

	switch num {
	case 13:
		if needs := checkTing0(counts); len(needs) > 0 {
			count, tiles := needs.parseZH()
			fmt.Printf("已听牌: %v, %d 张", tiles, count)
			improve := checkTing0Improve(counts, needs)
			goodTiles := improve.calcGoodImprove(counts)
			improveCount := goodTiles.allCount()
			fmt.Printf(" (%d 默听改良, %.2f 改良率)", improveCount, float64(improveCount)/float64(count))
			fmt.Println()
			improve.print()
			break
		}

		if needs, ting1Detail := checkTing1(counts, true); len(needs) > 0 {
			count, tiles := needs.parseZH()
			fmt.Println("一向听:", count, tiles)
			ting1Detail.print()
			//flushBuffer()
			break
		}

		if needs := checkTing2(counts); len(needs) > 0 {
			count, tiles := needs.parseZH()
			fmt.Println("两向听:", count, tiles)

			setTing2MinCount(count)
			break
		}

		fmt.Println("尚未两向听")
	case 14:
		defer resetTing2MinCount()

		if checkWin(counts) {
			fmt.Println("已胡牌")
			break
		}

		if ting0DiscardList := checkTing0Discard(counts); len(ting0DiscardList) > 0 {
			color.New(color.FgRed).Print("【已听牌！】")
			fmt.Println()
			ting0DiscardList.print()

			// 这里不 break，保留倒退回一向听的选择
		}

		if ting1DiscardList := checkTing1Discard(counts); len(ting1DiscardList) > 0 {
			ting1DiscardList.print()
			// TODO: 倒退回两向听的选择？
			break
		}

		if ting2DiscardList := checkTing2Discard(counts); len(ting2DiscardList) > 0 {
			ting2DiscardList.print()
			break
		}

		fmt.Println("尚未两向听")
	default:
		err = fmt.Errorf("参数错误: %s（%d 张牌）", raw, num)
		return
	}

	//fmt.Println("checkWin", checkWinCount)
	fmt.Println()

	return
}

func interact(raw string) {
	num, counts, err := analysis(raw)
	if err != nil {
		_errorExit(err)
	}
	printed := true

	var tile string
	for {
		for {
			if num < 14 {
				num = 999
				break
			}
			printed = false
			fmt.Print("> 切 ")
			fmt.Scanf("%s\n", &tile)
			idx, err := _convert(tile)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			} else {
				if counts[idx] == 0 {
					fmt.Fprintln(os.Stderr, "切掉的牌不存在")
				} else {
					counts[idx]--
					break
				}
			}
		}

		if !printed {
			// 交互模式时，13张牌的一向听分析显示改良具体情况
			detailFlag = true
			raw, _ = countsToString(counts)
			if _, _, err := analysis(raw); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
			detailFlag = false

			printed = true
		}

		for {
			printed = false

			fmt.Print("> 摸 ")
			fmt.Scanf("%s\n", &tile)
			idx, err := _convert(tile)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			} else {
				if counts[idx] == 4 {
					fmt.Fprintln(os.Stderr, "不可能摸更多的牌了")
				} else {
					counts[idx]++
					break
				}
			}
		}

		if !printed {
			raw, _ = countsToString(counts)
			if _, _, err := analysis(raw); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}

			printed = true
		}
	}
}

func main() {
	if len(os.Args) <= 1 {
		// 服务器模式
		runServer()
		return
	}

	if os.Args[len(os.Args)-1] == "-i" {
		// （一向听）交互模式
		interactFlag = true

		raw := strings.Join(os.Args[1:len(os.Args)-1], " ")
		interact(raw)
	}

	raw := strings.Join(os.Args[1:], " ")
	if os.Args[len(os.Args)-1] == "-d" {
		// 显示改良细节
		detailFlag = true
		raw = strings.Join(os.Args[1:len(os.Args)-1], " ")
	}

	t0 := time.Now()
	analysis(raw)
	fmt.Printf("耗时 %.2f 秒\n", float64(time.Now().UnixNano()-t0.UnixNano())/float64(time.Second))
}
