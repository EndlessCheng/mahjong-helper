package main

import "testing"

func Test_checkNewVersion(t *testing.T) {
	newVersionTag, err := checkNewVersion()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(newVersionTag)
}
