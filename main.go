package main

import (
	"strings"
	"fmt"
)

var mahjong = [...]string{
	"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m",
	"1p", "2p", "3p", "4p", "5p", "6p", "7p", "8p", "9p",
	"1s", "2s", "3s", "4s", "5s", "6s", "7s", "8s", "9s",
	"dong", "nan", "xi", "bei",
	"zhong", "bai", "fa",
}

var buffer = strings.Builder{}

// 13张牌，检查是否听牌，返回听牌的名称
func checkTing0(cnt []int) (tiles []string) {
	for i := range mahjong {
		if cnt[i] == 4 {
			continue
		}
		cnt[i]++ // 摸牌
		if checkWin(cnt) { // 和牌
			tiles = append(tiles, mahjong[i])
		}
		cnt[i]--
	}
	return tiles
}

// 13张牌，检查一向听
func checkTing1(cnt []int, recur bool) needTiles {
	needs := needTiles{}
	betterNeedsMap := map[int]map[int]needTiles{}
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
				if ting := checkTing0(cnt); len(ting) > 0 {
					// 若能听牌，则换的这张牌为一向听的进张
					if _, ok := needs[j]; !ok {
						needs[j] = 4 - (cnt[j] - 1)
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

	// TODO: 振听?
	if allCount, tiles := needs.parse(); allCount > 0 {
		maxBetterAllCount := -1
		for _, tmpNeedsMap := range betterNeedsMap {
			for drawIdx, betterNeeds := range tmpNeedsMap {
				if betterAllCount, _ := betterNeeds.parse(); betterAllCount > maxBetterAllCount && !in(mahjong[drawIdx], tiles) {
					maxBetterAllCount = betterAllCount
				}
			}
		}
		if maxBetterAllCount > allCount {
			// 进张数变多，则为一向听的改良
			for discardIdx, tmpNeedsMap := range betterNeedsMap {
				for drawIdx, betterNeeds := range tmpNeedsMap {
					if in(mahjong[drawIdx], tiles) {
						continue
					}
					if betterAllCount, betterTiles := betterNeeds.parse(); betterAllCount == maxBetterAllCount {
						buffer.WriteString(fmt.Sprintln(fmt.Sprintf("\t摸 %s 切 %s 改良:", mahjong[drawIdx], mahjong[discardIdx]), betterAllCount, betterTiles, ))
					}
				}
			}
		}
	}

	return needs
}

// 14张牌，一向听，何切
// 检查能进入一向听的牌，对比 1.进张数，2.是否能改良
func checkTing1Discard(cnt []int) {
	for i := range mahjong {
		if cnt[i] >= 1 {
			cnt[i]-- // 切牌
			if allCount, ans := checkTing1(cnt, true).parse(); allCount > 0 {
				fmt.Println(fmt.Sprintf("切 %s:", mahjong[i]), allCount, ans)
				fmt.Println(buffer.String())
				buffer.Reset()
			}
			cnt[i]++
		}
	}
}

func main() {
	//if len(os.Args) <= 1 {
	//	_errorExit("参数错误")
	//}
	//raw := os.Args[1]
	//raw := "11222333789s fa fa"
	//raw := "2355789p 356778s"
	raw := "4578999m 45p 11145s"
	num, cnt := convert(raw)
	switch num {
	case 13:
		if ans := checkTing0(cnt); len(ans) > 0 {
			fmt.Println("听牌:", strings.Join(ans, " "))
		} else {
			allCount, ans := checkTing1(cnt, true).parse()
			fmt.Println("一向听:", allCount, ans)
		}
	case 14:
		checkTing1Discard(cnt)
	default:
		_errorExit("参数错误")
	}
}
