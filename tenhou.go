package main

import (
	"strings"
	"strconv"
	"fmt"
	"github.com/fatih/color"
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
	Seed string `json:"seed"` // 本局信息：场数，连庄棒数，立直棒数，骰子A减一，骰子B减一，宝牌指示牌 1,0,0,3,2,92
	Ten  string `json:"ten"`  // 各家点数 280,230,240,250
	Oya  string `json:"oya"`  // 庄家 0=自家, 1=下家, 2=对家, 3=上家
	Hai  string `json:"hai"`  // 初始手牌 30,114,108,31,78,107,25,23,2,14,122,44,49

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

type playerInfo struct {
	name string // 自家 下家 对家 上家

	// 副露
	meldTiles            []int
	meldDiscardsAtGlobal []int
	meldDiscardsAt       []int

	// 全局舍牌
	globalDiscardTiles *[]int
	discardTiles       []int

	isReached bool
	// 立直宣言牌在 globalDiscardTiles 中的下标，初始为 -1
	reachTileAtGlobal int
	reachTileAt       int
}

func newPlayerInfo(name string, globalDiscardTiles *[]int) playerInfo {
	return playerInfo{
		name:               name,
		globalDiscardTiles: globalDiscardTiles,
		reachTileAtGlobal:  -1,
		reachTileAt:        -1,
	}
}

//

type tenhouRoundData struct {
	msg *tenhouMessage

	// 宝牌指示牌
	doraIndicators []int

	// 自家手牌
	counts []int

	// 全局舍牌
	// 按舍牌顺序，负数表示摸切(-)，非负数表示手切(+)
	// 可以理解成：- 表示不要/暗色，+ 表示进张/亮色
	globalDiscardTiles []int
	// 0=自家, 1=下家, 2=对家, 3=上家
	players [4]playerInfo
}

func newTenhouRoundData() *tenhouRoundData {
	globalDiscardTiles := []int{}
	return &tenhouRoundData{
		counts:             make([]int, 34),
		globalDiscardTiles: globalDiscardTiles,
		players: [4]playerInfo{
			newPlayerInfo("自家", &globalDiscardTiles),
			newPlayerInfo("下家", &globalDiscardTiles),
			newPlayerInfo("对家", &globalDiscardTiles),
			newPlayerInfo("上家", &globalDiscardTiles),
		},
	}
}

func (d *tenhouRoundData) reset() {
	newData := newTenhouRoundData()
	d.doraIndicators = newData.doraIndicators
	d.counts = newData.counts
	d.globalDiscardTiles = newData.globalDiscardTiles
	d.players = newData.players
}

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

const (
	meldTypeChi = iota
	meldTypePon
	meldTypeKan
	meldTypeKakan
)

func (*tenhouRoundData) _parseChi(data int) (meldType int, tiles []int, calledTile int) {
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

func (*tenhouRoundData) _parseKan(data int) (meldType int, tiles []int, calledTile int) {
	meldType = meldTypeKan
	baseAndCalled := data >> 8
	base, called := baseAndCalled/4, baseAndCalled%4
	tiles = []int{(4 * base) / 4, (1 + 4*base) / 4, (2 + 4*base) / 4, (3 + 4*base) / 4}
	calledTile = tiles[called]
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

// TODO: 临时用
func (d *tenhouRoundData) _fillZi() {
	for i, c := range d.counts[27:] {
		if c == 0 {
			d.counts[i+27] = 3
			break
		}
	}
}

// 分析各个牌的铳率
func (d *tenhouRoundData) analysisSafeTiles() map[int]float64 {
	table := map[int]float64{}
	for _, player := range d.players[1:] {
		// TODO: 根据该玩家的副露情况、手切数、巡目计算其听牌率

		if !player.isReached {
			continue
		}

		// 计算该玩家的巡目=切过的牌的数目
		_ = len(player.discardTiles)

		for _, tile := range player.discardTiles {
			// 安牌
			table[tile] = 0
		}
		for _, tile := range d.globalDiscardTiles[player.reachTileAtGlobal+1:] {
			// 安牌
			table[tile] = 0
		}

		// 双筋

		// 非立直宣言牌的筋牌

		// 立直宣言牌的筋牌

		// 壁 (No Chance)

		// Double One Chance

		// One Chance

		// ? Double Two Chance

		// 早外

		// 字牌

		// ? 有早外的半筋（早巡打过8m时，3m的半筋6m）

		// 半筋

		// 无筋

		// TODO: 多人立直的判断
		break
	}
	return table
}

func (d *tenhouRoundData) analysis() error {
	msg := d.msg
	fmt.Println("收到", msg.Tag)

	switch msg.Tag {
	case "INIT", "REINIT":
		// round 开始/重连
		fmt.Println("new round")
		d.reset()

		splits := strings.Split(msg.Seed, ",")
		if len(splits) != 5 {
			panic(fmt.Sprintln("seed 解析失败", msg.Seed))
		}
		doraIndicator := d._parseTenhouTile(splits[5])
		color.Yellow("宝牌指示牌是 %s", mahjongZH[doraIndicator])
		d.doraIndicators = []int{doraIndicator}

		for _, pai := range strings.Split(msg.Hai, ",") {
			tile := d._parseTenhouTile(pai)
			d.counts[tile]++
		}
		return _analysis(13, d.counts)
	case "N":
		// 某人已副露
		who, _ := strconv.Atoi(msg.Who)
		meldType, meldTiles, calledTile := d._parseTenhouMeld(msg.Meld)
		if meldType == meldTypeKakan {
			// TODO: 修改副露情况
			break
		}

		// TODO: 添加副露
		//d.players[who].meldTiles = append(d.players[who].meldTiles, meldTiles...)

		if who == 0 {
			// 简化，修改副露牌为字牌
			if meldType == meldTypeKan && d.counts[meldTiles[0]] == 4 { // 也可以判断手牌是否为 14 张
				// 暗杠
				d.counts[meldTiles[0]] = 0
			} else {
				d.counts[calledTile]++
				for _, tile := range meldTiles {
					d.counts[tile]--
				}
			}
			d._fillZi()
		}
	case "DORA":
		// 杠宝牌
		// 1. 能摸的牌减少
		// 2. 打点提高
		kanDoraIndicator := d._parseTenhouTile(msg.Hai)
		color.Yellow("杠宝牌指示牌是 %s", mahjongZH[kanDoraIndicator])
		d.doraIndicators = append(d.doraIndicators, kanDoraIndicator)
	case "REACH":
		// 如果是他家立直，进入攻守判断模式
		if msg.Step == "1" {
			// 立直宣告
			who, _ := strconv.Atoi(msg.Who)
			d.players[who].isReached = true
		} else {
			// 立直成功，扣1000点
		}
	case "AGARI", "RYUUKYOKU":
		// round 结束
	case "PROF":
		// 游戏结束
	case "BYE":
		// 某人退出
	case "REJOIN", "GO":
		// 重连
	case "FURITEN":
		// 振听
	case "U", "V", "W":
		//（下家,对家,上家 不要其上家的牌）摸牌
	case "TAIKYOKU", "UN", "LN":
		// 其他
	default:
		rawTile := msg.Tag[1:]
		tile := d._parseTenhouTile(rawTile)
		switch msg.Tag[0] {
		case 'T':
			// 自家摸牌

			// 他家舍牌信息
			for _, player := range d.players[1:] {
				fmt.Printf("%s:", player.name)
				for _, disTile := range player.discardTiles {
					if disTile >= 0 {
						// 手切
						fmt.Printf(mahjongZH[disTile] + " ")
					} else {
						// 摸切
						fmt.Printf("- ")
					}
				}
				fmt.Println()
			}

			// TODO: 若有危险牌信息，则排序后输出
			if dangerousTable := d.analysisSafeTiles(); len(dangerousTable) > 0 {

			}

			// 何切
			d.counts[tile]++
			return _analysis(14, d.counts)
		case 'D':
			// 自家舍牌
			d.globalDiscardTiles = append(d.globalDiscardTiles, tile)
			d.players[0].discardTiles = append(d.players[0].discardTiles, tile)

			d.counts[tile]--
			return _analysis(13, d.counts)
		case 'E', 'F', 'G', 'e', 'f', 'g':
			// 他家舍牌, e=下家, f=对家, g=上家
			who := lower(msg.Tag[0]) - 'd'
			isTsumogiri := msg.Tag[0] >= 'a' // 是否摸切

			disTile := tile
			if isTsumogiri {
				disTile = ^disTile
			}
			d.globalDiscardTiles = append(d.globalDiscardTiles, disTile)
			d.players[who].discardTiles = append(d.players[who].discardTiles, disTile)

			if d.players[who].isReached && d.players[who].reachTileAtGlobal == -1 {
				// 标记立直宣言牌
				d.players[who].reachTileAtGlobal = len(d.globalDiscardTiles) - 1
				d.players[who].reachTileAt = len(d.players[who].discardTiles) - 1

				if isTsumogiri {
					color.Yellow("%s 模切立直！", d.players[who].name)
				}
			}

			if msg.T != "" { // 是否副露
				// TODO: 消除海底/避免河底/型听提醒

				// TODO: 若有危险牌信息，则排序后输出
				if dangerousTable := d.analysisSafeTiles(); len(dangerousTable) > 0 {

				}

				// 何切
				d.counts[tile]++
				err := _analysis(14, d.counts)
				d.counts[tile]--
				return err
			}
		default:
		}
	}

	return nil
}
