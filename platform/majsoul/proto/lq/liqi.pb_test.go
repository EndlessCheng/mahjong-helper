package lq

import (
	"testing"
	"os"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/gorilla/websocket"
	"github.com/golang/protobuf/proto"
	"net/http"
	"encoding/binary"
	"github.com/EndlessCheng/mahjong-helper/tool"
)

const (
	messageTypeNotify   = 1
	messageTypeRequest  = 2
	messageTypeResponse = 3
)

func wrapMessage(name string, message proto.Message) (data []byte, err error) {
	data, err = proto.Marshal(message)
	if err != nil {
		return
	}
	wrap := Wrapper{
		Name: name,
		Data: data,
	}
	return proto.Marshal(&wrap)
}

func unwrapData(data []byte, message proto.Message) error {
	msgType := data[0]
	switch msgType {
	case messageTypeNotify:
	case messageTypeRequest:
	case messageTypeResponse:
		// 请求序号
		reqOrder := binary.LittleEndian.Uint16(data[1:3])
		fmt.Println("reqOrder:", reqOrder)

		wrapper := Wrapper{}
		if err := proto.Unmarshal(data[3:], &wrapper); err != nil {
			return err
		}
		if err := proto.Unmarshal(wrapper.Data, message); err != nil {
			return err
		}
	default:
		return fmt.Errorf("[unwrapData] 收到了异常的数据，请检查 %v %s", data, string(data))
	}
	return nil
}

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
	const key = "lailai" // from code.js
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(password))
	password = fmt.Sprintf("%x", mac.Sum(nil))

	// TODO: UUID 最好固定住，生成后保存到本地
	rawRandomKey, err := uuid.NewV4()
	randomKey := rawRandomKey.String()

	// 获取并连接雀魂 WebSocket 服务器
	endPoint, err := tool.GetMajsoulWebSocketURL() // wss://mj-srv-7.majsoul.com:4131/
	if err != nil {
		t.Fatal(err)
	}
	header := http.Header{}
	header.Set("origin", tool.MajsoulOriginURL) // 模拟来源
	ws, _, err := websocket.DefaultDialer.Dial(endPoint, header)
	if err != nil {
		t.Fatal(err)
	}
	defer ws.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)

		// 读取雀魂 WebSocket 服务器返回的消息
		_, message, err := ws.ReadMessage()
		if err != nil {
			t.Fatal(err)
		}

		// 解析登录请求结果
		respLogin := ResLogin{}
		if err := unwrapData(message, &respLogin); err != nil {
			t.Fatal(err)
		}
		t.Log(respLogin)
	}()

	// 生成登录请求参数
	version, err := tool.GetMajsoulVersion(tool.ApiGetVersionZH)
	if err != nil {
		t.Fatal(err)
	}
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
		RandomKey:         randomKey,          // 例如 e6b3fafb-aa11-11e9-8323-f45c89a43cff
		ClientVersion:     version.ResVersion, // 0.5.162.w
		GenAccessToken:    true,
		CurrencyPlatforms: []uint32{2}, // 1-inGooglePlay, 2-inChina
	}
	data, err := wrapMessage(".lq.Lobby.login", &reqLogin)
	if err != nil {
		t.Fatal(err)
	}

	// 填写消息序号后，发送登录请求给雀魂 WebSocket 服务器
	reqOrderBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(reqOrderBytes, 1) // [0,60007) 的一个数，from code.js
	msgHead := append([]byte{messageTypeRequest}, reqOrderBytes...)
	if err := ws.WriteMessage(websocket.BinaryMessage, append(msgHead, data...)); err != nil {
		t.Fatal(err)
	}

	<-done
}
