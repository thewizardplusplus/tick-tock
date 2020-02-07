package expressions

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// FunctionApplying ...
type FunctionApplying struct {
	arguments []Expression
	handler   FunctionHandler
}

// NewFunctionApplying ...
func NewFunctionApplying(arguments []Expression, handler FunctionHandler) FunctionApplying {
	return FunctionApplying{arguments, handler}
}

// Evaluate ...
func (expression FunctionApplying) Evaluate(
	context context.Context,
) (result interface{}, err error) {
	var arguments []interface{}
	for index, argument := range expression.arguments {
		result, err = argument.Evaluate(context)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to evaluate the argument #%d for the function", index)
		}

		arguments = append(arguments, result)
	}

	result, err = expression.handler(arguments)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to call the function handler")
	}

	return result, nil
}
