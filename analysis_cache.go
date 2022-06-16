package main

import (
	"fmt"

	"github.com/EndlessCheng/mahjong-helper/Console"
	"github.com/EndlessCheng/mahjong-helper/util"
	"github.com/fatih/color"
)

type analysisOpType int

const (
	AnalysisOpTypeTsumo     analysisOpType = iota
	AnalysisOpTypeChiPonKan                // 吃 碰 明杠
	AnalysisOpTypeKan                      // 加杠 暗杠
)

// TODO: 提醒「此处应该副露，不应跳过」

type AnalysisCache struct {
	analysisOpType analysisOpType

	selfDiscardTile int
	//isSelfDiscardRedFive bool
	selfDiscardTileRisk float64
	isRiichiWhenDiscard bool
	meldType            int

	// 用手牌中的什么牌去鳴牌，空就是跳过不鸣
	selfOpenTiles []int

	aiAttackDiscardTile      int
	aiDefenceDiscardTile     int
	aiAttackDiscardTileRisk  float64
	aiDefenceDiscardTileRisk float64

	tenpaiRate []float64 // TODO: 三家听牌率
}

type RoundAnalysisCache struct {
	IsStart bool
	IsEnd   bool
	Cache   []*AnalysisCache

	analysisCacheBeforeChiPon *AnalysisCache
}

func (roundAnalysisCache *RoundAnalysisCache) Print() {
	const (
		BaseInfo  = "助手正在计算推荐舍牌，请稍等……（计算结果仅供参考）"
		EmptyInfo = "--"
		Sep       = "  "
	)

	done := roundAnalysisCache != nil && roundAnalysisCache.IsEnd
	if !done {
		color.HiGreen(BaseInfo)
	} else {
		// 检查最后的是否自摸，若为自摸则去掉推荐
		if len(roundAnalysisCache.Cache) > 0 {
			latestCache := roundAnalysisCache.Cache[len(roundAnalysisCache.Cache)-1]
			if latestCache.selfDiscardTile == -1 {
				latestCache.aiAttackDiscardTile = -1
				latestCache.aiDefenceDiscardTile = -1
			}
		}
	}

	fmt.Print("巡目　　")
	if done {
		for i := range roundAnalysisCache.Cache {
			fmt.Printf("%s%2d", Sep, i+1)
		}
	}
	fmt.Println()

	printTileInfo := func(tile int, risk float64, suffix string) {
		info := EmptyInfo
		if tile != -1 {
			info = util.Mahjong[tile]
		}
		fmt.Print(Sep)
		if info == EmptyInfo || risk < 5 {
			fmt.Print(info)
		} else {
			color.New(GetNumRiskColor(risk)).Print(info)
		}
		fmt.Print(suffix)
	}

	fmt.Print("自家切牌")
	if done {
		for i, c := range roundAnalysisCache.Cache {
			suffix := ""
			if c.isRiichiWhenDiscard {
				suffix = "[立直]"
			} else if c.selfDiscardTile == -1 && i == len(roundAnalysisCache.Cache)-1 {
				//suffix = "[自摸]"
				// TODO: 流局
			}
			printTileInfo(c.selfDiscardTile, c.selfDiscardTileRisk, suffix)
		}
	}
	fmt.Println()

	fmt.Print("進攻推薦")
	if done {
		for _, c := range roundAnalysisCache.Cache {
			printTileInfo(c.aiAttackDiscardTile, c.aiAttackDiscardTileRisk, "")
		}
	}
	fmt.Println()

	fmt.Print("防守推薦")
	if done {
		for _, c := range roundAnalysisCache.Cache {
			printTileInfo(c.aiDefenceDiscardTile, c.aiDefenceDiscardTileRisk, "")
		}
	}
	fmt.Println()

	fmt.Println()
}

// （摸牌后、鳴牌后的）实际舍牌
func (rc *RoundAnalysisCache) AddSelfDiscardTile(tile int, risk float64, isRiichiWhenDiscard bool) {
	latestCache := rc.Cache[len(rc.Cache)-1]
	latestCache.selfDiscardTile = tile
	latestCache.selfDiscardTileRisk = risk
	latestCache.isRiichiWhenDiscard = isRiichiWhenDiscard
}

// 摸牌时的切牌推荐
func (rc *RoundAnalysisCache) AddAIDiscardTileWhenDrawTile(attackTile int, defenceTile int, attackTileRisk float64, defenceDiscardTileRisk float64) {
	// 摸牌，巡目+1
	rc.Cache = append(rc.Cache, &AnalysisCache{
		analysisOpType:           AnalysisOpTypeTsumo,
		selfDiscardTile:          -1,
		aiAttackDiscardTile:      attackTile,
		aiDefenceDiscardTile:     defenceTile,
		aiAttackDiscardTileRisk:  attackTileRisk,
		aiDefenceDiscardTileRisk: defenceDiscardTileRisk,
	})
	rc.analysisCacheBeforeChiPon = nil
}

// 加杠 暗杠
func (rc *RoundAnalysisCache) AddKan(meldType int) {
	// latestCache 是摸牌
	latestCache := rc.Cache[len(rc.Cache)-1]
	latestCache.analysisOpType = AnalysisOpTypeKan
	latestCache.meldType = meldType
	// 杠完之后又会摸牌，巡目+1
}

