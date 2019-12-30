package util

import (
	"fmt"
	"testing"
)

func TestGetTenpaiRate3(t *testing.T) {
	for _, data := range tenpaiRate[1:] {
		for turn, turnData := range data {
			fmt.Println(turn, GetTenpaiRate3(turnData[0]))
		}
		fmt.Println()
	}
}
