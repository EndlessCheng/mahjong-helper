package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/EndlessCheng/mahjong-helper/util"
	"github.com/EndlessCheng/mahjong-helper/util/model"
)

var debugMode = false

const (
	dataSourceTypeTenhou = iota
	dataSourceTypeMajsoul
)

type DataParser interface {
	GetDataSourceType() int

	GetMessage() string

	// 解析前，根据消息内容来决定是否要进行后续解析
	CheckMessage() bool

	// 登录成功
	// 目前是在 server 逻辑上解析的
	//IsLogin() bool
	//HandleLogin()

	// round 开始/重连
	// roundNumber: 场数（如东1为0，东2为1，...，南1为4，...）
	// dealer: 庄家 0-3
	// doraIndicator: 宝牌指示牌
	// handTiles: 手牌
	// numRedFive: 赤5个数
	IsInit() bool
	ParseInit() (roundNumber int, dealer int, doraIndicator int, handTiles []int, numRedFive int)

	// 自家摸牌
	// tile: 0-33
	// isRedFive: 是否为赤5
	// kanDoraIndicator: 摸牌时，若为暗杠摸的岭上牌，则可以翻出杠宝牌指示牌，否则返回 -1（目前恒为 -1，见 IsNewDora）
	IsSelfDraw() bool
	ParseSelfDraw() (tile int, isRedFive bool, kanDoraIndicator int)

	// 舍牌
	// who: 0=自家, 1=下家, 2=对家, 3=上家
	// isTsumogiri: 是否为摸切（who=0 时忽略该值）
	// isReach: 是否为立直宣言（isReach 对于天凤来说恒为 false，见 IsReach）
	// canBeMeld: 是否可以鸣牌（who=0 时忽略该值）
	// kanDoraIndicator: 大明杠/加杠的杠宝牌指示牌，在切牌后出现，没有则返回 -1（天凤恒为-1，见 IsNewDora）
	IsDiscard() bool
	ParseDiscard() (who int, discardTile int, isRedFive bool, isTsumogiri bool, isReach bool, canBeMeld bool, kanDoraIndicator int)

	// 鸣牌（含暗杠、加杠）
	// kanDoraIndicator: 暗杠的杠宝牌指示牌，在他家暗杠时出现，没有则返回 -1（天凤恒为-1，见 IsNewDora）
	IsOpen() bool
	ParseOpen() (who int, meld *model.Meld, kanDoraIndicator int)

	// 立直声明（IsReach 对于雀魂来说恒为 false，见 ParseDiscard）
	IsReach() bool
	ParseReach() (who int)

	// 振听
	IsFuriten() bool

	// 本局是否和牌
	IsRoundWin() bool
	ParseRoundWin() (whos []int, points []int)

	// 是否流局
	//IsRyuukyoku() bool

	// 这一项放在末尾处理
	// 杠宝牌（雀魂在暗杠后的摸牌时出现）
	// kanDoraIndicator: 0-33
	IsNewDora() bool
	ParseNewDora() (kanDoraIndicator int)
}

//

// 负数变正数
func normalDiscardTiles(discardTiles []int) []int {
	newD := make([]int, len(discardTiles))
	copy(newD, discardTiles)
	for i, discardTile := range newD {
		if discardTile < 0 {
			newD[i] = ^discardTile
		}
	}
	return newD
}

//

const (
	meldTypeChi    = iota // 吃
	meldTypePon           // 碰
	meldTypeAnkan         // 暗杠
	meldTypeMinkan        // 大明杠
	meldTypeKakan         // 加杠
)

//

