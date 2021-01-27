package commands

import (
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
)

//go:generate mockery --name=Context --inpackage --case=underscore --testonly

// Context ...
//
// It's used only for mock generating.
//
type Context interface {
	context.Context
}

//go:generate mockery --name=Expression --inpackage --case=underscore --testonly

// Expression ...
//
// It's used only for mock generating.
//
type Expression interface {
	expressions.Expression
}
