package main

import (
	"testing"
)

func TestCheckWin(t *testing.T) {
	_, cnt, err := convert("111234678m 11122z")
	if err != nil {
		t.Error(err)
	}
	if !checkWin(cnt) {
		t.Error("checkWin")
	}
}
