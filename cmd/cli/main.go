package main

import (
	dbUtils "fiveLettersHelper/internal/db"
	"fiveLettersHelper/internal/game"
	"fiveLettersHelper/internal/user"
	wordsUtils "fiveLettersHelper/internal/words"
	"fiveLettersHelper/packages/cliUtils"
	"fmt"
	"log"
)

func main() {
	db, err := dbUtils.GetDB()
	if err != nil {
		log.Fatal("Can't open database:", err)
	}
	defer db.Close()

	err = dbUtils.PrepareDB(db)
	if err != nil {
		log.Fatal("Error preparing DB:", err)
	}

	user, err := user.CreateAndGetCLIUser(db)
	if err != nil {
		log.Fatal("Error getting user:", err)
	}

	game, err := game.NewGame(user.Id, db)
	if err != nil {
		log.Fatal("Error creating new game:", err)
	}

	words, err := wordsUtils.GetFiveLettersWords()

	if err != nil {
		log.Fatal("Error getting words from file:", err)
	}

	err = game.InProgress(db)
	if err != nil {
		log.Fatal("Error setting game status to 'in progress':", err)
	}

	turnNumber := 0
	for {
		turnNumber++
		fmt.Printf("Ход №: %v\n", turnNumber)

		// TODO: validate user input
		word, err := cliUtils.UserInput("Введи слово: ")
		if err != nil {
			log.Fatal("Error getting word from user:", err)
		}

		// TODO: validate user input
		result, err := cliUtils.UserInput("Введи результат (0, 1, 2): ")
		if err != nil {
			log.Fatal("Error getting result from user:", err)
		}

		_, err = game.NewGuess(turnNumber, word, result, db)
		if err != nil {
			log.Fatal("Error creating guess:", err)
		}

		newWords, additionalResults, err := game.FilterWords(words, db)
		if err != nil {
			log.Fatal("Error filtering words:", err)
		}

		if len(newWords) == 1 {
			fmt.Printf("Игра закончена. Загаданное слово: %v\n", newWords[0])
			err = game.Complete(db)
			if err != nil {
				log.Fatal("Error setting game status to 'complete':", err)
			}

			cliUtils.UserInput("Нажми ENTER, чтобы закрыть окно...")
			break
		}

		words = newWords

		fmt.Printf("Осталось %v слов для выбора. Первые из них: %v\n", len(words), cliUtils.FormatListWithSeparator(words[:min(len(words), 10)], ", "))
		fmt.Printf("Известные положения букв: %v\n", cliUtils.FormatListWithSeparator(additionalResults.LetterPositions, " "))
		fmt.Printf("Неиспользуемые буквы: %v\n", cliUtils.FormatListWithSeparator(additionalResults.UnwantedLetters, ", "))
	}

	err = game.Cancel(db)
	if err != nil {
		log.Fatal("Error setting game status to 'cancel':", err)
	}
}
