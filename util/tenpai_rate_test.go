package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCalcTenpaiRate(t *testing.T) {
	assert := assert.New(t)
	const eps = 1e-3
	assert.Equal(0.0, CalcTenpaiRate(0, nil, nil))
	assert.Equal(100.0, CalcTenpaiRate(4, nil, nil))
	assert.InDelta(19.88, CalcTenpaiRate(1, []int{1, 2, 3, 4, 5}, []int{2}), eps)
	assert.InDelta(23.24, CalcTenpaiRate(2, []int{1, 2, 3, 4, 5}, []int{2, 4}), eps)
	assert.InDelta(98.15, CalcTenpaiRate(3, []int{1, 2, 3, 4, 5, 23, 16, 12, -4, -6, 7, 2}, []int{2, 4, 6}), eps)
}
