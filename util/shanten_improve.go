package util

import (
	"fmt"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"math"
	"sort"
)

// map[改良牌]进张（选择进张数最大的）
type Improves map[int]Waits

// 3k+1 张手牌的分析结果
type Hand13AnalysisResult struct {
	// 原手牌
	Tiles34 []int

	// 剩余牌
	LeftTiles34 []int

	// 是否已鸣牌（非门清状态）
	// 用于判断是否无役等
	IsNaki bool

	// 向听数
	Shanten int

	// 进张
	// 考虑了剩余枚数
	// 若某个进张牌 4 枚都可见，则该进张的 value 值为 0
	Waits Waits

	// 默听时的进张
	DamaWaits Waits

	// TODO: 鸣牌进张：他家打出这张牌，可以鸣牌，且能让向听数前进
	//MeldWaits Waits

	// map[进张牌]向听前进后的(最大)进张数
	NextShantenWaitsCountMap map[int]int

	// 向听前进后的(最大)进张数的加权均值
	AvgNextShantenWaitsCount float64

	// 综合了进张与向听前进后进张的评分
	MixedWaitsScore float64

	// 改良：摸到这张牌虽不能让向听数前进，但可以让进张变多
	// len(Improves) 即为改良的牌的种数
	Improves Improves

	// 改良情况数，这里计算的是有多少种使进张增加的摸牌-切牌方式
	ImproveWayCount int

	// 摸到非进张牌时的进张数的加权均值（非改良+改良。对于非改良牌，其进张数为 Waits.AllCount()）
	// 这里只考虑一巡的改良均值
	// TODO: 在考虑改良的情况下，如何计算向听前进所需要的摸牌次数的期望值？蒙特卡罗方法？
	AvgImproveWaitsCount float64

	// 听牌时的手牌和率
	// TODO: 未听牌时的和率？
	AvgAgariRate float64

	// 振听可能率（一向听和听牌时）
	FuritenRate float64

	// 役种
	YakuTypes map[int]struct{}

	// （鸣牌时）是否片听
	IsPartWait bool

	// 宝牌个数（手牌+副露）
	DoraCount int

	// 非立直状态下的打点期望（副露或默听）
	DamaPoint float64

	// 立直状态下的打点期望
	RiichiPoint float64

	// 局收支
	MixedRoundPoint float64

	// TODO: 赤牌改良提醒
}

// 进张和向听前进后进张的评分
// 这里粗略地近似为向听前进两次的概率
func (r *Hand13AnalysisResult) speedScore() float64 {
	if r.Waits.AllCount() == 0 || r.AvgNextShantenWaitsCount == 0 {
		return 0
	}
	leftCount := float64(CountOfTiles34(r.LeftTiles34))
	p2 := float64(r.Waits.AllCount()) / leftCount
	//p2 := r.AvgImproveWaitsCount / leftCount
	p1 := r.AvgNextShantenWaitsCount / leftCount
	//if r.AvgAgariRate > 0 { // TODO: 用和率需要考虑巡目
	//	p1 = r.AvgAgariRate / 100
	//}
	p2_, p1_ := 1-p2, 1-p1
	const leftTurns = 10.0 // math.Max(5.0, leftCount/4)
	sumP2 := p2_ * (1 - math.Pow(p2_, leftTurns)) / p2
	sumP1 := p1_ * (1 - math.Pow(p1_, leftTurns)) / p1
	result := p2 * p1 * (sumP2 - sumP1) / (p2_ - p1_)
	return result * 100
}

func (r *Hand13AnalysisResult) mixedRoundPoint() float64 {
	const weight = -1500
	if r.RiichiPoint > 0 {
		return r.AvgAgariRate/100*(r.RiichiPoint+1500) + weight
	}
	return r.AvgAgariRate/100*(r.DamaPoint+1500) + weight
}

