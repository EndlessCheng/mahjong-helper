package model

const (
	TileTypeMan = 0
	TileTypePin = 1
	TileTypeSou = 2
)

// TODO: 其他的也移过来
func InitLeftTiles34WithTiles34(tiles34 []int) []int {
	leftTiles34 := make([]int, 34)
	for i, count := range tiles34 {
		leftTiles34[i] = 4 - count
	}
	return leftTiles34
}
