package bot

import (
	"fiveLettersHelper/pkg/httpclient"
	"fiveLettersHelper/pkg/logging"
	"strings"
)

var logger = logging.NewLogger("telegram_bot")

type Botter interface {
	SendMessage(chatId int, text string) error
}

type Bot struct {
	Token string
	url   string
}

func NewBot(token string) Bot {
	return Bot{Token: token, url: "https://api.telegram.org/bot" + token}
}

func (b Bot) SendMessage(chatId int, text string) error {
	logger := logging.NewLoggerFromParent("SendMessage", &logger)

	req := SendMessageRequest{ChatId: chatId, Text: text}

	resp, err := httpclient.Post(b.url+"/sendMessage", req)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(resp.Status, "2") {
		logger.Infof("Send message status: %v. Reponse body: %q", resp.Status, resp.Body)
	}

	return nil
}
