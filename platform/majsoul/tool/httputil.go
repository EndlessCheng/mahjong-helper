package tool

import (
	"fmt"
	"github.com/levigross/grequests"
	"io/ioutil"
)

func newReqOpt() *grequests.RequestOptions {
	return &grequests.RequestOptions{
		UserAgent: "Mozilla/5.0 AppleWebKit/530.00 (KHTML, like Gecko) Chrome/75.0.3120.123 Safari/531.21",
	}
}

func get(url string, userStruct interface{}) (err error) {
	resp, err := grequests.Get(url, newReqOpt())
	if err != nil {
		return
	}
	defer resp.Close()
	if !resp.Ok {
		return fmt.Errorf("%s", resp.RawResponse.Status)
	}
	return resp.JSON(userStruct)
}

func Fetch(url string) (content []byte, err error) {
	resp, err := grequests.Get(url, newReqOpt())
	if err != nil {
		return
	}
	defer resp.Close()
	if !resp.Ok {
		return nil, fmt.Errorf("%s", resp.RawResponse.Status)
	}
	return ioutil.ReadAll(resp.RawResponse.Body)
}