// 调试用
func (r *Hand13AnalysisResult) String() string {
	s := fmt.Sprintf("%d 进张 %s\n%.2f 改良进张 [%d(%d) 种]",
		r.Waits.AllCount(),
		//r.Waits.AllCount()+r.MeldWaits.AllCount(),
		TilesToStrWithBracket(r.Waits.indexes()),
		r.AvgImproveWaitsCount,
		len(r.Improves),
		r.ImproveWayCount,
	)
	if len(r.DamaWaits) > 0 {
		s += fmt.Sprintf("（默听进张 %s）", TilesToStrWithBracket(r.DamaWaits.indexes()))
	}
	if r.Shanten >= 1 {
		mixedScore := r.MixedWaitsScore
		//for i := 2; i <= r.Shanten; i++ {
		//	mixedScore /= 4
		//}
		s += fmt.Sprintf(" %.2f %s进张（%.2f 综合分）",
			r.AvgNextShantenWaitsCount,
			NumberToChineseShanten(r.Shanten-1),
			mixedScore,
		)
	}
	if r.AvgAgariRate > 0 {
		s += fmt.Sprintf("[%.2f%% 和率] ", r.AvgAgariRate)
	}
	if r.MixedRoundPoint > 0 {
		s += fmt.Sprintf(" [局收支%d]", int(math.Round(r.MixedRoundPoint)))
	}
	if r.DamaPoint > 0 {
		s += fmt.Sprintf("[默听%d]", int(math.Round(r.DamaPoint)))
	}
	if r.RiichiPoint > 0 {
		s += fmt.Sprintf("[立直%d]", int(math.Round(r.RiichiPoint)))
	}
	if r.Shanten >= 0 && r.Shanten <= 1 {
		if r.FuritenRate > 0 {
			if r.FuritenRate < 1 {
				s += "[可能振听]"
			} else {
				s += "[振听]"
			}
		}
	}
	if len(r.YakuTypes) > 0 {
		s += YakuTypesWithDoraToStr(r.YakuTypes, r.DoraCount)
	}
	return s
}

