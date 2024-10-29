package main

import (
	"log"
)

type FakeBot struct {
}

func (fb FakeBot) SendMessage(chatId int, text string) error {
	log.Printf("Message sent to %v: %q", chatId, text)
	return nil
}

// func TestHandleCommands(t *testing.T) {
// 	log.SetFlags(0)

// 	testDB, err := dbUtils.GetTestDB()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	dbUtils.PrepareDB(testDB)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	user, err := user.CreateAndGetTelegramUser(123, testDB)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	tests := []struct {
// 		update     telegram.Update
// 		chatId     int
// 		botMessage string
// 	}{
// 		{
// 			update: telegram.Update{Message: telegram.Message{
// 				Chat: telegram.Chat{Id: 123},
// 				Text: "/start",
// 				From: telegram.User{Username: "Test user"},
// 			}},
// 			chatId:     123,
// 			botMessage: "Здравствуй странник",
// 		},
// 	}

// 	for index, test := range tests {
// 		output := caprureOutput(func() { handleCommands(test.update, user, testDB, FakeBot{}) })
// 		target := fmt.Sprintf("Message sent to %v: %q\n", test.chatId, test.botMessage)
// 		if output != target {
// 			t.Fatalf("Error for test %v. Output: %q. Target: %q", index, output, target)
// 		}
// 	}
// }
