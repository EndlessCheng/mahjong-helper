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
	Seed   string `json:"seed"` // 本局信息：场数，连庄棒数，立直棒数，骰子A减一，骰子B减一，宝牌指示牌 1,0,0,3,2,92
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

type playerInfo struct {
	name string // 自家 下家 对家 上家

	selfWindTile int

	// 副露
	melds                [][]int
	meldDiscardsAtGlobal []int
	meldDiscardsAt       []int

	// 全局舍牌
	// 注意负数要^
	globalDiscardTiles *[]int
	discardTiles       []int

	isReached bool
	// 立直宣言牌在 globalDiscardTiles 中的下标，初始为 -1
	reachTileAtGlobal int
	reachTileAt       int
}

func newPlayerInfo(name string, selfWindTile int, globalDiscardTiles *[]int) playerInfo {
	return playerInfo{
		name:               name,
		selfWindTile:       selfWindTile,
		globalDiscardTiles: globalDiscardTiles,
		reachTileAtGlobal:  -1,
		reachTileAt:        -1,
	}
}

func (p *playerInfo) printDiscards() {
	// TODO: 高亮不合理的舍牌或危险舍牌，如
	// - 一开始就切中张
	// - 开始切中张后，手切了幺九牌（也有可能是有人碰了牌，比如 133m 有人碰了 2m）
	// - 切了 dora，提醒一下
	// - 切了赤宝牌
	// - 有人立直的情况下，多次切出危险度高的牌（有可能是对方读准了牌，或者对方手里的牌与牌河加起来产生了安牌）

	fmt.Printf(p.name + ":")
	for _, disTile := range p.discardTiles {
		fmt.Printf(" ")
		// TODO: 显示 dora, 赤宝牌
		// FIXME: 颜色需要调整
		if disTile >= 0 { // 手切
			if len(p.melds) == 0 {
				// 未副露
				if disTile >= 27 {
					color.New(color.FgYellow).Printf(mahjong[disTile])
				} else {
					fmt.Printf(mahjong[disTile])
				}
			} else {
				// 副露者高亮摸切
				color.New(getDiscardAlertColor(disTile)).Printf(mahjong[disTile])
			}
		} else { // 摸切
			fmt.Printf("--")
		}
	}
	fmt.Println()
}

//

type tenhouRoundData struct {
	msg *tenhouMessage

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
	players [4]playerInfo
}

func newTenhouRoundData(roundNumber int, dealer int) *tenhouRoundData {
	roundWindTile := 27 + roundNumber/4
	playerWindTile := make([]int, 4)
	for i := 0; i < 4; i++ {
		playerWindTile[i] = 27 + (4-dealer+i)%4
	}
	globalDiscardTiles := []int{}
	d := &tenhouRoundData{
		roundNumber:        roundNumber,
		roundWindTile:      roundWindTile,
		counts:             make([]int, 34),
		leftCounts:         make([]int, 34),
		globalDiscardTiles: globalDiscardTiles,
		players: [4]playerInfo{
			newPlayerInfo("自家", playerWindTile[0], &globalDiscardTiles),
			newPlayerInfo("下家", playerWindTile[1], &globalDiscardTiles),
			newPlayerInfo("对家", playerWindTile[2], &globalDiscardTiles),
			newPlayerInfo("上家", playerWindTile[3], &globalDiscardTiles),
		},
	}
	for i := range d.leftCounts {
		d.leftCounts[i] = 4
	}
	return d
}

