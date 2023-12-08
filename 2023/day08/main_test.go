package main

import (
	"strings"
	"testing"
)

func TestGetResults(t *testing.T) {
	t.Run("Example 1", func(t *testing.T) {
		input := `RL

AAA = (BBB, CCC)
BBB = (DDD, EEE)
CCC = (ZZZ, GGG)
DDD = (DDD, DDD)
EEE = (EEE, EEE)
GGG = (GGG, GGG)
ZZZ = (ZZZ, ZZZ)
`

		result := getResult(strings.NewReader(input))
		if result != 2 {
			t.Errorf("Expected result to be 2, got %d", result)
		}
	})
	t.Run("Example 2", func(t *testing.T) {
		input := `LLR

AAA = (BBB, BBB)
BBB = (AAA, ZZZ)
ZZZ = (ZZZ, ZZZ)
`

		result := getResult(strings.NewReader(input))
		if result != 6 {
			t.Errorf("Expected result to be 6, got %d", result)
		}
	})
}
