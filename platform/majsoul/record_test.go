package majsoul

import (
	"os"
	"testing"
)

func TestDownloadRecords(t *testing.T) {
	username, ok := os.LookupEnv("USERNAME")
	if !ok {
		t.Skip("未配置环境变量 USERNAME，退出")
	}

	password, ok := os.LookupEnv("PASSWORD")
	if !ok {
		t.Skip("未配置环境变量 PASSWORD，退出")
	}

	if err := DownloadRecords(username, password, RecordTypeAll); err != nil {
		t.Fatal(err)
	}
}
