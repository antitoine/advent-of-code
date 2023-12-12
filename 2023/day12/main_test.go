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
		if result != 525152 {
			t.Errorf("Expected result to be 525152, got %d", result)
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
		if result != 16384 {
			t.Errorf("Expected result to be 16384, got %d, for line %s", result, input)
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
		if result != 16 {
			t.Errorf("Expected result to be 16, got %d, for line %s", result, input)
		}
	})
	t.Run("line 5", func(t *testing.T) {
		input := `????.######..#####. 1,6,5
`

		result := getResult(strings.NewReader(input))
		if result != 2500 {
			t.Errorf("Expected result to be 2500, got %d, for line %s", result, input)
		}
	})
	t.Run("line 6", func(t *testing.T) {
		input := `?###???????? 3,2,1
`

		result := getResult(strings.NewReader(input))
		if result != 506250 {
			t.Errorf("Expected result to be 506250, got %d, for line %s", result, input)
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
