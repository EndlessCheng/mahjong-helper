package main

import (
	"fmt"
	"time"
	"github.com/fatih/color"
	"github.com/EndlessCheng/mahjong-helper/util"
)

type majsoulMessage struct {
	// 对应到服务器用户数据库中的ID，该值越小表示您的注册时间越早
	AccountID int `json:"account_id"`

	Friends []*majsoulFriend `json:"friends"`

	// ResAuthGame
	IsGameStart *bool `json:"is_game_start"` // false=新游戏，true=重连
	SeatList    []int `json:"seat_list"`
	ReadyIDList []int `json:"ready_id_list"`

	// NotifyPlayerLoadGameReady
	//ReadyIDList []int `json:"ready_id_list"`

	// ActionNewRound
	// {"chang":0,"ju":0,"ben":0,"tiles":["1m","3m","7m","3p","6p","7p","6s","1z","1z","2z","3z","4z","7z"],"dora":"6m","scores":[25000,25000,25000,25000],"liqibang":0,"al":false,"md5":"","left_tile_count":69}
	MD5   string      `json:"md5"`
	Chang *int        `json:"chang"`
	Ju    *int        `json:"ju"`
	Tiles interface{} `json:"tiles"` // 一般情况下为 []interface{}, interface{} 即 string，但是暗杠的情况下，该值为一个 string
	Dora  string      `json:"dora"`

	// ActionDealTile
	// {"seat":1,"tile":"5m","left_tile_count":23,"operation":{"seat":1,"operation_list":[{"type":1}],"time_add":0,"time_fixed":60000},"zhenting":false}
	// 他家暗杠后的摸牌
	// {"seat":1,"left_tile_count":3,"doras":["7m","0p"],"zhenting":false}
	Seat          *int     `json:"seat"`
	Tile          string   `json:"tile"`
	Doras         []string `json:"doras"` // 暗杠摸牌了，同时翻出杠宝牌指示牌
	LeftTileCount *int     `json:"left_tile_count"`

	// ActionDiscardTile
	// {"seat":0,"tile":"5z","is_liqi":false,"moqie":true,"zhenting":false,"is_wliqi":false}
	// {"seat":0,"tile":"1z","is_liqi":false,"operation":{"seat":1,"operation_list":[{"type":3,"combination":["1z|1z"]}],"time_add":0,"time_fixed":60000},"moqie":false,"zhenting":false,"is_wliqi":false}
	// 吃 碰 和
	// {"seat":0,"tile":"6p","is_liqi":false,"operation":{"seat":1,"operation_list":[{"type":2,"combination":["7p|8p"]},{"type":3,"combination":["6p|6p"]},{"type":9}],"time_add":0,"time_fixed":60000},"moqie":false,"zhenting":true,"is_wliqi":false}
	IsLiqi    *bool     `json:"is_liqi"`
	IsWliqi   *bool     `json:"is_wliqi"`
	Moqie     *bool     `json:"moqie"`
	Operation *struct{} `json:"operation"`

	// ActionChiPengGang || ActionAnGangAddGang
	// 他家吃 {"seat":0,"type":0,"tiles":["2s","3s","4s"],"froms":[0,0,3],"zhenting":false}
	// 他家碰 {"seat":1,"type":1,"tiles":["1z","1z","1z"],"froms":[1,1,0],"operation":{"seat":1,"operation_list":[{"type":1,"combination":["1z"]}],"time_add":0,"time_fixed":60000},"zhenting":false,"tingpais":[{"tile":"4m","zhenting":false,"infos":[{"tile":"6s","haveyi":true},{"tile":"6p","haveyi":true}]},{"tile":"7m","zhenting":false,"infos":[{"tile":"6s","haveyi":true},{"tile":"6p","haveyi":true}]}]}
	// 他家大明杠 {"seat":2,"type":2,"tiles":["3z","3z","3z","3z"],"froms":[2,2,2,0],"zhenting":false}
	// 他家加杠 {"seat":2,"type":2,"tiles":"3z"}
	// 他家暗杠 {"seat":2,"type":3,"tiles":"3s"}
	Type  int   `json:"type"`
	Froms []int `json:"froms"`

	// ActionLiqi

	// ActionHule
	Hules []struct {
		Seat          int  `json:"seat"`
		Zimo          bool `json:"zimo"`
		PointRong     int  `json:"point_rong"`
		PointZimoQin  int  `json:"point_zimo_qin"`
		PointZimoXian int  `json:"point_zimo_xian"`
	} `json:"hules"`

	// ActionLiuJu

	// ActionBabei
}

const (
	majsoulMeldTypeChi = iota
	majsoulMeldTypePon
	majsoulMeldTypeMinKanOrKaKan
	majsoulMeldTypeAnKan
)

type majsoulRoundData struct {
	*roundData

	originJSON string
	accountID  int
	seat       int // 初始座位：0-第一局的东家 1-第一局的南家 2-第一局的西家 3-第一局的北家
	msg        *majsoulMessage
}

func (d *majsoulRoundData) fatalParse(info string, msg string) {
	panic(fmt.Sprintln(info, len(msg), msg, []byte(msg)))
}

