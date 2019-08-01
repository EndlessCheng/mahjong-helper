package main

import (
	"testing"
	"fmt"
	"encoding/json"
	"strings"
	"io/ioutil"
	"github.com/EndlessCheng/mahjong-helper/platform/tenhou"
	"github.com/EndlessCheng/mahjong-helper/debug"
)

func init() {
	debugMode = true
}

type _logContent struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

func Test_mjHandler_handleTenhouMessage(t *testing.T) {
	logFile := "log/gamedata-20190801-113014-c.log"
	logData, err := ioutil.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	// config
	startLo := -1
	endLo := -1
	lines := strings.Split(string(logData), "\n")
	if startLo == -1 {
		// 取最近游戏的日志
		for i := len(lines) - 1; i >= 0; i-- {
			if strings.Contains(lines[i], "==============") {
				startLo = i + 3
				break
			}
		}
	}
	if endLo == -1 {
		endLo = len(lines)
	}

	h := &mjHandler{
		tenhouRoundData: newGame(nil),
	}

	mr := tenhou.NewMessageReceiverWithSize(10000)
	mr.SkipOrderCheck = true
	for _, line := range lines[startLo-1 : endLo] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		d := _logContent{}
		if err := json.Unmarshal([]byte(line), &d); err != nil {
			fmt.Println(err)
			continue
		}

		if d.Level != "INFO" {
			fmt.Println(d.Message)
			continue
		}

		debug.Lo++
		mr.Put([]byte(d.Message))
		msg := mr.Get()
		h.handleTenhouMessage(msg)
	}
}

//func Test_mjHandler_runAnalysisMajsoulMessageTask(t *testing.T) {
//	logFile := "log/gamedata.log"
//	startLo := 33020
//	endLo := 33369
//
//	h := &mjHandler{
//		majsoulMessageQueue:  make(chan []byte, 10000),
//		majsoulRoundData:     &majsoulRoundData{},
//		majsoulRecordGameMap: map[string]*lq.RecordGame{},
//	}
//	h.majsoulRoundData.roundData = newGame(h.majsoulRoundData)
//
//	s := struct {
//		Level   string `json:"level"`
//		Message string `json:"message"`
//	}{}
//	logData, err := ioutil.ReadFile(logFile)
//	if err != nil {
//		t.Fatal(err)
//	}
//	lines := strings.Split(string(logData), "\n")
//	for i, line := range lines[startLo-1 : endLo] {
//		debug.Lo = i + 1
//		if line == "" {
//			continue
//		}
//		if err := json.Unmarshal([]byte(line), &s); err != nil {
//			fmt.Println(err)
//			continue
//		}
//		if s.Level != "INFO" {
//			fmt.Println(s.Level, s.Message)
//			//t.Fatal(s.Message)
//			break
//		}
//		h.majsoulMessageQueue <- []byte([]byte(s.Message))
//	}
//
//	go h.runAnalysisMajsoulMessageTask()
//
//	for {
//		if len(h.majsoulMessageQueue) == 0 {
//			break
//		}
//		time.Sleep(time.Second)
//	}
//}
