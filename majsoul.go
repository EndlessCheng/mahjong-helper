package main

import (
	"fmt"
)

type majsoulMessage struct {
	// NotifyPlayerLoadGameReady
	ReadyIDList []int `json:"ready_id_list,omitempty"`

	// ActionNewRound
	// {"chang":0,"ju":0,"ben":0,"tiles":["4m","6m","7m","3p","6p","7p","6s","1z","1z","2z","3z","4z","7z"],"dora":"6m","scores":[25000,25000,25000,25000],"liqibang":0,"al":false,"md5":"7527BD6868BBAB75B02A80CEA7CB4405","left_tile_count":69}
	MD5   string      `json:"md5,omitempty"`
	Chang *int        `json:"chang,omitempty"`
	Ju    *int        `json:"ju,omitempty"`
	Tiles interface{} `json:"tiles,omitempty"` // 一般情况下为 []interface{}, interface{} 即 string，但是暗杠的情况下，该值为一个 string
	Dora  string      `json:"dora,omitempty"`

	// ActionDealTile
	// {"seat":1,"tile":"5m","left_tile_count":64,"operation":{"seat":1,"operation_list":[{"type":1}],"time_add":0,"time_fixed":60000},"zhenting":false}
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

	// ActionLiuJu

	// ActionBabei
}

type majsoulRoundData struct {
	*roundData
	playerID   int
	originJSON string
	msg        *majsoulMessage
	seat       int // 初始座位：0-第一局的东家 1-第一局的南家 2-第一局的西家 3-第一局的北家
}

func (d *majsoulRoundData) fatalParse(info string, msg string) {
	panic(fmt.Sprintln(info, len(msg), msg, []byte(msg)))
}

func (d *majsoulRoundData) normalTiles(tiles interface{}) []string {
	_tiles, ok := tiles.([]interface{})
	if !ok {
		// 是否为暗杠？
		_tile, ok := tiles.(string)
		if !ok {
			panic(fmt.Sprintln("[anKanTile] 解析错误", tiles))
		}
		return []string{_tile}
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

//var tileReg = regexp.MustCompile("[0-9][mps]|[1-7]z")
//
//func (d *majsoulRoundData) extractTiles(msg string) (positions []int, tiles []int) {
//	indexPairs := tileReg.FindAllStringSubmatchIndex(msg, -1)
//	for _, pair := range indexPairs {
//		positions = append(positions, pair[0])
//		rawTile := msg[pair[0]:pair[1]]
//		tile := d.mustParseMajsoulTile(rawTile)
//		tiles = append(tiles, tile)
//	}
//	return
//}

func (d *majsoulRoundData) GetDataSourceType() int {
	return dataSourceTypeMajsoul
}

func (d *majsoulRoundData) GetMessage() string {
	return d.originJSON
}

func (d *majsoulRoundData) IsInit() bool {
	msg := d.msg

	const playersCount = 4

	if len(msg.ReadyIDList) > 0 && len(msg.ReadyIDList) < playersCount {
		fmt.Printf("等待玩家准备 (%d/%d)\n", len(msg.ReadyIDList), playersCount)
		return false
	}

	if d.playerID == -1 && len(msg.ReadyIDList) == playersCount {
		// 判断是否为人机对战，获取玩家ID
		if !inIntSlice(0, msg.ReadyIDList) {
			fmt.Println("尚未获取到玩家服务器端 ID，请先开启一局人机对战")
			return false
		}
		for _, playerID := range msg.ReadyIDList {
			if playerID != 0 {
				d.playerID = playerID
				fmt.Println("玩家服务器端 ID 获取成功:", d.playerID)
				break
			}
		}
	}

	if d.playerID == -1 {
		return false
	}

	// NotifyPlayerLoadGameReady || ActionNewRound
	return len(msg.ReadyIDList) == playersCount || msg.MD5 != ""
}

func (d *majsoulRoundData) ParseInit() (roundNumber int, dealer int, doraIndicator int, hands []int) {
	msg := d.msg

	if len(msg.ReadyIDList) == 4 {
		// dealer: 0=自家, 1=下家, 2=对家, 3=上家
		dealer = 1
		for i := len(msg.ReadyIDList) - 1; i >= 0; i-- {
			if msg.ReadyIDList[i] == d.playerID {
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

	if d.playerID == -1 {
		return false
	}

	if msg.Seat == nil || msg.Moqie != nil {
		return false
	}

	// ActionDealTile
	who := d.parseWho(*msg.Seat)
	return who == 0
}

func (d *majsoulRoundData) ParseSelfDraw() (tile int, kanDoraIndicator int) {
	msg := d.msg
	tile = d.mustParseMajsoulTile(msg.Tile)
	kanDoraIndicator = -1
	if len(msg.Doras) > 0 {
		kanDoraIndicator = d.mustParseMajsoulTile(msg.Doras[0])
	}
	return
}

func (d *majsoulRoundData) IsDiscard() bool {
	msg := d.msg

	if d.playerID == -1 {
		return false
	}

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
	if len(msg.Doras) > 0 {
		kanDoraIndicator = d.mustParseMajsoulTile(msg.Doras[0])
	}
	return
}

func (d *majsoulRoundData) IsOpen() bool {
	msg := d.msg

	if d.playerID == -1 {
		return false
	}

	// ActionChiPengGang || ActionAnGangAddGang
	return msg.Tiles != nil && len(d.normalTiles(msg.Tiles)) <= 4
}

func (d *majsoulRoundData) ParseOpen() (who int, meldType int, meldTiles []int, calledTile int, kanDoraIndicator int) {
	msg := d.msg

	who = d.parseWho(*msg.Seat)
	meldTiles = d.mustParseMajsoulTiles(d.normalTiles(msg.Tiles))
	kanDoraIndicator = -1
	if len(msg.Doras) > 0 {
		kanDoraIndicator = d.mustParseMajsoulTile(msg.Doras[0])
		calledTile = d.mustParseMajsoulTile(d.normalTiles(msg.Tiles)[0])
		if d.leftCounts[calledTile] == 4 {
			meldType = meldTypeAnKan
		} else {
			meldType = meldTypeKakan
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
		// 通过判断 calledTile 的来源来是否为上一张舍牌，来判断是明杠还是暗杠
		if len(d.globalDiscardTiles) > 0 && calledTile == d.globalDiscardTiles[len(d.globalDiscardTiles)-1] {
			// 明杠
			meldType = meldTypeMinKan
		} else {
			// 暗杠
			meldType = meldTypeAnKan
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

func (d *majsoulRoundData) IsNewDora() bool {
	return false
}

func (d *majsoulRoundData) ParseNewDora() (kanDoraIndicator int) {
	return 0
}
