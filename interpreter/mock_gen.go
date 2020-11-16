package interpreter

import (
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
