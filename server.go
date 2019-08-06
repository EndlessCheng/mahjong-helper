package main

import (
	"github.com/labstack/echo"
	"net/http"
	"time"
	"github.com/labstack/echo/middleware"
	"io/ioutil"
	"os"
	"github.com/labstack/gommon/log"
	"fmt"
	"encoding/json"
	"crypto/tls"
	"net"
	"github.com/fatih/color"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"github.com/EndlessCheng/mahjong-helper/util"
	"path/filepath"
	"strconv"
	"github.com/EndlessCheng/mahjong-helper/platform/tenhou"
	"github.com/EndlessCheng/mahjong-helper/platform/majsoul"
	"github.com/EndlessCheng/mahjong-helper/platform/majsoul/proto/lq"
)

const defaultPort = 12121

func newLogFilePath() (filePath string, err error) {
	const logDir = "log"
	if err = os.MkdirAll(logDir, os.ModePerm); err != nil {
		return
	}
	fileName := fmt.Sprintf("gamedata-%s.log", time.Now().Format("20060102-150405"))
	filePath = filepath.Join(logDir, fileName)
	return filepath.Abs(filePath)
}

type mjHandler struct {
	log echo.Logger

	analysing bool

	tenhouMessageReceiver *tenhou.MessageReceiver
	tenhouRoundData       *roundData

	majsoulMessageReceiver *majsoul.MessageReceiver
	majsoulUIMessageQueue  chan []byte
	majsoulRoundData       *majsoulRoundData

	majsoulRecordGameMap            map[string]*lq.RecordGame
	majsoulCurrentRecordUUID        string
	majsoulCurrentRecordActionsList []majsoulRoundActions
	majsoulCurrentRoundIndex        int
	majsoulCurrentActionIndex       int

	majsoulCurrentRoundActions majsoulRoundActions
}

func (h *mjHandler) logError(err error) {
	fmt.Fprintln(os.Stderr, err)
	if !debugMode {
		h.log.Error(err)
	}
}

