package tool

import (
	"math/rand"
	"strconv"
)

// for WebSocket
const (
	MajsoulOriginURL = "https://www.majsoul.com"
)

// for normal requests
const (
	majsoulJSURLPrefixZH = "https://www.majsoul.com/1/"
	majsoulJSURLPrefixEN = "https://mahjongsoul.game.yo-star.com/"
	majsoulJSURLPrefixJP = "https://game.mahjongsoul.com/"

	ApiGetVersionZH = majsoulJSURLPrefixZH + "version.json"

	apiGetResVersionFormatZH = majsoulJSURLPrefixZH + "resversion%s.json"
	apiGetConfigFormatZH     = majsoulJSURLPrefixZH + "v%s/config.json"
	apiGetLiqiJsonFormatZH   = majsoulJSURLPrefixZH + "%s/res/proto/liqi.json"
)

func appendRandv(apiGetVersionURL string) string {
	randVal1 := rand.Intn(1e9)
	randVal2 := rand.Intn(1e9)
	return apiGetVersionURL + "?randv" + strconv.Itoa(randVal1) + strconv.Itoa(randVal2)
}
