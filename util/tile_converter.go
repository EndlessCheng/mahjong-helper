package util

import (
	"strings"
	"errors"
	"fmt"
	"github.com/EndlessCheng/mahjong-helper/util/model"
)

func Tiles34ToTiles(tiles34 []int) (tiles []int) {
	for i, c := range tiles34 {
		for j := 0; j < c; j++ {
			tiles = append(tiles, i)
		}
	}
	return
}

func TilesToTiles34(tiles []int) (tiles34 []int) {
	tiles34 = make([]int, 34)
	for _, tile := range tiles {
		tiles34[tile]++
	}
	return
}

// e.g. "3m" => 2
// 可以接收赤5，如 0p
func StrToTile34(humanTile string) (tile34 int, isRedFive bool, err error) {
	defer func() {
		if er := recover(); er != nil {
			err = fmt.Errorf("[StrToTile34] %#v 参数错误: %s", er, humanTile)
		}
	}()

	wrongHumanTileError := errors.New("[StrToTile34] 参数错误: " + humanTile)

	humanTile = strings.TrimSpace(humanTile)
	if len(humanTile) != 2 {
		return -1, false, wrongHumanTileError
	}

	idx := ByteAtStr(Lower(humanTile[1]), "mpsz")
	if idx == -1 {
		return -1, false, wrongHumanTileError
	}

	i := humanTile[0]
	if i == '0' {
		if idx == 3 { // 没有 0z 这种东西
			return -1, false, wrongHumanTileError
		}
		i = '5'
		isRedFive = true
	}

	tile34 = 9*idx + int(i-'1')
	if tile34 >= 34 {
		return -1, false, wrongHumanTileError
	}

	return
}

// 调试用
func MustStrToTile34(humanTile string) int {
	tile34, _, err := StrToTile34(humanTile)
	if err != nil {
		panic(err)
	}
	return tile34
}

// e.g. "224m 24p" => [0, 2, 0, 1, 0, ..., 1, 0, 1, ...]
// 也可以传入不含空格的手牌，如 "224m24p"
// 可以接收赤5，如 0p
func StrToTiles34(humanTiles string) (tiles34 []int, numRedFives []int, err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New("[StrToTiles34] 参数错误: " + humanTiles)
		}
	}()

	// 在 mpsz 后面加上空格方便解析不含空格的 humanTiles
	for _, tileType := range []string{"m", "p", "s", "z"} {
		humanTiles = strings.Replace(humanTiles, tileType, tileType+" ", -1)
	}
	humanTiles = strings.TrimSpace(humanTiles)
	if humanTiles == "" {
		return nil, nil, errors.New("[StrToTiles34] 参数错误: 处理的手牌不能为空")
	}

	tiles34 = make([]int, 34)
	numRedFives = make([]int, 3)
	for _, split := range strings.Split(humanTiles, " ") {
		split = strings.TrimSpace(split)
		if split == "" {
			continue
		}
		if len(split) < 2 {
			return nil, nil, errors.New("[StrToTiles34] 参数错误: " + humanTiles)
		}
		tileType := split[len(split)-1:]
		for _, c := range split[:len(split)-1] {
			tile := string(c) + tileType
			tile34, isRedFive, er := StrToTile34(tile)
			if er != nil {
				return nil, nil, er
			}
			tiles34[tile34]++
			if tiles34[tile34] > 4 {
				return nil, nil, fmt.Errorf("[StrToTiles34] 参数错误: %s 有超过 4 张一样的牌", humanTiles)
			}
			if isRedFive {
				numRedFives[tile34/9]++
			}
		}
	}
	return
}

// 调试用
func MustStrToTiles34(humanTiles string) []int {
	tiles34, _, err := StrToTiles34(humanTiles)
	if err != nil {
		panic(err)
	}
	return tiles34
}

// e.g. "11122z" => [27, 27, 27, 28, 28]
func StrToTiles(humanTiles string) (tiles []int, numRedFives []int, err error) {
	tiles34, numRedFives, err := StrToTiles34(humanTiles)
	if err != nil {
		return
	}
	tiles = Tiles34ToTiles(tiles34)
	return
}

// 调试用
func MustStrToTiles(humanTiles string) []int {
	tiles, _, err := StrToTiles(humanTiles)
	if err != nil {
		panic(err)
	}
	return tiles
}

//

func Tiles34ToStr(tiles34 []int) (humanTiles string) {
	merge := func(lowerIndex, upperIndex int, endsWith string) {
		found := false
		for i, c := range tiles34[lowerIndex:upperIndex] {
			for j := 0; j < c; j++ {
				found = true
				humanTiles += string('1' + byte(i))
			}
		}
		if found {
			humanTiles += endsWith
		}
	}
	merge(0, 9, "m ")
	merge(9, 18, "p ")
	merge(18, 27, "s ")
	merge(27, 34, "z")
	return strings.TrimSpace(humanTiles)
}

// e.g. [9, 11, 27] => "13p 1z"
func TilesToStr(tiles []int) (humanTiles string) {
	return Tiles34ToStr(TilesToTiles34(tiles))
}

func Tile34ToStr(tile34 int) string {
	return TilesToStr([]int{tile34})
}

// e.g. [9, 11, 27] => "[13p 1z]"
func TilesToStrWithBracket(tiles []int) string {
	return "[" + TilesToStr(tiles) + "]"
}

func Tiles34ToStrWithBracket(tiles34 []int) string {
	return "[" + Tiles34ToStr(tiles34) + "]"
}

func ParseHumanTilesWithMelds(humanTilesWithMelds string) (playerInfo *model.PlayerInfo, err error) {
	humanTilesInfo := model.NewSimpleHumanTilesInfo(humanTilesWithMelds)

	if err = humanTilesInfo.SelfParse(); err != nil {
		return
	}

	tiles34, numRedFives, err := StrToTiles34(humanTilesInfo.HumanTiles)
	if err != nil {
		return
	}
	tileCount := CountOfTiles34(tiles34)
	if tileCount%3 == 0 {
		return nil, fmt.Errorf("输入错误: %s 是 %d 张牌", humanTilesInfo.HumanTiles, tileCount)
	}

	melds := []model.Meld{}
	for _, humanMeld := range humanTilesInfo.HumanMelds {
		tiles, _numRedFives, er := StrToTiles(humanMeld)
		if er != nil {
			return nil, er
		}
		isUpper := humanMeld[len(humanMeld)-1] <= 'Z'
		var meldType int
		switch {
		case len(tiles) == 3 && tiles[0] != tiles[1]:
			meldType = model.MeldTypeChi
		case len(tiles) == 3 && tiles[0] == tiles[1]:
			meldType = model.MeldTypePon
		case len(tiles) == 4 && isUpper:
			meldType = model.MeldTypeAnkan
		case len(tiles) == 4 && !isUpper:
			meldType = model.MeldTypeMinkan
		default:
			return nil, fmt.Errorf("输入错误: %s", humanMeld)
		}
		containRedFive := false
		for i, c := range _numRedFives {
			if c > 0 {
				containRedFive = true
				numRedFives[i] += c
			}
		}
		melds = append(melds, model.Meld{
			MeldType:       meldType,
			Tiles:          tiles,
			ContainRedFive: containRedFive,
		})
	}

	playerInfo = model.NewSimplePlayerInfo(tiles34, melds)
	playerInfo.NumRedFives = numRedFives

	return
}

func MustParseHumanTilesWithMelds(humanTilesWithMelds string) *model.PlayerInfo {
	playerInfo, err := ParseHumanTilesWithMelds(humanTilesWithMelds)
	if err != nil {
		panic(err)
	}
	return playerInfo
}
