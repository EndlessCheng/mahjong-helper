package main

import (
	"testing"
)

func TestCheckWin(t *testing.T) {
	_, cnt, err := convert("111234678m dong dong dong xi xi")
	if err != nil {
		t.Error(err)
	}
	if !checkWin(cnt) {
		t.Error("checkWin")
	}
}
