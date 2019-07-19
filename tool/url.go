package tool

import (
	"math/rand"
	"strconv"
)

const (
	majsoulJSURLPrefixZH = "https://majsoul.union-game.com/0/"
	majsoulJSURLPrefixEN = "https://mahjongsoul.game.yo-star.com/"
	majsoulJSURLPrefixJP = "https://game.mahjongsoul.com/"

	apiGetVersionZH = majsoulJSURLPrefixZH + "version.json"
	apiGetVersionEN = majsoulJSURLPrefixEN + "version.json"
	apiGetVersionJP = majsoulJSURLPrefixJP + "version.json"

	apiGetResVersionFormatZH = majsoulJSURLPrefixZH + "resversion%s.json"
	apiGetConfigFormatZH     = majsoulJSURLPrefixZH + "%s/config.json"
	apiGetLiqiJsonFormatZH   = majsoulJSURLPrefixZH + "%s/res/proto/liqi.json"
)

func appendRandv(apiGetVersionURL string) string {
	rand1 := rand.Intn(1e9)
	rand2 := rand.Intn(1e9)
	return apiGetVersionURL + "?randv" + strconv.Itoa(rand1) + strconv.Itoa(rand2)
}
