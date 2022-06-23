package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	stdLog "log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/EndlessCheng/mahjong-helper/Console"
	"github.com/EndlessCheng/mahjong-helper/platform/tenhou"
	"github.com/EndlessCheng/mahjong-helper/util"
	"github.com/EndlessCheng/mahjong-helper/util/debug"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const DefaultPort = 12121

func NewLogFilePath() (filePath string, err error) {
	const logDir = "log"
	if err = os.MkdirAll(logDir, os.ModePerm); err != nil {
		return
	}
	fileName := fmt.Sprintf("gamedata-%s.log", time.Now().Format("20060102-150405"))
	filePath = filepath.Join(logDir, fileName)
	return filepath.Abs(filePath)
}

type MahJongHandler struct {
	Log echo.Logger

	Analysing bool

	TenhouMessageReceiver *tenhou.MessageReceiver
	TenHouRoundData       *TenHouRoundData

	majsoulMessageQueue  chan []byte
	MahJongSoulRoundData *MahJongSoulRoundData

	MahJongSoulRecordMap                map[string]*MahJongSoulRecordBaseInfo
	MahJongSoulCurrentRecordUUID        string
	MahJongSoulCurrentRecordActionsList []MahJongSoulRoundActions
	MahJongSoulCurrentRoundIndex        int
	MahJongSoulCurrentActionIndex       int

	MahJongSoulCurrentRoundActions MahJongSoulRoundActions
}

// write error to log
func (handler *MahJongHandler) LogError(err error) {
	fmt.Fprintln(os.Stderr, err)
	if !DebugMode {
		handler.Log.Error(err)
	}
}

// 调试用
func (handler *MahJongHandler) Index(echo_context echo.Context) error {
	data, err := ioutil.ReadAll(echo_context.Request().Body)
	if err != nil {
		handler.Log.Error("[MahJongHandler.index.ioutil.ReadAll]", err)
		return echo_context.NoContent(http.StatusInternalServerError)
	}

	fmt.Println(data, string(data))
	handler.Log.Info(data)
	return echo_context.String(http.StatusOK, time.Now().Format("2006-01-02 15:04:05"))
}

// 打一摸一分析器
func (handler *MahJongHandler) Analysis(echo_context echo.Context) error {
	if handler.Analysing {
		return echo_context.NoContent(http.StatusForbidden)
	}

	handler.Analysing = true
	defer func() { handler.Analysing = false }()

	data := struct {
		Reset bool   `json:"reset"`
		Tiles string `json:"tiles"`
	}{}
	if err := echo_context.Bind(&data); err != nil {
		fmt.Println(err)
		return echo_context.String(http.StatusBadRequest, err.Error())
	}

	if _, err := AnalysisHumanTiles(model.NewSimpleHumanTilesInfo(data.Tiles)); err != nil {
		fmt.Println(err)
		return echo_context.String(http.StatusBadRequest, err.Error())
	}

	return echo_context.NoContent(http.StatusOK)
}

// 分析天凤 WebSocket 数据
func (handler *MahJongHandler) AnalysisTenHou(echo_context echo.Context) error {
	data, err := ioutil.ReadAll(echo_context.Request().Body)
	if err != nil {
		handler.LogError(err)
		return echo_context.String(http.StatusBadRequest, err.Error())
	}

	handler.TenhouMessageReceiver.Put(data)
	return echo_context.NoContent(http.StatusOK)
}

// run analysis TenHou Message
func (handler *MahJongHandler) RunAnalysisTenHouMessageTask() {
	if !DebugMode {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("内部错误：", err)
			}
		}()
	}

	for {
		msg := handler.TenhouMessageReceiver.Get()
		data := TenhouMessage{}
		if err := json.Unmarshal(msg, &data); err != nil {
			handler.LogError(err)
			continue
		}

		originJSON := string(msg)
		if handler.Log != nil {
			handler.Log.Info(originJSON)
		}

		handler.TenHouRoundData.Msg = &data
		handler.TenHouRoundData.OriginJSON = originJSON
		if err := handler.TenHouRoundData.Analysis(); err != nil {
			handler.LogError(err)
		}
	}
}

