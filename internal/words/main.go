package words

import (
	"os"
	"strings"
)

func GetAllWords() ([]string, error) {
	originalData, err := os.ReadFile("./data/russian_nouns.txt")
	if err != nil {
		return []string{}, err
	}

	return strings.Split(string(originalData), "\r\n"), nil
}

func GetFiveLettersWords() ([]string, error) {
	originalData, err := os.ReadFile("./data/five_letters_russian_nouns.txt")
	if err != nil {
		return []string{}, err
	}

	return strings.Split(string(originalData), "\r\n"), nil
}

func WordRemains(word string, unwantedLetters []rune, letterPositions []rune, amountOfLetters map[rune]int, wrongPositions map[int][]rune) bool {
	for _, letter := range unwantedLetters {
		if strings.ContainsRune(word, letter) {
			return false
		}
	}

	for index, letter := range letterPositions {
		if letter != '_' && []rune(word)[index] != letter {
			return false
		}
	}

	for letter, amount := range amountOfLetters {
		if strings.Count(word, string(letter)) < amount {
			return false
		}
	}

	for position, letters := range wrongPositions {
		for _, letter := range letters {
			if []rune(word)[position] == letter {
				return false
			}
		}
	}

	return true
}
