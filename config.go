package main

import (
	"io/ioutil"
	"encoding/json"
	"bytes"
	"os"
)

const (
	configFile = "config.json"
)

type gameConfig struct {
	MajsoulAccountIDs []int `json:"majsoul_account_ids"`

	currentActiveMajsoulAccountID int    `json:"-"`
	currentActiveTenhouUsername   string `json:"-"`
}

var gameConf = &gameConfig{
	MajsoulAccountIDs:             []int{},
	currentActiveMajsoulAccountID: -1,
}

func init() {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return
	}

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

func (c *gameConfig) saveConfigToFile() error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(configFile, data, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (c *gameConfig) isIDExist(majsoulAccountID int) bool {
	for _, id := range c.MajsoulAccountIDs {
		if id == majsoulAccountID {
			return true
		}
	}
	return false
}

func (c *gameConfig) addMajsoulAccountID(majsoulAccountID int) error {
	if c.isIDExist(majsoulAccountID) {
		return nil
	}
	gameConf.MajsoulAccountIDs = append(gameConf.MajsoulAccountIDs, majsoulAccountID)
	return c.saveConfigToFile()
}

func (c *gameConfig) setMajsoulAccountID(majsoulAccountID int) {
	c.currentActiveMajsoulAccountID = majsoulAccountID
}