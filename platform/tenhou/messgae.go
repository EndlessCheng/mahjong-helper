package tenhou

type message struct {
	Tag string `json:"tag" xml:"-"`

	//Name string `json:"name"` // id
	//Sex  string `json:"sx"`

	UserName string `json:"uname" xml:"-"`
	//RatingScale string `json:"ratingscale"`

	//N string `json:"n"`
	//J string `json:"j"`
	//G string `json:"g"`

	// round 开始 tag=INIT
	Seed   string `json:"seed" xml:"seed,attr"` // 本局信息：场数，场棒数，立直棒数，骰子A减一，骰子B减一，宝牌指示牌 1,0,0,3,2,92
	Ten    string `json:"ten" xml:"ten,attr"`   // 各家点数 280,230,240,250
	Dealer string `json:"oya" xml:"oya,attr"`   // 庄家 0=自家, 1=下家, 2=对家, 3=上家
	Hai    string `json:"hai" xml:"hai,attr"`   // 初始手牌 30,114,108,31,78,107,25,23,2,14,122,44,49
	Hai0   string `json:"-" xml:"hai0,attr"`
	Hai1   string `json:"-" xml:"hai1,attr"`
	Hai2   string `json:"-" xml:"hai2,attr"`
	Hai3   string `json:"-" xml:"hai3,attr"`

	// 摸牌 tag=T编号，如 T68

	// 副露 tag=N
	Who  string `json:"who" xml:"who,attr"` // 副露者 0=自家, 1=下家, 2=对家, 3=上家
	Meld string `json:"m" xml:"m,attr"`     // 副露编号 35914

	// 杠宝牌指示牌 tag=DORA
	// `json:"hai"` // 杠宝牌指示牌 39

	// 立直声明 tag=REACH, step=1
	// `json:"who"` // 立直者
	Step string `json:"step" xml:"step,attr"` // 1

	// 立直成功，扣1000点 tag=REACH, step=2
	// `json:"who"` // 立直者
	// `json:"ten"` // 立直成功后的各家点数 250,250,240,250
	// `json:"step"` // 2

	// 自摸/有人放铳 tag=牌, t>=8
	T string `json:"t"` // 选项

	// 和牌 tag=AGARI
	// ba, hai, m, machi, ten, yaku, doraHai, who, fromWho, sc
	//Ba string `json:"ba"` // 0,0
	// `json:"hai"` // 和牌型 8,9,11,14,19,125,126,127
	// `json:"m"` // 副露编号 13527,50794
	//Machi string `json:"machi"` // (待ち) 自摸/荣和的牌 126
	// `json:"ten"` // 符数,点数,这张牌的来源 30,7700,0
	//Yaku        string `json:"yaku"`       // 役（编号，翻数） 18,1,20,1,34,2
	//DoraTile    string `json:"doraHai"`    // 宝牌 123
	//UraDoraTile string `json:"doraHaiUra"` // 里宝牌 77
	// `json:"who"` // 和牌者
	//FromWho string `json:"fromWho"` // 自摸/荣和牌的来源
	//Score   string `json:"sc"`      // 各家增减分 260,-77,310,77,220,0,210,0

	// 游戏结束 tag=PROF

	// 重连 tag=GO
	// type, lobby, gpid
	//Type  string `json:"type"`
	//Lobby string `json:"lobby"`
	//GPID  string `json:"gpid"`

	// 重连 tag=REINIT
	// `json:"seed"`
	// `json:"ten"`
	// `json:"oya"`
	// `json:"hai"`
	//Meld1    string `json:"m1"` // 各家副露编号 17450
	//Meld2    string `json:"m2"`
	//Meld3    string `json:"m3"`
	//Kawa0 string `json:"kawa0"` // 各家牌河 112,73,3,131,43,98,78,116
	//Kawa1 string `json:"kawa1"`
	//Kawa2 string `json:"kawa2"`
	//Kawa3 string `json:"kawa3"`
}
