package majsoul

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"
	"github.com/EndlessCheng/mahjong-helper/platform/majsoul/proto/lq"
	"github.com/EndlessCheng/mahjong-helper/tool"
)

const (
	RecordTypeAll   = 0
	RecordTypeMatch = 4
)

func genReqLogin(username string, password string) (*lq.ReqLogin, error) {
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

// TODO: add token
func DownloadRecords(username string, password string, recordType uint32) error {
	endpoint, err := tool.GetMajsoulWebSocketURL()
	if err != nil {
		return err
	}
	c := NewWebSocketClient()
	if err := c.Connect(endpoint, tool.MajsoulOriginURL); err != nil {
		return err
	}
	defer c.Close()

	// 登录
	reqLogin, err := genReqLogin(username, password)
	if err != nil {
		return err
	}
	if _, err := c.Login(reqLogin); err != nil {
		return err
	}
	defer c.Logout(&lq.ReqLogout{})

	// 分页获取牌谱列表
	recordList := []*lq.RecordGame{}
	const pageSize = 10
	for i := uint32(1); ; i += pageSize {
		reqGameRecordList := lq.ReqGameRecordList{
			Start: i,
			Count: pageSize,
			Type:  recordType, // 全部/友人/段位/比赛/收藏
		}
		respGameRecordList, err := c.FetchGameRecordList(&reqGameRecordList)
		if err != nil {
			return err
		}
		recordList = append(recordList, respGameRecordList.RecordList...)
		if len(respGameRecordList.RecordList) < pageSize {
			break
		}
	}

	// TODO: 若牌谱数量巨大，可以使用协程增加下载速度
	for i, gameRecord := range recordList {
		fmt.Printf("%d/%d %s\n", i+1, len(recordList), gameRecord.Uuid)

		// 获取具体牌谱内容
		reqGameRecord := lq.ReqGameRecord{
			GameUuid: gameRecord.Uuid,
		}
		respGameRecord, err := c.FetchGameRecord(&reqGameRecord)
		if err != nil {
			return err
		}

		// 解析
		data := respGameRecord.Data
		if len(data) == 0 {
			dataURL := respGameRecord.DataUrl
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
		if err := c.UnwrapMessage(data, &detailRecords); err != nil {
			return err
		}

		type messageWithType struct {
			Name string        `json:"name"`
			Data proto.Message `json:"data"`
		}
		details := []messageWithType{}
		for _, detailRecord := range detailRecords.GetRecords() {
			name, data, err := c.UnwrapData(detailRecord)
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
