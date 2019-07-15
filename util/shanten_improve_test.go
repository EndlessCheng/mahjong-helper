package util

import (
	"testing"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"github.com/stretchr/testify/assert"
	"strings"
)

func Test_calculateIsolatedTileValue(t *testing.T) {
	assert := assert.New(t)

	newPI := func(selfWindTile int, roundWindTile int, discardedHumanTiles string) *model.PlayerInfo {
		return &model.PlayerInfo{
			SelfWindTile:  selfWindTile,
			RoundWindTile: roundWindTile,
			LeftTiles34:   InitLeftTiles34WithTiles34(MustStrToTiles34(discardedHumanTiles)),
		}
	}

	const eps = 1e-3

	assert.InDelta(100, float64(calculateIsolatedTileValue(MustStrToTile34("9m"), newPI(27, 27, "2s"))), eps)
	assert.InDelta(130, float64(calculateIsolatedTileValue(MustStrToTile34("1z"), newPI(27, 27, "2s"))), eps)
	assert.InDelta(117, float64(calculateIsolatedTileValue(MustStrToTile34("1z"), newPI(27, 27, "2s11z"))), eps)
	assert.InDelta(97, float64(calculateIsolatedTileValue(MustStrToTile34("2z"), newPI(27, 27, "2s"))), eps)
	assert.InDelta(98, float64(calculateIsolatedTileValue(MustStrToTile34("3z"), newPI(27, 27, "2s"))), eps)
	assert.InDelta(99, float64(calculateIsolatedTileValue(MustStrToTile34("4z"), newPI(27, 27, "2s"))), eps)
	assert.InDelta(114.9, float64(calculateIsolatedTileValue(MustStrToTile34("5z"), newPI(27, 27, "2s"))), eps)
	assert.InDelta(114.8, float64(calculateIsolatedTileValue(MustStrToTile34("6z"), newPI(27, 27, "2s"))), eps)
	assert.InDelta(115, float64(calculateIsolatedTileValue(MustStrToTile34("7z"), newPI(27, 27, "2s"))), eps)
	assert.InDelta(103.5, float64(calculateIsolatedTileValue(MustStrToTile34("7z"), newPI(27, 27, "2s77z"))), eps)
	assert.InDelta(23, float64(calculateIsolatedTileValue(MustStrToTile34("7z"), newPI(27, 27, "2s777z"))), eps)

	assert.InDelta(114, float64(calculateIsolatedTileValue(MustStrToTile34("1z"), newPI(29, 27, "2s"))), eps)
	assert.InDelta(102.6, float64(calculateIsolatedTileValue(MustStrToTile34("1z"), newPI(29, 27, "2s11z"))), eps)
	assert.InDelta(99, float64(calculateIsolatedTileValue(MustStrToTile34("2z"), newPI(29, 27, "2s"))), eps)
	assert.InDelta(116, float64(calculateIsolatedTileValue(MustStrToTile34("3z"), newPI(29, 27, "2s"))), eps)
	assert.InDelta(97, float64(calculateIsolatedTileValue(MustStrToTile34("4z"), newPI(29, 27, "2s"))), eps)
}

//

var exampleMelds = []model.Meld{{MeldType: model.MeldTypePon, Tiles: MustStrToTiles("666z")}}

