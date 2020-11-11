package translator

import (
	syncutils "github.com/thewizardplusplus/go-sync-utils"
	"github.com/thewizardplusplus/tick-tock/runtime"
)

//go:generate mockery -name=Waiter -inpkg -case=underscore -testonly

// Waiter ...
//
// It's used only for mock generating.
//
type Waiter interface {
	syncutils.WaitGroup
}

//go:generate mockery -name=ErrorHandler -inpkg -case=underscore -testonly

// ErrorHandler ...
//
// It's used only for mock generating.
//
type ErrorHandler interface {
	runtime.ErrorHandler
}
