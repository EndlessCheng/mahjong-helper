package main

import (
	"fmt"
	"time"
	"github.com/fatih/color"
)

type majsoulMessage struct {
	// 对应到服务器用户数据库中的ID，该值越小表示您的注册时间越早
	AccountID int `json:"account_id"`

	Friends []*majsoulFriend `json:"friends"`

	// NotifyPlayerLoadGameReady
	ReadyIDList []int `json:"ready_id_list,omitempty"`

	// ActionNewRound
	// {"chang":0,"ju":0,"ben":0,"tiles":["1m","3m","7m","3p","6p","7p","6s","1z","1z","2z","3z","4z","7z"],"dora":"6m","scores":[25000,25000,25000,25000],"liqibang":0,"al":false,"md5":"","left_tile_count":69}
	MD5   string      `json:"md5,omitempty"`
	Chang *int        `json:"chang,omitempty"`
	Ju    *int        `json:"ju,omitempty"`
	Tiles interface{} `json:"tiles,omitempty"` // 一般情况下为 []interface{}, interface{} 即 string，但是暗杠的情况下，该值为一个 string
	Dora  string      `json:"dora,omitempty"`

	// ActionDealTile
	// {"seat":1,"tile":"5m","left_tile_count":23,"operation":{"seat":1,"operation_list":[{"type":1}],"time_add":0,"time_fixed":60000},"zhenting":false}
	// 他家暗杠后的摸牌
	// {"seat":1,"left_tile_count":3,"doras":["7m","0p"],"zhenting":false}
	Seat          *int     `json:"seat,omitempty"`
	Tile          string   `json:"tile,omitempty"`
	Doras         []string `json:"doras,omitempty"` // 暗杠摸牌了，同时翻出杠宝牌指示牌
	LeftTileCount *int     `json:"left_tile_count,omitempty"`

	// ActionDiscardTile
	// {"seat":0,"tile":"5z","is_liqi":false,"moqie":true,"zhenting":false,"is_wliqi":false}
	// {"seat":0,"tile":"1z","is_liqi":false,"operation":{"seat":1,"operation_list":[{"type":3,"combination":["1z|1z"]}],"time_add":0,"time_fixed":60000},"moqie":false,"zhenting":false,"is_wliqi":false}
	// 吃 碰 和
	// {"seat":0,"tile":"6p","is_liqi":false,"operation":{"seat":1,"operation_list":[{"type":2,"combination":["7p|8p"]},{"type":3,"combination":["6p|6p"]},{"type":9}],"time_add":0,"time_fixed":60000},"moqie":false,"zhenting":true,"is_wliqi":false}
	IsLiqi    *bool     `json:"is_liqi,omitempty"`
	IsWliqi   *bool     `json:"is_wliqi,omitempty"`
	Moqie     *bool     `json:"moqie,omitempty"`
	Operation *struct{} `json:"operation,omitempty"`

	// ActionChiPengGang || ActionAnGangAddGang
	// {"seat":1,"type":1,"tiles":["1z","1z","1z"],"froms":[1,1,0],"operation":{"seat":1,"operation_list":[{"type":1,"combination":["1z"]}],"time_add":0,"time_fixed":60000},"zhenting":false,"tingpais":[{"tile":"4m","zhenting":false,"infos":[{"tile":"6s","haveyi":true},{"tile":"6p","haveyi":true}]},{"tile":"7m","zhenting":false,"infos":[{"tile":"6s","haveyi":true},{"tile":"6p","haveyi":true}]}]}
	Froms []int `json:"froms,omitempty"`

	// ActionLiqi

	// ActionHule
	Hules []struct {
		Seat      int `json:"seat"`
		PointRong int `json:"point_rong"`
	} `json:"hules"`

	// ActionLiuJu

	// ActionBabei
}

type majsoulFriend struct {
	PlayerBaseView struct {
		AccountID int    `json:"account_id"`
		Nickname  string `json:"nickname"`
	} `json:"base"`

	AccountActiveState struct {
		LoginTime  int64 `json:"login_time"`
		LogoutTime int64 `json:"logout_time"`
		IsOnline   bool  `json:"is_online"`
	} `json:"state"`
}

