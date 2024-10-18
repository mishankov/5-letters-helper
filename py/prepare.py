if __name__ == "__main__":
    f = open("./data/russian_nouns.txt")
    lines = f.readlines()
    f.close()

    five_letters_lines = []
    for line in lines:
        if len(line.replace("\n", "")) == 5:
            five_letters_lines.append(line)

    f = open("./data/five_letters_russian_nouns.txt", "w")
    f.writelines(five_letters_lines)
    f.close()
