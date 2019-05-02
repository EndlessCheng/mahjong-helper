package util

import "github.com/EndlessCheng/mahjong-helper/util/model"

type _handInfo struct {
	*model.PlayerInfo
	divideResult  *DivideResult // 手牌解析结果
	_containHonor *bool
	_isNaki       *bool
}
