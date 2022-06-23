package main

import "github.com/EndlessCheng/mahjong-helper/util"

func newPlayerInfo(name string, selfWindTile int) *PlayerInfo {
	return &PlayerInfo{
		name:                  name,
		selfWindTile:          selfWindTile,
		latestDiscardAtGlobal: -1,
		reachTileAtGlobal:     -1,
		reachTileAt:           -1,
	}
}

func modifySanninPlayerInfoList(lst []*PlayerInfo, roundNumber int) []*PlayerInfo {
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

func newRoundData(parser DataParser, roundNumber int, benNumber int, dealer int) *RoundData {
	// 无论是三麻还是四麻，都视作四个人
	const playerNumber = 4
	roundWindTile := 27 + roundNumber/playerNumber
	playerWindTile := make([]int, playerNumber)
	for i := 0; i < playerNumber; i++ {
		playerWindTile[i] = 27 + (playerNumber-dealer+i)%playerNumber
	}
	return &RoundData{
		parser:      parser,
		roundNumber: roundNumber,
		benNumber:   benNumber,

		roundWindTile:      roundWindTile,
		dealer:             dealer,
		counts:             make([]int, 34),
		leftCounts:         util.InitLeftTiles34(),
		globalDiscardTiles: []int{},
		players: []*PlayerInfo{
			newPlayerInfo("自家", playerWindTile[0]),
			newPlayerInfo("下家", playerWindTile[1]),
			newPlayerInfo("對家", playerWindTile[2]),
			newPlayerInfo("上家", playerWindTile[3]),
		},
	}
}

func NewGame(parser DataParser) *RoundData {
	return newRoundData(parser, 0, 0, 0)
}
