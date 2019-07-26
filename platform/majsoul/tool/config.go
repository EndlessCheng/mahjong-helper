package tool

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type majsoulConfig struct {
	IP []struct {
		RegionURLs struct {
			Mainland string `json:"mainland"` // https://lb-mainland.majsoul.com:2901/api/v0/recommend_list
			HK       string `json:"hk"`       // https://lb-hk.majsoul.com:7891/api/v0/recommend_list
		} `json:"region_urls"`
	} `json:"ip"`
}

func (c *majsoulConfig) apiGetMainlandRecommendListURL() string {
	if len(c.IP) == 0 {
		return ""
	}
	return c.IP[0].RegionURLs.Mainland + "?service=ws-gateway&protocol=ws&ssl=true"
}

func (c *majsoulConfig) apiGetMainlandRecommendListURLWithLocation(location string) string {
	if len(c.IP) == 0 {
		return ""
	}
	return c.IP[0].RegionURLs.Mainland + "?service=ws-game-gateway&protocol=ws&ssl=true&location=" + location
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

func getRecommendServers(apiGetRecommendServersURL string) (servers []string, err error) {
	recommendServers := recommendServers{}
	if err = get(apiGetRecommendServersURL, &recommendServers); err != nil {
		return
	}
	return recommendServers.Servers, nil
}

// 获取雀魂 WebSocket 服务器地址
func GetMajsoulWebSocketURL() (url string, err error) {
	version, err := GetMajsoulVersion(ApiGetVersionZH)
	if err != nil {
		return
	}

	apiGetConfigURL := fmt.Sprintf(apiGetConfigFormatZH, version.ResVersion)
	config, err := getConfig(apiGetConfigURL)
	if err != nil {
		return
	}

	apiGetMainlandRecommendListURL := config.apiGetMainlandRecommendListURL()
	servers, err := getRecommendServers(apiGetMainlandRecommendListURL)
	if err != nil {
		return
	}
	if len(servers) == 0 {
		return "", fmt.Errorf("维护中，没有可用的服务器地址")
	}

	// 随机取一个
	host := servers[rand.Intn(len(servers))]
	url = "wss://" + host + "/"
	return
}
