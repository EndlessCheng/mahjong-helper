package tool

type majsoulVersion struct {
	Code       string `json:"code"`    // code.js 路径 v0.5.81.w/code.js
	ResVersion string `json:"version"` // 资源文件版本  0.5.82.w（注意开头没有 v）
}

func GetMajsoulVersion(apiGetVersionURL string) (version *majsoulVersion, err error) {
	version = &majsoulVersion{}
	err = get(apiGetVersionURL, version)
	return
}