func TestCalculateShantenWithImproves13Closed(t *testing.T) {
	t.Skip()
	for _, tiles := range []string{
		//"11357m 13579p 135s",
		//"123456789m 1135s",
		//"123456789m 1134s",
		//"123456789m 1234z",
		//"3m 12668p 5678s 222z",
		//"6m 12668p 5678s 222z",
		//"557m 34789p 26s 111z",
		//"111333555m 23p 23s",
		//"23m 234p 234888s 44z",
		//"23467m 234p 23488s",
		//"34567m 22334455p",
		//"1199m 112235566z",
		//"123456789m 23p 88s",
		//"12399m 123p 12999s",
		//"3577m 345p 345678s",
		//"23467m 234p 23488s",
		//"13789m 11p 345s 555z",
		//"12346789m 123p 88s",
		//"3456m 111s 999p 777z",
		//"123m 44p 34888s 777z",
		//"13789m 111789p 77z",
		//"23467m 222p 23488s",
		//"13789m 111789p 11s",
		//"12346789m 123p 88s",
		//"56778p 112345s 77z",
		//"56778p 122345s 77z",
		"223446m 345p 1178s",
		"122344m 345p 1178s",
	} {
		tiles34 := MustStrToTiles34(tiles)
		if CountOfTiles34(tiles34) != 13 {
			t.Error(tiles, "不是13张牌")
			continue
		}
		playerInfo := model.NewSimplePlayerInfo(tiles34, nil)
		playerInfo.DoraTiles = MustStrToTiles("8s")
		//playerInfo.DiscardTiles = []int{MustStrToTile34("4s")}
		//playerInfo.IsRiichi = true
		//playerInfo.DoraCount = 2
		result := CalculateShantenWithImproves13(playerInfo)
		t.Log(tiles, "=\n"+result.String())
		for tile, left := range result.Waits {
			if left > 0 {
				playerInfo.HandTiles34[tile]++
				t.Log(Tiles34ToStr(playerInfo.HandTiles34))
				_, results, _ := CalculateShantenWithImproves14(playerInfo)
				for _, result := range results {
					t.Log(result)
				}
				t.Log()
				playerInfo.HandTiles34[tile]--
			}
		}
	}
}

func TestCalculateShantenWithImproves13Open(t *testing.T) {
	t.Skip()
	for _, tiles := range []string{
		"1234m",
		"1135m",
		"5p",
		"5555m", // 一向听，132 进张
	} {
		tiles34 := MustStrToTiles34(tiles)
		result := CalculateShantenWithImproves13(model.NewSimplePlayerInfo(tiles34, exampleMelds))
		t.Logf("%s = %s\n%s", tiles, NumberToChineseShanten(result.Shanten), result.String())
	}
}

func TestCalculateShantenWithImproves14Closed(t *testing.T) {
	t.Skip()
	tiles := "124679m 3678p 2366s"
	tiles = "11379m 347p 277s 111z"
	tiles = "334578m 11468p 235s"
	tiles = "478m 33588p 457899s"
	tiles = "2233688m 1234p 379s"
	tiles = "1233347m 23699p 88s"
	tiles = "334m 22457p 23456s 1z"
	tiles = "334m 122478p 23456s"
	tiles = "1m 258p 258s 1234567z"
	tiles = "4567m 4579p 344588s"
	tiles = "2479999m 45667p 13s" // 切任何一张都不会向听倒退
	tiles = "25667m 27789p 37s 44z"
	tiles = "111444777m 11177s"
	tiles = "2468m 33578p 22356s"
	tiles = "57m 4455p 12345699s"
	tiles = "57m 3445667p 12399s"
	tiles = "2335578899m 5677p"
	tiles = "123p 3445668m 6799s"
	tiles = "455678m 11566p 234s" // TODO 振听 9m 的场合，切 6p 振听听牌的概率比切 5m 低
	tiles = "1245m 12789p 34588s"
	tiles = "4456778p 2245s 111z"
	tiles = "388m 113668p 56s 456z"
	tiles = "56778p 1122345s 77z"
	tiles = "66778p 1122345s 77z"
	tiles = "67778p 1122345s 77z"
	tiles = "3336888m 678p 5678s"
	tiles = "11456678m 567p 235s" // 振听两面还是坎张
	tiles = "123m 1234789p 3388s"
	tiles = "789m 123467789p 11z"
	tiles = "11122m 199p 2455s 56z"
	tiles = "347m 579p 246s 12345z"
	tiles = "13m 344579p 5699s 15z"
	tiles = "145599m 18p 124s 157z"
	tiles = "2389m 4566p 2289s 44z" // 垃圾进张考虑改良
	tiles = "11233456789m 333p"     // 默听进张和立直不一样
	tiles = "136688m 456p 55677s"
	playerInfo := model.NewSimplePlayerInfo(MustStrToTiles34(tiles), nil)
	//playerInfo.SelfWindTile = MustStrToTile34("4z")
	playerInfo.LeftTiles34 = InitLeftTiles34WithTiles34(MustStrToTiles34("1366668888m 456p 55666677s")) // 注意手牌也算上
	//playerInfo.DiscardTiles = []int{MustStrToTile34("1p")}
	//playerInfo.DoraTiles = MustStrToTiles("3m")
	shanten, results, incShantenResults := CalculateShantenWithImproves14(playerInfo)
	t.Log(NumberToChineseShanten(shanten))
	for _, result := range results {
		t.Log(result)
	}
	if len(incShantenResults) > 0 {
		t.Log(NumberToChineseShanten(shanten + 1))
		for _, result := range incShantenResults {
			t.Log(result)
		}
	} else {
		t.Log("无向听倒退的切牌")
	}
}

