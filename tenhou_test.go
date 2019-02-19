package main

import "testing"

func Test_parseTenhouMeld(t *testing.T) {
	d := &tenhouRoundData{}
	for _, s := range []string{"43595", "17511"} {
		t.Log(d._parseTenhouMeld(s))
	}
}