func (n *shantenSearchNode13) analysis(playerInfo *model.PlayerInfo, considerImprove bool) (result13 *Hand13AnalysisResult) {
	tiles34 := playerInfo.HandTiles34
	leftTiles34 := playerInfo.LeftTiles34
	shanten13 := n.shanten
	waits := n.waits
	waitsCount := waits.AllCount()

	nextShantenWaitsCountMap := map[int]int{} // map[进张牌]听多少张牌
	improves := Improves{}
	improveWayCount := 0
	// 对于每张牌，摸到之后的手牌进张数（如果摸到的是 waits 中的牌，则进张数视作 waitsCount）
	maxImproveWaitsCount34 := make([]int, 34)
	for i := 0; i < 34; i++ {
		maxImproveWaitsCount34[i] = waitsCount // 初始化成基本进张
	}
	avgRoundPoint := 0.0
	roundPointWeight := 0
	yakuTypes := map[int]struct{}{}

	for i := 0; i < 34; i++ {
		// 从剩余牌中摸牌
		if leftTiles34[i] == 0 {
			continue
		}
		leftTiles34[i]--
		tiles34[i]++

		if node14, ok := n.children[i]; ok && node14 != nil { // 摸到的是进张
			// 计算最大向听前进后的进张
			maxNextShantenWaitsCount := 0
			for _, node13 := range node14.children {
				maxNextShantenWaitsCount = MaxInt(maxNextShantenWaitsCount, node13.waits.AllCount())
			}
			nextShantenWaitsCountMap[i] = maxNextShantenWaitsCount

			//const minRoundPoint = -1e10
			//maxRoundPoint := minRoundPoint

			if results14 := node14.analysis(playerInfo, false); len(results14) > 0 {
				bestResult14 := results14[0]

				// 加权：进张牌的剩余枚数*局收支
				w := leftTiles34[i] + 1
				avgRoundPoint += float64(w) * bestResult14.Result13.MixedRoundPoint
				roundPointWeight += w

				// 添加役种
				for t := range bestResult14.Result13.YakuTypes {
					yakuTypes[t] = struct{}{}
				}
			}

			//for discardTile, node13 := range node14.children {
			//
			//
			//	// 切牌，然后分析 3k+1 张牌下的手牌情况
			//	// 若这张是5，在只有赤5的情况下才会切赤5（TODO: 考虑赤5骗37）
			//	_isRedFive := playerInfo.IsOnlyRedFive(discardTile)
			//	playerInfo.DiscardTile(discardTile, _isRedFive)
			//
			//	// 听牌了
			//	if newShanten13 == 0 {
			//		// 听牌一般切局收支最高的，这里若为副露状态用副露局收支，否则用立直局收支
			//		_avgAgariRate := CalculateAvgAgariRate(newWaits, playerInfo) / 100
			//		var _roundPoint float64
			//		if isNaki {
			//			// FIXME: 后附时，应该只计算役牌的和率
			//			_avgPoint, _ := CalcAvgPoint(*playerInfo, newWaits)
			//			if _avgPoint == 0 { // 无役
			//				_avgAgariRate = 0
			//			}
			//			_roundPoint = _avgAgariRate*(_avgPoint+1500) - 1500
			//		} else {
			//			_avgRiichiPoint, _ := CalcAvgRiichiPoint(*playerInfo, newWaits)
			//			_roundPoint = _avgAgariRate*(_avgRiichiPoint+1500) - 1500
			//		}
			//		maxRoundPoint = math.Max(maxRoundPoint, _roundPoint)
			//		// 计算可能的役种
			//		//fillYakuTypes(newShanten13, newWaits)
			//	}
			//
			//	playerInfo.UndoDiscardTile(discardTile, _isRedFive)
			//}
			//// 加权：进张牌的剩余枚数*局收支
			//w := leftTiles34[i] + 1
			////avgAgariRate += maxAgariRate * float64(w)
			//if maxRoundPoint > minRoundPoint {
			//	avgRoundPoint += float64(w) * maxRoundPoint
			//	roundPointWeight += w
			//}
			//fmt.Println(i, maxAvgRiichiRonPoint)
			//avgRiichiPoint += maxAvgRiichiRonPoint * float64(w)
		} else if considerImprove { // 摸到的不是进张，但可能有改良
			for j := 0; j < 34; j++ {
				if tiles34[j] == 0 || j == i {
					continue
				}
				// 切牌，然后分析 3k+1 张牌下的手牌情况
				// 若这张是5，在只有赤5的情况下才会切赤5（TODO: 考虑赤5骗37）
				_isRedFive := playerInfo.IsOnlyRedFive(j)
				playerInfo.DiscardTile(j, _isRedFive)
				// 正确的切牌
				if newShanten13, improveWaits := CalculateShantenAndWaits13(tiles34, leftTiles34); newShanten13 == shanten13 {
					// 若进张数变多，则为改良
					// TODO: 若打点上升，也算改良
					if improveWaitsCount := improveWaits.AllCount(); improveWaitsCount > waitsCount {
						improveWayCount++
						if improveWaitsCount > maxImproveWaitsCount34[i] {
							maxImproveWaitsCount34[i] = improveWaitsCount
							// improves 选的是进张数最大的改良
							improves[i] = improveWaits
						}
						//fmt.Println(fmt.Sprintf("    摸 %s 切 %s 改良:", MahjongZH[i], MahjongZH[j]), improveWaitsCount, TilesToStrWithBracket(improveWaits.indexes()))
					}
				}
				playerInfo.UndoDiscardTile(j, _isRedFive)
			}
		}

		tiles34[i]--
		leftTiles34[i]++
	}

	_tiles34 := make([]int, 34)
	copy(_tiles34, tiles34)
	result13 = &Hand13AnalysisResult{
		Tiles34:                  _tiles34,
		LeftTiles34:              leftTiles34,
		IsNaki:                   playerInfo.IsNaki(),
		Shanten:                  shanten13,
		Waits:                    waits,
		DamaWaits:                Waits{},
		NextShantenWaitsCountMap: nextShantenWaitsCountMap,
		Improves:                 improves,
		ImproveWayCount:          improveWayCount,
		AvgImproveWaitsCount:     float64(waitsCount),
		YakuTypes:                yakuTypes,
		DoraCount:                playerInfo.CountDora(),
	}

	// 计算局收支、打点、和率和役种
	if waitsCount > 0 {
		//avgAgariRate /= float64(waitsCount)
		if roundPointWeight > 0 {
			avgRoundPoint /= float64(roundPointWeight)
			//if shanten13 == 1 {
			//	avgRoundPoint /= 6 // TODO: 待调整？
			//} else if shanten13 == 2 {
			//	avgRoundPoint /= 18 // TODO: 待调整？
			//}
		}
		//avgRiichiPoint /= float64(waitsCount)
		if shanten13 == shantenStateTenpai {
			// TODO: 考虑默听时的自摸
			avgRonPoint, pointResults := CalcAvgPoint(*playerInfo, waits)
			result13.DamaPoint = avgRonPoint
			// 计算默听进张
			for _, pr := range pointResults {
				result13.DamaWaits[pr.winTile] = leftTiles34[pr.winTile]
			}

			if !result13.IsNaki {
				avgRiichiPoint, riichiPointResults := CalcAvgRiichiPoint(*playerInfo, waits)
				result13.RiichiPoint = avgRiichiPoint
				result13.AvgAgariRate = CalculateAvgAgariRate(waits, playerInfo)
				for _, pr := range riichiPointResults {
					for _, yakuType := range pr.yakuTypes {
						result13.YakuTypes[yakuType] = struct{}{}
					}
				}
			} else {
				// 副露时，考虑到存在某些侍牌无法和牌（如后附、片听），不计算这些侍牌的和率
				agariRate := 0.0
				for _, pr := range pointResults { // pointResults 不包含无法和牌的情况
					agariRate = agariRate + pr.agariRate - agariRate*pr.agariRate/100
					for _, yakuType := range pr.yakuTypes {
						result13.YakuTypes[yakuType] = struct{}{}
					}
				}
				result13.AvgAgariRate = agariRate

				// 是否片听
				result13.IsPartWait = len(pointResults) < len(waits.AvailableTiles())
			}
		}
	}

	// 三向听七对子特殊提醒
	if len(playerInfo.Melds) == 0 && shanten13 == 3 && CountPairsOfTiles34(tiles34)+shanten13 == 6 {
		// 对于三向听，除非进张很差才会考虑七对子
		if waitsCount <= 21 {
			result13.YakuTypes[YakuChiitoi] = struct{}{}
		}
	}

	// 对于听牌及一向听，判断是否有振听可能
	if shanten13 <= 1 {
		for _, discardTile := range playerInfo.DiscardTiles {
			if _, ok := waits[discardTile]; ok {
				result13.FuritenRate = 0.5 // TODO: 待完善
				if shanten13 == shantenStateTenpai {
					result13.FuritenRate = 1
				}
			}
		}
	}

	// 计算局收支
	//if shanten13 <= 1 {
	//result13.DamaPoint = avgRonPoint
	//if !result13.IsNaki {
	//	result13.RiichiPoint = avgRiichiPoint
	//}
	// 振听时若能立直则只考虑立直
	//if result13.FuritenRate == 1 && result13.RiichiPoint > 0 {
	//	result13.DamaPoint = 0
	//}
	if shanten13 == shantenStateTenpai {
		result13.MixedRoundPoint = result13.mixedRoundPoint()
	} else {
		result13.MixedRoundPoint = avgRoundPoint
	}
	//}

	// 计算手牌速度
	if len(nextShantenWaitsCountMap) > 0 {
		nextShantenWaitsSum := 0
		weight := 0
		for tile, c := range nextShantenWaitsCountMap {
			w := leftTiles34[tile]
			nextShantenWaitsSum += w * c
			weight += w
		}
		result13.AvgNextShantenWaitsCount = float64(nextShantenWaitsSum) / float64(weight)
	}
	if len(improves) > 0 {
		improveWaitsSum := 0
		weight := 0
		for i := 0; i < 34; i++ {
			w := leftTiles34[i]
			improveWaitsSum += w * maxImproveWaitsCount34[i]
			weight += w
		}
		result13.AvgImproveWaitsCount = float64(improveWaitsSum) / float64(weight)
	}
	result13.MixedWaitsScore = result13.speedScore()

	// 特殊处理，方便提示向听倒退！
	if shanten13 == 2 {
		result13.MixedWaitsScore /= 4 // TODO: 待调整
	}

	return
}

