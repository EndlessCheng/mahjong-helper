package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/EndlessCheng/mahjong-helper/platform/tenhou"
	"github.com/EndlessCheng/mahjong-helper/util"
	"github.com/EndlessCheng/mahjong-helper/util/debug"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	stdLog "log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
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
	tenhouRoundData       *tenhouRoundData

	majsoulMessageQueue chan []byte
	majsoulRoundData    *majsoulRoundData

	majsoulRecordMap                map[string]*majsoulRecordBaseInfo
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

// 分析天凤 WebSocket 数据
func (h *mjHandler) analysisTenhou(c echo.Context) error {
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		h.logError(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	h.tenhouMessageReceiver.Put(data)
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

	for {
		msg := h.tenhouMessageReceiver.Get()
		d := tenhouMessage{}
		if err := json.Unmarshal(msg, &d); err != nil {
			h.logError(err)
			continue
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
			// 好友列表
			fmt.Println(d.Friends)
		case len(d.RecordBaseInfoList) > 0:
			// 牌谱基本信息列表
			for _, record := range d.RecordBaseInfoList {
				h.majsoulRecordMap[record.UUID] = record
			}
			color.HiGreen("收到 %2d 个雀魂牌谱（已收集 %d 个），请在网页上点击「查看」", len(d.RecordBaseInfoList), len(h.majsoulRecordMap))
		case d.SharedRecordBaseInfo != nil:
			// 处理分享的牌谱基本信息
			// FIXME: 观看自己的牌谱也会有 d.SharedRecordBaseInfo
			record := d.SharedRecordBaseInfo
			h.majsoulRecordMap[record.UUID] = record
			if err := h._loadMajsoulRecordBaseInfo(record.UUID); err != nil {
				h.logError(err)
				break
			}
		case d.CurrentRecordUUID != "":
			// 载入某个牌谱
			resetAnalysisCache()
			h.majsoulCurrentRecordActionsList = nil

			if err := h._loadMajsoulRecordBaseInfo(d.CurrentRecordUUID); err != nil {
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
			if h.majsoulCurrentRecordActionsList != nil {
				// TODO: 网页发送更恰当的信息？
				break
			}

			if h.majsoulCurrentRecordUUID == "" {
				h.logError(fmt.Errorf("错误：程序未收到所观看的雀魂牌谱的 UUID"))
				break
			}

			baseInfo, ok := h.majsoulRecordMap[h.majsoulCurrentRecordUUID]
			if !ok {
				h.logError(fmt.Errorf("错误：找不到雀魂牌谱 %s", h.majsoulCurrentRecordUUID))
				break
			}

			selfAccountID := gameConf.currentActiveMajsoulAccountID
			if selfAccountID == -1 {
				h.logError(fmt.Errorf("错误：当前雀魂账号为空"))
				break
			}

			h.majsoulRoundData.newGame()
			h.majsoulRoundData.gameMode = gameModeRecord

			// 获取并设置主视角初始座位
			selfSeat, err := baseInfo.getSelfSeat(selfAccountID)
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
			gameConf.setMajsoulAccountID(1) // TODO: 重构
			h.majsoulRoundData.newGame()
			h.majsoulRoundData.selfSeat = 0 // 观战进来后看的是东起的玩家
			h.majsoulRoundData.gameMode = gameModeLive
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
			if h.majsoulRoundData.gameMode == gameModeLive { // 观战
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

func (h *mjHandler) _loadMajsoulRecordBaseInfo(majsoulRecordUUID string) error {
	baseInfo, ok := h.majsoulRecordMap[majsoulRecordUUID]
	if !ok {
		return fmt.Errorf("错误：找不到雀魂牌谱 %s", majsoulRecordUUID)
	}

	// 标记当前正在观看的牌谱
	h.majsoulCurrentRecordUUID = majsoulRecordUUID
	clearConsole()
	fmt.Printf("正在解析雀魂牌谱：%s", baseInfo.String())

	// 标记古役模式
	isGuyiMode := baseInfo.Config.isGuyiMode()
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

	h.majsoulRoundData.skipOutput = isFast
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
	h.majsoulRoundData.skipOutput = true
	// 留最后三个刷新，这样确保会刷新界面
	for _, action := range actions[:fastRecordEnd] {
		h._analysisMajsoulRoundData(action.Action, "")
	}
	h.majsoulRoundData.skipOutput = false
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

		h.majsoulRoundData.skipOutput = true
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
		h.majsoulRoundData.skipOutput = false

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

	// 移除 echo.Echo 和 http.Server 在控制台上打印的信息
	e.HideBanner = true
	e.HidePort = true
	e.StdLogger = stdLog.New(ioutil.Discard, "", 0)

	// 默认是 log.ERROR
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
		tenhouRoundData:       &tenhouRoundData{isRoundEnd: true},
		majsoulMessageQueue:   make(chan []byte, 100),
		majsoulRoundData:      &majsoulRoundData{selfSeat: -1},
		majsoulRecordMap:      map[string]*majsoulRecordBaseInfo{},
	}

	h.majsoulRoundData.roundData = newGame(h.majsoulRoundData)


	go h.runAnalysisMajsoulMessageTask()

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.GET("/", h.index)
	e.POST("/debug", h.index)
	e.POST("/analysis", h.analysis)
	e.POST("/tenhou", h.analysisTenhou)
	e.POST("/majsoul", h.analysisMajsoul)

	// code.js 也用的该端口
	if port == 0 {
		port = defaultPort
	}
	addr := ":" + strconv.Itoa(port)

	e.POST("/", h.analysisMajsoul)
	err = startTLS(e, addr)

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
MIIEcjCCAtqgAwIBAgIRAKRU8OQwtD5xe2UzG+QmQRowDQYJKoZIhvcNAQELBQAw
gY8xHjAcBgNVBAoTFW1rY2VydCBkZXZlbG9wbWVudCBDQTEyMDAGA1UECwwpb3Jp
Z2luQHFpbmdjaGVuZ2RlQWlyLm1zaG9tZS5uZXQgKOa4heapmSkxOTA3BgNVBAMM
MG1rY2VydCBvcmlnaW5AcWluZ2NoZW5nZGVBaXIubXNob21lLm5ldCAo5riF5qmZ
KTAeFw0yMTAxMDUxNjExMDZaFw0yMzA0MDUxNjExMDZaMFsxJzAlBgNVBAoTHm1r
Y2VydCBkZXZlbG9wbWVudCBjZXJ0aWZpY2F0ZTEwMC4GA1UECwwnb3JpZ2luQHFp
bmdjaGVuZ2RlQWlyLmJicm91dGVyICjmuIXmqZkpMIIBIjANBgkqhkiG9w0BAQEF
AAOCAQ8AMIIBCgKCAQEAmLfvXUvBhffPgGDhJhWJNawrCAGAnaI+36V6qqAei/En
uwaet0J/1YXcmNhRe8GT+xm1NhtIbRSMl+XJdWr8YOQnaK2UxKw5tUaBoJC6xUj1
NAge7btrjtr06oeg2azS81DR8aHA1FY+e63deXVrmZ/qkW0FMh/ru7py+m3Y8TB7
Qpvka0mWqR8DrcXG0vpDNvEQL007YulAXtPE2GTL+l27w9aaLP7qvixn8CGlKuyr
d6NZ55xvj1xACDbGSfi+YqS/E7SR7FcDOK6DfjeQn2zFdsxxk/64wFx3CzUft+cy
ZqQGXYsxfHwQ+tIbMON4eSaUrKH8R720KcVfOJgIZwIDAQABo3wwejAOBgNVHQ8B
Af8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUHAwEwHwYDVR0jBBgwFoAUr+bJE8/w
vkHn0TzTCNFjjXX4BZgwMgYDVR0RBCswKYIJbG9jYWxob3N0ghZsb2NhbGhvc3Qu
bWFqLXNvdWwuY29thwR/AAABMA0GCSqGSIb3DQEBCwUAA4IBgQBL+2uXs76scVwl
mJomnFdB/aFOudWrzhsYzcoznubkNPWuywpKjmZj4mnOCNUEul/0dnwgKp3OTuxf
+r1+ugA6h6gK57S6IeaoVYEmAVLjrwAQERqD+ZobxcBAzMlJXD85fDyLyOUp2YuI
daLi4/7Q8m5aejcVOcHeLjsyOi3eidEdmQsAz9WWoHzQF0SJ7btarT5uXR5N1k3z
QLg65geFWIkZFqHYs2NAJUorOk4ddHPkPy1G3DkX4K9D7VyOpJRAIMce5JlNUlKq
0ez221hpiTk7hw7TZN13RAUHDkIC4+ELwFJNYjBN+wEsuEslmwlqY0vMHu8Yi0C/
QT+6g4mEMXw0oYFVKulWMCW3HMKRwOmTnN8MWUBwhzop7xbhDSBByG2nPHWGRN5v
YQEMbW0dH5+JBEjgxA6T9fTY8tAuM44WaSSqhNbpd862G1Oo6d7cQyLa3BKTjX+V
TFdCQlqh2/yrRYLt5s4rOpNKMSCdPIUqdQZzg4c5+sBIZ44Ctuc=
-----END CERTIFICATE-----
`

	keyText = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCYt+9dS8GF98+A
YOEmFYk1rCsIAYCdoj7fpXqqoB6L8Se7Bp63Qn/VhdyY2FF7wZP7GbU2G0htFIyX
5cl1avxg5CdorZTErDm1RoGgkLrFSPU0CB7tu2uO2vTqh6DZrNLzUNHxocDUVj57
rd15dWuZn+qRbQUyH+u7unL6bdjxMHtCm+RrSZapHwOtxcbS+kM28RAvTTti6UBe
08TYZMv6XbvD1pos/uq+LGfwIaUq7Kt3o1nnnG+PXEAINsZJ+L5ipL8TtJHsVwM4
roN+N5CfbMV2zHGT/rjAXHcLNR+35zJmpAZdizF8fBD60hsw43h5JpSsofxHvbQp
xV84mAhnAgMBAAECggEAUvUXgYZ1SWCjxwjZRObN8enkNiUayIYpwTMSvzzbWwlU
m7Os47+r0UE145EeMiulRvFgDlQjs51GAf1AwheroLZO8f3Yoj0r29zA7Yew7RoE
mI8QvfLhKmimQgAK0DOhI6rzYF6NcMqohmabuC6glILZ2MVv3RqZ4xAVBRRGlDz6
5ORPOUP6ASSYFAyPks674rQYRsJxGiO6z2jfMPHUwp05u2c2uMceAL98Yr32tcN7
IRSo52ELAKt6Hi8sgtfZ8Cqx1U/Cs9u1d6QeA/9eF4QL9GAj9Ay2PItHG+YB7bT1
tCGnn0mGUlcUvnAMvkQYF8ll9WYnLYzsZskMHMCIgQKBgQDHgMMRW4NKXZCGT6AM
vub1AOXkowj9gHG73H31+HWFkqcqNMNh4otJt8Fw4sbZdBlQwclAVpqtenXGO0th
tmkBPOXznhkHff5SD7VpG9FQaMhRfFxFJggjPhz4e5z4tfwGJ5+YyvD64GVvSqEe
Ogdh1XoG5PfM6zIqZL3xvW3fawKBgQDD93gIWZ4esjJg4wD/djsCYPWw0enduOWX
PZ7Dh1WE2UJjW/u7r18ZqZXm0C93d1s6sSLbz0u43WWJCJEcptwLU/bthBfb3WMY
V2MHdPE5QrU57ezVpaVHx4A+pzOb4LuhTteWLkPhf3VFmoa9Wa/cc1b33tGshBhL
nNXRTuRl9QKBgQDDa7mokv+0JJqhNfYNBiKt88c9gwYXa239Gyq3ej2ELfdZPH32
sDbIaxstPLaT40m49Vnxj+PL8pzTJNneSRPqhoCpdkAGOsCYGZMV9o2+OiWezDaF
9Y8bFojCTjOg3IKWdNG8lW4gERbLQUs5lJYOm1IA1uB09h4ZsLzuwyORKQKBgQCZ
X+1NM77irXtqgyC70HA822BQBOryQw1Ggs7on9paAKTKGSr76TUYY7dUECqmaP84
/3yV5zePt5AJYXAZqardHtlLajA2P56YZYS3SFqoA0LN1R6g1GV4uXbxEnH9FTYk
+Q0YmJs+OUCyuk+skS5n7snpdDZMvJI7U8OxbvqppQKBgEYBZfDEBGUdE84iRnyV
fjQFj2D1Ei/cuI9nm51l34gVrKui7xVlIALZvNEv3lCv4b0D7VTcfaNeIr0tweW/
a6wJEKtvtWth68cNkxQudHdrELMKRY67ee4CXLz37VCzxFGtlj0Q2WpXqHsby7IZ
RPKiFQUw+KrQVyxQ0OusrzWj
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
