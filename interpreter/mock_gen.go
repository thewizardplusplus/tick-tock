package interpreter

import (
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

//go:generate mockery -name=Context -inpkg -case=underscore -testonly

// Context ...
//
// It's used only for mock generating.
//
type Context interface {
	context.Context
}

//go:generate mockery -name=ErrorHandler -inpkg -case=underscore -testonly

// ErrorHandler ...
//
// It's used only for mock generating.
//
type ErrorHandler interface {
	runtime.ErrorHandler
}
