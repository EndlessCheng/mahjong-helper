package main

import (
	"strings"
	"fmt"
	"os"
	"time"
	"github.com/fatih/color"
	"math/rand"
	"github.com/EndlessCheng/mahjong-helper/util"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var debug = false

var detailFlag = false
var interactFlag = false // 交互模式

func _printIncShantenResults14(shanten int, results14, incShantenResults14 util.WaitsWithImproves14List) {
	if len(incShantenResults14) == 0 {
		return
	}

	bestWaitsCount := results14[0].Result13.Waits.AllCount()
	bestIncShantenWaitsCount := incShantenResults14[0].Result13.Waits.AllCount()

	// TODO: 待调整
	// 1 - 12
	// 2 - 24
	// 3 - 36
	// ...
	incShantenWaitsCountLimit := 12
	for i := 1; i < shanten; i++ {
		incShantenWaitsCountLimit *= 2
	}

	needPrintIncShanten := bestWaitsCount <= incShantenWaitsCountLimit && bestIncShantenWaitsCount >= 2*bestWaitsCount
	if shanten == 0 {
		needPrintIncShanten = bestIncShantenWaitsCount >= 18
	}
	if !needPrintIncShanten {
		return
	}

	if len(incShantenResults14[0].OpenTiles) > 0 {
		fmt.Print("鸣牌后")
	}
	fmt.Println(util.NumberToChineseShanten(shanten+1) + "：")
	for _, result := range incShantenResults14 {
		printWaitsWithImproves13(result.Result13, result.DiscardTile, result.OpenTiles)
	}
}

func _analysis(num int, tiles34 []int, leftTiles34 []int, isOpen bool) error {
	raw := util.Tiles34ToMergedStr(tiles34)
	fmt.Println(raw)
	fmt.Println(strings.Repeat("=", len(raw)))

	switch {
	case num%3 == 1:
		result := util.CalculateShantenWithImproves13(tiles34, isOpen)
		result.Waits.FixCountsWithLeftCounts(leftTiles34)
		fmt.Println(util.NumberToChineseShanten(result.Shanten) + "：")
		printWaitsWithImproves13(result, -1, nil)
	case num%3 == 2:
		if util.CheckWin(tiles34) {
			color.Red("【已胡牌】")
			break
		}

		shanten, results14, incShantenResults14 := util.CalculateShantenWithImproves14(tiles34, isOpen)

		if shanten == 0 {
			color.Red("【已听牌】")
		}

		// TODO: 若两向听的进张<=15，则添加向听倒退的提示（拒绝做七对子）

		// TODO: 改良没算，还是放到解析中去处理比较好
		results14.FilterWithLeftTiles34(leftTiles34)
		incShantenResults14.FilterWithLeftTiles34(leftTiles34)

		fmt.Println(util.NumberToChineseShanten(shanten) + "：")
		for _, result := range results14 {
			printWaitsWithImproves13(result.Result13, result.DiscardTile, result.OpenTiles)
		}

		// 不好的牌会打印出向听倒退的分析
		_printIncShantenResults14(shanten, results14, incShantenResults14)
	default:
		err := fmt.Errorf("参数错误: %d 张牌", num)
		return err
	}

	//fmt.Println("checkWin", checkWinCount)
	fmt.Println()

	return nil
}

func analysisMeld(tiles34 []int, leftTiles34 []int, targetTile34 int, allowChi bool) {
	raw := util.Tiles34ToMergedStr(tiles34) + " + " + util.Tile34ToMergedStr(targetTile34) + "?"
	fmt.Println(raw)
	fmt.Println(strings.Repeat("=", len(raw)))

	result := util.CalculateShantenWithImproves13(tiles34, true)
	result.Waits.FixCountsWithLeftCounts(leftTiles34)
	fmt.Println("原手牌" + util.NumberToChineseShanten(result.Shanten) + "：")
	printWaitsWithImproves13(result, -1, nil)

	// 副露分析
	shanten, results14, incShantenResults14 := util.CalculateMeld(tiles34, targetTile34, allowChi)
	results14.FilterWithLeftTiles34(leftTiles34)
	incShantenResults14.FilterWithLeftTiles34(leftTiles34)

	if shanten == -1 {
		color.Red("【已胡牌】")
		return
	}

	// 打印结果
	fmt.Println("鸣牌后" + util.NumberToChineseShanten(shanten) + "：")
	for _, result := range results14 {
		printWaitsWithImproves13(result.Result13, result.DiscardTile, result.OpenTiles)
	}
	_printIncShantenResults14(shanten, results14, incShantenResults14)
}

func analysis(raw string) (num int, counts []int, err error) {
	num, counts, err = convert(raw)
	if err != nil {
		return
	}

	err = _analysis(num, counts, nil, false)
	return
}

func interact(raw string) {
	num, counts, err := analysis(raw)
	if err != nil {
		_errorExit(err)
	}
	printed := true

	var tile string
	for {
		for {
			if num < 14 {
				num = 999
				break
			}
			printed = false
			fmt.Print("> 切 ")
			fmt.Scanf("%s\n", &tile)
			idx, err := _convert(tile)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			} else {
				if counts[idx] == 0 {
					fmt.Fprintln(os.Stderr, "切掉的牌不存在")
				} else {
					counts[idx]--
					break
				}
			}
		}

		if !printed {
			// 交互模式时，13张牌的一向听分析显示改良具体情况
			detailFlag = true
			raw, _ = countsToString(counts)
			if _, _, err := analysis(raw); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
			detailFlag = false

			printed = true
		}

		for {
			printed = false

			fmt.Print("> 摸 ")
			fmt.Scanf("%s\n", &tile)
			idx, err := _convert(tile)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			} else {
				if counts[idx] == 4 {
					fmt.Fprintln(os.Stderr, "不可能摸更多的牌了")
				} else {
					counts[idx]++
					break
				}
			}
		}

		if !printed {
			raw, _ = countsToString(counts)
			if _, _, err := analysis(raw); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}

			printed = true
		}
	}
}