func _stopShanten(shanten int) int {
	if shanten >= 3 {
		return shanten - 1
	}
	return shanten - 2
}

// 3k+1 张牌，计算向听数、进张、改良等（考虑了剩余枚数）
func CalculateShantenWithImproves13(playerInfo *model.PlayerInfo) (r *Hand13AnalysisResult) {
	if len(playerInfo.LeftTiles34) == 0 {
		playerInfo.FillLeftTiles34()
	}

	shanten := CalculateShanten(playerInfo.HandTiles34)
	shantenSearchRoot := _search13(shanten, playerInfo, _stopShanten(shanten))
	return shantenSearchRoot.analysis(playerInfo, true)
}

//

const (
	honorRiskRoundWind = 4
	honorRiskYaku      = 3
	honorRiskOtakaze   = 2
	honorRiskSelfWind  = 1
)

type tileValue float64

const (
	doraValue                tileValue = 10000
	doraFirstNeighbourValue  tileValue = 1000
	doraSecondNeighbourValue tileValue = 100
	honoredValue             tileValue = 15
)

func calculateIsolatedTileValue(tile int, playerInfo *model.PlayerInfo) tileValue {
	value := tileValue(100)

	// 是否为宝牌
	for _, doraTile := range playerInfo.DoraTiles {
		if tile == doraTile {
			value += doraValue
			//} else if doraTile < 27 {
			//	if tile/3 != doraTile/3 {
			//		continue
			//	}
			//	t9 := tile % 9
			//	dt9 := doraTile % 9
			//	if t9+1 == dt9 || t9-1 == dt9 {
			//		value += doraFirstNeighbourValue
			//	} else if t9+2 == dt9 || t9-2 == dt9 {
			//		value += doraSecondNeighbourValue
			//	}
		}
	}

	if tile >= 27 {
		if tile == playerInfo.SelfWindTile || tile == playerInfo.RoundWindTile || tile >= 31 {
			// 役牌
			value += honoredValue
			if playerInfo.SelfWindTile == playerInfo.RoundWindTile && tile == playerInfo.SelfWindTile {
				value += honoredValue // 连风
			} else if tile == playerInfo.SelfWindTile {
				value++ // 自风 +1
			} else if tile == playerInfo.RoundWindTile {
				value-- // 场风 -1
			}
			if tile == 31 {
				value -= 0.1
			}
			if tile == 32 {
				value -= 0.2
			}
		} else {
			// 客风
			for i := 1; i <= 3; i++ {
				otakazeTile := playerInfo.SelfWindTile + i
				if otakazeTile > 30 {
					otakazeTile -= 4
				}
				if tile == otakazeTile {
					// 下家 -3  对家 -2  上家 -1
					value -= tileValue(4 - i)
					break
				}
			}
		}
		left := playerInfo.LeftTiles34[tile]
		if left == 2 {
			value *= 0.9
		} else if left == 1 {
			value *= 0.2
		} else if left == 0 {
			value = 0
		}
	}

	return value
}

