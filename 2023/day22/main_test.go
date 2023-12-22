package main

import (
	"strings"
	"testing"
)

const testingInput = `1,0,1~1,2,1
0,0,2~2,0,2
0,2,3~2,2,3
0,0,4~0,2,4
2,0,5~2,2,5
0,1,6~2,1,6
1,1,8~1,1,9
`

func TestGetResults(t *testing.T) {
	t.Run("part1", func(t *testing.T) {
		const testingExpectedResult = 5
		result := getResultPart1(strings.NewReader(testingInput))
		if result != testingExpectedResult {
			t.Errorf("Expected result to be %d, got %d", testingExpectedResult, result)
		}
	})
	t.Run("part2", func(t *testing.T) {
		const testingExpectedResult = 7
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
