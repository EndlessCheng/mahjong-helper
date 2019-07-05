package handler

import "testing"

func Test_checkNewVersion(t *testing.T) {
	latestVersionTag, err := fetchLatestVersionTag()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(latestVersionTag)
}