// 分析雀魂 WebSocket 数据
func (handler *MahJongHandler) AnalysisMajsoul(echo_context echo.Context) error {
	data, err := ioutil.ReadAll(echo_context.Request().Body)
	if err != nil {
		handler.LogError(err)
		return echo_context.String(http.StatusBadRequest, err.Error())
	}

	handler.majsoulMessageQueue <- data
	return echo_context.NoContent(http.StatusOK)
}

// run analysis MahJongSoul Message
func (handler *MahJongHandler) RunAnalysisMahJongSoulMessageTask() {
	if !DebugMode {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("内部错误：", err)
			}
		}()
	}

	for msg := range handler.majsoulMessageQueue {
		d := &MaJongSoulMessage{}
		if err := json.Unmarshal(msg, d); err != nil {
			handler.LogError(err)
			continue
		}

		originJSON := string(msg)
		if handler.Log != nil && debug.Lo == 0 {
			handler.Log.Info(originJSON)
		} else {
			if len(originJSON) > 500 {
				originJSON = originJSON[:500]
			}
			fmt.Println(originJSON)
		}

		switch {
		case len(d.Friends) > 0:
			// 好友列表
			fmt.Println(d.Friends)
		case len(d.RecordBaseInfoList) > 0:
			// 牌谱基本信息列表
			for _, record := range d.RecordBaseInfoList {
				handler.MahJongSoulRecordMap[record.UUID] = record
			}
			color.HiGreen("收到 %2d 个雀魂牌谱（已收集 %d 个），请在网页上点击「查看」", len(d.RecordBaseInfoList), len(handler.MahJongSoulRecordMap))
		case d.SharedRecordBaseInfo != nil:
			// 处理分享的牌谱基本信息
			// FIXME: 观看自己的牌谱也会有 d.SharedRecordBaseInfo
			record := d.SharedRecordBaseInfo
			handler.MahJongSoulRecordMap[record.UUID] = record
			if err := handler.loadMahJongSoulRecordBaseInfo(record.UUID); err != nil {
				handler.LogError(err)
				break
			}
		case d.CurrentRecordUUID != "":
			// 载入某个牌谱
			ResetAnalysisCache()
			handler.MahJongSoulCurrentRecordActionsList = nil

			if err := handler.loadMahJongSoulRecordBaseInfo(d.CurrentRecordUUID); err != nil {
				// 看的是分享的牌谱（先收到 CurrentRecordUUID 和 AccountID，然后收到 SharedRecordBaseInfo）
				// 或者是比赛场的牌谱
				// 记录主视角 ID（可能是 0）
				gameConf.setMajsoulAccountID(d.AccountID)
				break
			}

			// 看的是自己的牌谱
			// 更新当前使用的账号
			gameConf.addMajsoulAccountID(d.AccountID)
			if gameConf.currentActiveMajsoulAccountID != d.AccountID {
				fmt.Println()
				printAccountInfo(d.AccountID)
				gameConf.setMajsoulAccountID(d.AccountID)
			}
		case len(d.RecordActions) > 0:
			if handler.MahJongSoulCurrentRecordActionsList != nil {
				// TODO: 网页发送更恰当的信息？
				break
			}

			if handler.MahJongSoulCurrentRecordUUID == "" {
				handler.LogError(fmt.Errorf("错误：程序未收到所观看的雀魂牌谱的 UUID"))
				break
			}

			baseInfo, ok := handler.MahJongSoulRecordMap[handler.MahJongSoulCurrentRecordUUID]
			if !ok {
				handler.LogError(fmt.Errorf("错误：找不到雀魂牌谱 %s", handler.MahJongSoulCurrentRecordUUID))
				break
			}

			selfAccountID := gameConf.currentActiveMajsoulAccountID
			if selfAccountID == -1 {
				handler.LogError(fmt.Errorf("错误：当前雀魂账号为空"))
				break
			}

			handler.MahJongSoulRoundData.newGame()
			handler.MahJongSoulRoundData.GameMode = gameModeRecord

			// 获取并设置主视角初始座位
			selfSeat, err := baseInfo.getSelfSeat(selfAccountID)
			if err != nil {
				handler.LogError(err)
				break
			}
			handler.MahJongSoulRoundData.SelfSeat = selfSeat

			// 准备分析……
			majsoulCurrentRecordActions, err := parseMajsoulRecordAction(d.RecordActions)
			if err != nil {
				handler.LogError(err)
				break
			}
			handler.MahJongSoulCurrentRecordActionsList = majsoulCurrentRecordActions
			handler.MahJongSoulCurrentRoundIndex = 0
			handler.MahJongSoulCurrentActionIndex = 0

			actions := handler.MahJongSoulCurrentRecordActionsList[handler.MahJongSoulCurrentRoundIndex]

			// 创建分析任务
			analysisCache := newGameAnalysisCache(handler.MahJongSoulCurrentRecordUUID, selfSeat)
			SetAnalysisCache(analysisCache)
			go analysisCache.RunMahJongSoulRecordAnalysisTask(actions)

			// 分析第一局的起始信息
			data := actions[0].Action
			handler.analysisMahJongSoulRoundData(data, originJSON)
		case d.RecordClickAction != "":
			// 处理网页上的牌谱点击：上一局/跳到某局/下一局/上一巡/跳到某巡/下一巡/上一步/播放/暂停/下一步/点击桌面
			// 暂不能分析他家手牌
			handler.onRecordClick(d.RecordClickAction, d.RecordClickActionIndex, d.FastRecordTo)
		case d.LiveBaseInfo != nil:
			// 观战
			gameConf.setMajsoulAccountID(1) // TODO: 重构
			handler.MahJongSoulRoundData.newGame()
			handler.MahJongSoulRoundData.SelfSeat = 0 // 观战进来后看的是东起的玩家
			handler.MahJongSoulRoundData.GameMode = gameModeLive
			Console.ClearScreen()
			fmt.Printf("正在載入對戰：%s", d.LiveBaseInfo.String())
		case d.LiveFastAction != nil:
			if err := handler.loadLiveAction(d.LiveFastAction, true); err != nil {
				handler.LogError(err)
				break
			}
		case d.LiveAction != nil:
			if err := handler.loadLiveAction(d.LiveAction, false); err != nil {
				handler.LogError(err)
				break
			}
		case d.ChangeSeatTo != nil:
			// 切换座位
			changeSeatTo := *(d.ChangeSeatTo)
			handler.MahJongSoulRoundData.SelfSeat = changeSeatTo
			if DebugMode {
				fmt.Println("座位已切换至", changeSeatTo)
			}

			var actions MahJongSoulRoundActions
			if handler.MahJongSoulRoundData.GameMode == gameModeLive { // 观战
				actions = handler.MahJongSoulCurrentRoundActions
			} else { // 牌谱
				fullActions := handler.MahJongSoulCurrentRecordActionsList[handler.MahJongSoulCurrentRoundIndex]
				actions = fullActions[:handler.MahJongSoulCurrentActionIndex+1]
				analysisCache := GetAnalysisCache(changeSeatTo)
				if analysisCache == nil {
					analysisCache = newGameAnalysisCache(handler.MahJongSoulCurrentRecordUUID, changeSeatTo)
				}
				SetAnalysisCache(analysisCache)
				// 创建分析任务
				go analysisCache.RunMahJongSoulRecordAnalysisTask(fullActions)
			}

			handler.fastLoadActions(actions)
		case len(d.SyncGameActions) > 0:
			handler.fastLoadActions(d.SyncGameActions)
		default:
			// 其他：AI 分析
			handler.analysisMahJongSoulRoundData(d, originJSON)
		}
	}
}

