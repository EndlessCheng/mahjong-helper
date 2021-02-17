package main

import (
	"testing"
	"github.com/EndlessCheng/mahjong-helper/util"
	"github.com/EndlessCheng/mahjong-helper/util/model"
)

func Test_parseTenhouMeld(t *testing.T) {
	d := &tenhouRoundData{}
	for _, s := range []string{
		"43595",
		"17511",
	} {
		t.Log(d._parseTenhouMeld(s))
	}
}

func TestAnalysisTilesRisk(t *testing.T) {
	debugMode = true

	d := newGame(&tenhouRoundData{})
	handsTiles34, _, err := util.StrToTiles34("123456789m 123456789p 123456789s 1234567z")
	if err != nil {
		t.Fatal(err)
	}
	globalDiscardTiles34, _, err := util.StrToTiles34("22m 158p 123789s 6z")
	if err != nil {
		t.Fatal(err)
	}
	for i, c := range handsTiles34 {
		if c == 0 {
			continue
		}
		d.leftCounts[i] -= c
		if d.leftCounts[c] < 0 {
			t.Fatal("参数有误: ", util.Mahjong[c])
		}
	}
	for i, c := range globalDiscardTiles34 {
		if c == 0 {
			continue
		}
		d.leftCounts[i] -= c
		if d.leftCounts[c] < 0 {
			t.Fatal("参数有误: ", util.Mahjong[c])
		}
		d.globalDiscardTiles = append(d.globalDiscardTiles, i)
	}

	d.players[1].isReached = true
	d.players[1].reachTileAtGlobal = 0
	d.players[1].discardTiles = []int{1, 1, 1, 1, 12, 4, 5, 6}
	d.players[2].isReached = true
	d.players[2].reachTileAtGlobal = 7
	d.players[2].discardTiles = []int{1, 1, 1, 1, 1}
	d.players[3].discardTiles = []int{1, 1, 1, 1, 1}
	d.players[3].melds = make([]*model.Meld, 2)
	d.players[3].meldDiscardsAt = []int{2, 3}
	d.players[3].latestDiscardAtGlobal = 10

	table := d.analysisTilesRisk()
	table.printWithHands(handsTiles34, d.leftCounts)
}

func TestReg(t *testing.T) {
	d := &tenhouRoundData{
		msg: &tenhouMessage{
			Tag: "T123",
		},
	}
	t.Log(d.IsSelfDraw() == true)
	d.msg.Tag = "TATA"
	t.Log(d.IsSelfDraw() == false)
	d.msg.Tag = "T"
	t.Log(d.IsSelfDraw() == false)
	d.msg.Tag = "T1234"
	t.Log(d.IsSelfDraw() == false)

	d.msg.Tag = "D123"
	t.Log(d.IsDiscard() == true)
	d.msg.Tag = "E123"
	t.Log(d.IsDiscard() == true)
	d.msg.Tag = "EAAA"
	t.Log(d.IsDiscard() == false)
	d.msg.Tag = "E"
	t.Log(d.IsDiscard() == false)
	d.msg.Tag = "E123123"
	t.Log(d.IsDiscard() == false)
}
