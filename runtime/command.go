package runtime

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// Command ...
//go:generate mockery -name=Command -case=underscore
type Command interface {
	Run(context context.Context) error
}

// CommandGroup ...
type CommandGroup []Command

// Run ...
func (commands CommandGroup) Run(context context.Context) error {
	for index, command := range commands {
		if err := command.Run(context); err != nil {
			return errors.Wrapf(err, "unable to run the command #%d", index)
		}
	}

	return nil
}
