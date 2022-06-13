package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_checkNewVersion(t *testing.T) {
	latestVersionTag, err := FetchLatestVersionTag()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEmpty(t, latestVersionTag)
}
