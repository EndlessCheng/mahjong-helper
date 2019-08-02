package majsoul

import (
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"github.com/EndlessCheng/mahjong-helper/platform/common"
	"github.com/EndlessCheng/mahjong-helper/platform/majsoul/proto/lq"
	"github.com/golang/protobuf/proto"
	"github.com/EndlessCheng/mahjong-helper/util"
	"sort"
	"fmt"
)

type Action struct {
	parsedAction proto.Message

	// 用于辅助计算
	CurrentDealer         int
	CurrentRoundNumber    int
	CurrentDoraIndicators []int
}

func NewAction(ap *lq.ActionPrototype) (*Action, error) {
	parsedAction, err := ap.ParseData()
	if err != nil {
		return nil, err
	}
	return &Action{parsedAction: parsedAction}, nil
}

func (*Action) mustParseMajsoulTile(majsoulTile string) (tile34 int, isRedFive bool) {
	tile34, isRedFive, err := util.StrToTile34(majsoulTile)
	if err != nil {
		panic(err)
	}
	return
}

func (a *Action) mustParseMajsoulTiles(majsoulTiles []string) (tiles []int, numRedFive int) {
	tiles = make([]int, len(majsoulTiles))
	for i, majsoulTile := range majsoulTiles {
		tile, isRedFive := a.mustParseMajsoulTile(majsoulTile)
		tiles[i] = tile
		if isRedFive {
			numRedFive++
		}
	}
	return
}

func (a *Action) parseWho(seat uint32) int {
	// 转换成 0=自家, 1=下家, 2=对家, 3=上家
	// 对三麻四麻均适用
	return (int(seat) + a.CurrentDealer - a.CurrentRoundNumber%4 + 4) % 4
}

func (a *Action) isNewDora(majsoulDoras []string) bool {
	return len(majsoulDoras) > len(a.CurrentDoraIndicators)
}

func (a *Action) parseKanDoraIndicator(doras []string) (kanDoraIndicator int) {
	if a.isNewDora(doras) {
		kanDoraIndicator, _ = a.mustParseMajsoulTile(doras[len(doras)-1])
	} else {
		kanDoraIndicator = -1
	}
	return
}

func (a *Action) GetDataSourceType() int {
	return common.DataSourceTypeMajsoul
}

func (a *Action) IsInit() bool {
	switch a.parsedAction.(type) {
	case *lq.ActionNewRound, *lq.RecordNewRound:
		return true
	}
	return false
}

// 吐槽：雀魂就不能把相同的信息重构一下吗
func (a *Action) ParseInit() (roundNumber int, benNumber int, dealer int, doraIndicator int, handTiles []int, numRedFives []int) {
	var chang, ju, ben uint32
	var majsoulDora string
	var majsoulTiles []string
	switch msg := a.parsedAction.(type) {
	case *lq.ActionNewRound:
		chang, ju, ben, majsoulDora = msg.Chang, msg.Ju, msg.Ben, msg.Dora
		majsoulTiles = msg.Tiles
	case *lq.RecordNewRound:
		chang, ju, ben, majsoulDora = msg.Chang, msg.Ju, msg.Ben, msg.Dora
		// TODO: 添加 selfSeat 参数
		majsoulTiles = [][]string{msg.Tiles0, msg.Tiles1, msg.Tiles2, msg.Tiles3}[0]
	}

	roundNumber = int(4*(chang) + ju)
	benNumber = int(ben)
	dealer = -1 // TODO
	doraIndicator, _ = a.mustParseMajsoulTile(majsoulDora)
	numRedFives = make([]int, 3)
	for _, majsoulTile := range majsoulTiles {
		tile, isRedFive := a.mustParseMajsoulTile(majsoulTile)
		handTiles = append(handTiles, tile)
		if isRedFive {
			numRedFives[tile/9]++
		}
	}
	return
}

func (a *Action) IsSelfDraw() bool {
	switch msg := a.parsedAction.(type) {
	case *lq.ActionDealTile:
		return a.parseWho(msg.Seat) == 0
	case *lq.RecordDealTile:
		return a.parseWho(msg.Seat) == 0
	}
	return false
}

