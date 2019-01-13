package main

import "testing"

func Test_calcAgariRate(t *testing.T) {
	t.Log(calcAgariRate(needTiles{0:3}))
	t.Log(calcAgariRate(needTiles{0:3,1:4}))
}
