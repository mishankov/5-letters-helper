def filter_words(unwanted_letters: list[str], letters_position: list[str], amount_of_letters: dict[str, int], wrong_positions: dict[int, list[str]]):
    def inner(word: str):
        for letter in unwanted_letters:
            if letter in word:
                return False

        for letter, position in zip(letters_position, range(0, 5)):
            if letter != "_":
                if word[position] != letter:
                    return False

        for letter, amount in amount_of_letters.items():
            if word.count(letter) < amount:
                return False

        for position, letters in wrong_positions.items():
            for letter in letters:
                if word[position] == letter:
                    return False
        
        return True

    return inner


if __name__ == "__main__":
    f = open("./data/five_letters_russian_nouns.txt")
    remaining_variants = list(map(lambda x: x.replace("\n", ""), f.readlines()))
    f.close() 

    amount_of_letters: dict[str, int] = {}
    letters_position: list[str] = ["_", "_", "_", "_", "_"]
    unwanted_letters: list[str] = []
    wrong_positions: dict[int, list[str]] = {}
    turn_number = 0
    while True:
        turn_number += 1
        print(f"Ход №: {turn_number}")
        print(f"Осталось {len(remaining_variants)} слов(а) для выбора. Первые из них: {', '.join(remaining_variants[:10])}")
        print(f"Известные положения букв: {' '.join(letters_position)}")
        print(f"Неиспользуемые буквы: {', '.join(sorted(unwanted_letters))}")
        word = input("Введи слово: ")
        result = input("Введи результат (0, 1, 2): ")

        local_amount_of_letters = {}
        for letter, letter_result, position in zip(word, result, range(0, 5)):
            match letter_result:
                case "0": 
                    if letter not in unwanted_letters: unwanted_letters.append(letter)
                case "1": 
                    if letter in local_amount_of_letters.keys():
                        local_amount_of_letters[letter] += 1
                    else:
                        local_amount_of_letters[letter] = 1

                    if position in wrong_positions.keys() and letter not in wrong_positions[position]:
                        wrong_positions[position].append(letter)
                        wrong_positions[position].sort()
                    else:
                        wrong_positions[position] = [letter]

                case "2": letters_position[position] = letter

        for letter, amount in local_amount_of_letters.items():
            if letter in amount_of_letters.keys() and amount_of_letters[letter] < amount:
                amount_of_letters[letter] += amount
            else:
                amount_of_letters[letter] = amount

        remaining_variants = list(filter(filter_words(unwanted_letters, letters_position, amount_of_letters, wrong_positions), remaining_variants))

        if len(remaining_variants) == 1:
            print(f"\nЗагаданное слово: {remaining_variants[0]}. Игра окончена!")
            break

        print("============================\n\n")
