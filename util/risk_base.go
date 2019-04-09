package util

type RiskTiles34 []float64

// 根据巡目（对于对手而言）、现物、立直后通过的牌的、NC，来计算基础铳率
// 至于早外、Dora、OC 和读牌交给后续的计算
// turns: 巡目，这里是对于对手而言的，也就是该玩家舍牌的次数
// safeTiles34: 现物及立直后通过的牌
// leftTiles34: 各个牌在山中剩余的枚数
// roundWindTile: 场风
// playerWindTile: 自风
func CalculateRiskTiles34(turns int, safeTiles34 []bool, leftTiles34 []int, roundWindTile int, playerWindTile int) (risk34 RiskTiles34) {
	// 生成用来计算筋牌的「安牌」
	sujiSafeTiles34 := make([]int, 34)
	for i, safe := range safeTiles34 {
		if safe {
			sujiSafeTiles34[i] = 1
		}
	}
	for i := 0; i < 3; i++ {
		// 2断，当做打过1
		if leftTiles34[9*i+1] == 0 {
			sujiSafeTiles34[9*i] = 1
		}
		// 3断，当做打过12
		if leftTiles34[9*i+2] == 0 {
			sujiSafeTiles34[9*i] = 1
			sujiSafeTiles34[9*i+1] = 1
		}
		// 4断，当做打过23
		if leftTiles34[9*i+3] == 0 {
			sujiSafeTiles34[9*i+1] = 1
			sujiSafeTiles34[9*i+2] = 1
		}
		// 6断，当做打过78
		if leftTiles34[9*i+5] == 0 {
			sujiSafeTiles34[9*i+6] = 1
			sujiSafeTiles34[9*i+7] = 1
		}
		// 7断，当做打过89
		if leftTiles34[9*i+6] == 0 {
			sujiSafeTiles34[9*i+7] = 1
			sujiSafeTiles34[9*i+8] = 1
		}
		// 8断，当做打过9
		if leftTiles34[9*i+7] == 0 {
			sujiSafeTiles34[9*i+8] = 1
		}
	}

	risk34 = make(RiskTiles34, 34)

	// 利用「安牌」计算无筋、筋、半筋、双筋的铳率
	// TODO: 单独处理宣言牌的筋牌、宣言牌的同色牌的铳率
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			t := TileTypeTable[j][sujiSafeTiles34[9*i+j+3]]
			risk34[9*i+j] = RiskData[turns][t]
		}
		for j := 3; j < 6; j++ {
			mixSafeTile := sujiSafeTiles34[9*i+j-3]<<1 | sujiSafeTiles34[9*i+j+3]
			t := TileTypeTable[j][mixSafeTile]
			risk34[9*i+j] = RiskData[turns][t]
		}
		for j := 6; j < 9; j++ {
			t := TileTypeTable[j][sujiSafeTiles34[9*i+j-3]]
			risk34[9*i+j] = RiskData[turns][t]
		}
		// 5断，37视作安牌筋
		if leftTiles34[9*i+4] == 0 {
			t := tileTypeSuji37
			risk34[9*i+2] = RiskData[turns][t]
			risk34[9*i+6] = RiskData[turns][t]
		}
	}
	for i := 27; i < 34; i++ {
		if leftTiles34[i] > 0 {
			// 该玩家的役牌 = 场风/其自风/白/发/中
			isYakuHai := i == roundWindTile || i == playerWindTile || i >= 31
			t := HonorTileType[boolToInt(isYakuHai)][leftTiles34[i]-1]
			risk34[i] = RiskData[turns][t]
		} else {
			// 剩余数为0可以视作安牌（只输国士）
			risk34[i] = 0
		}
	}

	// 更新铳率表：NC牌的安牌
	// 12和筋1差不多（2比1多10%）
	// 3和筋2差不多
	// 456和两筋差不多（存疑？）
	ncSafeTile34 := CalcNCSafeTiles(leftTiles34)
	for _, ncSafeTile := range ncSafeTile34 {
		switch ncSafeTile.Tile34 % 9 {
		case 1, 9:
			risk34[ncSafeTile.Tile34] = RiskData[turns][tileTypeSuji19]
		case 2, 8:
			risk34[ncSafeTile.Tile34] = RiskData[turns][tileTypeSuji19] * 1.1
		case 3, 7:
			risk34[ncSafeTile.Tile34] = RiskData[turns][tileTypeSuji28]
		case 4, 6:
			risk34[ncSafeTile.Tile34] = RiskData[turns][tileTypeDoubleSuji46]
		case 5:
			risk34[ncSafeTile.Tile34] = RiskData[turns][tileTypeDoubleSuji5]
		}
	}

	// 更新铳率表：DNC且剩余枚数为0的也当作安牌（忽略国士）
	dncSafeTiles := CalcDNCSafeTiles(leftTiles34)
	for _, dncSafeTile := range dncSafeTiles {
		if leftTiles34[dncSafeTile.Tile34] == 0 {
			risk34[dncSafeTile.Tile34] = 0
		}
	}

	// 更新铳率表：现物的铳率为0
	for i, isSafe := range safeTiles34 {
		if isSafe {
			risk34[i] = 0
		}
	}

	return
}

// TODO: 利用剩余牌是否为 0 或者 1 计算 No Chance, One Chance, Double One Chance, Double Two Chance(待定) 等
// TODO: 利用舍牌计算无筋早外
// TODO:（待定）有早外的半筋（早巡打过8m时，3m的半筋6m）
// TODO:（待定）利用赤宝牌计算铳率
// TODO: 宝牌周边牌的危险度要增加一点
// TODO:（待定）切过5的情况
