package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/EndlessCheng/mahjong-helper/util/model"
)

func TestCalculateAvgAgariRate(t *testing.T) {
	assert := assert.New(t)
	const eps = 1e-3
	assert.InDelta(62.1166, CalculateAvgAgariRate(Waits{0: 4, 3: 4}, nil), eps)
	assert.InDelta(57.715203, CalculateAvgAgariRate(Waits{0: 3, 3: 3}, nil), eps)
	assert.InDelta(34.6678, CalculateAvgAgariRate(Waits{0: 3, 3: 4}, &model.PlayerInfo{DiscardTiles: []int{0}}), eps) // 振听
	assert.InDelta(65.8944, CalculateAvgAgariRate(Waits{0: 2, 9: 2}, nil), eps)
	assert.InDelta(71.058, CalculateAvgAgariRate(Waits{0: 3, 1: 4}, nil), eps)
	assert.InDelta(71.058, CalculateAvgAgariRate(Waits{8: 3, 7: 4}, nil), eps)
	assert.InDelta(96.2222, CalculateAvgAgariRate(Waits{0: 1, 1: 3, 2: 3, 3: 3, 4: 3, 5: 3, 6: 3, 7: 3, 9: 1}, nil), eps)
	assert.InDelta(71.968, CalculateAvgAgariRate(Waits{9: 2, 27: 2}, nil), eps)
	assert.InDelta(49.5, CalculateAvgAgariRate(Waits{27: 3}, nil), eps)
	assert.InDelta(58, CalculateAvgAgariRate(Waits{27: 2}, nil), eps)
	assert.InDelta(47.5, CalculateAvgAgariRate(Waits{27: 1}, nil), eps)
	assert.InDelta(0, CalculateAvgAgariRate(Waits{27: 0}, nil), eps)
	assert.InDelta(49.6629, CalculateAvgAgariRate(Waits{0: 1, 7: 2}, nil), eps)
	assert.InDelta(51.31672, CalculateAvgAgariRate(Waits{2: 4, 5: 4}, nil), eps)
	assert.InDelta(54.5818, CalculateAvgAgariRate(Waits{4: 4, 7: 4}, nil), eps)
	assert.InDelta(61.744, CalculateAvgAgariRate(Waits{5: 2, 31: 2}, nil), eps)
	assert.InDelta(48.93434, CalculateAvgAgariRate(Waits{1: 4, 4: 2}, nil), eps)
	assert.InDelta(54.5818, CalculateAvgAgariRate(Waits{1: 4, 4: 4}, nil), eps)
}
