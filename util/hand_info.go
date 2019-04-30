package util

import "github.com/EndlessCheng/mahjong-helper/util/model"

type _handInfo struct {
	HandTiles34   []int
	Melds         []model.Meld // 副露
	IsTsumo       bool         // 是否自摸
	WinTile       int          // 自摸/荣和的牌
	RoundWindTile int          // 场风
	SelfWindTile  int          // 自风
	DoraCount     int          // 宝牌个数
	IsParent      bool         // 是否为亲家
	IsRiichi      bool         // 是否立直
}

type HandInfo struct {
	*_handInfo
	Divide        *DivideResult // 手牌解析结果
	_containHonor *bool
	_isNaki       *bool
}
