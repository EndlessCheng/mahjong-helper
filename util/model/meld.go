package model

const (
	MeldTypeChi    = iota // 吃
	MeldTypePon           // 碰
	MeldTypeAnkan         // 暗杠
	MeldTypeMinkan        // 大明杠
	MeldTypeKakan         // 加杠
)

type Meld struct {
	MeldType   int   // 鸣牌类型（吃、碰、暗杠、大明杠、加杠）
	Tiles      []int // 副露的牌 = sort(SelfTiles + CalledTile)
	SelfTiles  []int // 手牌中组成副露的牌（用于鸣牌分析）
	CalledTile int   // 被鸣的牌
	// TODO: 重构
	ContainRedFive    bool // 是否包含赤5
	RedFiveFromOthers bool // 赤5是否来自他家（用于获取宝牌数）（待重构）
}
