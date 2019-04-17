package main

import (
	"testing"
	"io/ioutil"
	"strings"
	"fmt"
	"encoding/json"
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

	//
	accountID := -1
	startLo := 1

	majsoulRoundData := &majsoulRoundData{accountID: accountID}
	majsoulRoundData.roundData = newRoundData(majsoulRoundData, 0, 0)

	for lo, line := range strings.Split(string(logData), "\n")[startLo-1:] {
		fmt.Println(lo + 1)
		json.Unmarshal([]byte(line), &s)

		msg := s.Message
		d := majsoulMessage{}
		if err := json.Unmarshal([]byte(msg), &d); err != nil {
			fmt.Println(err)
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

	//
	startLo := 9131

	tenhouRoundData := &tenhouRoundData{isRoundEnd: true}
	tenhouRoundData.roundData = newRoundData(tenhouRoundData, 0, 0)

	for lo, line := range strings.Split(string(logData), "\n")[startLo-1:] {
		fmt.Println(lo + 1)
		json.Unmarshal([]byte(line), &s)

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
