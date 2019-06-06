package util

import (
	"fmt"
	"sort"
)

var considerOldYaku bool

func SetConsiderOldYaku(b bool) {
	considerOldYaku = b
}

//

const (
	// https://en.wikipedia.org/wiki/Japanese_Mahjong_yaku
	// Special criteria
	YakuRiichi int = iota
	YakuChiitoi

	// Yaku based on luck
	YakuTsumo
	//YakuIppatsu
	//YakuHaitei
	//YakuHoutei
	//YakuRinshan
	//YakuChankan
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

	// Yakuman
	//YakuKokushi
	//YakuKokushi13
	YakuSuuAnkou
	YakuSuuAnkouTanki
	YakuDaisangen
	YakuShousuushii
	YakuDaisuushii
	YakuTsuuiisou
	YakuChinroutou
	YakuRyuuiisou
	YakuChuuren
	YakuChuuren9
	YakuSuuKantsu
	//YakuTenhou
	//YakuChiihou

	// 古役
	YakuShiiaruraotai
	YakuUumensai
	YakuSanrenkou
	YakuIsshokusanjun

	// 古役役满
	YakuDaisuurin
	YakuDaisharin
	YakuDaichikurin
	YakuDaichisei

	//_endYakuType  // 标记 enum 结束，方便计算有多少个 YakuType
)

//const maxYakuType = _endYakuType

var YakuNameMap = map[int]string{
	// Special criteria
	YakuRiichi:  "立直",
	YakuChiitoi: "七对",

	// Yaku based on luck
	YakuTsumo: "自摸",
	//YakuIppatsu: "一发",
	//YakuHaitei:  "海底",
	//YakuHoutei:  "河底",
	//YakuRinshan: "岭上",
	//YakuChankan: "抢杠",
	YakuDaburii: "w立",

	// Yaku based on sequences
	YakuPinfu:          "平和",
	YakuRyanpeikou:     "两杯口",
	YakuIipeikou:       "一杯口",
	YakuSanshokuDoujun: "三色",
	YakuIttsuu:         "一通", // 一气

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

	// Yakuman
	//YakuKokushi:       "国士",
	//YakuKokushi13:     "国士十三面",
	YakuSuuAnkou:      "四暗刻",
	YakuSuuAnkouTanki: "四暗刻单骑",
	YakuDaisangen:     "大三元",
	YakuShousuushii:   "小四喜",
	YakuDaisuushii:    "大四喜",
	YakuTsuuiisou:     "字一色",
	YakuChinroutou:    "清老头",
	YakuRyuuiisou:     "绿一色",
	YakuChuuren:       "九莲",
	YakuChuuren9:      "纯正九莲",
	YakuSuuKantsu:     "四杠子",
	//YakuTenhou:        "天和",
	//YakuChiihou:       "地和",
}

var OldYakuNameMap = map[int]string{
	YakuShiiaruraotai: "十二落抬",
	YakuUumensai:      "五门齐",
	YakuSanrenkou:     "三连刻",
	YakuIsshokusanjun: "一色三顺",

	YakuDaisuurin:   "大数邻",
	YakuDaisharin:   "大车轮",
	YakuDaichikurin: "大竹林",
	YakuDaichisei:   "大七星",
}

func YakuTypesToStr(yakuTypes []int) string {
	if len(yakuTypes) == 0 {
		return "[无役]"
	}
	names := []string{}
	for _, t := range yakuTypes {
		if name, ok := YakuNameMap[t]; ok {
			names = append(names, name)
		}
	}

	if considerOldYaku {
		for _, t := range yakuTypes {
			if name, ok := OldYakuNameMap[t]; ok {
				names = append(names, name)
			}
		}
	}

	return fmt.Sprint(names)
}

func YakuTypesWithDoraToStr(yakuTypes map[int]struct{}, numDora int) string {
	if len(yakuTypes) == 0 {
		return "[无役]"
	}
	yt := []int{}
	for t := range yakuTypes {
		yt = append(yt, t)
	}
	sort.Ints(yt)
	names := []string{}
	for _, t := range yt {
		names = append(names, YakuNameMap[t])
	}
	// TODO: old yaku
	if numDora > 0 {
		names = append(names, fmt.Sprintf("宝牌%d", numDora))
	}
	return fmt.Sprint(names)
}

//

type _yakuHanMap map[int]int
type _yakumanTimesMap map[int]int

var YakuHanMap = _yakuHanMap{
	YakuRiichi:  1,
	YakuChiitoi: 2,

	YakuTsumo: 1,
	//YakuIppatsu: 1,
	//YakuHaitei:  1,
	//YakuHoutei:  1,
	//YakuRinshan: 1,
	//YakuChankan: 1,
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
	//YakuHaitei:  1,
	//YakuHoutei:  1,
	//YakuRinshan: 1,
	//YakuChankan: 1,

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

var OldYakuHanMap = _yakuHanMap{
	YakuUumensai:      2,
	YakuSanrenkou:     2,
	YakuIsshokusanjun: 3,
}

var OldNakiYakuHanMap = _yakuHanMap{
	YakuShiiaruraotai: 1, // 四副露大吊车
	YakuUumensai:      2,
	YakuSanrenkou:     2,
	YakuIsshokusanjun: 2,
}

// 计算 yakuTypes(非役满) 累积的番数
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

	if considerOldYaku {
		if !isNaki {
			yakuHanMap = OldYakuHanMap
		} else {
			yakuHanMap = OldNakiYakuHanMap
		}

		for _, yakuType := range yakuTypes {
			if han, ok := yakuHanMap[yakuType]; ok {
				cntHan += han
			}
		}
	}

	return
}

//

var YakumanTimesMap = map[int]int{
	//YakuKokushi:       1,
	//YakuKokushi13:     2,
	YakuSuuAnkou:      1,
	YakuSuuAnkouTanki: 2,
	YakuDaisangen:     1,
	YakuShousuushii:   1,
	YakuDaisuushii:    2,
	YakuTsuuiisou:     1,
	YakuChinroutou:    1,
	YakuRyuuiisou:     1,
	YakuChuuren:       1,
	YakuChuuren9:      2,
	YakuSuuKantsu:     1,
	//YakuTenhou:        1,
	//YakuChiihou:       1,
}

var NakiYakumanTimesMap = map[int]int{
	YakuDaisangen:   1,
	YakuShousuushii: 1,
	YakuDaisuushii:  2,
	YakuTsuuiisou:   1,
	YakuChinroutou:  1,
	YakuRyuuiisou:   1,
	YakuSuuKantsu:   1,
}

var OldYakumanTimesMap = map[int]int{
	YakuDaisuurin:   1,
	YakuDaisharin:   1,
	YakuDaichikurin: 1,
	YakuDaichisei:   1, // 复合字一色，实际为两倍役满
}

// 计算役满倍数
func CalcYakumanTimes(yakuTypes []int, isNaki bool) (times int) {
	var yakumanTimesMap _yakumanTimesMap
	if !isNaki {
		yakumanTimesMap = YakumanTimesMap
	} else {
		yakumanTimesMap = NakiYakumanTimesMap
	}

	for _, yakuman := range yakuTypes {
		if t, ok := yakumanTimesMap[yakuman]; ok {
			times += t
		}
	}

	if considerOldYaku && !isNaki {
		for _, yakuman := range yakuTypes {
			if t, ok := OldYakumanTimesMap[yakuman]; ok {
				times += t
			}
		}
	}

	return
}
