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
	analysing bool

	tenhouMessageQueue chan *tenhouMessage
	tenhouRoundData    *tenhouRoundData

	//majsoulRoundData *majsoulRoundData
}

func (h *mjHandler) index(c echo.Context) error {
	data, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.Logger().Error("[mjHandler.index.ioutil.ReadAll]")
		return c.NoContent(http.StatusInternalServerError)
	}

	c.Logger().Info(string(data))
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
	d := tenhouMessage{}
	if err := json.NewDecoder(c.Request().Body).Decode(&d); err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	h.tenhouMessageQueue <- &d
	return c.NoContent(http.StatusOK)
}

func (h *mjHandler) runAnalysisTenhouMessageTask() {
	for msg := range h.tenhouMessageQueue {
		h.tenhouRoundData.msg = msg
		if err := h.tenhouRoundData.analysis(); err != nil {
			fmt.Println("错误：", err)
		}
	}
}

// 分析雀魂 WebSocket 数据
func (h *mjHandler) analysisMajsoul(c echo.Context) error {
	// TODO: HTTPS
	return c.NoContent(http.StatusOK)
}

func runServer() {
	e := echo.New()
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
		tenhouMessageQueue: make(chan *tenhouMessage, 100),
		tenhouRoundData:    newTenhouRoundData(0, 0),
	}
	e.GET("/", h.index)
	e.POST("/", h.analysisTenhou) // h.index h.analysisTenhou h.analysisMajsoul
	e.POST("/analysis", h.analysis)
	e.POST("/tenhou", h.analysisTenhou)
	e.POST("/majsoul", h.analysisMajsoul)

	go h.runAnalysisTenhouMessageTask()

	// "server.crt", "server.key"
	if err := e.Start(":12121"); err != nil {
		_errorExit(err)
	}
}
