package util

import "testing"

func TestCalcPointRon(t *testing.T) {
	t.Log(CalcPointRon(3, 40, 0, false) == 5200)
	t.Log(CalcPointRon(3, 40, 0, true) == 7700)
	t.Log(CalcPointRon(4, 40, 0, true) == 12000)
}

func TestCalcPointTsumoSum(t *testing.T) {
	t.Log(CalcPointTsumoSum(3, 40, 0, false) == 5200)
	t.Log(CalcPointTsumoSum(3, 40, 0, true) == 7800)
	t.Log(CalcPointTsumoSum(4, 40, 0, true) == 12000)
}