func TestCalculateShantenWithImproves14Open(t *testing.T) {
	t.Skip()
	tiles := "35m"
	tiles = "13m 456s 778p"
	tiles = "6888m 678p 5678s"
	//leftTiles34 := InitLeftTiles34WithTiles34(MustStrToTiles34(tiles))
	//leftTiles34[1] = 0
	playerInfo := model.NewSimplePlayerInfo(MustStrToTiles34(tiles), exampleMelds)
	//playerInfo.LeftTiles34 = leftTiles34
	shanten, results, incShantenResults := CalculateShantenWithImproves14(playerInfo)
	t.Log(NumberToChineseShanten(shanten))
	for _, result := range results {
		t.Log(result)
	}
	t.Log(NumberToChineseShanten(shanten + 1))
	for _, result := range incShantenResults {
		t.Log(result)
	}
}

func TestCalculateMeld(t *testing.T) {
	t.Skip()
	tiles := "1234m 112z"
	tiles = "23445667m 11z"
	tiles = "112356799m 1233z"
	tiles = "78m 12355p 789s" // ***
	tiles = "245689s 1z"
	tiles = "466m 234467p 77s 77z"
	tiles34 := MustStrToTiles34(tiles)
	pi := model.NewSimplePlayerInfo(tiles34, nil)
	if len(pi.HandTiles34) < 13 {
		pi.Melds = exampleMelds
	}
	result := CalculateShantenWithImproves13(pi)
	t.Log("原手牌" + NumberToChineseShanten(result.Shanten))
	t.Log(result)

	tile := "1m"
	tile = "3m" // "1z"
	tile = "4m"
	tile = "4p"
	tile = "3s"
	tile = "7z"
	shanten, results, incShantenResults := CalculateMeld(pi, MustStrToTile34(tile), false, true)
	t.Log("鸣牌后" + NumberToChineseShanten(shanten))
	for _, result := range results {
		t.Log(result)
	}
	t.Log("鸣牌后" + NumberToChineseShanten(shanten+1))
	for _, result := range incShantenResults {
		t.Log(result)
	}
}

//

func bestHumanDiscardTile(t *testing.T, humanTiles string, doraHumanTiles string) string {
	playerInfo := model.NewSimplePlayerInfo(MustStrToTiles34(humanTiles), nil)
	if doraHumanTiles != "" {
		playerInfo.DoraTiles = MustStrToTiles(doraHumanTiles)
	}
	_, results, _ := CalculateShantenWithImproves14(playerInfo)
	if true {
		t.Log(humanTiles, doraHumanTiles)
		for _, result := range results {
			t.Log(result)
		}
		t.Log()
	}
	tile := results[0].DiscardTile
	return Tile34ToStr(tile)
}

func bestHumanDiscardTile2(t *testing.T, humanTiles string, doraIndicatorHumanTiles string) string {
	// 根据 0 来记录赤宝牌
	numRedFives := make([]int, 3)
	for _, split := range strings.Split(humanTiles, " ") {
		for _, c := range split {
			if c == '0' {
				numRedFives[ByteAtStr(split[len(split)-1], "mps")]++
				break
			}
		}
	}
	humanTiles = strings.Replace(humanTiles, "0", "5", -1)
	playerInfo := model.NewSimplePlayerInfo(MustStrToTiles34(humanTiles), nil)
	if doraIndicatorHumanTiles != "" {
		doraIndicators := MustStrToTiles(doraIndicatorHumanTiles)
		for _, tile := range doraIndicators {
			playerInfo.LeftTiles34[tile]--
		}
		playerInfo.DoraTiles = model.DoraList(doraIndicators, false)
	}
	playerInfo.NumRedFives = numRedFives
	_, results, _ := CalculateShantenWithImproves14(playerInfo)
	if true {
		t.Log(humanTiles, doraIndicatorHumanTiles)
		for _, result := range results {
			t.Log(result)
		}
		t.Log()
	}
	tile := results[0].DiscardTile
	return Tile34ToStr(tile)
}

