package main

import (
	"strings"
	"fmt"
	"os"
	"time"
)

var detailFlag = false
var interactFlag = false // 交互模式

// 13张牌，检查是否听牌，返回听牌的具体情况
// （不考虑国士无双）
func checkTing0(cnt []int) needTiles {
	needs := needTiles{}

	// 剪枝：检测浮牌
	// 此处优化提升了 7-10 倍的性能
	for i := 0; i < 3; i++ {
		cnt0 := 2
		for j := 0; j < 9; j++ {
			idx := 9*i + j
			c := cnt[idx]
			switch {
			case c == 0:
				cnt0++
			case c == 1:
				if cnt0 < 2 {
					cnt0 = 0
				} else {
					if j+1 < 9 && cnt[idx+1] > 0 {
						j++
						cnt0 = 0
					} else if j+2 < 9 && cnt[idx+2] > 0 {
						j += 2
						cnt0 = 0
					} else {
						// 这是一张浮牌，要想和牌只能单骑这一张
						cnt[idx]++ // 摸牌
						if checkWin(cnt) { // 单骑和牌
							needs[idx] = 4 - (cnt[idx] - 1)
						}
						cnt[idx]--
						return needs
					}
				}
			case c >= 2:
				cnt0 = 0
			}
		}
	}
	for i := 27; i < len(mahjong); i++ {
		if cnt[i] == 1 {
			// 这是一张浮牌，要想和牌只能单骑这一张
			cnt[i]++ // 摸牌
			if checkWin(cnt) { // 单骑和牌
				needs[i] = 4 - (cnt[i] - 1)
			}
			cnt[i]--
			return needs
		}
	}

	// 剪枝：下面这段可以加快 35%~40%
	// 只计算可能胡的牌，对于 345m ... 这样的牌来说，摸到 1m 是肯定不能胡牌的
	needChecks := make([]int, 0, len(mahjong))
	for i := 0; i < 3; i++ {
		for j := 0; j < 9; j++ {
			idx := 9*i + j
			if cnt[idx] == 4 {
				continue
			}
			if j == 0 || j == 9 {
				if cnt[idx] > 0 || cnt[idx+1] > 0 {
					needChecks = append(needChecks, idx)
				}
			} else {
				if cnt[idx-1] > 0 || cnt[idx] > 0 || cnt[idx+1] > 0 {
					needChecks = append(needChecks, idx)
				}
			}
		}
	}
	for i := 27; i < 34; i++ {
		if cnt[i] == 4 {
			continue
		}
		if cnt[i] > 0 {
			needChecks = append(needChecks, i)
		}
	}
	for _, idx := range needChecks {
		cnt[idx]++ // 摸牌
		if checkWin(cnt) { // 和牌
			needs[idx] = 4 - (cnt[idx] - 1)
		}
		cnt[idx]--
	}

	//for i := range mahjong {
	//	if cnt[i] == 4 {
	//		continue
	//	}
	//	cnt[i]++ // 摸牌
	//	if checkWin(cnt) { // 和牌
	//		needs[i] = 4 - (cnt[i] - 1)
	//	}
	//	cnt[i]--
	//}

	return needs
}

// 检查切掉某张牌后是否听牌
func checkTing0Discard(cnt []int) bool {
	ok := false
	for i := range mahjong {
		if cnt[i] >= 1 {
			cnt[i]-- // 切牌
			if needs := checkTing0(cnt); len(needs) > 0 {
				ok = true
				fmt.Println("【已听牌！】 切", mahjongZH[i], needs.String())
			}
			cnt[i]++
		}
	}
	return ok
}

var (
	buffer       = strings.Builder{}
	detailBuffer = strings.Builder{}
)

func flushBuffer() {
	fmt.Print(buffer.String())
	if detailFlag {
		fmt.Print(detailBuffer.String())
	}
	fmt.Println()

	buffer.Reset()
	detailBuffer.Reset()
}

