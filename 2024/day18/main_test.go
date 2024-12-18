package main

import (
	"image"
	"strings"
	"testing"
)

const testingInput = `5,4
4,2
4,5
3,0
2,1
6,3
2,4
1,5
0,6
3,3
2,6
5,1
1,2
5,5
2,5
6,5
1,4
0,4
6,4
1,1
6,1
1,0
0,5
1,6
2,0
`

var testingInputSpace = image.Rectangle{Min: image.Pt(0, 0), Max: image.Pt(7, 7)}

const testingInputNbCorruptedBytes = 12

const testingExpectedResult = 22

func TestGetResults(t *testing.T) {
	result := getResult(strings.NewReader(testingInput), testingInputSpace, testingInputNbCorruptedBytes)
	if result != testingExpectedResult {
		t.Errorf("Expected result to be %d, got %d", testingExpectedResult, result)
	}
}

func BenchmarkGetResult(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			getResult(strings.NewReader(testingInput), testingInputSpace, testingInputNbCorruptedBytes)
		}
	})

	b.Run("large", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			inputFile := loadFile()
			getResult(inputFile, image.Rectangle{Min: image.Pt(0, 0), Max: image.Pt(69, 69)}, 1024)
			inputFile.Close()
		}
	})
}
