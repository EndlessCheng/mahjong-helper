package main

import (
	"fmt"
	"github.com/EndlessCheng/mahjong-helper/util"
	"github.com/fatih/color"
)

type analysisOpType int

const (
	analysisOpTypeTsumo     analysisOpType = iota
	analysisOpTypeChiPonKan  // 吃 碰 明杠
	analysisOpTypeKan        // 加杠 暗杠
)

// TODO: 提醒「此处应该副露，不应跳过」

type analysisCache struct {
	analysisOpType analysisOpType

	selfDiscardTile int
	//isSelfDiscardRedFive bool
	selfDiscardTileRisk float64
	isRiichiWhenDiscard bool
	meldType            int

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

	analysisCacheBeforeChiPon *analysisCache
}

func (rc *roundAnalysisCache) print() {
	const (
		baseInfo  = "助手正在计算推荐舍牌，请稍等……（计算结果仅供参考）"
		emptyInfo = "--"
		sep       = "  "
	)

	done := rc != nil && rc.isEnd
	if !done {
		color.HiGreen(baseInfo)
	} else {
		// 检查最后的是否自摸，若为自摸则去掉推荐
		if len(rc.cache) > 0 {
			latestCache := rc.cache[len(rc.cache)-1]
			if latestCache.selfDiscardTile == -1 {
				latestCache.aiAttackDiscardTile = -1
				latestCache.aiDefenceDiscardTile = -1
			}
		}
	}

	fmt.Print("巡目　　")
	if done {
		for i := range rc.cache {
			fmt.Printf("%s%2d", sep, i+1)
		}
	}
	fmt.Println()

	printTileInfo := func(tile int, risk float64, suffix string) {
		info := emptyInfo
		if tile != -1 {
			info = util.Mahjong[tile]
		}
		fmt.Print(sep)
		if info == emptyInfo || risk < 5 {
			fmt.Print(info)
		} else {
			color.New(getNumRiskColor(risk)).Print(info)
		}
		fmt.Print(suffix)
	}

	fmt.Print("自家切牌")
	if done {
		for i, c := range rc.cache {
			suffix := ""
			if c.isRiichiWhenDiscard {
				suffix = "[立直]"
			} else if c.selfDiscardTile == -1 && i == len(rc.cache)-1 {
				//suffix = "[自摸]"
				// TODO: 流局
			}
			printTileInfo(c.selfDiscardTile, c.selfDiscardTileRisk, suffix)
		}
	}
	fmt.Println()

	fmt.Print("进攻推荐")
	if done {
		for _, c := range rc.cache {
			printTileInfo(c.aiAttackDiscardTile, c.aiAttackDiscardTileRisk, "")
		}
	}
	fmt.Println()

	fmt.Print("防守推荐")
	if done {
		for _, c := range rc.cache {
			printTileInfo(c.aiDefenceDiscardTile, c.aiDefenceDiscardTileRisk, "")
		}
	}
	fmt.Println()

	fmt.Println()
}

// （摸牌后、鸣牌后的）实际舍牌
func (rc *roundAnalysisCache) addSelfDiscardTile(tile int, risk float64, isRiichiWhenDiscard bool) {
	latestCache := rc.cache[len(rc.cache)-1]
	latestCache.selfDiscardTile = tile
	latestCache.selfDiscardTileRisk = risk
	latestCache.isRiichiWhenDiscard = isRiichiWhenDiscard
}

// 摸牌时的切牌推荐
func (rc *roundAnalysisCache) addAIDiscardTileWhenDrawTile(attackTile int, defenceTile int, attackTileRisk float64, defenceDiscardTileRisk float64) {
	// 摸牌，巡目+1
	rc.cache = append(rc.cache, &analysisCache{
		analysisOpType:           analysisOpTypeTsumo,
		selfDiscardTile:          -1,
		aiAttackDiscardTile:      attackTile,
		aiDefenceDiscardTile:     defenceTile,
		aiAttackDiscardTileRisk:  attackTileRisk,
		aiDefenceDiscardTileRisk: defenceDiscardTileRisk,
	})
	rc.analysisCacheBeforeChiPon = nil
}

// 加杠 暗杠
func (rc *roundAnalysisCache) addKan(meldType int) {
	// latestCache 是摸牌
	latestCache := rc.cache[len(rc.cache)-1]
	latestCache.analysisOpType = analysisOpTypeKan
	latestCache.meldType = meldType
	// 杠完之后又会摸牌，巡目+1
}

