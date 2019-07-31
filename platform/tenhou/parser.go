package tenhou

import (
	"regexp"
	"github.com/EndlessCheng/mahjong-helper/platform/tenhou/ws"
	"encoding/json"
	"strconv"
	"reflect"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"github.com/EndlessCheng/mahjong-helper/platform/common"
	"sort"
)

type message struct {
	originJSON string
	tag        string
	metadata   ws.Message
}

// T=自家, U=下家, V=对家, W=上家
var isDraw = regexp.MustCompile("^[TUVW][0-9]{0,3}$").MatchString

func isSelfDraw(tag string) bool {
	return isDraw(tag) && len(tag) > 1
}

// D/d=自家, E/e=下家, F/f=对家, G/g=上家
// 大写为手切，小写为摸切
// 对战模式下自家一律为 D
var isDiscard = regexp.MustCompile("^[DEFGdefg][0-9]{1,3}$").MatchString

func parse(data []byte) (msg *message, err error) {
	d := struct {
		Tag string `json:"tag"`
		Op  *int   `json:"t,string"`
	}{}
	if err = json.Unmarshal(data, &d); err != nil {
		return
	}
	tag := d.Tag
	msg = &message{
		originJSON: string(data),
		tag:        tag,
	}

	if isDraw(tag) {
		tile := -1
		if len(tag) > 1 {
			tile, err = strconv.Atoi(tag[1:])
			if err != nil {
				return
			}
		}
		msg.metadata = &ws.Draw{
			Who:  int(tag[0] - 'T'),
			Tile: tile,
			Op:   d.Op,
		}
		return
	}

	if isDiscard(tag) {
		who := int(tag[0] - 'D')
		tile, er := strconv.Atoi(tag[1:])
		if er != nil {
			return nil, er
		}
		isTsumogiri := tag[0] > 'a'
		if isTsumogiri {
			who = int(tag[0] - 'd')
		}
		msg.metadata = &ws.Discard{
			Who:         who,
			Tile:        tile,
			IsTsumogiri: isTsumogiri,
			Op:          d.Op,
		}
		return
	}

	mt := ws.MessageType(tag)
	if mt == nil {
		return
	}
	messagePtr := reflect.New(mt.Elem())
	if err = json.Unmarshal(data, messagePtr.Interface().(ws.Message)); err != nil {
		return
	}
	msg.metadata = messagePtr.Interface().(ws.Message)

	return
}

// 0-35 m
// 36-71 p
// 72-107 s
// 108- z
const (
	redFiveMan = 16
	redFivePin = 52
	redFiveSou = 88
)

func (*message) isRedFive(tenhouTile int) bool {
	return tenhouTile == redFiveMan || tenhouTile == redFivePin || tenhouTile == redFiveSou
}

func (m *message) containRedFive(tenhouTiles []int) bool {
	for _, tenhouTile := range tenhouTiles {
		if m.isRedFive(tenhouTile) {
			return true
		}
	}
	return false
}

func (*message) tenhouTileToTile34(tenhouTile int) int {
	return tenhouTile / 4
}