func calculateTileValue(tile int, playerInfo *model.PlayerInfo) (value tileValue) {
	// 是否为宝牌或宝牌周边
	for _, doraTile := range playerInfo.DoraTiles {
		if tile == doraTile {
			value += doraValue
		} else if doraTile < 27 {
			if tile/3 != doraTile/3 {
				continue
			}
			t9 := tile % 9
			dt9 := doraTile % 9
			if t9+1 == dt9 || t9-1 == dt9 {
				value += doraFirstNeighbourValue
			} else if t9+2 == dt9 || t9-2 == dt9 {
				value += doraSecondNeighbourValue
			}
		}
	}
	return
}

type Hand14AnalysisResult struct {
	// 需要切的牌
	DiscardTile int

	// 切的是否为宝牌
	IsDiscardDoraTile bool

	// 切的牌的价值（宝牌或宝牌周边）
	DiscardTileValue tileValue

	// 切的牌是否为幺九浮牌
	isIsolatedYaochuDiscardTile bool

	// 切牌后的手牌分析结果
	Result13 *Hand13AnalysisResult

	DiscardHonorTileRisk int

	// 剩余可以摸的牌数
	LeftDrawTilesCount int

	// 副露信息（没有副露就是 nil）
	// 比如用 23m 吃了牌，OpenTiles 就是 [1,2]
	OpenTiles []int
}

func (r *Hand14AnalysisResult) String() string {
	meldInfo := ""
	if len(r.OpenTiles) > 0 {
		meldType := "吃"
		if r.OpenTiles[0] == r.OpenTiles[1] {
			meldType = "碰"
		}
		meldInfo = fmt.Sprintf("用 %s%s %s，", string([]rune(MahjongZH[r.OpenTiles[0]])[:1]), MahjongZH[r.OpenTiles[1]], meldType)
	}
	return meldInfo + fmt.Sprintf("切 %s: %s", MahjongZH[r.DiscardTile], r.Result13.String())
}

type Hand14AnalysisResultList []*Hand14AnalysisResult

