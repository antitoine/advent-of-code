package main

import (
	"strings"
	"testing"
)

func TestGetResults(t *testing.T) {
	t.Run("full", func(t *testing.T) {
		input := `???.### 1,1,3
.??..??...?##. 1,1,3
?#?#?#?#?#?#?#? 1,3,1,6
????.#...#... 4,1,1
????.######..#####. 1,6,5
?###???????? 3,2,1
`

		result := getResult(strings.NewReader(input))
		if result != 21 {
			t.Errorf("Expected result to be 21, got %d", result)
		}
	})
	t.Run("line 1", func(t *testing.T) {
		input := `???.### 1,1,3
`

		result := getResult(strings.NewReader(input))
		if result != 1 {
			t.Errorf("Expected result to be 1, got %d, for line %s", result, input)
		}
	})
	t.Run("line 2", func(t *testing.T) {
		input := `.??..??...?##. 1,1,3
`

		result := getResult(strings.NewReader(input))
		if result != 4 {
			t.Errorf("Expected result to be 4, got %d, for line %s", result, input)
		}
	})
	t.Run("line 3", func(t *testing.T) {
		input := `?#?#?#?#?#?#?#? 1,3,1,6
`

		result := getResult(strings.NewReader(input))
		if result != 1 {
			t.Errorf("Expected result to be 1, got %d, for line %s", result, input)
		}
	})
	t.Run("line 4", func(t *testing.T) {
		input := `????.#...#... 4,1,1
`

		result := getResult(strings.NewReader(input))
		if result != 1 {
			t.Errorf("Expected result to be 1, got %d, for line %s", result, input)
		}
	})
	t.Run("line 5", func(t *testing.T) {
		input := `????.######..#####. 1,6,5
`

		result := getResult(strings.NewReader(input))
		if result != 4 {
			t.Errorf("Expected result to be 4, got %d, for line %s", result, input)
		}
	})
	t.Run("line 6", func(t *testing.T) {
		input := `?###???????? 3,2,1
`

		result := getResult(strings.NewReader(input))
		if result != 10 {
			t.Errorf("Expected result to be 10, got %d, for line %s", result, input)
		}
	})
}

func BenchmarkGetResult(b *testing.B) {
	b.Run("small", func(b *testing.B) {
		input := `???.### 1,1,3
.??..??...?##. 1,1,3
?#?#?#?#?#?#?#? 1,3,1,6
????.#...#... 4,1,1
????.######..#####. 1,6,5
?###???????? 3,2,1`

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
