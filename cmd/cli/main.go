package main

import (
	"fiveLettersHelper/internal/db"
	"fiveLettersHelper/internal/game"
	"fiveLettersHelper/internal/user"
	"fiveLettersHelper/internal/words"
	"fmt"
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

	words, err := words.GetFiveLettersWords()

	if err != nil {
		log.Fatal("Error getting words from file:", err)
	}

	turnNumber := 0
	game.InProgress()
	for {
		turnNumber++
		fmt.Printf("Ход №: %v\n", turnNumber)
		fmt.Printf("Осталось %v слов для выбора\n", len(words))
		break
	}

	game.Cancel()
}
