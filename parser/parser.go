package parser

import (
	"github.com/alecthomas/participle"
	"github.com/pkg/errors"
)

// ParseToAST parses a code string to a structure representing an AST.
func ParseToAST(code string, ast interface{}) error {
	parser, err := participle.Build(ast)
	if err != nil {
		return errors.Wrap(err, "unable to build the parser")
	}

	if err := parser.ParseString(code, ast); err != nil {
		return errors.Wrap(err, "unable to parse the code")
	}

	return nil
}
