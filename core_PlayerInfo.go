package main

import (
	"fmt"

	"github.com/EndlessCheng/mahjong-helper/util"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"github.com/fatih/color"
)

type PlayerInfo struct {
	name string // 自家/下家/对家/上家

	selfWindTile int // 自风

	melds                []*model.Meld // 副露
	meldDiscardsAtGlobal []int
	meldDiscardsAt       []int
	isNaki               bool // 是否鳴牌（暗杠不算鳴牌）

	// 注意负数（自摸切）要^
	discardTiles          []int // 该玩家的舍牌
	latestDiscardAtGlobal int   // 该玩家最近一次舍牌在 globalDiscardTiles 中的下标，初始为 -1
	earlyOutsideTiles     []int // 立直前的1-5巡的外侧牌

	isReached  bool // 是否立直
	canIppatsu bool // 是否有一发

	reachTileAtGlobal int // 立直宣言牌在 globalDiscardTiles 中的下标，初始为 -1
	reachTileAt       int // 立直宣言牌在 discardTiles 中的下标，初始为 -1

	nukiDoraNum int // 拔北宝牌数
}

func (p *PlayerInfo) doraNum(doraList []int) (doraCount int) {
	for _, meld := range p.melds {
		for _, tile := range meld.Tiles {
			for _, doraTile := range doraList {
				if tile == doraTile {
					doraCount++
				}
			}
		}
		if meld.ContainRedFive {
			doraCount++
		}
	}
	if p.nukiDoraNum > 0 {
		doraCount += p.nukiDoraNum
		// 特殊：西为指示牌
		for _, doraTile := range doraList {
			if doraTile == 30 {
				doraCount += p.nukiDoraNum
			}
		}
	}
	return
}

func (p *PlayerInfo) printDiscards() {
	// TODO: 高亮不合理的舍牌或危险舍牌，如
	// - 一开始就切中张
	// - 开始切中张后，手切了幺九牌（也有可能是有人碰了牌，比如 133m 有人碰了 2m）
	// - 切了 dora，提醒一下
	// - 切了赤宝牌
	// - 有人立直的情况下，多次切出危险度高的牌（有可能是对方读准了牌，或者对方手里的牌与牌河加起来产生了安牌）
	// - 其余可以参考贴吧的《魔神之眼》翻译 https://tieba.baidu.com/p/3311909701
	//      举个简单的例子,如果出现手切了一个对子的情况的话那么基本上就不可能是七对子。
	//      如果对方早巡手切了一个两面搭子的话，那么就可以推理出他在做染手或者牌型是对子型，如果他立直或者鳴牌的话，也比较容易读出他的手牌。
	// https://tieba.baidu.com/p/3311909701
	//      鳴牌之后和终盘的手切牌要尽量记下来，别人手切之前的安牌应该先切掉
	// https://tieba.baidu.com/p/3372239806
	//      吃牌时候打出来的牌的颜色是危险的；碰之后全部的牌都是危险的

	fmt.Printf(p.name + ":")
	for i, disTile := range p.discardTiles {
		fmt.Printf(" ")
		// TODO: 显示 dora, 赤宝牌
		bgColor := color.BgBlack
		fgColor := color.FgWhite
		var tile string
		if disTile >= 0 { // 手切
			tile = util.Mahjong[disTile]
			if disTile >= 27 {
				tile = util.MahjongU[disTile] // 关注字牌的手切
			}
			if p.isNaki { // 副露
				fgColor = getOtherDiscardAlertColor(disTile) // 高亮中张手切
				if util.InInts(i, p.meldDiscardsAt) {
					bgColor = color.BgWhite // 鳴牌时切的那张牌要背景高亮
					fgColor = color.FgBlack
				}
			}
		} else { // 摸切
			disTile = ^disTile
			tile = util.Mahjong[disTile]
			fgColor = color.FgHiBlack // 暗色显示
		}
		color.New(bgColor, fgColor).Print(tile)
	}
	fmt.Println()
}
