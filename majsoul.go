package main

import (
	"strings"
)

type majsoulMessage string

type majsoulRoundData struct {
	*roundData
	msg *majsoulMessage
}

func (d *majsoulRoundData) _mustParseMajsoulTile(tile string) int {
	if tile[0] == '0' {
		tile = "5" + tile[1:]
	}
	idx, err := _convert(tile)
	if err != nil {
		panic(err)
	}
	return idx
}

func (d *majsoulRoundData) _parseMajsoulTile(tile string) (int, error) {
	if tile[0] == '0' {
		tile = "5" + tile[1:]
	}
	return _convert(tile)
}

func (d *majsoulRoundData) GetDataSourceType() int {
	return dataSourceTypeMajsoul
}

func (d *majsoulRoundData) GetMessage() string {
	return string(*d.msg)
}

func (d *majsoulRoundData) IsInit() bool {
	msg := string(*d.msg)
	return strings.Contains(msg, "ActionNewRound") || strings.Contains(msg, "NotifyPlayerLoadGameReady")
}

func (d *majsoulRoundData) ParseInit() (roundNumber int, dealer int, doraIndicator int, hands []int) {
	msg := string(*d.msg)

	if strings.Contains(msg, "NotifyPlayerLoadGameReady") {
		// dealer: 0=自家, 1=下家, 2=对家, 3=上家
		dealer = 1
		for i := len(msg) - 1; i >= 0; i-- {
			if msg[i] != 0 {
				break
			}
			dealer++
		}
		dealer %= 4
		return
	}
	dealer = -1

	shift := 0
	if msg[45] == 1 {
		shift = 1
	}
	isDealer := shift == 1

	// TODO: 开局九种九牌、w立、暗杠、天和
	roundType := int(msg[46+shift])
	normalRoundNumber := int(msg[48+shift])
	roundNumber = 4*roundType + normalRoundNumber

	handShift := 0
	if isDealer {
		// 庄家开局有14张牌
		handShift = 5
	}
	for _, rawTile := range strings.Split(msg[53+handShift:103+handShift], string([]byte{34, 2})) {
		tile := d._mustParseMajsoulTile(rawTile)
		hands = append(hands, tile)
	}

	doraIndicator = d._mustParseMajsoulTile(msg[105+handShift : 107+handShift])

	return
}

func (d *majsoulRoundData) IsSelfDraw() bool {
	msg := string(*d.msg)
	const otherDrawMsgLength = 60 // 50
	return len(msg) > otherDrawMsgLength && strings.Contains(msg, "ActionDealTile")
}

func (d *majsoulRoundData) ParseSelfDraw() (tile int) {
	msg := string(*d.msg)
	// 含有摸到的牌，若有加杠、暗杠选项会更长
	var err error
	tile, err = d._parseMajsoulTile(msg[48:50])
	if err != nil {
		tile = d._mustParseMajsoulTile(msg[49:51])
	}
	return
}

func (d *majsoulRoundData) IsDiscard() bool {
	msg := string(*d.msg)
	return strings.Contains(msg, "ActionDiscardTile")
}

func (d *majsoulRoundData) ParseDiscard() (who int, tile int, isTsumogiri bool, isReach bool, canBeMeld bool, kanDoraIndicator int) {
	msg := string(*d.msg)

	//splits := strings.Split(msg, "ActionDiscardTile")

	// 0-第一局的东家 1-第一局的南家 2-第一局的西家 3-第一局的北家
	majsoulWho := int(msg[48])
	// 转换成 0=自家, 1=下家, 2=对家, 3=上家
	who = (majsoulWho + d.dealer - d.roundNumber) % 4

	var err error
	shift := 0
	tile, err = d._parseMajsoulTile(msg[51:53])
	if err != nil {
		tile = d._mustParseMajsoulTile(msg[52:54])
		shift = 1
	}
	isTsumogiri = msg[len(msg)-5] == 1
	isReach = msg[54+shift] == 1
	kanDoraIndicator = -1

	const normalActionDiscardLength = 61
	if len(msg) > normalActionDiscardLength {
		// TODO: 用who是否为自家来判断
		// TODO: 如果明杠后的舍牌又恰好能鸣牌，会是什么样的格式？
		rawTile := msg[len(msg)-4 : len(msg)-2]
		if dora, err := d._parseMajsoulTile(rawTile); err == nil {
			kanDoraIndicator = dora
		}

		// TODO: 待验证
		if len(msg) > 67+shift && msg[67+shift] == '|' {
			canBeMeld = true
		}
	}

	return
}

// 他家暗杠
// .lq.ActionPrototypeActionAnGangAddGan8m
// 57
// [1 10 19 46 108 113 46 65 99 116 105 111 110 80 114 111 116 111 116 121 112 101 18 33 8 28 18 19 65 99 116 105 111 110 65 110 71 97 110 103 65 100 100 71 97 110 103 26 8 8 1 16 3 26 2 56 109]
// 自家加杠
//.lq.ActionPrototypevActionAnGangAddGang7zB6pB/3p
//73
//[1 10 19 46 108 113 46 65 99 116 105 111 110 80 114 111 116 111 116 121 112 101 18 49 8 118 18 19 65 99 116 105 111 110 65 110 71 97 110 103 65 100 100 71 97 110 103 26 24 8 2 16 2 26 2 55 122 66 6 10 2 54 112 16 1 66 6 10 2 51 112 16 1]
func (d *majsoulRoundData) IsOpen() bool {
	msg := string(*d.msg)
	return strings.Contains(msg, "ActionChiPengGang") || strings.Contains(msg, "ActionAnGangAddGan")
}

func (d *majsoulRoundData) ParseOpen() (who int, meldType int, meldTiles []int, calledTile int) {
	msg := string(*d.msg)

	if strings.Contains(msg, "ActionAnGangAddGan") {
		// 0-第一局的东家 1-第一局的南家 2-第一局的西家 3-第一局的北家
		majsoulWho := int(msg[50])
		// 转换成 0=自家, 1=下家, 2=对家, 3=上家
		who = (majsoulWho + d.dealer - d.roundNumber) % 4

		calledTile = d._mustParseMajsoulTile(msg[len(msg)-2:])
		if d.leftCounts[calledTile] == 4 {
			meldType = meldTypeAnKan
		} else {
			meldType = meldTypeKakan
		}
		return
	}

	// 0-第一局的东家 1-第一局的南家 2-第一局的西家 3-第一局的北家
	majsoulWho := int(msg[48])
	// 转换成 0=自家, 1=下家, 2=对家, 3=上家
	who = (majsoulWho + d.dealer - d.roundNumber) % 4

	var rawMeldTiles string
	if msg[63] == '"' {
		rawMeldTiles = msg[53:63]
	} else if msg[64] == '"' {
		rawMeldTiles = msg[54:64]
	} else if msg[67] == '"' {
		rawMeldTiles = msg[53:67]
	} else if msg[68] == '"' {
		rawMeldTiles = msg[54:68]
	} else {
		meldType = meldTypeKakan
		panic("解析失败（可能是加杠？）")
	}
	for _, rawTile := range strings.Split(rawMeldTiles, string([]byte{26, 2})) {
		tile := d._mustParseMajsoulTile(rawTile)
		meldTiles = append(meldTiles, tile)
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
	// TODO
	return false
}

func (d *majsoulRoundData) IsNewDora() bool {
	// TODO
	return false
}

func (d *majsoulRoundData) ParseNewDora() (kanDoraIndicator int) {
	// TODO
	return 0
}
