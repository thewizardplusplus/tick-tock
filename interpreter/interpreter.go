package interpreter

import (
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/translator"
)

// Dependencies ...
type Dependencies struct {
	translator.Dependencies
	ReaderDependencies
}

// Interpret ...
func Interpret(
	context context.Context,
	filename string,
	inboxSize int,
	dependencies Dependencies,
) error {
	code, err := readCode(filename, dependencies.ReaderDependencies)
	if err != nil {
		return err
	}

	program, err := parser.Parse(code)
	if err != nil {
		return err
	}

	actors, err := translator.Translate(program.Actors, inboxSize, dependencies.Dependencies)
	if err != nil {
		return err
	}

	actors.Start(context)
	actors.SendMessage("__initialize__")

	return nil
}
