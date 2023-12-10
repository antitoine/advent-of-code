package main

import (
	"strings"
	"testing"
)

func TestGetResults(t *testing.T) {
	t.Run("example 1", func(t *testing.T) {
		input := `.....
.S-7.
.|.|.
.L-J.
.....
`

		result := getResult(strings.NewReader(input))
		if result != 4 {
			t.Errorf("Expected result to be 4, got %d", result)
		}
	})
	t.Run("example 2", func(t *testing.T) {
		input := `..F7.
.FJ|.
SJ.L7
|F--J
LJ...
`

		result := getResult(strings.NewReader(input))
		if result != 8 {
			t.Errorf("Expected result to be 8, got %d", result)
		}
	})
}

func BenchmarkGetResult(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		input := `..F7.
.FJ|.
SJ.L7
|F--J
LJ...
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
