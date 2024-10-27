package game

import (
	"database/sql"
	"fiveLettersHelper/internal/guess"
	wordsUtils "fiveLettersHelper/internal/words"
	"fiveLettersHelper/pkg/logging"
	"log"
	"slices"
	"time"

	"github.com/google/uuid"
)

var logger = logging.NewLogger("game")

type Game struct {
	Id     string
	User   string
	Status string
}

type Status struct {
	New        string
	InProgress string
	Cancelled  string
	Completed  string
	Failed     string
}

var status = Status{New: "new", InProgress: "in progress", Cancelled: "cancelled", Completed: "completed", Failed: "failed"}

func NewGame(user string, db *sql.DB) (Game, error) {
	game := Game{Id: uuid.NewString(), User: user, Status: status.New}
	_, err := db.Exec("INSERT INTO game (id, user, status, created) VALUES (?, ?, ?, ?)", game.Id, game.User, game.Status, time.Now())
	if err != nil {
		return Game{}, err
	}

	return game, nil
}

func (g *Game) InProgress(db *sql.DB) error {
	g.Status = status.InProgress
	_, err := db.Exec("UPDATE game SET status = ?, updated = ? WHERE id = ?", g.Status, time.Now(), g.Id)
	return err
}

func (g *Game) Complete(db *sql.DB) error {
	if !g.StatusIsFinal() {
		g.Status = status.Completed
		_, err := db.Exec("UPDATE game SET status = ?, updated = ? WHERE id = ?", g.Status, time.Now(), g.Id)
		return err
	}
	return nil
}

func (g *Game) Cancel(db *sql.DB) error {
	if !g.StatusIsFinal() {
		g.Status = status.Cancelled
		_, err := db.Exec("UPDATE game SET status = ?, updated = ? WHERE id = ?", g.Status, time.Now(), g.Id)
		return err
	}

	return nil
}

func (g *Game) Fail(db *sql.DB) error {
	if !g.StatusIsFinal() {
		g.Status = status.Failed
		_, err := db.Exec("UPDATE game SET status = ?, updated = ? WHERE id = ?", g.Status, time.Now(), g.Id)
		return err
	}

	return nil
}

func (g *Game) StatusIsFinal() bool {
	return slices.Contains([]string{status.Cancelled, status.Completed, status.Failed}, g.Status)
}

func CancelAllGamesForUser(user string, db *sql.DB) error {
	_, err := db.Exec("UPDATE game SET status = ?, updated = ? WHERE user = ? AND status NOT IN (?, ?)", status.Cancelled, time.Now(), user, status.Completed, status.Cancelled)
	return err
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
	UnwantedWords   []string
	WrongPositions  map[int][]rune
	AmountOfLetters map[rune]int
}

func FilterWords(words []string, guesses []guess.Guess) (filteredWords []string, additionalResults FWAdditionalResults, error error) {
	logger := logging.NewLoggerFromParent("FilterWords", &logger)

	letterPositions := []rune{'_', '_', '_', '_', '_'}
	unwantedLetters := []rune{}
	unwantedWords := []string{}
	wrongPositions := map[int][]rune{}
	amountOfLetters := map[rune]int{}

	for _, guess := range guesses {
		logger.Debugf("Current guess %v - %v", guess.Word, guess.Result)
		localAmountOfLetters := map[rune]int{}

		if guess.Result != "22222" {
			unwantedWords = append(unwantedWords, guess.Word)
		}

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
				localAmountOfLetters[currentLetter] += 1
			default:
				log.Fatal("Unexpected result rune: ", string(currentResult))
			}
		}

		for letter, localAmount := range localAmountOfLetters {
			amount, ok := amountOfLetters[letter]
			if (ok && amount < localAmount) || !ok {
				amountOfLetters[letter] = localAmount
			}
		}

		for _, letter := range unwantedLetters {
			if _, ok := amountOfLetters[letter]; ok || slices.Contains(letterPositions, letter) {
				logger.Debugf("Letter %q is in amountOfLetters or letterPositions so it must be deleted from unwantedLetters %q", letter, unwantedLetters)
				// Need to recalculate this, because after unwantedLetters shrink, index of loop is no more valid letter position in unwantedLetters
				currentLetterIndex := slices.Index(unwantedLetters, letter)
				if currentLetterIndex < len(unwantedLetters) {
					unwantedLetters[currentLetterIndex] = unwantedLetters[len(unwantedLetters)-1]
				}

				unwantedLetters = unwantedLetters[:len(unwantedLetters)-1]

				logger.Debugf("unwantedLetters after deletion of letter %q: %q", letter, unwantedLetters)
			}
		}
	}

	newWords := []string{}
	for _, word := range words {
		if len(word) == 0 {
			continue
		}

		if wordsUtils.WordRemains(word, unwantedLetters, unwantedWords, letterPositions, amountOfLetters, wrongPositions) {
			logger.Debugf("Word remains %v", word)
			newWords = append(newWords, word)
		}
	}

	slices.Sort(unwantedLetters)

	return newWords, FWAdditionalResults{LetterPositions: letterPositions, UnwantedLetters: unwantedLetters, UnwantedWords: unwantedWords, WrongPositions: wrongPositions, AmountOfLetters: amountOfLetters}, nil
}
