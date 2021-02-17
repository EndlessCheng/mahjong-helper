package main

import (
	"fmt"
	"github.com/EndlessCheng/mahjong-helper/util"
	"strconv"
	"time"
	"sort"
)

type _majsoulRecordAccount struct {
	AccountID int `json:"account_id"`
	// 初始座位：0-第一局的东家 1-第一局的南家 2-第一局的西家 3-第一局的北家
	Seat     int    `json:"seat"` // *重点是拿到自己的座位
	Nickname string `json:"nickname"`
}

// 牌谱基本信息
type majsoulRecordBaseInfo struct {
	UUID      string `json:"uuid"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`

	Config *majsoulGameConfig `json:"config"`

	Accounts []_majsoulRecordAccount `json:"accounts"`
}

func (i *majsoulRecordBaseInfo) sort() {
	sort.Slice(i.Accounts, func(i_, j int) bool {
		return i.Accounts[i_].Seat < i.Accounts[j].Seat
	})
}

var seatNameZH = []string{"东", "南", "西", "北"}

func (i *majsoulRecordBaseInfo) String() string {
	i.sort()

	const timeFormat = "2006-01-02 15:04:05"
	output := fmt.Sprintf("%s\n从 %s\n到 %s\n\n", i.UUID, time.Unix(i.StartTime, 0).Format(timeFormat), time.Unix(i.EndTime, 0).Format(timeFormat))

	maxAccountID := 0
	for _, account := range i.Accounts {
		maxAccountID = util.MaxInt(maxAccountID, account.AccountID)
	}
	accountShownWidth := len(strconv.Itoa(maxAccountID))
	for _, account := range i.Accounts {
		output += fmt.Sprintf("%s %*d %s\n", seatNameZH[account.Seat], accountShownWidth, account.AccountID, account.Nickname)
	}
	return output
}

func (i *majsoulRecordBaseInfo) getSelfSeat(accountID int) (int, error) {
	if len(i.Accounts) == 0 {
		return -1, fmt.Errorf("牌谱基本信息为空")
	}
	for _, account := range i.Accounts {
		if account.AccountID == accountID {
			return account.Seat, nil
		}
	}
	// 若没有，则以东家为主视角
	return 0, nil
}

//

// 牌谱、观战中的单个操作信息
type majsoulRecordAction struct {
	Name   string          `json:"name"`
	Action *majsoulMessage `json:"data"`
}

type majsoulRoundActions []*majsoulRecordAction

func (l majsoulRoundActions) append(action *majsoulRecordAction) (majsoulRoundActions, error) {
	if action == nil {
		return nil, fmt.Errorf("数据异常：拿到的操作内容为空")
	}
	newL := l

	if action.Name == "RecordNewRound" {
		newL = majsoulRoundActions{action}
	} else {
		if len(newL) == 0 {
			return nil, fmt.Errorf("数据异常：未收到 RecordNewRound")
		}
		newL = append(newL, action)
	}

	return newL, nil
}

func parseMajsoulRecordAction(actions []*majsoulRecordAction) (roundActionsList []majsoulRoundActions, err error) {
	if len(actions) == 0 {
		return nil, fmt.Errorf("数据异常：拿到的牌谱内容为空")
	}

	var currentRoundActions majsoulRoundActions
	for _, action := range actions {
		if action.Name == "RecordNewRound" {
			if len(currentRoundActions) > 0 {
				roundActionsList = append(roundActionsList, currentRoundActions)
			}
			currentRoundActions = []*majsoulRecordAction{action}
		} else {
			if len(currentRoundActions) == 0 {
				return nil, fmt.Errorf("数据异常：未收到 RecordNewRound")
			}
			currentRoundActions = append(currentRoundActions, action)
		}
	}
	if len(currentRoundActions) > 0 {
		roundActionsList = append(roundActionsList, currentRoundActions)
	}
	return
}
