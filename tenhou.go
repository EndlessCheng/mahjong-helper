package main

import (
	"strings"
	"strconv"
	"fmt"
	"regexp"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"sort"
	"github.com/EndlessCheng/mahjong-helper/util"
	"net/url"
	"github.com/fatih/color"
)

/*
<SHUFFLE>
 - seed         Seed for RNG for generating walls and dice rolls.
 - ref          ?
<GO>            Start of game
 - type             Lobby type.
 - lobby            Lobby number.
<UN>            User list or user reconnect
 - n[0-3]           Names for each player as URLEncoded UTF-8.
 - dan              List of ranks for each player.
 - rate             List of rates for each player.
 - sx               List of sex ("M" or "F") for each player.
<BYE>           User disconnect
 - who              Player who disconnected.
<TAIKYOKU>      Start of round
 - oya              Dealer
<INIT>          Start of hand
 - seed             Six element list:
                        Round number,
                        Number of combo sticks,
                        Number of riichi sticks,
                        First dice minus one,
                        Second dice minus one,
                        Dora indicator.
 - ten              List of scores for each player
 - oya              Dealer
 - hai[0-3]         Starting hands as a list of tiles for each player.
<[T-W][0-9]*>   Player draws a tile.
<[D-G][0-9]*>   Player discards a tile.
<N>             Player calls a tile.
 - who              The player who called the tile.
 - m                The meld.
<REACH>         Player declares riichi.
 - who              The player who declared riichi
 - step             Where the player is in declaring riichi:
                        1 -> Called "riichi"
                        2 -> Placed point stick on table after discarding.
 - ten              List of current scores for each player.
<DORA>          New dora indicator.
 - hai              The new dora indicator tile.
<AGARI>         A player won the hand
 - who              The player who won.
 - fromwho          Who the winner won from: themselves for tsumo, someone else for ron.
 - hai              The closed hand of the winner as a list of tiles.
 - m                The open melds of the winner as a list of melds.
 - machi            The waits of the winner as a list of tiles.
 - doraHai          The dora as a list of tiles.
 - dorahaiUra       The ura dora as a list of tiles.
 - yaku             List of yaku and their han values.
                            0 -> tsumo
                            1 -> riichi
                            2 -> ippatsu
                            3 -> chankan
                            4 -> rinshan
                            5 -> haitei
                            6 -> houtei
                            7 -> pinfu
                            8 -> tanyao
                            9 -> ippeiko
                        10-17 -> fanpai
                        18-20 -> yakuhai
                           21 -> daburi
                           22 -> chiitoi
                           23 -> chanta
                           24 -> itsuu
                           25 -> sanshokudoujin
                           26 -> sanshokudou
                           27 -> sankantsu
                           28 -> toitoi
                           29 -> sanankou
                           30 -> shousangen
                           31 -> honrouto
                           32 -> ryanpeikou
                           33 -> junchan
                           34 -> honitsu
                           35 -> chinitsu
                           52 -> dora
                           53 -> uradora
                           54 -> akadora
 - yakuman          List of yakuman.
                           36 -> renhou
                           37 -> tenhou
                           38 -> chihou
                           39 -> daisangen
                        40,41 -> suuankou
                           42 -> tsuiisou
                           43 -> ryuuiisou
                           44 -> chinrouto
                        45,46 -> chuurenpooto
                        47,48 -> kokushi
                           49 -> daisuushi
                           50 -> shousuushi
                           51 -> suukantsu
 - ten              Three element list:
                        The fu points in the hand,
                        The point value of the hand,
                        The limit value of the hand:
                            0 -> No limit
                            1 -> Mangan
                            2 -> Haneman
                            3 -> Baiman
                            4 -> Sanbaiman
                            5 -> Yakuman
 - ba               Two element list of stick counts:
                        The number of combo sticks,
                        The number of riichi sticks.
 - sc               List of scores and the changes for each player.
 - owari            Final scores including uma at the end of the game.
<RYUUKYOKU>     The hand ended with a draw
 - type             The type of draw:
                        "yao9"   -> 9 ends
                        "reach4" -> Four riichi calls
                        "ron3"   -> Triple ron
                        "kan4"   -> Four kans
                        "kaze4"  -> Same wind discard on first round
                        "nm"     -> Nagashi mangan.
 - hai[0-3]         The hands revealed by players as a list of tiles.
 - ba               Two element list of stick counts:
                        The number of combo sticks,
                        The number of riichi sticks.
 - sc               List of scores and the changes for each player.
 - owari            Final scores including uma at the end of the game.
*/