// 吃 碰 明杠
func (rc *roundAnalysisCache) addChiPonKan(meldType int) {
	if meldType == meldTypeMinkan {
		// 暂时忽略明杠，巡目不+1，留给摸牌时+1
		return
	}
	// 巡目+1
	var newCache *analysisCache
	if rc.analysisCacheBeforeChiPon != nil {
		newCache = rc.analysisCacheBeforeChiPon // 见 addPossibleChiPonKan
		newCache.analysisOpType = analysisOpTypeChiPonKan
		newCache.meldType = meldType
		rc.analysisCacheBeforeChiPon = nil
	} else {
		// 此处代码应该不会触发
		if debugMode {
			panic("rc.analysisCacheBeforeChiPon == nil")
		}
		newCache = &analysisCache{
			analysisOpType:       analysisOpTypeChiPonKan,
			selfDiscardTile:      -1,
			aiAttackDiscardTile:  -1,
			aiDefenceDiscardTile: -1,
			meldType:             meldType,
		}
	}
	rc.cache = append(rc.cache, newCache)
}

// 吃 碰 杠 跳过
func (rc *roundAnalysisCache) addPossibleChiPonKan(attackTile int, attackTileRisk float64) {
	rc.analysisCacheBeforeChiPon = &analysisCache{
		analysisOpType:          analysisOpTypeChiPonKan,
		selfDiscardTile:         -1,
		aiAttackDiscardTile:     attackTile,
		aiDefenceDiscardTile:    -1,
		aiAttackDiscardTileRisk: attackTileRisk,
	}
}

//

type gameAnalysisCache struct {
	// 局数 本场数
	wholeGameCache [][]*roundAnalysisCache

	majsoulRecordUUID string

	selfSeat int
}

func newGameAnalysisCache(majsoulRecordUUID string, selfSeat int) *gameAnalysisCache {
	cache := make([][]*roundAnalysisCache, 3*4) // 最多到西四
	for i := range cache {
		cache[i] = make([]*roundAnalysisCache, 100) // 最多连庄
	}
	return &gameAnalysisCache{
		wholeGameCache:    cache,
		majsoulRecordUUID: majsoulRecordUUID,
		selfSeat:          selfSeat,
	}
}

//

// TODO: 重构成 struct
var (
	_analysisCacheList = make([]*gameAnalysisCache, 4)
	_currentSeat       int
)

func resetAnalysisCache() {
	_analysisCacheList = make([]*gameAnalysisCache, 4)
}

func setAnalysisCache(analysisCache *gameAnalysisCache) {
	_analysisCacheList[analysisCache.selfSeat] = analysisCache
	_currentSeat = analysisCache.selfSeat
}

func getAnalysisCache(seat int) *gameAnalysisCache {
	if seat == -1 {
		return nil
	}
	return _analysisCacheList[seat]
}

func getCurrentAnalysisCache() *gameAnalysisCache {
	return getAnalysisCache(_currentSeat)
}

func (c *gameAnalysisCache) runMajsoulRecordAnalysisTask(actions majsoulRoundActions) error {
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
		if debugMode {
			fmt.Println("无需重复计算")
		}
		return nil
	}

	// 遍历自家舍牌，找到舍牌前的操作
	// 若为摸牌操作，计算出此时的 AI 进攻舍牌和防守舍牌
	// 若为鸣牌操作，计算出此时的 AI 进攻舍牌（无进攻舍牌则设为 -1），防守舍牌设为 -1
	// TODO: 玩家跳过，但是 AI 觉得应鸣牌？
	majsoulRoundData := &majsoulRoundData{selfSeat: c.selfSeat} // 注意这里是用的一个新的 majsoulRoundData 去计算的，不会有数据冲突
	majsoulRoundData.roundData = newGame(majsoulRoundData)
	majsoulRoundData.roundData.gameMode = gameModeRecordCache
	majsoulRoundData.skipOutput = true
	for i, action := range actions[:len(actions)-1] {
		if c.majsoulRecordUUID != getMajsoulCurrentRecordUUID() {
			if debugMode {
				fmt.Println("用户退出该牌谱")
			}
			// 提前退出，减少不必要的计算
			return nil
		}
		if debugMode {
			fmt.Println("助手正在计算推荐舍牌…… action", i)
		}
		majsoulRoundData.msg = action.Action
		majsoulRoundData.analysis()
	}
	roundCache.isEnd = true

	if c.majsoulRecordUUID != getMajsoulCurrentRecordUUID() {
		if debugMode {
			fmt.Println("用户退出该牌谱")
		}
		return nil
	}

	clearConsole()
	roundCache.print()

	return nil
}
