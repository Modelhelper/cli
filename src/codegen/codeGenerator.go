package codegen

import "modelhelper/cli/modelhelper/models"

// var datatypeMap map[string]code.LangDefDataType

// type CodeGenerator interface {
// 	Generate(model interface{}) (string, error)
// }

// type CodeGen interface {
// 	Generate(ctx context.Context, m model.ModelConverter) (modelhelper.CodeGeneratorResult, error)
// }

type CodeContextValue struct {
	TemplateName  string
	Template      string
	Blocks        map[string]string
	Datatypes     map[string]string
	NullableTypes map[string]string
	Project       models.ProjectConfig

	// AlternativeNullableTypes map[string]string
	// Templates                map[string]string // is this really needed
}

// type Statistics struct {
// 	FilesExported    int
// 	TemplatesUsed    int
// 	EntitiesUsed     int
// 	SnippetsInserted int
// 	FilesCreated     int
// 	SnippetsCreated  int
// 	Chars            int
// 	Lines            int
// 	Words            int
// 	Duration         time.Duration
// 	TimeSaved        int
// }

// type Result struct {
// 	Stat    Statistics
// 	Content string
// }
