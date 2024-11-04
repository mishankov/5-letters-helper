package telegram

type Chat struct {
	Id   int
	Type string
}

type User struct {
	Id         int
	First_name string
	Last_name  string
	Username   string
}

type Message struct {
	Message_id int
	From       User
	Date       int
	Chat       Chat
	Text       string
}

type Update struct {
	Update_id int
	Message   Message
}
