package interpreter

import (
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/translator"
)

// Options ...
type Options struct {
	Filename       string
	InboxSize      int
	InitialState   string
	InitialMessage string
}

// Dependencies ...
type Dependencies struct {
	Reader  ReaderDependencies
	Runtime runtime.Dependencies
}

// Interpret ...
func Interpret(ctx context.Context, options Options, dependencies Dependencies) error {
	code, err := readCode(options.Filename, dependencies.Reader)
	if err != nil {
		return err
	}

	program, err := parser.Parse(code)
	if err != nil {
		return err
	}

	definitions, initialFactories, err := translator.Translate(
		program,
		ctx.ValuesNames(),
		translator.Options{
			InboxSize:    options.InboxSize,
			InitialState: context.State{Name: options.InitialState},
		},
		dependencies.Runtime,
	)
	if err != nil {
		return err
	}

	context.SetValues(ctx, definitions)

	actors := runtime.NewConcurrentActorGroup(ctx)
	for _, factory := range initialFactories {
		actors.RegisterActor(factory.CreateActor())
	}
	actors.SendMessage(context.Message{Name: options.InitialMessage})

	return nil
}
