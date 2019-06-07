package main

import (
	"testing"
	"github.com/EndlessCheng/mahjong-helper/util/model"
)

func TestAnalysis(t *testing.T) {
	var raw string
	raw = "11222333789s 11z"
	raw = "2355789p 356778s"
	raw = "4578999m 45p 11145s"
	raw = "123345567m 34p 345s"
	raw = "123m 2378p 34599s 1z"
	raw = "2334567788s 5699p"
	raw = "123m 22378p 345899s"
	raw = "123m 22378p 345899s"
	raw = "1234m 22277p 3456s"
	raw = "123m 2378p 234999s"
	raw = "45689m 1189p 22256s" // 41775557 => 7800426
	raw = "12367m 123667p 556s"
	raw = "12378m 12378p 123s"
	raw = "123m 2378p 34599s 1z" // 5180198 => 416416

	// http://blog.sina.com.cn/s/blog_7f78b76f0100s0nl.html
	raw = "11379m 347p 277s 777z"
	raw = "334578m 11468p 235s"
	raw = "478m 33588p 457899s"
	raw = "2233688m 1234p 378s"
	raw = "1233347m 23699p 88s"

	raw = "56778m 1245s 23388p"
	raw = "23m 22456p 1156899s"
	//raw = "2379m 22399s 23479p"

	raw = "24m 133479p 226778s"

	raw = "1234689m 468p 4699s"
	raw = "13579m 1357p 44789s"

	raw = "4589m 1345677p 458s"
	raw = "233m 11335p 44789s"

	raw = "34568m 678p 13567s"

	raw = "123345m 234p 24s 44z"
	raw = "123345m 23468p 44z"

	raw = "24567m 24456p 229s"

	raw = "3456667m 345566p"
	raw = "3456667m 34566p 5s"
	raw = "24688m 34s # 6666P 234p + 3m"
	raw = "24688m 34s # 111p 234p + 3m"
	raw = "1112234567999m"
	raw = "123567m 3334688p + 7z"
	raw = "23777m 45677s # 777p + 7s" // *片听
	raw = "456789m 1123678p 6z"
	raw = "44779m 889p 78s # 666z + 8p?"
	if _, err := analysisHumanTiles(model.NewSimpleHumanTilesInfo(raw)); err != nil {
		t.Fatal(err)
	}
}
