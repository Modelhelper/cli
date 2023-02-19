package codegen

import (
	"testing"
)

func TestTrainCase(t *testing.T) {

	tests := []struct {
		input string
		want  string
	}{
		{"ThisIsTheCase", "This_Is_The_Case"},
		{"thisIsTheCase", "This_Is_The_Case"},
		{"this-is-the-case", "This_Is_The_Case"},
		{"this_is_the_case", "This_Is_The_Case"},
		{"this is the case", "This_Is_The_Case"},
	}

	for _, testCase := range tests {

		actual := TrainCase(testCase.input)
		if actual != testCase.want {
			t.Errorf("Expected %s but got %s", testCase.want, actual)
		}
	}
}
func TestKebabCase(t *testing.T) {

	tests := []struct {
		input string
		want  string
	}{
		{"ThisIsTheCase", "this-is-the-case"},
		{"thisIsTheCase", "this-is-the-case"},
		{"this-is-the-case", "this-is-the-case"},
		{"this_is_the_case", "this-is-the-case"},
		{"this is the case", "this-is-the-case"},
	}

	for _, testCase := range tests {

		actual := kebabCase(testCase.input)
		if actual != testCase.want {
			t.Errorf("Expected %s but got %s", testCase.want, actual)
		}
	}

}
func TestPascalCase(t *testing.T) {

	tests := []struct {
		input string
		want  string
	}{
		{"ThisIsTheCase", "ThisIsTheCase"},
		{"thisIsTheCase", "ThisIsTheCase"},
		{"this-is-the-case", "ThisIsTheCase"},
		{"this_is_the_case", "ThisIsTheCase"},
		{"this is the case", "ThisIsTheCase"},
	}

	for _, testCase := range tests {

		actual := pascalCase(testCase.input)
		if actual != testCase.want {
			t.Errorf("Expected %s but got %s", testCase.want, actual)
		}
	}

}
func TestCamelCase(t *testing.T) {

	tests := []struct {
		input string
		want  string
	}{
		{"ThisIsTheCase", "thisIsTheCase"},
		{"thisIsTheCase", "thisIsTheCase"},
		{"this-is-the-case", "thisIsTheCase"},
		{"this_is_the_case", "thisIsTheCase"},
		{"this is the case", "thisIsTheCase"},
	}

	for _, testCase := range tests {

		actual := camelCase(testCase.input)
		if actual != testCase.want {
			t.Errorf("Expected %s but got %s", testCase.want, actual)
		}
	}

}
func TestSnakeCase(t *testing.T) {

	tests := []struct {
		input string
		want  string
	}{
		{"ThisIsTheCase", "this_is_the_case"},
		{"thisIsTheCase", "this_is_the_case"},
		{"this-is-the-case", "this_is_the_case"},
		{"this_is_the_case", "this_is_the_case"},
		{"this is the case", "this_is_the_case"},
	}

	for _, testCase := range tests {

		actual := snakeCase(testCase.input)
		if actual != testCase.want {
			t.Errorf("Expected %s but got %s", testCase.want, actual)
		}
	}

}
func TestDotCase(t *testing.T) {

	tests := []struct {
		input string
		want  string
	}{
		{"ThisIsTheCase", "This.Is.The.Case"},
		{"thisIsTheCase", "This.Is.The.Case"},
		{"this-is-the-case", "This.Is.The.Case"},
		{"this_is_the_case", "This.Is.The.Case"},
		{"this is the case", "This.Is.The.Case"},
	}

	for _, testCase := range tests {

		actual := DotCase(testCase.input)
		if actual != testCase.want {
			t.Errorf("Expected %s but got %s", testCase.want, actual)
		}
	}

}
func TestSentence(t *testing.T) {

	tests := []struct {
		input string
		want  string
	}{
		{"ThisIsTheCase", "This is the case"},
		{"thisIsTheCase", "This is the case"},
		{"This Is The Case", "This is the case"},
		{"this is the case", "This is the case"},
		{"this-is-the-case", "This is the case"},
		{"this_is_the_case", "This is the case"},
	}

	for _, testCase := range tests {

		actual := asSentence(testCase.input)
		if actual != testCase.want {
			t.Errorf("Expected %s but got %s", testCase.want, actual)
		}
	}

}
func TestTitleCasing(t *testing.T) {

	tests := []struct {
		input string
		want  string
	}{
		{"ThisIsTheCase", "This Is The Case"},
		{"This is the case", "This Is The Case"},
		{"This-is-the-case", "This Is The Case"},
		{"ThisIsTheCase", "This Is The Case"},
		{"TitleOfManAndMice", "Title Of Man And Mice"},
	}

	for _, testCase := range tests {

		actual := titleCase(testCase.input)
		if actual != testCase.want {
			t.Errorf("Expected %s but got %s", testCase.want, actual)
		}
	}

}
func TestCapitalCasing(t *testing.T) {

	tests := []struct {
		input string
		want  string
	}{
		{"ThisIsTheCase", "This Is The Case"},
		{"This is the case", "This Is The Case"},
		{"This-is-the-case", "This Is The Case"},
	}

	for _, testCase := range tests {

		actual := Captial(testCase.input)
		if actual != testCase.want {
			t.Errorf("Expected %s but got %s", testCase.want, actual)
		}
	}

}
func TestMacroCase(t *testing.T) {

	tests := []struct {
		input string
		want  string
	}{
		{"ThisIsTheCase", "THIS_IS_THE_CASE"},
	}

	for _, testCase := range tests {

		actual := macroCase(testCase.input)
		if actual != testCase.want {
			t.Errorf("Expected %s but got %s", testCase.want, actual)
		}
	}

}
func TestUpperCaseSplitter(t *testing.T) {

	type tc struct {
		word string
		exp  []string
	}

	testTable := []tc{
		{"PascalAPIController", []string{"Pascal", "API", "Controller"}},
		{"PascalHTTPControllerAPI", []string{"Pascal", "HTTP", "Controller", "API"}},
		{"PascalCase", []string{"Pascal", "Case"}},
		{"API", []string{"API"}},
	}

	for _, testCase := range testTable {
		exp := testCase.exp
		actual := splitOnCasing(testCase.word)

		for idx, ex := range exp {
			a := actual[idx]
			if a != ex {
				t.Errorf("Expected %s but got %s", ex, a)
			}
		}
	}
}
func TestSplitter(t *testing.T) {

	type tc struct {
		word string
		exp  []string
	}

	testTable := []tc{
		{"PascalCase", []string{"PascalCase"}},
		{"camelCase", []string{"camelCase"}},
		{"Kebab-case", []string{"Kebab", "case"}},
		{"snake_case", []string{"snake", "case"}},
		{"space_case", []string{"space", "case"}},
	}

	for _, testCase := range testTable {
		exp := testCase.exp
		actual := splitOnSplitter(testCase.word)

		for idx, ex := range exp {
			a := actual[idx]
			if a != ex {
				t.Errorf("Expected %s but got %s", ex, a)
			}
		}
	}
}
func TestFullSplitter(t *testing.T) {

	type tc struct {
		word string
		exp  []string
	}

	testTable := []tc{
		{"PascalCase", []string{"Pascal", "Case"}},
		{"camelCase", []string{"camel", "Case"}},
		{"Kebab-case", []string{"Kebab", "case"}},
		{"snake_case", []string{"snake", "case"}},
		{"space_case", []string{"space", "case"}},
		{"PascalAPIController", []string{"Pascal", "API", "Controller"}},
		{"PascalHTTPControllerAPI", []string{"Pascal", "HTTP", "Controller", "API"}},
		{"PascalCase", []string{"Pascal", "Case"}},
		{"API", []string{"API"}},
	}

	for _, testCase := range testTable {
		exp := testCase.exp
		actual := asWordArray(testCase.word)

		for idx, ex := range exp {
			a := actual[idx]
			if a != ex {
				t.Errorf("Expected %s but got %s", ex, a)
			}
		}
	}
}
