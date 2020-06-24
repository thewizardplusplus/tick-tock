package commands

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

// SetCommand ...
type SetCommand struct {
	name      string
	arguments []expressions.Expression
}

// NewSetCommand ...
func NewSetCommand(name string, arguments []expressions.Expression) SetCommand {
	return SetCommand{name, arguments}
}

// Run ...
func (command SetCommand) Run(ctx context.Context) (result interface{}, err error) {
	var arguments []interface{}
	for index, argument := range command.arguments {
		result, err := argument.Evaluate(ctx)
		if err != nil {
			return nil, errors.Wrapf(
				err,
				"unable to evaluate the argument #%d for the set command",
				index,
			)
		}

		arguments = append(arguments, result)
	}

	if err := ctx.SetState(context.State{
		Name:      command.name,
		Arguments: arguments,
	}); err != nil {
		return nil, errors.Wrapf(err, "unable to set the state %s", command.name)
	}

	return types.Nil{}, nil
}
