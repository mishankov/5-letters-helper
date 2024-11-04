package main

import (
	"database/sql"
	"fiveLettersHelper/internal/config"
	"fiveLettersHelper/internal/dbUtils"
	"fiveLettersHelper/internal/game"
	"fiveLettersHelper/internal/guess"
	"fiveLettersHelper/internal/telegram"
	"fiveLettersHelper/internal/telegram/bot"
	"fiveLettersHelper/internal/user"
	"fiveLettersHelper/internal/words"
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
		err := handleCommands(update, user, db, bot)
		if err != nil {
			bot.SendMessage(update.Message.Chat.Id, errorHappendMessage())
		}
	default:
		err := handleCurrentGameState(update, user, db, bot)
		if err != nil {
			bot.SendMessage(update.Message.Chat.Id, errorHappendMessage())
		}
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
	logger.Debug("handleCurrentGameState begin")
	latestGame, err := game.GetLatestGameGorUser(user.Id, db)
	if err != nil {
		logger.Error("Error getting latest game for user:", err)
		return err
	}

	if latestGame.Id == "" {
		logger.Warn("User did not started game manualy. Creating game")
		latestGame, err = game.NewGame(user.Id, db)
		if err != nil {
			logger.Error("Error creating new game for user:", err)
			return err
		}
	}
	logger.Debug("Latest game id:", latestGame.Id)

	guesses, err := latestGame.GetGuesses(db)
	if err != nil {
		logger.Error("Error getting latest game guesses:", err)
		return err
	}

	if len(guesses) == 0 {
		logger.Warn("No guesses founded for game")
		guess, err := guess.NewEmptyGuess(latestGame.Id, 1, db)
		if err != nil {
			logger.Error("Error creating new guess:", err)
			return err
		}

		allWords, err := words.GetFiveLettersWords()
		if err != nil {
			logger.Error("Error getting words:", err)
			return err
		}

		wordsAmount := len(allWords)
		wordsCutted := words.GetFirstNWords(words.RankWords(allWords, 1), 10)

		bot.SendMessage(update.Message.Chat.Id, newRoundInfo(guess.Number, wordsAmount, wordsCutted))
		bot.SendMessage(update.Message.Chat.Id, askForWord())

		return nil
	}
	lastGuess := guesses[len(guesses)-1]

	if lastGuess.Word == "" {
		word := update.Message.Text

		if len([]rune(word)) != 5 {
			bot.SendMessage(update.Message.Chat.Id, invalidWord(word))
			bot.SendMessage(update.Message.Chat.Id, askForWord())

			return nil
		} else {
			err := lastGuess.AddWord(word, db)
			if err != nil {
				logger.Error("Error adding word to guess:", err)
				return err
			}

			err = latestGame.InProgress(db)
			if err != nil {
				logger.Error("Error inprogressing game:", err)
				return err
			}

			bot.SendMessage(update.Message.Chat.Id, askForResult())

			return nil
		}
	}

	if lastGuess.Result == "" {
		result := update.Message.Text

		if len([]rune(result)) != 5 {
			bot.SendMessage(update.Message.Chat.Id, invalidResultLen(result))
			bot.SendMessage(update.Message.Chat.Id, askForResult())
			return nil
		}

		for i, r := range result {
			if !slices.Contains([]rune{'0', '1', '2'}, r) {
				bot.SendMessage(update.Message.Chat.Id, invalidResultContent(i+1, r))
				bot.SendMessage(update.Message.Chat.Id, askForResult())
				return nil
			}
		}

		err := lastGuess.AddResult(result, db)
		if err != nil {
			logger.Error("Error adding resilt to guess:", err)
		}

		guesses, err := latestGame.GetGuesses(db)
		if err != nil {
			logger.Error("Error getting guesses for game:", err)
			return err
		}

		allWords, err := words.GetFiveLettersWords()
		if err != nil {
			logger.Error("Error getting words:", err)
			return err
		}

		filteredWords, _, err := game.FilterWords(allWords, guesses)
		if err != nil {
			logger.Error("Error filtering words:", err)
			return err
		}

		wordsAmount := len(filteredWords)
		if wordsAmount == 1 {
			err := latestGame.Complete(db)
			if err != nil {
				logger.Error("Error completing game:", err)
				return err
			}

			bot.SendMessage(update.Message.Chat.Id, gameCompleted(filteredWords[0]))

			return nil
		}

		if wordsAmount == 0 {
			err := latestGame.Fail(db)
			if err != nil {
				logger.Error("Error failing game:", err)
				return err
			}

			bot.SendMessage(update.Message.Chat.Id, gameFailed())

			return nil
		}

		filteredWords = words.GetFirstNWords(words.RankWords(filteredWords, 1), 10)

		guess, err := guess.NewEmptyGuess(latestGame.Id, 1, db)
		if err != nil {
			logger.Error("Error creating new guess:", err)
			return err
		}

		bot.SendMessage(update.Message.Chat.Id, newRoundInfo(guess.Number, wordsAmount, filteredWords))
		bot.SendMessage(update.Message.Chat.Id, askForWord())

		return nil
	}

	return nil
}
