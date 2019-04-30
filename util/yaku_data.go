package util

type YakuType int

const (
	// https://en.wikipedia.org/wiki/Japanese_Mahjong_yaku
	// Special criteria
	YakuRiichi YakuType = iota
	YakuChiitoi

	// Yaku based on luck
	YakuTsumo
	YakuIppatsu
	YakuHaitei
	YakuHoutei
	YakuRinshan
	YakuChankan
	YakuDaburii

	// Yaku based on sequences
	YakuPinfu
	YakuRyanpeikou
	YakuIipeikou
	YakuSanshokuDoujun  // *
	YakuIttsuu          // *

	// Yaku based on triplets and/or quads
	YakuToitoi
	YakuSanAnkou
	YakuSanshokuDoukou
	YakuSanKantsu

	// Yaku based on terminal or honor tiles
	YakuTanyao
	YakuYakuhai
	YakuChanta     // * 必须有顺子
	YakuJunchan    // * 必须有顺子
	YakuHonroutou  // 七对也算
	YakuShousangen

	// Yaku based on suits
	YakuHonitsu   // *
	YakuChinitsu  // *

	// TODO: 役满
)

var YakuHanMap = map[YakuType]int{
	YakuRiichi:  1,
	YakuChiitoi: 2,

	YakuTsumo:   1,
	YakuIppatsu: 1,
	YakuHaitei:  1,
	YakuHoutei:  1,
	YakuRinshan: 1,
	YakuChankan: 1,
	YakuDaburii: 2,

	YakuPinfu:          1,
	YakuRyanpeikou:     3,
	YakuIipeikou:       1,
	YakuSanshokuDoujun: 2,
	YakuIttsuu:         2,

	YakuToitoi:         2,
	YakuSanAnkou:       2,
	YakuSanshokuDoukou: 2,
	YakuSanKantsu:      2,

	YakuTanyao:     1,
	YakuYakuhai:    1,
	YakuChanta:     2,
	YakuJunchan:    3,
	YakuHonroutou:  2,
	YakuShousangen: 2,

	YakuHonitsu:  3,
	YakuChinitsu: 6,
}

var YakumanTimesMap = map[YakuType]int{
	// TODO
}