func (d *tenhouRoundData) reset(roundNumber int, dealer int) {
	newData := newTenhouRoundData(roundNumber, dealer)
	*d = *newData
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

// 分析34种牌的危险度，可以用来判断自家手牌的安全度，以及他家是否在进攻（多次切出危险度高的牌）
func (d *tenhouRoundData) analysisTilesRisk() (tables riskTables) {
	tables = make([]riskTable, 3)
	for who, player := range d.players[1:] {
		// TODO: 对于副露者，根据他的副露情况、手切数、巡目计算其听牌率
		// TODO: 若某人一直摸切，然后突然切了一张字牌，那他很有可能听牌了（含默听）
		// 目前暂时简化成「三副露 = 立直」（暗杠算副露）
		if !player.isReached && len(player.melds) < 3 {
			continue
		}

		table := make([]float64, 34)

		// 该玩家的巡目 = 为其切过的牌的数目
		turns := minInt(len(player.discardTiles), 19)
		if turns == 0 {
			continue
		}

		// 收集安牌
		safeTiles := make([]uint8, 34)
		for _, tile := range player.discardTiles {
			if tile < 0 {
				tile = ^tile
			}
			safeTiles[tile] = 1
		}
		if player.reachTileAtGlobal != -1 {
			for _, tile := range d.globalDiscardTiles[player.reachTileAtGlobal:] {
				if tile < 0 {
					tile = ^tile
				}
				safeTiles[tile] = 1
			}
		}

		// 利用安牌计算双筋、筋、半筋、无筋等
		// TODO: 单独处理宣言牌的筋牌、宣言牌的同色牌的危险度
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				t := tileTypeTable[j][safeTiles[9*i+j+3]]
				table[9*i+j] = riskData[turns][t]
			}
			for j := 3; j < 6; j++ {
				mixSafeTile := safeTiles[9*i+j-3]<<1 | safeTiles[9*i+j+3]
				t := tileTypeTable[j][mixSafeTile]
				table[9*i+j] = riskData[turns][t]
			}
			for j := 6; j < 9; j++ {
				t := tileTypeTable[j][safeTiles[9*i+j-3]]
				table[9*i+j] = riskData[turns][t]
			}
		}
		for i := 27; i < 34; i++ {
			if d.leftCounts[i] > 0 {
				isYakuHai := i == d.roundWindTile || i == player.selfWindTile || i >= 31
				t := ziTileType[boolToInt(isYakuHai)][d.leftCounts[i]-1]
				table[i] = riskData[turns][t]
			}
		}

		// 利用剩余牌是否为 0 或者 1 计算 No Chance, One Chance, Double One Chance, Double Two Chance(待定) 等
		// 利用舍牌计算无筋早外
		//（待定）有早外的半筋（早巡打过8m时，3m的半筋6m）
		//（待定）利用赤宝牌计算危险度
		// 宝牌周边牌的危险度要增加一点
		//（待定）切过5的情况

		for i, isSafe := range safeTiles {
			if isSafe == 1 {
				table[i] = 0
			}
		}

		tables[who] = table
	}
	return tables
}

