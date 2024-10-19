package main

import (
	"database/sql"
	"fiveLettersHelper/internal/dbUtils"
	"fiveLettersHelper/internal/game"
	"fiveLettersHelper/internal/telegram"
	"fiveLettersHelper/internal/user"
	"log"
	"slices"
)

type Commands struct {
	start       string
	newGame     string
	cancelGames string
}

func (c Commands) isCommand(s string) bool {
	return slices.Contains([]string{c.start, c.newGame, c.cancelGames}, s)
}

var commands = Commands{start: "/start", newGame: "/newgame", cancelGames: "/cancelgames"}

func handleTelegramUpdate(u telegram.Update) error {
	db, err := dbUtils.GetDB()
	if err != nil {
		log.Fatal("Can't open database:", err)
	}
	defer db.Close()

	user, err := user.CreateAndGetTelegramUser(u.Message.From.Id, db)
	if err != nil {
		log.Fatal("Can't create telegram user:", err)
	}

	switch {
	case commands.isCommand(u.Message.Text):
		return handleCommands(u, user, db)
	}

	return nil
}

func handleCommands(u telegram.Update, user user.User, db *sql.DB) error {
	switch u.Message.Text {
	case commands.start:
		// TODO: send greetings message
	case commands.newGame:
		game.CancelAllGamesForUser(user.Id, db)
		game.NewGame(user.Id, db)
		// TODO: send new game message
	case commands.cancelGames:
		game.CancelAllGamesForUser(user.Id, db)
		// TOOD : send game cancelation message
	}

	return nil
}