// 按照特定规则排序
// 若 improveFirst 为 true，则优先按照 AvgImproveWaitsCount 排序（对于三向听及以上来说）
func (l Hand14AnalysisResultList) Sort(improveFirst bool) {
	if len(l) <= 1 {
		return
	}

	shanten := l[0].Result13.Shanten

	sort.Slice(l, func(i, j int) bool {
		ri, rj := l[i].Result13, l[j].Result13
		riWaitsCount, rjWaitsCount := ri.Waits.AllCount(), rj.Waits.AllCount()

		// 首先，无论怎样，进张数为 0，无条件排在后面，也不看改良
		// 进张数都为 0 才看改良
		if riWaitsCount == 0 || rjWaitsCount == 0 {
			if riWaitsCount == 0 && rjWaitsCount == 0 {
				return ri.AvgImproveWaitsCount > rj.AvgImproveWaitsCount
			}
			return riWaitsCount > rjWaitsCount
		}

		switch shanten {
		case 0:
			// 听牌的话：局收支 - 和率
			// 局收支，有明显差异
			if !InDelta(ri.MixedRoundPoint, rj.MixedRoundPoint, 100) {
				return ri.MixedRoundPoint > rj.MixedRoundPoint
			}
			// 和率优先
			if !Equal(ri.AvgAgariRate, rj.AvgAgariRate) {
				return ri.AvgAgariRate > rj.AvgAgariRate
			}
		case 1:
			// 一向听：进张*局收支
			var riScore, rjScore float64
			if shanten >= 2 && improveFirst {
				// 对于两向听，若需要改良的话以改良为主
				//riScore = float64(ri.AvgImproveWaitsCount) * ri.MixedRoundPoint
				//rjScore = float64(rj.AvgImproveWaitsCount) * rj.MixedRoundPoint
				break
			} else {
				// 负数要调整
				wi := float64(riWaitsCount)
				if ri.MixedRoundPoint < 0 {
					wi = 1 / wi
				}
				wj := float64(rjWaitsCount)
				if rj.MixedRoundPoint < 0 {
					wj = 1 / wj
				}
				riScore = wi * ri.MixedRoundPoint
				rjScore = wj * rj.MixedRoundPoint
			}
			if !Equal(riScore, rjScore) {
				return riScore > rjScore
			}
		}

		if shanten >= 2 {
			// 两向听及以上时，若存在幺九浮牌，则根据价值来单独比较浮牌
			if l[i].isIsolatedYaochuDiscardTile && l[j].isIsolatedYaochuDiscardTile {
				// 优先切掉价值最低的浮牌，这里直接比较浮点数
				if l[i].DiscardTileValue != l[j].DiscardTileValue {
					return l[i].DiscardTileValue < l[j].DiscardTileValue
				}
			} else if l[i].isIsolatedYaochuDiscardTile && l[i].DiscardTileValue < 500 {
				return true
			} else if l[j].isIsolatedYaochuDiscardTile && l[j].DiscardTileValue < 500 {
				return false
			}
		}

		//if improveFirst {
		//	// 优先按照 AvgImproveWaitsCount 排序
		//	if !Equal(ri.AvgImproveWaitsCount, rj.AvgImproveWaitsCount) {
		//		return ri.AvgImproveWaitsCount > rj.AvgImproveWaitsCount
		//	}
		//}

		// 排序规则：综合评分（速度） - 进张 - 前进后的进张 - 和率 - 改良 - 价值低 - 好牌先走
		// 必须注意到的一点是，随着游戏的进行，进张会被他家打出，所以进张是有减少的趋势的
		// 对于一向听，考虑到未听牌之前要听的牌会被他家打出而造成听牌时的枚数降低，所以听牌枚数比和率更重要
		// 对比当前进张与前进后的进张，在二者综合评分相近的情况下（注意这个前提），由于进张越多听牌速度越快，听牌时的进张数也就越接近预期进张数，所以进张越多越好（再次强调是在二者综合评分相近的情况下）

		if !Equal(ri.MixedWaitsScore, rj.MixedWaitsScore) {
			return ri.MixedWaitsScore > rj.MixedWaitsScore
		}

		if riWaitsCount != rjWaitsCount {
			return riWaitsCount > rjWaitsCount
		}

		if !Equal(ri.AvgNextShantenWaitsCount, rj.AvgNextShantenWaitsCount) {
			return ri.AvgNextShantenWaitsCount > rj.AvgNextShantenWaitsCount
		}

		// shanten == 1
		if !Equal(ri.AvgAgariRate, rj.AvgAgariRate) {
			return ri.AvgAgariRate > rj.AvgAgariRate
		}

		if !Equal(ri.AvgImproveWaitsCount, rj.AvgImproveWaitsCount) {
			return ri.AvgImproveWaitsCount > rj.AvgImproveWaitsCount
		}

		if l[i].DiscardTileValue != l[j].DiscardTileValue {
			// 价值低的先走
			return l[i].DiscardTileValue < l[j].DiscardTileValue
		}

		// 好牌先走
		idxI, idxJ := l[i].DiscardTile, l[j].DiscardTile
		if idxI < 27 && idxJ < 27 {
			idxI %= 9
			if idxI > 4 {
				idxI = 8 - idxI
			}
			idxJ %= 9
			if idxJ > 4 {
				idxJ = 8 - idxJ
			}
			return idxI > idxJ
		}
		if idxI < 27 || idxJ < 27 {
			// 数牌先走
			return idxI < idxJ
		}
		// 场风 - 三元牌 - 他家客风 - 自风
		return l[i].DiscardHonorTileRisk > l[j].DiscardHonorTileRisk

		//// 改良种类、方式多的优先
		//if len(ri.Improves) != len(rj.Improves) {
		//	return len(ri.Improves) > len(rj.Improves)
		//}
		//if ri.ImproveWayCount != rj.ImproveWayCount {
		//	return ri.ImproveWayCount > rj.ImproveWayCount
		//}
	})
}

