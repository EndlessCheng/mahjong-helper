package model

type PlayerInfo struct {
	HandTiles34   []int
	Melds         []Meld // 副露
	IsTsumo       bool   // 是否自摸
	WinTile       int    // 自摸/荣和的牌
	RoundWindTile int    // 场风
	SelfWindTile  int    // 自风
	DoraCount     int    // 宝牌个数
	IsParent      bool   // 是否为亲家
	IsDaburii     bool   // 是否双立直
	IsRiichi      bool   // 是否立直
	DiscardTiles  []int  // 注意初始化的时候把负数调整成正的！
	LeftTiles34   []int  // 剩余牌
}

func NewSimplePlayerInfo(tiles34 []int, melds []Meld) *PlayerInfo {
	return &PlayerInfo{
		HandTiles34:   tiles34,
		Melds:         melds,
		RoundWindTile: 27,
		SelfWindTile:  27,
		LeftTiles34:   InitLeftTiles34WithTiles34(tiles34),
	}
}

func (pi *PlayerInfo) FillLeftTiles34() {
	pi.LeftTiles34 = InitLeftTiles34WithTiles34(pi.HandTiles34)
}
