package tool

import "testing"

func Test_appendRandv(t *testing.T) {
	for i := 0; i < 30; i++ {
		t.Log(appendRandv(""))
	}
}
