package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/EndlessCheng/mahjong-helper/util"
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
	meldTypeMinKan        // 大明杠
	meldTypeKakan         // 加杠
)

type DataParser interface {
	GetDataSourceType() int

	GetMessage() string

	// 解析前，根据消息内容来决定是否要进行后续解析
	CheckMessage() bool

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
	// kanDoraIndicator: 大明杠/加杠的杠宝牌指示牌，在切牌后出现，没有则返回 -1（天凤恒为-1，见 IsNewDora）
	IsDiscard() bool
	ParseDiscard() (who int, tile int, isTsumogiri bool, isReach bool, canBeMeld bool, kanDoraIndicator int)

	// 鸣牌（含暗杠、加杠）
	// meldType: 鸣牌类型（吃、碰、暗杠、明杠、加杠）
	// meldTiles: 副露的牌 [0-33]
	// calledTile: 被鸣的牌 0-33
	// kanDoraIndicator: 暗杠的杠宝牌指示牌，在他家暗杠时出现，没有则返回 -1（天凤恒为-1，见 IsNewDora）
	IsOpen() bool
	ParseOpen() (who int, meldType int, meldTiles []int, calledTile int, kanDoraIndicator int)

	// 立直声明（IsReach 对于雀魂来说恒为 false，见 ParseDiscard）
	IsReach() bool
	ParseReach() (who int)

	// 振听
	IsFuriten() bool

	// 本局是否和牌
	IsRoundWin() bool
	ParseRoundWin() (whos []int, points []int)

	// 这一项放在末尾处理
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
		bgColor := color.BgBlack
		fgColor := color.FgWhite
		var tile string
		if disTile >= 0 { // 手切
			tile = util.Mahjong[disTile]
			if disTile >= 27 {
				tile = util.MahjongU[disTile] // 关注字牌的手切
			}
			if len(p.melds) == 0 { // 未副露
			} else { // 副露
				fgColor = getDiscardAlertColor(disTile) // 高亮中张手切
				if util.InInts(i, p.meldDiscardsAt) {
					bgColor = color.BgWhite // 鸣牌时切的那张牌要背景高亮
					fgColor = color.FgBlack
				}
			}
		} else { // 摸切
			disTile = ^disTile
			tile = util.Mahjong[disTile]
			fgColor = color.FgHiBlack // 暗色显示
		}
		color.New(bgColor, fgColor).Print(tile)
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
	counts    []int
	meldCount int

	// 牌山剩余牌量
	leftCounts []int

	// 全局舍牌
	// 按舍牌顺序，负数表示摸切(-)，非负数表示手切(+)
	// 可以理解成：- 表示不要/暗色，+ 表示进张/亮色
	globalDiscardTiles []int

	// 0=自家, 1=下家, 2=对家, 3=上家
	players []*playerInfo
}

