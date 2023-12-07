package main

import (
	"strings"
	"testing"
)

func TestGetSumOfWinningCards(t *testing.T) {
	input := `32T3K 765
T55J5 684
KK677 28
KTJJT 220
QQQJA 483`

	result := getResult(strings.NewReader(input))
	if result != 6440 {
		t.Errorf("Expected result to be 6440, got %d", result)
	}
}
