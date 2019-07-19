package lq

import (
	"testing"
	"os"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/satori/go.uuid"
	//"golang.org/x/net/websocket"
	"github.com/gorilla/websocket"
	"github.com/golang/protobuf/proto"
	"net/http"
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
	const key = "lailai" // 提取于 code.js 源码
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(password))
	password = fmt.Sprintf("%x", mac.Sum(nil))

	// UUID 最好固定住，生成后保存到本地
	rawRandomKey, err := uuid.NewV4()
	randomKey := rawRandomKey.String()

	const resVersion = "v0.5.163.w"

	const endPoint = "wss://mj-srv-7.majsoul.com:4131/"
	const originZH = "https://majsoul.union-game.com" // 模拟来源
	//ws, err := websocket.Dial(endPoint, "", originZH)
	header := http.Header{}
	header.Set("originZH", originZH)
	ws, _, err := websocket.DefaultDialer.Dial(endPoint, header)
	if err != nil {
		t.Fatal(err)
	}
	defer ws.Close()

	done := make(chan bool)
	go func() {
		//var msg string
		_, message, err := ws.ReadMessage()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(message)
		fmt.Println(string(message))
		done <- true
	}()

	reqLogin := ReqLogin{
		Account:   username,
		Password:  password,
		Reconnect: false,
		Device: &ClientDeviceInfo{
			DeviceType: "pc",
			Os:         "",
			OsVersion:  "",
			Browser:    "safari",
		},
		RandomKey:         randomKey,
		ClientVersion:     resVersion,
		GenAccessToken:    true,
		CurrencyPlatforms: []uint32{2}, // 1-inGooglePlay, 2-inChina
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
	err = ws.WriteMessage(websocket.BinaryMessage, append(msgHead, data...))
	if err != nil {
		t.Fatal(err)
	}
	//t.Logf("%#v", n)
	<-done
}