const (
	redFiveMan = 16
	redFivePin = 52
	redFiveSou = 88
)

type tenhouMessage struct {
	Tag string `json:"tag" xml:"-"`

	//Name string `json:"name"` // id
	//Sex  string `json:"sx"`

	UserName string `json:"uname" xml:"-"`
	//RatingScale string `json:"ratingscale"`

	//N string `json:"n"`
	//J string `json:"j"`
	//G string `json:"g"`

	// round 开始 tag=INIT
	// 注意无论是三麻还是四麻，南1的场数都是4
	Seed   string `json:"seed" xml:"seed,attr"` // 本局信息：场数，场棒数，立直棒数，骰子A减一，骰子B减一，宝牌指示牌 1,0,0,3,2,92
	Ten    string `json:"ten" xml:"ten,attr"`   // 各家点数 280,230,240,250
	Dealer string `json:"oya" xml:"oya,attr"`   // 庄家 0=自家, 1=下家, 2=对家, 3=上家
	Hai    string `json:"hai" xml:"hai,attr"`   // 初始手牌 30,114,108,31,78,107,25,23,2,14,122,44,49
	Hai0   string `json:"-" xml:"hai0,attr"`
	Hai1   string `json:"-" xml:"hai1,attr"`
	Hai2   string `json:"-" xml:"hai2,attr"`
	Hai3   string `json:"-" xml:"hai3,attr"`

	// 摸牌 tag=T编号，如 T68

	// 副露 tag=N
	Who  string `json:"who" xml:"who,attr"` // 副露者 0=自家, 1=下家, 2=对家, 3=上家
	Meld string `json:"m" xml:"m,attr"`     // 副露编号 35914

	// 杠宝牌指示牌 tag=DORA
	// `json:"hai"` // 杠宝牌指示牌 39

	// 立直声明 tag=REACH, step=1
	// `json:"who"` // 立直者
	Step string `json:"step" xml:"step,attr"` // 1

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

	isRoundEnd bool // 某人和牌或流局。初始值为 true
}

func (*tenhouRoundData) _tenhouTileToTile34(tenhouTile int) int {
	return tenhouTile / 4
}

// 0-35 m
// 36-71 p
// 72-107 s
// 108- z
func (d *tenhouRoundData) _parseTenhouTile(tenhouTile string) (tile int, isRedFive bool) {
	t, err := strconv.Atoi(tenhouTile)
	if err != nil {
		panic(err)
	}
	return d._tenhouTileToTile34(t), d.isRedFive(t)
}

/*
CHI

 0                   1
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
| Base Tile | |   |   |   | |   |
|    and    |0| T2| T1| T0|1|Who|
|Called Tile| |   |   |   | |   |
+-----------+-+---+---+---+-+---+

    Base Tile and Called Tile:
        ((Base / 9) * 7 + Base % 9) * 3 + Chi
    T[0-2]:
        Tile[i] - 4 * i - Base * 4
    Who:
        Offset of player the tile was called from.
    Tile[0-2]:
        The tiles in the chi.
    Base:
        The lowest tile in the chi / 4.
    Called:
        Which tile out of the three was called.
*/
func (*tenhouRoundData) _parseChi(data int) (meldType int, tenhouMeldTiles []int, tenhouCalledTile int) {
	// 吃
	meldType = meldTypeChi
	t0, t1, t2 := (data>>3)&0x3, (data>>5)&0x3, (data>>7)&0x3
	baseAndCalled := data >> 10
	base, called := baseAndCalled/3, baseAndCalled%3
	base = (base/7)*9 + base%7
	tenhouMeldTiles = []int{t0 + 4*base, t1 + 4*(base+1), t2 + 4*(base+2)}
	tenhouCalledTile = tenhouMeldTiles[called]
	return
}