func (d *majsoulRoundData) normalTiles(tiles interface{}) (majsoulTiles []string) {
	_tiles, ok := tiles.([]interface{})
	if !ok {
		_tile, ok := tiles.(string)
		if !ok {
			panic(fmt.Sprintln("[normalTiles] 解析错误", tiles))
		}
		return []string{_tile}
	}

	majsoulTiles = make([]string, len(_tiles))
	for i, _tile := range _tiles {
		_t, ok := _tile.(string)
		if !ok {
			panic(fmt.Sprintln("[normalTiles] 解析错误", tiles))
		}
		majsoulTiles[i] = _t
	}
	return majsoulTiles
}

func (d *majsoulRoundData) parseWho(seat int) int {
	// 转换成 0=自家, 1=下家, 2=对家, 3=上家
	who := (seat + d.dealer - d.roundNumber%4 + 4) % 4
	return who
}

func (d *majsoulRoundData) mustParseMajsoulTile(humanTile string) (tile34 int, isRedFive bool) {
	if humanTile[0] == '0' {
		humanTile = "5" + humanTile[1:]
		isRedFive = true
	}
	tile34, err := util.StrToTile34(humanTile)
	if err != nil {
		panic(err)
	}
	return
}

func (d *majsoulRoundData) mustParseMajsoulTiles(majsoulTiles []string) (tiles []int, containRedFive bool) {
	var isRedFive bool
	tiles = make([]int, len(majsoulTiles))
	for i, majsoulTile := range majsoulTiles {
		tiles[i], isRedFive = d.mustParseMajsoulTile(majsoulTile)
		if isRedFive {
			containRedFive = true
		}
	}
	return
}

func (d *majsoulRoundData) isNewDora(doras []string) bool {
	return len(doras) > len(d.doraIndicators)
}

func (d *majsoulRoundData) GetDataSourceType() int {
	return dataSourceTypeMajsoul
}

func (d *majsoulRoundData) GetMessage() string {
	return d.originJSON
}

func (d *majsoulRoundData) CheckMessage() bool {
	msg := d.msg

	// 首先，获取玩家账号
	if msg.SeatList != nil {
		if d.accountID > 0 {
			// 有 accountID 时，检查 accountID 是否正确
			if !util.InInts(d.accountID, msg.SeatList) {
				color.HiRed("尚未正确获取到玩家账号 ID，请您刷新网页，或开启一局人机对战（错误信息：您的账号 ID %d 不在对战列表 %v 中）", d.accountID, msg.SeatList)
				return false
			}
		} else {
			// 判断是否为人机对战，若为人机对战，则获取 accountID
			if util.InInts(0, msg.SeatList) {
				for _, accountID := range msg.SeatList {
					if accountID > 0 {
						d.accountID = accountID
						printAccountInfo(accountID)
						time.Sleep(2 * time.Second)
					}
				}
			}
		}
	}

	// 没有账号直接 return false
	if d.accountID == -1 {
		return false
	}

	// 当自家准备好时（msg.SeatList == nil），打印准备信息
	if msg.SeatList == nil && msg.ReadyIDList != nil {
		fmt.Printf("等待玩家准备 (%d/%d) %v\n", len(msg.ReadyIDList), 4, msg.ReadyIDList)
	}

	// 筛去重连的消息，目前的程序不考虑重连的情况
	if msg.IsGameStart != nil && *msg.IsGameStart {
		return false
	}

	return true
}

func (d *majsoulRoundData) IsInit() bool {
	msg := d.msg
	// ResAuthGame || ActionNewRound
	const playerNumber = 4
	return len(msg.SeatList) == playerNumber || msg.MD5 != ""
}

func (d *majsoulRoundData) ParseInit() (roundNumber int, dealer int, doraIndicator int, hands []int) {
	msg := d.msg
	const playerNumber = 4

	if len(msg.SeatList) == playerNumber {
		// dealer: 0=自家, 1=下家, 2=对家, 3=上家
		dealer = 1
		for i := len(msg.SeatList) - 1; i >= 0; i-- {
			if msg.SeatList[i] == d.accountID {
				break
			}
			dealer++
		}
		dealer %= playerNumber
		return
	}
	dealer = -1

	roundNumber = playerNumber*(*msg.Chang) + *msg.Ju
	doraIndicator, _ = d.mustParseMajsoulTile(msg.Dora)
	majsoulTiles := d.normalTiles(msg.Tiles)
	hands, _ = d.mustParseMajsoulTiles(majsoulTiles)
	return
}

func (d *majsoulRoundData) IsSelfDraw() bool {
	msg := d.msg

	if msg.Seat == nil || msg.Moqie != nil || msg.Tile == "" {
		return false
	}

	// FIXME: 更好的判断？
	// ActionDealTile
	who := d.parseWho(*msg.Seat)
	return who == 0
}

func (d *majsoulRoundData) ParseSelfDraw() (tile int, kanDoraIndicator int) {
	msg := d.msg
	tile, _ = d.mustParseMajsoulTile(msg.Tile)
	kanDoraIndicator = -1
	if d.isNewDora(msg.Doras) {
		kanDoraIndicator, _ = d.mustParseMajsoulTile(msg.Doras[len(msg.Doras)-1])
	}
	return
}

