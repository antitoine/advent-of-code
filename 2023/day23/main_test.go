package main

import (
	"strings"
	"testing"
)

const testingInput = `#.#####################
#.......#########...###
#######.#########.#.###
###.....#.>.>.###.#.###
###v#####.#v#.###.#.###
###.>...#.#.#.....#...#
###v###.#.#.#########.#
###...#.#.#.......#...#
#####.#.#.#######.#.###
#.....#.#.#.......#...#
#.#####.#.#.#########v#
#.#...#...#...###...>.#
#.#.#v#######v###.###v#
#...#.>.#...>.>.#.###.#
#####v#.#.###v#.#.###.#
#.....#...#...#.#.#...#
#.#########.###.#.#.###
#...###...#...#...#.###
###.###.#.###v#####v###
#...#...#.#.>.>.#.>.###
#.###.###.#.###.#.#v###
#.....###...###...#...#
#####################.#
`

func TestGetResults(t *testing.T) {
	t.Run("part1", func(t *testing.T) {
		const testingExpectedResult = 94
		result := getResultPart1(strings.NewReader(testingInput))
		if result != testingExpectedResult {
			t.Errorf("Expected result to be %d, got %d", testingExpectedResult, result)
		}
	})
	t.Run("part2", func(t *testing.T) {
		const testingExpectedResult = 154
		result := getResultPart2(strings.NewReader(testingInput))
		if result != testingExpectedResult {
			t.Errorf("Expected result to be %d, got %d", testingExpectedResult, result)
		}
	})
}

func BenchmarkGetResult(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			getResultPart1(strings.NewReader(testingInput))
		}
	})

	b.Run("large", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			inputFile := loadFile()
			getResultPart1(inputFile)
			inputFile.Close()
		}
	})
}
