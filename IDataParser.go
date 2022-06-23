package main

import "github.com/EndlessCheng/mahjong-helper/util/model"

type DataParser interface {
	// 数据来源（是天凤还是雀魂）
	GetDataSourceType() int

	// 获取自家初始座位：0-第一局的东家 1-第一局的南家 2-第一局的西家 3-第一局的北家
	// 仅处理雀魂数据，天凤返回 -1
	GetSelfSeat() int

	// 原始 JSON
	GetMessage() string

	// 解析前，根据消息内容来决定是否要进行后续解析
	SkipMessage() bool

	// 尝试解析用户名
	IsLogin() bool
	HandleLogin()

	// round 开始/重连
	// roundNumber: 场数（如东1为0，东2为1，...，南1为4，...，南4为7，...），对于三麻来说南1也是4
	// benNumber: 本场数
	// dealer: 庄家 0-3
	// doraIndicators: 宝牌指示牌
	// handTiles: 手牌
	// numRedFives: 按照 mps 的顺序，赤5个数
	IsInit() bool
	ParseInit() (roundNumber int, benNumber int, dealer int, doraIndicators []int, handTiles []int, numRedFives []int)

	// 自家摸牌
	// tile: 0-33
	// isRedFive: 是否为赤5
	// kanDoraIndicator: 摸牌时，若为暗杠摸的岭上牌，则可以翻出杠宝牌指示牌，否则返回 -1（目前恒为 -1，见 IsNewDora）
	IsSelfDraw() bool
	ParseSelfDraw() (tile int, isRedFive bool, kanDoraIndicator int)

	// 舍牌
	// who: 0=自家, 1=下家, 2=对家, 3=上家
	// isTsumogiri: 是否为摸切（who=0 时忽略该值）
	// isReach: 是否为立直宣言（isReach 对于天凤来说恒为 false，见 IsReach）
	// canBeMeld: 是否可以鳴牌（who=0 时忽略该值）
	// kanDoraIndicator: 大明杠/加杠的杠宝牌指示牌，在切牌后出现，没有则返回 -1（天凤恒为-1，见 IsNewDora）
	IsDiscard() bool
	ParseDiscard() (who int, discardTile int, isRedFive bool, isTsumogiri bool, isReach bool, canBeMeld bool, kanDoraIndicator int)

	// 鳴牌（含暗杠、加杠）
	// kanDoraIndicator: 暗杠的杠宝牌指示牌，在他家暗杠时出现，没有则返回 -1（天凤恒为-1，见 IsNewDora）
	IsOpen() bool
	ParseOpen() (who int, meld *model.Meld, kanDoraIndicator int)

	// 立直声明（IsReach 对于雀魂来说恒为 false，见 ParseDiscard）
	IsReach() bool
	ParseReach() (who int)

	// 振听
	IsFuriten() bool

	// 本局是否和牌
	IsRoundWin() bool
	ParseRoundWin() (whos []int, points []int)

	// 是否流局
	// 四风连打 四家立直 四杠散了 九种九牌 三家和了 | 流局听牌 流局未听牌 | 流局满贯
	// 三家和了
	IsRyuukyoku() bool
	ParseRyuukyoku() (type_ int, whos []int, points []int)

	// 拔北宝牌
	IsNukiDora() bool
	ParseNukiDora() (who int, isTsumogiri bool)

	// 这一项放在末尾处理
	// 杠宝牌（雀魂在暗杠后的摸牌时出现）
	// kanDoraIndicator: 0-33
	IsNewDora() bool
	ParseNewDora() (kanDoraIndicator int)
}
