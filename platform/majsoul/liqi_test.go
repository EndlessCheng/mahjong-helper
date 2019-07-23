package majsoul

import (
	"testing"
	"os"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/EndlessCheng/mahjong-helper/tool"
	"github.com/EndlessCheng/mahjong-helper/platform/majsoul/proto/lq"
)

func TestLogin(t *testing.T) {
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

	// randomKey 最好是个固定值
	randomKey, ok := os.LookupEnv("RANDOM_KEY")
	if !ok {
		rawRandomKey, _ := uuid.NewV4()
		randomKey = rawRandomKey.String()
	}

	// 获取并连接雀魂 WebSocket 服务器
	endpoint, err := tool.GetMajsoulWebSocketURL() // wss://mj-srv-7.majsoul.com:4131/
	if err != nil {
		t.Fatal(err)
	}
	rpcCh := newRpcChannel()
	if err := rpcCh.connect(endpoint, tool.MajsoulOriginURL); err != nil {
		t.Fatal(err)
	}
	defer rpcCh.close()

	// 生成登录请求参数
	version, err := tool.GetMajsoulVersion(tool.ApiGetVersionZH)
	if err != nil {
		t.Fatal(err)
	}
	reqLogin := lq.ReqLogin{
		Account:   username,
		Password:  password,
		Reconnect: false,
		Device: &lq.ClientDeviceInfo{
			DeviceType: "pc",
			Os:         "",
			OsVersion:  "",
			Browser:    "safari",
		},
		RandomKey:         randomKey,          // 例如 aa566cfc-547e-4cc0-a36f-2ebe6269109b
		ClientVersion:     version.ResVersion, // 0.5.162.w
		GenAccessToken:    true,
		CurrencyPlatforms: []uint32{2}, // 1-inGooglePlay, 2-inChina
	}
	respLoginChan := make(chan *lq.ResLogin)
	if err := rpcCh.send(".lq.Lobby.login", &reqLogin, respLoginChan); err != nil {
		t.Fatal(err)
	}
	respLogin := <-respLoginChan
	t.Log(respLogin)
}
