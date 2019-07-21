package model

import "fmt"

type PlayerInfo struct {
	HandTiles34 []int  // 手牌，不含副露
	Melds       []Meld // 副露
	DoraTiles   []int  // 宝牌指示牌产生的宝牌，可以重复
	NumRedFives []int  // 按照 mps 的顺序，各个赤5的个数（手牌和副露中的）

	IsTsumo       bool // 是否自摸
	WinTile       int  // 自摸/荣和的牌
	RoundWindTile int  // 场风
	SelfWindTile  int  // 自风
	IsParent      bool // 是否为亲家
	IsDaburii     bool // 是否双立直
	IsRiichi      bool // 是否立直

	DiscardTiles []int // 自家舍牌，用于判断和率，是否振听等  *注意创建 PlayerInfo 的时候把负数调整成正的！
	LeftTiles34  []int // 剩余牌

	LeftDrawTilesCount int // 剩余可以摸的牌数

	//LeftRedFives []int // 剩余赤5个数，用于估算打点
	//AvgUraDora float64 // 平均里宝牌个数，用于计算立直时的打点

	NukiDoraNum int // 拔北宝牌数
}

func NewSimplePlayerInfo(tiles34 []int, melds []Meld) *PlayerInfo {
	leftTiles34 := InitLeftTiles34WithTiles34(tiles34)
	for _, meld := range melds {
		for _, tile := range meld.Tiles {
			leftTiles34[tile]--
			if leftTiles34[tile] < 0 {
				panic(fmt.Sprint("副露数据不合法", melds))
			}
		}
	}
	return &PlayerInfo{
		HandTiles34:   tiles34,
		Melds:         melds,
		NumRedFives:   make([]int, 3),
		RoundWindTile: 27,
		SelfWindTile:  27,
		LeftTiles34:   leftTiles34,
	}
}

// 根据手牌、副露、赤5，结合哪些是宝牌，计算出拥有的宝牌个数
func (pi *PlayerInfo) CountDora() (count int) {
	for _, doraTile := range pi.DoraTiles {
		count += pi.HandTiles34[doraTile]
		for _, m := range pi.Melds {
			for _, tile := range m.Tiles {
				if tile == doraTile {
					count++
				}
			}
		}
	}
	// 手牌和副露中的赤5
	for _, num := range pi.NumRedFives {
		count += num
	}
	// 拔北宝牌
	if pi.NukiDoraNum > 0 {
		count += pi.NukiDoraNum
		// 特殊：西为指示牌
		for _, doraTile := range pi.DoraTiles {
			if doraTile == 30 {
				count += pi.NukiDoraNum
			}
		}
	}
	return
}

// 立直时，根据牌山计算和了时的里宝牌个数
// TODO: 考虑 WinTile
//func (pi *PlayerInfo) CountUraDora() (count float64) {
//	if !pi.IsRiichi || pi.IsNaki() {
//		return 0
//	}
//	uraDoraTileLeft := make([]int, len(pi.LeftTiles34))
//	for tile, left := range pi.LeftTiles34 {
//		uraDoraTileLeft[DoraTile(tile)] = left
//	}
//	sum := 0
//	weight := 0
//	for tile, c := range pi.HandTiles34 {
//		w := uraDoraTileLeft[tile]
//		sum += w * c
//		weight += w
//	}
//	for _, meld := range pi.Melds {
//		for tile, c := range meld.Tiles {
//			w := uraDoraTileLeft[tile]
//			sum += w * c
//			weight += w
//		}
//	}
//	// 简化计算，直接乘上宝牌指示牌的个数
//	return float64(len(pi.DoraTiles)*sum) / float64(weight)
//}

// 是否已鸣牌（暗杠不算）
// 可以用来判断该玩家能否立直，计算门清加符、役种番数等
func (pi *PlayerInfo) IsNaki() bool {
	for _, meld := range pi.Melds {
		if meld.MeldType != MeldTypeAnkan {
			return true
		}
	}
	return false
}

// 是否振听
// 仅限听牌时调用
// TODO: Waits 移进来
func (pi *PlayerInfo) IsFuriten(waits map[int]int) bool {
	for _, discardTile := range pi.DiscardTiles {
		if _, ok := waits[discardTile]; ok {
			return true
		}
	}
	return false
}

/************* 以下接口暂为内部调用 ************/

func (pi *PlayerInfo) FillLeftTiles34() {
	pi.LeftTiles34 = InitLeftTiles34WithTiles34(pi.HandTiles34)
}

// 手上的这种牌只有赤5
func (pi *PlayerInfo) IsOnlyRedFive(tile int) bool {
	return tile < 27 && tile%9 == 4 && pi.HandTiles34[tile] > 0 && pi.HandTiles34[tile] == pi.NumRedFives[tile/9]
}

func (pi *PlayerInfo) DiscardTile(tile int, isRedFive bool) {
	// 从手牌中舍去一张牌到牌河
	pi.HandTiles34[tile]--
	if isRedFive {
		pi.NumRedFives[tile/9]--
	}
	pi.DiscardTiles = append(pi.DiscardTiles, tile)
}

func (pi *PlayerInfo) UndoDiscardTile(tile int, isRedFive bool) {
	// 复原从手牌中舍去一张牌到牌河的动作，即把这张牌从牌河移回手牌
	pi.DiscardTiles = pi.DiscardTiles[:len(pi.DiscardTiles)-1]
	pi.HandTiles34[tile]++
	if isRedFive {
		pi.NumRedFives[tile/9]++
	}
}

//func (pi *PlayerInfo) DrawTile(tile int) {
//	// 从牌山中摸牌
//}
//
//func (pi *PlayerInfo) UndoDrawTile(tile int) {
//	// 复原从牌山中摸牌的动作，即把这张牌放回牌山
//}

func (pi *PlayerInfo) AddMeld(meld Meld) {
	// 用手牌中的牌去鸣牌
	// 原有的宝牌数量并未发生变化
	for _, tile := range meld.SelfTiles {
		pi.HandTiles34[tile]--
	}
	pi.Melds = append(pi.Melds, meld)
	if meld.RedFiveFromOthers {
		tile := meld.Tiles[0]
		pi.NumRedFives[tile/9]++
	}
}

func (pi *PlayerInfo) UndoAddMeld() {
	// 复原鸣牌动作
	latestMeld := pi.Melds[len(pi.Melds)-1]
	for _, tile := range latestMeld.SelfTiles {
		pi.HandTiles34[tile]++
	}
	pi.Melds = pi.Melds[:len(pi.Melds)-1]
	if latestMeld.RedFiveFromOthers {
		tile := latestMeld.Tiles[0]
		pi.NumRedFives[tile/9]--
	}
}
