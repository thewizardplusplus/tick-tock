package commands

import (
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// ExitCommand ...
type ExitCommand struct{}

// Run ...
func (command ExitCommand) Run(context context.Context) error {
	return runtime.ErrUserExit
}
