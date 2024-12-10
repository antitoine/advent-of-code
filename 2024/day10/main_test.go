package main

import (
	"strings"
	"testing"
)

const testingInput = `89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732
`

const testingExpectedResult = 81

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
