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
	t.Run("Example 3", func(t *testing.T) {
		input := `LR

11A = (11B, XXX)
11B = (XXX, 11Z)
11Z = (11B, XXX)
22A = (22B, XXX)
22B = (22C, 22C)
22C = (22Z, 22Z)
22Z = (22B, 22B)
XXX = (XXX, XXX)
`

		result := getResult(strings.NewReader(input))
		if result != 6 {
			t.Errorf("Expected result to be 6, got %d", result)
		}
	})
}
