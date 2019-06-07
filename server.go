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
	"github.com/EndlessCheng/mahjong-helper/util/debug"
	"github.com/EndlessCheng/mahjong-helper/util"
)

var logFile = "gamedata.log"

func init() {
	if version == versionDev {
		const logDir = "log"
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			panic(err)
		}
		logFile = fmt.Sprintf(logDir+"/gamedata-%s.log", time.Now().Format("20060102-150405"))
	}
}

type mjHandler struct {
	log echo.Logger

	analysing bool

	tenhouMessageQueue chan []byte
	tenhouRoundData    *tenhouRoundData

	majsoulMessageQueue chan []byte
	majsoulRoundData    *majsoulRoundData

	majsoulRecordMap            map[string]*majsoulRecordBaseInfo
	majsoulCurrentRecordUUID    string
	majsoulCurrentRecordActions [][]*majsoulRecordAction
	majsoulCurrentRoundIndex    int
	majsoulCurrentActionIndex   int
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

	h.log.Info(string(data))
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

// 分析天凤 WebSocket 数据
func (h *mjHandler) analysisTenhou(c echo.Context) error {
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		h.logError(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	h.tenhouMessageQueue <- data
	return c.NoContent(http.StatusOK)
}
func (h *mjHandler) runAnalysisTenhouMessageTask() {
	if !debugMode {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("内部错误：", err)
			}
		}()
	}

	for msg := range h.tenhouMessageQueue {
		d := tenhouMessage{}
		if err := json.Unmarshal(msg, &d); err != nil {
			h.logError(err)
			continue
		}

		// FIX: 如果在 isRoundEnd 为 true 时收到了 IsSelfDraw()，则（等待收到 INIT）将其放入后面再解析
		if h.tenhouRoundData.isRoundEnd {
			if isTenhouSelfDraw(d.Tag) {
				time.Sleep(100 * time.Millisecond)
				h.tenhouMessageQueue <- msg
				continue
			}
		}

		originJSON := string(msg)
		if h.log != nil {
			h.log.Info(originJSON)
		}

		h.tenhouRoundData.msg = &d
		h.tenhouRoundData.originJSON = originJSON
		if err := h.tenhouRoundData.analysis(); err != nil {
			h.logError(err)
		}
	}
}

