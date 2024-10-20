package bot

type SendMessageRequest struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}