// private function to load MahJongSoul record base information
func (handler *MahJongHandler) loadMahJongSoulRecordBaseInfo(mahjongsoulRecordUUID string) error {
	baseInfo, ok := handler.MahJongSoulRecordMap[mahjongsoulRecordUUID]
	if !ok {
		return fmt.Errorf("错误：找不到雀魂牌谱 %s", mahjongsoulRecordUUID)
	}

	// 标记当前正在观看的牌谱
	handler.MahJongSoulCurrentRecordUUID = mahjongsoulRecordUUID
	Console.ClearScreen()
	fmt.Printf("正在解析雀魂牌谱：%s", baseInfo.String())

	// 标记古役模式
	isOldYaKuMode := baseInfo.Config.isGuyiMode()
	util.SetConsiderOldYaku(isOldYaKuMode)
	if isOldYaKuMode {
		fmt.Println()
		color.HiGreen("古役模式已開啟")
	}

	return nil
}

// private function to load live action
func (handler *MahJongHandler) loadLiveAction(action *MahJongSoulRecordAction, isFast bool) error {
	if DebugMode {
		fmt.Println("[_loadLiveAction] 收到", action, isFast)
	}

	newActions, err := handler.MahJongSoulCurrentRoundActions.Append(action)
	if err != nil {
		return err
	}
	handler.MahJongSoulCurrentRoundActions = newActions

	handler.MahJongSoulRoundData.SkipOutput = isFast
	handler.analysisMahJongSoulRoundData(action.Action, "")
	return nil
}