func welcome() int {
	platforms := map[int]string{
		0: "天凤",
		1: "雀魂",
	}

	fmt.Println("使用前，请确认相关配置已完成，详见 https://github.com/EndlessCheng/mahjong-helper")
	fmt.Println("请输入数字，以选择对应的平台：")
	for k, v := range platforms {
		fmt.Printf("%d - %s\n", k, v)
	}

	choose := 1
	fmt.Scanf("%d", &choose)
	if choose < 0 || choose > 1 {
		choose = 1
	}

	clearConsole()
	platformName := platforms[choose]
	if choose == 1 {
		platformName += "（水晶杠杠版）"
	}
	color.Magenta("已选择 - %s", platformName)
	if choose == 1 {
		color.Yellow("提醒：若您已登录游戏，请刷新网页，或者开启一局人机对战\n" +
			"该步骤用于获取您的账号 ID，便于在游戏开始时分析自风，否则程序将无法解析后续数据")
	}

	return choose
}

func main() {
	// TODO: github.com/urfave/cli
	// TODO: flag 比较下二者
	// TODO: -majsoul -tenhou 快捷方式
	if len(os.Args) <= 1 {
		// 服务器模式
		isHTTPS := welcome() == 1
		runServer(isHTTPS)
		return
	}

	if os.Args[len(os.Args)-1] == "-analysis" {
		runServer(false)
		return
	}

	if os.Args[len(os.Args)-1] == "-i" {
		// （一向听）交互模式
		interactFlag = true

		raw := strings.Join(os.Args[1:len(os.Args)-1], " ")
		interact(raw)
	}

	raw := strings.Join(os.Args[1:], " ")
	if os.Args[len(os.Args)-1] == "-d" {
		// 显示改良细节
		detailFlag = true
		raw = strings.Join(os.Args[1:len(os.Args)-1], " ")
	}

	t0 := time.Now()
	analysis(raw)
	fmt.Printf("耗时 %.2f 秒\n", float64(time.Now().UnixNano()-t0.UnixNano())/float64(time.Second))
}
