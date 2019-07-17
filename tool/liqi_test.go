package tool

import "testing"

func TestFetchLatestLiqiJson(t *testing.T) {
	if err := FetchLatestLiqiJson("../platform/majsoul/proto/lq/liqi.json"); err != nil {
		t.Fatal(err)
	}
}
