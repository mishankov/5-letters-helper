package main

import (
	"fiveLettersHelper/internal/words"
	"log"
	"slices"
)

// TODO: move logic to external package
// TODO: claculate score between remaining words
// TODO: which is better: higher or lower score? Test for every word, which is faster
func main() {
	words, err := words.GetFiveLettersWords()
	if err != nil {
		log.Fatal("Error getting words:", err)
	}

	lettersCount := map[rune]int{}

	for _, word := range words {
		for _, letter := range word {
			lettersCount[letter] += 1
		}
	}

	type WordScore struct {
		word  string
		score int
	}
	wordScores := []WordScore{}

	for _, word := range words {
		uniqeLetters := []rune{}
		score := 0

		for _, letter := range word {
			if !slices.Contains(uniqeLetters, letter) {
				score += lettersCount[letter]
				uniqeLetters = append(uniqeLetters, letter)
			}
		}

		wordScores = append(wordScores, WordScore{word: word, score: score})
	}

	slices.SortFunc(wordScores, func(w1, w2 WordScore) int {
		if w1.score < w2.score {
			return 1
		}

		if w1.score > w2.score {
			return -1
		}

		return 0
	})

	log.Println(wordScores[:10])
}
