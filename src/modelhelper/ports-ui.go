package modelhelper

type UI struct {
	Prompts          UIPromptService
	SingleLineString UIPrompter[string]
}

type UIPromptService interface {
	SingleLineString(input string) string
	MulitiLineString(input string) string
	YesNo(input, defaultValue string) bool
	TrueFalse(input string) bool
	Secret(input string) string
}

type UIPrompter[TReturn any] interface {
	Prompt(input string) TReturn
}
