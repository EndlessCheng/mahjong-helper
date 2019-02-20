package main

import "testing"

func Test(t *testing.T) {
	for i, risks := range riskData {
		t.Log(i, len(risks))
	}
}
