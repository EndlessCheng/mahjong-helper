package updater

import (
	"net/http"
	"time"
	"fmt"
	"encoding/json"
	"github.com/fatih/color"
)

func fetchLatestVersionTag() (latestVersionTag string, err error) {
	const getLatestReleaseAPI = "https://api.github.com/repos/EndlessCheng/mahjong-helper/releases/latest"
	const timeout = 10 * time.Second

	c := &http.Client{Timeout: timeout}
	resp, err := c.Get(getLatestReleaseAPI)
	if err != nil {
		return
	}

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

// TODO: move out
func AlertNewVersion(currentVersionTag string) {
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