// 分析雀魂 WebSocket 数据
func (h *mjHandler) analysisMajsoul(c echo.Context) error {
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		h.logError(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	h.majsoulMessageQueue <- data
	return c.NoContent(http.StatusOK)
}
func (h *mjHandler) runAnalysisMajsoulMessageTask() {
	if !debugMode {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("内部错误：", err)
			}
		}()
	}

	for msg := range h.majsoulMessageQueue {
		d := &majsoulMessage{}
		if err := json.Unmarshal(msg, d); err != nil {
			h.logError(err)
			continue
		}

		originJSON := string(msg)
		if h.log != nil && debug.Lo == 0 {
			h.log.Info(originJSON)
		} else {
			if len(originJSON) > 500 {
				originJSON = originJSON[:500]
			}
			fmt.Println(originJSON)
		}

		switch {
		case len(d.Friends) > 0:
			// 处理好友列表
			fmt.Println("好友账号ID   好友上次登录时间        好友上次登出时间       好友昵称")
			for _, friend := range d.Friends {
				fmt.Println(friend)
			}
		case len(d.RecordBaseInfoList) > 0:
			// 处理牌谱基本信息列表
			for _, record := range d.RecordBaseInfoList {
				h.majsoulRecordMap[record.UUID] = record
			}
			color.HiGreen("收到 %2d 个雀魂牌谱（已收集 %d 个），请在网页上点击「查看」", len(d.RecordBaseInfoList), len(h.majsoulRecordMap))
		case d.CurrentRecordUUID != "":
			baseInfo, ok := h.majsoulRecordMap[d.CurrentRecordUUID]
			if !ok {
				h.logError(fmt.Errorf("错误：找不到雀魂牌谱 %s", h.majsoulCurrentRecordUUID))
				break
			}
			// 标记当前正在观看的牌谱
			h.majsoulCurrentRecordUUID = d.CurrentRecordUUID
			clearConsole()
			fmt.Printf("正在解析雀魂牌谱：%s", baseInfo.String())

			// 标记古役模式
			isGuyiMode := baseInfo.Config.isGuyiMode()
			util.SetConsiderOldYaku(isGuyiMode)
			if isGuyiMode {
				fmt.Println()
				color.HiGreen("古役模式已开启")
			}

			gameConf.addMajsoulAccountID(d.AccountID)
			if d.AccountID != gameConf.currentActiveMajsoulAccountID {
				fmt.Println()
				printAccountInfo(d.AccountID)
				gameConf.currentActiveMajsoulAccountID = d.AccountID
			}
		case len(d.RecordActions) > 0:
			if h.majsoulCurrentRecordUUID == "" {
				h.logError(fmt.Errorf("错误：程序未收到所观看的雀魂牌谱的 UUID"))
				break
			}

			baseInfo, ok := h.majsoulRecordMap[h.majsoulCurrentRecordUUID]
			if !ok {
				h.logError(fmt.Errorf("错误：找不到雀魂牌谱 %s", h.majsoulCurrentRecordUUID))
				break
			}

			if gameConf.currentActiveMajsoulAccountID == -1 {
				h.logError(fmt.Errorf("错误：当前雀魂账号为空"))
				break
			}

			// 获取并设置第一局的庄家
			firstRoundDealer, err := baseInfo.getFistRoundDealer(gameConf.currentActiveMajsoulAccountID)
			if err != nil {
				h.logError(err)
				break
			}
			h.majsoulRoundData.reset(0, 0, firstRoundDealer)
			h.majsoulRoundData.gameMode = gameModeRecord

			// 设置初始座位
			h.majsoulRoundData.selfSeat, err = baseInfo.getSelfSeat(gameConf.currentActiveMajsoulAccountID)
			if err != nil {
				h.logError(err)
				break
			}

			// 准备分析……
			majsoulCurrentRecordActions, err := parseMajsoulRecordAction(d.RecordActions)
			if err != nil {
				h.logError(err)
				break
			}
			h.majsoulCurrentRecordActions = majsoulCurrentRecordActions
			h.majsoulCurrentRoundIndex = 0
			h.majsoulCurrentActionIndex = 0

			actions := h.majsoulCurrentRecordActions[0]

			// 创建分析任务
			globalAnalysisCache = newGameAnalysisCache(h.majsoulCurrentRecordUUID, h.majsoulRoundData.selfSeat)
			go globalAnalysisCache.runMajsoulRecordAnalysisTask(actions)

			// 分析第一局的起始信息
			data := actions[0].Action
			h._analysisMajsoulRoundData(data, originJSON)
		case d.RecordClickAction != "":
			// 处理网页上的牌谱点击：上一局/跳到某局/下一局/上一巡/跳到某巡/下一巡/上一步/播放/暂停/下一步/点击桌面
			// 暂不能分析他家手牌
			h._onRecordClick(d.RecordClickAction, d.RecordClickActionIndex, d.FastRecordTo)
		default:
			// TODO: 重构
			//if d.AccountID > 0 {
			//	break
			//}

			//if debugMode {
			//	color.HiGreen("加入对战……")
			//}
			//h.majsoulCurrentRecordUUID = ""

			// 其他：场况分析
			h._analysisMajsoulRoundData(d, originJSON)
		}
	}
}

func (h *mjHandler) _analysisMajsoulRoundData(data *majsoulMessage, originJSON string) {
	h.majsoulRoundData.msg = data
	h.majsoulRoundData.originJSON = originJSON
	if err := h.majsoulRoundData.analysis(); err != nil {
		h.logError(err)
	}
}

