package runtime

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

// Command ...
//go:generate mockery -name=Command -case=underscore
type Command interface {
	Run(context context.Context) (result interface{}, err error)
}

// CommandGroup ...
type CommandGroup []Command

// Run ...
func (commands CommandGroup) Run(context context.Context) (result interface{}, err error) {
	result = types.Nil{}
	for index, command := range commands {
		if result, err = command.Run(context); err != nil {
			return nil, errors.Wrapf(err, "unable to run the command #%d", index)
		}
	}

	return result, nil
}
