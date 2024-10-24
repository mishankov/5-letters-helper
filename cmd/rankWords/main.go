package main

import (
	"fiveLettersHelper/internal/words"
	"log"
)

// TODO: which is better: higher or lower score? Test for every word, which is faster
// TODO: which is better: overall or remaining words score? Test for every word, which is faster
func main() {
	fiveLettersWords, err := words.GetFiveLettersWords()
	if err != nil {
		log.Fatal("Error getting words:", err)
	}

	// for index, order := range []int{-1, 1} {
	// 	var amounts []int
	// 	totalAmount := 0

	// 	for _, targetWord := range fiveLettersWords {
	// 		remainigVariants := words.RankWords(fiveLettersWords, order)
	// 		status := ""
	// 		amount := 0

	// 		for status != "22222" {
	// 			amount++

	// 			guess := remainigVariants[0].Word
	// 			status = game.GetWordStatus(guess, targetWord)
	// 		}

	// 		amounts = append(amounts, amount)
	// 	}
	// }

	wordScores := words.RankWords(fiveLettersWords, -1)

	log.Println(words.GetFirstNWords(wordScores, 10))
}
