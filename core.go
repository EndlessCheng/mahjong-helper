package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/EndlessCheng/mahjong-helper/util"
	"github.com/EndlessCheng/mahjong-helper/util/model"
)

type DataParser interface {
	// 数据来源是天凤还是雀魂
	GetDataSourceType() int

	// 获取自家初始座位：0-第一局的东家 1-第一局的南家 2-第一局的西家 3-第一局的北家
	// 仅处理雀魂数据，天凤返回 -1
	GetSelfSeat() int

	// 原始 JSON
	GetMessage() string

	// 解析前，根据消息内容来决定是否要进行后续解析
	CheckMessage() bool

	// 登录成功
	// TODO: 重构，目前是在 server 逻辑上解析的
	//IsLogin() bool
	//HandleLogin()

	// round 开始/重连
	// roundNumber: 场数（如东1为0，东2为1，...，南1为4，...，南4为7，...）
	// benNumber: 本场数
	// dealer: 庄家 0-3
	// doraIndicator: 宝牌指示牌
	// handTiles: 手牌
	// numRedFives: 按照 mps 的顺序，赤5个数
	IsInit() bool
	ParseInit() (roundNumber int, benNumber int, dealer int, doraIndicator int, handTiles []int, numRedFives []int)

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
	// 四风连打 四家立直 四杠散了 九种九牌 三家和了 | 流局听牌 流局未听牌 | 流局满贯
	// 三家和了
	// "{\"tag\":\"RYUUKYOKU\",\"type\":\"ron3\",\"ba\":\"1,1\",\"sc\":\"290,0,228,0,216,0,256,0\",\"hai0\":\"18,19,30,32,33,41,43,94,95,114,115,117,119\",\"hai2\":\"29,31,74,75\",\"hai3\":\"8,13,17,25,35,46,48,53,78,79\"}"
	//IsRyuukyoku() bool
	//ParseRyuukyoku() (type_ int, whos []int, points []int)

	// 这一项放在末尾处理
	// 杠宝牌（雀魂在暗杠后的摸牌时出现）
	// kanDoraIndicator: 0-33
	IsNewDora() bool
	ParseNewDora() (kanDoraIndicator int)
}

