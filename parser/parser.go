package parser

import (
	"github.com/alecthomas/participle"
	"github.com/pkg/errors"
)

// Program ...
type Program struct {
	Actors []*Actor `parser:"{ @@ }"`
}

// Actor ...
type Actor struct {
	States []*State `parser:"\"actor\" { @@ } \";\""`
}

// State ...
type State struct {
	Name     string     `parser:"\"state\" @Ident"`
	Messages []*Message `parser:"{ @@ } \";\""`
}

// Message ...
type Message struct {
	Name     string     `parser:"\"message\" @Ident"`
	Commands []*Command `parser:"{ @@ } \";\""`
}

// Command ...
type Command struct {
	Send *string `parser:"\"send\" @Ident"`
	Set  *string `parser:"| \"set\" @Ident"`
	Out  *string `parser:"| \"out\" ( @String | @RawString )"`
	Exit bool    `parser:"| @\"exit\""`
}

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
