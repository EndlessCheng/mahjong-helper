package main

import (
	"strings"
	"fmt"
	"os"
	"time"
	"github.com/fatih/color"
	"math/rand"
)

func init() {
	rand.Seed(time.Now().UnixNano())
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
	flags, restArgs := parseArgs(os.Args)
	isMajsoul := flags.Bool("majsoul")
	isTenhou := flags.Bool("tenhou")
	isAnalysis := flags.Bool("analysis")
	isInteractive := flags.Bool("i", "interactive")
	//isDetail := flags.Bool("d", "detail")

	switch {
	case len(flags) == 0 && len(restArgs) == 0:
		// 服务器模式
		isHTTPS := welcome() == 1
		runServer(isHTTPS)
	case isMajsoul:
		runServer(true)
	case isTenhou || isAnalysis:
		runServer(false)
	case isInteractive:
		// 交互模式
		raw := strings.Join(os.Args[1:len(os.Args)-1], " ")
		interact(raw)
	case len(restArgs) > 0:
		humanTiles := strings.Join(restArgs, " ")
		t0 := time.Now()
		if _, err := analysisHumanTiles(humanTiles); err != nil {
			fmt.Println(err)
		}
		fmt.Printf("耗时 %.2f 秒\n", float64(time.Now().UnixNano()-t0.UnixNano())/float64(time.Second))
	}
}