func (f *majsoulFriend) String() string {
	return fmt.Sprintf("%9d   %s   %s   %s",
		f.PlayerBaseView.AccountID,
		time.Unix(f.AccountActiveState.LoginTime, 0).Format("2006-01-02 15:04:05"),
		time.Unix(f.AccountActiveState.LogoutTime, 0).Format("2006-01-02 15:04:05"),
		f.PlayerBaseView.Nickname,
	)
}

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

func (d *majsoulRoundData) normalTiles(tiles interface{}) []string {
	_tiles, ok := tiles.([]interface{})
	if !ok {
		_tile, ok := tiles.(string)
		if !ok {
			panic(fmt.Sprintln("[normalTiles] 解析错误", tiles))
		}
		return []string{_tile, _tile, _tile, _tile}
	}

	results := make([]string, len(_tiles))
	for i, _tile := range _tiles {
		_t, ok := _tile.(string)
		if !ok {
			panic(fmt.Sprintln("[normalTiles] 解析错误", tiles))
		}
		results[i] = _t
	}
	return results
}

func (d *majsoulRoundData) parseWho(seat int) int {
	// 转换成 0=自家, 1=下家, 2=对家, 3=上家
	who := (seat + d.dealer - d.roundNumber%4 + 4) % 4
	return who
}

func (d *majsoulRoundData) mustParseMajsoulTile(tile string) int {
	if tile[0] == '0' {
		tile = "5" + tile[1:]
	}
	idx, err := _convert(tile)
	if err != nil {
		panic(err)
	}
	return idx
}

func (d *majsoulRoundData) parseMajsoulTile(tile string) (int, error) {
	if tile[0] == '0' {
		tile = "5" + tile[1:]
	}
	return _convert(tile)
}

func (d *majsoulRoundData) mustParseMajsoulTiles(tiles []string) []int {
	hands := make([]int, len(tiles))
	for i, tile := range tiles {
		hands[i] = d.mustParseMajsoulTile(tile)
	}
	return hands
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

	if msg.ReadyIDList != nil {
		fmt.Printf("等待玩家准备 (%d/%d) %v\n", len(msg.ReadyIDList), 4, msg.ReadyIDList)
	}

	if d.accountID > 0 {
		if msg.ReadyIDList != nil && !inIntSlice(d.accountID, msg.ReadyIDList) {
			color.Red("尚未正确获取到玩家数据库账号 ID，请您刷新网页，或开启一局人机对战（错误信息：您的 ID %d 不在对战列表 %v 中）", d.accountID, msg.ReadyIDList)
			return false
		}
		return true
	}
	// 尚未获取到玩家数据库账号 ID

	if len(msg.ReadyIDList) < 4 {
		return false
	}

	// 判断是否为人机对战，获取玩家ID
	if !inIntSlice(0, msg.ReadyIDList) {
		color.Red("尚未获取到玩家数据库账号 ID，请您刷新网页，或开启一局人机对战")
		return false
	}

	for _, accountID := range msg.ReadyIDList {
		if accountID > 0 {
			d.accountID = accountID
			printAccountInfo(accountID)
			time.Sleep(2 * time.Second)
			return true
		}
	}

	return false
}

func (d *majsoulRoundData) IsInit() bool {
	msg := d.msg
	// NotifyPlayerLoadGameReady || ActionNewRound
	return len(msg.ReadyIDList) == 4 || msg.MD5 != ""
}

func (d *majsoulRoundData) ParseInit() (roundNumber int, dealer int, doraIndicator int, hands []int) {
	msg := d.msg

	if len(msg.ReadyIDList) == 4 {
		// dealer: 0=自家, 1=下家, 2=对家, 3=上家
		dealer = 1
		for i := len(msg.ReadyIDList) - 1; i >= 0; i-- {
			if msg.ReadyIDList[i] == d.accountID {
				break
			}
			dealer++
		}
		dealer %= 4
		return
	}
	dealer = -1

	roundNumber = 4**msg.Chang + *msg.Ju
	doraIndicator = d.mustParseMajsoulTile(msg.Dora)
	hands = d.mustParseMajsoulTiles(d.normalTiles(msg.Tiles))
	return
}

