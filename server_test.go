package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/EndlessCheng/mahjong-helper/util/debug"
)

func Test_mjHandler_runAnalysisMajsoulMessageTask(t *testing.T) {
	DebugMode = true

	logFile := "log/gamedata.log"
	startLo := 33020
	endLo := 33369

	h := &MahJongHandler{
		majsoulMessageQueue:  make(chan []byte, 10000),
		MahJongSoulRoundData: &MahJongSoulRoundData{},
		MahJongSoulRecordMap: map[string]*MahJongSoulRecordBaseInfo{},
	}
	h.MahJongSoulRoundData.RoundData = NewGame(h.MahJongSoulRoundData)

	s := struct {
		Level   string `json:"level"`
		Message string `json:"message"`
	}{}
	logData, err := ioutil.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(logData), "\n")
	for i, line := range lines[startLo-1 : endLo] {
		debug.Lo = i + 1
		if line == "" {
			continue
		}
		if err := json.Unmarshal([]byte(line), &s); err != nil {
			fmt.Println(err)
			continue
		}
		if s.Level != "INFO" {
			fmt.Println(s.Level, s.Message)
			//t.Fatal(s.Message)
			break
		}
		h.majsoulMessageQueue <- []byte([]byte(s.Message))
	}

	go h.RunAnalysisMahJongSoulMessageTask()

	for {
		if len(h.majsoulMessageQueue) == 0 {
			break
		}
		time.Sleep(time.Second)
	}
}
