package game

import (
	"database/sql"

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

	return Game{Id: uuid.NewString(), User: user, Status: "new"}, nil
}
