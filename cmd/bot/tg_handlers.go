package main

import (
	"database/sql"
	"fiveLettersHelper/internal/config"
	"fiveLettersHelper/internal/dbUtils"
	"fiveLettersHelper/internal/game"
	"fiveLettersHelper/internal/telegram"
	"fiveLettersHelper/internal/telegram/bot"
	"fiveLettersHelper/internal/user"
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
		logger.Error("Can't open database:", err)
		return err
	}
	defer db.Close()

	user, err := user.CreateAndGetTelegramUser(update.Message.From.Id, db)
	if err != nil {
		logger.Error("Can't create telegram user:", err)
		return err
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
		err := bot.SendMessage(update.Message.Chat.Id, greetingMessage(update.Message.From.Username))
		if err != nil {
			return err
		}
	case commands.newGame:
		err := game.CancelAllGamesForUser(user.Id, db)
		if err != nil {
			return err
		}

		_, err = game.NewGame(user.Id, db)
		if err != nil {
			return err
		}

		err = bot.SendMessage(update.Message.Chat.Id, newGameMessage())
		if err != nil {
			return err
		}

		err = handleCurrentGameState(update, user, db, bot)
		if err != nil {
			return err
		}
	case commands.cancelGames:
		err := game.CancelAllGamesForUser(user.Id, db)
		if err != nil {
			return err
		}

		err = bot.SendMessage(update.Message.Chat.Id, cancelGameMessage())
		if err != nil {
			return err
		}
	}

	return nil
}

func handleCurrentGameState(update telegram.Update, user user.User, db *sql.DB, bot bot.Botter) error {
	// get latest game
	// get latest game guesses

	// if no guesses, create first one, and ask for first word

	// if has guesses, check last
	// if has no word, process update for word
	// else if has no result, process update for result

	return nil
}
