package model

const (
	MeldTypeChi    = iota // 吃
	MeldTypePon           // 碰
	MeldTypeAnkan         // 暗杠
	MeldTypeMinkan        // 大明杠
	MeldTypeKakan         // 加杠
)

type Meld struct {
	MeldType       int   // 鸣牌类型（吃、碰、暗杠、大明杠、加杠）
	Tiles          []int // 副露的牌 [0-33]
	CalledTile     int   // 被鸣的牌 0-33
	ContainRedFive bool  // 是否包含赤5
}