// 13张牌，检查一向听
func checkTing1(cnt []int, recur bool) needTiles {
	needs := needTiles{}
	betterNeedsMap := map[int]map[int]needTiles{}
	tingCntMap := map[int]int{} // map[摸到idx]听多少张牌

	for i := range mahjong {
		if cnt[i] >= 1 {
			tmpNeedsMap := map[int]needTiles{}
			cnt[i]-- // 切掉其中一张牌
			for j := range mahjong {
				if j == i {
					continue
				}
				if cnt[j] == 4 {
					continue
				}
				cnt[j]++ // 换成其他牌
				if nd := checkTing0(cnt); len(nd) > 0 {
					// 若能听牌，则换的这张牌为一向听的进张
					if _, ok := needs[j]; !ok {
						needs[j] = 4 - (cnt[j] - 1)
					} else {
						// 比如说 57m22566s，切 5s/6s 来 8m 都听牌，但是听牌的数量有区别
					}
					if recur {
						tingCnt, _ := nd.parse()
						if tingCnt > tingCntMap[j] {
							// 听牌一般听数量最多的
							tingCntMap[j] = tingCnt
						}
					}
				} else if recur {
					if betterNeeds := checkTing1(cnt, false); len(betterNeeds) > 0 {
						// 换成这张牌也是一向听，可能是改良型，记录一下
						tmpNeedsMap[j] = betterNeeds
					}
				}
				cnt[j]--
			}
			cnt[i]++
			betterNeedsMap[i] = tmpNeedsMap
		}
	}

	if !recur {
		return needs
	}

	buffer.Reset()
	detailBuffer.Reset()

	// TODO: 振听?
	if allCount, tiles := needs.parse(); allCount > 0 {
		improveCount := make([]int, len(mahjong))
		for i := range mahjong {
			improveCount[i] = allCount
		}
		impWay := 0
		for discardIdx, tmpNeedsMap := range betterNeedsMap {
			for drawIdx, betterNeeds := range tmpNeedsMap {
				if in(mahjong[drawIdx], tiles) {
					// 跳过改良牌就是一向听的进张的情况
					continue
				}
				if betterAllCount, betterTiles := betterNeeds.parseZH(); betterAllCount > allCount {
					// 进张数变多，则为一向听的改良
					impWay++
					if betterAllCount > improveCount[drawIdx] {
						improveCount[drawIdx] = betterAllCount
					}
					detailBuffer.WriteString(fmt.Sprintln(fmt.Sprintf("    摸 %s 切 %s 改良:", mahjongZH[drawIdx], mahjongZH[discardIdx]), betterAllCount, betterTiles))
				}
			}
		}

		if detailBuffer.Len() > 0 {
			improveScore := 0
			weight := 0
			for i := range mahjong {
				w := 4 - cnt[i]
				improveScore += w * improveCount[i]
				weight += w
			}
			avgImproveNum := float64(improveScore) / float64(weight)
			buffer.WriteString(fmt.Sprintf("%.2f [%d 改良]", avgImproveNum, impWay))
		} else {
			buffer.WriteString(strings.Repeat(" ", 14))
		}

		avgTingSum := 0
		weight := 0
		for idx, c := range tingCntMap {
			w := 4 - cnt[idx]
			avgTingSum += w * c
			weight += w
		}
		avgTingNum := float64(avgTingSum) / float64(weight)
		avgTingStr := fmt.Sprintf("%.2f 听牌数", avgTingNum)
		// TODO: 根据1-9的牌来计算综合和牌率
		buffer.WriteString("  " + avgTingStr + "\n")
	}

	return needs
}

// 14张牌，可以一向听，何切
// 检查能一向听的切牌，对比：
// 1. 进张数
// 2. 改良之后的（加权）平均进张数
// 3. 听牌后的（加权）平均听牌数
// 4. 听牌后所听牌的名称（就是一向听的进张名称）（一般来说 14m 优于 25m。不过还是要根据场况来判断）
func checkTing1Discard(cnt []int) bool {
	ok := false
	for i := range mahjong {
		if cnt[i] >= 1 {
			cnt[i]-- // 切牌
			if allCount, ans := checkTing1(cnt, true).parseZH(); allCount > 0 {
				ok = true

				colorNumber1(allCount)
				fmt.Printf("    切 %s %v\n", mahjongZH[i], ans)
				flushBuffer()
			}
			cnt[i]++
		}
	}
	return ok
}

