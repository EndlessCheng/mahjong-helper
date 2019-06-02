package main

import (
	"fmt"
	"github.com/EndlessCheng/mahjong-helper/util"
	"github.com/fatih/color"
)

type analysisOpType int

const (
	analysisOpTypeTsumo analysisOpType = iota
	analysisOpTypeNaki
	analysisOpTypeKan  // 加杠和暗杠
)

type analysisCache struct {
	analysisOpType analysisOpType

	selfDiscardTile int
	//isSelfDiscardRedFive bool
	selfDiscardTileRisk float64
	isRiichiWhenDiscard bool

	// 用手牌中的什么牌去鸣牌，空就是跳过不鸣
	selfOpenTiles []int

	aiAttackDiscardTile      int
	aiDefenceDiscardTile     int
	aiAttackDiscardTileRisk  float64
	aiDefenceDiscardTileRisk float64

	tenpaiRate []float64 // TODO: 三家听牌率
}

type roundAnalysisCache struct {
	isStart bool
	isEnd   bool
	cache   []*analysisCache
}

func (rc *roundAnalysisCache) print() {
	const baseInfo = "助手正在计算推荐舍牌，请稍等……（计算结果仅供参考）"
	const sep = "  "
	done := rc != nil && rc.isEnd
	if !done {
		color.HiGreen(baseInfo)
	}
	fmt.Print("巡目　　")
	if done {
		for i := range rc.cache {
			fmt.Printf("%s%2d", sep, i+1)
		}
	}
	fmt.Println()
	fmt.Print("自家切牌")
	if done {
		for _, c := range rc.cache {
			info := ""
			if c.selfDiscardTile != -1 {
				info = util.Mahjong[c.selfDiscardTile]
			}
			if c.isRiichiWhenDiscard {
				info += "[立直]"
			}
			fmt.Print(sep)
			if c.selfDiscardTileRisk < 5 {
				fmt.Print(info)
			} else {
				color.New(getNumRiskColor(c.selfDiscardTileRisk)).Print(info)
			}
		}
	}
	fmt.Println()
	fmt.Print("进攻推荐")
	if done {
		for _, c := range rc.cache {
			info := ""
			if c.aiAttackDiscardTile != -1 {
				info = util.Mahjong[c.aiAttackDiscardTile]
			}
			fmt.Print(sep)
			if c.aiAttackDiscardTileRisk < 5 {
				fmt.Print(info)
			} else {
				color.New(getNumRiskColor(c.aiAttackDiscardTileRisk)).Print(info)
			}
		}
	}
	fmt.Println()
	fmt.Print("防守推荐")
	if done {
		for _, c := range rc.cache {
			info := "--"
			if c.aiDefenceDiscardTile != -1 {
				info = util.Mahjong[c.aiDefenceDiscardTile]
			}
			fmt.Print(sep)
			if c.aiDefenceDiscardTileRisk < 5 {
				fmt.Print(info)
			} else {
				color.New(getNumRiskColor(c.aiDefenceDiscardTileRisk)).Print(info)
			}
		}
	}
	fmt.Println()
	fmt.Println()
}

func (rc *roundAnalysisCache) addSelfDiscardTile(tile int, risk float64, isRiichiWhenDiscard bool) {
	latestCache := rc.cache[len(rc.cache)-1]
	latestCache.selfDiscardTile = tile
	latestCache.selfDiscardTileRisk = risk
	latestCache.isRiichiWhenDiscard = isRiichiWhenDiscard
}

func (rc *roundAnalysisCache) addAIDiscardTile(attackTile int, defenceTile int, attackTileRisk float64, defenceDiscardTileRisk float64) {
	rc.cache = append(rc.cache, &analysisCache{
		aiAttackDiscardTile:      attackTile,
		aiDefenceDiscardTile:     defenceTile,
		aiAttackDiscardTileRisk:  attackTileRisk,
		aiDefenceDiscardTileRisk: defenceDiscardTileRisk,
	})
}

//

type gameAnalysisCache struct {
	// roundNumber ben 巡目
	wholeGameCache [][]*roundAnalysisCache

	majsoulRecordUUID string

	selfSeat int
}

func newGameAnalysisCache(majsoulRecordUUID string, selfSeat int) *gameAnalysisCache {
	cache := make([][]*roundAnalysisCache, 12)
	for i := range cache {
		cache[i] = make([]*roundAnalysisCache, 20)
	}
	return &gameAnalysisCache{
		wholeGameCache:    cache,
		majsoulRecordUUID: majsoulRecordUUID,
		selfSeat:          selfSeat,
	}
}

var globalAnalysisCache *gameAnalysisCache

func (c *gameAnalysisCache) runMajsoulRecordAnalysisTask(actions []*majsoulRecordAction) error {
	// 从第一个 action 中取出局和场
	if len(actions) == 0 {
		return fmt.Errorf("数据异常：此局数据为空")
	}

	newRoundAction := actions[0]
	data := newRoundAction.Action
	roundNumber := 4*(*data.Chang) + *data.Ju
	ben := *data.Ben
	roundCache := c.wholeGameCache[roundNumber][ben] // TODO: 建议用原子操作
	if roundCache == nil {
		roundCache = &roundAnalysisCache{isStart: true}
		if debugMode {
			fmt.Println("助手正在计算推荐舍牌…… 创建 roundCache")
		}
		c.wholeGameCache[roundNumber][ben] = roundCache
	} else if roundCache.isStart {
		return nil
	}

	// 遍历自家舍牌，找到舍牌前的操作
	// 若为摸牌操作，计算出此时的 AI 进攻舍牌和防守舍牌
	// 若为鸣牌操作，计算出此时的 AI 进攻舍牌（无进攻舍牌则设为 -1），防守舍牌设为 -1
	// TODO: 玩家跳过，但是 AI 觉得应鸣牌？
	majsoulRoundData := &majsoulRoundData{
		accountID: gameConf.MajsoulAccountID,
		selfSeat:  c.selfSeat,
	}
	majsoulRoundData.roundData = newGame(majsoulRoundData)
	majsoulRoundData.roundData.gameMode = gameModeRecordCache
	majsoulRoundData.skipOutput = true
	for i, action := range actions[:len(actions)-1] {
		h := getGlobalMJHandler()
		if c.majsoulRecordUUID != h.majsoulCurrentRecordUUID {
			if debugMode {
				fmt.Println("用户退出该牌谱")
			}
			return nil
		}
		if debugMode {
			fmt.Println("助手正在计算推荐舍牌…… action", i)
		}
		majsoulRoundData.msg = action.Action
		majsoulRoundData.analysis()
	}
	roundCache.isEnd = true

	if c.majsoulRecordUUID != h.majsoulCurrentRecordUUID {
		if debugMode {
			fmt.Println("用户退出该牌谱")
		}
		return nil
	}

	clearConsole()
	roundCache.print()

	return nil
}
