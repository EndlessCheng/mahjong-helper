package model

const (
	MeldTypeChi    = iota // 吃
	MeldTypePon           // 碰
	MeldTypeAnkan         // 暗杠
	MeldTypeMinkan        // 大明杠
	MeldTypeKakan         // 加杠
)

type Meld struct {
	MeldType int // 鸣牌类型（吃、碰、暗杠、大明杠、加杠）

	// Tiles == sort(SelfTiles + CalledTile)
	Tiles      []int // 副露的牌
	SelfTiles  []int // 手牌中组成副露的牌（用于鸣牌分析）
	CalledTile int   // 被鸣的牌

	// TODO: 重构 ContainRedFive RedFiveFromOthers
	ContainRedFive    bool // 是否包含赤5
	RedFiveFromOthers bool // 赤5是否来自他家（用于获取宝牌数）
}

// 是否为杠子
func (m *Meld) IsKan() bool {
	return m.MeldType == MeldTypeAnkan || m.MeldType == MeldTypeMinkan || m.MeldType == MeldTypeKakan
}
