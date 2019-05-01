package util

import "fmt"

const (
	// https://en.wikipedia.org/wiki/Japanese_Mahjong_yaku
	// Special criteria
	YakuRiichi int = iota
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

const maxYakuType = 64

var YakuNameMap = map[int]string{
	// Special criteria
	YakuRiichi:  "立直",
	YakuChiitoi: "七对",

	// Yaku based on luck
	YakuTsumo:   "自摸",
	YakuIppatsu: "一发",
	YakuHaitei:  "海底",
	YakuHoutei:  "河底",
	YakuRinshan: "岭上",
	YakuChankan: "抢杠",
	YakuDaburii: "w立",

	// Yaku based on sequences
	YakuPinfu:          "平和",
	YakuRyanpeikou:     "两杯口",
	YakuIipeikou:       "一杯口",
	YakuSanshokuDoujun: "三色",
	YakuIttsuu:         "一气",

	// Yaku based on triplets and/or quads
	YakuToitoi:         "对对",
	YakuSanAnkou:       "三暗刻",
	YakuSanshokuDoukou: "三色同刻",
	YakuSanKantsu:      "三杠子",

	// Yaku based on terminal or honor tiles
	YakuTanyao:     "断幺",
	YakuYakuhai:    "役牌",
	YakuChanta:     "混全",
	YakuJunchan:    "纯全",
	YakuHonroutou:  "混老头", // 七对也算
	YakuShousangen: "小三元",

	// Yaku based on suits
	YakuHonitsu:  "混一色",
	YakuChinitsu: "清一色",
}

// 调试用
func YakuTypesToStr(yakuTypes []int) string {
	names := []string{}
	for _, t := range yakuTypes {
		names = append(names, YakuNameMap[t])
	}
	if len(names) == 0 {
		return "[无役]"
	}
	return fmt.Sprint(names)
}

//

type _yakuHanMap map[int]int

var YakuHanMap = _yakuHanMap{
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

var NakiYakuHanMap = _yakuHanMap{
	YakuHaitei:  1,
	YakuHoutei:  1,
	YakuRinshan: 1,
	YakuChankan: 1,

	YakuSanshokuDoujun: 1,
	YakuIttsuu:         1,

	YakuToitoi:         2,
	YakuSanAnkou:       2,
	YakuSanshokuDoukou: 2,
	YakuSanKantsu:      2,

	YakuTanyao:     1,
	YakuYakuhai:    1,
	YakuChanta:     1,
	YakuJunchan:    2,
	YakuHonroutou:  2,
	YakuShousangen: 2,

	YakuHonitsu:  2,
	YakuChinitsu: 5,
}

var YakumanTimesMap = map[int]int{
	// TODO
}

// 计算 yakuTypes 累积的番数
func CalcYakuHan(yakuTypes []int, isNaki bool) (cntHan int) {
	var yakuHanMap _yakuHanMap
	if !isNaki {
		yakuHanMap = YakuHanMap
	} else {
		yakuHanMap = NakiYakuHanMap
	}

	for _, yakuType := range yakuTypes {
		if han, ok := yakuHanMap[yakuType]; ok {
			cntHan += han
		}
	}
	return
}

func CalcYakumanTimes() int {
	// TODO
	return 0
}
