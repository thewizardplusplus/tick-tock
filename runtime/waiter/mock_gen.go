package waiter

import (
	syncutils "github.com/thewizardplusplus/go-sync-utils"
)

// Waiter ...
//go:generate mockery -name=Waiter -case=underscore
type Waiter interface {
	syncutils.WaitGroup
}
