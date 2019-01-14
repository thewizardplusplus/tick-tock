package interpreter

import (
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/translator"
)

// Options ...
type Options struct {
	translator.Options

	Filename       string
	InitialMessage string
}

// Dependencies ...
type Dependencies struct {
	translator.Dependencies
	ReaderDependencies
}

// Interpret ...
func Interpret(context context.Context, options Options, dependencies Dependencies) error {
	code, err := readCode(options.Filename, dependencies.ReaderDependencies)
	if err != nil {
		return err
	}

	program, err := parser.Parse(code)
	if err != nil {
		return err
	}

	actors, err := translator.Translate(program.Actors, options.Options, dependencies.Dependencies)
	if err != nil {
		return err
	}

	actors.Start(context)
	actors.SendMessage(options.InitialMessage)

	return nil
}
