package main

import "testing"

func TestName(t *testing.T) {
	var raw string
	//raw = "11222333789s fa fa"
	//raw = "2355789p 356778s"
	//raw = "4578999m 45p 11145s"
	//raw = "123345567m 34p 345s"
	//raw = "123m 2378p 34599s bei"
	//raw = "2334567788s 5699p"
	//raw = "123m 22378p 345899s"
	raw = "123m 22 378p 345899s"
	analysis(raw)
}
