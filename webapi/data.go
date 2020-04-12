package webapi

type ApiData struct {
	// 数据更新时间戳
	Timestamp int `json:"timestamp"`

	// 自家手牌 一个长度为 34 的整数数组
	Counts []int `json:"counts"`

	// 手牌危险度 一个长度为 34 的浮点数组
	RiskTable []float64 `json:"risk"`
}


