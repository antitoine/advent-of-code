package main

import (
	"strings"
	"testing"
)

const testingInput = `3-5
10-14
16-20
12-18

1
5
8
11
17
32
`

const testingExpectedResult = 3

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
