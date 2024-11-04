package words

import (
	"fiveLettersHelper/pkg/logging"
	"os"
	"slices"
	"strings"
)

var logger = logging.NewLogger("words")

func GetAllWords() ([]string, error) {
	originalData, err := os.ReadFile("./data/russian_nouns.txt")
	if err != nil {
		return []string{}, err
	}

	return strings.Split(string(originalData), "\n"), nil
}

func GetFiveLettersWords() ([]string, error) {
	originalData, err := os.ReadFile("./data/five_letters_russian_nouns.txt")
	if err != nil {
		return []string{}, err
	}

	return strings.Split(string(originalData), "\n"), nil
}

func WordRemains(word string, unwantedLetters []rune, unwantedWords []string, letterPositions []rune, amountOfLetters map[rune]int, wrongPositions map[int][]rune) bool {
	// logger := logging.NewLoggerFromParent("WordRemains", &logger)

	logger.Debugf("Word: %v", word)

	if slices.Contains(unwantedWords, word) {
		logger.Debug("Is in unwanted words")
		return false
	}

	for _, letter := range unwantedLetters {
		if strings.ContainsRune(word, letter) {
			logger.Debugf("Has unwanted letter %q", letter)
			return false
		}
	}

	for index, letter := range letterPositions {
		if letter != '_' && []rune(word)[index] != letter {
			logger.Debugf("Does not have letter %q at position %v", letter, index)
			return false
		}
	}

	for letter, amount := range amountOfLetters {
		if strings.Count(word, string(letter)) < amount {
			logger.Debugf("Does not have %v amount of letter %q", amount, letter)
			return false
		}
	}

	for position, letters := range wrongPositions {
		for _, letter := range letters {
			if []rune(word)[position] == letter {
				logger.Debugf("Has letter %q at known wrong position %v", letter, position)
				return false
			}
		}
	}

	logger.Debug("Remains")
	return true
}

type WordScore struct {
	Word  string
	Score int
}

func GetFirstNWords(ws []WordScore, n int) []string {
	result := []string{}

	for i, w := range ws {
		result = append(result, w.Word)

		if i == n-1 {
			return result
		}

	}

	return result
}

func RankWords(words []string, order int) []WordScore {
	lettersCount := map[rune]int{}

	for _, word := range words {
		for _, letter := range word {
			lettersCount[letter] += 1
		}
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

		wordScores = append(wordScores, WordScore{Word: word, Score: score})
	}

	slices.SortFunc(wordScores, func(w1, w2 WordScore) int {
		if w1.Score < w2.Score {
			return order
		}

		if w1.Score > w2.Score {
			return -order
		}

		return 0
	})

	return wordScores
}
