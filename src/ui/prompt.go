package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

func PromptForString(question string) (answer string) {

	validate := func(input string) error {
		return nil
	}
	// reader := bufio.NewReader(os.Stdin)

	// fmt.Print(question)
	// text, _ := reader.ReadString('\n')
	// text = strings.ReplaceAll(text, "\r\n", "")

	// return text
	prompt := promptui.Prompt{
		Label:    question,
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		return ""
	}

	return result
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
func PromptForPassword(question string) string {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:       "Password",
		Validate:    validate,
		Mask:        '*',
		HideEntered: true,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return result
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

func PromptForEditor(question string) (answer string) {
	return question
}

func PromptForLanguage(question string) (answer string) {
	lang := make(map[string]string)
	lang["C#"] = "cs"
	lang["Go"] = "go"
	lang["TypeScript"] = "ts"
	lang["JavaScript"] = "js"
	lang["Python"] = "py"
	lang["Java"] = "java"

	items := []string{}

	for k, _ := range lang {
		items = append(items, k)
	}

	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label:    question,
			Items:    items,
			AddLabel: "Other",
		}

		index, result, err = prompt.Run()

		if index == -1 {
			items = append(items, result)
		}
	}

	if err != nil {
		// fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	selected, found := lang[result]
	if found {
		return selected
	}

	return result
}

func PromptForYesNoList(question string) bool {
	answers := make(map[string]bool)
	answers["Yes"] = true
	answers["No"] = false

	items := []string{}

	for k, _ := range answers {
		items = append(items, k)
	}

	index := -1
	var key string
	var err error

	for index < 0 {
		prompt := promptui.Select{
			Label: question,
			Items: items,
		}

		index, key, err = prompt.Run()

		if index == -1 {
			items = append(items, key)
		}
	}

	if err != nil {
		// fmt.Printf("Prompt failed %v\n", err)
		return false
	}

	selected, found := answers[key]
	if found {
		return selected
	}

	return false
}

func ClearScreen() {

	cmd := "cls"
	// switch runtime.GOOS {
	// case "winddows":
	// }

	runner := exec.Command(cmd)
	runner.Stdout = os.Stdout
	runner.Run()

}