// private function analysis MahJongSoul Round Data
func (handler *MahJongHandler) analysisMahJongSoulRoundData(data *MaJongSoulMessage, originJSON string) {
	//if originJSON == "{}" {
	//	return
	//}
	handler.MahJongSoulRoundData.Message = data
	handler.MahJongSoulRoundData.OriginJSON = originJSON
	if err := handler.MahJongSoulRoundData.Analysis(); err != nil {
		handler.LogError(err)
	}
}

// private function fast load actions
func (handler *MahJongHandler) fastLoadActions(actions []*MahJongSoulRecordAction) {
	if len(actions) == 0 {
		return
	}
	fastRecordEnd := util.MaxInt(0, len(actions)-3)
	handler.MahJongSoulRoundData.SkipOutput = true
	// 留最后三个刷新，这样确保会刷新界面
	for _, action := range actions[:fastRecordEnd] {
		handler.analysisMahJongSoulRoundData(action.Action, "")
	}
	handler.MahJongSoulRoundData.SkipOutput = false
	for _, action := range actions[fastRecordEnd:] {
		handler.analysisMahJongSoulRoundData(action.Action, "")
	}
}

// private function on Record Click
func (handler *MahJongHandler) onRecordClick(clickAction string, clickActionIndex int, fastRecordTo int) {
	if DebugMode {
		fmt.Println("[_onRecordClick] 收到", clickAction, clickActionIndex, fastRecordTo)
	}

	analysisCache := GetCurrentAnalysisCache()

	switch clickAction {
	case "nextStep", "update":
		newActionIndex := handler.MahJongSoulCurrentActionIndex + 1
		if newActionIndex >= len(handler.MahJongSoulCurrentRecordActionsList[handler.MahJongSoulCurrentRoundIndex]) {
			return
		}
		handler.MahJongSoulCurrentActionIndex = newActionIndex
	case "nextRound":
		handler.MahJongSoulCurrentRoundIndex = (handler.MahJongSoulCurrentRoundIndex + 1) %
			len(handler.MahJongSoulCurrentRecordActionsList)
		handler.MahJongSoulCurrentActionIndex = 0
		go analysisCache.RunMahJongSoulRecordAnalysisTask(
			handler.MahJongSoulCurrentRecordActionsList[handler.MahJongSoulCurrentRoundIndex])
	case "preRound":
		handler.MahJongSoulCurrentRoundIndex = (handler.MahJongSoulCurrentRoundIndex - 1 +
			len(handler.MahJongSoulCurrentRecordActionsList)) % len(handler.MahJongSoulCurrentRecordActionsList)
		handler.MahJongSoulCurrentActionIndex = 0
		go analysisCache.RunMahJongSoulRecordAnalysisTask(
			handler.MahJongSoulCurrentRecordActionsList[handler.MahJongSoulCurrentRoundIndex])
	case "jumpRound":
		handler.MahJongSoulCurrentRoundIndex = clickActionIndex % len(handler.MahJongSoulCurrentRecordActionsList)
		handler.MahJongSoulCurrentActionIndex = 0
		go analysisCache.RunMahJongSoulRecordAnalysisTask(
			handler.MahJongSoulCurrentRecordActionsList[handler.MahJongSoulCurrentRoundIndex])
	case "nextXun", "preXun", "jumpXun", "preStep", "jumpToLastRoundXun":
		if clickAction == "jumpToLastRoundXun" {
			handler.MahJongSoulCurrentRoundIndex = (handler.MahJongSoulCurrentRoundIndex - 1 + len(
				handler.MahJongSoulCurrentRecordActionsList)) % len(handler.MahJongSoulCurrentRecordActionsList)
			go analysisCache.RunMahJongSoulRecordAnalysisTask(
				handler.MahJongSoulCurrentRecordActionsList[handler.MahJongSoulCurrentRoundIndex])
		}

		handler.MahJongSoulRoundData.SkipOutput = true
		currentRoundActions := handler.MahJongSoulCurrentRecordActionsList[handler.MahJongSoulCurrentRoundIndex]
		startActionIndex := 0
		endActionIndex := fastRecordTo
		if clickAction == "nextXun" {
			startActionIndex = handler.MahJongSoulCurrentActionIndex + 1
		}
		if DebugMode {
			fmt.Printf("快速处理牌谱中的操作：局 %d 动作 %d-%d\n", handler.MahJongSoulCurrentRoundIndex,
				startActionIndex, endActionIndex)
		}
		for i, action := range currentRoundActions[startActionIndex : endActionIndex+1] {
			if DebugMode {
				fmt.Printf("快速处理牌谱中的操作：局 %d 动作 %d\n", handler.MahJongSoulCurrentRoundIndex,
					startActionIndex+i)
			}
			handler.analysisMahJongSoulRoundData(action.Action, "")
		}
		handler.MahJongSoulRoundData.SkipOutput = false

		handler.MahJongSoulCurrentActionIndex = endActionIndex + 1
	default:
		return
	}

	if DebugMode {
		fmt.Printf("处理牌谱中的操作：局 %d 动作 %d\n", handler.MahJongSoulCurrentRoundIndex, handler.MahJongSoulCurrentActionIndex)
	}
	action := handler.MahJongSoulCurrentRecordActionsList[handler.MahJongSoulCurrentRoundIndex][handler.MahJongSoulCurrentActionIndex]
	handler.analysisMahJongSoulRoundData(action.Action, "")

	if action.Name == "RecordHule" || action.Name == "RecordLiuJu" || action.Name == "RecordNoTile" {
		// 播放和牌/流局动画，进入下一局或显示终局动画
		handler.MahJongSoulCurrentRoundIndex++
		handler.MahJongSoulCurrentActionIndex = 0
		if handler.MahJongSoulCurrentRoundIndex == len(handler.MahJongSoulCurrentRecordActionsList) {
			handler.MahJongSoulCurrentRoundIndex = 0
			return
		}

		time.Sleep(time.Second)

		actions := handler.MahJongSoulCurrentRecordActionsList[handler.MahJongSoulCurrentRoundIndex]
		go analysisCache.RunMahJongSoulRecordAnalysisTask(actions)
		// 分析下一局的起始信息
		data := actions[handler.MahJongSoulCurrentActionIndex].Action
		handler.analysisMahJongSoulRoundData(data, "")
	}
}

