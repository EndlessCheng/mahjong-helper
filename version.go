package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/color"
)

// declare const
const VersionDev string = "dev"

// 编译时自动写入版本号
// go build -ldflags "-X main.Version=$(git describe --abbrev=0 --tags)" -o mahjong-helper
// declare Version
var Version string = VersionDev

// if check github has new version and remind update
func CheckNewVersion(currentVersionTag string) {
	// target to github
	const latestReleasePage = "https://github.com/EndlessCheng/mahjong-helper/releases/latest"

	latestVersionTag, err := fetchLatestVersionTag()
	if err != nil {
		// 下次再说~
		return
	}

	if latestVersionTag > currentVersionTag {
		color.HiGreen("检测到新版本: %s！请前往 %s 下载", latestVersionTag, latestReleasePage)
	}
}

// fetch the lastest version
func fetchLatestVersionTag() (latestVersionTag string, err error) {
	const apiGetLatestRelease string = "https://api.github.com/repos/EndlessCheng/mahjong-helper/releases/latest"
	const timeout time.Duration = 10 * time.Second

	// declare a client from http.Client and set Timeout to timeout
	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(apiGetLatestRelease)
	// if err unqual to nil mean get error
	if err != nil {
		return
	}
	// deffered close 
	defer resp.Body.Close()

	// if response is not OK
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("[fetchLatestVersionTag] 返回 %s", resp.Status)
	}

	//anonymous struct
	data := struct {
		TagName string `json:"tag_name"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	return data.TagName, nil
}
