package codegen

// var datatypeMap map[string]code.LangDefDataType

type CodeGenerator interface {
	Generate(model interface{}) (string, error)
}
