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

type _logContent struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

var logs []_logContent

func init() {
	debugMode = true

	logFile := "log/gamedata-20190801-113014-c.log"
	logData, err := ioutil.ReadFile(logFile)
	if err != nil {
		panic(err)
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

	lines = lines[startLo-1 : endLo]
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		d := _logContent{}
		if err := json.Unmarshal([]byte(line), &d); err != nil {
			panic(err)
			continue
		}
		logs = append(logs, d)
	}
}

func Test_mjHandler_handleTenhouWebSocketMessage(t *testing.T) {
	h := &mjHandler{
		tenhouRoundData: newGame(nil),
	}

	mr := tenhou.NewMessageReceiverWithQueueSize(10000)
	mr.SkipOrderCheck = true
	for _, l := range logs {
		if l.Level != "INFO" {
			fmt.Println(l.Message)
			continue
		}
		debug.Lo++
		mr.Put([]byte(l.Message))
		msg := mr.Get()
		h.handleTenhouWebSocketMessage(msg)
	}
}

//func Test_mjHandler_runAnalysisMajsoulMessageTask(t *testing.T) {
//	logFile := "log/gamedata.log"
//	startLo := 33020
//	endLo := 33369
//
//	h := &mjHandler{
//		majsoulUIMessageQueue:  make(chan []byte, 10000),
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
//		h.majsoulUIMessageQueue <- []byte([]byte(s.Message))
//	}
//
//	go h.runAnalysisMajsoulWebSocketMessageTask()
//
//	for {
//		if len(h.majsoulUIMessageQueue) == 0 {
//			break
//		}
//		time.Sleep(time.Second)
//	}
//}