/*
PON or KAKAN

 0                   1
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|  Base Tile  |   |   |K|P| |   |
|     and     | 0 | T4|A|O|0|Who|
| Called Tile |   |   |N|N| |   |
+---------------+-+---+-+-+-+---+

    Base Tile and Called Tile:
        Base * 3 + Called
    T4:
        Tile4 - Base * 4
    PON:
        Set iff the meld is a pon.
    KAN:
        Set iff the meld is a pon upgraded to a kan.
    Who:
        Offset of player the tile was called from.
    Tile4:
        The tile which is not part of the pon.
    Base:
        A tile in the pon / 4.
    Called:
        Which tile out of the three was called.
*/
func (*tenhouRoundData) _parsePonOrKakan(data int) (meldType int, tenhouMeldTiles []int, tenhouCalledTile int) {
	t4 := (data >> 5) & 0x3
	_t := [4][3]int{{1, 2, 3}, {0, 2, 3}, {0, 1, 3}, {0, 1, 2}}[t4]
	t0, t1, t2 := _t[0], _t[1], _t[2]
	baseAndCalled := data >> 9
	base, called := baseAndCalled/3, baseAndCalled%3
	if data&0x8 > 0 {
		// 碰
		meldType = meldTypePon
		tenhouMeldTiles = []int{t0 + 4*base, t1 + 4*base, t2 + 4*base}
		tenhouCalledTile = tenhouMeldTiles[called]
	} else { // data&0x16 > 0
		// 加杠
		meldType = meldTypeKakan
		tenhouMeldTiles = []int{t0 + 4*base, t1 + 4*base, t2 + 4*base, t4 + 4*base}
		tenhouCalledTile = tenhouMeldTiles[3]
	}
	return
}

/*
KAN

 0                   1
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|   Base Tile   |           |   |
|      and      |     0     |Who|
|   Called Tile |           |   |
+---------------+-+---+-+-+-+---+

    Base Tile and Called Tile:
        Base * 4 + Called
    Who:
        Offset of player the tile was called from or 0 for a closed kan.
    Base:
        A tile in the kan / 4.
    Called:
        Which tile out of the four was called.
*/
func (*tenhouRoundData) _parseKan(data int) (meldType int, tenhouMeldTiles []int, tenhouCalledTile int) {
	baseAndCalled := data >> 8
	base, called := baseAndCalled/4, baseAndCalled%4
	tenhouMeldTiles = []int{4 * base, 1 + 4*base, 2 + 4*base, 3 + 4*base}
	tenhouCalledTile = tenhouMeldTiles[called]
	if offsetFromWho := data & 0x3; offsetFromWho == 0 {
		// 暗杠
		meldType = meldTypeAnkan
	} else {
		// 大明杠，offsetFromWho=1即为下家，=2为对家，=3为上家
		meldType = meldTypeMinkan
	}
	return
}

func (d *tenhouRoundData) _parseTenhouMeld(data string) (meldType int, tenhouMeldTiles []int, tenhouCalledTile int) {
	bits, err := strconv.Atoi(data)
	if err != nil {
		panic(err)
	}

	switch {
	case bits&0x4 > 0:
		return d._parseChi(bits)
	case bits&0x18 > 0:
		return d._parsePonOrKakan(bits)
	case bits&0x20 > 0:
		// 拔北
		panic("[_parseTenhouMeld] 代码有误")
	default:
		return d._parseKan(bits)
	}
}

