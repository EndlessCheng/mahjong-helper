package main

import (
	"flag"
	"fmt"
	"github.com/EndlessCheng/mahjong-helper/util"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"github.com/fatih/color"
	"math/rand"
	"strings"
	"time"
)

var (
	considerOldYaku bool

	isMajsoul     bool
	isTenhou      bool
	isAnalysis    bool
	isInteractive bool

	showImproveDetail      bool
	showAgariAboveShanten1 bool
	showScore              bool
	showAllYakuTypes       bool

	humanDoraTiles string

	port int
)

func init() {
	rand.Seed(time.Now().UnixNano())

	flag.BoolVar(&considerOldYaku, "old", false, "允许古役")
	flag.BoolVar(&isMajsoul, "majsoul", false, "雀魂助手")
	flag.BoolVar(&isTenhou, "tenhou", false, "天凤助手")
	flag.BoolVar(&isAnalysis, "analysis", false, "分析模式")
	flag.BoolVar(&isInteractive, "interactive", false, "交互模式")
	flag.BoolVar(&isInteractive, "i", false, "同 -interactive")
	flag.BoolVar(&showImproveDetail, "detail", false, "显示改良细节")
	flag.BoolVar(&showAgariAboveShanten1, "agari", false, "显示听牌前的估计和率")
	flag.BoolVar(&showAgariAboveShanten1, "a", false, "同 -agari")
	flag.BoolVar(&showScore, "score", false, "显示局收支")
	flag.BoolVar(&showScore, "s", false, "同 -score")
	flag.BoolVar(&showAllYakuTypes, "yaku", false, "显示所有役种")
	flag.BoolVar(&showAllYakuTypes, "y", false, "同 -yaku")
	flag.StringVar(&humanDoraTiles, "dora", "", "指定哪些牌是宝牌")
	flag.StringVar(&humanDoraTiles, "d", "", "同 -dora")
	flag.IntVar(&port, "port", 12121, "指定服务端口")
	flag.IntVar(&port, "p", 12121, "同 -port")
}

const (
	platformTenhou  = 0
	platformMajsoul = 1

	defaultPlatform = platformMajsoul
)

var platforms = map[int][]string{
	platformTenhou: {
		"天凤",
		"Web",
		"4K",
	},
	platformMajsoul: {
		"雀魂",
		"国际中文服",
		"日服",
		"国际服",
	},
}

const readmeURL = "https://github.com/EndlessCheng/mahjong-helper/blob/master/README.md"
const issueURL = "https://github.com/EndlessCheng/mahjong-helper/issues"
const issueCommonQuestions = "https://github.com/EndlessCheng/mahjong-helper/issues/104"
const qqGroupNum = "375865038"

func welcome() int {
	fmt.Println("使用说明：" + readmeURL)
	fmt.Println("问题反馈：" + issueURL)
	fmt.Println("吐槽群：" + qqGroupNum)
	fmt.Println()

	fmt.Println("请输入数字，选择对应网站：")
	for i, cnt := 0, 0; cnt < len(platforms); i++ {
		if platformInfo, ok := platforms[i]; ok {
			info := platformInfo[0] + " [" + strings.Join(platformInfo[1:], ",") + "]"
			fmt.Printf("%d - %s\n", i, info)
			cnt++
		}
	}

	choose := defaultPlatform
	fmt.Scanln(&choose) // 直接回车也无妨
	platformInfo, ok := platforms[choose]
	var platformName string
	if ok {
		platformName = platformInfo[0]
	}
	if !ok {
		choose = defaultPlatform
		platformName = platforms[choose][0]
	}

	clearConsole()
	color.HiGreen("已选择 - %s", platformName)

	if choose == platformMajsoul {
		if len(gameConf.MajsoulAccountIDs) == 0 {
			color.HiYellow(`
提醒：首次启用时，请开启一局人机对战，或者重登游戏。
该步骤用于获取您的账号 ID，便于在游戏开始时获取自风，否则程序将无法解析后续数据。

若助手无响应，请确认您已按步骤安装完成。
相关链接 ` + issueCommonQuestions)
		}
	}

	return choose
}

func main() {
	flag.Parse()

	color.HiGreen("日本麻将助手 %s (by EndlessCheng)", version)
	if version != versionDev {
		go checkNewVersion(version)
	}

	util.SetConsiderOldYaku(considerOldYaku)

	humanTiles := strings.Join(flag.Args(), " ")
	humanTilesInfo := &model.HumanTilesInfo{
		HumanTiles:     humanTiles,
		HumanDoraTiles: humanDoraTiles,
	}

	var err error
	switch {
	case isMajsoul:
		err = runServer(true, port)
	case isTenhou || isAnalysis:
		err = runServer(true, port)
	case isInteractive: // 交互模式
		err = interact(humanTilesInfo)
	case len(flag.Args()) > 0: // 静态分析
		_, err = analysisHumanTiles(humanTilesInfo)
	default: // 服务器模式
		choose := welcome()
		isHTTPS := choose == platformMajsoul
		err = runServer(isHTTPS, port)
	}
	if err != nil {
		errorExit(err)
	}
}
