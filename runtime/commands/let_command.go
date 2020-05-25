package commands

import (
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
)

// LetCommand ...
type LetCommand struct {
	identifier string
	expression expressions.Expression
}

// NewLetCommand ...
func NewLetCommand(identifier string, expression expressions.Expression) LetCommand {
	return LetCommand{identifier, expression}
}

// Run ...
func (command LetCommand) Run(context context.Context) (result interface{}, err error) {
	result, err = command.expression.Evaluate(context)
	if err != nil {
		return nil, err
	}

	context.SetValue(command.identifier, result)
	return result, nil
}
