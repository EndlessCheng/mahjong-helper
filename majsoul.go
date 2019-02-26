package main

type majsoulMessage string

type majsoulRoundData struct {
	roundData
	msg *majsoulMessage
}

func (d *majsoulRoundData) ParseInit() {

}

func (d *majsoulRoundData) ParseDraw() {

}

func (d *majsoulRoundData) ParseDiscard() {

}
