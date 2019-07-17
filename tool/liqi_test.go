package tool

import "testing"

func TestFetchLatestLiqiJson(t *testing.T) {
	if err := FetchLatestLiqiJson("../res/liqi.json"); err != nil {
		t.Fatal(err)
	}
}
