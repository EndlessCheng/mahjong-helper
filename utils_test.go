package main

import (
	"testing"
)

func TestCountToString(t *testing.T) {
	raw := countToString([]int{
		2, 1, 1, 1, 0, 0, 0, 0, 0,
		0, 1, 1, 1, 2, 1, 1, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1,
		1, 0, 0,
	})
	if raw != "11234m 2345567p bei zhong" {
		t.Error(raw)
	}
}
