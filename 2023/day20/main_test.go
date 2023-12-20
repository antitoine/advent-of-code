package main

import (
	"strings"
	"testing"
)

const testingInput1 = `broadcaster -> a, b, c
%a -> b
%b -> c
%c -> inv
&inv -> a
`

const testingInput2 = `broadcaster -> a
%a -> inv, con
&inv -> b
%b -> con
&con -> output
`

func TestGetResultForPart1(t *testing.T) {
	t.Run("first", func(t *testing.T) {
		const testingExpectedResult1 = 32000000
		result := getResultForPart1(strings.NewReader(testingInput1))
		if result != testingExpectedResult1 {
			t.Errorf("Expected result to be %d, got %d", testingExpectedResult1, result)
		}
	})
	t.Run("second", func(t *testing.T) {
		const testingExpectedResult2 = 11687500
		result := getResultForPart1(strings.NewReader(testingInput2))
		if result != testingExpectedResult2 {
			t.Errorf("Expected result to be %d, got %d", testingExpectedResult2, result)
		}
	})
}

func BenchmarkGetResult(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			getResultForPart1(strings.NewReader(testingInput1))
		}
	})

	b.Run("large", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			inputFile := loadFile()
			getResultForPart1(inputFile)
			inputFile.Close()
		}
	})
}
