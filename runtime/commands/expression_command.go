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
func (command ExpressionCommand) Run(context context.Context) error {
	_, err := command.expression.Evaluate(context)
	return err
}
