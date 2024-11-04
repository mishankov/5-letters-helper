package main

import (
	"encoding/json"
	"fiveLettersHelper/internal/config"
	"fiveLettersHelper/internal/dbUtils"
	"fiveLettersHelper/internal/telegram"
	"fiveLettersHelper/pkg/logging"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var logger = logging.NewLogger("bot")

func healthcheck(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Time: %v", time.Now())
}

func handleBot(w http.ResponseWriter, req *http.Request) {
	logger.Debug("Bot request handle begin")
	body, err := io.ReadAll(req.Body)
	if err != nil {
		logger.Error(err)
	}

	logger.Debug("Bot request body:", string(body))

	var update telegram.Update

	err = json.Unmarshal(body, &update)
	if err != nil {
		logger.Error(err)
	}

	handleTelegramUpdate(update)
	logger.Debug("Bot request handle end")
}

func handleGetDB(w http.ResponseWriter, req *http.Request) {
	db, err := dbUtils.GetDBFile()
	if err != nil {
		log.Fatal(err)
	}

	w.Write(db)
}

func main() {
	// logger := logging.NewLoggerFromParent("main", &logger)

	http.HandleFunc("/healthcheck", healthcheck)
	http.HandleFunc(fmt.Sprintf("/bot/%v", config.BotSecret), handleBot)
	http.HandleFunc(fmt.Sprintf("/get_db/%v", config.BotSecret), handleGetDB)

	db, err := dbUtils.GetDB()
	if err != nil {
		logger.Fatalf("Can't open database: %v", err)
	}
	defer db.Close()

	err = dbUtils.PrepareDB(db)
	if err != nil {
		logger.Fatalf("Error preparing DB: %v", err)
	}

	logger.Infof("Starting server: http://localhost%v", config.Port)
	logger.Infof("Bot token:", config.BotSecret)
	http.ListenAndServe(config.Port, nil)
}
