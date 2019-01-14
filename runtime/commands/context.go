package commands

import "github.com/thewizardplusplus/tick-tock/runtime"

// Context ...
//go:generate mockery -name=Context -inpkg -case=underscore -testonly
type Context interface {
	runtime.Context
}
