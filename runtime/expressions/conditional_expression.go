package expressions

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

// ConditionalCase ...
type ConditionalCase struct {
	Condition Expression
	Command   runtime.Command
}

// ConditionalExpression ...
type ConditionalExpression struct {
	conditionalCases []ConditionalCase
}

// NewConditionalExpression ...
func NewConditionalExpression(conditionalCases []ConditionalCase) ConditionalExpression {
	return ConditionalExpression{conditionalCases}
}

// Evaluate ...
func (expression ConditionalExpression) Evaluate(
	context context.Context,
) (result interface{}, err error) {
	for index, conditionalCase := range expression.conditionalCases {
		conditionResult, err := conditionalCase.Condition.Evaluate(context)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to evaluate the condition #%d", index)
		}

		conditionBooleanResult, err := types.NewBoolean(conditionResult)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to convert the condition #%d to boolean", index)
		}

		if conditionBooleanResult == types.True {
			commandResult, err := conditionalCase.Command.Run(context.Copy())
			if err != nil {
				return nil, errors.Wrapf(err, "unable to evaluate the command of the condition #%d", index)
			}

			return commandResult, nil
		}
	}

	return types.Nil{}, nil
}
