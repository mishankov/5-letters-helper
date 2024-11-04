package guess

import (
	"database/sql"
	"time"

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
	_, err := db.Exec("INSERT INTO guess (id, game, number, word, result, created) VALUES (?, ?, ?, ?, ?, ?)", guess.Id, guess.Game, guess.Number, guess.Word, guess.Result, time.Now())
	if err != nil {
		return Guess{}, err
	}

	return guess, nil
}

func NewEmptyGuess(game string, number int, db *sql.DB) (Guess, error) {
	guess := Guess{Id: uuid.NewString(), Game: game, Number: number}
	_, err := db.Exec("INSERT INTO guess (id, game, number, created) VALUES (?, ?, ?, ?)", guess.Id, guess.Game, guess.Number, time.Now())
	if err != nil {
		return Guess{}, err
	}

	return guess, nil
}

func GetGuessesForGame(game string, db *sql.DB) ([]Guess, error) {
	guesses := []Guess{}

	rows, err := db.Query("SELECT id, game, number, word, result FROM guess WHERE game = ? ORDER BY number", game)
	if err != nil {
		return []Guess{}, err
	}
	defer rows.Close()

	for rows.Next() {
		guess := Guess{}
		word := sql.NullString{}
		result := sql.NullString{}
		if err := rows.Scan(&guess.Id, &guess.Game, &guess.Number, &word, &result); err != nil {
			return []Guess{}, err
		}

		guess.Word = word.String
		guess.Result = result.String

		guesses = append(guesses, guess)
	}

	return guesses, nil
}

func (g *Guess) AddWord(word string, db *sql.DB) error {
	_, err := db.Exec("UPDATE guess SET word = ?, updated = ? WHERE id = ?", word, time.Now(), g.Id)
	return err
}

func (g *Guess) AddResult(result string, db *sql.DB) error {
	_, err := db.Exec("UPDATE guess SET result = ?, updated = ? WHERE id = ?", result, time.Now(), g.Id)
	return err
}
