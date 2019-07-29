package lq

import (
	"fmt"
	"time"
	"reflect"
	"strings"
	"github.com/golang/protobuf/proto"
	"sort"
	"github.com/EndlessCheng/mahjong-helper/util"
	"strconv"
)

var (
	lobbyClientMethodMap    = map[string]reflect.Type{}
	fastTestClientMethodMap = map[string]reflect.Type{}
)

func init() {
	t := reflect.TypeOf((*LobbyClient)(nil)).Elem()
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		lobbyClientMethodMap[method.Name] = method.Type
	}

	t = reflect.TypeOf((*FastTestClient)(nil)).Elem()
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		fastTestClientMethodMap[method.Name] = method.Type
	}
}

func FindMethod(clientName string, methodName string) reflect.Type {
	methodName = strings.Title(methodName)
	if clientName == "Lobby" {
		return lobbyClientMethodMap[methodName]
	} else { // clientName == "FastTest"
		return fastTestClientMethodMap[methodName]
	}
}

// 下面补充一些功能

func (m *Friend) CLIString() string {
	return fmt.Sprintf("%9d   %s   %s   %s",
		m.Base.AccountId,
		time.Unix(int64(m.State.LoginTime), 0).Format("2006-01-02 15:04:05"),
		time.Unix(int64(m.State.LogoutTime), 0).Format("2006-01-02 15:04:05"),
		m.Base.Nickname,
	)
}

type FriendList []*Friend

func (l FriendList) String() string {
	out := "好友账号ID   好友上次登录时间        好友上次登出时间       好友昵称\n"
	for _, friend := range l {
		out += friend.CLIString() + "\n"
	}
	return out
}

func (m *ActionPrototype) ParseData() (proto.Message, error) {
	name := "lq." + m.Name
	mt := proto.MessageType(name)
	if mt == nil {
		return nil, fmt.Errorf("ActionPrototype.ParseData 未找到 %s，请检查！", name)
	}
	messagePtr := reflect.New(mt.Elem())
	if err := proto.Unmarshal(m.Data, messagePtr.Interface().(proto.Message)); err != nil {
		return nil, err
	}
	return messagePtr.Interface().(proto.Message), nil
}

func (m *GameConfig) IsGuyiMode() bool {
	return m != nil && m.Mode != nil && m.Mode.DetailRule != nil && m.Mode.DetailRule.GuyiMode == 1
}

func (m *RecordGame) GetSelfSeat(accountID int) (int, error) {
	if len(m.Accounts) == 0 {
		return -1, fmt.Errorf("牌谱基本信息为空")
	}
	for _, account := range m.Accounts {
		if int(account.AccountId) == accountID {
			return int(account.Seat), nil
		}
	}
	// 若没有，则以东家为主视角
	return 0, nil
}

func (m *RecordGame) sort() {
	sort.Slice(m.Accounts, func(i_, j int) bool {
		return m.Accounts[i_].Seat < m.Accounts[j].Seat
	})
}

var seatNameZH = []string{"东", "南", "西", "北"}

func (m *RecordGame) CLIString() string {
	m.sort()

	const timeFormat = "2006-01-02 15:04:05"
	output := fmt.Sprintf("%s\n从 %s\n到 %s\n\n", m.Uuid, time.Unix(int64(m.StartTime), 0).Format(timeFormat), time.Unix(int64(m.EndTime), 0).Format(timeFormat))

	maxAccountID := 0
	for _, account := range m.Accounts {
		maxAccountID = util.MaxInt(maxAccountID, int(account.AccountId))
	}
	accountShownWidth := len(strconv.Itoa(maxAccountID))
	for _, account := range m.Accounts {
		output += fmt.Sprintf("%s %*d %s\n", seatNameZH[account.Seat], accountShownWidth, account.AccountId, account.Nickname)
	}
	return output
}
