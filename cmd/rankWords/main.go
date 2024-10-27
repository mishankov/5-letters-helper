package main

import (
	"fiveLettersHelper/internal/game"
	"fiveLettersHelper/internal/guess"
	"fiveLettersHelper/internal/words"
	"fiveLettersHelper/pkg/logging"
)

var logger = logging.NewLogger("rankWords")

// TODO: which is better: higher or lower score? Test for every word, which is faster
// TODO: which is better: overall or remaining words score? Test for every word, which is faster
func main() {
	var logger = logging.NewLoggerFromParent("main", &logger)
	fiveLettersWords, err := words.GetFiveLettersWords()
	if err != nil {
		logger.Fatalf("Error getting words: %v", err)
	}

	for _, order := range []int{-1, 1} {
		totalAmount := 0

		for _, targetWord := range fiveLettersWords {
			remainigVariants := words.RankWords(fiveLettersWords, order)
			status := ""
			amount := 0
			guesses := []guess.Guess{}

			for status != "22222" {
				amount++

				if len(remainigVariants) == 0 {
					logger.Infof("No remaining variants for order %v, target %v", order, targetWord)
					break
				}

				guessWord := remainigVariants[0].Word
				status = game.GetWordStatus(guessWord, targetWord)

				guesses = append(guesses, guess.Guess{Word: guessWord, Result: status})

				onlyWords, _, _ := game.FilterWords(words.GetFirstNWords(remainigVariants, len(remainigVariants)), guesses)
				remainigVariants = words.RankWords(onlyWords, order)
			}

			totalAmount += amount
		}

		logger.Infof("%v %v\n", order, totalAmount)
	}

	wordScores := words.RankWords(fiveLettersWords, -1)

	logger.Info(words.GetFirstNWords(wordScores, 10))
}