func (a *Action) ParseSelfDraw() (tile int, isRedFive bool, kanDoraIndicator int) {
	var majsoulTile string
	var doras []string
	switch msg := a.parsedAction.(type) {
	case *lq.ActionDealTile:
		majsoulTile, doras = msg.Tile, msg.Doras
	case *lq.RecordDealTile:
		majsoulTile, doras = msg.Tile, msg.Doras
	}

	tile, isRedFive = a.mustParseMajsoulTile(majsoulTile)
	kanDoraIndicator = a.parseKanDoraIndicator(doras)
	return
}

func (a *Action) IsDiscard() bool {
	switch a.parsedAction.(type) {
	case *lq.ActionDiscardTile, *lq.RecordDiscardTile:
		return true
	}
	return false
}

func (a *Action) ParseDiscard() (who int, discardTile int, isRedFive bool, isTsumogiri bool, isReach bool, canBeMeld bool, kanDoraIndicator int) {
	var seat uint32
	var majsoulTile string
	var isWLiqi bool
	var ops *lq.OptionalOperationList
	var doras []string
	switch msg := a.parsedAction.(type) {
	case *lq.ActionDiscardTile:
		seat, majsoulTile, isTsumogiri, isReach, isWLiqi, doras = msg.Seat, msg.Tile, msg.Moqie, msg.IsLiqi, msg.IsWliqi, msg.Doras
		ops = msg.Operation
	case *lq.RecordDiscardTile:
		seat, majsoulTile, isTsumogiri, isReach, isWLiqi, doras = msg.Seat, msg.Tile, msg.Moqie, msg.IsLiqi, msg.IsWliqi, msg.Doras
	}

	who = a.parseWho(seat)
	discardTile, isRedFive = a.mustParseMajsoulTile(majsoulTile)
	if isWLiqi {
		isReach = true
	}
	canBeMeld = len(ops.GetOperationList()) > 0 // TODO check 注意：观战模式下无此选项
	kanDoraIndicator = a.parseKanDoraIndicator(doras)
	return
}

func (a *Action) IsOpen() bool {
	switch a.parsedAction.(type) {
	case *lq.ActionChiPengGang, *lq.RecordChiPengGang,
	*lq.ActionAnGangAddGang, *lq.RecordAnGangAddGang:
		return true
	}
	return false
}

const (
	majsoulMeldTypeChi uint32 = iota
	majsoulMeldTypePon
	majsoulMeldTypeMinkanOrKakan
	majsoulMeldTypeAnkan
)

func (a *Action) ParseOpen() (who int, meld *model.Meld, kanDoraIndicator int) {
	var seat, majsoulMeldType uint32
	var majsoulTiles, doras []string
	var froms []uint32
	switch msg := a.parsedAction.(type) {
	case *lq.ActionChiPengGang: // 吃、碰、大明杠
		seat, majsoulMeldType, majsoulTiles, froms = msg.Seat, msg.Type, msg.Tiles, msg.Froms
	case *lq.RecordChiPengGang:
		seat, majsoulMeldType, majsoulTiles, froms = msg.Seat, msg.Type, msg.Tiles, msg.Froms
	case *lq.ActionAnGangAddGang: // 暗杠、加杠
		seat, majsoulMeldType, majsoulTiles, doras = msg.Seat, msg.Type, []string{msg.Tiles}, msg.Doras
	case *lq.RecordAnGangAddGang:
		seat, majsoulMeldType, majsoulTiles, doras = msg.Seat, msg.Type, []string{msg.Tiles}, msg.Doras
	}

	who = a.parseWho(seat)
	kanDoraIndicator = a.parseKanDoraIndicator(doras)

	var meldType int
	switch majsoulMeldType {
	case majsoulMeldTypeChi:
		meldType = common.MeldTypeChi // 吃
	case majsoulMeldTypePon:
		meldType = common.MeldTypePon // 碰
	case majsoulMeldTypeMinkanOrKakan:
		if len(froms) > 0 {
			meldType = common.MeldTypeMinkan // 大明杠
		} else {
			meldType = common.MeldTypeMinkan // 加杠
		}
	case majsoulMeldTypeAnkan:
		meldType = common.MeldTypeAnkan // 暗杠
	default:
		panic(fmt.Sprintf("鸣牌数据解析失败: 无法识别的 majsoulMeldType %d", majsoulMeldType))
	}

	if len(majsoulTiles) > 1 { // 吃、碰、大明杠
		meldTiles, numRedFive := a.mustParseMajsoulTiles(majsoulTiles)
		if meldType == common.MeldTypeChi {
			sort.Ints(meldTiles)
		}

		var majsoulCalledTile string
		for i, seat := range froms {
			fromWho := a.parseWho(seat)
			if fromWho != who {
				majsoulCalledTile = majsoulTiles[i]
			}
		}
		if majsoulCalledTile == "" {
			panic("鸣牌数据解析失败: 未找到 majsoulCalledTile")
		}
		calledTile, redFiveFromOthers := a.mustParseMajsoulTile(majsoulCalledTile)

		meld = &model.Meld{
			MeldType:          meldType,
			Tiles:             meldTiles,
			CalledTile:        calledTile,
			ContainRedFive:    numRedFive > 0,
			RedFiveFromOthers: redFiveFromOthers,
		}
	} else { // 暗杠、加杠
		majsoulTiles = []string{majsoulTiles[0], majsoulTiles[0], majsoulTiles[0], majsoulTiles[0]}
		meldTiles, numRedFive := a.mustParseMajsoulTiles(majsoulTiles)
		meld = &model.Meld{
			MeldType:       meldType,
			Tiles:          meldTiles,
			CalledTile:     meldTiles[0],
			ContainRedFive: numRedFive > 0 || meldTiles[0] < 27 && meldTiles[0]%9 == 4, // 杠5也意味着有赤5
		}
	}
	return
}

