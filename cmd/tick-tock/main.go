package main

import (
	"os"
	"sync"

	"github.com/spf13/afero"
	"github.com/thewizardplusplus/tick-tock/interpreter"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/translator"
)

func main() {
	options := parseOptions()
	waiter := new(sync.WaitGroup)
	errorHandler := runtime.NewDefaultErrorHandler(os.Stderr, os.Exit)
	dependencies := interpreter.Dependencies{
		Dependencies: translator.Dependencies{
			Dependencies: runtime.Dependencies{Waiter: waiter, ErrorHandler: errorHandler},
			OutWriter:    os.Stdout,
		},
		ReaderDependencies: interpreter.ReaderDependencies{
			DefaultReader: os.Stdin,
			FileSystem:    afero.NewOsFs(),
		},
	}
	if err := interpreter.Interpret(new(context.DefaultContext), options, dependencies); err != nil {
		errorHandler.HandleError(err)
	}

	waiter.Wait()
}
