package guess

import (
	"database/sql"

	"github.com/google/uuid"
)

type Guess struct {
	Id     string
	Game   string
	Number int
	Word   string
	Result string
}

func NewGuess(game string, number int, word string, result string, db *sql.DB) (Guess, error) {
	guess := Guess{Id: uuid.NewString(), Game: game, Number: number, Word: word, Result: result}
	_, err := db.Exec("INSERT INTO guess (id, game, number, word, result) VALUES (?, ?, ?, ?, ?)", guess.Id, guess.Game, guess.Number, guess.Word, guess.Result)
	if err != nil {
		return Guess{}, err
	}

	return guess, nil
}

func GetGuesseForGame(game string, db *sql.DB) ([]Guess, error) {
	guesses := []Guess{}

	rows, err := db.Query("SELECT id, game, number, word, result FROM guess WHERE game = ? ORDER BY number", game)
	if err != nil {
		return []Guess{}, err
	}
	defer rows.Close()

	for rows.Next() {
		guess := Guess{}
		if err := rows.Scan(&guess.Id, &guess.Game, &guess.Number, &guess.Word, &guess.Result); err != nil {
			return []Guess{}, err
		}
		guesses = append(guesses, guess)
	}

	return guesses, nil
}
