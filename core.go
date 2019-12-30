package main

import (
	"fmt"
	"github.com/EndlessCheng/mahjong-helper/util"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"github.com/fatih/color"
)

type DataParser interface {
	// 数据来源（是天凤还是雀魂）
	GetDataSourceType() int

	// 获取自家初始座位：0-第一局的东家 1-第一局的南家 2-第一局的西家 3-第一局的北家
	// 仅处理雀魂数据，天凤返回 -1
	GetSelfSeat() int

	// 原始 JSON
	GetMessage() string

	// 解析前，根据消息内容来决定是否要进行后续解析
	SkipMessage() bool

	// 尝试解析用户名
	IsLogin() bool
	HandleLogin()

	// round 开始/重连
	// roundNumber: 场数（如东1为0，东2为1，...，南1为4，...，南4为7，...），对于三麻来说南1也是4
	// benNumber: 本场数
	// dealer: 庄家 0-3
	// doraIndicators: 宝牌指示牌
	// handTiles: 手牌
	// numRedFives: 按照 mps 的顺序，赤5个数
	IsInit() bool
	ParseInit() (roundNumber int, benNumber int, dealer int, doraIndicators []int, handTiles []int, numRedFives []int)

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
	IsRyuukyoku() bool
	ParseRyuukyoku() (type_ int, whos []int, points []int)

	// 拔北宝牌
	IsNukiDora() bool
	ParseNukiDora() (who int, isTsumogiri bool)

	// 这一项放在末尾处理
	// 杠宝牌（雀魂在暗杠后的摸牌时出现）
	// kanDoraIndicator: 0-33
	IsNewDora() bool
	ParseNewDora() (kanDoraIndicator int)
}

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

	nukiDoraNum int // 拔北宝牌数
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

func modifySanninPlayerInfoList(lst []*playerInfo, roundNumber int) []*playerInfo {
	windToIdxMap := map[int]int{}
	for i, pi := range lst {
		windToIdxMap[pi.selfWindTile] = i
	}

	idxS, idxW, idxN := windToIdxMap[28], windToIdxMap[29], windToIdxMap[30]
	switch roundNumber % 4 {
	case 0:
	case 1:
		// 北和西交换
		lst[idxN].selfWindTile, lst[idxW].selfWindTile = lst[idxW].selfWindTile, lst[idxN].selfWindTile
	case 2:
		// 北和西交换，再和南交换
		lst[idxN].selfWindTile, lst[idxW].selfWindTile, lst[idxS].selfWindTile = lst[idxW].selfWindTile, lst[idxS].selfWindTile, lst[idxN].selfWindTile
	default:
		panic("[modifySanninPlayerInfoList] 代码有误")
	}
	return lst
}

func (p *playerInfo) doraNum(doraList []int) (doraCount int) {
	for _, meld := range p.melds {
		for _, tile := range meld.Tiles {
			for _, doraTile := range doraList {
				if tile == doraTile {
					doraCount++
				}
			}
		}
		if meld.ContainRedFive {
			doraCount++
		}
	}
	if p.nukiDoraNum > 0 {
		doraCount += p.nukiDoraNum
		// 特殊：西为指示牌
		for _, doraTile := range doraList {
			if doraTile == 30 {
				doraCount += p.nukiDoraNum
			}
		}
	}
	return
}

//

type roundData struct {
	parser DataParser

	gameMode gameMode

	skipOutput bool

	// 玩家数，3 为三麻，4 为四麻
	playerNumber int

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
	// 无论是三麻还是四麻，都视作四个人
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
	playerNumber := d.playerNumber
	newData := newRoundData(d.parser, roundNumber, benNumber, dealer)
	newData.skipOutput = skipOutput
	newData.gameMode = gameMode
	newData.playerNumber = playerNumber
	if playerNumber == 3 {
		// 三麻没有 2-8m
		for i := 1; i <= 7; i++ {
			newData.leftCounts[i] = 0
		}
		newData.players = modifySanninPlayerInfoList(newData.players, roundNumber)
	}
	*d = *newData
}

func (d *roundData) newGame() {
	d.reset(0, 0, 0)
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

// 杠！
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
	return model.DoraList(d.doraIndicators, d.playerNumber == 3)
}

func (d *roundData) printDiscards() {
	// 三麻的北家是不需要打印的
	for i := len(d.players) - 1; i >= 1; i-- {
		if player := d.players[i]; d.playerNumber != 3 || player.selfWindTile != 30 {
			player.printDiscards()
		}
	}
}

