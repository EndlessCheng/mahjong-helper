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

const versionDev = "dev"

// go build -ldflags "-X main.version=$(git describe --abbrev=0 --tags)" -o mahjong-helper
var version = versionDev

var (
	showImproveDetail      bool
	showAgariAboveShanten1 bool
	showScore              bool
	showAllYakuTypes       bool
)

const (
	platformTenhou  = 0
	platformMajsoul = 1

	defaultPlatform = platformMajsoul
)

var platforms map[int]string

func init() {
	rand.Seed(time.Now().UnixNano())

	platforms = map[int]string{
		platformTenhou:  "天凤",
		platformMajsoul: "雀魂",
	}
}

func welcome() int {
	fmt.Println("使用说明：https://github.com/EndlessCheng/mahjong-helper/blob/master/README.md")
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

	if choose == platformMajsoul {
		if len(gameConf.MajsoulAccountIDs) == 0 {
			color.HiYellow("提醒：若您是第一次使用助手，请重新登录游戏，或者开启一局人机对战\n" +
				"该步骤用于获取您的账号 ID，便于在游戏开始时获取自风，否则程序将无法解析后续数据")
		}
	}

	return choose
}

func main() {
	color.HiGreen("日本麻将助手 %s (by EndlessCheng)", version)
	if version != versionDev {
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

	var err error
	switch {
	case isMajsoul:
		err = runServer(true)
	case isTenhou || isAnalysis:
		err = runServer(false)
	case isInteractive: // 交互模式
		err = interact(humanTilesInfo)
	case len(restArgs) > 0: // 静态分析
		_, err = analysisHumanTiles(humanTilesInfo)
	default: // 服务器模式
		choose := welcome()
		isHTTPS := choose == platformMajsoul
		err = runServer(isHTTPS)
	}
	if err != nil {
		errorExit(err)
	}
}
