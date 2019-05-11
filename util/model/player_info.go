package model

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

	AvgUraDora float64 // 平均里宝牌个数，用于计算立直时的打点
}

func NewSimplePlayerInfo(tiles34 []int, melds []Meld) *PlayerInfo {
	return &PlayerInfo{
		HandTiles34:   tiles34,
		Melds:         melds,
		NumRedFives:   []int{0, 0, 0},
		RoundWindTile: 27,
		SelfWindTile:  27,
		LeftTiles34:   InitLeftTiles34WithTiles34(tiles34),
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
	return
}

// 是否已鸣牌（暗杠不算）
// 可以用来判断该玩家能否立直，计算门清加符等
func (pi *PlayerInfo) IsNaki() bool {
	for _, meld := range pi.Melds {
		if meld.MeldType != MeldTypeAnkan {
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
	return tile < 27 && tile%9 == 4 && pi.HandTiles34[tile] == pi.NumRedFives[tile/9]
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
