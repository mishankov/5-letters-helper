package main

import (
	"fiveLettersHelper/internal/db"
	"fiveLettersHelper/internal/game"
	"fiveLettersHelper/internal/user"
	"log"
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

	log.Println("Game created with id:", game.Id)
}
