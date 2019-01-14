package main

import (
	"os"
	"sync"

	"github.com/spf13/afero"
	"github.com/thewizardplusplus/tick-tock/internal/options"
	"github.com/thewizardplusplus/tick-tock/interpreter"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/translator"
)

func main() {
	errorHandler := runtime.NewDefaultErrorHandler(os.Stderr, os.Exit)
	options, err := options.Parse(os.Args, options.Dependencies{
		UsageWriter: os.Stdout,
		ErrorWriter: os.Stderr,
		Exiter:      os.Exit,
	})
	if err != nil {
		errorHandler.HandleError(err)
	}

	var waiter sync.WaitGroup
	if err := interpreter.Interpret(new(context.DefaultContext), options, interpreter.Dependencies{
		Reader: interpreter.ReaderDependencies{DefaultReader: os.Stdin, FileSystem: afero.NewOsFs()},
		Translator: translator.Dependencies{
			OutWriter: os.Stdout,
			Runtime:   runtime.Dependencies{Waiter: &waiter, ErrorHandler: errorHandler},
		},
	}); err != nil {
		errorHandler.HandleError(err)
	}

	waiter.Wait()
}
