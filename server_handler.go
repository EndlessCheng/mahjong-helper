package main

import (
	"github.com/EndlessCheng/mahjong-helper/platform/tenhou/ws"
	"fmt"
	"github.com/EndlessCheng/mahjong-helper/debug"
	"net/url"
	"github.com/fatih/color"
	"github.com/EndlessCheng/mahjong-helper/platform/tenhou"
	"github.com/EndlessCheng/mahjong-helper/platform/majsoul"
	"github.com/EndlessCheng/mahjong-helper/platform/majsoul/proto/lq"
	"github.com/EndlessCheng/mahjong-helper/util"
)

func (h *mjHandler) handleTenhouWebSocketMessage(message *tenhou.Message) (err error) {
	if !debugMode {
		defer func() {
			if er := recover(); er != nil {
				err = fmt.Errorf("内部错误: %v", er)
			}
		}()
	}

	if !debugMode {
		h.log.Info(message.OriginJSON)
	} else {
		fmt.Println(debug.Lo, message.OriginJSON)
	}

	switch meta := message.Metadata.(type) {
	case *ws.Helo: // 登录
		username, er := url.QueryUnescape(meta.UserName)
		if er != nil {
			return er
		}
		if username != userConf.currentActiveTenhouUserName {
			userConf.currentActiveTenhouUserName = username
			color.HiGreen("%s 登录成功", username)
		}
	case *ws.UN: // 对战前的各家用户信息
		// 游戏配置：三麻/四麻
		h.tenhouRoundData.config.playerNumber = meta.PlayerNumber()
	default:
		h.tenhouRoundData.parser = message
		if er := h.tenhouRoundData.analysis(); er != nil {
			return fmt.Errorf("analysis: %v", er)
		}
	}
	return
}

func (h *mjHandler) loadAccountID(accountID uint32) {
	_accountID := int(accountID)
	if userConf.currentActiveMajsoulAccountID == _accountID {
		return
	}
	userConf.addMajsoulAccountID(_accountID)
	userConf.setMajsoulAccountID(_accountID)
	fmt.Print("您的账号 ID 为 ")
	color.New(color.FgHiGreen).Printf("%d", _accountID)
	fmt.Print("，该数字为雀魂服务器账号数据库中的 ID，该值越小表示您的注册时间越早\n")
}

func (h *mjHandler) handleMajsoulWebSocketMessage(message *majsoul.Message) (err error) {
	if !debugMode {
		defer func() {
			if er := recover(); er != nil {
				err = fmt.Errorf("内部错误: %v", er)
			}
		}()
	}

	if !debugMode {
		h.log.Info(message.JSON())
	} else {
		fmt.Println(debug.Lo, message.JSON())
	}

	// 把 ResponseMessage 和 NotifyMessage 混在一起，方便用同一个 switch 判断
	recvMsg := message.ResponseMessage
	if recvMsg == nil {
		recvMsg = message.NotifyMessage
	}

	switch msg := recvMsg.(type) {
	case *lq.ResLogin: // 登录
		h.loadAccountID(msg.AccountId)
	case *lq.ResFriendList: // 好友列表
		fmt.Println(lq.FriendList(msg.Friends))
	case *lq.ResAuthGame: // 对战前的各家用户信息
		defer func() {
			// 获取自家初始座位：0-第一局的东家 1-第一局的南家 2-第一局的西家 3-第一局的北家
			for i, accountID := range msg.SeatList {
				if int(accountID) == userConf.currentActiveMajsoulAccountID {
					h.majsoulRoundData.selfSeat = i
					break
				}
			}
			// dealer: 0=自家, 1=下家, 2=对家, 3=上家
			dealer := (4 - h.majsoulRoundData.selfSeat) % 4
			h.majsoulRoundData.reset(0, 0, dealer)
			h.majsoulRoundData.config.gameMode = gameModeMatch
			fmt.Printf("游戏即将开始，您分配到的座位是：")
			color.HiGreen(util.MahjongZH[h.majsoulRoundData.players[0].selfWindTile])
		}()

		// 游戏配置：三麻/四麻
		h.majsoulRoundData.config.playerNumber = len(msg.SeatList)

		// 尝试从中找到缓存账号 ID
		for _, accountID := range msg.SeatList {
			if accountID > 0 && userConf.isIDExist(int(accountID)) {
				h.loadAccountID(accountID)
				return
			}
		}

		// 未找到缓存 ID。若为人机对战，则获取账号 ID
		if msg.IsComputer() {
			for _, accountID := range msg.SeatList {
				if accountID > 0 {
					h.loadAccountID(accountID)
					return
				}
			}
		}

		color.HiRed("尚未获取到您的账号 ID，请您刷新网页，或开启一局人机对战")
	case *lq.NotifyPlayerLoadGameReady: // 玩家准备
		fmt.Printf("等待玩家准备 (%d/%d) %v\n", len(msg.ReadyIdList), h.majsoulRoundData.config.playerNumber, msg.ReadyIdList)
	case *lq.ActionPrototype: // 对战操作信息
		if userConf.currentActiveMajsoulAccountID == -1 {
			clearConsole()
			color.HiRed("尚未获取到您的账号 ID，请您刷新网页，或开启一局人机对战")
			return
		}
		action, er := majsoul.NewAction(msg)
		if er != nil {
			return er
		}
		h.majsoulRoundData.parser = action
		if er := h.majsoulRoundData.analysis(); er != nil {
			return fmt.Errorf("analysis: %v", er)
		}
	case *lq.ResGameRecordList: // 牌谱基本信息列表
		for _, record := range msg.RecordList {
			h.majsoulRecordGameMap[record.Uuid] = record
		}
		color.HiGreen("收到 %2d 个雀魂牌谱（已收集 %d 个），请在网页上点击「查看」", len(msg.RecordList), len(h.majsoulRecordGameMap))
	case *lq.ResGameRecord: // 载入某个牌谱（含分享）
		// 处理基础信息
		record := msg.Head
		h.majsoulRecordGameMap[record.Uuid] = record
		if er := h._loadMajsoulRecordBaseInfo(record.Uuid); er != nil {
			return er
		}
		// 处理东一局
	default:
		//
	}
	return
}

func (h *mjHandler) handleMajsoulUIMessage() (err error) {
	if !debugMode {
		defer func() {
			if er := recover(); er != nil {
				err = fmt.Errorf("内部错误: %v", er)
			}
		}()
	}

	return
}
