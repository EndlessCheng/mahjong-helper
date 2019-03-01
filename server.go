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
)

type mjHandler struct {
	log echo.Logger

	analysing bool

	tenhouMessageQueue chan []byte
	tenhouRoundData    *tenhouRoundData

	majsoulMessageQueue chan []byte
	majsoulRoundData    *majsoulRoundData
}

func (h *mjHandler) index(c echo.Context) error {
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		h.log.Error("[mjHandler.index.ioutil.ReadAll]", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	h.log.Info(string(data))
	return c.String(http.StatusOK, time.Now().Format("2006-01-02 15:04:05"))
}

func (h *mjHandler) analysis(c echo.Context) error {
	if h.analysing {
		return c.NoContent(http.StatusForbidden)
	}

	h.analysing = true
	defer func() { h.analysing = false }()

	d := struct {
		Reset      bool   `json:"reset"`
		Tiles      string `json:"tiles"`
		ShowDetail bool   `json:"show_detail"`
	}{}
	if err := c.Bind(&d); err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	if d.ShowDetail {
		detailFlag = true
		defer func() { detailFlag = false }()
	}
	if _, _, err := analysis(d.Tiles); err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// 分析天凤 WebSocket 数据
func (h *mjHandler) analysisTenhou(c echo.Context) error {
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	h.tenhouMessageQueue <- data
	return c.NoContent(http.StatusOK)
}

func (h *mjHandler) runAnalysisTenhouMessageTask() {
	for msg := range h.tenhouMessageQueue {
		d := tenhouMessage{}
		if err := json.Unmarshal(msg, &d); err != nil {
			fmt.Println(err)
			continue
		}

		originJSON := string(msg)
		h.log.Info(originJSON)

		h.tenhouRoundData.msg = &d
		h.tenhouRoundData.originJSON = originJSON
		if err := h.tenhouRoundData.analysis(); err != nil {
			fmt.Println("错误：", err)
		}
	}
}

// 分析雀魂 WebSocket 数据
func (h *mjHandler) analysisMajsoul(c echo.Context) error {
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	h.majsoulMessageQueue <- data
	return c.NoContent(http.StatusOK)
}

func (h *mjHandler) runAnalysisMajsoulMessageTask() {
	for msg := range h.majsoulMessageQueue {
		d := majsoulMessage{}
		if err := json.Unmarshal(msg, &d); err != nil {
			fmt.Println(err)
			continue
		}

		originJSON := string(msg)
		h.log.Info(originJSON)

		if d.AccountID > 0 && h.majsoulRoundData.accountID != d.AccountID {
			h.majsoulRoundData.accountID = d.AccountID
			printAccountInfo(d.AccountID)
			continue
		}

		if d.Friends != nil {
			fmt.Println("好友账号ID   好友上次登录时间        好友上次登出时间       好友昵称")
			for _, friend := range d.Friends {
				fmt.Println(friend)
			}
			continue
		}

		h.majsoulRoundData.msg = &d
		h.majsoulRoundData.originJSON = originJSON
		if err := h.majsoulRoundData.analysis(); err != nil {
			fmt.Println("错误：", err)
		}
	}
}

func runServer(isHTTPS bool) {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// 默认是 log.ERROR
	e.Logger.SetLevel(log.INFO)
	go func() {
		// 等待服务启动再设置输出
		time.Sleep(time.Second)
		logFile, err := os.OpenFile("gamedata.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
		e.Logger.SetOutput(logFile)
	}()

	h := &mjHandler{
		log: e.Logger,

		tenhouMessageQueue:  make(chan []byte, 100),
		tenhouRoundData:     &tenhouRoundData{},
		majsoulMessageQueue: make(chan []byte, 100),
		majsoulRoundData:    &majsoulRoundData{accountID: -1},
	}
	h.tenhouRoundData.roundData = newRoundData(h.tenhouRoundData, 0, 0)
	h.majsoulRoundData.roundData = newRoundData(h.majsoulRoundData, 0, 0)

	go h.runAnalysisTenhouMessageTask()
	go h.runAnalysisMajsoulMessageTask()

	e.GET("/", h.index)
	e.POST("/analysis", h.analysis)
	e.POST("/tenhou", h.analysisTenhou)
	e.POST("/majsoul", h.analysisMajsoul)

	addr := ":12121"
	if !isHTTPS {
		e.POST("/", h.analysisTenhou)
		if err := e.Start(addr); err != nil {
			_errorExit(err)
		}
	} else {
		e.POST("/", h.analysisMajsoul)
		if err := e.StartTLS(addr, "selfsigned.crt", "selfsigned.key"); err != nil {
			_errorExit(err)
		}
	}
}