// 调试用
func (h *mjHandler) index(c echo.Context) error {
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		h.log.Error("[mjHandler.index.ioutil.ReadAll]", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	fmt.Println(data, string(data))
	h.log.Info(data)
	return c.String(http.StatusOK, time.Now().Format("2006-01-02 15:04:05"))
}

// 打一摸一分析器
func (h *mjHandler) analysis(c echo.Context) error {
	if h.analysing {
		return c.NoContent(http.StatusForbidden)
	}

	h.analysing = true
	defer func() { h.analysing = false }()

	d := struct {
		Reset bool   `json:"reset"`
		Tiles string `json:"tiles"`
	}{}
	if err := c.Bind(&d); err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	if _, err := analysisHumanTiles(model.NewSimpleHumanTilesInfo(d.Tiles)); err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// 接收天凤 WebSocket 数据
func (h *mjHandler) saveTenhouWebSocketMessage(c echo.Context) error {
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		h.logError(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	h.tenhouMessageReceiver.Put(data)
	return c.NoContent(http.StatusOK)
}

// 处理天凤 WebSocket 数据
func (h *mjHandler) runAnalysisTenhouWebSocketMessageTask() {
	for {
		msg := h.tenhouMessageReceiver.Get()
		if err := h.handleTenhouWebSocketMessage(msg); err != nil {
			h.logError(fmt.Errorf("handleTenhouWebSocketMessage: %v", err))
		}
	}
}

// 接收雀魂 WebSocket 数据
func (h *mjHandler) saveMajsoulWebSocketMessage(c echo.Context) error {
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		h.logError(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	h.majsoulMessageReceiver.Put(data)
	return c.NoContent(http.StatusOK)
}

// 接收雀魂 UI 数据：牌谱点击事件等
func (h *mjHandler) saveMajsoulUIMessage(c echo.Context) error {
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		h.logError(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	h.majsoulUIMessageQueue <- data
	return c.NoContent(http.StatusOK)
}

// 处理雀魂 WebSocket 数据
func (h *mjHandler) runAnalysisMajsoulWebSocketMessageTask() {
	for {
		msg := h.majsoulMessageReceiver.Get()
		if err := h.handleMajsoulWebSocketMessage(msg); err != nil {
			h.logError(fmt.Errorf("runAnalysisMajsoulWebSocketMessageTask: %v", err))
		}
	}

	for msg := range h.majsoulUIMessageQueue {
		d := &majsoulMessage{}
		if err := json.Unmarshal(msg, d); err != nil {
			h.logError(err)
			continue
		}

		originJSON := string(msg)
		if h.log != nil {
			h.log.Info(originJSON)
		}

		switch {
		case d.CurrentRecordUUID != "":
			// 载入某个牌谱
			resetAnalysisCache()
			h.majsoulCurrentRecordActionsList = nil

			if err := h._loadMajsoulRecordBaseInfo(d.CurrentRecordUUID); err != nil {
				// 看的是分享的牌谱（先收到 CurrentRecordUUID 和 AccountID，然后收到 SharedRecordBaseInfo）
				// 或者是比赛场的牌谱
				// 记录主视角 ID（可能是 0）
				userConf.setMajsoulAccountID(d.AccountID)
				break
			}

			// 看的是自己的牌谱
			// 更新当前使用的账号
			userConf.addMajsoulAccountID(d.AccountID)
			if userConf.currentActiveMajsoulAccountID != d.AccountID {
				fmt.Println()
				printAccountInfo(d.AccountID)
				userConf.setMajsoulAccountID(d.AccountID)
			}
		case len(d.RecordActions) > 0:
			if h.majsoulCurrentRecordActionsList != nil {
				// TODO: 网页发送更恰当的信息？
				break
			}

			if h.majsoulCurrentRecordUUID == "" {
				h.logError(fmt.Errorf("错误：程序未收到所观看的雀魂牌谱的 UUID"))
				break
			}

			baseInfo, ok := h.majsoulRecordGameMap[h.majsoulCurrentRecordUUID]
			if !ok {
				h.logError(fmt.Errorf("错误：找不到雀魂牌谱 %s", h.majsoulCurrentRecordUUID))
				break
			}

			selfAccountID := userConf.currentActiveMajsoulAccountID
			if selfAccountID == -1 {
				h.logError(fmt.Errorf("错误：当前雀魂账号为空"))
				break
			}

			h.majsoulRoundData.newGame()
			h.majsoulRoundData.config.gameMode = gameModeRecord

			// 获取并设置主视角初始座位
			selfSeat, err := baseInfo.GetSelfSeat(selfAccountID)
			if err != nil {
				h.logError(err)
				break
			}
			h.majsoulRoundData.selfSeat = selfSeat

			// 准备分析……
			majsoulCurrentRecordActions, err := parseMajsoulRecordAction(d.RecordActions)
			if err != nil {
				h.logError(err)
				break
			}
			h.majsoulCurrentRecordActionsList = majsoulCurrentRecordActions
			h.majsoulCurrentRoundIndex = 0
			h.majsoulCurrentActionIndex = 0

			actions := h.majsoulCurrentRecordActionsList[h.majsoulCurrentRoundIndex]

			// 创建分析任务
			analysisCache := newGameAnalysisCache(h.majsoulCurrentRecordUUID, selfSeat)
			setAnalysisCache(analysisCache)
			go analysisCache.runMajsoulRecordAnalysisTask(actions)

			// 分析第一局的起始信息
			data := actions[0].Action
			h._analysisMajsoulRoundData(data, originJSON)
		case d.RecordClickAction != "":
			// 处理网页上的牌谱点击：上一局/跳到某局/下一局/上一巡/跳到某巡/下一巡/上一步/播放/暂停/下一步/点击桌面
			// 暂不能分析他家手牌
			h._onRecordClick(d.RecordClickAction, d.RecordClickActionIndex, d.FastRecordTo)
		case d.LiveBaseInfo != nil:
			// 观战
			userConf.setMajsoulAccountID(1) // TODO: 重构
			h.majsoulRoundData.newGame()
			h.majsoulRoundData.selfSeat = 0 // 观战进来后看的是东起的玩家
			h.majsoulRoundData.config.gameMode = gameModeLive
			clearConsole()
			fmt.Printf("正在载入对战：%s", d.LiveBaseInfo.String())
		case d.LiveFastAction != nil:
			if err := h._loadLiveAction(d.LiveFastAction, true); err != nil {
				h.logError(err)
				break
			}
		case d.LiveAction != nil:
			if err := h._loadLiveAction(d.LiveAction, false); err != nil {
				h.logError(err)
				break
			}
		case d.ChangeSeatTo != nil:
			// 切换座位
			changeSeatTo := *(d.ChangeSeatTo)
			h.majsoulRoundData.selfSeat = changeSeatTo
			if debugMode {
				fmt.Println("座位已切换至", changeSeatTo)
			}

			var actions majsoulRoundActions
			if h.majsoulRoundData.config.gameMode == gameModeLive { // 观战
				actions = h.majsoulCurrentRoundActions
			} else { // 牌谱
				fullActions := h.majsoulCurrentRecordActionsList[h.majsoulCurrentRoundIndex]
				actions = fullActions[:h.majsoulCurrentActionIndex+1]
				analysisCache := getAnalysisCache(changeSeatTo)
				if analysisCache == nil {
					analysisCache = newGameAnalysisCache(h.majsoulCurrentRecordUUID, changeSeatTo)
				}
				setAnalysisCache(analysisCache)
				// 创建分析任务
				go analysisCache.runMajsoulRecordAnalysisTask(fullActions)
			}

			h._fastLoadActions(actions)
		case len(d.SyncGameActions) > 0:
			h._fastLoadActions(d.SyncGameActions)
		default:
			// 其他：AI 分析
			h._analysisMajsoulRoundData(d, originJSON)
		}
	}
}

// 处理雀魂 UI 数据
func (h *mjHandler) runAnalysisMajsoulUIMessageTask() {
}

func (h *mjHandler) _loadMajsoulRecordBaseInfo(majsoulRecordUUID string) error {
	baseInfo, ok := h.majsoulRecordGameMap[majsoulRecordUUID]
	if !ok {
		return fmt.Errorf("错误：找不到雀魂牌谱 %s", majsoulRecordUUID)
	}

	// TODO: 处理玩家个数

	// 标记当前正在观看的牌谱
	h.majsoulCurrentRecordUUID = majsoulRecordUUID
	clearConsole()
	fmt.Printf("正在解析雀魂牌谱：%s", baseInfo.CLIString())

	// 标记古役模式
	isGuyiMode := baseInfo.Config.IsGuyiMode()
	util.SetConsiderOldYaku(isGuyiMode)
	if isGuyiMode {
		fmt.Println()
		color.HiGreen("古役模式已开启")
	}

	return nil
}

func (h *mjHandler) _loadLiveAction(action *majsoulRecordAction, isFast bool) error {
	if debugMode {
		fmt.Println("[_loadLiveAction] 收到", action, isFast)
	}

	newActions, err := h.majsoulCurrentRoundActions.append(action)
	if err != nil {
		return err
	}
	h.majsoulCurrentRoundActions = newActions

	h.majsoulRoundData.config.skipOutput = isFast
	h._analysisMajsoulRoundData(action.Action, "")
	return nil
}

func (h *mjHandler) _analysisMajsoulRoundData(data *majsoulMessage, originJSON string) {
	//if originJSON == "{}" {
	//	return
	//}
	h.majsoulRoundData.msg = data
	h.majsoulRoundData.originJSON = originJSON
	if err := h.majsoulRoundData.analysis(); err != nil {
		h.logError(err)
	}
}

func (h *mjHandler) _fastLoadActions(actions []*majsoulRecordAction) {
	if len(actions) == 0 {
		return
	}
	fastRecordEnd := util.MaxInt(0, len(actions)-3)
	h.majsoulRoundData.config.skipOutput = true
	// 留最后三个刷新，这样确保会刷新界面
	for _, action := range actions[:fastRecordEnd] {
		h._analysisMajsoulRoundData(action.Action, "")
	}
	h.majsoulRoundData.config.skipOutput = false
	for _, action := range actions[fastRecordEnd:] {
		h._analysisMajsoulRoundData(action.Action, "")
	}
}

func (h *mjHandler) _onRecordClick(clickAction string, clickActionIndex int, fastRecordTo int) {
	if debugMode {
		fmt.Println("[_onRecordClick] 收到", clickAction, clickActionIndex, fastRecordTo)
	}

	analysisCache := getCurrentAnalysisCache()

	switch clickAction {
	case "nextStep", "update":
		newActionIndex := h.majsoulCurrentActionIndex + 1
		if newActionIndex >= len(h.majsoulCurrentRecordActionsList[h.majsoulCurrentRoundIndex]) {
			return
		}
		h.majsoulCurrentActionIndex = newActionIndex
	case "nextRound":
		h.majsoulCurrentRoundIndex = (h.majsoulCurrentRoundIndex + 1) % len(h.majsoulCurrentRecordActionsList)
		h.majsoulCurrentActionIndex = 0
		go analysisCache.runMajsoulRecordAnalysisTask(h.majsoulCurrentRecordActionsList[h.majsoulCurrentRoundIndex])
	case "preRound":
		h.majsoulCurrentRoundIndex = (h.majsoulCurrentRoundIndex - 1 + len(h.majsoulCurrentRecordActionsList)) % len(h.majsoulCurrentRecordActionsList)
		h.majsoulCurrentActionIndex = 0
		go analysisCache.runMajsoulRecordAnalysisTask(h.majsoulCurrentRecordActionsList[h.majsoulCurrentRoundIndex])
	case "jumpRound":
		h.majsoulCurrentRoundIndex = clickActionIndex % len(h.majsoulCurrentRecordActionsList)
		h.majsoulCurrentActionIndex = 0
		go analysisCache.runMajsoulRecordAnalysisTask(h.majsoulCurrentRecordActionsList[h.majsoulCurrentRoundIndex])
	case "nextXun", "preXun", "jumpXun", "preStep", "jumpToLastRoundXun":
		if clickAction == "jumpToLastRoundXun" {
			h.majsoulCurrentRoundIndex = (h.majsoulCurrentRoundIndex - 1 + len(h.majsoulCurrentRecordActionsList)) % len(h.majsoulCurrentRecordActionsList)
			go analysisCache.runMajsoulRecordAnalysisTask(h.majsoulCurrentRecordActionsList[h.majsoulCurrentRoundIndex])
		}

		h.majsoulRoundData.config.skipOutput = true
		currentRoundActions := h.majsoulCurrentRecordActionsList[h.majsoulCurrentRoundIndex]
		startActionIndex := 0
		endActionIndex := fastRecordTo
		if clickAction == "nextXun" {
			startActionIndex = h.majsoulCurrentActionIndex + 1
		}
		if debugMode {
			fmt.Printf("快速处理牌谱中的操作：局 %d 动作 %d-%d\n", h.majsoulCurrentRoundIndex, startActionIndex, endActionIndex)
		}
		for i, action := range currentRoundActions[startActionIndex : endActionIndex+1] {
			if debugMode {
				fmt.Printf("快速处理牌谱中的操作：局 %d 动作 %d\n", h.majsoulCurrentRoundIndex, startActionIndex+i)
			}
			h._analysisMajsoulRoundData(action.Action, "")
		}
		h.majsoulRoundData.config.skipOutput = false

		h.majsoulCurrentActionIndex = endActionIndex + 1
	default:
		return
	}

	if debugMode {
		fmt.Printf("处理牌谱中的操作：局 %d 动作 %d\n", h.majsoulCurrentRoundIndex, h.majsoulCurrentActionIndex)
	}
	action := h.majsoulCurrentRecordActionsList[h.majsoulCurrentRoundIndex][h.majsoulCurrentActionIndex]
	h._analysisMajsoulRoundData(action.Action, "")

	if action.Name == "RecordHule" || action.Name == "RecordLiuJu" || action.Name == "RecordNoTile" {
		// 播放和牌/流局动画，进入下一局或显示终局动画
		h.majsoulCurrentRoundIndex++
		h.majsoulCurrentActionIndex = 0
		if h.majsoulCurrentRoundIndex == len(h.majsoulCurrentRecordActionsList) {
			h.majsoulCurrentRoundIndex = 0
			return
		}

		time.Sleep(time.Second)

		actions := h.majsoulCurrentRecordActionsList[h.majsoulCurrentRoundIndex]
		go analysisCache.runMajsoulRecordAnalysisTask(actions)
		// 分析下一局的起始信息
		data := actions[h.majsoulCurrentActionIndex].Action
		h._analysisMajsoulRoundData(data, "")
	}
}

var h *mjHandler

func getMajsoulCurrentRecordUUID() string {
	return h.majsoulCurrentRecordUUID
}

func runServer(isHTTPS bool, port int) (err error) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// 改成 INFO 等级（默认是 ERROR）
	e.Logger.SetLevel(log.INFO)

	// 设置日志输出到 log/gamedata-xxx.log
	filePath, err := newLogFilePath()
	if err != nil {
		return
	}
	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	e.Logger.SetOutput(logFile)

	e.Logger.Info("============================================================================================")
	e.Logger.Info("服务启动")

	h = &mjHandler{
		log: e.Logger,

		tenhouMessageReceiver: tenhou.NewMessageReceiver(),
		tenhouRoundData:       newGame(nil),

		majsoulMessageReceiver: majsoul.NewMessageReceiver(),
		majsoulUIMessageQueue:  make(chan []byte, 100),
		majsoulRoundData:       &majsoulRoundData{selfSeat: -1},
		majsoulRecordGameMap:   map[string]*lq.RecordGame{},
	}
	h.majsoulRoundData.roundData = newGame(h.majsoulRoundData)

	go h.runAnalysisTenhouWebSocketMessageTask()
	go h.runAnalysisMajsoulWebSocketMessageTask()
	go h.runAnalysisMajsoulUIMessageTask()

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.GET("/", h.index)
	e.POST("/debug", h.index)
	e.POST("/analysis", h.analysis)
	e.POST("/tenhou/ws", h.saveTenhouWebSocketMessage)
	e.POST("/majsoul/ws", h.saveMajsoulWebSocketMessage)
	e.POST("/majsoul/ui", h.saveMajsoulUIMessage)

	// code.js 也用的该端口
	if port == 0 {
		port = defaultPort
	}
	addr := ":" + strconv.Itoa(port)
	if !isHTTPS {
		err = e.Start(addr)
	} else {
		err = startTLS(e, addr)
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

func startTLS(e *echo.Echo, address string) (err error) {
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
