package main

import (
	"fmt"

	"github.com/mishankov/go-utlz/cliutils"
)

func greetingMessage(username string) string {
	return fmt.Sprintf(`Привет, %v! Это бот-помошник для игры "5 букв" или любого другого русскоязычного аналога игры Wordle
Для запуска новой игры используй команду %v. Во время игры введи сначала слово из 5-ти букв и следующим сообщением результат этого слова в игре используя цифры 0, 1 и 2, где:
0 - буквы нет в слове
1 - буква есть в слове, но не на том месте
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

func newRoundInfo(roundNumber int, remainingWordsAmount int, wordsToShow []string) string {
	var wordsString string

	if remainingWordsAmount >= len(wordsToShow) {
		wordsString = fmt.Sprintf("Осталось %v слов. Первые %v из них: %v", remainingWordsAmount, len(wordsToShow), cliutils.FormatListWithSeparator(wordsToShow, ", "))
	} else {
		wordsString = fmt.Sprintf("Осталось %v слов: %v", remainingWordsAmount, cliutils.FormatListWithSeparator(wordsToShow, ", "))
	}

	return fmt.Sprintf("Ход №%v\n%v", roundNumber, wordsString)
}

func askForWord() string {
	return "Введи слово (5 букв):"
}

func invalidWord(word string) string {
	return fmt.Sprintf("В слове должно быть 5 букв. В слове %q их %v. Попробуй ещё раз", word, len([]rune(word)))
}

func askForResult() string {
	return "Введи результат (0, 1, 2):"
}

func invalidResultLen(result string) string {
	return fmt.Sprintf("В результате должно быть 5 символов. В результате %q их %v. Попробуй ещё раз", result, len([]rune(result)))
}

func invalidResultContent(position int, symbol rune) string {
	return fmt.Sprintf("В результате должно быть только символы 0, 1 и 2. На %v позиции находится символ %q. Попробуй ещё раз\n", position, symbol)
}

func gameCompleted(word string) string {
	return fmt.Sprintf("Игра закончена! Загаданное слово: %v", word)
}

func gameFailed() string {
	return "Игра закончена! Не удалось найти загаданное слово"
}
