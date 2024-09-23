package text

import (
	"bufio"
	"strings"
	"unicode"

	"github.com/gertd/go-pluralize"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GetStat(body []byte) (chars, lines, words int) {

	s := string(body)
	return len(s), GetLines(s), GetWords(s)
}

func GetWords(input string) int {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	count := 0
	for scanner.Scan() {
		count++
	}

	return count
}
func GetLines(input string) int {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanLines)

	count := 0
	for scanner.Scan() {
		count++
	}

	return count
}

func PluralForm(input string) string {

	if strings.ToLower(input) == "data" {
		return input
	}

	pluralize := pluralize.NewClient()
	output := pluralize.Plural(input)

	return output
}
func SingularForm(input string) string {
	if strings.ToLower(input) == "data" {
		return input
	}

	pluralize := pluralize.NewClient()
	output := pluralize.Singular(input)

	return output
}

func SnakeCase(input string) string {
	snake := wordJoiner(asWordArray(input), "_")
	return strings.ToLower(snake)
}

func MacroCase(input string) string {
	snake := wordJoiner(asWordArray(input), "_")
	return strings.ToUpper(snake)
}

func TrainCase(input string) string {
	casing := wordJoiner(asWordArray(Captial(input)), "_")
	return casing
}

func DotCase(input string) string {
	casing := wordJoiner(asWordArray(Captial(input)), ".")
	return casing
}

func KebabCase(input string) string {
	kebab := wordJoiner(asWordArray(input), "-")
	return strings.ToLower(kebab)
}
func Captial(input string) string {
	words := asWordArray(input)

	for idx, word := range words {
		word = strings.ToLower(word)
		word = strings.ToUpper(word[0:1]) + word[1:]

		words[idx] = word
	}

	return wordJoiner(words, " ")
}

func UpperCase(input string) string {
	return strings.ToUpper(input)
}

func LowerCase(input string) string {
	return strings.ToLower(input)
}

func AsSentence(input string) string {
	sentence := AsWords(input)
	sentence = strings.ToUpper(sentence[0:1]) + strings.ToLower(sentence[1:])

	return sentence
}
func TitleCase(input string) string {
	w := AsWords(input)
	c := cases.Title(language.AmericanEnglish)
	return c.String(w)
}

func AsWords(input string) string {

	return wordJoiner(asWordArray(input), " ")
}

func PascalCase(input string) string {
	w := asWordArray(input)

	var sb strings.Builder

	for i, str := range w {

		c := strings.Title(str)
		if i == 0 {

		}
		sb.WriteString(c)
	}
	return sb.String()
}

func CamelCase(input string) string {
	w := asWordArray(input)

	var sb strings.Builder

	for i, str := range w {

		c := strings.Title(str)
		if i == 0 {
			c = strings.ToLower(c)
		}
		sb.WriteString(c)
	}
	return sb.String()
}

func wordJoiner(input []string, separator string) string {
	var sb strings.Builder
	l := len(input) - 1

	for i, str := range input {
		if l == i {
			separator = ""
		}

		sb.WriteString(str + separator)
	}
	return sb.String()
}

func splitOnCasing(input string) []string {
	var words []string
	var splitPos []int
	var letterMap []int

	// nextSplitPos := 0

	// wrd := strings.Split(input, " ")
	for _, c := range input {
		val := 0
		if unicode.IsUpper(c) {
			val = 1
		}

		letterMap = append(letterMap, val)

	}

	for idx, val := range letterMap {
		if idx == 0 {
			splitPos = append(splitPos, idx)
			continue
		}

		addPos := (val == 1 && letterMap[idx-1] == 0)

		if val == 1 && idx+1 < len(letterMap) && letterMap[idx+1] == 0 {
			addPos = true
		}

		if addPos {
			splitPos = append(splitPos, idx)
		}
	}

	for idx, start := range splitPos {
		end := len(input)
		if len(splitPos) > idx+1 {
			end = splitPos[idx+1]
		}
		words = append(words, input[start:end])
	}

	return words
}

func splitOnSplitter(input string) []string {

	words := strings.FieldsFunc(input, Split)

	return words
}
func Split(r rune) bool {
	return r == ' ' || r == '_' || r == '-'
}

func asWordArray(input string) []string {
	var words []string

	split := splitOnSplitter(input)

	for _, word := range split {
		caseSplit := splitOnCasing(word)

		for _, caseWord := range caseSplit {
			words = append(words, caseWord)
		}
	}

	return words
}

func AddWord(what string, input string) string {
	output := input
	if len(what) > 0 {
		output += what
	}

	return output
}

func Abbreviate(s string) string {
	abr := ""
	for i, c := range s {
		if i == 0 || unicode.IsUpper(c) {
			abr = abr + string(c)
		}
	}

	return strings.ToLower(abr)
}
