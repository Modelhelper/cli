package parser

type Parser interface {
	Parse(input interface{}) (string, error)
}