// 13张牌，检查一向听（简化版）
func _simpleCheckTing1(cnt []int) needTiles {
	needs := needTiles{}
	for i := range mahjong {
		if cnt[i] >= 1 {
			cnt[i]-- // 切掉其中一张牌
			for j := range mahjong {
				if j == i {
					continue
				}
				if cnt[j] == 4 {
					continue
				}
				cnt[j]++ // 换成其他牌
				if nd := checkTing0(cnt); len(nd) > 0 {
					// 若能听牌，则换的这张牌为一向听的进张
					if _, ok := needs[j]; !ok {
						needs[j] = 4 - (cnt[j] - 1)
					} else {
						// 比如说 57m22566s，切 5s/6s 来 8m 都听牌
					}
				}
				cnt[j]--
			}
			cnt[i]++
		}
	}
	return needs
}

// 13张牌，检查两向听
func checkTing2(cnt []int) needTiles {
	needs := needTiles{}
	for i := range mahjong {
		if cnt[i] >= 1 {
			cnt[i]-- // 切掉其中一张牌
			for j := range mahjong {
				if j == i {
					continue
				}
				if cnt[j] == 4 {
					continue
				}
				cnt[j]++ // 换成其他牌
				if nd := _simpleCheckTing1(cnt); len(nd) > 0 {
					// 若能一向听，则换的这张牌为两向听的进张
					if _, ok := needs[j]; !ok {
						needs[j] = 4 - (cnt[j] - 1)
					}
				}
				cnt[j]--
			}
			cnt[i]++
		}
	}
	return needs
}

// 交互模式下，两向听的最低值
var ting2MinCount = -1

func reset() {
	ting2MinCount = -1
}

// 14张牌，可以两向听，何切
func checkTing2Discard(cnt []int) bool {
	ok := false
	for i := range mahjong {
		if cnt[i] >= 1 {
			cnt[i]-- // 切牌
			if allCount, ans := checkTing2(cnt).parse(); allCount > 0 {
				ok = true

				if allCount >= ting2MinCount {
					colorNumber2(allCount)
					fmt.Printf("   切 %s %v\n", mahjongZH[i], ans)
				}
			}
			cnt[i]++
		}
	}
	fmt.Println()
	return ok
}

func analysis(raw string) (num int, cnt []int, err error) {
	fmt.Println(raw)
	fmt.Println(strings.Repeat("=", len(raw)))

	num, cnt, err = convert(raw)
	if err != nil {
		return
	}

	if countDui(cnt) >= 4 {
		alert("对子手可能")
	}

	switch num {
	case 13:
		if needs := checkTing0(cnt); len(needs) > 0 {
			fmt.Println("已听牌:", needs.String())
		} else {
			allCount, ans := checkTing1(cnt, true).parseZH()
			if allCount > 0 {
				fmt.Println("一向听:", allCount, ans)
				flushBuffer()
			} else {
				allCount, ans := checkTing2(cnt).parseZH()
				if allCount > 0 {
					fmt.Println("两向听:", allCount, ans)
					flushBuffer()

					// 13设置
					ting2MinCount = allCount
				} else {
					fmt.Println("尚未两向听")
				}
			}
		}
	case 14:
		if checkWin(cnt) {
			fmt.Println("已胡牌")
		} else {
			if !checkTing0Discard(cnt) {
				if !checkTing1Discard(cnt) {
					if !checkTing2Discard(cnt) {
						fmt.Println("尚未两向听")
					}
				}
			}
		}
		// 14失效
		reset()
	default:
		err = fmt.Errorf("参数错误: %s（%d 张牌）", raw, num)
		return
	}

	//fmt.Println("checkWin", checkWinCount)

	return
}

func interact(raw string) {
	num, cnt, err := analysis(raw)
	if err != nil {
		_errorExit(err.Error())
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
				if cnt[idx] == 0 {
					fmt.Fprintln(os.Stderr, "切掉的牌不存在")
				} else {
					cnt[idx]--
					break
				}
			}
		}

		if !printed {
			// 交互模式时，13张牌的一向听分析显示改良具体情况
			detailFlag = true
			raw, _ = countToString(cnt)
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
				if cnt[idx] == 4 {
					fmt.Fprintln(os.Stderr, "不可能摸更多的牌了")
				} else {
					cnt[idx]++
					break
				}
			}
		}

		if !printed {
			raw, _ = countToString(cnt)
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
