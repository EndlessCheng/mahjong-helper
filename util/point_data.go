package util

// 子家荣和点数均值
// 参考：「統計学」のマージャン戦術
// 亲家按 x1.5 算
// TODO: 剩余 dora 数对失点的影响
const (
	RonPointRiichiHiIppatsu = 5172.0 // 基准
	RonPointRiichiIppatsu   = 7445.0
	//RonPointHonitsu         = 6603.0
	//RonPointToitoi          = 7300.0
	RonPointOtherNaki = 3000.0 // *fixed
	// TODO: 考虑双东的影响
	RonPointDama = 4536.0
)

// 简单地判断子家副露者的打点 3000-4200-5880-8232
// 亲家按 x1.5 算
// TODO: 暗杠对打点的提升？
func RonPointOtherNakiWithDora(doraCount int) (point float64) {
	if doraCount > 3 {
		doraCount = 3
	}
	point = RonPointOtherNaki
	const doraMulti = 1.4 // TODO: 待调整？
	for i := 0; i < doraCount; i++ {
		point *= doraMulti
	}
	return point
}
