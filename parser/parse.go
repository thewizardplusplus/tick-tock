package parser

import (
	"github.com/alecthomas/participle"
	"github.com/pkg/errors"
)

// Parse ...
func Parse(code string) (*Program, error) {
	program := new(Program)
	if err := parseToAST(code, program); err != nil {
		return nil, err
	}

	return program, nil
}

func parseToAST(code string, ast interface{}) error {
	parser, err := participle.Build(ast)
	if err != nil {
		return errors.Wrap(err, "unable to build the parser")
	}

	if err := parser.ParseString(code, ast); err != nil {
		return errors.Wrap(err, "unable to parse the code")
	}

	return nil
}
