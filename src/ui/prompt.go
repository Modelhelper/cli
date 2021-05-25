package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func PromptForString(question string) (answer string) {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print(question)
	text, _ := reader.ReadString('\n')
	text = strings.ReplaceAll(text, "\r\n", "")

	return text
}
func PromptForMultilineString(question string) (answer string) {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print(question)
	text, _ := reader.ReadString('|')
	// text = strings.ReplaceAll(text, "\r\n", "")

	return text
}
func PromptForYesNo(question string, defaultVal string) (answer bool) {
	text := PromptForString(question)

	if len(text) == 0 {
		text = defaultVal
	}
	a := strings.ToLower(text)
	if len(a) > 0 && a[0] == 'y' {
		return true
	}
	// strconv.ParseBool(text)
	return false
}

func PromptForBool(question string) (answer bool) {

	text := PromptForString(question)
	res, err := strconv.ParseBool(text)

	if err != nil {
		return false
	}

	return res

}
func PromptForInt(question string) (answer int) {

	text := PromptForString(question)
	fmt.Println(text)
	res, err := strconv.Atoi(text)

	if err != nil {
		fmt.Println(err)
		return 0
	}

	return res
}
