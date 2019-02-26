package main

import (
	"strings"
	"fmt"
)

type majsoulMessage string

type majsoulRoundData struct {
	roundData
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
	return strings.Contains(msg, "ActionNewRound")
}

func (d *majsoulRoundData) ParseInit() (roundNumber int, dealer int, doraIndicator int, hands []int) {
	msg := string(*d.msg)
	shift := 0
	if msg[45] == 1 {
		shift = 1
	}
	isDealer := shift == 1
	fmt.Println("isDealer", isDealer)

	// TODO
	dealer = -1

	// TODO: 开局九种九牌、w立、暗杠、天和
	roundType := int(msg[46+shift])
	normalRoundNumber := int(msg[48+shift])
	roundNumber = 4*roundType + normalRoundNumber

	handShift := 0
	if isDealer {
		// 庄家开局有14张牌
		handShift = 5
	}
	hands = make([]int, 34)
	for _, rawTile := range strings.Split(msg[53+handShift:103+handShift], string([]byte{34, 2})) {
		tile := d._mustParseMajsoulTile(rawTile)
		hands[tile]++
	}

	doraIndicator = d._mustParseMajsoulTile(msg[105+handShift : 107+handShift])

	return
}

func (d *majsoulRoundData) IsSelfDraw() bool {
	msg := string(*d.msg)
	const otherDrawMsgLength = 50
	return len(msg) > otherDrawMsgLength && strings.Contains(msg, "ActionDealTile")
}

func (d *majsoulRoundData) ParseSelfDraw() (tile int) {
	msg := string(*d.msg)
	// 含有摸到的牌，若有加杠、暗杠选项会更长
	tile = d._mustParseMajsoulTile(msg[48:50])
	return
}

func (d *majsoulRoundData) IsDiscard() bool {
	msg := string(*d.msg)
	return strings.Contains(msg, "ActionDiscardTile")
}

func (d *majsoulRoundData) ParseDiscard() (who int, tile int, isTsumogiri bool, isReach bool, canBeMeld bool, kanDoraIndicator int) {
	msg := string(*d.msg)

	who = int(msg[48])
	tile = d._mustParseMajsoulTile(msg[51:53])
	isTsumogiri = msg[len(msg)-5] == 1
	isReach = msg[54] == 1

	const normalActionDiscardLength = 61
	if len(msg) > normalActionDiscardLength {
		// TODO: 用who是否为自家来判断
		// TODO: 如果明杠后的舍牌又恰好能鸣牌，会是什么样的格式？
		if dora, err := d._parseMajsoulTile(msg[len(msg)-4 : len(msg)-2]); err == nil {
			kanDoraIndicator = dora
		}

		// TODO: 待验证
		if len(msg) > 67 && msg[67] == '|' {
			canBeMeld = true
		}
	}

	return
}

func (d *majsoulRoundData) IsOpen() bool {
	msg := string(*d.msg)
	return strings.Contains(msg, "ActionChiPengGang")
}

func (d *majsoulRoundData) ParseOpen() (who int, meldType int, meldTiles []int, calledTile int) {
	msg := string(*d.msg)

	// TODO: 待验证
	who = int(msg[48])

	// TODO: 如何判断加杠？
	var rawMeldTiles string
	if msg[63] == '"' {
		rawMeldTiles = msg[53:63]
	} else if msg[67] == '"' {
		rawMeldTiles = msg[53:67]
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
