package translator

import (
	syncutils "github.com/thewizardplusplus/go-sync-utils"
)

//go:generate mockery -name=Waiter -inpkg -case=underscore -testonly

// Waiter ...
//
// It's used only for mock generating.
//
type Waiter interface {
	syncutils.WaitGroup
}
