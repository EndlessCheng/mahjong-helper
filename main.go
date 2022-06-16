package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/EndlessCheng/mahjong-helper/Console"

	"github.com/EndlessCheng/mahjong-helper/util"
	"github.com/EndlessCheng/mahjong-helper/util/model"
	"github.com/fatih/color"
)

// Enum
const (
	TenHou int = iota
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
			Code: TenHou,
		},
		{
			Name: "雀魂",
			Type: []string{
				"國際中文服",
				"日服",
				"國際服"},
			Code: MahJongSoul,
		},
	}
)

func init() {
	rand.Seed(time.Now().UnixNano())

	// this program can run different mode by command-line with argument "old, majsoul, tenhou, analysis, interactive,
	// i detail agari, a, score, s, yaku, y, dora, d, port and p" to use different mode"
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

// Print description of README, question issues, and community and 
// let player can choose witch platform and reminder
func welcome() int {
	fmt.Println("使用说明：" + "https://github.com/EndlessCheng/mahjong-helper/blob/master/README.md")
	fmt.Println("问题反馈：" + "https://github.com/EndlessCheng/mahjong-helper/issues")
	fmt.Println("吐槽群：" + "375865038")
	fmt.Println()

RenterPlatform: // wrong enter goto label
	// print platforms
	for _, element := range Platforms {
		
		fmt.Printf("%d - %s %v\n\n", element.Code, element.Name, "[" + strings.Join(element.Type,`,` ) +"]")
	}
	fmt.Print("請選擇對應的網站(0或1)，如未選擇則預設雀魂(1): ")

	// set default value to int MahJongSoul(1) can exclude not int type
	choose := TenHou
	// choose := MahJongSoul
	fmt.Scanln(&choose)

	Console.ClearScreen()
	if choose == TenHou { // choose TenHou
		color.HiGreen("已選擇 - %s", Platforms[0].Name)
	} else if choose == MahJongSoul { // choose MahJongSoul
		color.HiGreen("已選擇 - %s", Platforms[1].Name)
		if len(gameConf.MajsoulAccountIDs) == 0 {
			color.HiYellow(`
提醒：首次启用时，请开启一局人机对战，或者重登游戏。
该步骤用于获取您的账号 ID，便于在游戏开始时获取自风，否则程序将无法解析后续数据。

若助手无响应，请确认您已按步骤安装完成。
相关链接 https://github.com/EndlessCheng/mahjong-helper/issues/104`)
		}
	} else { // the choice not in selection
		fmt.Printf("輸入錯誤，請重新輸入選擇\n\n")
		// goto RenterPlatform label
		goto RenterPlatform
	}
	return choose
}

func main() {
	Console.ClearScreen()
	flag.Parse()

	// print text with green color
	color.HiGreen("日本麻将助手 ver.%s (by EndlessCheng)", Version)
	if Version != VersionDev {
		go CheckNewVersion(Version)
	}

	// set consider old yaku to false
	util.SetConsiderOldYaku(ConsiderOldYaku)

	HumanTiles := strings.Join(flag.Args(), " ")
	HumanTilesInfo := &model.HumanTilesInfo{
		HumanTiles:     HumanTiles,
		HumanDoraTiles: HumanDoraTiles,
	}

	var err error
	// switch to different mode by command-line argument
	switch {
	case IsMajsoul:
		err = RunServer(true, Port)
	case IsTenhou || IsAnalysis:
		err = RunServer(true, Port)
	case IsInteractive: // 交互模式
		err = Interact(HumanTilesInfo)
	case len(flag.Args()) > 0: // 静态分析
		_, err = AnalysisHumanTiles(HumanTilesInfo)
	default: // 服务器模式
		choose := welcome()
		IsHTTPS := choose == MahJongSoul
		err = RunServer(IsHTTPS, Port)
	}
	if err != nil {
		ErrorExit(err)
	}
}
