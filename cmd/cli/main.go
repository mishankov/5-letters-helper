package main

import (
	"fiveLettersHelper/internal/db"
	"fiveLettersHelper/internal/game"
	"fiveLettersHelper/internal/user"
	"fiveLettersHelper/internal/words"
	"fiveLettersHelper/packages/cliUtils"
	"fmt"
	"log"
)

func main() {
	db, err := db.GetDB()
	if err != nil {
		log.Fatal("Can't open database:", err)
	}
	defer db.Close()

	user, err := user.CreateAndGetCLIUser(db)
	if err != nil {
		log.Fatal("Error getting user:", err)
	}

	game, err := game.NewGame(user.Id, db)
	if err != nil {
		log.Fatal("Error creating new game:", err)
	}

	words, err := words.GetFiveLettersWords()

	if err != nil {
		log.Fatal("Error getting words from file:", err)
	}

	turnNumber := 0
	letterPositions := []rune{'_', '_', '_', '_', '_'}
	unwantedLetters := []rune{}
	game.InProgress()
	for {
		turnNumber++
		fmt.Printf("Ход №: %v\n", turnNumber)
		fmt.Printf("Осталось %v слов для выбора. Первые из них: %v\n", len(words), words[:10])
		fmt.Printf("Известные положения букв: %q\n", letterPositions)
		fmt.Printf("Неиспользуемые буквы: %q\n", unwantedLetters)

		word, err := cliUtils.UserInput("Введи слово: ")
		if err != nil {
			log.Fatal("Error getting word from user:", err)
		}

		result, err := cliUtils.UserInput("Введи результат (0, 1, 2): ")
		if err != nil {
			log.Fatal("Error getting result from user:", err)
		}

		log.Println(word, result)

		break
	}

	game.Cancel()
}