func (*tenhouRoundData) isRedFive(tenhouTile int) bool {
	return tenhouTile == redFiveMan || tenhouTile == redFivePin || tenhouTile == redFiveSou
}

func (d *tenhouRoundData) containRedFive(tenhouTiles []int) bool {
	for _, tenhouTile := range tenhouTiles {
		if d.isRedFive(tenhouTile) {
			return true
		}
	}
	return false
}

func (d *tenhouRoundData) GetDataSourceType() int {
	return dataSourceTypeTenhou
}

func (d *tenhouRoundData) GetSelfSeat() int {
	return -1
}

func (d *tenhouRoundData) GetMessage() string {
	return d.originJSON
}

func (d *tenhouRoundData) SkipMessage() bool {
	// 注意：即使没有获取到用户名也能正常进行游戏
	return false
}

func (d *tenhouRoundData) IsLogin() bool {
	// TODO: 重连时要填入 gameConf.currentActiveTenhouUsername
	return d.msg.Tag == "HELO"
}

func (d *tenhouRoundData) HandleLogin() {
	username, err := url.QueryUnescape(d.msg.UserName)
	if err != nil {
		h.logError(err)
	}
	if username != gameConf.currentActiveTenhouUsername {
		color.HiGreen("%s 登录成功", username)
		gameConf.currentActiveTenhouUsername = username
	}
}

func (d *tenhouRoundData) IsInit() bool {
	return d.msg.Tag == "INIT" || d.msg.Tag == "REINIT"
}

func (d *tenhouRoundData) ParseInit() (roundNumber int, benNumber int, dealer int, doraIndicators []int, handTiles []int, numRedFives []int) {
	d.isRoundEnd = false

	seedSplits := strings.Split(d.msg.Seed, ",")
	if len(seedSplits) != 6 {
		panic(fmt.Sprintln("seed 解析失败", d.msg.Seed))
	}

	roundNumber, _ = strconv.Atoi(seedSplits[0])
	benNumber, _ = strconv.Atoi(seedSplits[1])
	// TODO: 重构至 core。parser 不要修改任何东西
	if roundNumber == 0 && benNumber == 0 {
		if util.InStrings("0", strings.Split(d.msg.Ten, ",")) {
			d.playerNumber = 3
		} else {
			d.playerNumber = 4
		}
	}

	dealer, _ = strconv.Atoi(d.msg.Dealer)
	doraIndicator, _ := d._parseTenhouTile(seedSplits[5])
	doraIndicators = append(doraIndicators, doraIndicator)
	numRedFives = make([]int, 3)
	tenhouTiles := strings.Split(d.msg.Hai, ",")
	for _, tenhouTile := range tenhouTiles {
		tile, isRedFive := d._parseTenhouTile(tenhouTile)
		handTiles = append(handTiles, tile)
		if isRedFive {
			numRedFives[tile/9]++
		}
	}
	return
}

var _selfDrawReg = regexp.MustCompile("^T[0-9]{1,3}$")

func isTenhouSelfDraw(tag string) bool {
	return _selfDrawReg.MatchString(tag)
}

func (d *tenhouRoundData) IsSelfDraw() bool {
	return isTenhouSelfDraw(d.msg.Tag)
}

func (d *tenhouRoundData) ParseSelfDraw() (tile int, isRedFive bool, kanDoraIndicator int) {
	rawTile := d.msg.Tag[1:]
	tile, isRedFive = d._parseTenhouTile(rawTile)
	kanDoraIndicator = -1
	return
}

var _discardReg = regexp.MustCompile("^[DEFGefg][0-9]{1,3}$")

func (d *tenhouRoundData) IsDiscard() bool {
	return _discardReg.MatchString(d.msg.Tag)
}

