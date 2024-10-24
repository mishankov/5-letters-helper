package main

import (
	"database/sql"
	"fiveLettersHelper/internal/config"
	"fiveLettersHelper/internal/dbUtils"
	"fiveLettersHelper/internal/game"
	"fiveLettersHelper/internal/telegram"
	"fiveLettersHelper/internal/telegram/bot"
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

func handleTelegramUpdate(update telegram.Update) error {
	db, err := dbUtils.GetDB()
	if err != nil {
		log.Fatal("Can't open database:", err)
	}
	defer db.Close()

	user, err := user.CreateAndGetTelegramUser(update.Message.From.Id, db)
	if err != nil {
		log.Fatal("Can't create telegram user:", err)
	}

	bot := bot.NewBot(config.BotSecret)

	switch {
	case commands.isCommand(update.Message.Text):
		return handleCommands(update, user, db, bot)
	}

	return nil
}

func handleCommands(update telegram.Update, user user.User, db *sql.DB, bot bot.Botter) error {
	switch update.Message.Text {
	case commands.start:
		bot.SendMessage(update.Message.Chat.Id, "Здравствуй странник")
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
