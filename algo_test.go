package main

import (
	"testing"
)

func TestCheckWin(t *testing.T) {
	if !checkWin(convert("111234678m dong dong dong xi xi")) {
		t.Error("checkWin")
	}
}