func (d *majsoulRoundData) IsSelfDraw() bool {
	msg := d.msg

	if msg.Seat == nil || msg.Moqie != nil {
		return false
	}

	// FIXME: 更好的判断？
	// ActionDealTile
	who := d.parseWho(*msg.Seat)
	return who == 0
}

func (d *majsoulRoundData) ParseSelfDraw() (tile int, kanDoraIndicator int) {
	msg := d.msg
	tile = d.mustParseMajsoulTile(msg.Tile)
	kanDoraIndicator = -1
	if d.isNewDora(msg.Doras) {
		kanDoraIndicator = d.mustParseMajsoulTile(msg.Doras[len(msg.Doras)-1])
	}
	return
}

func (d *majsoulRoundData) IsDiscard() bool {
	msg := d.msg
	// ActionDiscardTile
	return msg.Moqie != nil
}

func (d *majsoulRoundData) ParseDiscard() (who int, tile int, isTsumogiri bool, isReach bool, canBeMeld bool, kanDoraIndicator int) {
	msg := d.msg
	who = d.parseWho(*msg.Seat)
	tile = d.mustParseMajsoulTile(msg.Tile)
	isTsumogiri = *msg.Moqie
	isReach = *msg.IsLiqi || *msg.IsWliqi
	canBeMeld = msg.Operation != nil
	kanDoraIndicator = -1
	if d.isNewDora(msg.Doras) {
		kanDoraIndicator = d.mustParseMajsoulTile(msg.Doras[len(msg.Doras)-1])
	}
	return
}

func (d *majsoulRoundData) IsOpen() bool {
	msg := d.msg
	// FIXME: 更好的判断？
	// ActionChiPengGang || ActionAnGangAddGang
	return msg.Tiles != nil && len(d.normalTiles(msg.Tiles)) <= 4
}

func (d *majsoulRoundData) ParseOpen() (who int, meldType int, meldTiles []int, calledTile int, kanDoraIndicator int) {
	msg := d.msg

	who = d.parseWho(*msg.Seat)
	meldTiles = d.mustParseMajsoulTiles(d.normalTiles(msg.Tiles))
	kanDoraIndicator = -1
	if d.isNewDora(msg.Doras) {
		kanDoraIndicator = d.mustParseMajsoulTile(msg.Doras[len(msg.Doras)-1])

		meldType = meldTypeAnKan
		calledTile = d.mustParseMajsoulTile(d.normalTiles(msg.Tiles)[0])
		if d.leftCounts[calledTile] != 4 {
			// TODO: 改成 panic
			fmt.Println("暗杠数据解析错误！")
		}
		return
	}

	if len(meldTiles) == 3 {
		if meldTiles[0] == meldTiles[1] {
			meldType = meldTypePon
			calledTile = meldTiles[0]
		} else {
			meldType = meldTypeChi
			calledTile = d.globalDiscardTiles[len(d.globalDiscardTiles)-1]
		}
	} else if len(meldTiles) == 4 {
		calledTile = meldTiles[0]
		// 通过判断 calledTile 的来源来是否为上一张舍牌，来判断是大明杠还是加杠
		if len(d.globalDiscardTiles) > 0 && calledTile == d.globalDiscardTiles[len(d.globalDiscardTiles)-1] {
			// 大明杠
			meldType = meldTypeMinKan
		} else {
			// 加杠
			meldType = meldTypeKakan
		}
	} else {
		panic("鸣牌数据解析失败！")
	}
	if calledTile < 0 {
		calledTile = ^calledTile
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
		whos = append(whos, d.parseWho(result.Seat))
		points = append(points, result.PointRong)
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

	kanDoraIndicator = d.mustParseMajsoulTile(msg.Doras[len(msg.Doras)-1])
	return
}