// 吃 碰 明杠
func (rc *RoundAnalysisCache) AddChiPonKan(meldType int) {
	if meldType == meldTypeMinkan {
		// 暂时忽略明杠，巡目不+1，留给摸牌时+1
		return
	}
	// 巡目+1
	var newCache *AnalysisCache
	if rc.analysisCacheBeforeChiPon != nil {
		newCache = rc.analysisCacheBeforeChiPon // 见 addPossibleChiPonKan
		newCache.analysisOpType = AnalysisOpTypeChiPonKan
		newCache.meldType = meldType
		rc.analysisCacheBeforeChiPon = nil
	} else {
		// 此处代码应该不会触发
		if DebugMode {
			panic("rc.analysisCacheBeforeChiPon == nil")
		}
		newCache = &AnalysisCache{
			analysisOpType:       AnalysisOpTypeChiPonKan,
			selfDiscardTile:      -1,
			aiAttackDiscardTile:  -1,
			aiDefenceDiscardTile: -1,
			meldType:             meldType,
		}
	}
	rc.Cache = append(rc.Cache, newCache)
}

// 吃 碰 杠 跳过
func (rc *RoundAnalysisCache) AddPossibleChiPonKan(attackTile int, attackTileRisk float64) {
	rc.analysisCacheBeforeChiPon = &AnalysisCache{
		analysisOpType:          AnalysisOpTypeChiPonKan,
		selfDiscardTile:         -1,
		aiAttackDiscardTile:     attackTile,
		aiDefenceDiscardTile:    -1,
		aiAttackDiscardTileRisk: attackTileRisk,
	}
}

//

type GameAnalysisCache struct {
	// 局数 本场数
	WholeGameCache [][]*RoundAnalysisCache

	MahJongSoulRecordUUID string

	SelfSeat int
}

func newGameAnalysisCache(majsoulRecordUUID string, selfSeat int) *GameAnalysisCache {
	cache := make([][]*RoundAnalysisCache, 3*4) // 最多到西四
	for i := range cache {
		cache[i] = make([]*RoundAnalysisCache, 100) // 最多连庄
	}
	return &GameAnalysisCache{
		WholeGameCache:        cache,
		MahJongSoulRecordUUID: majsoulRecordUUID,
		SelfSeat:              selfSeat,
	}
}

//

// TODO: 重构成 struct
var (
	analysisCacheList = make([]*GameAnalysisCache, 4)
	currentSeat       int
)

func ResetAnalysisCache() {
	analysisCacheList = make([]*GameAnalysisCache, 4)
}

func SetAnalysisCache(analysisCache *GameAnalysisCache) {
	analysisCacheList[analysisCache.SelfSeat] = analysisCache
	currentSeat = analysisCache.SelfSeat
}

func GetAnalysisCache(seat int) *GameAnalysisCache {
	if seat == -1 {
		return nil
	}
	return analysisCacheList[seat]
}

func GetCurrentAnalysisCache() *GameAnalysisCache {
	return GetAnalysisCache(currentSeat)
}

func (cache *GameAnalysisCache) RunMahJongSoulRecordAnalysisTask(actions MahJongSoulRoundActions) error {
	// 从第一个 action 中取出局和场
	if len(actions) == 0 {
		return fmt.Errorf("数据异常：此局数据为空")
	}

	newRoundAction := actions[0]
	data := newRoundAction.Action
	roundNumber := 4*(*data.Chang) + *data.Ju
	ben := *data.Ben
	roundCache := cache.WholeGameCache[roundNumber][ben] // TODO: 建议用原子操作
	if roundCache == nil {
		roundCache = &RoundAnalysisCache{IsStart: true}
		if DebugMode {
			fmt.Println("助手正在计算推荐舍牌…… 创建 roundCache")
		}
		cache.WholeGameCache[roundNumber][ben] = roundCache
	} else if roundCache.IsStart {
		if DebugMode {
			fmt.Println("无需重复计算")
		}
		return nil
	}

	// 遍历自家舍牌，找到舍牌前的操作
	// 若为摸牌操作，计算出此时的 AI 进攻舍牌和防守舍牌
	// 若为鳴牌操作，计算出此时的 AI 进攻舍牌（无进攻舍牌则设为 -1），防守舍牌设为 -1
	// TODO: 玩家跳过，但是 AI 觉得应鳴牌？
	majsoulRoundData := &MahJongSoulRoundData{SelfSeat: cache.SelfSeat} // 注意这里是用的一个新的 majsoulRoundData 去计算的，不会有数据冲突
	majsoulRoundData.RoundData = NewGame(majsoulRoundData)
	majsoulRoundData.RoundData.GameMode = GameModeRecordCache
	majsoulRoundData.SkipOutput = true
	for i, action := range actions[:len(actions)-1] {
		if cache.MahJongSoulRecordUUID != GetMajsoulCurrentRecordUUID() {
			if DebugMode {
				fmt.Println("用户退出该牌谱")
			}
			// 提前退出，减少不必要的计算
			return nil
		}
		if DebugMode {
			fmt.Println("助手正在计算推荐舍牌…… action", i)
		}
		majsoulRoundData.Message = action.Action
		majsoulRoundData.Analysis()
	}
	roundCache.IsEnd = true

	if cache.MahJongSoulRecordUUID != GetMajsoulCurrentRecordUUID() {
		if DebugMode {
			fmt.Println("用户退出该牌谱")
		}
		return nil
	}

	Console.ClearScreen()
	roundCache.Print()

	return nil
}