func (d *tenhouRoundData) ParseDiscard() (who int, discardTile int, isRedFive bool, isTsumogiri bool, isReach bool, canBeMeld bool, kanDoraIndicator int) {
	// D=自家, e/E=下家, f/F=对家, g/G=上家
	who = int(util.Lower(d.msg.Tag[0]) - 'd')
	rawTile := d.msg.Tag[1:]
	discardTile, isRedFive = d._parseTenhouTile(rawTile)
	if d.msg.Tag[0] != 'D' {
		isTsumogiri = d.msg.Tag[0] >= 'a'
		canBeMeld = d.msg.T != ""
	}
	kanDoraIndicator = -1
	return
}

func (*tenhouRoundData) isNukiOperator(data string) bool {
	bits, err := strconv.Atoi(data)
	if err != nil {
		panic(err)
	}
	return bits&0x4 == 0 && bits&0x18 == 0 && bits&0x20 > 0
}

func (d *tenhouRoundData) IsOpen() bool {
	if d.msg.Tag != "N" {
		return false
	}

	// 除去拔北
	return !d.isNukiOperator(d.msg.Meld)
}

func (d *tenhouRoundData) ParseOpen() (who int, meld *model.Meld, kanDoraIndicator int) {
	who, _ = strconv.Atoi(d.msg.Who)
	meldType, tenhouMeldTiles, tenhouCalledTile := d._parseTenhouMeld(d.msg.Meld)
	meldTiles := make([]int, len(tenhouMeldTiles))
	for i, tenhouTile := range tenhouMeldTiles {
		meldTiles[i] = d._tenhouTileToTile34(tenhouTile)
	}
	sort.Ints(meldTiles)
	calledTile := d._tenhouTileToTile34(tenhouCalledTile)
	isCalledTileRedFive := d.isRedFive(tenhouCalledTile)
	meld = &model.Meld{
		MeldType:          meldType,
		Tiles:             meldTiles,
		CalledTile:        calledTile,
		ContainRedFive:    d.containRedFive(tenhouMeldTiles),
		RedFiveFromOthers: isCalledTileRedFive && (meldType == model.MeldTypeChi || meldType == model.MeldTypePon || meldType == model.MeldTypeMinkan),
	}
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
	return d.msg.Tag == "AGARI"
}

func (d *tenhouRoundData) ParseRoundWin() (whos []int, points []int) {
	d.isRoundEnd = true

	who, _ := strconv.Atoi(d.msg.Who)
	splits := strings.Split(d.msg.Ten, ",")
	if len(splits) < 2 {
		return
	}
	point, _ := strconv.Atoi(splits[1])
	return []int{who}, []int{point}
}

func (d *tenhouRoundData) IsRyuukyoku() bool {
	return d.msg.Tag == "RYUUKYOKU"
}

// "{\"tag\":\"RYUUKYOKU\",\"type\":\"ron3\",\"ba\":\"1,1\",\"sc\":\"290,0,228,0,216,0,256,0\",\"hai0\":\"18,19,30,32,33,41,43,94,95,114,115,117,119\",\"hai2\":\"29,31,74,75\",\"hai3\":\"8,13,17,25,35,46,48,53,78,79\"}"
func (d *tenhouRoundData) ParseRyuukyoku() (type_ int, whos []int, points []int) {
	d.isRoundEnd = true

	// TODO
	return
}

func (d *tenhouRoundData) IsNukiDora() bool {
	if d.msg.Tag != "N" {
		return false
	}

	return d.isNukiOperator(d.msg.Meld)
}

func (d *tenhouRoundData) ParseNukiDora() (who int, isTsumogiri bool) {
	// TODO: isTsumogiri
	who, _ = strconv.Atoi(d.msg.Who)
	return
}

func (d *tenhouRoundData) IsNewDora() bool {
	return d.msg.Tag == "DORA"
}

func (d *tenhouRoundData) ParseNewDora() (kanDoraIndicator int) {
	kanDoraIndicator, _ = d._parseTenhouTile(d.msg.Hai)
	return
}
