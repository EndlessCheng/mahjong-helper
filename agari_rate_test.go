package main

import "testing"

func Test_calcAgariRate(t *testing.T) {
	t.Log(calcAgariRate(needTiles{0: 3}, nil))
	t.Log(calcAgariRate(needTiles{0: 3, 1: 4}, nil))
	t.Log(calcAgariRate(needTiles{8: 3, 7: 4}, nil))
	t.Log(calcAgariRate(needTiles{0: 1, 1: 3, 2: 3, 3: 3, 4: 3, 5: 3, 6: 3, 7: 3, 9: 1}, nil))
	t.Log(calcAgariRate(needTiles{9: 2, 27: 2}, nil))
	t.Log(calcAgariRate(needTiles{27: 3}, nil))
	t.Log(calcAgariRate(needTiles{27: 2}, nil))
	t.Log(calcAgariRate(needTiles{27: 1}, nil))
	t.Log(calcAgariRate(needTiles{27: 0}, nil))
}
