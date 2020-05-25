package commands

import (
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
)

// ExpressionCommand ...
type ExpressionCommand struct {
	expression expressions.Expression
}

// NewExpressionCommand ...
func NewExpressionCommand(expression expressions.Expression) ExpressionCommand {
	return ExpressionCommand{expression}
}

// Run ...
func (command ExpressionCommand) Run(context context.Context) (result interface{}, err error) {
	return command.expression.Evaluate(context)
}
