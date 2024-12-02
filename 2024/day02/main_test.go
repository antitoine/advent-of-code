package main

import (
	"strings"
	"testing"
)

const testingInput = `7 6 4 2 1
1 2 7 8 9
9 7 6 2 1
1 3 2 4 5
8 6 4 4 1
1 3 6 7 9
`

const testingExpectedResult = 2

func TestGetResults(t *testing.T) {
	result := getResult(strings.NewReader(testingInput))
	if result != testingExpectedResult {
		t.Errorf("Expected result to be %d, got %d", testingExpectedResult, result)
	}
}

func BenchmarkGetResult(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			getResult(strings.NewReader(testingInput))
		}
	})

	b.Run("large", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			inputFile := loadFile()
			getResult(inputFile)
			inputFile.Close()
		}
	})
}
