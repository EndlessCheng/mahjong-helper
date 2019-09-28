package util

import (
	"fmt"
	"sort"
	"math/rand"
)

var Mahjong = [...]string{
	"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m",
	"1p", "2p", "3p", "4p", "5p", "6p", "7p", "8p", "9p",
	"1s", "2s", "3s", "4s", "5s", "6s", "7s", "8s", "9s",
	"1z", "2z", "3z", "4z", "5z", "6z", "7z",
}

var MahjongU = [...]string{
	"1M", "2M", "3M", "4M", "5M", "6M", "7M", "8M", "9M",
	"1p", "2p", "3p", "4p", "5p", "6p", "7p", "8p", "9p",
	"1S", "2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S",
	"1Z", "2Z", "3Z", "4Z", "5Z", "6Z", "7Z",
}

var MahjongZH = [...]string{
	"1万", "2万", "3万", "4万", "5万", "6万", "7万", "8万", "9万",
	"1饼", "2饼", "3饼", "4饼", "5饼", "6饼", "7饼", "8饼", "9饼",
	"1索", "2索", "3索", "4索", "5索", "6索", "7索", "8索", "9索",
	"东", "南", "西", "北", "白", "发", "中",
}

var YaochuTiles = [...]int{0, 8, 9, 17, 18, 26, 27, 28, 29, 30, 31, 32, 33}

func TilesToMahjongZH(tiles []int) (words []string) {
	for _, tile := range tiles {
		words = append(words, MahjongZH[tile])
	}
	return
}

func TilesToMahjongZHInterface(tiles []int) (words []interface{}) {
	for _, tile := range tiles {
		words = append(words, MahjongZH[tile])
	}
	return
}

// 进张
// map[进张牌]剩余数
type Waits map[int]int

func (w Waits) AllCount() (count int) {
	for _, cnt := range w {
		count += cnt
	}
	return count
}

// 剩余数不为零的进张
func (w Waits) AvailableTiles() []int {
	if len(w) == 0 {
		return nil
	}

	tileIndexes := []int{}
	for idx, left := range w {
		if left > 0 {
			tileIndexes = append(tileIndexes, idx)
		}
	}
	sort.Ints(tileIndexes)

	return tileIndexes
}

func (w Waits) indexes() []int {
	if len(w) == 0 {
		return nil
	}

	tileIndexes := make([]int, 0, len(w))
	for idx := range w {
		tileIndexes = append(tileIndexes, idx)
	}
	sort.Ints(tileIndexes)

	return tileIndexes
}

func (w Waits) ParseIndex() (allCount int, indexes []int) {
	return w.AllCount(), w.indexes()
}

func (w Waits) _parse(template [34]string) (allCount int, tiles []string) {
	if len(w) == 0 {
		return 0, nil
	}

	tileIndexes := make([]int, 0, len(w))
	for idx, cnt := range w {
		tileIndexes = append(tileIndexes, idx)
		allCount += cnt
	}
	sort.Ints(tileIndexes)

	tiles = make([]string, len(tileIndexes))
	for i, idx := range tileIndexes {
		tiles[i] = template[idx]
	}

	return allCount, tiles
}

func (w Waits) parse() (allCount int, tiles []string) {
	return w._parse(Mahjong)
}

func (w Waits) parseZH() (allCount int, tilesZH []string) {
	return w._parse(MahjongZH)
}

func (w Waits) tilesZH() []string {
	_, tiles := w.parseZH()
	return tiles
}

func (w Waits) String() string {
	return fmt.Sprintf("%d 进张 %s", w.AllCount(), TilesToStrWithBracket(w.indexes()))
}

func (w Waits) Equals(w1 Waits) bool {
	tiles0, tiles1 := w.AvailableTiles(), w1.AvailableTiles()
	if len(tiles0) != len(tiles1) {
		return false
	}
	for i := range tiles0 {
		if tiles0[i] != tiles1[i] {
			return false
		}
	}
	return true
}

func isMan(tile int) bool {
	return tile < 9
}

func isPin(tile int) bool {
	return tile >= 9 && tile < 18
}

func isSou(tile int) bool {
	return tile >= 18 && tile < 27
}

func isYaochupai(tile int) bool {
	if tile >= 27 {
		return true
	}
	t := tile % 9
	return t == 0 || t == 8
}

// tiles34 为 13 张牌，判断 tile 若置于 tiles34 中是否是孤张
func isIsolatedTile(tile int, tiles34 []int) bool {
	if tile >= 27 {
		return tiles34[tile] == 0
	}
	t := tile % 9
	l := tile - t + MaxInt(0, t-2)
	r := tile - t + MinInt(8, t+2)
	for i := l; i <= r; i++ {
		if tiles34[i] > 0 {
			return false
		}
	}
	return true
}

// 计算手牌枚数
func CountOfTiles34(tiles34 []int) (count int) {
	for _, c := range tiles34 {
		count += c
	}
	return
}

// 计算手牌对子数
func CountPairsOfTiles34(tiles34 []int) (count int) {
	for _, c := range tiles34 {
		if c >= 2 {
			count++
		}
	}
	return
}

func InitLeftTiles34() []int {
	leftTiles34 := make([]int, 34)
	for i := range leftTiles34 {
		leftTiles34[i] = 4
	}
	return leftTiles34
}

// 根据传入的牌，返回移除这些牌后剩余的牌
func InitLeftTiles34WithTiles34(tiles34 []int) []int {
	leftTiles34 := make([]int, 34)
	for i, count := range tiles34 {
		leftTiles34[i] = 4 - count
	}
	return leftTiles34
}

// 计算外侧牌
func OutsideTiles(tile int) (outsideTiles []int) {
	if tile >= 27 {
		return
	}
	switch tile%9 + 1 {
	case 1, 9:
		return
	case 2, 3, 4:
		for i := tile - tile%9; i < tile; i++ {
			outsideTiles = append(outsideTiles, i)
		}
	case 5:
		// 早巡切5，37 比较安全（TODO 还有片筋A 46）
		outsideTiles = append(outsideTiles, tile-2, tile+2)
	case 6, 7, 8:
		for i := tile - tile%9 + 8; i > tile; i-- {
			outsideTiles = append(outsideTiles, i)
		}
	default:
		panic(fmt.Errorf("[OutsideTiles] 代码有误: tile = %d", tile))
	}
	return
}

// 随机补充一张牌
func RandomAddTile(tiles34 []int) {
	for {
		if tile := rand.Intn(34); tiles34[tile] < 4 {
			tiles34[tile]++
			break
		}
	}
}
