package main

import (
	"strings"
	"fmt"
	"os"
	"time"
)

var detailFlag = false

var mahjong = [...]string{
	"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m",
	"1p", "2p", "3p", "4p", "5p", "6p", "7p", "8p", "9p",
	"1s", "2s", "3s", "4s", "5s", "6s", "7s", "8s", "9s",
	"dong", "nan", "xi", "bei",
	"zhong", "bai", "fa",
}

// 13张牌，检查是否听牌，返回听牌的具体情况
func checkTing0(cnt []int) needTiles {
	needs := needTiles{}
	for i := range mahjong {
		if cnt[i] == 4 {
			continue
		}
		cnt[i]++ // 摸牌
		if checkWin(cnt) { // 和牌
			needs[i] = 4 - (cnt[i] - 1)
		}
		cnt[i]--
	}
	return needs
}

var (
	buffer       = strings.Builder{}
	detailBuffer = strings.Builder{}
)

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
					tingCnt, _ := nd.parse()
					if tingCnt > tingCntMap[j] {
						// 听牌一般听数量最多的
						tingCntMap[j] = tingCnt
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
				if betterAllCount, betterTiles := betterNeeds.parse(); betterAllCount > allCount {
					// 进张数变多，则为一向听的改良
					impWay++
					if betterAllCount > improveCount[drawIdx] {
						improveCount[drawIdx] = betterAllCount
					}
					detailBuffer.WriteString(fmt.Sprintln(fmt.Sprintf("    摸 %s 切 %s 改良:", mahjong[drawIdx], mahjong[discardIdx]), betterAllCount, betterTiles))
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
			buffer.WriteString(fmt.Sprintf("%.2f [%d 变化]", avgImproveNum, impWay))
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
			if allCount, ans := checkTing1(cnt, true).parse(); allCount > 0 {
				colorNumber(allCount)
				fmt.Printf("    切 %s %v\n", mahjong[i], ans)
				fmt.Print(buffer.String())
				if detailFlag {
					fmt.Print(detailBuffer.String())
				}
				fmt.Println()

				ok = true
			}
			cnt[i]++
		}
	}
	return ok
}

func analysis(raw string) {
	fmt.Println(raw)
	fmt.Println(strings.Repeat("=", len(raw)))
	switch num, cnt := convert(raw); num {
	case 13:
		if needs := checkTing0(cnt); len(needs) > 0 {
			fmt.Println("已听牌:", needs.String())
		} else {
			allCount, ans := checkTing1(cnt, true).parse()
			fmt.Println("一向听:", allCount, ans)
			fmt.Println(buffer.String())
			buffer.Reset()
		}
	case 14:
		if checkWin(cnt) {
			fmt.Println("已胡牌")
		} else {
			if !checkTing1Discard(cnt) {
				fmt.Println("尚未一向听")
				// TODO
			}
		}
	default:
		_errorExit("参数错误")
	}

	//fmt.Println("checkWin", checkWinCount)
}

func main() {
	if len(os.Args) <= 1 {
		_errorExit("参数错误")
	}

	raw := strings.Join(os.Args[1:], " ")
	if os.Args[len(os.Args)-1] == "-d" {
		detailFlag = true
		raw = strings.Join(os.Args[1:len(os.Args)-1], " ")
	}

	t0 := time.Now()
	analysis(raw)
	fmt.Printf("耗时 %.2f 秒\n", float64(time.Now().UnixNano()-t0.UnixNano())/float64(time.Second))
}