func (a *Action) IsRiichi() bool {
	return false
}

func (a *Action) ParseRiichi() (who int) {
	panic("should not be happen")
}

func (a *Action) IsRoundWin() bool {
	switch a.parsedAction.(type) {
	case *lq.ActionHule, *lq.RecordHule:
		return true
	}
	return false
}

func (a *Action) ParseRoundWin() (whos []int, points []int) {
	var hules []*lq.HuleInfo
	switch msg := a.parsedAction.(type) {
	case *lq.ActionHule:
		hules = msg.Hules
	case *lq.RecordHule:
		hules = msg.Hules
	}

	for _, result := range hules {
		whos = append(whos, a.parseWho(result.Seat))
		points = append(points, int(result.PointSum)) // TODO: check PointSum
	}
	return
}

func (a *Action) IsRyuukyoku() bool {
	switch a.parsedAction.(type) {
	case *lq.ActionLiuJu, *lq.RecordLiuJu:
		return true
	}
	return false
}

func (a *Action) ParseRyuukyoku() (type_ int, whos []int, points []int) {
	// TODO
	return
}

func (a *Action) IsNukiDora() bool {
	switch a.parsedAction.(type) {
	case *lq.ActionBaBei, *lq.RecordBaBei:
		return true
	}
	return false
}

func (a *Action) ParseNukiDora() (who int, isTsumogiri bool) {
	var seat uint32
	switch msg := a.parsedAction.(type) {
	case *lq.ActionBaBei:
		seat, isTsumogiri = msg.Seat, msg.Moqie
	case *lq.RecordBaBei:
		seat, isTsumogiri = msg.Seat, msg.Moqie
	}

	return a.parseWho(seat), isTsumogiri
}

func (a *Action) IsNewDora() bool {
	// 非自家摸牌时可能有新的宝牌产生（暗杠）
	var majsoulDoras []string
	switch msg := a.parsedAction.(type) {
	case *lq.ActionDealTile:
		majsoulDoras = msg.Doras
	case *lq.RecordDealTile:
		majsoulDoras = msg.Doras
	}

	return a.isNewDora(majsoulDoras)
}

func (a *Action) ParseNewDora() (kanDoraIndicator int) {
	var majsoulDoras []string
	switch msg := a.parsedAction.(type) {
	case *lq.ActionDealTile:
		majsoulDoras = msg.Doras
	case *lq.RecordDealTile:
		majsoulDoras = msg.Doras
	}

	kanDoraIndicator, _ = a.mustParseMajsoulTile(majsoulDoras[len(majsoulDoras)-1])
	return
}
