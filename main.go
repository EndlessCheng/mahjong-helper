package main

import (
	"flag"
	"github.com/EndlessCheng/mahjong-helper/util"
	"github.com/EndlessCheng/mahjong-helper/util/model"
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

	defaultPlatform = platformTenhou
)

var platforms = map[int]string{
	platformTenhou:  "天凤",
	platformMajsoul: "雀魂",
}


func welcome() int {
	return 0
}

func main() {
	flag.Parse()
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
