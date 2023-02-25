package code

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"modelhelper/cli/modelhelper"
	"modelhelper/cli/modelhelper/models"
	"modelhelper/cli/utils/funcmap"
	txt "modelhelper/cli/utils/text"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

type codeGenerator struct {
	templateService modelhelper.CodeTemplateService
	languageService modelhelper.LanguageDefinitionService
}

func NewCodeGenerator(templateService modelhelper.CodeTemplateService, langService modelhelper.LanguageDefinitionService) modelhelper.TemplateGenerator[*models.CodeTemplate] {
	return &codeGenerator{templateService, langService}
}

func (g *codeGenerator) Generate(ctx context.Context, tpl *models.CodeTemplate, model interface{}) (*models.TemplateGeneratorResult, error) {
	start := time.Now()

	res := models.TemplateGeneratorResult{}

	template := g.fromFiles(tpl)
	buf := new(bytes.Buffer)
	err := template.ExecuteTemplate(buf, tpl.Name, model)
	if err != nil {
		return nil, err
	}

	res.Body = buf.Bytes()
	if len(res.Body) > 0 {
		c, l, w := txt.GetStat(res.Body)
		res.Statistics.Chars = c
		res.Statistics.Lines = l
		res.Statistics.Words = w
		res.Statistics.Duration = time.Since(start)
	}

	return &res, nil

}

// func fromFiles(cv CodeContextValue) *template.Template {
func (g *codeGenerator) fromFiles(currentTemplate *models.CodeTemplate) *template.Template {
	templates := make(map[string]string)
	o := &models.CodeTemplateListOptions{FilterTypes: []string{"block"}}
	tl := g.templateService.List(o)
	for k, b := range tl {
		templates[k] = b.Body
	}
	templates[currentTemplate.Name] = currentTemplate.Body

	dir := createTempDir()
	defer os.RemoveAll(dir)

	err := writeTempFiles(dir, templates)

	if err != nil {
		return nil
	}

	lang := g.languageService.GetDefinition(currentTemplate.Language)
	dt := make(map[string]string)
	ndt := make(map[string]string)

	for dtk, dtv := range lang.DataTypes {
		dt[dtk] = dtv.NotNull
		ndt[dtk] = dtv.Nullable
	}
	pattern := filepath.Join(dir, "*")

	drivers := template.Must(template.New(currentTemplate.Name).Funcs(
		funcmap.FullFuncMap(dt, ndt)).ParseGlob(pattern))

	return drivers
}

func createTempDir() string {

	dir, err := ioutil.TempDir("", "template") //os.MkdirTemp("", "template")
	if err != nil {
		log.Fatal(err)
	}

	return dir

}

func writeTempFiles(dir string, templates map[string]string) error {
	for name, body := range templates {
		fp := filepath.Join(dir, name)
		err := ioutil.WriteFile(fp, []byte(body), 0777)
		if err != nil {
			return err
		}
	}

	return nil
}
