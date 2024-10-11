package main

import (
	dbUtils "fiveLettersHelper/internal/db"
	"fiveLettersHelper/internal/game"
	"fiveLettersHelper/internal/user"
	wordsUtils "fiveLettersHelper/internal/words"
	"fiveLettersHelper/packages/cliUtils"
	"fmt"
	"log"
	"slices"
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

		var word string
		for ok := false; !ok; {
			word, err = cliUtils.UserInput("Введи слово: ")
			if err != nil {
				log.Fatal("Error getting word from user:", err)
			}

			if len([]rune(word)) != 5 {
				fmt.Printf("В слове должно быть 5 букв. В слове %q их %v. Попробуй ещё раз\n", word, len([]rune(word)))
			} else {
				ok = true
			}
		}

		var result string
		for ok := false; !ok; {
			result, err = cliUtils.UserInput("Введи результат (0, 1, 2): ")
			if err != nil {
				log.Fatal("Error getting result from user:", err)
			}

			if len([]rune(result)) != 5 {
				fmt.Printf("В результате должно быть 5 символов. В результате %q их %v. Попробуй ещё раз\n", result, len([]rune(result)))
				continue
			}

			continueOuter := false
			for i, r := range result {
				if !slices.Contains([]rune{'0', '1', '2'}, r) {
					fmt.Printf("В результате должно быть только символы 0, 1 и 2. На %v позиции находится символ %q. Попробуй ещё раз\n", i+1, r)
					continueOuter = true
					break
				}
			}

			if continueOuter {
				continue
			}

			ok = true
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
