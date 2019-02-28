package main

import (
	"strings"
	"strconv"
	"fmt"
	"regexp"
)

type tenhouMessage struct {
	Tag string `json:"tag"`

	//Name string `json:"name"` // id
	//Sex  string `json:"sx"`

	//UserName    string `json:"uname"`
	//RatingScale string `json:"ratingscale"`

	//N string `json:"n"`
	//J string `json:"j"`
	//G string `json:"g"`

	// round 开始 tag=INIT
	Seed   string `json:"seed"` // 本局信息：场数，场棒数，立直棒数，骰子A减一，骰子B减一，宝牌指示牌 1,0,0,3,2,92
	Ten    string `json:"ten"`  // 各家点数 280,230,240,250
	Dealer string `json:"oya"`  // 庄家 0=自家, 1=下家, 2=对家, 3=上家
	Hai    string `json:"hai"`  // 初始手牌 30,114,108,31,78,107,25,23,2,14,122,44,49

	// 摸牌 tag=T编号，如 T68

	// 副露 tag=N
	Who  string `json:"who"` // 副露者 0=自家, 1=下家, 2=对家, 3=上家
	Meld string `json:"m"`   // 副露编号 35914

	// 杠宝牌指示牌 tag=DORA
	// `json:"hai"` // 杠宝牌指示牌 39

	// 立直声明 tag=REACH, step=1
	// `json:"who"` // 立直者
	Step string `json:"step"` // 1

	// 立直成功，扣1000点 tag=REACH, step=2
	// `json:"who"` // 立直者
	// `json:"ten"` // 立直成功后的各家点数 250,250,240,250
	// `json:"step"` // 2

	// 自摸/有人放铳 tag=牌, t>=8
	T string `json:"t"` // 选项

	// 和牌 tag=AGARI
	// ba, hai, m, machi, ten, yaku, doraHai, who, fromWho, sc
	//Ba string `json:"ba"` // 0,0
	// `json:"hai"` // 和牌型 8,9,11,14,19,125,126,127
	// `json:"m"` // 副露编号 13527,50794
	//Machi string `json:"machi"` // (待ち) 自摸/荣和的牌 126
	// `json:"ten"` // 符数,点数,这张牌的来源 30,7700,0
	//Yaku        string `json:"yaku"`       // 役（编号，翻数） 18,1,20,1,34,2
	//DoraTile    string `json:"doraHai"`    // 宝牌 123
	//UraDoraTile string `json:"doraHaiUra"` // 里宝牌 77
	// `json:"who"` // 和牌者
	//FromWho string `json:"fromWho"` // 自摸/荣和牌的来源
	//Score   string `json:"sc"`      // 各家增减分 260,-77,310,77,220,0,210,0

	// 游戏结束 tag=PROF

	// 重连 tag=GO
	// type, lobby, gpid
	//Type  string `json:"type"`
	//Lobby string `json:"lobby"`
	//GPID  string `json:"gpid"`

	// 重连 tag=REINIT
	// `json:"seed"`
	// `json:"ten"`
	// `json:"oya"`
	// `json:"hai"`
	//Meld1    string `json:"m1"` // 各家副露编号 17450
	//Meld2    string `json:"m2"`
	//Meld3    string `json:"m3"`
	//Kawa0 string `json:"kawa0"` // 各家牌河 112,73,3,131,43,98,78,116
	//Kawa1 string `json:"kawa1"`
	//Kawa2 string `json:"kawa2"`
	//Kawa3 string `json:"kawa3"`
}

//

type tenhouRoundData struct {
	*roundData

	originJSON string
	msg        *tenhouMessage
}

//func (d *tenhouRoundData) mergeCachedTile() {
//	if cachedTile == -1 {
//		return
//	}
//
//
//
//	cachedTile = -1
//}

// 0-35 m
// 36-71 p
// 72-107 s
// 108- z
func (*tenhouRoundData) _parseTenhouTile(tile string) int {
	t, err := strconv.Atoi(tile)
	if err != nil {
		panic(err)
	}
	return t / 4
}

func (*tenhouRoundData) _parseChi(data int) (meldType int, tiles []int, calledTile int) {
	// 吃
	meldType = meldTypeChi
	t0, t1, t2 := (data>>3)&0x3, (data>>5)&0x3, (data>>7)&0x3
	baseAndCalled := data >> 10
	base, called := baseAndCalled/3, baseAndCalled%3
	base = (base/7)*9 + base%7
	tiles = []int{(t0 + 4*(base+0)) / 4, (t1 + 4*(base+1)) / 4, (t2 + 4*(base+2)) / 4}
	calledTile = tiles[called]
	return
}

func (*tenhouRoundData) _parsePon(data int) (meldType int, tiles []int, calledTile int) {
	t4 := (data >> 5) & 0x3
	_t := [4][3]int{{1, 2, 3}, {0, 2, 3}, {0, 1, 3}, {0, 1, 2}}[t4]
	t0, t1, t2 := _t[0], _t[1], _t[2]
	baseAndCalled := data >> 9
	base, called := baseAndCalled/3, baseAndCalled%3
	if data&0x8 > 0 {
		// 碰
		meldType = meldTypePon
		tiles = []int{(t0 + 4*base) / 4, (t1 + 4*base) / 4, (t2 + 4*base) / 4}
		calledTile = tiles[called]
	} else {
		// 加杠
		meldType = meldTypeKakan
		tiles = []int{(t0 + 4*base) / 4, (t1 + 4*base) / 4, (t2 + 4*base) / 4, (t4 + 4*base) / 4}
		calledTile = tiles[3]
	}
	return
}

