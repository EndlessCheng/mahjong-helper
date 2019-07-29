package main

import (
	"fmt"
					)

type _majsoulRecordAccount struct {
	AccountID int `json:"account_id"`
	// 初始座位：0-第一局的东家 1-第一局的南家 2-第一局的西家 3-第一局的北家
	Seat     int    `json:"seat"` // *重点是拿到自己的座位
	Nickname string `json:"nickname"`
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
