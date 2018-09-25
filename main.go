package main

import (
	"strings"
	"fmt"
	"os"
	"sort"
)

var mahjong = [...]string{
	"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m",
	"1p", "2p", "3p", "4p", "5p", "6p", "7p", "8p", "9p",
	"1s", "2s", "3s", "4s", "5s", "6s", "7s", "8s", "9s",
	"dong", "nan", "xi", "bei",
	"zhong", "bai", "fa",
}

func _convert(single string) int {
	for i, m := range mahjong {
		if m == single {
			return i
		}
	}
	fmt.Fprintln(os.Stderr, "参数错误:", single)
	os.Exit(1)
	return -1
}

func convert(raw string) []int {
	cnt := make([]int, 34)

	raw = strings.TrimSpace(raw)
	splits := strings.Split(raw, " ")

	var result []int
	for _, split := range splits {
		if split[0] >= '1' && split[0] <= '9' {
			for i := range split[:len(split)-1] {
				single := split[i:i+1] + split[len(split)-1:]
				result = append(result, _convert(single))
			}
		} else {
			result = append(result, _convert(split))
		}
	}

	for _, m := range result {
		cnt[m]++
		if cnt[m] > 4 {
			fmt.Fprintln(os.Stderr, "参数错误: 超过4张一样的牌！")
			os.Exit(1)
		}
	}

	return cnt
}

var cnt = make([]int, 34)

// 13张牌，检查是否听牌
func checkTing0() (ans []string) {
	for i := range mahjong {
		if cnt[i] == 4 {
			continue
		}
		cnt[i]++ // 摸牌
		if checkWin(cnt) { // 和牌
			ans = append(ans, mahjong[i])
		}
		cnt[i]--
	}
	return ans
}

// 13张牌，检查一向听的进张数，进张名称
func checkTing1() (allCount int, ans []string) {
	cards := map[int]int{}
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
				if ting := checkTing0(); len(ting) > 0 {
					// 若能听牌，则换的这张牌为一向听的进张
					if _, ok := cards[j]; !ok {
						cards[j] = 4 - (cnt[j] - 1)
					}
				} else {
					// 若依然为一向听，但进张数变多，则为一向听的改良
					// TODO anotherAllCount,anotherAns:=checkTing1()
				}
				cnt[j]--
			}
			cnt[i]++
		}
	}

	idxAns := make([]int, 0, len(cards))
	for k, v := range cards {
		idxAns = append(idxAns, k)
		allCount += v
	}
	sort.Ints(idxAns)

	ans = make([]string, len(idxAns))
	for i, idx := range idxAns {
		ans[i] = mahjong[idx]
	}

	return allCount, ans
}

// 一向听，何切，14张牌
// 检查能进入一向听的牌，对比 1.进张数，2.是否能改良
func checkTing1What() {

}

func main() {
	//if len(os.Args) <= 1 {
	//	fmt.Fprintln(os.Stderr, "参数错误")
	//	os.Exit(1)
	//}
	//raw := os.Args[1]
	//raw := "11222333789s fa fa"
	raw := "2355789p 356778s"
	cnt = convert(raw)

	if ans := checkTing0(); len(ans) > 0 {
		fmt.Println(strings.Join(ans, " "))
	} else {
		fmt.Println(checkTing1())
	}
}