type playerInfo struct {
	name string // 自家/下家/对家/上家

	//turn int // 该玩家的巡目（从13张牌的状态开始算，每「得到」一张牌，巡目就+1。比如：亲家一开始就是第一巡、副露后巡目加一 ）

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

type roundData struct {
	parser DataParser

	gameMode gameMode

	skipOutput bool

	// 场数（如东1为0，东2为1，...，南1为4，...）
	roundNumber int

	// 本场数，从 0 开始算
	benNumber int

	// 场风
	roundWindTile int

	// 庄家 0=自家, 1=下家, 2=对家, 3=上家
	// 请用 reset 设置
	dealer int

	// 宝牌指示牌
	doraIndicators []int

	// 自家手牌
	counts []int

	// 按照 mps 的顺序记录自家赤5数量，包含副露的赤5
	// 比如有 0p 和 0s 就是 [1, 0, 1]
	numRedFives []int

	// 牌山剩余牌量
	leftCounts []int

	// 全局舍牌
	// 按舍牌顺序，负数表示摸切(-)，非负数表示手切(+)
	// 可以理解成：- 表示不要/暗色，+ 表示进张/亮色
	globalDiscardTiles []int

	// 0=自家, 1=下家, 2=对家, 3=上家
	players []*playerInfo
}

func newRoundData(parser DataParser, roundNumber int, benNumber int, dealer int) *roundData {
	const playerNumber = 4
	roundWindTile := 27 + roundNumber/playerNumber
	playerWindTile := make([]int, playerNumber)
	for i := 0; i < playerNumber; i++ {
		playerWindTile[i] = 27 + (playerNumber-dealer+i)%playerNumber
	}
	return &roundData{
		parser:      parser,
		roundNumber: roundNumber,
		benNumber:   benNumber,

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

func newGame(parser DataParser) *roundData {
	return newRoundData(parser, 0, 0, 0)
}

// 新的一局
func (d *roundData) reset(roundNumber int, benNumber int, dealer int) {
	skipOutput := d.skipOutput
	gameMode := d.gameMode
	newData := newRoundData(d.parser, roundNumber, benNumber, dealer)
	newData.skipOutput = skipOutput
	newData.gameMode = gameMode
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
	d.doraIndicators = append(d.doraIndicators, kanDoraIndicator)
	d.descLeftCounts(kanDoraIndicator)

	if d.skipOutput {
		return
	}

	color.Yellow("杠宝牌指示牌是 %s", util.MahjongZH[kanDoraIndicator])
}

// 根据宝牌指示牌计算出宝牌
func (d *roundData) doraList() (dl []int) {
	return model.DoraList(d.doraIndicators)
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

	const self = 0
	selfPlayer := d.players[self]

	return &model.PlayerInfo{
		HandTiles34: d.counts,
		Melds:       melds,
		DoraTiles:   d.doraList(),
		NumRedFives: d.numRedFives,

		RoundWindTile: d.roundWindTile,
		SelfWindTile:  selfPlayer.selfWindTile,
		IsParent:      d.dealer == self,
		//IsDaburii:     d.isPlayerDaburii(self), // FIXME PLS，应该在立直时就判断
		IsRiichi: selfPlayer.isReached,

		DiscardTiles: normalDiscardTiles(selfPlayer.discardTiles),
		LeftTiles34:  d.leftCounts,
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
		if msg := d.parser.GetMessage(); len(msg) > 0 {
			const printLimit = 500
			if len(msg) > printLimit {
				msg = msg[:printLimit]
			}
			fmt.Println("收到", msg)
		}
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
		if !debugMode && !d.skipOutput {
			clearConsole()
		}

		roundNumber, benNumber, dealer, doraIndicator, hands, numRedFives := d.parser.ParseInit()
		switch d.parser.GetDataSourceType() {
		case dataSourceTypeTenhou:
			d.reset(roundNumber, 0, dealer)
		case dataSourceTypeMajsoul:
			playerNumber := len(d.players)
			if dealer != -1 { // 先就坐，还没洗牌呢~
				// 设置第一局的 dealer
				d.reset(0, 0, dealer)
				d.gameMode = gameModeMatch

				fmt.Printf("游戏即将开始，您分配到的座位是：")
				windTile := 27 + (playerNumber-dealer)%playerNumber
				color.HiGreen(util.MahjongZH[windTile])

				return nil
			} else {
				// 根据当前的 roundNumber 和 selfSeat 计算当前局的 dealer
				newDealer := (len(d.players) - d.parser.GetSelfSeat() + roundNumber) % len(d.players)
				// 新的一局
				d.reset(roundNumber, benNumber, newDealer)
			}
		default:
			panic("not impl!")
		}

		d.doraIndicators = []int{doraIndicator}
		d.descLeftCounts(doraIndicator)
		for _, tile := range hands {
			d.counts[tile]++
			d.descLeftCounts(tile)
		}
		d.numRedFives = numRedFives

		playerInfo := d.newModelPlayerInfo()

		// 牌谱分析模式下，记录舍牌推荐
		if d.gameMode == gameModeRecordCache && len(hands) == 14 {
			currentRoundCache := globalAnalysisCache.wholeGameCache[d.roundNumber][d.benNumber]
			currentRoundCache.addAIDiscardTileWhenDrawTile(simpleBestDiscardTile(playerInfo), -1, 0, 0)
		}

		if d.skipOutput {
			return nil
		}

		// 牌谱模式下，打印舍牌推荐
		if d.gameMode == gameModeRecord {
			currentRoundCache := globalAnalysisCache.wholeGameCache[d.roundNumber][d.benNumber]
			currentRoundCache.print()
		}

		color.New(color.FgHiGreen).Printf("%s", util.MahjongZH[d.roundWindTile])
		fmt.Printf("%d局开始，自风为", roundNumber%4+1)
		color.New(color.FgHiGreen).Printf("%s", util.MahjongZH[d.players[0].selfWindTile])
		fmt.Println()
		color.HiYellow("宝牌指示牌是 %s", util.MahjongZH[doraIndicator])
		fmt.Println()
		// TODO: 显示地和概率
		return analysisPlayerWithRisk(playerInfo, nil)
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

		// 加杠单独处理
		if meldType == meldTypeKakan {
			if who != 0 {
				// （不是自家时）修改牌山剩余量
				d.descLeftCounts(calledTile)
			} else {
				// 自家加杠成功，修改手牌
				d.counts[calledTile]--
				// 由于均为自家操作，宝牌数是不变的

				// 牌谱分析模式下，记录加杠操作
				if d.gameMode == gameModeRecordCache {
					currentRoundCache := globalAnalysisCache.wholeGameCache[d.roundNumber][d.benNumber]
					currentRoundCache.addKan(meldType)
				}
			}
			// 修改原副露
			for _, _meld := range player.melds {
				// 找到原有的碰副露
				if _meld.Tiles[0] == calledTile {
					_meld.MeldType = meldTypeKakan
					_meld.Tiles = append(_meld.Tiles, calledTile)
					_meld.ContainRedFive = meld.ContainRedFive
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

				// 牌谱分析模式下，记录暗杠操作
				if d.gameMode == gameModeRecordCache {
					currentRoundCache := globalAnalysisCache.wholeGameCache[d.roundNumber][d.benNumber]
					currentRoundCache.addKan(meldType)
				}
			} else {
				d.counts[calledTile]++
				for _, tile := range meldTiles {
					d.counts[tile]--
				}
				if meld.RedFiveFromOthers {
					tileType := meldTiles[0] / 9
					d.numRedFives[tileType]++
				}

				// 牌谱分析模式下，记录吃碰明杠操作
				if d.gameMode == gameModeRecordCache {
					currentRoundCache := globalAnalysisCache.wholeGameCache[d.roundNumber][d.benNumber]
					currentRoundCache.addChiPonKan(meldType)
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
		if d.skipOutput {
			return nil
		}
		color.HiYellow("振听")
		//case "U", "V", "W":
		//	//（下家,对家,上家 不要其上家的牌）摸牌
		//case "HELO", "RANKING", "TAIKYOKU", "UN", "LN", "SAIKAI":
		//	// 其他
	case d.parser.IsSelfDraw():
		if !debugMode && !d.skipOutput {
			clearConsole()
		}
		// 自家（从牌山 d.leftCounts）摸牌（至手牌 d.counts）
		tile, isRedFive, kanDoraIndicator := d.parser.ParseSelfDraw()
		d.descLeftCounts(tile)
		d.counts[tile]++
		if isRedFive {
			d.numRedFives[tile/9]++
		}
		if kanDoraIndicator != -1 {
			d.newDora(kanDoraIndicator)
		}

		playerInfo := d.newModelPlayerInfo()

		// 安全度分析
		riskTables := d.analysisTilesRisk()
		mixedRiskTable := riskTables.mixedRiskTable()

		// 牌谱分析模式下，记录舍牌推荐
		if d.gameMode == gameModeRecordCache {
			currentRoundCache := globalAnalysisCache.wholeGameCache[d.roundNumber][d.benNumber]
			bestAttackDiscardTile := simpleBestDiscardTile(playerInfo)
			bestDefenceDiscardTile := mixedRiskTable.getBestDefenceTile()
			bestAttackDiscardTileRisk, bestDefenceDiscardTileRisk := 0.0, 0.0
			if bestDefenceDiscardTile >= 0 {
				bestAttackDiscardTileRisk = mixedRiskTable[bestAttackDiscardTile]
				bestDefenceDiscardTileRisk = mixedRiskTable[bestDefenceDiscardTile]
			}
			currentRoundCache.addAIDiscardTileWhenDrawTile(bestAttackDiscardTile, bestDefenceDiscardTile, bestAttackDiscardTileRisk, bestDefenceDiscardTileRisk)
		}

		if d.skipOutput {
			return nil
		}

		// 牌谱模式下，打印舍牌推荐
		if d.gameMode == gameModeRecord {
			currentRoundCache := globalAnalysisCache.wholeGameCache[d.roundNumber][d.benNumber]
			currentRoundCache.print()
		}

		// 打印他家舍牌信息
		d.printDiscards()
		fmt.Println()

		// 打印手牌对各家的安全度
		riskTables.printWithHands(d.counts, d.leftCounts)

		// 打印何切推荐
		// TODO: 根据是否听牌/一向听、打点、巡目、和率等进行攻守判断
		return analysisPlayerWithRisk(playerInfo, mixedRiskTable)
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
			// 特殊处理自家舍牌的情况
			riskTables := d.analysisTilesRisk()
			mixedRiskTable := riskTables.mixedRiskTable()

			// 自家（从手牌 d.counts）舍牌（至牌河 d.globalDiscardTiles）
			d.counts[discardTile]--

			d.globalDiscardTiles = append(d.globalDiscardTiles, discardTile)
			player.discardTiles = append(player.discardTiles, discardTile)
			player.latestDiscardAtGlobal = len(d.globalDiscardTiles) - 1

			if isRedFive {
				d.numRedFives[discardTile/9]--
			}

			// 牌谱分析模式下，记录自家舍牌
			if d.gameMode == gameModeRecordCache {
				currentRoundCache := globalAnalysisCache.wholeGameCache[d.roundNumber][d.benNumber]
				currentRoundCache.addSelfDiscardTile(discardTile, mixedRiskTable[discardTile], isReach)
			}

			return nil
		}

		// 他家舍牌
		d.descLeftCounts(discardTile)

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
			if isTsumogiri && !d.skipOutput {
				color.HiYellow("%s 摸切立直！", player.name)
			}
		} else if len(player.meldDiscardsAt) != len(player.melds) {
			// 标记鸣牌的舍牌
			// 注意这里会标记到暗杠后的舍牌上
			// 注意对于连续开杠的情况，len(player.meldDiscardsAt) 和 len(player.melds) 是不等的
			player.meldDiscardsAt = append(player.meldDiscardsAt, len(player.discardTiles)-1)
			player.meldDiscardsAtGlobal = append(player.meldDiscardsAtGlobal, len(d.globalDiscardTiles)-1)
		}

		// 若玩家在立直后摸牌舍牌，则没有一发
		if player.reachTileAt < len(player.discardTiles)-1 {
			player.canIppatsu = false
		}

		playerInfo := d.newModelPlayerInfo()

		// 安全度分析
		riskTables := d.analysisTilesRisk()
		mixedRiskTable := riskTables.mixedRiskTable()

		// 牌谱分析模式下，记录可能的鸣牌
		if d.gameMode == gameModeRecordCache {
			currentRoundCache := globalAnalysisCache.wholeGameCache[d.roundNumber][d.benNumber]
			allowChi := who == 3
			_, results14, incShantenResults14 := util.CalculateMeld(playerInfo, discardTile, isRedFive, allowChi)
			bestAttackDiscardTile := -1
			if len(results14) > 0 {
				bestAttackDiscardTile = results14[0].DiscardTile
			} else if len(incShantenResults14) > 0 {
				bestAttackDiscardTile = incShantenResults14[0].DiscardTile
			}
			if bestAttackDiscardTile != -1 {
				bestDefenceDiscardTile := mixedRiskTable.getBestDefenceTile()
				bestAttackDiscardTileRisk := 0.0
				if bestDefenceDiscardTile >= 0 {
					bestAttackDiscardTileRisk = mixedRiskTable[bestAttackDiscardTile]
				}
				currentRoundCache.addPossibleChiPonKan(bestAttackDiscardTile, bestAttackDiscardTileRisk)
			}
		}

		if d.skipOutput {
			return nil
		}

		// 上家舍牌时若无法鸣牌则跳过显示
		if who == 3 && !canBeMeld {
			return nil
		}

		if !debugMode {
			clearConsole()
		}

		// 牌谱模式下，打印舍牌推荐
		if d.gameMode == gameModeRecord {
			currentRoundCache := globalAnalysisCache.wholeGameCache[d.roundNumber][d.benNumber]
			currentRoundCache.print()
		}

		// 打印他家舍牌信息
		d.printDiscards()
		fmt.Println()
		riskTables.printWithHands(d.counts, d.leftCounts)

		// 天凤人机对战时，偶尔会有先收到他家舍牌消息然后才收到自家舍牌消息的情况
		// 这时 analysisMeld 会因手牌数量异常而失败
		// TODO: 可以考虑在绘制动画时才发送消息给客户端？
		if d.parser.GetDataSourceType() == dataSourceTypeTenhou && !canBeMeld {
			return nil
		}

		// 为了方便解析牌谱，这里尽可能地解析副露
		// TODO: 提醒: 消除海底/避免河底/型听
		// FIXME: 最后一张牌是无法鸣牌的
		allowChi := who == 3
		return analysisMeld(playerInfo, discardTile, isRedFive, allowChi, mixedRiskTable)
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