// 分析34种牌的危险度
// 可以用来判断自家手牌的安全度，以及他家是否在进攻（多次切出危险度高的牌）
func (d *roundData) analysisTilesRisk() (riList riskInfoList) {
	riList = make(riskInfoList, len(d.players))
	for who := range riList {
		riList[who] = &riskInfo{
			playerNumber: d.playerNumber,
			safeTiles34:  make([]bool, 34),
		}
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

		// 特殊：杠产生的安牌
		// 很难想象一个人会在有 678888 的时候去开杠（即使有这个可能，本程序也是不防的）
		for _, meld := range player.melds {
			if meld.IsKan() {
				riList[who].safeTiles34[meld.Tiles[0]] = true
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
			if player.reachTileAtGlobal < len(d.globalDiscardTiles) { // 天凤可能有数据漏掉
				riList[who].isTsumogiriRiichi = d.globalDiscardTiles[player.reachTileAtGlobal] < 0
			}
		} else {
			rate := util.CalcTenpaiRate(player.melds, player.discardTiles, player.meldDiscardsAt)
			if d.playerNumber == 3 {
				rate = util.GetTenpaiRate3(rate)
			}
			riList[who].tenpaiRate = rate
		}

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
			doraCount := player.doraNum(d.doraList())
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
		if len(player.melds) < 4 {
			riList[who].leftNoSujiTiles = util.CalculateLeftNoSujiTiles(riList[who].safeTiles34, d.leftCounts)
		} else {
			// 大吊车：愚型听牌
		}
	}

	return riList
}

// TODO: 特殊处理w立直
func (d *roundData) isPlayerDaburii(who int) bool {
	// w立直成立的前提是没有任何玩家副露
	for _, p := range d.players {
		if len(p.melds) > 0 {
			return false
		}
		// 对于三麻来说，还不能有拔北
		if p.nukiDoraNum > 0 {
			return false
		}
	}
	return d.players[who].reachTileAt == 0
}

// 自家的 PlayerInfo
func (d *roundData) newModelPlayerInfo() *model.PlayerInfo {
	const wannpaiTilesCount = 14
	leftDrawTilesCount := util.CountOfTiles34(d.leftCounts) - (wannpaiTilesCount - len(d.doraIndicators))
	for _, player := range d.players[1:] {
		leftDrawTilesCount -= 13 - 3*len(player.melds)
	}
	if d.playerNumber == 3 {
		leftDrawTilesCount += 13
	}

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

		LeftDrawTilesCount: leftDrawTilesCount,

		NukiDoraNum: selfPlayer.nukiDoraNum,
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

	// 先获取用户信息
	if d.parser.IsLogin() {
		d.parser.HandleLogin()
	}

	if d.parser.SkipMessage() {
		return nil
	}

	// 若自家立直，则进入看戏模式
	// TODO: 见逃判断
	if !d.parser.IsInit() && !d.parser.IsRoundWin() && !d.parser.IsRyuukyoku() && d.players[0].isReached {
		return nil
	}

	if debugMode {
		fmt.Println("当前座位为", d.parser.GetSelfSeat())
	}

	var currentRoundCache *roundAnalysisCache
	if analysisCache := getAnalysisCache(d.parser.GetSelfSeat()); analysisCache != nil {
		currentRoundCache = analysisCache.wholeGameCache[d.roundNumber][d.benNumber]
	}

	switch {
	case d.parser.IsInit():
		// round 开始/重连
		if !debugMode && !d.skipOutput {
			clearConsole()
		}

		roundNumber, benNumber, dealer, doraIndicators, hands, numRedFives := d.parser.ParseInit()
		switch d.parser.GetDataSourceType() {
		case dataSourceTypeTenhou:
			d.reset(roundNumber, benNumber, dealer)
			d.gameMode = gameModeMatch // TODO: 牌谱模式？
		case dataSourceTypeMajsoul:
			if dealer != -1 { // 先就坐，还没洗牌呢~
				// 设置第一局的 dealer
				d.reset(0, 0, dealer)
				d.gameMode = gameModeMatch
				fmt.Printf("游戏即将开始，您分配到的座位是：")
				color.HiGreen(util.MahjongZH[d.players[0].selfWindTile])
				return nil
			} else {
				// 根据 selfSeat 和当前的 roundNumber 计算当前局的 dealer
				newDealer := (4 - d.parser.GetSelfSeat() + roundNumber) % 4
				// 新的一局
				d.reset(roundNumber, benNumber, newDealer)
			}
		default:
			panic("not impl!")
		}

		// 由于 reset 了，重新获取 currentRoundCache
		if analysisCache := getAnalysisCache(d.parser.GetSelfSeat()); analysisCache != nil {
			currentRoundCache = analysisCache.wholeGameCache[d.roundNumber][d.benNumber]
		}

		d.doraIndicators = doraIndicators
		for _, dora := range doraIndicators {
			d.descLeftCounts(dora)
		}
		for _, tile := range hands {
			d.counts[tile]++
			d.descLeftCounts(tile)
		}
		d.numRedFives = numRedFives

		playerInfo := d.newModelPlayerInfo()

		// 牌谱分析模式下，记录舍牌推荐
		if d.gameMode == gameModeRecordCache && len(hands) == 14 {
			currentRoundCache.addAIDiscardTileWhenDrawTile(simpleBestDiscardTile(playerInfo), -1, 0, 0)
		}

		if d.skipOutput {
			return nil
		}

		// 牌谱模式下，打印舍牌推荐
		if d.gameMode == gameModeRecord {
			currentRoundCache.print()
		}

		color.New(color.FgHiGreen).Printf("%s", util.MahjongZH[d.roundWindTile])
		fmt.Printf("%d局开始，自风为", roundNumber%4+1)
		color.New(color.FgHiGreen).Printf("%s", util.MahjongZH[d.players[0].selfWindTile])
		fmt.Println()
		info := fmt.Sprintln(util.TilesToMahjongZHInterface(d.doraIndicators)...)
		info = info[:len(info)-1]
		color.HiYellow("宝牌指示牌是 " + info)
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

			if debugMode {
				if who == 0 {
					if handsCount := util.CountOfTiles34(d.counts); handsCount%3 != 1 {
						return fmt.Errorf("手牌错误：%d 张牌 %v", handsCount, d.counts)
					}
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
					currentRoundCache.addChiPonKan(meldType)
				}
			}

			if debugMode {
				if meldType == meldTypeMinkan || meldType == meldTypeAnkan {
					if handsCount := util.CountOfTiles34(d.counts); handsCount%3 != 1 {
						return fmt.Errorf("手牌错误：%d 张牌 %v", handsCount, d.counts)
					}
				} else {
					if handsCount := util.CountOfTiles34(d.counts); handsCount%3 != 2 {
						return fmt.Errorf("手牌错误：%d 张牌 %v", handsCount, d.counts)
					}
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
			bestAttackDiscardTile := simpleBestDiscardTile(playerInfo)
			bestDefenceDiscardTile := mixedRiskTable.getBestDefenceTile(playerInfo.HandTiles34)
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
				currentRoundCache.addSelfDiscardTile(discardTile, mixedRiskTable[discardTile], isReach)
			}

			if debugMode {
				if handsCount := util.CountOfTiles34(d.counts); handsCount%3 != 1 {
					return fmt.Errorf("手牌错误：%d 张牌 %v", handsCount, d.counts)
				}
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
			allowChi := who == 3
			_, results14, incShantenResults14 := util.CalculateMeld(playerInfo, discardTile, isRedFive, allowChi)
			bestAttackDiscardTile := -1
			if len(results14) > 0 {
				bestAttackDiscardTile = results14[0].DiscardTile
			} else if len(incShantenResults14) > 0 {
				bestAttackDiscardTile = incShantenResults14[0].DiscardTile
			}
			if bestAttackDiscardTile != -1 {
				bestDefenceDiscardTile := mixedRiskTable.getBestDefenceTile(playerInfo.HandTiles34)
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
		//if d.gameMode == gameModeMatch && who == 3 && !canBeMeld {
		//	return nil
		//}

		if !debugMode {
			clearConsole()
		}

		// 牌谱模式下，打印舍牌推荐
		if d.gameMode == gameModeRecord {
			currentRoundCache.print()
		}

		// 打印他家舍牌信息
		d.printDiscards()
		fmt.Println()
		riskTables.printWithHands(d.counts, d.leftCounts)

		if d.gameMode == gameModeMatch && !canBeMeld {
			return nil
		}

		// 为了方便解析牌谱，这里尽可能地解析副露
		// TODO: 提醒: 消除海底/避免河底
		allowChi := d.playerNumber != 3 && who == 3 && playerInfo.LeftDrawTilesCount > 0
		return analysisMeld(playerInfo, discardTile, isRedFive, allowChi, mixedRiskTable)
	case d.parser.IsRoundWin():
		// TODO: 解析天凤牌谱 - 注意 skipOutput

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
	case d.parser.IsRyuukyoku():
		// TODO
		d.parser.ParseRyuukyoku()
	case d.parser.IsNukiDora():
		who, isTsumogiri := d.parser.ParseNukiDora()
		player := d.players[who]
		player.nukiDoraNum++
		if who != 0 {
			// 减少北的数量
			d.descLeftCounts(30)
			// TODO
			_ = isTsumogiri
		} else {
			// 减少自己手牌中北的数量
			d.counts[30]--
		}
		// 消除一发
		for _, player := range d.players {
			player.canIppatsu = false
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
