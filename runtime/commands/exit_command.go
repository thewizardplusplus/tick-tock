package commands

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime"
)

// ExitCommand ...
type ExitCommand struct{}

// Run ...
func (command ExitCommand) Run(context runtime.Context) error {
	panic(errors.New("user exit"))
}
