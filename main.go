package main

import (
	"strings"
	"fmt"
	"os"
	"github.com/fatih/color"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"github.com/EndlessCheng/mahjong-helper/util"
	"math/rand"
	"time"
)

// go build -ldflags "-X main.version=$(git describe --abbrev=0 --tags)" -o mahjong-helper
var version = "dev"

var (
	showImproveDetail      bool
	showAgariAboveShanten1 bool
	showScore              bool
	showAllYakuTypes       bool
)

const (
	platformTenhou         = 0
	platformMajsoul        = 1
	platformMajsoulOldYaku = 9

	defaultPlatform = platformMajsoul
)

var platforms map[int]string

func init() {
	rand.Seed(time.Now().UnixNano())

	platforms = map[int]string{
		platformTenhou:  "天凤",
		platformMajsoul: "雀魂",
	}

	now := time.Now()
	const tf = "2006-01-02 15:04:05"
	start, err := time.Parse(tf, "2019-06-07 05:00:00")
	if err != nil {
		panic(err)
	}
	end, err := time.Parse(tf, "2019-06-10 05:00:00")
	if now.After(start) && now.Before(end) {
		platforms[platformMajsoulOldYaku] = "雀魂-乱斗之间"
	}
}

func welcome() int {
	fmt.Println("使用说明：https://github.com/EndlessCheng/mahjong-helper")
	fmt.Println("问题反馈：https://github.com/EndlessCheng/mahjong-helper/issues")
	fmt.Println("吐槽群：375865038")
	fmt.Println()

	fmt.Println("请输入数字，以选择对应的平台：")
	for i, cnt := 0, 0; cnt < len(platforms); i++ {
		if platformName, ok := platforms[i]; ok {
			fmt.Printf("%d - %s\n", i, platformName)
			cnt++
		}
	}

	choose := defaultPlatform
	fmt.Scanf("%d", &choose)
	platformName, ok := platforms[choose]
	if !ok {
		choose = defaultPlatform
		platformName = platforms[choose]
	}

	clearConsole()
	if choose == platformMajsoul {
		platformName += "（水晶杠杠版）"
	}
	color.HiGreen("已选择 - %s", platformName)
	if choose == platformMajsoul || choose == platformMajsoulOldYaku {
		color.HiYellow("提醒：若您已登录游戏，请刷新网页，或者开启一局人机对战\n" +
			"该步骤用于获取您的账号 ID，便于在游戏开始时分析自风，否则程序将无法解析后续数据")
	}

	return choose
}

func main() {
	color.HiGreen("日本麻将助手 %s (by EndlessCheng)", version)
	if version != "dev" {
		go alertNewVersion(version)
	}

	flags, restArgs := parseArgs(os.Args[1:])

	isMajsoul := flags.Bool("majsoul")
	isTenhou := flags.Bool("tenhou")
	isAnalysis := flags.Bool("analysis")
	isInteractive := flags.Bool("i", "interactive")
	showImproveDetail = flags.Bool("detail")
	showAgariAboveShanten1 = flags.Bool("a", "agari")
	showScore = flags.Bool("s", "score")
	considerOldYaku := flags.Bool("old")
	showAllYakuTypes = flags.Bool("y", "yaku")

	util.SetConsiderOldYaku(considerOldYaku)

	humanDoraTiles := flags.String("d", "dora")
	humanTiles := strings.Join(restArgs, " ")
	humanTilesInfo := &model.HumanTilesInfo{
		HumanTiles:     humanTiles,
		HumanDoraTiles: humanDoraTiles,
	}

	switch {
	case isMajsoul:
		runServer(true)
	case isTenhou || isAnalysis:
		runServer(false)
	case isInteractive:
		// 交互模式
		if err := interact(humanTilesInfo); err != nil {
			errorExit(err)
		}
	case len(restArgs) > 0:
		// 静态分析
		if _, err := analysisHumanTiles(humanTilesInfo); err != nil {
			errorExit(err)
		}
	default:
		// 服务器模式
		choose := welcome()
		considerOldYaku = choose == platformMajsoulOldYaku
		util.SetConsiderOldYaku(considerOldYaku)
		isHTTPS := choose == platformMajsoul || choose == platformMajsoulOldYaku
		runServer(isHTTPS)
	}
}
