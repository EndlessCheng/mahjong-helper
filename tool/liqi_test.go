package tool

import "testing"

func TestFetchLatestLiqiJson(t *testing.T) {
	if err := FetchLatestLiqiJson("../platform/majsoul/proto/lq/liqi.json"); err != nil {
		t.Fatal(err)
	}
}

func TestLiqiJsonToProto3(t *testing.T) {
	content, err := fetchLatestLiqiJson()
	if err != nil {
		t.Fatal(err)
	}
	if err := LiqiJsonToProto3(content, "../platform/majsoul/proto/lq/liqi.proto"); err != nil {
		t.Fatal(err)
	}
}
