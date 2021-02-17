package main

const (
	majsoulGameConfigCategoryFriends = 1 // 友人
	majsoulGameConfigCategoryMatch   = 2 // 段位 比赛
)

// 古役 {\"category\":2,\"mode\":{\"mode\":1,\"detail_rule\":{\"guyi_mode\":1}}
// 非古役{\"category\":2,\"mode\":{\"mode\":1}
type majsoulGameConfig struct {
	Category int `json:"category"`
	Mode     *struct {
		Mode       int `json:"mode"`
		DetailRule *struct {
			GuyiMode int `json:"guyi_mode"`
		} `json:"detail_rule"`
	} `json:"mode"`
}

func (c *majsoulGameConfig) isGuyiMode() bool {
	return c != nil && c.Mode != nil && c.Mode.DetailRule != nil && c.Mode.DetailRule.GuyiMode == 1
}