var handler *MahJongHandler

func GetMajsoulCurrentRecordUUID() string {
	return handler.MahJongSoulCurrentRecordUUID
}

func RunServer(isHTTPS bool, port int) (err error) {
	echo := echo.New()

	// 移除 echo.Echo 和 http.Server 在控制台上打印的信息
	echo.HideBanner = true
	echo.HidePort = true
	echo.StdLogger = stdLog.New(ioutil.Discard, "", 0)

	// 默认是 log.ERROR
	echo.Logger.SetLevel(log.INFO)

	// 设置日志输出到 log/gamedata-xxx.log
	filePath, err := NewLogFilePath()
	if err != nil {
		return
	}
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	echo.Logger.SetOutput(logFile)

	echo.Logger.Info("============================================================================================")
	echo.Logger.Info("服务启动")

	handler = &MahJongHandler{
		Log: echo.Logger,

		TenhouMessageReceiver: tenhou.NewMessageReceiver(),
		TenHouRoundData:       &TenHouRoundData{IsRoundEnd: true},
		majsoulMessageQueue:   make(chan []byte, 100),
		MahJongSoulRoundData:  &MahJongSoulRoundData{SelfSeat: -1},
		MahJongSoulRecordMap:  map[string]*MahJongSoulRecordBaseInfo{},
	}
	handler.TenHouRoundData.RoundData = NewGame(handler.TenHouRoundData)
	handler.MahJongSoulRoundData.RoundData = NewGame(handler.MahJongSoulRoundData)

	go handler.RunAnalysisTenHouMessageTask()
	go handler.RunAnalysisMahJongSoulMessageTask()

	echo.Use(middleware.Recover())
	echo.Use(middleware.CORS())
	echo.GET("/", handler.Index)
	echo.POST("/debug", handler.Index)
	echo.POST("/analysis", handler.Analysis)
	echo.POST("/tenhou", handler.AnalysisTenHou)
	echo.POST("/majsoul", handler.AnalysisMajsoul)

	// code.js 也用的该端口
	if port == 0 {
		port = DefaultPort
	}
	addr := ":" + strconv.Itoa(port)
	if !isHTTPS {
		echo.POST("/", handler.AnalysisTenHou)
		err = echo.Start(addr)
	} else {
		echo.POST("/", handler.AnalysisMajsoul)
		err = StartTLS(echo, addr)
	}
	if err != nil {
		// 检查是否为端口占用错误
		if opErr, ok := err.(*net.OpError); ok && opErr.Op == "listen" {
			if syscallErr, ok := opErr.Err.(*os.SyscallError); ok && syscallErr.Syscall == "bind" {
				color.HiRed(addr + " 端口已被占用，程序无法启动（是否已经开启了本程序？）")
			}
		}
		return
	}
	return nil
}