func (l *Hand14AnalysisResultList) filterOutDiscard(cantDiscardTile int) {
	newResults := Hand14AnalysisResultList{}
	for _, r := range *l {
		if r.DiscardTile != cantDiscardTile {
			newResults = append(newResults, r)
		}
	}
	*l = newResults
}

func (l Hand14AnalysisResultList) addOpenTile(openTiles []int) {
	for _, r := range l {
		r.OpenTiles = openTiles
	}
}

func (n *shantenSearchNode14) analysis(playerInfo *model.PlayerInfo, considerImprove bool) (results Hand14AnalysisResultList) {
	for discardTile, node13 := range n.children {
		isRedFive := playerInfo.IsOnlyRedFive(discardTile)

		// 切牌，然后分析 3k+1 张牌下的手牌情况
		// 若这张是5，在只有赤5的情况下才会切赤5（TODO: 考虑赤5骗37）
		playerInfo.DiscardTile(discardTile, isRedFive)
		result13 := node13.analysis(playerInfo, considerImprove)

		// 记录切牌后的分析结果
		r14 := &Hand14AnalysisResult{
			DiscardTile:        discardTile,
			IsDiscardDoraTile:  InInts(discardTile, playerInfo.DoraTiles),
			Result13:           result13,
			LeftDrawTilesCount: playerInfo.LeftDrawTilesCount,
		}
		results = append(results, r14)

		if n.shanten >= 2 {
			if isYaochupai(discardTile) && isIsolatedTile(discardTile, playerInfo.HandTiles34) {
				r14.isIsolatedYaochuDiscardTile = true
				r14.DiscardTileValue = calculateIsolatedTileValue(discardTile, playerInfo)
			} else {
				r14.DiscardTileValue = calculateTileValue(discardTile, playerInfo)
			}
		}

		if discardTile >= 27 {
			switch discardTile {
			case playerInfo.RoundWindTile:
				r14.DiscardHonorTileRisk = honorRiskRoundWind
			case 31, 32, 33:
				r14.DiscardHonorTileRisk = honorRiskYaku
			case playerInfo.SelfWindTile:
				r14.DiscardHonorTileRisk = honorRiskSelfWind
			default:
				r14.DiscardHonorTileRisk = honorRiskOtakaze
			}
		}

		playerInfo.UndoDiscardTile(discardTile, isRedFive)
	}

	// 下面这一逻辑被「综合速度」取代
	//improveFirst := func(l []*Hand14AnalysisResult) bool {
	//	if !considerImprove || len(l) <= 1 {
	//		return false
	//	}
	//
	//	shanten := l[0].Result13.Shanten
	//	// 一向听及以下着眼于进张，改良其次
	//	if shanten <= 1 {
	//		return false
	//	}
	//
	//	// 判断七对和一般型的向听数是否相同，若七对更小则改良优先
	//	tiles34 := playerInfo.HandTiles34
	//	shantenChiitoi := CalculateShantenOfChiitoi(tiles34)
	//	shantenNormal := CalculateShantenOfNormal(tiles34, CountOfTiles34(tiles34))
	//	return shantenChiitoi < shantenNormal
	//}
	//improveFst := improveFirst(results)

	results.Sort(false)

	return
}

// 3k+2 张牌，计算向听数、进张、改良、向听倒退等
func CalculateShantenWithImproves14(playerInfo *model.PlayerInfo) (shanten int, results Hand14AnalysisResultList, incShantenResults Hand14AnalysisResultList) {
	if len(playerInfo.LeftTiles34) == 0 {
		playerInfo.FillLeftTiles34()
	}

	shanten = CalculateShanten(playerInfo.HandTiles34)
	stopAtShanten := _stopShanten(shanten)
	shantenSearchRoot := searchShanten14(shanten, playerInfo, stopAtShanten)
	results = shantenSearchRoot.analysis(playerInfo, true)
	incShantenSearchRoot := searchShanten14(shanten+1, playerInfo, stopAtShanten+1)
	incShantenResults = incShantenSearchRoot.analysis(playerInfo, true)
	return
}