func (d *majsoulRoundData) IsDiscard() bool {
	msg := d.msg
	// ActionDiscardTile
	return msg.Moqie != nil
}

func (d *majsoulRoundData) ParseDiscard() (who int, discardTile int, isTsumogiri bool, isReach bool, canBeMeld bool, kanDoraIndicator int) {
	msg := d.msg
	who = d.parseWho(*msg.Seat)
	discardTile, _ = d.mustParseMajsoulTile(msg.Tile)
	isTsumogiri = *msg.Moqie
	isReach = *msg.IsLiqi || *msg.IsWliqi
	canBeMeld = msg.Operation != nil
	kanDoraIndicator = -1
	if d.isNewDora(msg.Doras) {
		kanDoraIndicator, _ = d.mustParseMajsoulTile(msg.Doras[len(msg.Doras)-1])
	}
	return
}

func (d *majsoulRoundData) IsOpen() bool {
	msg := d.msg
	// FIXME: 更好的判断？
	// ActionChiPengGang || ActionAnGangAddGang
	if msg.Tiles == nil {
		return false
	}
	majsoulTiles := d.normalTiles(msg.Tiles)
	return len(majsoulTiles) <= 4
}

func (d *majsoulRoundData) ParseOpen() (who int, meld *mjMeld, kanDoraIndicator int) {
	msg := d.msg

	who = d.parseWho(*msg.Seat)

	kanDoraIndicator = -1
	if d.isNewDora(msg.Doras) { // 暗杠（有时会在玩家摸牌后才发送 doras，可能是因为需要考虑抢暗杠的情况）
		kanDoraIndicator, _ = d.mustParseMajsoulTile(msg.Doras[len(msg.Doras)-1])
	}

	var meldType, calledTile int

	majsoulTiles := d.normalTiles(msg.Tiles)
	isSelfKan := len(majsoulTiles) == 1 // 加杠/暗杠
	if isSelfKan {
		majsoulTile := majsoulTiles[0]
		majsoulTiles = []string{majsoulTile, majsoulTile, majsoulTile, majsoulTile}
	}
	meldTiles, containRedFive := d.mustParseMajsoulTiles(majsoulTiles)
	if isSelfKan {
		calledTile = meldTiles[0]
		// 也可以通过副露来判断是加杠还是暗杠，这里简单地用 msg.Type 判断
		if msg.Type == majsoulMeldTypeMinKanOrKaKan {
			meldType = meldTypeKakan // 加杠
		} else if msg.Type == majsoulMeldTypeAnKan {
			meldType = meldTypeAnKan // 暗杠
			if d.leftCounts[calledTile] != 4 {
				// TODO: 改成 panic?
				fmt.Println("暗杠数据解析错误！")
			}
		}
		meld = &mjMeld{
			meldType:       meldType,
			tiles:          meldTiles,
			calledTile:     calledTile,
			containRedFive: containRedFive,
		}
		return
	}

	if len(meldTiles) == 3 {
		if meldTiles[0] == meldTiles[1] {
			meldType = meldTypePon // 碰
			calledTile = meldTiles[0]
		} else {
			meldType = meldTypeChi // 吃
			calledTile = d.globalDiscardTiles[len(d.globalDiscardTiles)-1]
			if calledTile < 0 {
				calledTile = ^calledTile
			}
		}
	} else if len(meldTiles) == 4 {
		meldType = meldTypeMinKan // 大明杠
		calledTile = meldTiles[0]
	} else {
		panic("鸣牌数据解析失败！")
	}
	meld = &mjMeld{
		meldType:       meldType,
		tiles:          meldTiles,
		calledTile:     calledTile,
		containRedFive: containRedFive,
	}
	return
}

func (d *majsoulRoundData) IsReach() bool {
	return false
}

func (d *majsoulRoundData) ParseReach() (who int) {
	return 0
}

func (d *majsoulRoundData) IsFuriten() bool {
	return false
}

func (d *majsoulRoundData) IsRoundWin() bool {
	msg := d.msg
	// ActionHule
	return msg.Hules != nil
}

func (d *majsoulRoundData) ParseRoundWin() (whos []int, points []int) {
	msg := d.msg

	for _, result := range msg.Hules {
		who := d.parseWho(result.Seat)
		whos = append(whos, d.parseWho(result.Seat))
		point := result.PointRong
		if result.Zimo {
			if who == d.dealer {
				point = 3 * result.PointZimoXian
			} else {
				point = result.PointZimoQin + 2*result.PointZimoXian
			}
		}
		points = append(points, point)
	}
	return
}

func (d *majsoulRoundData) IsNewDora() bool {
	msg := d.msg
	// 在最后处理该项
	// ActionDealTile
	return d.isNewDora(msg.Doras)
}

func (d *majsoulRoundData) ParseNewDora() (kanDoraIndicator int) {
	msg := d.msg

	kanDoraIndicator, _ = d.mustParseMajsoulTile(msg.Doras[len(msg.Doras)-1])
	return
}
