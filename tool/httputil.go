package tool

import (
	"fmt"
	"github.com/levigross/grequests"
	"io/ioutil"
)

func get(url string, userStruct interface{}) (err error) {
	resp, err := grequests.Get(url, nil)
	if err != nil {
		return
	}

	if !resp.Ok {
		return fmt.Errorf("%s", resp.RawResponse.Status)
	}

	return resp.JSON(userStruct)
}

func fetch(url string) (content []byte, err error) {
	resp, err := grequests.Get(url, nil)
	if err != nil {
		return
	}
	defer resp.Close()

	if !resp.Ok {
		return nil, fmt.Errorf("%s", resp.RawResponse.Status)
	}

	return ioutil.ReadAll(resp.RawResponse.Body)
}
