package util

import (
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"fmt"
)

type shantenSearchNode14 struct {
	shanten  int
	children map[int]*shantenSearchNode13 // 向听不变的舍牌-node13
}

func (n *shantenSearchNode14) printWithPrefix(prefix string) string {
	if n == nil || n.shanten == -1 {
		return prefix + "end\n"
	}
	output := ""
	for discardTile, node13 := range n.children {
		output += prefix + fmt.Sprintln("舍", Mahjong[discardTile]) + node13.printWithPrefix(prefix+"  ")
	}
	return output
}

func (n *shantenSearchNode14) String() string {
	return n.printWithPrefix("")
}

type shantenSearchNode13 struct {
	shanten  int
	waits    Waits
	children map[int]*shantenSearchNode14 // 向听前进的摸牌-node14
}

func (n *shantenSearchNode13) printWithPrefix(prefix string) string {
	output := ""
	for drawTile, node14 := range n.children {
		output += prefix + fmt.Sprintln("摸", Mahjong[drawTile]) + node14.printWithPrefix(prefix+"  ")
	}
	return output
}

func (n *shantenSearchNode13) String() string {
	return n.printWithPrefix("")
}

func _search14(shanten int, playerInfo *model.PlayerInfo) *shantenSearchNode14 {
	// 不需要判断 shanten 是否为 -1：因为_search13 中用的是 IsAgari，所以 shanten 是 >=0 的
	children := map[int]*shantenSearchNode13{}
	tiles34 := playerInfo.HandTiles34
	for i := 0; i < 34; i++ {
		if tiles34[i] == 0 {
			continue
		}
		tiles34[i]--
		// 向听不变的舍牌
		if CalculateShanten(tiles34) == shanten {
			children[i] = _search13(shanten, playerInfo)
		}
		tiles34[i]++
	}

	return &shantenSearchNode14{
		shanten:  shanten,
		children: children,
	}
}

func _search13(shanten int, playerInfo *model.PlayerInfo) *shantenSearchNode13 {
	waits := Waits{}
	children := map[int]*shantenSearchNode14{}
	tiles34 := playerInfo.HandTiles34
	leftTiles34 := playerInfo.LeftTiles34
	isTenpai := shanten == 0
	for i := 0; i < 34; i++ {
		if tiles34[i] == 4 {
			continue
		}
		tiles34[i]++
		if isTenpai {
			// 优化：听牌时改用更为快速的 IsAgari
			if IsAgari(tiles34) {
				waits[i] = leftTiles34[i]
				children[i] = nil
			}
		} else {
			if CalculateShanten(tiles34) < shanten {
				// 向听前进了，则换的这张牌为进张，进张数即剩余枚数
				// 有可能为 0，但考虑到判断振听时需要进张种类，所以记录
				waits[i] = leftTiles34[i]
				if leftTiles34[i] > 0 {
					leftTiles34[i]--
					children[i] = _search14(shanten-1, playerInfo)
					leftTiles34[i]++
				} else {
					children[i] = nil
				}
			}
		}
		tiles34[i]--
	}

	return &shantenSearchNode13{
		shanten:  shanten,
		waits:    waits,
		children: children,
	}
}

// 通过快速地搜索得出进张数据，供后续分析用
func searchShanten14(shanten int, playerInfo *model.PlayerInfo) *shantenSearchNode14 {
	if shanten == -1 {
		return &shantenSearchNode14{
			shanten:  shanten,
			children: map[int]*shantenSearchNode13{},
		}
	}
	return _search14(shanten, playerInfo)
}