func newRoundData(parser DataParser, roundNumber int, dealer int) *roundData {
	const playerNumber = 4
	roundWindTile := 27 + roundNumber/playerNumber
	playerWindTile := make([]int, playerNumber)
	for i := 0; i < playerNumber; i++ {
		playerWindTile[i] = 27 + (playerNumber-dealer+i)%playerNumber
	}
	globalDiscardTiles := []int{}
	return &roundData{
		parser:             parser,
		roundNumber:        roundNumber,
		roundWindTile:      roundWindTile,
		dealer:             dealer,
		counts:             make([]int, 34),
		leftCounts:         util.InitLeftTiles34(),
		globalDiscardTiles: globalDiscardTiles,
		players: []*playerInfo{
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

func (d *roundData) descLeftCounts(tile int) {
	d.leftCounts[tile]--
	if d.leftCounts[tile] < 0 {
		info := fmt.Sprintf("数据异常: %s 数量为 %d", util.MahjongZH[tile], d.leftCounts[tile])
		if debugMode {
			panic(info)
		} else {
			fmt.Println(info)
		}
	}
}

func (d *roundData) newDora(kanDoraIndicator int) {
	color.Yellow("杠宝牌指示牌是 %s", util.MahjongZH[kanDoraIndicator])
	d.doraIndicators = append(d.doraIndicators, kanDoraIndicator)
	d.descLeftCounts(kanDoraIndicator)
}

func (d *roundData) printDiscards() {
	for i := len(d.players) - 1; i >= 1; i-- {
		d.players[i].printDiscards()
	}
}

// 分析34种牌的危险度
// 可以用来判断自家手牌的安全度，以及他家是否在进攻（多次切出危险度高的牌）
func (d *roundData) analysisTilesRisk() (tables riskTables) {
	tables = make(riskTables, len(d.players))

	for who, player := range d.players {
		// TODO: 对于副露者，根据他的副露情况、手切数、巡目计算其听牌率
		// TODO: 若某人一直摸切，然后突然手切了一张字牌，那他很有可能默听/一向听
		// 目前暂时简化成「三副露=听牌，晚巡两副露=听牌」（暗杠算副露）
		if !player.isReached && (len(player.melds) < 2 || len(player.melds) == 2 && len(player.discardTiles) < 13) {
			continue
		}

		// 该玩家的巡目 = 为其切过的牌的数目
		turns := util.MinInt(len(player.discardTiles), util.MaxTurns)
		if turns == 0 {
			continue
		}

		// 收集安牌
		safeTiles34 := make([]bool, 34)
		for _, tile := range player.discardTiles {
			// 该玩家的舍牌
			if tile < 0 {
				tile = ^tile
			}
			safeTiles34[tile] = true
		}
		if player.reachTileAtGlobal != -1 {
			// 立直后其他家切出的牌
			for _, tile := range d.globalDiscardTiles[player.reachTileAtGlobal:] {
				if tile < 0 {
					tile = ^tile
				}
				safeTiles34[tile] = true
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

		risk34 := util.CalculateRiskTiles34(turns, safeTiles34, d.leftCounts, d.roundWindTile, player.selfWindTile)
		tables[who] = riskTable(risk34)
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
		fmt.Println("收到", d.parser.GetMessage())
	}

	if !d.parser.CheckMessage() {
		return nil
	}

	// 若自家立直，则进入看戏模式
	// TODO: 见逃判断
	if !d.parser.IsInit() && !d.parser.IsRoundWin() && d.players[0].isReached {
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
			playerNumber := len(d.players)
			if dealer != -1 {
				d.dealer = dealer

				fmt.Printf("游戏即将开始，您分配到的座位是：")
				windTile := 27 + (playerNumber-dealer)%playerNumber
				color.Yellow(util.MahjongZH[windTile])

				return nil
			} else {
				dealer = d.dealer
				if roundNumber > 0 && roundNumber != d.roundNumber {
					dealer = (dealer + 1) % playerNumber
				}
				d.reset(roundNumber, dealer)
			}
		default:
			panic("not impl!")
		}

		fmt.Printf("%s%d局开始，自风为%s\n", util.MahjongZH[d.roundWindTile], roundNumber%4+1, util.MahjongZH[d.players[0].selfWindTile])

		color.Yellow("宝牌指示牌是 %s", util.MahjongZH[doraIndicator])
		d.doraIndicators = []int{doraIndicator}
		d.descLeftCounts(doraIndicator)

		for _, tile := range hands {
			d.counts[tile]++
			d.descLeftCounts(tile)
		}

		if len(hands) == 14 {
			return analysisTiles34(d.counts, d.leftCounts, false)
		}
	case d.parser.IsOpen():
		// 某家鸣牌（含暗杠、加杠）
		who, meldType, meldTiles, calledTile, kanDoraIndicator := d.parser.ParseOpen()
		if kanDoraIndicator != -1 {
			d.newDora(kanDoraIndicator)
		}
		if meldType == meldTypeKakan {
			// TODO: 修改副露情况
			if who != 0 {
				d.descLeftCounts(calledTile)
			} else {
				// 自家加杠成功
				d.counts[calledTile]--
			}
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
			if meldType == meldTypeAnKan {
				d.counts[meldTiles[0]] = 0
			} else {
				d.counts[calledTile]++
				for _, tile := range meldTiles {
					d.counts[tile]--
				}
			}
			d.meldCount++
		}
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
		isOpen := len(d.players[0].melds) > 0
		return analysisTiles34(d.counts, d.leftCounts, isOpen)
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

		// 天凤：为防止先收到自家摸牌，然后收到上家摸牌，上家舍牌时不刷新
		if d.parser.GetDataSourceType() != dataSourceTypeTenhou || who != 3 {
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
			// 注意这里会标记到暗杠的舍牌上
			if len(player.meldDiscardsAt)+1 != len(player.melds) {
				fmt.Printf("玩家数据异常 %#v", *player)
			}
			player.meldDiscardsAt = append(player.meldDiscardsAt, len(player.discardTiles)-1)
			player.meldDiscardsAtGlobal = append(player.meldDiscardsAtGlobal, len(d.globalDiscardTiles)-1)
		}

		if d.parser.GetDataSourceType() != dataSourceTypeTenhou || who != 3 {
			// 打印他家舍牌信息
			d.printDiscards()
			fmt.Println()

			// 安全度分析
			riskTables := d.analysisTilesRisk()
			riskTables.printWithHands(d.counts, d.leftCounts)
		}

		// 若能副露，计算何切
		if canBeMeld {
			// TODO: 消除海底/避免河底/型听提醒
			allowChi := who == 3
			analysisMeld(d.counts, d.leftCounts, tile, allowChi)
		}
	case d.parser.IsRoundWin():
		if !debugMode {
			clearConsole()
		}
		fmt.Println("和牌，本局结束")
		whos, points := d.parser.ParseRoundWin()
		if len(whos) == 3 {
			color.Yellow("凤 凰 级 避 铳")
			if d.parser.GetDataSourceType() == dataSourceTypeMajsoul {
				color.Yellow("（快醒醒，这是雀魂）")
			}
		}
		for i, who := range whos {
			fmt.Println(d.players[who].name, points[i])
		}
	case d.parser.IsNewDora():
		// 杠宝牌
		// 1. 剩余牌减少
		// 2. 打点提高
		kanDoraIndicator := d.parser.ParseNewDora()
		d.newDora(kanDoraIndicator)
	default:
	}

	return nil
}
