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

const testingExpectedResult = 2

var testZone = Zone{
	min: Coordinates{
		x: 7,
		y: 7,
	},
	max: Coordinates{
		x: 27,
		y: 27,
	},
}

func TestGetResults(t *testing.T) {
	result := getResult(strings.NewReader(testingInput), testZone)
	if result != testingExpectedResult {
		t.Errorf("Expected result to be %d, got %d", testingExpectedResult, result)
	}
}

func BenchmarkGetResult(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			getResult(strings.NewReader(testingInput), testZone)
		}
	})

	b.Run("large", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			inputFile := loadFile()
			getResult(inputFile, testZone)
			inputFile.Close()
		}
	})
}
