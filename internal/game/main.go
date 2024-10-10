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
