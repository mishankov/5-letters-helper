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
