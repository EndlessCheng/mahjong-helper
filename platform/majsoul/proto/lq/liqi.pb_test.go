package lq

import (
	"testing"
	"os"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
	"github.com/golang/protobuf/proto"
)

func TestReqLogin(t *testing.T) {
	username, ok := os.LookupEnv("USERNAME")
	if !ok {
		t.Log("未配置环境变量 USERNAME，退出")
		t.Skip()
	}
	password, ok := os.LookupEnv("PASSWORD")
	if !ok {
		t.Log("未配置环境变量 PASSWORD，退出")
		t.Skip()
	}
	const key = "lailai"
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(password))
	password = fmt.Sprintf("%x", mac.Sum(nil))

	rawRandomKey, err := uuid.NewV1()
	randomKey := rawRandomKey.String()

	const clientVersion = "v0.5.43.w"

	const endPoint = "wss://mj-srv-7.majsoul.com:4131/"
	const origin = "https://majsoul.union-game.com"
	ws, err := websocket.Dial(endPoint, "", origin)
	if err != nil {
		t.Fatal(err)
	}
	defer ws.Close()

	done := make(chan bool)
	go func() {
		var msg string
		if err := websocket.Message.Receive(ws, &msg); err != nil {
			t.Fatal(err)
		}
		done <- true
	}()

	reqLogin := ReqLogin{
		Account:   username,
		Password:  password,
		Reconnect: true,
		Device: &ClientDeviceInfo{
			DeviceType: "pc",
			Os:         "",
			OsVersion:  "",
			Browser:    "safari",
		},
		RandomKey:         randomKey,
		ClientVersion:     clientVersion,
		GenAccessToken:    true,
		CurrencyPlatforms: []uint32{2},
	}
	data, _ := proto.Marshal(&reqLogin)
	fmt.Println(string(data))

	wrap := Wrapper{
		Name: ".lq.Lobby.login",
		Data: data,
	}
	data, _ = proto.Marshal(&wrap)
	fmt.Println(string(data))

	msgHead := []byte{0x02, 0x01, 0x00}
	n, err := ws.Write(append(msgHead, data...))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", n)
	<-done
}
