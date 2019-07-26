package main

import (
	"strings"
	"fmt"
	"github.com/fatih/color"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"github.com/EndlessCheng/mahjong-helper/util"
	"math/rand"
	"time"
	"flag"
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

var platforms = map[int]string{
	platformTenhou:  "天凤",
	platformMajsoul: "雀魂",
}

const readmeURL = "https://github.com/EndlessCheng/mahjong-helper/blob/master/README.md"

func welcome() int {
	fmt.Println("使用说明：" + readmeURL)
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
			color.HiYellow(`提醒：若您是第一次使用助手，请重新登录游戏，或者开启一局人机对战
该步骤用于获取您的账号 ID，便于在游戏开始时获取自风，否则程序将无法解析后续数据

若助手无响应，请确认您已按步骤安装完成
安装及使用说明：` + readmeURL)
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
		err = runServer(false, port)
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
