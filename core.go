package main

import (
	"fmt"
	"github.com/fatih/color"
)

var debugMode = false

const (
	dataSourceTypeTenhou = iota
	dataSourceTypeMajsoul
)

const (
	meldTypeChi    = iota // 吃
	meldTypePon           // 碰
	meldTypeAnKan         // 暗杠
	meldTypeMinKan        // 明杠
	meldTypeKakan         // 加杠
)

type DataParser interface {
	GetDataSourceType() int

	GetMessage() string

	// round 开始/重连
	// roundNumber: 场数（如东1为0，东2为1，...，南1为4，...）
	// dealer: 庄家 0-3
	// doraIndicator: 宝牌指示牌 0-33
	// hands: 手牌 [0-33]
	IsInit() bool
	ParseInit() (roundNumber int, dealer int, doraIndicator int, hands []int)

	// 自家摸牌
	// tile: 0-33
	// kanDoraIndicator: 摸牌时，若为暗杠摸的岭上牌，则可以翻出杠宝牌指示牌，否则返回 -1 （天凤恒为 -1，见 IsNewDora）
	IsSelfDraw() bool
	ParseSelfDraw() (tile int, kanDoraIndicator int)

	// 舍牌
	// who: 0=自家, 1=下家, 2=对家, 3=上家
	// isTsumogiri: 是否为摸切（who=0 时忽略该值）
	// isReach: 是否为立直宣言（isReach 对于天凤来说恒为 false，见 IsReach）
	// canBeMeld: 是否可以鸣牌（who=0 时忽略该值）
	// kanDoraIndicator: 明杠/加杠的杠宝牌指示牌，在切牌后出现，没有则返回 -1（天凤恒为-1）
	IsDiscard() bool
	ParseDiscard() (who int, tile int, isTsumogiri bool, isReach bool, canBeMeld bool, kanDoraIndicator int)

	// 鸣牌（含暗杠、加杠）
	// meldType: 鸣牌类型（吃、碰、暗杠、明杠、加杠）
	// meldTiles: 副露的牌 [0-33]
	// calledTile: 被鸣的牌 0-33
	// kanDoraIndicator: 暗杠的杠宝牌指示牌，在他家暗杠时出现，没有则返回 -1（天凤恒为-1）
	IsOpen() bool
	ParseOpen() (who int, meldType int, meldTiles []int, calledTile int, kanDoraIndicator int)

	// 立直声明（IsReach 对于雀魂来说恒为 false，见 ParseDiscard）
	IsReach() bool
	ParseReach() (who int)

	// 振听
	IsFuriten() bool

	// 杠宝牌（IsNewDora 对于雀魂来说恒为 false，见 ParseSelfDraw ParseDiscard ParseOpen）
	// kanDoraIndicator: 0-33
	IsNewDora() bool
	ParseNewDora() (kanDoraIndicator int)
}

//

