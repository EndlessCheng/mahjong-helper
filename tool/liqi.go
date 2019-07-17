package tool

import (
	"fmt"
	"io/ioutil"
)

func fetchLatestLiqiJson() (jsonContent []byte, err error) {
	apiGetVersionURL := appendRandv(apiGetVersionZH)
	version, err := getVersion(apiGetVersionURL)
	if err != nil {
		return
	}

	apiGetResJsonURL := fmt.Sprintf(apiGetResVersionFormatZH, version.ResVersion)
	resource, err := getResource(apiGetResJsonURL)
	if err != nil {
		return
	}

	apiGetLiqiJsonURL := fmt.Sprintf(apiGetLiqiJsonFormatZH, resource.Res.LiqiJson.Prefix)
	return fetch(apiGetLiqiJsonURL)
}

func FetchLatestLiqiJson(filePath string) (err error) {
	jsonContent, err := fetchLatestLiqiJson()
	if err != nil {
		return
	}
	return ioutil.WriteFile(filePath, jsonContent, 0644)
}
