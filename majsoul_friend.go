package main

import (
	"fmt"
	"time"
)

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
