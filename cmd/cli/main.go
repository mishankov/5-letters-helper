package main

import (
	"fiveLettersHelper/internal/dbUtils"
	"fiveLettersHelper/internal/game"
	"fiveLettersHelper/internal/user"
	wordsUtils "fiveLettersHelper/internal/words"
	"fmt"
	"log"
	"slices"

	"github.com/mishankov/go-utlz/cliutils"
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

	game.CancelAllGamesForUser(user.Id, db)

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
			word, err = cliutils.UserInput("Введи слово: ")
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
			result, err = cliutils.UserInput("Введи результат (0, 1, 2): ")
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

		guesses, err := game.GetGuesses(db)
		if err != nil {
			log.Fatal("Error getting guesses for game:", err)
		}

		newWords, additionalResults, err := game.FilterWords(words, guesses)
		if err != nil {
			log.Fatal("Error filtering words:", err)
		}

		if len(newWords) == 1 {
			fmt.Printf("Игра закончена. Загаданное слово: %v\n", newWords[0])
			err = game.Complete(db)
			if err != nil {
				log.Fatal("Error setting game status to 'complete':", err)
			}

			cliutils.UserInput("Нажми ENTER, чтобы закрыть окно...")
			break
		}

		if len(newWords) == 0 {
			fmt.Printf("Игра закончена. Не нашлось слов, удовлетворяющих условиям\n")
			err = game.Fail(db)
			if err != nil {
				log.Fatal("Error setting game status to 'failed':", err)
			}
			cliutils.UserInput("Нажми ENTER, чтобы закрыть окно...")
			break
		}

		words = wordsUtils.GetFirstNWords(wordsUtils.RankWords(newWords, 1), len(newWords))

		fmt.Printf("Осталось %v слов для выбора. Первые из них: %v\n", len(words), cliutils.FormatListWithSeparator(words[:min(len(words), 10)], ", "))
		fmt.Printf("Известные положения букв: %v\n", cliutils.FormatListWithSeparator(additionalResults.LetterPositions, " "))
		fmt.Printf("Неиспользуемые буквы: %v\n", cliutils.FormatListWithSeparator(additionalResults.UnwantedLetters, ", "))
	}

	err = game.Cancel(db)
	if err != nil {
		log.Fatal("Error setting game status to 'cancel':", err)
	}
}
