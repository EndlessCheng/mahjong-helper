package tool

type majsoulConfig struct {
	IP []struct {
		RegionURLs []struct {
			Mainland string `json:"mainland"` // https://lb-mainland.majsoul.com:2901/api/v0/recommend_list
			HK       string `json:"hk"`       // https://lb-hk.majsoul.com:7891/api/v0/recommend_list
		} `json:"region_urls"`
	} `json:"ip"`
}

func (c *majsoulConfig) apiGetMainlandRecommendListURL() string {
	if len(c.IP) == 0 || len(c.IP[0].RegionURLs) == 0 {
		return ""
	}
	return c.IP[0].RegionURLs[0].Mainland + "?service=ws-gateway&protocol=ws&ssl=true"
}

func (c *majsoulConfig) apiGetMainlandRecommendListURLWithLocation(location string) string {
	if len(c.IP) == 0 || len(c.IP[0].RegionURLs) == 0 {
		return ""
	}
	return c.IP[0].RegionURLs[0].Mainland + "?service=ws-game-gateway&protocol=ws&ssl=true&location=" + location
}

func getConfig(apiGetConfigURL string) (config *majsoulConfig, err error) {
	config = &majsoulConfig{}
	err = get(apiGetConfigURL, config)
	return
}

// {"servers":["mj-srv-7.majsoul.com:4130","mj-srv-7.majsoul.com:4131","mj-srv-7.majsoul.com:4132","mj-srv-7.majsoul.com:4133","mj-srv-5.majsoul.com:4100","mj-srv-5.majsoul.com:4102","mj-srv-5.majsoul.com:4101","mj-srv-5.majsoul.com:4103"]}
type recommendServers struct {
	Servers []string `json:"servers"`
}
