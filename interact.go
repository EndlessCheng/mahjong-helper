package main

import (
	"github.com/EndlessCheng/mahjong-helper/util"
	"fmt"
	"os"
	"github.com/EndlessCheng/mahjong-helper/util/model"
)

func interact(humanTilesInfo *model.HumanTilesInfo) (err error) {
	if !debugMode {
		defer func() {
			if er := recover(); er != nil {
				err = fmt.Errorf("参数错误: %v (%v)", er, *humanTilesInfo)
			}
		}()
	}

	playerInfo, err := analysisHumanTiles(humanTilesInfo)
	if err != nil {
		return err
	}
	tiles34 := playerInfo.HandTiles34
	leftTiles34 := playerInfo.LeftTiles34
	var tile string
	for {
		count := util.CountOfTiles34(tiles34)
		switch count % 3 {
		case 0:
			return fmt.Errorf("参数错误: %d 张牌", count)
		case 1:
			fmt.Print("> 摸 ")
			fmt.Scanf("%s\n", &tile)
			tile, isRedFive, er := util.StrToTile34(tile)
			if er != nil {
				// 让用户重新输入
				fmt.Fprintln(os.Stderr, er)
				continue
			}
			if tiles34[tile] == 4 {
				// 让用户重新输入
				fmt.Fprintln(os.Stderr, "不可能摸更多的牌了")
				continue
			}
			if isRedFive {
				playerInfo.NumRedFives[tile/9]++
			}
			leftTiles34[tile]--
			tiles34[tile]++
		case 2:
			fmt.Print("> 切 ")
			fmt.Scanf("%s\n", &tile)
			tile, isRedFive, er := util.StrToTile34(tile)
			if er != nil {
				// 让用户重新输入
				fmt.Fprintln(os.Stderr, er)
				continue
			}
			if tiles34[tile] == 0 {
				// 让用户重新输入
				fmt.Fprintln(os.Stderr, "切掉的牌不存在")
				continue
			}
			if isRedFive {
				playerInfo.NumRedFives[tile/9]--
			}
			tiles34[tile]--
			playerInfo.DiscardTiles = append(playerInfo.DiscardTiles, tile) // 仅判断振听用
		}
		if er := analysisPlayerWithRisk(playerInfo, nil); er != nil {
			fmt.Fprintln(os.Stderr, er)
			continue
		}
	}
}
