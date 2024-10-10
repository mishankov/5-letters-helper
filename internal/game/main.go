package game

import (
	"database/sql"
	"log"

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

func (g *Game) InProgress() {
	g.Status = "in progress"
	// TODO: запись в БД
}

func (g *Game) Complete() {
	if g.Status != "cancelled" && g.Status != "completed" {
		g.Status = "completed"
		// TODO: запись в БД
		log.Printf("Game %v completed\n", g.Id)
	}

}

func (g *Game) Cancel() {
	if g.Status != "cancelled" && g.Status != "completed" {
		g.Status = "canceled"
		// TODO: запись в БД
		log.Printf("Game %v cancelled\n", g.Id)
	}
}
