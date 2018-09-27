package main

import (
	"testing"
	"fmt"
)

func TestConvert(t *testing.T) {
	num, cnt, err := convert("11234m 2345567p 45z")
	if err != nil {
		t.Fatal(err)
	}
	if num != 14 {
		t.Fatal(num)
	}
	if fmt.Sprint(cnt) != "[2 1 1 1 0 0 0 0 0 0 1 1 1 2 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 1 0 0]" {
		t.Fatal(cnt)
	}
}

func TestCountToString(t *testing.T) {
	cnt := []int{
		2, 1, 1, 1, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 2, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 1, 0, 0,
	}
	raw := countToString(cnt)
	if raw != "11234m 2345567p 45z" {
		t.Error(raw)
	}
}
