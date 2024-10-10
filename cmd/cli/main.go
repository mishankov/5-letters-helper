package main

import (
	"fiveLettersHelper/internal/db"
	"fiveLettersHelper/internal/game"
	"fiveLettersHelper/internal/guess"
	"fiveLettersHelper/internal/user"
	wordsUtils "fiveLettersHelper/internal/words"
	"fiveLettersHelper/packages/cliUtils"
	"fmt"
	"log"
	"slices"
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

	words, err := wordsUtils.GetFiveLettersWords()

	if err != nil {
		log.Fatal("Error getting words from file:", err)
	}

	err = game.InProgress(db)
	if err != nil {
		log.Fatal("Error setting game status to 'in progress':", err)
	}

	turnNumber := 0
	letterPositions := []rune{'_', '_', '_', '_', '_'}
	unwantedLetters := []rune{}
	wrongPositions := map[int][]rune{}
	amountOfLetters := map[rune]int{}
	for {
		turnNumber++
		fmt.Printf("Ход №: %v\n", turnNumber)
		// TODO: format lists output better
		fmt.Printf("Осталось %v слов для выбора. Первые из них: %v\n", len(words), words[:min(len(words), 10)])
		fmt.Printf("Известные положения букв: %q\n", letterPositions)
		fmt.Printf("Неиспользуемые буквы: %q\n", unwantedLetters)

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

		localAmountOfLetters := map[rune]int{}
		for i := 0; i < 5; i++ {
			currentResult := []rune(result)[i]
			currentLetter := []rune(word)[i]

			switch currentResult {
			case '0':
				if !slices.Contains(unwantedLetters, currentLetter) {
					unwantedLetters = append(unwantedLetters, currentLetter)
				}
			case '1':
				localAmountOfLetters[currentLetter] += 1
				wrongPositions[i] = append(wrongPositions[i], currentLetter)
			case '2':
				letterPositions[i] = currentLetter
			default:
				log.Fatal("Unexpected result rune: ", string(currentResult))
			}
		}

		for letter, localAmount := range localAmountOfLetters {
			amount, ok := amountOfLetters[letter]
			if ok && amount < localAmount || !ok {
				amountOfLetters[letter] = localAmount
			}
		}

		for i, letter := range unwantedLetters {
			if _, ok := amountOfLetters[letter]; ok || slices.Contains(letterPositions, letter) {
				unwantedLetters[i] = unwantedLetters[len(unwantedLetters)-1]
				unwantedLetters = unwantedLetters[:len(unwantedLetters)-1]
			}
		}

		newWords := []string{}

		for _, word := range words {
			if len(word) == 0 {
				continue
			}

			if wordsUtils.WordRemains(word, unwantedLetters, letterPositions, amountOfLetters, wrongPositions) {
				newWords = append(newWords, word)
			}
		}

		_, err = guess.NewGuess(game.Id, turnNumber, word, result, db)
		if err != nil {
			log.Fatal("Error creating guess:", err)
		}

		if len(newWords) == 1 {
			fmt.Printf("Игра закончена. Загаданное слово: %v\n", newWords[0])
			err = game.Complete(db)
			if err != nil {
				log.Fatal("Error setting game status to 'complete':", err)
			}
			break
		}

		words = newWords
	}

	err = game.Cancel(db)
	if err != nil {
		log.Fatal("Error setting game status to 'cancel':", err)
	}
}
