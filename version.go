package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/color"
)

const VersionDev = "dev"

// 编译时自动写入版本号
// go build -ldflags "-X main.Version=$(git describe --abbrev=0 --tags)" -o mahjong-helper
var Version = VersionDev

func FetchLatestVersionTag() (latestVersionTag string, err error) {
	const apiGetLatestRelease = "https://api.github.com/repos/EndlessCheng/mahjong-helper/releases/latest"
	const timeout = 10 * time.Second

	c := &http.Client{Timeout: timeout}
	resp, err := c.Get(apiGetLatestRelease)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("[fetchLatestVersionTag] 返回 %s", resp.Status)
	}

	d := struct {
		TagName string `json:"tag_name"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return
	}

	return d.TagName, nil
}

func CheckNewVersion(currentVersionTag string) {
	const latestReleasePage = "https://github.com/EndlessCheng/mahjong-helper/releases/latest"

	latestVersionTag, err := FetchLatestVersionTag()
	if err != nil {
		// 下次再说~
		return
	}

	if latestVersionTag > currentVersionTag {
		color.HiGreen("检测到新版本: %s！请前往 %s 下载", latestVersionTag, latestReleasePage)
	}
}
