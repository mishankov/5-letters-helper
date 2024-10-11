package cliUtils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func UserInput(promt string) (string, error) {
	if len(promt) > 0 {
		fmt.Print(promt)
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(input), err

}

func FormatListWithSeparator[T string | rune](l []T, sep string) string {
	if len(l) == 0 {
		return ""
	}

	r := string(l[0])
	for _, el := range l[1:] {
		r += sep + string(el)
	}
	return r
}
