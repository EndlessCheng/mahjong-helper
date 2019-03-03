package util

import "testing"

func TestCalculateAgariRate(t *testing.T) {
	t.Log(CalculateAgariRate(Waits{0: 4, 3: 4}, nil))
	t.Log(CalculateAgariRate(Waits{0: 2, 9: 2}, nil))
	t.Log(CalculateAgariRate(Waits{0: 3, 1: 4}, nil))
	t.Log(CalculateAgariRate(Waits{8: 3, 7: 4}, nil))
	t.Log(CalculateAgariRate(Waits{0: 1, 1: 3, 2: 3, 3: 3, 4: 3, 5: 3, 6: 3, 7: 3, 9: 1}, nil))
	t.Log(CalculateAgariRate(Waits{9: 2, 27: 2}, nil))
	t.Log(CalculateAgariRate(Waits{27: 3}, nil))
	t.Log(CalculateAgariRate(Waits{27: 2}, nil))
	t.Log(CalculateAgariRate(Waits{27: 1}, nil))
	t.Log(CalculateAgariRate(Waits{27: 0}, nil))
}
