package expressions

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// BinaryArithmeticOperationHandler ...
type BinaryArithmeticOperationHandler func(float64, float64) float64

// BinaryArithmeticOperation ...
type BinaryArithmeticOperation struct {
	leftOperand  Expression
	rightOperand Expression
	handler      BinaryArithmeticOperationHandler
}

// NewBinaryArithmeticOperation ...
func NewBinaryArithmeticOperation(
	leftOperand Expression,
	rightOperand Expression,
	handler BinaryArithmeticOperationHandler,
) BinaryArithmeticOperation {
	return BinaryArithmeticOperation{leftOperand, rightOperand, handler}
}

// Evaluate ...
func (expression BinaryArithmeticOperation) Evaluate(
	context context.Context,
) (result interface{}, err error) {
	leftValue, err := evaluateFloat64Operand(context, expression.leftOperand)
	if err != nil {
		return nil, errors.Wrap(err, "unable to evaluate the left operand")
	}

	rightValue, err := evaluateFloat64Operand(context, expression.rightOperand)
	if err != nil {
		return nil, errors.Wrap(err, "unable to evaluate the right operand")
	}

	return expression.handler(leftValue, rightValue), nil
}
