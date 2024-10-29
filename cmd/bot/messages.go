package main

import "fmt"

func greetingMessage(username string) string {
	return fmt.Sprintf(`Привет, %v! Это бот-помошник для игры "5 букв" или любого другого рускоязычнго аналога игры Wordle 
Для запуска новой игры используй команду %v. Во время игры введи сначала слово из 5-ти букв и следующим сообщением результат этого слова в игре используя цифры 0, 1 и 2, где:
0 - буквы нет в слове
1 - буква есть в слове, но не натом месте
2 - буква есть в слове и на правильном месте`, username, commands.newGame)
}

func newGameMessage() string {
	return "Началась новая игра!"
}

func cancelGameMessage() string {
	return "Игра завершена!"
}

func errorHappendMessage() string {
	return "Что-то пошло не так! Мы уже знаем об этом и пытаемся что-то сделать"
}
