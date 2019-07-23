package majsoul

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/EndlessCheng/mahjong-helper/platform/majsoul/proto/lq"
	"github.com/EndlessCheng/mahjong-helper/tool"
	"github.com/satori/go.uuid"
	"os"
	"testing"
)

func genLoginReq(t *testing.T) *lq.ReqLogin {
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

func TestLogin(t *testing.T) {
	endpoint, err := tool.GetMajsoulWebSocketURL() // wss://mj-srv-7.majsoul.com:4131/
	if err != nil {
		t.Fatal(err)
	}
	t.Log("endpoint: " + endpoint)
	rpcCh := newRpcChannel()
	if err := rpcCh.connect(endpoint, tool.MajsoulOriginURL); err != nil {
		t.Fatal(err)
	}
	defer rpcCh.close()

	reqLogin := genLoginReq(t)
	respLoginChan := make(chan *lq.ResLogin)
	if err := rpcCh.callLobby("login", reqLogin, respLoginChan); err != nil {
		t.Fatal(err)
	}
	respLogin := <-respLoginChan
	if respLogin.GetAccountId() == 0 {
		t.Skip("登录失败")
	}
	t.Log(respLogin)
	t.Log(respLogin.GetAccessToken())

	reqLogout := lq.ReqLogout{}
	respLogoutChan := make(chan *lq.ResLogout)
	if err := rpcCh.callLobby("logout", &reqLogout, respLogoutChan); err != nil {
		t.Fatal(err)
	}
	respLogout := <-respLogoutChan
	t.Log(respLogout)
}

func TestFetchGameRecordList(t *testing.T) {
	endpoint, err := tool.GetMajsoulWebSocketURL()
	if err != nil {
		t.Fatal(err)
	}
	rpcCh := newRpcChannel()
	if err := rpcCh.connect(endpoint, tool.MajsoulOriginURL); err != nil {
		t.Fatal(err)
	}
	defer rpcCh.close()

	reqLogin := genLoginReq(t)
	respLoginChan := make(chan *lq.ResLogin)
	if err := rpcCh.callLobby("login", reqLogin, respLoginChan); err != nil {
		t.Fatal(err)
	}
	<-respLoginChan

	reqGameRecordList := lq.ReqGameRecordList{
		Start: 1,
		Count: 10,
		Type:  0, // ?
	}
	respGameRecordListChan := make(chan *lq.ResGameRecordList)
	if err := rpcCh.callLobby("fetchGameRecordList", &reqGameRecordList, respGameRecordListChan); err != nil {
		t.Fatal(err)
	}
	respGameRecordList := <-respGameRecordListChan

	//t.Log(respGameRecordList)
	records := respGameRecordList.GetRecordList()
	if len(records) == 0 {
		t.Skip("没有牌谱")
	}
	for i, record := range records {
		t.Log(i+1, record.Uuid)
	}

	reqGameRecord := lq.ReqGameRecord{
		GameUuid: records[0].Uuid,
	}
	respGameRecordChan := make(chan *lq.ResGameRecord)
	if err := rpcCh.callLobby("fetchGameRecord", &reqGameRecord, respGameRecordChan); err != nil {
		t.Fatal(err)
	}
	respGameRecord := <-respGameRecordChan
	//t.Log(respGameRecord)

	data := respGameRecord.GetData()
	if len(data) == 0 {
		// TODO: respGameRecord.DataUrl
		t.Skip(respGameRecord.GetDataUrl())
	}
	detailRecords := lq.GameDetailRecords{}
	if err := rpcCh.unwrapMessage(data, &detailRecords); err != nil {
		t.Fatal(err)
	}

	t.Log(detailRecords)
	for _, detailRecord := range detailRecords.GetRecords() {
		name, data, err := rpcCh.unwrapData(detailRecord)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(name, string(data))
		return
	}
}
