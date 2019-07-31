package ws

import "reflect"

// TODO
//case "AGARI", "RYUUKYOKU":
//	// 某人和牌或流局，round 结束
//case "PROF":
//	// 游戏结束
//case "BYE":
//	// 某人退出
//case "REJOIN", "GO":
//	// 重连
//case "HELO", "RANKING", "TAIKYOKU", "UN", "LN", "SAIKAI":
//	// 其他

var messageTypes = map[string]reflect.Type{
	"HELO":      reflect.TypeOf((*Helo)(nil)),
	"GO":        reflect.TypeOf((*Go)(nil)),
	"UN":        reflect.TypeOf((*UN)(nil)),
	"TAIKYOKU":  reflect.TypeOf((*Taikyoku)(nil)),
	"INIT":      reflect.TypeOf((*Init)(nil)),
	"REINIT":    reflect.TypeOf((*Init)(nil)),
	"N":         reflect.TypeOf((*Meld)(nil)),
	"REACH":     reflect.TypeOf((*Riichi)(nil)),
	"DORA":      reflect.TypeOf((*Dora)(nil)),
	"AGARI":     reflect.TypeOf((*Agari)(nil)),
	"RYUUKYOKU": reflect.TypeOf((*Ryuukyoku)(nil)),
}

func MessageType(tagName string) reflect.Type {
	return messageTypes[tagName]
}

type Message interface{}

// Login
type Helo struct {
	UserName string `json:"uname"`
}

// Start of game
type Go struct {
	Type  int `json:"type,string"`  // Lobby type
	Lobby int `json:"lobby,string"` // Lobby number
}

// User list or user reconnect
type UN struct {
	N0    string   `json:"n0"` // Player name (URLEncoded UTF-8)
	N1    string   `json:"n1"`
	N2    string   `json:"n2"`
	N3    string   `json:"n3"`
	Ranks Ints     `json:"dan"`  // List of ranks for each player.
	Rates Float64s `json:"rate"` // List of rates for each player
	Sexes Strings  `json:"sx"`   // List of sexes ("M" or "F") for each player.
}

// Start of round
type Taikyoku struct {
	Dealer int `json:"oya,string"`
}

// Start of hand
type Init struct {
	// Six element list:
	//     Round number,
	//     Number of combo sticks,
	//     Number of riichi sticks,
	//     First dice minus one,
	//     Second dice minus one,
	//     Dora indicator.
	Seed Ints `json:"seed"`

	// List of scores for each player
	Scores Ints `json:"ten"`

	Dealer int `json:"oya,string"`

	// Starting hand
	Tiles Ints `json:"hai"`
}

type RecordInit struct {
	// TODO
}

// Player draws a tile
// [T-W][0-9]*
type Draw struct {
	Who  int
	Tile int
	Op   *int `json:"t,string"` // 可以进行的操作，如暗杠、加杠、九种九牌、立直等
}

// Player discards a tile
// 手切 [D-G][0-9]+
// 摸切 [d-g][0-9]+
type Discard struct {
	Who         int
	Tile        int
	IsTsumogiri bool
	Op          *int `json:"t,string"` // 可以进行的操作，如吃、碰、明杠、荣和等
}

// Player calls a tile
type Meld struct {
	Who  int `json:"who,string"`
	Bits int `json:"m,string"`
}

// Player declares riichi
type Riichi struct {
	Who int `json:"who,string"`

	// Where the player is in declaring riichi:
	//     1 -> Called "riichi"
	//     2 -> Placed an 1000 point stick on table after discarding and no one called "ron".
	Step int `json:"step,string"`

	// List of current scores for each player
	Scores Ints `json:"ten"`
}

// New dora indicator
type Dora struct {
	Tile int `json:"hai,string"` // The new dora indicator tile
}

// A player won the hand
type Agari struct {
	// The player who won
	Who int `json:"who,string"`

	// Who the winner won from: themselves for tsumo, someone else for ron
	FromWho int `json:"fromwho,string"`

	// Three element list:
	//     The fu points in the hand,
	//     The point value of the hand,
	//     The limit value of the hand:
	//         0 -> No limit
	//         1 -> Mangan
	//         2 -> Haneman
	//         3 -> Baiman
	//         4 -> Sanbaiman
	//         5 -> Yakuman
	Ten Ints `json:"ten"`

	// TODO owari
}

// The hand ended with a draw
// "{\"tag\":\"RYUUKYOKU\",\"type\":\"ron3\",\"ba\":\"1,1\",\"sc\":\"290,0,228,0,216,0,256,0\",\"hai0\":\"18,19,30,32,33,41,43,94,95,114,115,117,119\",\"hai2\":\"29,31,74,75\",\"hai3\":\"8,13,17,25,35,46,48,53,78,79\"}"
// TODO: add more example
type Ryuukyoku struct {
	// The type of draw:
	//     "yao9"   -> 9 ends
	//     "reach4" -> Four riichi calls
	//     "ron3"   -> Triple ron
	//     "kan4"   -> Four kans
	//     "kaze4"  -> Same wind discard on first round
	//     "nm"     -> Nagashi mangan.
	Type string `json:"type"`

	// TODO owari
}