// CxQx 来源：知るだけで強くなる麻雀の2択
func TestBestDiscard(t *testing.T) {
	assert := assert.New(t)

	// 听牌
	assert.Equal("7m", bestHumanDiscardTile(t, "123667m 234p 345s 55z", ""))   // C3Q4 数牌字牌双碰优于两面
	assert.Equal("6m", bestHumanDiscardTile(t, "123667m 234p 345s 44z", ""))   // C3Q4 平和
	assert.Equal("4m", bestHumanDiscardTile(t, "134m 123567p 12355s", ""))     // C3Q5 三色
	assert.Equal("4m", bestHumanDiscardTile(t, "134m 123567p 12355s", "5p"))   // C3Q5 三色
	assert.Equal("4s", bestHumanDiscardTile(t, "234456m 11567p 468s", ""))     // C3Q8 筋引挂比赤5好
	assert.Equal("1m", bestHumanDiscardTile(t, "1234m 345789p 567s 3z", "3z")) // C3Q10 宝牌单骑比两面好
	assert.Equal("5s", bestHumanDiscardTile(t, "345m 345789p 3455s 4z", ""))   // C3Q15 三色比平和好
	assert.Equal("7m", bestHumanDiscardTile(t, "234788m 234567s 33z", "8m"))   // C3Q17 和率下降一点但是打点提升
	assert.Equal("8m", bestHumanDiscardTile(t, "234788m 234567s 33z", "3z"))   // C3Q17 打点充足时和率优先
	assert.Equal("7m", bestHumanDiscardTile(t, "334557m 222p 789s 33z", "9s")) // C3Q18 和率下降一点但是打点提升

	// 一向听 技术论
	//assert.Equal("1m", bestHumanDiscardTile(t, "1223446m 345p 1178s", "8s"))
	//assert.Equal("6m", bestHumanDiscardTile(t, "1223446m 345p 78s 77z", "8s"))
	//assert.Equal("2m", bestHumanDiscardTile(t, "1223446789m 1178s", "8s")) // 4m 也可以
	//assert.Equal("2m", bestHumanDiscardTile(t, "1223446789m 78s 77z", "8s"))

	//assert.Equal("4m", bestHumanDiscardTile(t, "334456788m 45p 456s", ""))

	// 两向听

	// 三向听
	assert.Equal("8p", bestHumanDiscardTile(t, "23668m 258p 4678s 77z", "2p"))
	assert.Equal("2p", bestHumanDiscardTile(t, "23668m 258p 4678s 77z", "8p"))
	assert.Equal("2p", bestHumanDiscardTile(t, "23668m 258p 4678s 77z", ""))
}

func TestBest2(t *testing.T) {
	assert := assert.New(t)

	// 传入手牌和宝牌指示牌，赤牌用 0 表示
	assert.Equal("5z", bestHumanDiscardTile2(t, "789m 23678p 1235s 25z", "3z"))
	assert.Equal("8p", bestHumanDiscardTile2(t, "24668m 2078p 23457s", "3z"))
}

func TestFuritenBestDiscard(t *testing.T) {
	assert := assert.New(t)

	bestHumanDiscardTileWhenFuriten := func(t *testing.T, humanTiles string, doraHumanTiles string, selfDiscardHumanTiles string) string {
		playerInfo := model.NewSimplePlayerInfo(MustStrToTiles34(humanTiles), nil)
		if doraHumanTiles != "" {
			playerInfo.DoraTiles = MustStrToTiles(doraHumanTiles)
		}
		playerInfo.DiscardTiles = MustStrToTiles(selfDiscardHumanTiles)
		for _, dis := range playerInfo.DiscardTiles {
			playerInfo.LeftTiles34[dis]--
		}
		_, results, _ := CalculateShantenWithImproves14(playerInfo)
		if true {
			t.Log(humanTiles, doraHumanTiles)
			for _, result := range results {
				t.Log(result)
			}
			t.Log()
		}
		tile := results[0].DiscardTile
		return Tile34ToStr(tile)
	}

	// 振听听牌
	assert.Equal("5s", bestHumanDiscardTileWhenFuriten(t, "11456678m 567p 235s", "8m", "1s")) // C3Q6
	assert.Equal("6p", bestHumanDiscardTileWhenFuriten(t, "455678m 11566p 234s", "", "9m"))
}

