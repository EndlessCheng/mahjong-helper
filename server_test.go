package main

import (
	"testing"
	"time"
	"fmt"
	"encoding/json"
	"strings"
	"io/ioutil"
	"github.com/EndlessCheng/mahjong-helper/util/debug"
)

func Test_mjHandler_runAnalysisTenhouMessageTask(t *testing.T) {
	debugMode = true

	h := &mjHandler{
		tenhouMessageQueue: make(chan []byte, 1000),
		tenhouRoundData:    &tenhouRoundData{isRoundEnd: true},
	}
	h.tenhouRoundData.roundData = newGame(h.tenhouRoundData)

	h.tenhouMessageQueue <- []byte("{\"tag\":\"T8\"}")
	h.tenhouMessageQueue <- []byte("{\"tag\":\"INIT\",\"seed\":\"1,1,0,2,0,27\",\"ten\":\"318,224,224,234\",\"oya\":\"0\",\"hai\":\"129,90,47,39,4,9,116,53,33,123,69,28,14\"}")

	go h.runAnalysisTenhouMessageTask()

	for {
		if len(h.tenhouMessageQueue) == 0 {
			break
		}
		time.Sleep(time.Second)
	}
}

func Test_mjHandler_runAnalysisMajsoulMessageTask(t *testing.T) {
	debugMode = true

	startLo := 33020
	endLo := 33369

	h := &mjHandler{
		majsoulMessageQueue: make(chan []byte, 10000),
		majsoulRoundData:    &majsoulRoundData{},
		majsoulRecordMap:    map[string]*majsoulRecordBaseInfo{},
	}
	h.majsoulRoundData.roundData = newGame(h.majsoulRoundData)

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

	go h.runAnalysisMajsoulMessageTask()

	for {
		if len(h.majsoulMessageQueue) == 0 {
			break
		}
		time.Sleep(time.Second)
	}
}
