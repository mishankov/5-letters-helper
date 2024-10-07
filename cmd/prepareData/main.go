package main

import (
	"io/fs"
	"log"
	"os"
	"strings"
)

func main() {
	originalData, errOpen := os.ReadFile("./data/russian_nouns.txt")
	if errOpen != nil {
		log.Panic(errOpen)
	}

	words := strings.Split(string(originalData), "\r\n")

	var fiveLettersWords string
	for _, word := range words {
		if len([]rune(word)) == 5 {
			fiveLettersWords = fiveLettersWords + word + "\r\n"
		}
	}

	errSave := os.WriteFile("./data/five_letters_russian_nouns.txt", []byte(fiveLettersWords), fs.ModeExclusive)
	if errSave != nil {
		log.Panic(errSave)
	}
}
