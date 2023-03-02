package project

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	tpl "text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type projectGeneratorService struct {
	config *models.Config
}

func (p *projectGeneratorService) BuildTemplateModel(options *models.ProjectTemplateCreateOptions, tpl *models.ProjectTemplate) *models.ProjectTemplateModel {

	v := coalesceString(tpl.Version, options.Version, "0.0.1")

	mdl := &models.ProjectTemplateModel{
		Name:    options.Name,
		Version: v,
	}

	return mdl
}

func (p *projectGeneratorService) GenerateRootDirectoryName(rootFolderTemplateName string, model *models.ProjectTemplateModel) (string, error) {
	b := generate(rootFolderTemplateName, model)
	return string(b), nil
}

// Generate implements modelhelper.ProjectGenerator
func (p *projectGeneratorService) Generate(ctx context.Context, tpl *models.ProjectTemplate, model *models.ProjectTemplateModel) ([]*models.SourceFile, error) {
	files := []*models.SourceFile{}

	for _, source := range tpl.Sources {
		path := filepath.Join(tpl.TemplateFilePath, source)
		files = append(files, getProjectSourceFiles(path)...)
	}

	parsedFiles := parseBody(model, files)

	return parsedFiles, nil
}

func NewProjectGeneratorService(cfg *models.Config) modelhelper.ProjectGenerator {
	return &projectGeneratorService{cfg}
}

func getProjectSourceFiles(path string) []*models.SourceFile {

	// abs, err := filepath.Abs(path)
	// if err != nil {
	// 	fmt.Printf("err %v", err)
	// }

	// fmt.Printf("input: %s\nabs: %s\n\n", path, abs)

	files := []*models.SourceFile{}

	filepath.Walk(path, func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			dirName, fname := filepath.Split(fullPath)

			rel := relativePath(path, dirName)

			b, _ := ioutil.ReadFile(fullPath)
			snippets := extractSnippetIdentifiers(b)

			file := models.SourceFile{
				DirectoryName: dirName,
				RelativePath:  rel,
				FileName:      fname,
				Content:       b,
				Snippets:      snippets,
			}

			files = append(files, &file)

		}

		return nil
	})

	return files
}

func extractSnippetIdentifiers(content []byte) []string {
	m := make(map[string]string)
	list := []string{}
	searchText := fmt.Sprintf(`%%%%(\w+)%%%%`)
	pattern := regexp.MustCompile(searchText)

	// Find the position of the search text using the regular expression
	matches := pattern.FindAllStringSubmatch(string(content), -1)
	if len(matches) < 2 {
		// fmt.Printf("Error: could not find \"%s\" in file\n", searchText)
		return list
	}

	// idxMatch := pattern.FindAllStringSubmatchIndex(string(content), -1)

	for _, match := range matches {
		id := match[1]
		// idx := idxMatch[midx][1]

		_, found := m[id]
		if !found {
			m[id] = id
			list = append(list, id)
		}
		// si := snippetIdentifier{id, idx}

	}

	return list
}

func coalesceString(name ...string) string {
	output := ""
	for _, n := range name {
		if len(n) > 0 {
			output = n
			break
		}
	}

	return output
}

func relativePath(path, file string) string {
	return strings.Replace(file, path, "", -1)
}
func parseBody(model *models.ProjectTemplateModel, files []*models.SourceFile) []*models.SourceFile {
	fout := []*models.SourceFile{}

	for _, f := range files {
		b := generate(string(f.Content), model)
		f.Content = b
		fout = append(fout, &models.SourceFile{
			DirectoryName: f.DirectoryName, RelativePath: f.RelativePath, FileName: f.FileName, Content: b})
	}

	return fout
}

func generate(tbody string, model *models.ProjectTemplateModel) []byte {
	t := tpl.New("project-source").Funcs(stringMap())
	t, err := t.Parse(tbody)
	if err != nil {
		// fmt.Printf("Err: %v\n", err)
		return []byte(tbody)
	}
	buf := new(bytes.Buffer)
	t.Execute(buf, model)

	return buf.Bytes()
}

func stringMap() tpl.FuncMap {
	return tpl.FuncMap{
		"plural":   pluralForm,
		"singular": SingularForm,
		// "datatype":  dataTypeConverter,
		"lower":    lowerCase,
		"upper":    upperCase,
		"words":    asWords,
		"sentence": asSentence,
		"snake":    snakeCase,
		"macro":    macroCase,
		"train":    TrainCase,
		"kebab":    kebabCase,
		"dot":      DotCase,
		"title":    titleCase,
		"pascal":   pascalCase,
		"camel":    camelCase,
		// "nullable":  nullableDatatype,
		// "datatypeN": dataTypeWithNullcheck,
		// "append": addWord,
	}

}

func pluralForm(input string) string {
	return input
}
func SingularForm(input string) string {

	return input
}

func snakeCase(input string) string {
	snake := wordJoiner(asWordArray(input), "_")
	return strings.ToLower(snake)
}

func macroCase(input string) string {
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

func kebabCase(input string) string {
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

func upperCase(input string) string {
	return strings.ToUpper(input)
}

func lowerCase(input string) string {
	return strings.ToLower(input)
}

func asSentence(input string) string {
	sentence := asWords(input)
	sentence = strings.ToUpper(sentence[0:1]) + strings.ToLower(sentence[1:])

	return sentence
}
func titleCase(input string) string {
	w := asWords(input)
	c := cases.Title(language.AmericanEnglish)
	return c.String(w)
}

func asWords(input string) string {

	return wordJoiner(asWordArray(input), " ")
}

func pascalCase(input string) string {
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

func camelCase(input string) string {
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
