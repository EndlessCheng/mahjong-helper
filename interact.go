package main

import (
	"github.com/EndlessCheng/mahjong-helper/util"
	"fmt"
	"os"
)

func interact(raw string) {
	tiles34, err := analysisHumanTiles(raw)
	if err != nil {
		errorExit(err)
	}
	printed := true
	countOfTiles := util.CountOfTiles34(tiles34)

	var tile string
	for {
		for {
			if countOfTiles < 14 {
				countOfTiles = 999
				break
			}
			printed = false
			fmt.Print("> 切 ")
			fmt.Scanf("%s\n", &tile)
			tile34, err := util.StrToTile34(tile)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			} else {
				if tiles34[tile34] == 0 {
					fmt.Fprintln(os.Stderr, "切掉的牌不存在")
				} else {
					tiles34[tile34]--
					break
				}
			}
		}

		if !printed {
			raw = util.Tiles34ToStr(tiles34)
			if _, err := analysisHumanTiles(raw); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}

			printed = true
		}

		for {
			printed = false

			fmt.Print("> 摸 ")
			fmt.Scanf("%s\n", &tile)
			tile34, err := util.StrToTile34(tile)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			} else {
				if tiles34[tile34] == 4 {
					fmt.Fprintln(os.Stderr, "不可能摸更多的牌了")
				} else {
					tiles34[tile34]++
					break
				}
			}
		}

		if !printed {
			raw = util.Tiles34ToStr(tiles34)
			if _, err := analysisHumanTiles(raw); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}

			printed = true
		}
	}
}
