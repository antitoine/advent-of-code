package main

import (
	"strings"
	"testing"
)

const testingInput = `19, 13, 30 @ -2,  1, -2
18, 19, 22 @ -1, -1, -2
20, 25, 34 @ -2, -2, -4
12, 31, 28 @ -1, -2, -1
20, 19, 15 @  1, -5, -3
`

var testingTestZone = Zone{
	min: Coordinates{
		x: 7,
		y: 7,
	},
	max: Coordinates{
		x: 27,
		y: 27,
	},
}

var finalTestZone = Zone{
	min: Coordinates{
		x: 200000000000000,
		y: 200000000000000,
	},
	max: Coordinates{
		x: 400000000000000,
		y: 400000000000000,
	},
}

func TestGetResults(t *testing.T) {
	t.Run("part1", func(t *testing.T) {
		t.Run("small", func(t *testing.T) {
			const testingExpectedResult = 2
			result := GetResultPart1(strings.NewReader(testingInput), testingTestZone)
			if result != testingExpectedResult {
				t.Errorf("Expected result to be %d, got %d", testingExpectedResult, result)
			}
		})
		t.Run("large", func(t *testing.T) {
			const finalResult = 17244
			inputFile := loadFile()
			defer inputFile.Close()
			result := GetResultPart1(inputFile, finalTestZone)
			if result != finalResult {
				t.Errorf("Expected result to be %d, got %d", finalResult, result)
			}
		})
	})
	t.Run("part2", func(t *testing.T) {
		t.Run("small", func(t *testing.T) {
			const testingExpectedResult = 47
			result := GetResultPart2(strings.NewReader(testingInput))
			if result != testingExpectedResult {
				t.Errorf("Expected result to be %d, got %d", testingExpectedResult, result)
			}
		})
	})
}

func BenchmarkGetResult(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			GetResultPart1(strings.NewReader(testingInput), testingTestZone)
		}
	})

	b.Run("large", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			inputFile := loadFile()
			GetResultPart1(inputFile, finalTestZone)
			inputFile.Close()
		}
	})
}
