package main

import (
	"fmt"
	"time"
	"github.com/EndlessCheng/mahjong-helper/util"
	"strconv"
)

// 观战基本信息
type majsoulLiveRecordBaseInfo struct {
	UUID      string `json:"uuid"`
	StartTime int64  `json:"start_time"`

	GameConfig *majsoulGameConfig `json:"game_config"`

	Players  []_majsoulRecordAccount `json:"players"`
	SeatList []int                   `json:"seat_list"`
}

func (i *majsoulLiveRecordBaseInfo) String() string {
	const timeFormat = "2006-01-02 15:04:05"
	output := fmt.Sprintf("%s\n开始于 %s\n\n", i.UUID, time.Unix(i.StartTime, 0).Format(timeFormat))

	maxAccountID := 0
	for _, account := range i.Players {
		maxAccountID = util.MaxInt(maxAccountID, account.AccountID)
	}
	accountShownWidth := len(strconv.Itoa(maxAccountID))
	for _, account := range i.Players {
		output += fmt.Sprintf("%*d %s\n", accountShownWidth, account.AccountID, account.Nickname)
	}
	return output
}
