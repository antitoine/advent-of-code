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
	if result != 71503 {
		t.Errorf("Expected result to be 71503, got %d", result)
	}
}
