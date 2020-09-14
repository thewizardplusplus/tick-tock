package expressions

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

// NilCoalescingOperator ...
type NilCoalescingOperator struct {
	leftOperand  Expression
	rightOperand Expression
}

// NewNilCoalescingOperator ...
func NewNilCoalescingOperator(
	leftOperand Expression,
	rightOperand Expression,
) NilCoalescingOperator {
	return NilCoalescingOperator{leftOperand, rightOperand}
}

// Evaluate ...
func (expression NilCoalescingOperator) Evaluate(
	context context.Context,
) (result interface{}, err error) {
	leftResult, err := expression.leftOperand.Evaluate(context)
	if err != nil {
		return nil, errors.Wrap(err, "unable to evaluate the left operand of the nil coalescing operator")
	}
	if leftResult != (types.Nil{}) {
		return leftResult, nil
	}

	rightResult, err := expression.rightOperand.Evaluate(context)
	if err != nil {
		return nil,
			errors.Wrap(err, "unable to evaluate the right operand of the nil coalescing operator")
	}

	return rightResult, nil
}
