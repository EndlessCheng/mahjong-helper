package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCalculateAvgAgariRate(t *testing.T) {
	const eps = 1e-3
	assert.InDelta(t, 68.95, CalculateAvgAgariRate(Waits{0: 4, 3: 4}, nil), eps)
	assert.InDelta(t, 68.95, CalculateAvgAgariRate(Waits{0: 4, 3: 4}, nil), eps)
	assert.InDelta(t, 65.8944, CalculateAvgAgariRate(Waits{0: 2, 9: 2}, nil), eps)
	assert.InDelta(t, 71.058, CalculateAvgAgariRate(Waits{0: 3, 1: 4}, nil), eps)
	assert.InDelta(t, 71.058, CalculateAvgAgariRate(Waits{8: 3, 7: 4}, nil), eps)
	assert.InDelta(t, 96.2222, CalculateAvgAgariRate(Waits{0: 1, 1: 3, 2: 3, 3: 3, 4: 3, 5: 3, 6: 3, 7: 3, 9: 1}, nil), eps)
	assert.InDelta(t, 75.472, CalculateAvgAgariRate(Waits{9: 2, 27: 2}, nil), eps)
	assert.InDelta(t, 49.5, CalculateAvgAgariRate(Waits{27: 3}, nil), eps)
	assert.InDelta(t, 58, CalculateAvgAgariRate(Waits{27: 2}, nil), eps)
	assert.InDelta(t, 47.5, CalculateAvgAgariRate(Waits{27: 1}, nil), eps)
	assert.InDelta(t, 0, CalculateAvgAgariRate(Waits{27: 0}, nil), eps)
}
