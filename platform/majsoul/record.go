package majsoul

import (
	"github.com/EndlessCheng/mahjong-helper/tool"
	"github.com/EndlessCheng/mahjong-helper/platform/majsoul/proto/lq"
	"github.com/golang/protobuf/proto"
	"reflect"
	"encoding/json"
	"io/ioutil"
	"os"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/satori/go.uuid"
)

const (
	RecordTypeAll   = 0
	RecordTypeMatch = 4
)

func genLoginReq(username string, password string) (*lq.ReqLogin, error) {
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
		return nil, err
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
	}, nil
}

func DownloadRecords(username string, password string, recordType uint32) error {
	endpoint, err := tool.GetMajsoulWebSocketURL()
	if err != nil {
		return err
	}
	rpcCh := newRpcChannel()
	if err := rpcCh.connect(endpoint, tool.MajsoulOriginURL); err != nil {
		return err
	}
	defer rpcCh.close()

	// 登录
	reqLogin, err := genLoginReq(username, password)
	respLoginChan := make(chan *lq.ResLogin)
	if err := rpcCh.callLobby("login", reqLogin, respLoginChan); err != nil {
		return err
	}
	respLogin := <-respLoginChan
	if respLogin.GetError() != nil {
		return fmt.Errorf("登录失败: %v", respLogin.Error)
	}
	defer func() {
		reqLogout := lq.ReqLogout{}
		respLogoutChan := make(chan *lq.ResLogout)
		rpcCh.callLobby("logout", &reqLogout, respLogoutChan)
		<-respLogoutChan
	}()

	// 分页获取牌谱列表
	recordList := []*lq.RecordGame{}
	const pageSize = 10
	for i := uint32(1); ; i += pageSize {
		reqGameRecordList := lq.ReqGameRecordList{
			Start: i,
			Count: pageSize,
			Type:  recordType, // 全部/友人/段位/比赛/收藏
		}
		respGameRecordListChan := make(chan *lq.ResGameRecordList)
		if err := rpcCh.callLobby("fetchGameRecordList", &reqGameRecordList, respGameRecordListChan); err != nil {
			return err
		}
		respGameRecordList := <-respGameRecordListChan
		l := respGameRecordList.GetRecordList()
		if len(l) < pageSize {
			break
		}
		recordList = append(recordList, l...)
	}

	// TODO: 若牌谱数量巨大，可以使用协程增加下载速度
	for i, gameRecord := range recordList {
		fmt.Printf("%d/%d %s\n", i+1, len(recordList), gameRecord.Uuid)

		// 获取具体牌谱内容
		reqGameRecord := lq.ReqGameRecord{
			GameUuid: gameRecord.Uuid,
		}
		respGameRecordChan := make(chan *lq.ResGameRecord)
		if err := rpcCh.callLobby("fetchGameRecord", &reqGameRecord, respGameRecordChan); err != nil {
			return err
		}
		respGameRecord := <-respGameRecordChan

		// 解析
		data := respGameRecord.GetData()
		if len(data) == 0 {
			dataURL := respGameRecord.GetDataUrl()
			if dataURL == "" {
				fmt.Fprintln(os.Stderr, "数据异常: dataURL 为空")
				continue
			}
			data, err = tool.Fetch(dataURL)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
		}
		detailRecords := lq.GameDetailRecords{}
		if err := rpcCh.unwrapMessage(data, &detailRecords); err != nil {
			return err
		}

		type messageWithType struct {
			Name string        `json:"name"`
			Data proto.Message `json:"data"`
		}
		details := []messageWithType{}
		for _, detailRecord := range detailRecords.GetRecords() {
			name, data, err := rpcCh.unwrapData(detailRecord)
			if err != nil {
				return err
			}

			name = name[1:] // 移除开头的 .
			mt := proto.MessageType(name)
			if mt == nil {
				fmt.Fprintf(os.Stderr, "未找到 %s，请检查！", name)
			}
			messagePtr := reflect.New(mt.Elem())
			if err := proto.Unmarshal(data, messagePtr.Interface().(proto.Message)); err != nil {
				return err
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
			return err
		}
		if err := ioutil.WriteFile(gameRecord.Uuid+".json", jsonData, 0644); err != nil {
			return err
		}
	}

	return nil
}