func (h *mjHandler) _onRecordClick(clickAction string, clickActionIndex int, fastRecordTo int) {
	if debugMode {
		fmt.Println("收到", clickAction, clickActionIndex)
	}

	switch clickAction {
	case "nextStep", "update":
		newActionIndex := h.majsoulCurrentActionIndex + 1
		if newActionIndex >= len(h.majsoulCurrentRecordActions[h.majsoulCurrentRoundIndex]) {
			return
		}
		h.majsoulCurrentActionIndex = newActionIndex
	case "nextRound":
		h.majsoulCurrentRoundIndex = (h.majsoulCurrentRoundIndex + 1) % len(h.majsoulCurrentRecordActions)
		h.majsoulCurrentActionIndex = 0
		go globalAnalysisCache.runMajsoulRecordAnalysisTask(h.majsoulCurrentRecordActions[h.majsoulCurrentRoundIndex])
	case "preRound":
		h.majsoulCurrentRoundIndex = (h.majsoulCurrentRoundIndex - 1 + len(h.majsoulCurrentRecordActions)) % len(h.majsoulCurrentRecordActions)
		h.majsoulCurrentActionIndex = 0
		go globalAnalysisCache.runMajsoulRecordAnalysisTask(h.majsoulCurrentRecordActions[h.majsoulCurrentRoundIndex])
	case "jumpRound":
		h.majsoulCurrentRoundIndex = clickActionIndex % len(h.majsoulCurrentRecordActions)
		h.majsoulCurrentActionIndex = 0
		go globalAnalysisCache.runMajsoulRecordAnalysisTask(h.majsoulCurrentRecordActions[h.majsoulCurrentRoundIndex])
	case "nextXun", "preXun", "jumpXun", "preStep", "jumpToLastRoundXun":
		if clickAction == "jumpToLastRoundXun" {
			h.majsoulCurrentRoundIndex = (h.majsoulCurrentRoundIndex - 1 + len(h.majsoulCurrentRecordActions)) % len(h.majsoulCurrentRecordActions)
			go globalAnalysisCache.runMajsoulRecordAnalysisTask(h.majsoulCurrentRecordActions[h.majsoulCurrentRoundIndex])
		}

		h.majsoulRoundData.skipOutput = true
		currentRoundActions := h.majsoulCurrentRecordActions[h.majsoulCurrentRoundIndex]
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
		h.majsoulRoundData.skipOutput = false

		h.majsoulCurrentActionIndex = endActionIndex + 1
	default:
		return
	}

	if debugMode {
		fmt.Printf("处理牌谱中的操作：局 %d 动作 %d\n", h.majsoulCurrentRoundIndex, h.majsoulCurrentActionIndex)
	}
	action := h.majsoulCurrentRecordActions[h.majsoulCurrentRoundIndex][h.majsoulCurrentActionIndex]
	h._analysisMajsoulRoundData(action.Action, "")

	if action.Name == "RecordHule" || action.Name == "RecordLiuJu" || action.Name == "RecordNoTile" {
		// 播放和牌/流局动画，进入下一局或显示终局动画
		h.majsoulCurrentRoundIndex++
		h.majsoulCurrentActionIndex = 0
		if h.majsoulCurrentRoundIndex == len(h.majsoulCurrentRecordActions) {
			h.majsoulCurrentRoundIndex = 0
			return
		}

		time.Sleep(time.Second)

		actions := h.majsoulCurrentRecordActions[h.majsoulCurrentRoundIndex]
		go globalAnalysisCache.runMajsoulRecordAnalysisTask(actions)
		// 分析下一局的起始信息
		data := actions[h.majsoulCurrentActionIndex].Action
		h._analysisMajsoulRoundData(data, "")
	}
}

var h *mjHandler

func getMajsoulCurrentRecordUUID() string {
	return h.majsoulCurrentRecordUUID
}

func runServer(isHTTPS bool) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// 默认是 log.ERROR
	e.Logger.SetLevel(log.INFO)

	// 设置输出
	logFile, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	e.Logger.SetOutput(logFile)
	e.Logger.Info("============================================================================================")
	e.Logger.Info("服务启动")

	h = &mjHandler{
		log: e.Logger,

		tenhouMessageQueue:  make(chan []byte, 100),
		tenhouRoundData:     &tenhouRoundData{isRoundEnd: true},
		majsoulMessageQueue: make(chan []byte, 100),
		majsoulRoundData:    &majsoulRoundData{},
		majsoulRecordMap:    map[string]*majsoulRecordBaseInfo{},
	}
	h.tenhouRoundData.roundData = newGame(h.tenhouRoundData)
	h.majsoulRoundData.roundData = newGame(h.majsoulRoundData)

	go h.runAnalysisTenhouMessageTask()
	go h.runAnalysisMajsoulMessageTask()

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.GET("/", h.index)
	e.POST("/analysis", h.analysis)
	e.POST("/tenhou", h.analysisTenhou)
	e.POST("/majsoul", h.analysisMajsoul)

	const addr = ":12121"
	if !isHTTPS {
		e.POST("/", h.analysisTenhou)
		err = e.Start(addr)
	} else {
		e.POST("/", h.analysisMajsoul)
		err = startTLS(e, addr)
	}
	if err != nil {
		// 检查是否为端口占用错误
		if opErr, ok := err.(*net.OpError); ok && opErr.Op == "listen" {
			if syscallErr, ok := opErr.Err.(*os.SyscallError); ok && syscallErr.Syscall == "bind" {
				color.HiRed(addr + " 端口已被占用，程序无法启动（是否已经开启了本程序？）")
			}
		}
		errorExit(err)
	}
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