type playerInfo struct {
	name string // 自家/下家/对家/上家

	selfWindTile int // 自风

	melds                []*model.Meld // 副露
	meldDiscardsAtGlobal []int
	meldDiscardsAt       []int
	isNaki               bool // 是否鸣牌（暗杠不算鸣牌）

	// 注意负数（自摸切）要^
	discardTiles          []int // 该玩家的舍牌
	latestDiscardAtGlobal int   // 该玩家最近一次舍牌在 globalDiscardTiles 中的下标，初始为 -1
	earlyOutsideTiles     []int // 立直前的1-5巡的外侧牌

	isReached  bool // 是否立直
	canIppatsu bool // 是否有一发

	reachTileAtGlobal int // 立直宣言牌在 globalDiscardTiles 中的下标，初始为 -1
	reachTileAt       int // 立直宣言牌在 discardTiles 中的下标，初始为 -1
}

func newPlayerInfo(name string, selfWindTile int) *playerInfo {
	return &playerInfo{
		name:                  name,
		selfWindTile:          selfWindTile,
		latestDiscardAtGlobal: -1,
		reachTileAtGlobal:     -1,
		reachTileAt:           -1,
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
			if p.isNaki { // 副露
				fgColor = getOtherDiscardAlertColor(disTile) // 高亮中张手切
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
	counts []int

	// 自家赤5数量，包含副露的赤5
	numRedFive int

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
	return &roundData{
		parser:             parser,
		roundNumber:        roundNumber,
		roundWindTile:      roundWindTile,
		dealer:             dealer,
		counts:             make([]int, 34),
		leftCounts:         util.InitLeftTiles34(),
		globalDiscardTiles: []int{},
		players: []*playerInfo{
			newPlayerInfo("自家", playerWindTile[0]),
			newPlayerInfo("下家", playerWindTile[1]),
			newPlayerInfo("对家", playerWindTile[2]),
			newPlayerInfo("上家", playerWindTile[3]),
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

// 根据 dora 指示牌计算出 dora
func (d *roundData) doraList() (dl []int) {
	for _, doraIndicator := range d.doraIndicators {
		var dora int
		if doraIndicator < 27 {
			if doraIndicator%9 < 8 {
				dora = doraIndicator + 1
			} else {
				dora = doraIndicator - 8
			}
		} else if doraIndicator < 31 {
			if doraIndicator < 30 {
				dora = doraIndicator + 1
			} else {
				dora = 27
			}
		} else {
			if doraIndicator < 33 {
				dora = doraIndicator + 1
			} else {
				dora = 31
			}
		}
		dl = append(dl, dora)
	}
	return
}

func (d *roundData) printDiscards() {
	for i := len(d.players) - 1; i >= 1; i-- {
		d.players[i].printDiscards()
	}
}

// 分析34种牌的危险度
// 可以用来判断自家手牌的安全度，以及他家是否在进攻（多次切出危险度高的牌）
func (d *roundData) analysisTilesRisk() (riList riskInfoList) {
	riList = make(riskInfoList, len(d.players))
	for who := range riList {
		riList[who].safeTiles34 = make([]bool, 34)
	}

	// 先利用振听规则收集各家安牌
	for who, player := range d.players {
		if who == 0 {
			// TODO: 暂时不计算自家的
			continue
		}

		// 舍牌振听产生的安牌
		for _, tile := range normalDiscardTiles(player.discardTiles) {
			riList[who].safeTiles34[tile] = true
		}
		if player.reachTileAtGlobal != -1 {
			// 立直后振听产生的安牌
			for _, tile := range normalDiscardTiles(d.globalDiscardTiles[player.reachTileAtGlobal:]) {
				riList[who].safeTiles34[tile] = true
			}
		} else if player.latestDiscardAtGlobal != -1 {
			// 同巡振听产生的安牌
			// 即该玩家在最近一次舍牌后，其他玩家的舍牌
			for _, tile := range normalDiscardTiles(d.globalDiscardTiles[player.latestDiscardAtGlobal:]) {
				riList[who].safeTiles34[tile] = true
			}
		}
	}

	// 计算各种数据
	for who, player := range d.players {
		if who == 0 {
			// TODO: 暂时不计算自家的
			continue
		}

		// 该玩家的巡目 = 为其切过的牌的数目
		turns := util.MinInt(len(player.discardTiles), util.MaxTurns)
		if turns == 0 {
			turns = 1
		}

		// TODO: 若某人一直摸切，然后突然手切了一张字牌，那他很有可能默听/一向听
		if player.isReached {
			riList[who].tenpaiRate = 100.0
		} else {
			riList[who].tenpaiRate = util.CalcTenpaiRate(len(player.melds), player.discardTiles, player.meldDiscardsAt)
		}

		// 收集可能的安牌
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
		// 食延的情况 https://tieba.baidu.com/p/3688516724

		// 估计该玩家荣和点数
		var ronPoint float64
		switch {
		case player.canIppatsu:
			// 立直一发巡的荣和点数
			ronPoint = util.RonPointRiichiIppatsu
		case player.isReached:
			// 立直非一发巡的荣和点数
			ronPoint = util.RonPointRiichiHiIppatsu
		case player.isNaki:
			// 副露时的荣和点数（非常粗略地估计）
			doraCount := 0
			doraList := d.doraList()
			for _, meld := range player.melds {
				for _, tile := range meld.Tiles {
					for _, dora := range doraList {
						if tile == dora {
							doraCount++
						}
					}
				}
				if meld.ContainRedFive {
					doraCount++
				}
			}
			ronPoint = util.RonPointOtherNakiWithDora(doraCount)
		default:
			// 默听时的荣和点数
			ronPoint = util.RonPointDama
		}
		// 亲家*1.5
		if who == d.dealer {
			ronPoint *= 1.5
		}
		riList[who]._ronPoint = ronPoint

		// 根据该玩家的巡目、现物、立直后通过的牌、NC、Dora、早外、荣和点数来计算每张牌的危险度
		risk34 := util.CalculateRiskTiles34(turns, riList[who].safeTiles34, d.leftCounts, d.doraList(), d.roundWindTile, player.selfWindTile).
			FixWithEarlyOutside(player.earlyOutsideTiles).
			FixWithPoint(ronPoint)
		riList[who].riskTable = riskTable(risk34)

		// 计算剩余筋牌
		riList[who].leftNoSujiTiles = util.CalculateLeftNoSujiTiles(riList[who].safeTiles34, d.leftCounts)
	}

	return riList
}

func (d *roundData) isPlayerDaburii(who int) bool {
	// w立直成立的前提是没有任何玩家副露
	for _, p := range d.players {
		if len(p.melds) > 0 {
			return false
		}
	}
	return d.players[who].reachTileAt == 0
}

func (d *roundData) newModelPlayerInfo() *model.PlayerInfo {
	melds := []model.Meld{}
	for _, m := range d.players[0].melds {
		melds = append(melds, *m)
	}

	doraCount := 0
	doraList := d.doraList()
	for _, dora := range doraList {
		doraCount += d.counts[dora]
		for _, m := range melds {
			for _, tile := range m.Tiles {
				if tile == dora {
					doraCount++
				}
			}
		}
	}
	// 手牌和副露中的赤5
	doraCount += d.numRedFive

	const self = 0
	selfPlayer := d.players[self]

	return &model.PlayerInfo{
		HandTiles34:   d.counts,
		Melds:         melds,
		RoundWindTile: d.roundWindTile,
		SelfWindTile:  selfPlayer.selfWindTile,
		DoraCount:     doraCount,
		IsParent:      d.dealer == self,
		IsDaburii:     d.isPlayerDaburii(self),
		IsRiichi:      selfPlayer.isReached,
		DiscardTiles:  normalDiscardTiles(selfPlayer.discardTiles),
		LeftTiles34:   d.leftCounts,
	}
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

		roundNumber, dealer, doraIndicator, hands, numRedFive := d.parser.ParseInit()
		switch d.parser.GetDataSourceType() {
		case dataSourceTypeTenhou:
			d.reset(roundNumber, dealer)
		case dataSourceTypeMajsoul:
			playerNumber := len(d.players)
			if dealer != -1 {
				d.dealer = dealer

				fmt.Printf("游戏即将开始，您分配到的座位是：")
				windTile := 27 + (playerNumber-dealer)%playerNumber
				color.HiGreen(util.MahjongZH[windTile])

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

		color.HiYellow("宝牌指示牌是 %s", util.MahjongZH[doraIndicator])
		d.doraIndicators = []int{doraIndicator}
		d.descLeftCounts(doraIndicator)

		for _, tile := range hands {
			d.counts[tile]++
			d.descLeftCounts(tile)
		}

		d.numRedFive = numRedFive

		if len(hands) == 14 {
			return analysisTiles34(d.newModelPlayerInfo(), nil)
		}
	case d.parser.IsOpen():
		// 某家鸣牌（含暗杠、加杠）
		who, meld, kanDoraIndicator := d.parser.ParseOpen()
		meldType := meld.MeldType
		meldTiles := meld.Tiles
		calledTile := meld.CalledTile

		// 任何形式的鸣牌都能破除一发
		for _, player := range d.players {
			player.canIppatsu = false
		}

		// 杠宝牌指示牌
		if kanDoraIndicator != -1 {
			d.newDora(kanDoraIndicator)
		}

		player := d.players[who]

		// 不是暗杠则标记该玩家鸣牌了
		if meldType != meldTypeAnkan {
			player.isNaki = true
		}

		// 加杠特殊处理
		if meldType == meldTypeKakan {
			if who != 0 {
				// （不是自家时）修改牌山剩余量
				d.descLeftCounts(calledTile)
			} else {
				// 自家加杠成功，修改手牌
				d.counts[calledTile]--
			}
			// 修改原副露
			for _, _meld := range player.melds {
				// 找到原有的碰副露
				if _meld.Tiles[0] == calledTile {
					_meld.MeldType = meldTypeKakan
					_meld.Tiles = append(_meld.Tiles, calledTile)
					break
				}
			}
			break
		}

		// 修改玩家副露数据
		d.players[who].melds = append(d.players[who].melds, meld)

		if who != 0 {
			// （不是自家时）修改牌山剩余量
			// 先增后减
			if meldType != meldTypeAnkan {
				d.leftCounts[calledTile]++
			}
			for _, tile := range meldTiles {
				d.descLeftCounts(tile)
			}
		} else {
			// 自家，修改手牌
			if meldType == meldTypeAnkan {
				d.counts[meldTiles[0]] = 0
			} else {
				d.counts[calledTile]++
				for _, tile := range meldTiles {
					d.counts[tile]--
				}
				if meld.RedFiveFromOthers {
					d.numRedFive++
				}
			}
		}
	case d.parser.IsReach():
		// 立直宣告
		// 如果是他家立直，进入攻守判断模式
		who := d.parser.ParseReach()
		d.players[who].isReached = true
		d.players[who].canIppatsu = true
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
		color.HiYellow("振听")
		//case "U", "V", "W":
		//	//（下家,对家,上家 不要其上家的牌）摸牌
		//case "HELO", "RANKING", "TAIKYOKU", "UN", "LN", "SAIKAI":
		//	// 其他
	case d.parser.IsSelfDraw():
		if !debugMode {
			clearConsole()
		}
		// 自家（从牌山 d.leftCounts）摸牌（至手牌 d.counts）
		tile, isRedFive, kanDoraIndicator := d.parser.ParseSelfDraw()
		d.descLeftCounts(tile)
		d.counts[tile]++
		if isRedFive {
			d.numRedFive++
		}
		if kanDoraIndicator != -1 {
			d.newDora(kanDoraIndicator)
		}

		// 打印他家舍牌信息
		d.printDiscards()
		fmt.Println()

		// 安全度分析
		riskTables := d.analysisTilesRisk()
		riskTables.printWithHands(d.counts, d.leftCounts)

		mixedRiskTable := riskTables.mixedRiskTable()

		// 何切
		// TODO: 根据是否听牌/一向听、打点、巡目、和率等进行攻守判断
		return analysisTiles34(d.newModelPlayerInfo(), mixedRiskTable)
	case d.parser.IsDiscard():
		who, discardTile, isRedFive, isTsumogiri, isReach, canBeMeld, kanDoraIndicator := d.parser.ParseDiscard()

		if kanDoraIndicator != -1 {
			d.newDora(kanDoraIndicator)
		}

		player := d.players[who]
		if isReach {
			player.isReached = true
			player.canIppatsu = true
		}

		if who == 0 {
			// 自家（从手牌 d.counts）舍牌（至牌河 d.globalDiscardTiles）
			d.counts[discardTile]--

			d.globalDiscardTiles = append(d.globalDiscardTiles, discardTile)
			player.discardTiles = append(player.discardTiles, discardTile)
			player.latestDiscardAtGlobal = len(d.globalDiscardTiles) - 1

			if isRedFive {
				d.numRedFive--
			}

			return nil
		}

		// 他家舍牌
		d.descLeftCounts(discardTile)

		// 天凤fix：为防止先收到自家摸牌，然后收到上家摸牌，上家舍牌时不刷新
		if d.parser.GetDataSourceType() != dataSourceTypeTenhou || who != 3 {
			if !debugMode {
				clearConsole()
			}
		}

		_disTile := discardTile
		if isTsumogiri {
			_disTile = ^_disTile
		}
		d.globalDiscardTiles = append(d.globalDiscardTiles, _disTile)
		player.discardTiles = append(player.discardTiles, _disTile)
		player.latestDiscardAtGlobal = len(d.globalDiscardTiles) - 1

		// 标记外侧牌
		if !player.isReached && len(player.discardTiles) <= 5 {
			player.earlyOutsideTiles = append(player.earlyOutsideTiles, util.OutsideTiles(discardTile)...)
		}

		if player.isReached && player.reachTileAtGlobal == -1 {
			// 标记立直宣言牌
			player.reachTileAtGlobal = len(d.globalDiscardTiles) - 1
			player.reachTileAt = len(player.discardTiles) - 1

			// 若该玩家摸切立直，打印提示信息
			if isTsumogiri {
				color.HiYellow("%s 摸切立直！", d.players[who].name)
			}
		} else if len(player.meldDiscardsAt) != len(player.melds) {
			// 标记鸣牌的舍牌
			// 注意这里会标记到暗杠后的舍牌上
			if len(player.meldDiscardsAt)+1 != len(player.melds) {
				fmt.Printf("玩家数据异常 %#v", *player)
			}
			player.meldDiscardsAt = append(player.meldDiscardsAt, len(player.discardTiles)-1)
			player.meldDiscardsAtGlobal = append(player.meldDiscardsAtGlobal, len(d.globalDiscardTiles)-1)
		}

		// 若玩家在立直后摸牌舍牌，则没有一发
		if player.reachTileAt < len(player.discardTiles)-1 {
			player.canIppatsu = false
		}

		// 安全度分析
		riskTables := d.analysisTilesRisk()

		if d.parser.GetDataSourceType() != dataSourceTypeTenhou || who != 3 {
			// 打印他家舍牌信息
			d.printDiscards()
			fmt.Println()
			riskTables.printWithHands(d.counts, d.leftCounts)
		}

		// 若能副露，计算何切
		if canBeMeld {
			// TODO: 消除海底/避免河底/型听提醒
			allowChi := who == 3
			mixedRiskTable := riskTables.mixedRiskTable()
			analysisMeld(d.newModelPlayerInfo(), discardTile, allowChi, mixedRiskTable)
		}
	case d.parser.IsRoundWin():
		if !debugMode {
			clearConsole()
		}
		fmt.Println("和牌，本局结束")
		whos, points := d.parser.ParseRoundWin()
		if len(whos) == 3 {
			color.HiYellow("凤 凰 级 避 铳")
			if d.parser.GetDataSourceType() == dataSourceTypeMajsoul {
				color.HiYellow("（快醒醒，这是雀魂）")
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