func (d *tenhouRoundData) analysis() error {
	//defer func() {
	//	if err := recover(); err != nil {
	//		fmt.Println("内部错误：", err)
	//	}
	//}()

	msg := d.msg
	//fmt.Println("收到", msg.Tag)

	// 若自家立直，则进入看戏模式
	if msg.Tag != "INIT" && msg.Tag != "REINIT" && d.players[0].isReached {
		return nil
	}

	switch msg.Tag {
	case "INIT", "REINIT":
		// round 开始/重连
		clearConsole()
		splits := strings.Split(msg.Seed, ",")
		if len(splits) != 6 {
			panic(fmt.Sprintln("seed 解析失败", msg.Seed))
		}
		roundNumber, _ := strconv.Atoi(splits[0])
		dealer, _ := strconv.Atoi(msg.Dealer)
		d.reset(roundNumber, dealer)

		fmt.Printf("%s%d局开始，自风为%s\n", mahjongZH[d.roundWindTile], roundNumber%4+1, mahjongZH[d.players[0].selfWindTile])

		doraIndicator := d._parseTenhouTile(splits[5])
		color.Yellow("宝牌指示牌是 %s", mahjongZH[doraIndicator])
		d.doraIndicators = []int{doraIndicator}
		d.leftCounts[doraIndicator]--

		for _, pai := range strings.Split(msg.Hai, ",") {
			tile := d._parseTenhouTile(pai)
			d.counts[tile]++
			d.leftCounts[tile]--
		}
		//fmt.Println("剩余", d.leftCounts)
	case "N":
		// 某人已副露
		who, _ := strconv.Atoi(msg.Who)
		meldType, meldTiles, calledTile := d._parseTenhouMeld(msg.Meld)
		if meldType == meldTypeKakan {
			// TODO: 修改副露情况
			d.leftCounts[calledTile]--
			if d.leftCounts[calledTile] != 0 {
				fmt.Printf("数据异常: %s 数量为 %d\n", mahjongZH[calledTile], d.leftCounts[calledTile])
			}
			break
		}

		// TODO: 添加 calledTile
		d.players[who].melds = append(d.players[who].melds, meldTiles)
		// FIXME: 处理他家暗杠的情况
		if who != 0 {
			d.leftCounts[calledTile]++
			for _, tile := range meldTiles {
				d.leftCounts[tile]--
			}
		}

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

		//fmt.Println("剩余", d.leftCounts)
	case "DORA":
		// 杠宝牌
		// 1. 能摸的牌减少
		// 2. 打点提高
		kanDoraIndicator := d._parseTenhouTile(msg.Hai)
		color.Yellow("杠宝牌指示牌是 %s", mahjongZH[kanDoraIndicator])
		d.doraIndicators = append(d.doraIndicators, kanDoraIndicator)
		d.leftCounts[kanDoraIndicator]--
		//fmt.Println("剩余", d.leftCounts)
	case "REACH":
		// 如果是他家立直，进入攻守判断模式
		if msg.Step == "1" {
			// 立直宣告
			who, _ := strconv.Atoi(msg.Who)
			d.players[who].isReached = true
		} else {
			// (立直成功，扣1000点)
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
	case "HELO", "RANKING", "TAIKYOKU", "UN", "LN", "SAIKAI":
		// 其他
	default:
		rawTile := msg.Tag[1:]
		tile := d._parseTenhouTile(rawTile)
		switch msg.Tag[0] {
		case 'T':
			clearConsole()
			// 自家（从牌山 d.leftCounts）摸牌（至手牌 d.counts）
			d.leftCounts[tile]--
			//fmt.Println("剩余", d.leftCounts)
			d.counts[tile]++

			// 打印他家舍牌信息
			for _, player := range d.players[1:] {
				player.printDiscards()
			}
			fmt.Println()

			// 安全度分析
			riskTables := d.analysisTilesRisk()
			riskTables.printWithHands(d.counts, d.leftCounts)

			// 何切
			// TODO: 根据是否听牌/完全一向听、和牌率、打点、巡目进行攻守判断
			return _analysis(14, d.counts, d.leftCounts)
		case 'D':
			// 自家（从手牌 d.counts）舍牌（至牌河 d.globalDiscardTiles）
			d.counts[tile]--

			d.globalDiscardTiles = append(d.globalDiscardTiles, tile)
			d.players[0].discardTiles = append(d.players[0].discardTiles, tile)
		case 'E', 'F', 'G', 'e', 'f', 'g':
			// 他家舍牌, e=下家, f=对家, g=上家
			clearConsole()
			d.leftCounts[tile]--
			//fmt.Println("剩余", d.leftCounts)

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

				// 若该玩家摸切立直，打印提示信息
				if isTsumogiri {
					color.Yellow("%s 摸切立直！", d.players[who].name)
				}
			}

			// 打印他家舍牌信息
			for _, player := range d.players[1:] {
				player.printDiscards()
			}
			fmt.Println()

			if msg.T != "" { // 是否副露
				d.counts[tile]++

				// TODO: 消除海底/避免河底/型听提醒

				// 安全度分析
				riskTables := d.analysisTilesRisk()
				riskTables.printWithHands(d.counts, d.leftCounts)

				// 何切
				err := _analysis(14, d.counts, d.leftCounts)
				d.counts[tile]--
				return err
			}
		default:
		}
	}

	return nil
}
