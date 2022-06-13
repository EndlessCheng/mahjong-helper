package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/EndlessCheng/mahjong-helper/util"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"github.com/fatih/color"
)

// Enum
const (
	Tenhou int = iota
	MahJongSoul
)

// define Platform Class parameter
type Platform struct {
	Name string
	Type []string
	Code int
}

// declare variable
var (
	// bool
	ConsiderOldYaku        bool = false
	IsMajsoul              bool = false
	IsTenhou               bool = false
	IsAnalysis             bool = false
	IsInteractive          bool = false
	ShowImproveDetail      bool = false
	ShowAgariAboveShanten1 bool = false
	ShowScore              bool = false
	ShowAllYakuTypes       bool = false

	//int
	Port int = 0

	// string
	HumanDoraTiles string = ""

	// struct
	Platforms []Platform = []Platform{
		{
			Name: "天鳳",
			Type: []string{
				"Web",
				"4K"},
			Code: Tenhou,
		},
		{
			Name: "雀魂",
			Type: []string{
				"國際中文服",
				"日服",
				"国际服"},
			Code: MahJongSoul,
		},
	}
)

func init() {
	rand.Seed(time.Now().UnixNano())

	flag.BoolVar(&ConsiderOldYaku, "old", false, "允许古役")
	flag.BoolVar(&IsMajsoul, "majsoul", false, "雀魂助手")
	flag.BoolVar(&IsTenhou, "tenhou", false, "天凤助手")
	flag.BoolVar(&IsAnalysis, "analysis", false, "分析模式")
	flag.BoolVar(&IsInteractive, "interactive", false, "交互模式")
	flag.BoolVar(&IsInteractive, "i", false, "同 -interactive")
	flag.BoolVar(&ShowImproveDetail, "detail", false, "显示改良细节")
	flag.BoolVar(&ShowAgariAboveShanten1, "agari", false, "显示听牌前的估计和率")
	flag.BoolVar(&ShowAgariAboveShanten1, "a", false, "同 -agari")
	flag.BoolVar(&ShowScore, "score", false, "显示局收支")
	flag.BoolVar(&ShowScore, "s", false, "同 -score")
	flag.BoolVar(&ShowAllYakuTypes, "yaku", false, "显示所有役种")
	flag.BoolVar(&ShowAllYakuTypes, "y", false, "同 -yaku")
	flag.StringVar(&HumanDoraTiles, "dora", "", "指定哪些牌是宝牌")
	flag.StringVar(&HumanDoraTiles, "d", "", "同 -dora")
	flag.IntVar(&Port, "port", 12121, "指定服务端口")
	flag.IntVar(&Port, "p", 12121, "同 -port")
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

RenterPlatform: // wrong enter goto label
	// print platforms
	for _, element := range Platforms {
		fmt.Printf("%d - %s %v\n", element.Code, element.Name, element.Type)
	}
	fmt.Print("請選擇對應的網站(0或1)，如未選擇則預設雀魂(1): ")

	// set default value to int MahJongSoul(1) can exclude not int type
	choose := MahJongSoul
	fmt.Scanln(&choose)

	clearConsole()
	if choose == Tenhou { // choose TenHou
		color.HiGreen("已選擇 - %s", Platforms[0].Name)
	} else if choose == MahJongSoul { // choose MahJongSoul
		color.HiGreen("已選擇 - %s", Platforms[1].Name)
		if len(gameConf.MajsoulAccountIDs) == 0 {
			color.HiYellow(`
提醒：首次启用时，请开启一局人机对战，或者重登游戏。
该步骤用于获取您的账号 ID，便于在游戏开始时获取自风，否则程序将无法解析后续数据。

若助手无响应，请确认您已按步骤安装完成。
相关链接 ` + issueCommonQuestions)
		}
	} else { // the choice not in selection
		fmt.Printf("輸入錯誤，請重新輸入選擇\n\n")
		// goto RenterPlatform label
		goto RenterPlatform
	}
	return choose
}

func main() {
	flag.Parse()

	color.HiGreen("日本麻将助手 %s (by EndlessCheng)", version)
	if version != versionDev {
		go checkNewVersion(version)
	}

	util.SetConsiderOldYaku(ConsiderOldYaku)

	humanTiles := strings.Join(flag.Args(), " ")
	humanTilesInfo := &model.HumanTilesInfo{
		HumanTiles:     humanTiles,
		HumanDoraTiles: HumanDoraTiles,
	}

	var err error
	switch {
	case IsMajsoul:
		err = runServer(true, Port)
	case IsTenhou || IsAnalysis:
		err = runServer(true, Port)
	case IsInteractive: // 交互模式
		err = interact(humanTilesInfo)
	case len(flag.Args()) > 0: // 静态分析
		_, err = analysisHumanTiles(humanTilesInfo)
	default: // 服务器模式
		choose := welcome()
		isHTTPS := choose == MahJongSoul
		err = runServer(isHTTPS, Port)
	}
	if err != nil {
		errorExit(err)
	}
}
