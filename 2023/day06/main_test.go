package main

import (
	"strings"
	"testing"
)

func TestGetSumOfWinningCards(t *testing.T) {
	input := `Time:      7  15   30
Distance:  9  40  200
`

	result := getResult(strings.NewReader(input))
	if result != 288 {
		t.Errorf("Expected result to be 288, got %d", result)
	}
}
