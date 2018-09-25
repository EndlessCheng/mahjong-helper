package main

import (
	"testing"
)

func TestCheckWin(t *testing.T) {
	_, cnt := convert("111234678m dong dong dong xi xi")
	if !checkWin(cnt) {
		t.Error("checkWin")
	}
}
