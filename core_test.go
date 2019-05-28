package main

import (
	"testing"
	"io/ioutil"
	"strings"
	"fmt"
	"encoding/json"
	"github.com/EndlessCheng/mahjong-helper/util/debug"
)

func Test_majsoul_analysis(t *testing.T) {
	debugMode = true

	logData, err := ioutil.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	s := struct {
		Message string `json:"message"`
	}{}

	// config
	accountID := gameConf.MajsoulAccountID
	startLo := -1
	endLo := -1

	majsoulRoundData := &majsoulRoundData{accountID: accountID}
	majsoulRoundData.roundData = newRoundData(majsoulRoundData, 0, 0)

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

	for lo, line := range lines[startLo-1 : endLo] {
		debug.Lo = lo + 1
		fmt.Println(debug.Lo)
		if line == "" {
			continue
		}

		if err := json.Unmarshal([]byte(line), &s); err != nil {
			fmt.Println(err)
			continue
		}

		msg := s.Message
		d := majsoulMessage{}
		if err := json.Unmarshal([]byte(msg), &d); err != nil {
			fmt.Println(err)
			continue
		}

		if d.AccountID > 0 {
			majsoulRoundData.accountID = d.AccountID
			printAccountInfo(d.AccountID)
			continue
		}

		majsoulRoundData.msg = &d
		majsoulRoundData.originJSON = msg
		if err := majsoulRoundData.analysis(); err != nil {
			fmt.Println("错误：", err)
		}
	}
}

func Test_tenhou_analysis(t *testing.T) {
	debugMode = true

	logData, err := ioutil.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	s := struct {
		Message string `json:"message"`
	}{}

	// config
	startLo := -1
	endLo := -1

	tenhouRoundData := &tenhouRoundData{isRoundEnd: true}
	tenhouRoundData.roundData = newRoundData(tenhouRoundData, 0, 0)

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

	for lo, line := range lines[startLo-1 : endLo] {
		debug.Lo = lo + 1
		fmt.Println(debug.Lo)
		if line == "" {
			continue
		}

		if err := json.Unmarshal([]byte(line), &s); err != nil {
			fmt.Println(err)
			continue
		}

		msg := s.Message
		d := tenhouMessage{}
		if err := json.Unmarshal([]byte(msg), &d); err != nil {
			fmt.Println(err)
			continue
		}

		tenhouRoundData.msg = &d
		tenhouRoundData.originJSON = msg
		if err := tenhouRoundData.analysis(); err != nil {
			fmt.Println("错误：", err)
		}
	}
}
