package main

import (
	"strings"
	"fmt"
	"os"
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
	return -1
}

func convert(raw string) []int {
	var result []int

	raw = strings.TrimSpace(raw)
	splits := strings.Split(raw, " ")

	for _, split := range splits {
		if split[0] >= '1' && split[0] <= '9' {
			for i := range split[:len(split)-1] {
				single := split[i:i+1] + split[len(split)-1:]
				r := _convert(single)
				if r == -1 {
					return nil
				}
				result = append(result, r)
			}
		} else {
			r := _convert(split)
			if r == -1 {
				return nil
			}
			result = append(result, r)
		}
	}

	return result
}

var cnt = make([]int, 34)

func search(dep int) bool {
	for i := range mahjong {
		if cnt[i] >= 3 { // 刻子
			if dep == 3 { // 4组面子
				return true
			}
			cnt[i] -= 3
			ok := search(dep + 1)
			cnt[i] += 3
			if ok {
				return true
			}
		}
		for i := 0; i <= 24; i++ { // 一直到 7s
			if i%9 <= 6 && cnt[i] >= 1 && cnt[i+1] >= 1 && cnt[i+2] >= 1 { // 顺子
				if dep == 3 { // 4组面子
					return true
				}
				cnt[i]--
				cnt[i+1]--
				cnt[i+2]--
				ok := search(dep + 1)
				cnt[i+2]++
				cnt[i+1]++
				cnt[i]++
				if ok {
					return true
				}
			}
		}
	}
	return false
}

// 检查是否和牌
func checkWin() bool {
	for i := range mahjong {
		if cnt[i] >= 2 { // 雀头
			cnt[i] -= 2
			ok := search(0)
			cnt[i] += 2
			if ok {
				return true
			}
		}
	}
	return false
}

func checkTing0() (ans []string) {
	for i := range mahjong {
		if cnt[i] == 4 {
			continue
		}
		cnt[i]++ // 摸牌
		if checkWin() { // 和牌
			ans = append(ans, mahjong[i])
		}
		cnt[i]--
	}
	return ans
}

// 一向听，13张牌
// 切掉其中一张牌，换成其他牌，若能听牌，则为一向听的进张
// 切掉其中一张牌，换成其他牌，若依然为一向听，但进张数变多，则为一向听的改良
func checkTing1() (allCount int, ans []string) {
	for i := range mahjong {
		if cnt[i] >= 1 {

			fmt.Println("i=",i)

			cnt[i]-- // 切掉其中一张牌
			for j := range mahjong { // 换成其他牌
				if j == i {
					continue
				}
				if cnt[j] == 4 {
					continue
				}

				fmt.Println("j=",j)

				cnt[j]++
				if ting := checkTing0(); len(ting) > 0 {
					allCount += 4 - (cnt[j] - 1)
					ans = append(ans, mahjong[j])
					fmt.Println(allCount, ans)
				} else {
					//anotherAllCount,anotherAns:=checkTing1()
				}
				cnt[j]--
			}
			cnt[i]++
		}
	}
	return
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
	mj := convert(raw)
	if mj == nil {
		fmt.Fprintln(os.Stderr, "参数错误")
		os.Exit(1)
	}
	for _, m := range mj {
		cnt[m]++
		if cnt[m] > 4 {
			fmt.Fprintln(os.Stderr, "超过4张一样的牌！")
			os.Exit(1)
		}
	}

	ans := checkTing0()
	if len(ans) > 0 {
		fmt.Println(strings.Join(ans, " "))
	} else {
		checkTing1()

	}
}
