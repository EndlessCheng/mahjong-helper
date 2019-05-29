package main

import (
	"github.com/EndlessCheng/mahjong-helper/util"
	"fmt"
	"os"
	"github.com/EndlessCheng/mahjong-helper/util/model"
)

func interact(humanTilesInfo *model.HumanTilesInfo) {
	tiles34, err := analysisHumanTiles(humanTilesInfo)
	if err != nil {
		errorExit(err)
	}
	var tile string
	for {
		count := util.CountOfTiles34(tiles34)
		switch count % 3 {
		case 0:
			errorExit("参数错误", count, "张牌")
		case 1:
			fmt.Print("> 摸 ")
			fmt.Scanf("%s\n", &tile)

			// TODO: 0p

			tile34, err := util.StrToTile34(tile)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				break
			}
			if tiles34[tile34] == 4 {
				fmt.Fprintln(os.Stderr, "不可能摸更多的牌了")
				break
			}
			tiles34[tile34]++
			humanTilesInfo.HumanTiles = util.Tiles34ToStr(tiles34)
			if _, err := analysisHumanTiles(humanTilesInfo); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		case 2:
			fmt.Print("> 切 ")
			fmt.Scanf("%s\n", &tile)
			tile34, err := util.StrToTile34(tile)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				break
			}
			if tiles34[tile34] == 0 {
				fmt.Fprintln(os.Stderr, "切掉的牌不存在")
				break
			}
			tiles34[tile34]--
			humanTilesInfo.HumanTiles = util.Tiles34ToStr(tiles34)
			if _, err := analysisHumanTiles(humanTilesInfo); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}
	}
}
