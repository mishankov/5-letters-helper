package guess

import (
	"database/sql"

	"github.com/google/uuid"
)

type Guess struct {
	Id                       string
	Game                     string
	Number                   int
	Word                     string
	Result                   string
	RemainingWordsAfterGuess []string
}

func NewGuess(game string, number int, word string, result string, db *sql.DB) (Guess, error) {
	guess := Guess{Id: uuid.NewString(), Game: game, Number: number, Word: word, Result: result}
	_, err := db.Exec("INSERT INTO guess (id, game, number, word, result) VALUES (?, ?, ?, ?, ?)", guess.Id, guess.Game, guess.Number, guess.Word, guess.Result)
	if err != nil {
		return Guess{}, err
	}

	return guess, nil
}
