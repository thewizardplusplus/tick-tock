package expressions

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// ArithmeticNegation ...
type ArithmeticNegation struct {
	operand Expression
}

// NewArithmeticNegation ...
func NewArithmeticNegation(operand Expression) ArithmeticNegation {
	return ArithmeticNegation{operand}
}

// Evaluate ...
func (expression ArithmeticNegation) Evaluate(
	context context.Context,
) (result interface{}, err error) {
	value, err := evaluateFloat64Operand(context, expression.operand)
	if err != nil {
		return nil, errors.Wrap(err, "unable to evaluate the operand of arithmetic negation")
	}

	return -value, nil
}

func evaluateFloat64Operand(context context.Context, operand Expression) (float64, error) {
	result, err := operand.Evaluate(context)
	if err != nil {
		return 0, errors.Wrap(err, "unable to evaluate the operand")
	}

	value, ok := result.(float64)
	if !ok {
		return 0, errors.Errorf("incorrect type of the operand (%T instead float64)", result)
	}

	return value, nil
}
