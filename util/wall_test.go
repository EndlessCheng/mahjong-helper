package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCalcWallTiles(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("189s", CalcWallTiles(InitLeftTiles34WithTiles34(MustStrToTiles34("2222777888s"))).String())
	assert.Equal("12589s", CalcWallTiles(InitLeftTiles34WithTiles34(MustStrToTiles34("33337777s"))).String())
	assert.Equal("12589s", CalcWallTiles(InitLeftTiles34WithTiles34(MustStrToTiles34("333777s"))).String())
	assert.Equal("1235689s", CalcWallTiles(InitLeftTiles34WithTiles34(MustStrToTiles34("333444777s"))).String())
	assert.Equal("9s", CalcWallTiles(InitLeftTiles34WithTiles34(MustStrToTiles34("8888s"))).String())
}

func TestCalcDNCSafeTiles(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("9s", CalcDNCSafeTiles(InitLeftTiles34WithTiles34(MustStrToTiles34("8888s"))).String())
	assert.Equal("1245s", CalcDNCSafeTiles(InitLeftTiles34WithTiles34(MustStrToTiles34("33336666s"))).String())
	assert.Equal("124s", CalcDNCSafeTiles(InitLeftTiles34WithTiles34(MustStrToTiles34("33335555s"))).String())
	assert.Equal("1289s", CalcDNCSafeTiles(InitLeftTiles34WithTiles34(MustStrToTiles34("33337777s"))).String())
	assert.Equal("124689s", CalcDNCSafeTiles(InitLeftTiles34WithTiles34(MustStrToTiles34("333355557777s"))).String())
}
