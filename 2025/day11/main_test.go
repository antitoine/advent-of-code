package main

import (
	"strings"
	"testing"
)

const testingInput = `svr: aaa bbb
aaa: fft
fft: ccc
bbb: tty
tty: ccc
ccc: ddd eee
ddd: hub
hub: fff
eee: dac
dac: fff
fff: ggg hhh
ggg: out
hhh: out
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
