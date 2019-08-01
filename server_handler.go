package main

import (
	"github.com/EndlessCheng/mahjong-helper/platform/tenhou/ws"
	"fmt"
	"github.com/EndlessCheng/mahjong-helper/debug"
	"net/url"
	"github.com/fatih/color"
	"github.com/EndlessCheng/mahjong-helper/platform/tenhou"
)

func (h *mjHandler) handleTenhouMessage(msg *tenhou.Message) {
	if !debugMode {
		h.log.Info(msg.OriginJSON)
	} else {
		fmt.Println(debug.Lo, msg.OriginJSON)
	}

	switch meta := msg.Metadata.(type) {
	case *ws.Helo: // 用户登录
		username, err := url.QueryUnescape(meta.UserName)
		if err != nil {
			h.logError(err)
		}
		if username != gameConf.currentActiveTenhouUserName {
			color.HiGreen("%s 登录成功", username)
			gameConf.currentActiveTenhouUserName = username
		}
	case *ws.UN: // 对战前的各家用户信息
		// 游戏配置：三麻/四麻
		h.tenhouRoundData.playerNumber = meta.PlayerNumber()
	default:
		h.tenhouRoundData.parser = msg
		if err := h.tenhouRoundData.analysis(); err != nil {
			h.logError(err)
		}
	}
}
