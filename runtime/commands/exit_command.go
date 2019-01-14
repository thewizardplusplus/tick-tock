package commands

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// ExitCommand ...
type ExitCommand struct{}

// Run ...
func (command ExitCommand) Run(context context.Context) error {
	panic(errors.New("user exit"))
}
