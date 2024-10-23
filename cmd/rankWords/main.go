package main

import (
	"fiveLettersHelper/internal/words"
	"log"
)

// TODO: which is better: higher or lower score? Test for every word, which is faster
// TODO: which is better: overall or remainong words score? Test for every word, which is faster
func main() {
	fiveLettersWords, err := words.GetFiveLettersWords()
	if err != nil {
		log.Fatal("Error getting words:", err)
	}

	wordScores := words.RankWords(fiveLettersWords)

	log.Println(words.GetFirstNWords(wordScores, 10))
}
