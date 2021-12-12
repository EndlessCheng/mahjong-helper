package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_checkNewVersion(t *testing.T) {
	latestVersionTag, err := fetchLatestVersionTag()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, latestVersionTag)
}
