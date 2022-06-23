package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/EndlessCheng/mahjong-helper/util/debug"
	"github.com/stretchr/testify/assert"
)

func Test_majsoul_analysis(t *testing.T) {
	DebugMode = true

	logFile := "log/gamedata-x.log"
	logData, err := ioutil.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	s := struct {
		Level   string `json:"level"`
		Message string `json:"message"`
	}{}

	// config
	startLo := -1
	endLo := -1

	majsoulRoundData := &MahJongSoulRoundData{}
	majsoulRoundData.RoundData = NewGame(majsoulRoundData)

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

		if s.Level != "INFO" {
			fmt.Println(s.Level, s.Message)
			continue
		}

		msg := s.Message
		d := MaJongSoulMessage{}
		if err := json.Unmarshal([]byte(msg), &d); err != nil {
			fmt.Println(err)
			continue
		}

		majsoulRoundData.Message = &d
		majsoulRoundData.OriginJSON = msg
		if err := majsoulRoundData.Analysis(); err != nil {
			fmt.Println("错误：", err)
		}
	}
}

func Test_tenhou_analysis(t *testing.T) {
	DebugMode = true

	logFile := "log/gamedata-20190715-201349.log"
	logData, err := ioutil.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	s := struct {
		Level   string `json:"level"`
		Message string `json:"message"`
	}{}

	// config
	startLo := -1
	endLo := -1

	tenhouRoundData := &TenHouRoundData{IsRoundEnd: true}
	tenhouRoundData.RoundData = NewGame(tenhouRoundData)

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

		if s.Level != "INFO" {
			fmt.Println(s.Message)
			continue
		}

		msg := s.Message
		d := TenhouMessage{}
		if err := json.Unmarshal([]byte(msg), &d); err != nil {
			fmt.Println(err)
			continue
		}

		tenhouRoundData.Msg = &d
		tenhouRoundData.OriginJSON = msg
		if err := tenhouRoundData.Analysis(); err != nil {
			fmt.Println("错误：", err)
		}
	}
}

func Test_modifySanninPlayerInfoList(t *testing.T) {
	assert := assert.New(t)

	roundNumber := 0
	dealer := 2
	rd := newRoundData(nil, roundNumber, 0, dealer)
	newPlayers := modifySanninPlayerInfoList(rd.players, roundNumber)
	assert.Equal(newPlayers[0].selfWindTile, 29)
	assert.Equal(newPlayers[1].selfWindTile, 30)
	assert.Equal(newPlayers[2].selfWindTile, 27)
	assert.Equal(newPlayers[3].selfWindTile, 28)

	roundNumber = 1
	dealer = 3
	rd = newRoundData(nil, roundNumber, 0, dealer)
	newPlayers = modifySanninPlayerInfoList(rd.players, roundNumber)
	assert.Equal(newPlayers[0].selfWindTile, 28)
	assert.Equal(newPlayers[1].selfWindTile, 30)
	assert.Equal(newPlayers[2].selfWindTile, 29)
	assert.Equal(newPlayers[3].selfWindTile, 27)

	roundNumber = 2
	dealer = 0
	rd = newRoundData(nil, roundNumber, 0, dealer)
	newPlayers = modifySanninPlayerInfoList(rd.players, roundNumber)
	assert.Equal(newPlayers[0].selfWindTile, 27)
	assert.Equal(newPlayers[1].selfWindTile, 30)
	assert.Equal(newPlayers[2].selfWindTile, 28)
	assert.Equal(newPlayers[3].selfWindTile, 29)
}
