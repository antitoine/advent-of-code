package main

import (
	"strings"
	"testing"
)

const testingInput = `Register A: 729
Register B: 0
Register C: 0

Program: 0,1,5,4,3,0
`

const testingExpectedResult = "4,6,3,5,6,3,5,2,1,0"

func TestGetResults(t *testing.T) {
	result := getResult(strings.NewReader(testingInput))
	if result != testingExpectedResult {
		t.Errorf("Expected result to be %s, got %s", testingExpectedResult, result)
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
