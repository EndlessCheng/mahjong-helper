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
		if err := startTLS(e, addr); err != nil {
			_errorExit(err)
		}
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
