package main

import (
	"fiveLettersHelper/internal/words"
	"io/fs"
	"log"
	"os"
)

func main() {
	words, err := words.GetAllWords()
	if err != nil {
		log.Fatal(err)
	}

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
