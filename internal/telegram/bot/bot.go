package bot

import (
	"fiveLettersHelper/pkg/httpclient"
	"log"
	"strings"
)

type Bot struct {
	Token string
	url   string
}

func NewBot(token string) Bot {
	return Bot{Token: token, url: "https://api.telegram.org/bot" + token}
}

func (b Bot) SendMessage(chatId int, text string) error {
	req := SendMessageRequest{Chat_id: chatId, Text: text}

	resp, err := httpclient.Post(b.url+"/sendMessage", req)
	if err != nil {
		return err
	}

	if !strings.HasPrefix(resp.Status, "2") {
		log.Printf("Send message status: %v. Reponse body: %v\n", resp.Status, resp.Body)
	}

	return nil
}