func (d *tenhouRoundData) _parseKan(data int) (meldType int, tiles []int, calledTile int) {
	baseAndCalled := data >> 8
	base, called := baseAndCalled/4, baseAndCalled%4
	tiles = []int{(4 * base) / 4, (1 + 4*base) / 4, (2 + 4*base) / 4, (3 + 4*base) / 4}
	calledTile = tiles[called]

	// 通过判断 calledTile 的来源来是否为上一张舍牌，来判断是明杠还是暗杠
	if len(d.globalDiscardTiles) > 0 && calledTile == d.globalDiscardTiles[len(d.globalDiscardTiles)-1] {
		// 明杠
		meldType = meldTypeMinKan
	} else {
		// 暗杠
		meldType = meldTypeAnKan
	}
	return
}

func (d *tenhouRoundData) _parseTenhouMeld(data string) (meldType int, tiles []int, calledTile int) {
	bits, err := strconv.Atoi(data)
	if err != nil {
		panic(err)
	}

	switch {
	case bits&0x4 > 0:
		return d._parseChi(bits)
	case bits&0x18 > 0:
		return d._parsePon(bits)
	case bits&0x20 > 0:
		// 拔北
		panic("暂不支持三人麻将")
	default:
		return d._parseKan(bits)
	}
}

func (d *tenhouRoundData) GetDataSourceType() int {
	return dataSourceTypeTenhou
}

func (d *tenhouRoundData) GetMessage() string {
	return d.originJSON
}

func (d *tenhouRoundData) CheckMessage() bool {
	return true
}

func (d *tenhouRoundData) IsInit() bool {
	return d.msg.Tag == "INIT" || d.msg.Tag == "REINIT"
}

func (d *tenhouRoundData) ParseInit() (roundNumber int, dealer int, doraIndicator int, hands []int) {
	splits := strings.Split(d.msg.Seed, ",")
	if len(splits) != 6 {
		panic(fmt.Sprintln("seed 解析失败", d.msg.Seed))
	}
	roundNumber, _ = strconv.Atoi(splits[0])
	dealer, _ = strconv.Atoi(d.msg.Dealer)
	doraIndicator = d._parseTenhouTile(splits[5])
	for _, tenhouTile := range strings.Split(d.msg.Hai, ",") {
		tile := d._parseTenhouTile(tenhouTile)
		hands = append(hands, tile)
	}
	return
}

var _selfDrawReg = regexp.MustCompile("^T[0-9]{1,3}$")

func (d *tenhouRoundData) IsSelfDraw() bool {
	return _selfDrawReg.MatchString(d.msg.Tag)
}

func (d *tenhouRoundData) ParseSelfDraw() (tile int, kanDoraIndicator int) {
	rawTile := d.msg.Tag[1:]
	tile = d._parseTenhouTile(rawTile)
	kanDoraIndicator = -1
	return
}

var _discardReg = regexp.MustCompile("^[DEFGefg][0-9]{1,3}$")

func (d *tenhouRoundData) IsDiscard() bool {
	return _discardReg.MatchString(d.msg.Tag)
}

func (d *tenhouRoundData) ParseDiscard() (who int, tile int, isTsumogiri bool, isReach bool, canBeMeld bool, kanDoraIndicator int) {
	// D=自家, e/E=下家, f/F=对家, g/G=上家
	who = int(lower(d.msg.Tag[0]) - 'd')
	rawTile := d.msg.Tag[1:]
	tile = d._parseTenhouTile(rawTile)
	if d.msg.Tag[0] != 'D' {
		isTsumogiri = d.msg.Tag[0] >= 'a'
		canBeMeld = d.msg.T != ""
	}
	kanDoraIndicator = -1
	return
}

func (d *tenhouRoundData) IsOpen() bool {
	return d.msg.Tag == "N"
}

func (d *tenhouRoundData) ParseOpen() (who int, meldType int, meldTiles []int, calledTile int, kanDoraIndicator int) {
	who, _ = strconv.Atoi(d.msg.Who)
	meldType, meldTiles, calledTile = d._parseTenhouMeld(d.msg.Meld)
	kanDoraIndicator = -1
	return
}

func (d *tenhouRoundData) IsReach() bool {
	// Step == "1" 立直宣告
	// Step == "2" 立直成功，扣1000点
	return d.msg.Tag == "REACH" && d.msg.Step == "1"
}

func (d *tenhouRoundData) ParseReach() (who int) {
	who, _ = strconv.Atoi(d.msg.Who)
	return
}

func (d *tenhouRoundData) IsFuriten() bool {
	return d.msg.Tag == "FURITEN"
}

func (d *tenhouRoundData) IsRoundWin() bool {
	return false
}

func (d *tenhouRoundData) ParseRoundWin() (whos []int, points []int) {
	return
}

func (d *tenhouRoundData) IsNewDora() bool {
	return d.msg.Tag == "DORA"
}

func (d *tenhouRoundData) ParseNewDora() (kanDoraIndicator int) {
	kanDoraIndicator = d._parseTenhouTile(d.msg.Hai)
	return
}
