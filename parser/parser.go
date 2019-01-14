package parser

import (
	"github.com/alecthomas/participle"
	"github.com/pkg/errors"
)

// Program represents an actor list.
type Program struct {
	Actors []*Actor `parser:"{ @@ }"`
}

// Actor represents a description of an actor (green thread) with several states.
type Actor struct {
	States []*State `parser:"\"actor\" { @@ } \";\""`
}

// State represents a description of a state in which several messages can be received.
// Also, it keeps additional state parameters.
type State struct {
	Initial  bool       `parser:"[ @\"initial\" ]"`
	Name     string     `parser:"\"state\" @Ident"`
	Messages []*Message `parser:"{ @@ } \";\""`
}

// Message represents a message handler containing several commands. Also, it keeps a handler name.
type Message struct {
	Name     string     `parser:"\"message\" @Ident"`
	Commands []*Command `parser:"{ @@ } \";\""`
}

// Command represents one of the supported commands and keeps their arguments if necessary.
type Command struct {
	Send *string `parser:"\"send\" @Ident"`
	Set  *string `parser:"| \"set\" @Ident"`
	Out  *string `parser:"| \"out\" @String"`
	Exit bool    `parser:"| @\"exit\""`
}

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
