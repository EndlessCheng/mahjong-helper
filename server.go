package main

import (
	"github.com/labstack/echo"
	"net/http"
	"time"
	"github.com/labstack/echo/middleware"
)

type echoHandler struct {
	analysing bool
	num       int
	cnt       []int
}

func (h *echoHandler) index(c echo.Context) error {
	return c.String(http.StatusOK, time.Now().Format("2006-01-02 15:04:05"))
}

func (h *echoHandler) interact(c echo.Context) error {
	if h.analysing {
		return c.NoContent(http.StatusForbidden)
	}

	h.analysing = true
	defer func() { h.analysing = false }()

	d := struct {
		Tiles   string `json:"tiles"`
		Discard string `json:"discard"`
		Draw    string `json:"draw"`
	}{}
	if err := c.Bind(&d); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var raw string
	switch {
	case d.Tiles != "":
		raw = d.Tiles
	case d.Discard != "":
		if h.num == 0 {
			return c.String(http.StatusBadRequest, "未填入手牌（需要先传入 tiles）")
		}

		idx, err := _convert(d.Discard)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		if h.cnt[idx] == 0 {
			return c.String(http.StatusBadRequest, "切掉的牌不存在")
		}
		h.cnt[idx]--
		raw = countToString(h.cnt)
		h.cnt[idx]++

		detailFlag = true
		defer func() { detailFlag = false }()
	case d.Draw != "":
		if h.num == 0 {
			return c.String(http.StatusBadRequest, "未填入手牌（需要先传入 tiles）")
		}

		idx, err := _convert(d.Draw)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		if h.cnt[idx] == 4 {
			return c.String(http.StatusBadRequest, "不可能摸更多的牌了")
		}
		h.cnt[idx]++
		raw = countToString(h.cnt)
		h.cnt[idx]--
	default:
		return c.String(http.StatusBadRequest, "Empty input!")
	}

	num, cnt, err := analysis(raw)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	h.num = num
	h.cnt = cnt

	return c.NoContent(http.StatusOK)
}

func runServer() {
	h := &echoHandler{}

	e := echo.New()
	e.Use(middleware.Recover())

	e.GET("/", h.index)
	e.POST("/interact", h.interact)

	if err := e.Start(":12121"); err != nil {
		_errorExit(err.Error())
	}
}