// 计算最小向听数，鸣牌方式
func calculateMeldShanten(tiles34 []int, calledTile int, isRedFive bool, allowChi bool) (minShanten int, meldCombinations []model.Meld) {
	// 是否能碰
	if tiles34[calledTile] >= 2 {
		meldCombinations = append(meldCombinations, model.Meld{
			MeldType:          model.MeldTypePon,
			Tiles:             []int{calledTile, calledTile, calledTile},
			SelfTiles:         []int{calledTile, calledTile},
			CalledTile:        calledTile,
			RedFiveFromOthers: isRedFive,
		})
	}
	// 是否能吃
	if allowChi && calledTile < 27 {
		checkChi := func(tileA, tileB int) {
			if tiles34[tileA] > 0 && tiles34[tileB] > 0 {
				_tiles := []int{tileA, tileB, calledTile}
				sort.Ints(_tiles)
				meldCombinations = append(meldCombinations, model.Meld{
					MeldType:          model.MeldTypeChi,
					Tiles:             _tiles,
					SelfTiles:         []int{tileA, tileB},
					CalledTile:        calledTile,
					RedFiveFromOthers: isRedFive,
				})
			}
		}
		t9 := calledTile % 9
		if t9 >= 2 {
			checkChi(calledTile-2, calledTile-1)
		}
		if t9 >= 1 && t9 <= 7 {
			checkChi(calledTile-1, calledTile+1)
		}
		if t9 <= 6 {
			checkChi(calledTile+1, calledTile+2)
		}
	}

	// 计算所有鸣牌下的最小向听数
	minShanten = 99
	for _, c := range meldCombinations {
		tiles34[c.SelfTiles[0]]--
		tiles34[c.SelfTiles[1]]--
		minShanten = MinInt(minShanten, CalculateShanten(tiles34))
		tiles34[c.SelfTiles[0]]++
		tiles34[c.SelfTiles[1]]++
	}

	return
}

// TODO 鸣牌的情况判断（待重构）
// 编程时注意他家切掉的这张牌是否算到剩余数中
//if isOpen {
//if newShanten, combinations, shantens := calculateMeldShanten(tiles34, i, true); newShanten < shanten {
//	// 向听前进了，说明鸣牌成功，则换的这张牌为鸣牌进张
//	// 计算进张数：若能碰则 =剩余数*3，否则 =剩余数
//	meldWaits[i] = leftTile - tiles34[i]
//	for i, comb := range combinations {
//		if comb[0] == comb[1] && shantens[i] == newShanten {
//			meldWaits[i] *= 3
//			break
//		}
//	}
//}
//}

// 计算鸣牌下的何切分析
// calledTile 他家出的牌，尝试鸣这张牌
// isRedFive 这张牌是否为赤5
// allowChi 是否允许吃这张牌
func CalculateMeld(playerInfo *model.PlayerInfo, calledTile int, isRedFive bool, allowChi bool) (minShanten int, results Hand14AnalysisResultList, incShantenResults Hand14AnalysisResultList) {
	if len(playerInfo.LeftTiles34) == 0 {
		playerInfo.FillLeftTiles34()
	}

	minShanten, meldCombinations := calculateMeldShanten(playerInfo.HandTiles34, calledTile, isRedFive, allowChi)

	for _, c := range meldCombinations {
		// 尝试鸣这张牌
		playerInfo.AddMeld(c)
		_shanten, _results, _incShantenResults := CalculateShantenWithImproves14(playerInfo)
		playerInfo.UndoAddMeld()

		// 去掉现物食替的情况
		_results.filterOutDiscard(calledTile)
		_incShantenResults.filterOutDiscard(calledTile)

		// 去掉筋食替的情况
		if c.MeldType == model.MeldTypeChi {
			cannotDiscardTile := -1
			if c.SelfTiles[0] < calledTile && c.SelfTiles[1] < calledTile && calledTile%9 >= 3 {
				cannotDiscardTile = calledTile - 3
			} else if c.SelfTiles[0] > calledTile && c.SelfTiles[1] > calledTile && calledTile%9 <= 5 {
				cannotDiscardTile = calledTile + 3
			}
			if cannotDiscardTile != -1 {
				_results.filterOutDiscard(cannotDiscardTile)
				_incShantenResults.filterOutDiscard(cannotDiscardTile)
			}
		}

		// 添加副露信息，用于输出
		_results.addOpenTile(c.SelfTiles)
		_incShantenResults.addOpenTile(c.SelfTiles)

		// 整理副露结果
		if _shanten == minShanten {
			results = append(results, _results...)
			incShantenResults = append(incShantenResults, _incShantenResults...)
		} else if _shanten == minShanten+1 {
			incShantenResults = append(incShantenResults, _results...)
		}
	}

	results.Sort(false)
	incShantenResults.Sort(false)

	return
}