const (
	certText = `-----BEGIN CERTIFICATE-----
MIIDHjCCAgYCCQDU2jXI1a7kizANBgkqhkiG9w0BAQsFADBRMQswCQYDVQQGEwJV
UzELMAkGA1UECAwCVVMxCzAJBgNVBAcMAkFBMQswCQYDVQQKDAJBQTEMMAoGA1UE
CwwDQUFBMQ0wCwYDVQQDDARBQUFBMB4XDTE5MDIyNjA2Mjc1OFoXDTIwMDIyNjA2
Mjc1OFowUTELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAlVTMQswCQYDVQQHDAJBQTEL
MAkGA1UECgwCQUExDDAKBgNVBAsMA0FBQTENMAsGA1UEAwwEQUFBQTCCASIwDQYJ
KoZIhvcNAQEBBQADggEPADCCAQoCggEBALHryqHQDhOjwfEhzAm7sfiMbFjLAY13
+oyQ+7dTFVe9h2ONYVQ3wvd0f/ncYrUc98n6K+X9c06/auHs0D/ruZa+XizSKyvB
/2vhmbus8mcm8NKZBC2JEi5YI4oIoD8af9kA+cnQ1diwWl60ic54HxSlLpC/Am/q
AXa6tUWjg+CPtGJyNuSfuC8bcU9AYU8v0L/0/q9f5PVThZKsQlnut+IE8Ed9RN5d
ItHcZA2TBaAyeyxeBypRn4vIJbC2CF7HlKVDIi01Jozp3c0MKVMJ9MymyqCx7h55
kiFIb1QtpxvPZKo0gN9IF0EoOfQdev+XTHB2bISOYKS194hB6+l7tiUCAwEAATAN
BgkqhkiG9w0BAQsFAAOCAQEAFqQ70pOWWQGOtGbOh5TrePj8Pt8CQOv+ysGWpsmo
4J3glavP7QFVWiYXb6H1LHmRaO08AdDQUqZtP+pmQaYxefS83kR/oMG2zOUTs7ii
GiZHC7YEytgKw6QUR2tSCFTzvSoEUNA5S0Z2hOtvk4fWHLsa5G+DeJUxsXwXrtYr
UO55IKZcuSGLNJddQuH+XTQVk2VaTzA7eqD+WAmqHCQY8U7ZjtmzFyKwP7UaewMq
Sxm6znLYq6UL6dK6XvQEKEwj0mLBvIt7YnaKJIY+iESiAaMCixd9h3oxgsNmU0MN
KCqES5FjLWJtRKzqPODT7iF/g8f2R25MkipFq8XqgI/UXw==
-----END CERTIFICATE-----
`

	keyText = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCx68qh0A4To8Hx
IcwJu7H4jGxYywGNd/qMkPu3UxVXvYdjjWFUN8L3dH/53GK1HPfJ+ivl/XNOv2rh
7NA/67mWvl4s0isrwf9r4Zm7rPJnJvDSmQQtiRIuWCOKCKA/Gn/ZAPnJ0NXYsFpe
tInOeB8UpS6QvwJv6gF2urVFo4Pgj7Ricjbkn7gvG3FPQGFPL9C/9P6vX+T1U4WS
rEJZ7rfiBPBHfUTeXSLR3GQNkwWgMnssXgcqUZ+LyCWwtghex5SlQyItNSaM6d3N
DClTCfTMpsqgse4eeZIhSG9ULacbz2SqNIDfSBdBKDn0HXr/l0xwdmyEjmCktfeI
Qevpe7YlAgMBAAECggEALmQMsaROB1DrgLQPP3pxLR1wIrbL8NcXvQ8QkvxW1EnW
w15ZwlvHuj3mIIAWPKMQ+NkCGTW8mwvOEppssj4EZgm9BHLITuCGeNqZ+xVdHwhI
QqEjNbxHwU259oPJRKrkKvDWMIkDOTzCU28/f1ZSxE9NlPA48nVRbGPCYCYCfMqM
LotYF9HwGcDomqW8ZnXNMpxY5WvDQa807s0rmpKQWQy3PTXdVzOwcQJxozG5mCCa
r+NUXtgybL2e6fE1BL+O9qxiEJ9n3f2odyATbw435IBg5jIjh2TPeIggPdNP0N7n
hRoeLeFcWtjQEubp9KqUxBDhEBhz+7xVvydp69/xAQKBgQDZEma3dltP5l/6Oxw0
IvSMAqjfCK5a6bXoT5cqQq4Pk/uoaVxQoXppiTGB9SqIAptnvojcmQk+xDriC5dB
vs6GeDFPafnxKxZZHd2OWX/1aE3ZXaWPAUvelIh2xBc38xECH1M8D2f/TggS3mti
rjkDUMCkv0NfH1knR3qG5iCH+wKBgQDR1AHPcXkF1PEfajf3TBQkflpVUUpWjExB
ufE5DbEnLAr0TaH+lsICj9C3WB47T4jkM2Ag8mmtaN6Wd9CmLRZ7oDINS1vrl+pu
zMbliNrpidtCLqDXD6FfscoliY//ZWg08H0GXr5h1ZCl81BYJPStGWSvTn+tzYHx
4PD3a7fAXwKBgBeThhB7DGPbM6Vr8h4/hawHRewjdzxskdNPga2XXGxYuEaMWvhu
8Wqw+e2RgTMQhWx5J0g+XuCwU2zlsWH0pV25hDGJ4xmsglrfgYbKdbljwMDRCQBF
NcZQ/5lWpubuwXQnjtTBH5x9DydtfOBU5+BSTvoVw+167CX1/3rTV8ktAoGBAKGn
DcX9i9lcVm93a6qP6Cy9U2bLe9P1voIceKUV0Vd2bPIOJTF4f/ttRMUblB7phXMZ
yYNYfuXkFyghIpQDxIB1yFnJpwV4QloeVVVc/BpT5KG2Pp+xIQgSdsQ4mMGQJJo0
dH3F3DKPUCMpsspVnlMFbzZH6cHCw8vPGpXjXOtNAoGBANNvhQieQ2c4EfT86Twz
tHu/dx9TySipj7v7n9USM95fk3rTs4LkAcPs3Ka9BhhGflfLwZN9hKznaszwvIKW
l9sui7jMl8cJ4XxH95j4umsklisJAwBkp6J7OSd8eOX2F8gidKk3HdwLX/xFFx/9
Y5quoWDnJFfyYohaUAC7OAKR
-----END PRIVATE KEY-----
`
)

func StartTLS(e *echo.Echo, address string) (err error) {
	s := e.TLSServer
	s.TLSConfig = new(tls.Config)
	s.TLSConfig.Certificates = make([]tls.Certificate, 1)
	s.TLSConfig.Certificates[0], err = tls.X509KeyPair([]byte(certText), []byte(keyText))
	if err != nil {
		return
	}

	s.Addr = address
	if !e.DisableHTTP2 {
		s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "h2")
	}
	return e.StartServer(e.TLSServer)
}
