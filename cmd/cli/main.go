package main

import (
	"database/sql"
	"fiveLettersHelper/internal/game"
	"fiveLettersHelper/internal/user"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./fiveLettersHelp.db")
	if err != nil {
		log.Fatal("Can't open database:", err)
	}
	defer db.Close()

	user := user.CreateAndGetCLIUser(db)
	game := game.NewGame(user.Id, db)

	log.Println(game)
}
