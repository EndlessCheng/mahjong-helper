package main

import (
	"io/ioutil"
	"encoding/json"
	"bytes"
	)

const (
	configFile = "config.json"
	logFile    = "gamedata.log"
)

type gameConfig struct {
	MajsoulAccountID int `json:"majsoul_account_id"`
}

var gameConf = &gameConfig{
	MajsoulAccountID: -1,
}

func init() {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		if debugMode {
			panic(err)
		}
		return
	}

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(gameConf); err != nil {
		if debugMode {
			panic(err)
		}
		return
	}

	//fmt.Println(*gameConf)
}