// TODO: 测试何切 300 和 301
func TestQ300(t *testing.T) {
	assert := assert.New(t)

	// 传入手牌和宝牌指示牌，赤牌用 0 表示
	assert.Equal("2s", bestHumanDiscardTile2(t, "06778p 1122345s 77z", "2z"))                // Q001
	assert.Equal("5s", bestHumanDiscardTile2(t, "66778p 1122345s 77z", "2z"))                // Q002
	assert.Equal("8p", bestHumanDiscardTile2(t, "67778p 1122345s 77z", "2z"), "尚未考虑自摸时的三暗刻") // Q003
	assert.Equal("5s", bestHumanDiscardTile2(t, "12388m 455679p 556s", "2z"))                // Q004
	assert.Equal("9p", bestHumanDiscardTile2(t, "23488m 455679p 556s", "2z"))                // Q005
	assert.Equal("8p", bestHumanDiscardTile2(t, "33455m 668p 345667s", "2p"))                // Q006
	assert.Equal("1p", bestHumanDiscardTile2(t, "4406m 134556p 3478s", "2z"))                // Q007
	assert.Equal("2p", bestHumanDiscardTile2(t, "135m 11240667p 789s", "9s"), "麻雀是自摸的游戏")    // Q008
	assert.Equal("5m", bestHumanDiscardTile2(t, "135m 12399p 123667s", "2z"))                // Q009
	assert.Equal("3p", bestHumanDiscardTile2(t, "40699m 1133p 34567s", "1m"))                // Q010
	assert.Equal("1m", bestHumanDiscardTile2(t, "1234m 5678p 122233s", "8s"))                // Q011
	assert.Equal("2s", bestHumanDiscardTile2(t, "55678m 3467p 24668s", "6p"), "6s 也可以")      // Q012
	assert.Equal("2p", bestHumanDiscardTile2(t, "678m 123306p 12378s", "2z"))                // Q013
	assert.Equal("8p", bestHumanDiscardTile2(t, "23468p 130777s 444z", "1p"))                // Q030
	assert.Equal("3m", bestHumanDiscardTile2(t, "3356m 23478p 56777s", "7p"), "6m 也可以")      // Q033
	assert.Equal("1p", bestHumanDiscardTile2(t, "3456m 137899p 4578s", "1m"), "3m 也可以")      // Q037
	assert.Equal("8s", bestHumanDiscardTile2(t, "34056m 2224p 23468s", "8m"), "考虑改良的影响")     // Q058
}

//

func BenchmarkCalculateShantenWithImproves14_Shanten0(b *testing.B) {
	pi := model.NewSimplePlayerInfo(MustStrToTiles34("134m 123567p 12355s"), nil)
	for i := 0; i < b.N; i++ {
		// 157,447,560 ns/op
		// 由于考虑了向听倒退，所以速度和一向听差不多
		CalculateShantenWithImproves14(pi)
	}
}

func BenchmarkCalculateShantenWithImproves14_Shanten1(b *testing.B) {
	pi := model.NewSimplePlayerInfo(MustStrToTiles34("12388m 455679p 556s"), nil)
	for i := 0; i < b.N; i++ {
		// 145,454,040 ns/op
		// FIXME: 向听倒退？
		CalculateShantenWithImproves14(pi)
	}
}

func BenchmarkCalculateShantenWithImproves14_Shanten2(b *testing.B) {
	pi := model.NewSimplePlayerInfo(MustStrToTiles34("3456m 137899p 4578s"), nil)
	for i := 0; i < b.N; i++ {
		// 751,249,800 ns/op
		CalculateShantenWithImproves14(pi)
	}
}

func BenchmarkCalculateShantenWithImproves14_Shanten3(b *testing.B) {
	pi := model.NewSimplePlayerInfo(MustStrToTiles34("12688m 33579p 24s 56z"), nil)
	for i := 0; i < b.N; i++ {
		// 136,756,870 ns/op
		CalculateShantenWithImproves14(pi)
	}
}
