package expressions

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

// BooleanOperator ...
type BooleanOperator struct {
	leftOperand       Expression
	rightOperand      Expression
	valueForEarlyExit types.Boolean
}

// NewBooleanOperator ...
func NewBooleanOperator(
	leftOperand Expression,
	rightOperand Expression,
	valueForEarlyExit types.Boolean,
) BooleanOperator {
	return BooleanOperator{leftOperand, rightOperand, valueForEarlyExit}
}

// Evaluate ...
func (expression BooleanOperator) Evaluate(
	context context.Context,
) (result interface{}, err error) {
	leftResult, err := expression.leftOperand.Evaluate(context)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to evaluate the left operand of the boolean operator")
	}

	leftBooleanResult, err := types.NewBoolean(leftResult)
	if err != nil {
		return nil,
			errors.Wrapf(err, "unable to convert the left operand of the boolean operator to boolean")
	}
	if leftBooleanResult == expression.valueForEarlyExit {
		return leftResult, nil
	}

	rightResult, err := expression.rightOperand.Evaluate(context)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to evaluate the right operand of the boolean operator")
	}

	return rightResult, nil
}
