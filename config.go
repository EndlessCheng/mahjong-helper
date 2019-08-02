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

type userConfig struct {
	MajsoulAccountIDs []int `json:"majsoul_account_ids"`

	currentActiveMajsoulAccountID int    `json:"-"`
	currentActiveTenhouUserName   string `json:"-"`
}

var userConf = &userConfig{
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

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(userConf); err != nil {
		if debugMode {
			panic(err)
		}
		return
	}
}

func (c *userConfig) saveConfigToFile() error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(configFile, data, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (c *userConfig) isIDExist(majsoulAccountID int) bool {
	for _, id := range c.MajsoulAccountIDs {
		if id == majsoulAccountID {
			return true
		}
	}
	return false
}

func (c *userConfig) addMajsoulAccountID(majsoulAccountID int) error {
	if c.isIDExist(majsoulAccountID) {
		return nil
	}
	userConf.MajsoulAccountIDs = append(userConf.MajsoulAccountIDs, majsoulAccountID)
	return c.saveConfigToFile()
}

func (c *userConfig) setMajsoulAccountID(majsoulAccountID int) {
	c.currentActiveMajsoulAccountID = majsoulAccountID
}
