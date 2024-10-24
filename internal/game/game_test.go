package game

import (
	"fiveLettersHelper/internal/guess"
	"slices"
	"testing"
)

func TestFilterWords(t *testing.T) {
	tests := []struct {
		name    string
		words   []string
		guesses []guess.Guess
		result  []string
	}{
		{
			name:  "алмаз-саван-выпад",
			words: []string{"алмаз", "саван", "выпад"},
			guesses: []guess.Guess{
				{Word: "алмаз", Result: "10020"},
			},
			result: []string{"саван"},
		},
		{
			name:  "аборт-аббат",
			words: []string{"аборт", "аббат"},
			guesses: []guess.Guess{
				{Word: "аббат", Result: "22002"},
			},
			result: []string{"аборт"},
		},
	}

	for _, test := range tests {
		result, _, _ := FilterWords(test.words, test.guesses)

		for _, want := range test.result {
			if !slices.Contains(result, want) {
				t.Fatalf("Test: %q. %q should be in results", test.name, want)
			}
		}

		for _, got := range result {
			if !slices.Contains(test.result, got) {
				t.Fatalf("Test: %q. %q should not be in results", test.name, got)
			}
		}
	}
}