type playerInfo struct {
	name string // 自家 下家 对家 上家

	selfWindTile int

	// 副露，鸣牌时的舍牌
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

func newPlayerInfo(name string, selfWindTile int, globalDiscardTiles *[]int) *playerInfo {
	return &playerInfo{
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
	// - 其余可以参考贴吧的《魔神之眼》翻译 https://tieba.baidu.com/p/3311909701
	//      举个简单的例子,如果出现手切了一个对子的情况的话那么基本上就不可能是七对子。
	//      如果对方早巡手切了一个两面搭子的话，那么就可以推理出他在做染手或者牌型是对子型，如果他立直或者鸣牌的话，也比较容易读出他的手牌。
	// https://tieba.baidu.com/p/3311909701
	//      鸣牌之后和终盘的手切牌要尽量记下来，别人手切之前的安牌应该先切掉
	// https://tieba.baidu.com/p/3372239806
	//      吃牌时候打出来的牌的颜色是危险的；碰之后全部的牌都是危险的

	fmt.Printf(p.name + ":")
	for i, disTile := range p.discardTiles {
		fmt.Printf(" ")
		// TODO: 显示 dora, 赤宝牌
		if disTile >= 0 { // 手切
			if len(p.melds) == 0 { // 未副露
				if disTile >= 27 {
					// 关注字牌的手切
					fmt.Printf(mahjongU[disTile])
				} else {
					fmt.Printf(mahjong[disTile])
				}
			} else { // 副露
				// 高亮中张和字牌的手切
				c := color.New(getDiscardAlertColor(disTile))
				if inInts(i, p.meldDiscardsAt) {
					// 鸣牌时切的那张牌要大写
					c.Printf(mahjongU[disTile])
				} else {
					c.Printf(mahjong[disTile])
				}
			}
		} else { // 摸切
			fmt.Printf("--")
		}
	}
	fmt.Println()
}

//

type roundData struct {
	parser DataParser

	// 场数（如东1为0，东2为1，...，南1为4，...）
	roundNumber int

	// 场风
	roundWindTile int

	// 庄家 0=自家, 1=下家, 2=对家, 3=上家
	dealer int

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
}

func newRoundData(parser DataParser, roundNumber int, dealer int) *roundData {
	roundWindTile := 27 + roundNumber/4
	playerWindTile := make([]int, 4)
	for i := 0; i < 4; i++ {
		playerWindTile[i] = 27 + (4-dealer+i)%4
	}
	leftCounts := make([]int, 34)
	for i := range leftCounts {
		leftCounts[i] = 4
	}
	globalDiscardTiles := []int{}
	return &roundData{
		parser:             parser,
		roundNumber:        roundNumber,
		roundWindTile:      roundWindTile,
		dealer:             dealer,
		counts:             make([]int, 34),
		leftCounts:         leftCounts,
		globalDiscardTiles: globalDiscardTiles,
		players: [4]*playerInfo{
			newPlayerInfo("自家", playerWindTile[0], &globalDiscardTiles),
			newPlayerInfo("下家", playerWindTile[1], &globalDiscardTiles),
			newPlayerInfo("对家", playerWindTile[2], &globalDiscardTiles),
			newPlayerInfo("上家", playerWindTile[3], &globalDiscardTiles),
		},
	}
}

func (d *roundData) reset(roundNumber int, dealer int) {
	newData := newRoundData(d.parser, roundNumber, dealer)
	*d = *newData
}

// TODO: 临时用
func (d *roundData) _fillZi() {
	for i, c := range d.counts[27:] {
		if c == 0 {
			d.counts[i+27] = 3
			break
		}
	}
}

func (d *roundData) descLeftCounts(tile int) {
	d.leftCounts[tile]--
	if d.leftCounts[tile] < 0 {
		fmt.Printf("数据异常: %s 数量为 %d\n", mahjongZH[tile], d.leftCounts[tile])
	}
}

func (d *roundData) newDora(kanDoraIndicator int) {
	color.Yellow("杠宝牌指示牌是 %s", mahjongZH[kanDoraIndicator])
	d.doraIndicators = append(d.doraIndicators, kanDoraIndicator)
	d.descLeftCounts(kanDoraIndicator)
}

func (d *roundData) printDiscards() {
	for i := 3; i > 0; i-- {
		d.players[i].printDiscards()
	}
}

// 分析34种牌的危险度，可以用来判断自家手牌的安全度，以及他家是否在进攻（多次切出危险度高的牌）
func (d *roundData) analysisTilesRisk() (tables riskTables) {
	tables = make([]riskTable, 3)
	for who, player := range d.players[1:] {
		// TODO: 对于副露者，根据他的副露情况、手切数、巡目计算其听牌率
		// TODO: 若某人一直摸切，然后突然手切了一张字牌，那他很有可能在默听，或者进入了完全一向听
		// 目前暂时简化成「三副露=听牌，晚巡两副露=听牌」（暗杠算副露）
		if !player.isReached && (len(player.melds) < 2 || len(player.melds) == 2 && len(player.discardTiles) < 13) {
			continue
		}

		// 该玩家的巡目 = 为其切过的牌的数目
		turns := minInt(len(player.discardTiles), 19)
		if turns == 0 {
			continue
		}

		// 收集安牌
		safeTiles := make([]uint8, 34)
		for _, tile := range player.discardTiles {
			// 该玩家的舍牌
			if tile < 0 {
				tile = ^tile
			}
			safeTiles[tile] = 1
		}
		if player.reachTileAtGlobal != -1 {
			// 立直后其他家切出的牌
			for _, tile := range d.globalDiscardTiles[player.reachTileAtGlobal:] {
				if tile < 0 {
					tile = ^tile
				}
				safeTiles[tile] = 1
			}
		} else {
			// TODO: 副露者三副露之后，其上家的舍牌大概率是安牌

			// https://tieba.baidu.com/p/3418094524
			// 副露家的上家的舍牌是重要的提示，副露家不鸣的牌也可以成为读牌的线索：
			// ① 副露家的上家切过的牌高概率能通过
			// ② 对于上一巡被切出来的牌（a）无反应，然后这一巡鸣牌后打牌（a）的情况，牌（a）的跨筋比较安全。
			// 举个例子，对于被切出来的7p毫无反应的对手，34s鸣2s后打7p。假定他听69p或者58p，那么之前的形状就是778p切7p和677p切7p，这样的话，7p被打出来的时候就应该被碰了，所以不会是听69p或者58p。
			// 顺带一提，这种情况并不限于鸣牌打7p的场合，其实在普通的手切7p的场合也是可以通用的。如果拿着677p或者778p这样的搭子的话，7p被切出来的时候就应该鸣了，如果是拿着67p或者78p的话，摸7p也不会特意手切一张7p来让别人注意防守7p的周边。（但是有的人可能会故意这样切牌，所以还是需要注意一下的）。
			// ③ 副露家鸣牌之后将上家切过的牌的周边牌切出来了
			// 鸣牌家的东家没有鸣北家切的7s，然后碰了南家的8p之后切8s。顺带一提，上家碰白打8m，吃4m打3m，3s是手切的。
			// 这样的例子从舍牌和副露看，并不能推理出他的待牌，但是上家没有鸣7s是一个线索，而且这个线索十分关键。
			// 东家2巡前切3s，上一巡摸切北，然后打的是8s。重视孤立牌靠张的话应该留3s，如果是需要安全牌的话应该留北才对，所以留8s的原因是他手里有和8s相关的搭子。
			// 然后我们知道他没鸣7s，而且和8s有关又鸣不了7s的搭子只有78s和788s（和8s有关的搭子有468s，688s，668s，68s，778s，788s，78s，889s，899s，89s）。仔细想一下的话，如果他拿着78s的搭子就不会特意鸣8p变成7s单骑了，所以能够推断出他鸣8p之前手里的搭子是788s。
			// 由此可知，东家是788s碰8p打8s听69s。像这样副露家不鸣哪一些牌也是一条挺重要的线索，所以请大家打牌的时候务必注意一下。

			// 空切·振替 https://tieba.baidu.com/p/3471413696
			//

			// 食延的情况 https://tieba.baidu.com/p/3688516724

		}

		table := make([]float64, 34)

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

func (d *roundData) analysis() error {
	if !debugMode {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("内部错误：", err)
			}
		}()
	}

	if debugMode {
		msg := d.parser.GetMessage()
		if d.parser.GetDataSourceType() == dataSourceTypeMajsoul {
			if len(msg) == 7 {
				return nil
			}
			fmt.Printf("收到 [len=%d] %s\n", len(msg), msg)
			fmt.Println(stringToSlice(msg))
			fmt.Println([]byte(msg))
		} else {
			fmt.Println("收到", msg)
		}
	}

	// 若自家立直，则进入看戏模式
	// TODO: 见逃判断
	if !d.parser.IsInit() && d.players[0].isReached {
		return nil
	}

	switch {
	case d.parser.IsInit():
		// round 开始/重连
		if !debugMode {
			clearConsole()
		}

		roundNumber, dealer, doraIndicator, hands := d.parser.ParseInit()
		switch d.parser.GetDataSourceType() {
		case dataSourceTypeTenhou:
			d.reset(roundNumber, dealer)
		case dataSourceTypeMajsoul:
			if dealer != -1 {
				d.dealer = dealer

				fmt.Printf("游戏即将开始，您分配到的座位是：")
				windTile := 27 + (4-dealer)%4
				color.Yellow(mahjongZH[windTile])

				return nil
			} else {
				dealer = d.dealer
				if roundNumber > 0 && roundNumber != d.roundNumber {
					dealer = (dealer + 1) % 4
				}
				d.reset(roundNumber, dealer)
			}
		default:
			panic("not impl!")
		}

		fmt.Printf("%s%d局开始，自风为%s\n", mahjongZH[d.roundWindTile], roundNumber%4+1, mahjongZH[d.players[0].selfWindTile])

		color.Yellow("宝牌指示牌是 %s", mahjongZH[doraIndicator])
		d.doraIndicators = []int{doraIndicator}
		d.descLeftCounts(doraIndicator)

		for _, tile := range hands {
			d.counts[tile]++
			d.descLeftCounts(tile)
		}
	case d.parser.IsOpen():
		// 某家鸣牌（含暗杠、加杠）
		who, meldType, meldTiles, calledTile, kanDoraIndicator := d.parser.ParseOpen()
		if kanDoraIndicator != -1 {
			d.newDora(kanDoraIndicator)
		}
		if meldType == meldTypeKakan {
			// TODO: 修改副露情况
			d.descLeftCounts(calledTile)
			break
		}

		// TODO: 添加 calledTile 等
		d.players[who].melds = append(d.players[who].melds, meldTiles)
		if who != 0 {
			// 处理牌山剩余量
			if meldType != meldTypeAnKan {
				d.leftCounts[calledTile]++
			}
			for _, tile := range meldTiles {
				d.descLeftCounts(tile)
			}
		}

		if who == 0 {
			// 自家副露
			// 简化，修改副露牌为字牌
			if meldType == meldTypeAnKan {
				d.counts[meldTiles[0]] = 0
			} else {
				d.counts[calledTile]++
				for _, tile := range meldTiles {
					d.counts[tile]--
				}
			}
			d._fillZi()
		}
	case d.parser.IsNewDora():
		// 杠宝牌
		// 1. 剩余牌减少
		// 2. 打点提高
		kanDoraIndicator := d.parser.ParseNewDora()
		d.newDora(kanDoraIndicator)
	case d.parser.IsReach():
		// 立直宣告
		// 如果是他家立直，进入攻守判断模式
		who := d.parser.ParseReach()
		d.players[who].isReached = true
		//case "AGARI", "RYUUKYOKU":
		//	// 某人和牌或流局，round 结束
		//case "PROF":
		//	// 游戏结束
		//case "BYE":
		//	// 某人退出
		//case "REJOIN", "GO":
		//	// 重连
	case d.parser.IsFuriten():
		// 振听
		color.Yellow("振听")
		//case "U", "V", "W":
		//	//（下家,对家,上家 不要其上家的牌）摸牌
		//case "HELO", "RANKING", "TAIKYOKU", "UN", "LN", "SAIKAI":
		//	// 其他
	case d.parser.IsSelfDraw():
		if !debugMode {
			clearConsole()
		}
		// 自家（从牌山 d.leftCounts）摸牌（至手牌 d.counts）
		// FIXME: 对于天凤，有一定概率在自己坐庄时，会先收到摸牌的消息，然后收到本局开始的消息
		tile, kanDoraIndicator := d.parser.ParseSelfDraw()
		d.descLeftCounts(tile)
		d.counts[tile]++
		if kanDoraIndicator != -1 {
			d.newDora(kanDoraIndicator)
		}

		// 打印他家舍牌信息
		d.printDiscards()
		fmt.Println()

		// 安全度分析
		riskTables := d.analysisTilesRisk()
		riskTables.printWithHands(d.counts, d.leftCounts)

		// 何切
		// TODO: 根据是否听牌/一向听、打点、巡目、和率等进行攻守判断
		return _analysis(14, d.counts, d.leftCounts)
	case d.parser.IsDiscard():
		who, tile, isTsumogiri, isReach, canBeMeld, kanDoraIndicator := d.parser.ParseDiscard()

		if kanDoraIndicator != -1 {
			d.newDora(kanDoraIndicator)
		}

		player := d.players[who]
		if isReach {
			player.isReached = isReach
		}

		if who == 0 {
			// 自家（从手牌 d.counts）舍牌（至牌河 d.globalDiscardTiles）
			d.counts[tile]--

			d.globalDiscardTiles = append(d.globalDiscardTiles, tile)
			player.discardTiles = append(player.discardTiles, tile)

			return nil
		}

		// 他家舍牌
		d.descLeftCounts(tile)

		if who != 3 {
			// 为防止先收到自家摸牌，然后收到上家摸牌，上家舍牌时不刷新
			if !debugMode {
				clearConsole()
			}
		}

		disTile := tile
		if isTsumogiri {
			disTile = ^disTile
		}
		d.globalDiscardTiles = append(d.globalDiscardTiles, disTile)
		player.discardTiles = append(player.discardTiles, disTile)

		if player.isReached && player.reachTileAtGlobal == -1 {
			// 标记立直宣言牌
			player.reachTileAtGlobal = len(d.globalDiscardTiles) - 1
			player.reachTileAt = len(player.discardTiles) - 1

			// 若该玩家摸切立直，打印提示信息
			if isTsumogiri {
				color.Yellow("%s 摸切立直！", d.players[who].name)
			}
		} else if len(player.meldDiscardsAt) != len(player.melds) {
			// 标记鸣牌的舍牌
			if len(player.meldDiscardsAt)+1 != len(player.melds) {
				fmt.Printf("玩家数据异常 %#v", *player)
			}
			player.meldDiscardsAt = append(player.meldDiscardsAt, len(player.discardTiles)-1)
			player.meldDiscardsAtGlobal = append(player.meldDiscardsAtGlobal, len(d.globalDiscardTiles)-1)
		}

		if who != 3 {
			// 打印他家舍牌信息
			d.printDiscards()
			fmt.Println()

			// 安全度分析
			riskTables := d.analysisTilesRisk()
			riskTables.printWithHands(d.counts, d.leftCounts)
		}

		if canBeMeld { // 是否副露
			d.counts[tile]++

			// TODO: 消除海底/避免河底/型听提醒

			// 何切
			err := _analysis(14, d.counts, d.leftCounts)
			d.counts[tile]--
			return err
		}
	default:
	}

	return nil
}
