package main

import (
	"strings"
	"testing"
)

func TestGetResults(t *testing.T) {
	input := `0 3 6 9 12 15
1 3 6 10 15 21
10 13 16 21 30 45
`

	result := getResult(strings.NewReader(input))
	if result != 114 {
		t.Errorf("Expected result to be 114, got %d", result)
	}
}

func BenchmarkGetResult(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		input := `0 3 6 9 12 15
1 3 6 10 15 21
10 13 16 21 30 45
`

		for n := 0; n < b.N; n++ {
			getResult(strings.NewReader(input))
		}
	})

	b.Run("large", func(b *testing.B) {
		inputFile := loadFile()
		defer inputFile.Close()

		for n := 0; n < b.N; n++ {
			getResult(inputFile)
		}
	})
}
