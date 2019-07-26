package tool

type majsoulResource struct {
	Res struct {
		LiqiJson struct {
			Prefix string `json:"prefix"` // v0.5.143.w
		} `json:"res/proto/liqi.json"`
	} `json:"res"`
}

func getResource(apiGetResourceURL string) (resource *majsoulResource, err error) {
	resource = &majsoulResource{}
	err = get(apiGetResourceURL, resource)
	return
}
