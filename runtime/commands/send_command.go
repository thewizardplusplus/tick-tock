package commands

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

// SendCommand ...
type SendCommand struct {
	name      string
	arguments []expressions.Expression
}

// NewSendCommand ...
func NewSendCommand(name string, arguments []expressions.Expression) SendCommand {
	return SendCommand{name, arguments}
}

// Run ...
func (command SendCommand) Run(ctx context.Context) (result interface{}, err error) {
	var arguments []interface{}
	for index, argument := range command.arguments {
		result, err := argument.Evaluate(ctx)
		if err != nil {
			return nil, errors.Wrapf(
				err,
				"unable to evaluate the argument #%d for the send command",
				index,
			)
		}

		arguments = append(arguments, result)
	}

	ctx.SendMessage(context.Message{
		Name:      command.name,
		Arguments: arguments,
	})

	return types.Nil{}, nil
}
