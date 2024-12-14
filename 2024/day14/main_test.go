package main

import (
	"strings"
	"testing"
)

const testingInput = `p=0,4 v=3,-3
p=6,3 v=-1,-3
p=10,3 v=-1,2
p=2,0 v=2,-1
p=0,0 v=1,3
p=3,0 v=-2,-2
p=7,6 v=-1,-3
p=3,0 v=-1,-2
p=9,3 v=2,3
p=7,3 v=-1,2
p=2,4 v=2,-3
p=9,5 v=-3,-3
`

const testingExpectedResult = 12

const testingSizeX = 11
const testingSizeY = 7
const testingNbSeconds = 100

func TestGetResults(t *testing.T) {
	result := getResult(strings.NewReader(testingInput), testingSizeX, testingSizeY, testingNbSeconds)
	if result != testingExpectedResult {
		t.Errorf("Expected result to be %d, got %d", testingExpectedResult, result)
	}
}

func BenchmarkGetResult(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			getResult(strings.NewReader(testingInput), testingSizeX, testingSizeY, testingNbSeconds)
		}
	})

	b.Run("large", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			inputFile := loadFile()
			getResult(inputFile, testingSizeX, testingSizeY, testingNbSeconds)
			inputFile.Close()
		}
	})
}
