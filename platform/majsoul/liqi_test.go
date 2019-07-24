package majsoul

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
	"github.com/EndlessCheng/mahjong-helper/platform/majsoul/proto/lq"
	"github.com/EndlessCheng/mahjong-helper/tool"
	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
)

func _genReqLogin(t *testing.T) *lq.ReqLogin {
	username, ok := os.LookupEnv("USERNAME")
	if !ok {
		t.Skip("未配置环境变量 USERNAME，退出")
	}

	password, ok := os.LookupEnv("PASSWORD")
	if !ok {
		t.Skip("未配置环境变量 PASSWORD，退出")
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

	version, err := tool.GetMajsoulVersion(tool.ApiGetVersionZH)
	if err != nil {
		t.Fatal(err)
	}
	return &lq.ReqLogin{
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
}

func _genReqOauth2Login(t *testing.T, accessToken string) *lq.ReqOauth2Login {
	randomKey, ok := os.LookupEnv("RANDOM_KEY")
	if !ok {
		rawRandomKey, _ := uuid.NewV4()
		randomKey = rawRandomKey.String()
	}

	version, err := tool.GetMajsoulVersion(tool.ApiGetVersionZH)
	if err != nil {
		t.Fatal(err)
	}
	return &lq.ReqOauth2Login{
		Type:        0, // ? 怀疑是账号/QQ/微信/微博
		AccessToken: accessToken,
		Reconnect:   false,
		Device: &lq.ClientDeviceInfo{
			DeviceType: "pc",
			Os:         "",
			OsVersion:  "",
			Browser:    "safari",
		},
		RandomKey:         randomKey,
		ClientVersion:     version.ResVersion,
		CurrencyPlatforms: []uint32{2}, // 1-inGooglePlay, 2-inChina
	}
}

func TestLogin(t *testing.T) {
	endpoint, err := tool.GetMajsoulWebSocketURL() // wss://mj-srv-7.majsoul.com:4131/
	if err != nil {
		t.Fatal(err)
	}
	t.Log("连接 endpoint: " + endpoint)
	c := NewWebSocketClient()
	if err := c.Connect(endpoint, tool.MajsoulOriginURL); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	reqLogin := _genReqLogin(t)
	respLoginChan := make(chan *lq.ResLogin)
	if err := c.callLobby("login", reqLogin, respLoginChan); err != nil {
		t.Fatal(err)
	}
	respLogin := <-respLoginChan
	if respLogin.GetError() != nil {
		t.Skip("登录失败:", respLogin.Error)
	}
	t.Log("登录成功:", respLogin)
	t.Log(respLogin.GetAccessToken())

	time.Sleep(time.Second)

	reqLogout := lq.ReqLogout{}
	respLogoutChan := make(chan *lq.ResLogout)
	if err := c.callLobby("logout", &reqLogout, respLogoutChan); err != nil {
		t.Fatal(err)
	}
	respLogout := <-respLogoutChan
	t.Log("登出", respLogout)
}

func TestReLogin(t *testing.T) {
	endpoint, err := tool.GetMajsoulWebSocketURL() // wss://mj-srv-7.majsoul.com:4131/
	if err != nil {
		t.Fatal(err)
	}
	t.Log("连接 endpoint: " + endpoint)
	c := NewWebSocketClient()
	if err := c.Connect(endpoint, tool.MajsoulOriginURL); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	accessToken, ok := os.LookupEnv("TOKEN")
	if !ok {
		t.Skip("未配置环境变量 TOKEN，退出")
	}
	reqOauth2Check := lq.ReqOauth2Check{
		Type:        0, // ? 怀疑是账号/QQ/微信/微博
		AccessToken: accessToken,
	}
	respOauth2CheckChan := make(chan *lq.ResOauth2Check)
	if err := c.callLobby("oauth2Check", &reqOauth2Check, respOauth2CheckChan); err != nil {
		t.Fatal(err)
	}
	respOauth2Check := <-respOauth2CheckChan
	if respOauth2Check.GetError() != nil {
		t.Skip("oauth2Check 失败:", respOauth2Check.Error)
	}
	t.Log(respOauth2Check)

	if !respOauth2Check.GetHasAccount() {
		t.Skip("无效的 token")
	}

	reqOauth2Login := _genReqOauth2Login(t, accessToken)
	respLoginChan := make(chan *lq.ResLogin)
	if err := c.callLobby("oauth2Login", reqOauth2Login, respLoginChan); err != nil {
		t.Fatal(err)
	}
	respLogin := <-respLoginChan
	if respLogin.GetError() != nil {
		t.Skip("登录失败:", respLogin.Error)
	}
	t.Log("登录成功:", respLogin)
	t.Log(respLogin.GetAccessToken())

	time.Sleep(time.Second)

	reqLogout := lq.ReqLogout{}
	respLogoutChan := make(chan *lq.ResLogout)
	if err := c.callLobby("logout", &reqLogout, respLogoutChan); err != nil {
		t.Fatal(err)
	}
	respLogout := <-respLogoutChan
	t.Log("登出", respLogout)
}

const (
	recordListTypeAll    = 0
	recordListTypeFriend = 1
	recordListTypeMatch  = 2
	recordListTypeCompte = 3
	recordListTypeFav    = 4
)

func TestFetchGameRecordList(t *testing.T) {
	endpoint, err := tool.GetMajsoulWebSocketURL()
	if err != nil {
		t.Fatal(err)
	}
	c := NewWebSocketClient()
	if err := c.Connect(endpoint, tool.MajsoulOriginURL); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	// 登录
	reqLogin := _genReqLogin(t)
	respLoginChan := make(chan *lq.ResLogin)
	if err := c.callLobby("login", reqLogin, respLoginChan); err != nil {
		t.Fatal(err)
	}
	respLogin := <-respLoginChan
	if respLogin.GetError() != nil {
		t.Skip("登录失败:", respLogin.Error)
	}
	defer func() {
		reqLogout := lq.ReqLogout{}
		respLogoutChan := make(chan *lq.ResLogout)
		if err := c.callLobby("logout", &reqLogout, respLogoutChan); err != nil {
			t.Fatal(err)
		}
		respLogout := <-respLogoutChan
		t.Log("登出", respLogout)
	}()

	// 分页获取牌谱列表
	// TODO: 若牌谱数量巨大，可以使用协程增加下载速度
	reqGameRecordList := lq.ReqGameRecordList{
		Start: 1,
		Count: 10,
		Type:  0, // 全部/友人/段位/比赛/收藏
	}
	respGameRecordListChan := make(chan *lq.ResGameRecordList)
	if err := c.callLobby("fetchGameRecordList", &reqGameRecordList, respGameRecordListChan); err != nil {
		t.Fatal(err)
	}
	respGameRecordList := <-respGameRecordListChan

	for i, gameRecord := range respGameRecordList.GetRecordList() {
		t.Log(i+1, gameRecord.Uuid)

		// 获取具体牌谱内容
		reqGameRecord := lq.ReqGameRecord{
			GameUuid: gameRecord.Uuid,
		}
		respGameRecordChan := make(chan *lq.ResGameRecord)
		if err := c.callLobby("fetchGameRecord", &reqGameRecord, respGameRecordChan); err != nil {
			t.Fatal(err)
		}
		respGameRecord := <-respGameRecordChan

		// 解析
		data := respGameRecord.GetData()
		if len(data) == 0 {
			dataURL := respGameRecord.GetDataUrl()
			if dataURL == "" {
				t.Error("数据异常: dataURL 为空")
				continue
			}
			data, err = tool.Fetch(dataURL)
			if err != nil {
				t.Error(err)
				continue
			}
		}
		detailRecords := lq.GameDetailRecords{}
		if err := c.unwrapMessage(data, &detailRecords); err != nil {
			t.Fatal(err)
		}

		type messageWithType struct {
			Name string        `json:"name"`
			Data proto.Message `json:"data"`
		}
		details := []messageWithType{}
		for _, detailRecord := range detailRecords.GetRecords() {
			name, data, err := c.unwrapData(detailRecord)
			if err != nil {
				t.Fatal(err)
			}

			name = name[1:] // 移除开头的 .
			mt := proto.MessageType(name)
			if mt == nil {
				t.Fatalf("未找到 %s，请检查！", name)
			}
			messagePtr := reflect.New(mt.Elem())
			if err := proto.Unmarshal(data, messagePtr.Interface().(proto.Message)); err != nil {
				t.Fatal(err)
			}

			details = append(details, messageWithType{
				Name: name[3:], // 移除开头的 lq.
				Data: messagePtr.Interface().(proto.Message),
			})
		}

		// 保存至本地（JSON 格式）
		parseResult := struct {
			Head    *lq.RecordGame    `json:"head"`
			Details []messageWithType `json:"details"`
		}{
			Head:    gameRecord,
			Details: details,
		}
		jsonData, err := json.MarshalIndent(&parseResult, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := ioutil.WriteFile(gameRecord.Uuid+".json", jsonData, 0644); err != nil {
			t.Fatal(err)
		}
	}
}
