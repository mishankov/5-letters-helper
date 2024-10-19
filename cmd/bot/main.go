package main

import (
	"encoding/json"
	"fiveLettersHelper/internal/config"
	"fiveLettersHelper/internal/telegram"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func healthcheck(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Time: %v", time.Now())
}

func handleBot(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	var update telegram.Update

	err = json.Unmarshal(body, &update)
	if err != nil {
		log.Fatal(err)
	}

	err = handleTelegramUpdate(update)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/healthcheck", healthcheck)
	http.HandleFunc(fmt.Sprintf("/bot/%v", config.BotSecret), handleBot)

	log.Printf("Starting server: http://localhost%v", config.Port)
	http.ListenAndServe(config.Port, nil)
}