func (m *message) parseTenhouTile(tenhouTile int) (tile int, isRedFive bool) {
	return m.tenhouTileToTile34(tenhouTile), m.isRedFive(tenhouTile)
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
func (*message) parseChi(bits int) (meldType int, tenhouMeldTiles []int, tenhouCalledTile int) {
	// 吃
	meldType = common.MeldTypeChi
	t0, t1, t2 := (bits>>3)&0x3, (bits>>5)&0x3, (bits>>7)&0x3
	baseAndCalled := bits >> 10
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
func (*message) parsePonOrKakan(bits int) (meldType int, tenhouMeldTiles []int, tenhouCalledTile int) {
	t4 := (bits >> 5) & 0x3
	_t := [4][3]int{{1, 2, 3}, {0, 2, 3}, {0, 1, 3}, {0, 1, 2}}[t4]
	t0, t1, t2 := _t[0], _t[1], _t[2]
	baseAndCalled := bits >> 9
	base, called := baseAndCalled/3, baseAndCalled%3
	if bits&0x8 > 0 {
		// 碰
		meldType = common.MeldTypePon
		tenhouMeldTiles = []int{t0 + 4*base, t1 + 4*base, t2 + 4*base}
		tenhouCalledTile = tenhouMeldTiles[called]
	} else { // bits&0x16 > 0
		// 加杠
		meldType = common.MeldTypeKakan
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
func (*message) parseKan(bits int) (meldType int, tenhouMeldTiles []int, tenhouCalledTile int) {
	baseAndCalled := bits >> 8
	base, called := baseAndCalled/4, baseAndCalled%4
	tenhouMeldTiles = []int{4 * base, 1 + 4*base, 2 + 4*base, 3 + 4*base}
	tenhouCalledTile = tenhouMeldTiles[called]
	if fromWho := bits & 0x3; fromWho == 0 {
		// 暗杠
		meldType = common.MeldTypeAnkan
	} else {
		// 大明杠
		meldType = common.MeldTypeMinkan
	}
	return
}

func (*message) isNukiOperator(bits int) bool {
	return bits&0x4 == 0 && bits&0x18 == 0 && bits&0x20 > 0
}

func (m *message) parseTenhouMeld(bits int) (meldType int, tenhouMeldTiles []int, tenhouCalledTile int) {
	switch {
	case bits&0x4 > 0:
		return m.parseChi(bits)
	case bits&0x18 > 0:
		return m.parsePonOrKakan(bits)
	case bits&0x20 > 0:
		// 拔北
		panic("[message.parseTenhouMeld] 代码有误！")
	default:
		return m.parseKan(bits)
	}
}

// TODO: 重构
func (m *message) GetDataSourceType() int {
	return 0
}

func (m *message) GetSelfSeat() int {
	return -1
}

func (m *message) GetMessage() string {
	return m.originJSON
}

func (m *message) SkipMessage() bool {
	return false
}

// TODO: remove this
func (m *message) IsLogin() bool {
	return false
}

// TODO: remove this
func (m *message) HandleLogin() {
}

func (m *message) IsInit() bool {
	_, ok := m.metadata.(*ws.Init)
	return ok
}

// TODO: 重构至 core。parser 不要修改任何东西
//if roundNumber == 0 && benNumber == 0 {
//	if util.InStrings("0", strings.Split(meta.Ten, ",")) {
//		d.playerNumber = 3
//	} else {
//		d.playerNumber = 4
//	}
//}
func (m *message) ParseInit() (roundNumber int, benNumber int, dealer int, doraIndicator int, handTiles []int, numRedFives []int) {
	meta := m.metadata.(*ws.Init)
	roundNumber = meta.Seed[0]
	benNumber = meta.Seed[1]
	dealer = meta.Dealer
	doraIndicator = m.tenhouTileToTile34(meta.Seed[5])
	numRedFives = make([]int, 3)
	for _, tenhouTile := range meta.Tiles {
		tile, isRedFive := m.parseTenhouTile(tenhouTile)
		handTiles = append(handTiles, tile)
		if isRedFive {
			numRedFives[tile/9]++
		}
	}
	return
}

func (m *message) IsSelfDraw() bool {
	meta, ok := m.metadata.(*ws.Draw)
	return ok && meta.Who == 0
}

func (m *message) ParseSelfDraw() (tile int, isRedFive bool, kanDoraIndicator int) {
	meta := m.metadata.(*ws.Draw)
	tile, isRedFive = m.parseTenhouTile(meta.Tile)
	kanDoraIndicator = -1
	return
}

func (m *message) IsDiscard() bool {
	_, ok := m.metadata.(*ws.Discard)
	return ok
}

func (m *message) ParseDiscard() (who int, discardTile int, isRedFive bool, isTsumogiri bool, isReach bool, canBeMeld bool, kanDoraIndicator int) {
	meta := m.metadata.(*ws.Discard)
	who = meta.Who
	discardTile, isRedFive = m.parseTenhouTile(meta.Tile)
	isTsumogiri = meta.IsTsumogiri
	canBeMeld = meta.Op != nil
	kanDoraIndicator = -1
	return
}

func (m *message) IsOpen() bool {
	meta, ok := m.metadata.(*ws.Meld)
	// 除去拔北
	return ok && !m.isNukiOperator(meta.Bits)
}

func (m *message) ParseOpen() (who int, meld *model.Meld, kanDoraIndicator int) {
	meta := m.metadata.(*ws.Meld)
	who = meta.Who
	meldType, tenhouMeldTiles, tenhouCalledTile := m.parseTenhouMeld(meta.Bits)
	meldTiles := make([]int, len(tenhouMeldTiles))
	for i, tenhouTile := range tenhouMeldTiles {
		meldTiles[i] = m.tenhouTileToTile34(tenhouTile)
	}
	sort.Ints(meldTiles)
	calledTile, isCalledTileRedFive := m.parseTenhouTile(tenhouCalledTile)
	meld = &model.Meld{
		MeldType:          meldType,
		Tiles:             meldTiles,
		CalledTile:        calledTile,
		ContainRedFive:    m.containRedFive(tenhouMeldTiles),
		RedFiveFromOthers: isCalledTileRedFive && (meldType == model.MeldTypeChi || meldType == model.MeldTypePon || meldType == model.MeldTypeMinkan),
	}
	kanDoraIndicator = -1
	return
}

func (m *message) IsRiichi() bool {
	meta, ok := m.metadata.(*ws.Riichi)
	return ok && meta.Step == 1
}

func (m *message) ParseRiichi() (who int) {
	meta := m.metadata.(*ws.Riichi)
	return meta.Who
}

func (m *message) IsRoundWin() bool {
	_, ok := m.metadata.(*ws.Agari)
	return ok
}

func (m *message) ParseRoundWin() (whos []int, points []int) {
	meta := m.metadata.(*ws.Agari)
	return []int{meta.Who}, []int{meta.Ten[1]}
}

func (m *message) IsRyuukyoku() bool {
	_, ok := m.metadata.(*ws.Ryuukyoku)
	return ok
}

func (m *message) ParseRyuukyoku() (type_ int, whos []int, points []int) {
	// TODO
	return
}

func (m *message) IsNukiDora() bool {
	meta, ok := m.metadata.(*ws.Meld)
	return ok && m.isNukiOperator(meta.Bits)
}

func (m *message) ParseNukiDora() (who int, isTsumogiri bool) {
	// TODO: isTsumogiri
	meta := m.metadata.(*ws.Meld)
	return meta.Who, true
}

func (m *message) IsNewDora() bool {
	_, ok := m.metadata.(*ws.Dora)
	return ok
}

func (m *message) ParseNewDora() (kanDoraIndicator int) {
	meta := m.metadata.(*ws.Dora)
	return m.tenhouTileToTile34(meta.Tile)
}
