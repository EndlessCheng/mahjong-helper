package main

type roundData struct {
	roundNumber int

	// 场风
	roundWindTile int

	// 宝牌指示牌
	doraIndicators []int

	// 自家手牌
	counts []int

	// 牌山剩余牌量
	leftCounts []int

	// 全局舍牌
	// 按舍牌顺序，负数表示摸切(-)，非负数表示手切(+)
	// 可以理解成：- 表示不要/暗色，+ 表示进张/亮色
	globalDiscardTiles []int
	// 0=自家, 1=下家, 2=对家, 3=上家
	players [4]*playerInfo

	parser ParseDataInterface
}

type ParseDataInterface interface {
	ParseInit()
	ParseDraw()
	ParseDiscard()
}

func (d *roundData) analysis() error {
	return nil
}
