package game

import (
	"database/sql"
	"fiveLettersHelper/internal/guess"
	wordsUtils "fiveLettersHelper/internal/words"
	"log"
	"slices"

	"github.com/google/uuid"
)

type Game struct {
	Id     string
	User   string
	Status string
}

func NewGame(user string, db *sql.DB) (Game, error) {
	game := Game{Id: uuid.NewString(), User: user, Status: "new"}
	_, err := db.Exec("INSERT INTO game (id, user, status) VALUES (?, ? ,?)", game.Id, game.User, game.Status)
	if err != nil {
		return Game{}, err
	}

	return game, nil
}

func (g *Game) InProgress(db *sql.DB) error {
	g.Status = "in progress"
	_, err := db.Exec("UPDATE game SET status = ? WHERE id = ?", g.Status, g.Id)
	return err
}

func (g *Game) Complete(db *sql.DB) error {
	if g.Status != "cancelled" && g.Status != "completed" {
		g.Status = "completed"
		_, err := db.Exec("UPDATE game SET status = ? WHERE id = ?", g.Status, g.Id)
		return err
	}
	return nil
}

func (g *Game) Cancel(db *sql.DB) error {
	if g.Status != "cancelled" && g.Status != "completed" {
		g.Status = "canceled"
		_, err := db.Exec("UPDATE game SET status = ? WHERE id = ?", g.Status, g.Id)
		return err
	}

	return nil
}

func (g *Game) GetGuesses(db *sql.DB) ([]guess.Guess, error) {
	return guess.GetGuesseForGame(g.Id, db)
}

func (g *Game) NewGuess(number int, word string, result string, db *sql.DB) (guess.Guess, error) {
	return guess.NewGuess(g.Id, number, word, result, db)
}

type FWAdditionalResults struct {
	LetterPositions []rune
	UnwantedLetters []rune
	WrongPositions  map[int][]rune
	AmountOfLetters map[rune]int
}

func (g *Game) FilterWords(words []string, db *sql.DB) (filteredWords []string, additionalResults FWAdditionalResults, error error) {
	guesses, err := g.GetGuesses(db)
	if err != nil {
		return []string{}, FWAdditionalResults{}, err
	}

	letterPositions := []rune{'_', '_', '_', '_', '_'}
	unwantedLetters := []rune{}
	wrongPositions := map[int][]rune{}
	amountOfLetters := map[rune]int{}

	for _, guess := range guesses {
		localAmountOfLetters := map[rune]int{}
		for i := 0; i < 5; i++ {
			currentResult := []rune(guess.Result)[i]
			currentLetter := []rune(guess.Word)[i]

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

		words = newWords
	}

	return words, FWAdditionalResults{LetterPositions: letterPositions, UnwantedLetters: unwantedLetters, WrongPositions: wrongPositions, AmountOfLetters: amountOfLetters}, nil
}
