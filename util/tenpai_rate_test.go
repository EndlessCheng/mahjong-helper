package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/EndlessCheng/mahjong-helper/util/model"
)

func TestCalcTenpaiRate(t *testing.T) {
	assert := assert.New(t)
	const eps = 1e-3
	assert.Equal(0.0, CalcTenpaiRate(nil, nil, nil))
	assert.Equal(0.0, CalcTenpaiRate([]*model.Meld{{MeldType: model.MeldTypeAnkan}}, nil, nil))
	assert.Equal(100.0, CalcTenpaiRate([]*model.Meld{{}, {}, {}, {}}, nil, nil))
	assert.InDelta(19.88, CalcTenpaiRate([]*model.Meld{{}}, []int{1, 2, 3, 4, 5}, []int{2}), eps)
	assert.InDelta(23.24, CalcTenpaiRate([]*model.Meld{{}, {}}, []int{1, 2, 3, 4, 5}, []int{2, 4}), eps)
	assert.InDelta(98.26, CalcTenpaiRate([]*model.Meld{{}, {}, {}}, []int{1, 2, 3, 4, 5, 23, 16, 12, -4, -6, 7, 2}, []int{2, 4, 6}), eps)
}
