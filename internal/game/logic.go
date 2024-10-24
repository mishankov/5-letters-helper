package game

import (
	"slices"
)

// GetWordStatus returns how guess compares to target word
func GetWordStatus(guess, target string) string {
	status := []rune{'_', '_', '_', '_', '_'}
	takenLetters := []bool{false, false, false, false, false}

	for i, letter := range []rune(guess) {
		if letter == []rune(target)[i] {
			status[i] = '2'
			takenLetters[i] = true
		}

		if !slices.Contains([]rune(target), letter) {
			status[i] = '0'
		}
	}

	for i, letter := range []rune(guess) {
		if status[i] != '_' {
			continue
		}

		letterIndex := -1

		for j, targetLetter := range []rune(target) {
			if !takenLetters[j] && targetLetter == letter {
				letterIndex = j
				break
			}
		}

		if letterIndex > -1 {
			status[i] = '1'
			takenLetters[letterIndex] = true
		} else {
			status[i] = '0'
		}
	}

	return string(status)
}
